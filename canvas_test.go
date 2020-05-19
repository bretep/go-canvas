package canvas

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/matryer/is"
)

func testToken() string {
	tok := os.Getenv("CANVAS_TEST_TOKEN")
	if tok == "" {
		panic("no testing token")
	}
	return tok
}

func init() {
	t := testToken()
	SetToken(t)
}

var (
	mu             sync.Mutex
	testingUser    *User
	testingCourses []*Course
	testingCourse  *Course
)

func testUser() (*User, error) {
	var err error
	if testingUser == nil {
		testingUser, err = CurrentUser()
	}
	return testingUser, err
}

func testCourse() *Course {
	if testingCourse == nil {
		var err error
		testingCourse, err = GetCourse(2056049)
		if err != nil {
			panic("could not get test course: " + err.Error())
		}
	}
	return testingCourse
}

func testCourses() ([]*Course, error) {
	var err error
	if testingCourses == nil {
		c := New(testToken())
		testingCourses, err = c.Courses()
	}
	return testingCourses, err
}

func Test(t *testing.T) {
}

func TestSetHost(t *testing.T) {
	trans := defaultCanvas.client.Transport
	auth, ok := trans.(*auth)
	if !ok {
		t.Fatalf("could not set a host for this transport: %T", trans)
	}
	host := auth.host

	if err := SetHost("test.host"); err != nil {
		t.Error(err)
	}
	if auth.host != "test.host" {
		t.Error("did not set correct host")
	}
	defaultCanvas.client.Transport = http.DefaultTransport
	if err := SetHost("test1.host"); err == nil {
		t.Errorf("expected an error for setting host on %T", defaultCanvas.client.Transport)
	}
	defaultCanvas.client.Transport = auth
	auth.host = host
}

func TestAnnouncements(t *testing.T) {
	is := is.New(t)
	c := New(testToken())
	_, err := c.Announcements([]string{})
	is.True(err != nil)

	_, err = c.Announcements([]string{"course_1"})
	is.NoErr(err)
}

func TestCanvas_Err(t *testing.T) {
	for _, c := range []*Canvas{
		WithHost(testToken(), ""),
		WithHost("", DefaultHost),
	} {
		_, err := c.CurrentUser()
		if err == nil {
			t.Error("expected an error")
		}
		courses, err := c.ActiveCourses()
		if err == nil {
			t.Error("expected an error")
		}
		if courses != nil {
			t.Error("expected nil courses")
		}
	}
}

func TestUser(t *testing.T) {
	t.Skip()
	is := is.New(t)
	u, err := testUser()
	is.NoErr(err)
	settings, err := u.Settings()
	is.NoErr(err)
	is.True(len(settings) > 0)

	profile, err := u.Profile()
	is.NoErr(err)
	is.True(profile.ID != 0)
	is.True(len(profile.Name) > 0)

	subs, err := u.GradedSubmissions()
	is.NoErr(err)
	is.True(len(subs) > 0)

	colors, err := u.Colors()
	is.NoErr(err)
	var col, val string
	for col, val = range colors {
		break
	}
	color, err := u.Color(col)
	is.NoErr(err)
	is.Equal(color.HexCode, val)

	user, err := GetUser(u.ID)
	if err != nil {
		t.Error(err)
	}
	is.Equal(user.Name, u.Name) // names should be the same
	is.Equal(user.ID, u.ID)
	is.Equal(user.Email, u.Email)
	is.True(user.CreatedAt.Equal(u.CreatedAt))
}

func TestUser_Err(t *testing.T) {
	is := is.New(t)
	u, err := testUser()
	is.NoErr(err)
	colors, err := u.Colors()
	is.NoErr(err)
	defer deauthorize(u.client)()

	settings, err := u.Settings()
	is.True(err != nil)
	is.True(settings == nil)
	is.True(len(settings) == 0)

	profile, err := u.Profile()
	is.True(err != nil)
	is.True(profile == nil)

	var col string
	for col = range colors {
		break
	}
	color, err := u.Color(col)
	is.True(err != nil)
	is.True(color == nil)

	err = u.SetColor(col, "#FFFFFF")
	is.True(err != nil)
	// _, ok := err.(*AuthError)
	// is.True(ok)
}

func TestCourse_Files(t *testing.T) {
	is := is.New(t)
	c := testCourse()

	c.SetErrorHandler(func(e error, quit chan int) {
		t.Fatal(e)
		quit <- 1
	})
	is.True(c.client != nil)

	var (
		file   *File
		folder *Folder
	)
	t.Run("Course.Files", func(t *testing.T) {
		is := is.New(t)
		files := c.Files()
		is.True(files != nil)
		for file = range files {
			is.True(file.client != nil)
			is.True(file.ID != 0)
		}
	})

	t.Run("Course.Folders", func(t *testing.T) {
		is := is.New(t)
		folders := c.Folders()
		is.True(folders != nil)
		for folder = range folders {
			is.True(folder.client != nil)
			is.True(folder.ID != 0)
		}
	})
}

func TestCourseFiles_Err(t *testing.T) {
	is := is.New(t)
	c := testCourse()

	errorCount := 0
	c.SetErrorHandler(func(e error, q chan int) {
		if e == nil {
			t.Error("expected an error")
		} else {
			errorCount++
		}
		q <- 1
	})

	t.Run("Files", func(t *testing.T) {
		is := is.New(t)
		all, err := c.ListFiles()
		is.NoErr(err)
		i := 0
		files := c.Files()
		defer deauthorize(c.client)() // deauthorize after goroutines started
		for f := range files {
			is.True(f.ID != 0) // these be valid
			i++
		}
		is.True(len(all) > i) // the channel should have been stopped early
		files = c.Files()
		is.True(files != nil)
		for range files {
			panic("this code should not execute")
		}
	})

	t.Run("Folders", func(t *testing.T) {
		is := is.New(t)
		all, err := c.ListFolders()
		is.NoErr(err)
		i := 0
		folders := c.Folders()
		defer deauthorize(c.client)()
		for f := range folders {
			is.True(f.ID > 0)
			is.True(f.ID == all[i].ID)
			i++
		}
		is.True(len(all) > i)
		for range folders {
			panic("this code should not execute")
		}
	})
	is.True(errorCount >= 2)
	c.errorHandler = defaultErrorHandler
}

func TestAccount(t *testing.T) {
	is := is.New(t)
	c := New(testToken())
	_, err := c.SearchAccounts(Opt("name", "UC Berkeley"))
	is.NoErr(err)
}

func TestErrChan(t *testing.T) {
	is := is.New(t)
	courses, err := testCourses()
	is.NoErr(err)
	c := courses[1]
	files, _ := c.FilesErrChan()
	for range files {
	}
	folders, _ := c.FoldersErrChan()
	for range folders {
	}
}

func TestBookmarks(t *testing.T) {
	is := is.New(t)

	c := testCourse()
	err := CreateBookmark(&Bookmark{
		Name: "test bookmark",
		URL:  fmt.Sprintf("https://%s/courses/%d/assignments", DefaultHost, c.ID),
	})
	if err != nil {
		t.Error(err)
	}

	bks, err := Bookmarks()
	is.NoErr(err)
	for _, b := range bks {
		if b.Name != "test bookmark" {
			t.Error("got the wrong bookmark")
		}
		is.NoErr(DeleteBookmark(&b))
	}

	defer deauthorize(defaultCanvas.client)()
	err = CreateBookmark(&Bookmark{
		Name: "test bookmark",
		URL:  fmt.Sprintf("https://%s/courses/%d/assignments", DefaultHost, c.ID),
	})
	if err == nil {
		t.Error("expected an error")
	}
}

func TestErrors(t *testing.T) {
	is := is.New(t)
	e := &AuthError{
		Status: "test",
		Errors: []errorMsg{{"one"}, {"two"}},
	}
	is.Equal(e.Error(), "test: one, two")
	e = &AuthError{
		Status: "",
		Errors: []errorMsg{{"one"}, {"two"}},
	}
	is.Equal(e.Error(), "one, two")
	is.Equal(checkErrors([]errorMsg{}), "")

	err := &Error{}
	json.Unmarshal([]byte(`{"errors":{"end_date":"no"},"message":"error"}`), err)
	is.Equal(err.Error(), "error")
	err = &Error{}
	json.Unmarshal([]byte(`{"errors":{"end_date":"no"}}`), err)
	is.Equal(err.Error(), "end_date: no")
}

func deauthorize(d doer) func() {
	mu.Lock()
	defer mu.Unlock()
	warning := func() {
		fmt.Println("warning: client no deauthorized")
	}
	var cli *http.Client

	switch c := d.(type) {
	case *client:
		cli = &c.Client
	case *http.Client:
		cli = c
	default:
		return warning
	}

	au, ok := cli.Transport.(*auth)
	if !ok {
		return warning
	}
	token := au.token
	au.token = ""
	return func() {
		mu.Lock()
		au.token = token
		mu.Unlock()
	}
}