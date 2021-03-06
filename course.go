package canvas

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/harrybrwn/errs"
	"github.com/harrybrwn/go-querystring/query"
)

// Course represents a canvas course.
//
// https://canvas.instructure.com/doc/api/courses.html
type Course struct {
	ID                   int           `json:"id"`
	Name                 string        `json:"name"`
	SisCourseID          int           `json:"sis_course_id"`
	UUID                 string        `json:"uuid"`
	IntegrationID        string        `json:"integration_id"`
	SisImportID          int           `json:"sis_import_id"`
	CourseCode           string        `json:"course_code"`
	WorkflowState        string        `json:"workflow_state"`
	AccountID            int           `json:"account_id"`
	RootAccountID        int           `json:"root_account_id"`
	EnrollmentTermID     int           `json:"enrollment_term_id"`
	GradingStandardID    int           `json:"grading_standard_id"`
	GradePassbackSetting string        `json:"grade_passback_setting"`
	CreatedAt            time.Time     `json:"created_at"`
	StartAt              time.Time     `json:"start_at"`
	EndAt                time.Time     `json:"end_at"`
	Locale               string        `json:"locale"`
	Enrollments          []*Enrollment `json:"enrollments"`
	TotalStudents        int           `json:"total_students"`
	Calendar             struct {
		// ICS Download is the download link for the calendar
		ICSDownload string `json:"ics"`
	} `json:"calendar"`
	DefaultView       string `json:"default_view"`
	SyllabusBody      string `json:"syllabus_body"`
	NeedsGradingCount int    `json:"needs_grading_count"`

	Term           Term           `json:"term"`
	CourseProgress CourseProgress `json:"course_progress"`

	ApplyAssignmentGroupWeights bool `json:"apply_assignment_group_weights"`
	UserPermissions             struct {
		CreateDiscussionTopic bool `json:"create_discussion_topic"`
		CreateAnnouncement    bool `json:"create_announcement"`
	} `json:"permissions"`
	IsPublic                         bool   `json:"is_public"`
	IsPublicToAuthUsers              bool   `json:"is_public_to_auth_users"`
	PublicSyllabus                   bool   `json:"public_syllabus"`
	PublicSyllabusToAuth             bool   `json:"public_syllabus_to_auth"`
	PublicDescription                string `json:"public_description"`
	StorageQuotaMb                   int    `json:"storage_quota_mb"`
	StorageQuotaUsedMb               int    `json:"storage_quota_used_mb"`
	HideFinalGrades                  bool   `json:"hide_final_grades"`
	License                          string `json:"license"`
	AllowStudentAssignmentEdits      bool   `json:"allow_student_assignment_edits"`
	AllowWikiComments                bool   `json:"allow_wiki_comments"`
	AllowStudentForumAttachments     bool   `json:"allow_student_forum_attachments"`
	OpenEnrollment                   bool   `json:"open_enrollment"`
	SelfEnrollment                   bool   `json:"self_enrollment"`
	RestrictEnrollmentsToCourseDates bool   `json:"restrict_enrollments_to_course_dates"`
	CourseFormat                     string `json:"course_format"`
	AccessRestrictedByDate           bool   `json:"access_restricted_by_date"`
	TimeZone                         string `json:"time_zone"`
	Blueprint                        bool   `json:"blueprint"`
	BlueprintRestrictions            struct {
		Content           bool `json:"content"`
		Points            bool `json:"points"`
		DueDates          bool `json:"due_dates"`
		AvailabilityDates bool `json:"availability_dates"`
	} `json:"blueprint_restrictions"`
	BlueprintRestrictionsByObjectType struct {
		Assignment struct {
			Content bool `json:"content"`
			Points  bool `json:"points"`
		} `json:"assignment"`
		WikiPage struct {
			Content bool `json:"content"`
		} `json:"wiki_page"`
	} `json:"blueprint_restrictions_by_object_type"`

	client       doer
	errorHandler errorHandlerFunc
}

// ContextCode will return the context code for this specific course.
func (c *Course) ContextCode() string {
	return fmt.Sprintf("course_%d", c.ID)
}

// Settings gets the course settings
func (c *Course) Settings(opts ...Option) (cs *CourseSettings, err error) {
	cs = &CourseSettings{}
	return cs, getjson(c.client, cs, optEnc(opts), "/courses/%d/settings", c.ID)
}

// Permissions get the current user's permissions with respect to
// the course object.
func (c *Course) Permissions() (*Permissions, error) {
	p := &Permissions{}
	return p, getjson(c.client, p, nil, "/courses/%d/permissions", c.ID)

}

// Permissions is a canvas user permissions object
type Permissions struct {
	Read                        bool `json:"read"`
	ReadOutcomes                bool `json:"read_outcomes"`
	ReadSyllabus                bool `json:"read_syllabus"`
	ManageCanvasnetCourses      bool `json:"manage_canvasnet_courses"`
	ProvisionCatalog            bool `json:"provision_catalog"`
	CreateAccounts              bool `json:"create_accounts"`
	ManageLinks                 bool `json:"manage_links"`
	SuspendAccounts             bool `json:"suspend_accounts"`
	ManageDemos                 bool `json:"manage_demos"`
	BecomeUser                  bool `json:"become_user"`
	ImportSis                   bool `json:"import_sis"`
	ManageAccountMemberships    bool `json:"manage_account_memberships"`
	ManageAccountSettings       bool `json:"manage_account_settings"`
	ManageAlerts                bool `json:"manage_alerts"`
	ManageCatalog               bool `json:"manage_catalog"`
	ManageCourses               bool `json:"manage_courses"`
	ManageDataServices          bool `json:"manage_data_services"`
	ManageCourseVisibility      bool `json:"manage_course_visibility"`
	ManageDeveloperKeys         bool `json:"manage_developer_keys"`
	ModerateUserContent         bool `json:"moderate_user_content"`
	ManageFeatureFlags          bool `json:"manage_feature_flags"`
	ManageFrozenAssignments     bool `json:"manage_frozen_assignments"`
	ManageGlobalOutcomes        bool `json:"manage_global_outcomes"`
	ManageJobs                  bool `json:"manage_jobs"`
	ManageMasterCourses         bool `json:"manage_master_courses"`
	ManageRoleOverrides         bool `json:"manage_role_overrides"`
	ManageStorageQuotas         bool `json:"manage_storage_quotas"`
	ManageSis                   bool `json:"manage_sis"`
	ManageSiteSettings          bool `json:"manage_site_settings"`
	ManageUserLogins            bool `json:"manage_user_logins"`
	ManageUserObservers         bool `json:"manage_user_observers"`
	ReadCourseContent           bool `json:"read_course_content"`
	ReadCourseList              bool `json:"read_course_list"`
	ReadMessages                bool `json:"read_messages"`
	ResetAnyMfa                 bool `json:"reset_any_mfa"`
	UndeleteCourses             bool `json:"undelete_courses"`
	ChangeCourseState           bool `json:"change_course_state"`
	CreateCollaborations        bool `json:"create_collaborations"`
	CreateConferences           bool `json:"create_conferences"`
	CreateForum                 bool `json:"create_forum"`
	GenerateObserverPairingCode bool `json:"generate_observer_pairing_code"`
	ImportOutcomes              bool `json:"import_outcomes"`
	LtiAddEdit                  bool `json:"lti_add_edit"`
	ManageAdminUsers            bool `json:"manage_admin_users"`
	ManageAssignments           bool `json:"manage_assignments"`
	ManageCalendar              bool `json:"manage_calendar"`
	ManageContent               bool `json:"manage_content"`
	ManageFiles                 bool `json:"manage_files"`
	ManageGrades                bool `json:"manage_grades"`
	ManageGroups                bool `json:"manage_groups"`
	ManageInteractionAlerts     bool `json:"manage_interaction_alerts"`
	ManageOutcomes              bool `json:"manage_outcomes"`
	ManageSections              bool `json:"manage_sections"`
	ManageStudents              bool `json:"manage_students"`
	ManageUserNotes             bool `json:"manage_user_notes"`
	ManageRubrics               bool `json:"manage_rubrics"`
	ManageWiki                  bool `json:"manage_wiki"`
	ManageWikiCreate            bool `json:"manage_wiki_create"`
	ManageWikiDelete            bool `json:"manage_wiki_delete"`
	ManageWikiUpdate            bool `json:"manage_wiki_update"`
	ModerateForum               bool `json:"moderate_forum"`
	PostToForum                 bool `json:"post_to_forum"`
	ReadAnnouncements           bool `json:"read_announcements"`
	ReadEmailAddresses          bool `json:"read_email_addresses"`
	ReadForum                   bool `json:"read_forum"`
	ReadQuestionBanks           bool `json:"read_question_banks"`
	ReadReports                 bool `json:"read_reports"`
	ReadRoster                  bool `json:"read_roster"`
	ReadSis                     bool `json:"read_sis"`
	SelectFinalGrade            bool `json:"select_final_grade"`
	SendMessages                bool `json:"send_messages"`
	SendMessagesAll             bool `json:"send_messages_all"`

	ViewAnalytics         bool `json:"view_analytics"`
	ViewAuditTrail        bool `json:"view_audit_trail"`
	ViewAllGrades         bool `json:"view_all_grades"`
	ViewGroupPages        bool `json:"view_group_pages"`
	ViewQuizAnswerAudits  bool `json:"view_quiz_answer_audits"`
	ViewFeatureFlags      bool `json:"view_feature_flags"`
	ViewUserLogins        bool `json:"view_user_logins"`
	ViewLearningAnalytics bool `json:"view_learning_analytics"`
	ViewUnpublishedItems  bool `json:"view_unpublished_items"`
	ViewCourseChanges     bool `json:"view_course_changes"`
	ViewErrorReports      bool `json:"view_error_reports"`
	ViewGradeChanges      bool `json:"view_grade_changes"`
	ViewJobs              bool `json:"view_jobs"`
	ViewNotifications     bool `json:"view_notifications"`
	ViewStatistics        bool `json:"view_statistics"`

	ParticipateAsStudent bool `json:"participate_as_student"`
	ReadGrades           bool `json:"read_grades"`
	Update               bool `json:"update"`
	Delete               bool `json:"delete"`
	ReadAsAdmin          bool `json:"read_as_admin"`
	Manage               bool `json:"manage"`
	UseStudentView       bool `json:"use_student_view"`
	ReadRubrics          bool `json:"read_rubrics"`
	ResetContent         bool `json:"reset_content"`
	ReadPriorRoster      bool `json:"read_prior_roster"`
	CreateToolManually   bool `json:"create_tool_manually"`
}

// UpdateSettings will update a user's settings based on a given settings struct and
// will return the updated settings struct.
func (c *Course) UpdateSettings(settings *CourseSettings) (*CourseSettings, error) {
	m := make(map[string]interface{})
	raw, err := json.Marshal(settings)
	if err = errs.Pair(err, json.Unmarshal(raw, &m)); err != nil {
		return nil, err
	}
	vals := make(params)
	for k, v := range m {
		vals[k] = []string{fmt.Sprintf("%v", v)}
	}
	resp, err := put(c.client, c.id("/courses/%d/settings"), vals)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	s := CourseSettings{}
	return &s, json.NewDecoder(resp.Body).Decode(&s)
}

// CourseSettings is a json struct for a course's settings.
type CourseSettings struct {
	AllowStudentDiscussionTopics  bool `json:"allow_student_discussion_topics"`
	AllowStudentForumAttachments  bool `json:"allow_student_forum_attachments"`
	AllowStudentDiscussionEditing bool `json:"allow_student_discussion_editing"`
	GradingStandardEnabled        bool `json:"grading_standard_enabled"`
	GradingStandardID             int  `json:"grading_standard_id"`
	AllowStudentOrganizedGroups   bool `json:"allow_student_organized_groups"`
	HideFinalGrades               bool `json:"hide_final_grades"`
	HideDistributionGraphs        bool `json:"hide_distribution_graphs"`
	LockAllAnnouncements          bool `json:"lock_all_announcements"`
	UsageRightsRequired           bool `json:"usage_rights_required"`
}

// Users will get a list of users in the course
func (c *Course) Users(opts ...Option) (users []*User, err error) {
	return c.collectUsers("/courses/%d/users", opts)
}

// SearchUsers will search for a user in the course
func (c *Course) SearchUsers(term string, opts ...Option) (users []*User, err error) {
	opts = append(opts, Opt("search_term", term))
	return c.collectUsers("/courses/%d/search_users", opts)
}

// User gets a specific user.
func (c *Course) User(id int, opts ...Option) (*User, error) {
	u := &User{client: c.client}
	return u, getjson(c.client, u, optEnc(opts), "/courses/%d/users/%d", c.ID, id)
}

// Assignment will get an assignment from the course given an id.
//
// https://canvas.instructure.com/doc/api/assignments.html#method.assignments_api.index
func (c *Course) Assignment(id int, opts ...Option) (ass *Assignment, err error) {
	ass = &Assignment{client: c.client, courseCode: c.CourseCode}
	return ass, getjson(c.client, &ass, optEnc(opts), "/courses/%d/assignments/%d", c.ID, id)
}

// Assignments send the courses assignments over a channel concurrently.
//
// https://canvas.instructure.com/doc/api/assignments.html#method.assignments_api.index
func (c *Course) Assignments(opts ...Option) <-chan *Assignment {
	ch := make(assignmentChan)
	pages := c.assignmentspager(ch, opts)
	go handleErrs(pages.start(), ch, c.errorHandler)
	return ch
}

// ListAssignments will get all the course assignments and put them in a slice.
func (c *Course) ListAssignments(opts ...Option) (asses []*Assignment, err error) {
	ch := make(assignmentChan)
	pages := c.assignmentspager(ch, opts)
	errs := pages.start()
	for {
		select {
		case as := <-ch:
			asses = append(asses, as)
		case err = <-errs:
			return asses, err
		}
	}
}

// CreateAssignment will create an assignment.
func (c *Course) CreateAssignment(a Assignment, opts ...Option) (*Assignment, error) {
	q, err := query.Values(&assignmentOptions{a})
	if err != nil {
		return nil, err
	}
	resp, err := post(c.client, c.id("/courses/%d/assignments"), q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	as := &Assignment{}
	return as, json.NewDecoder(resp.Body).Decode(as)
}

// DeleteAssignment will delete an assignment
func (c *Course) DeleteAssignment(a *Assignment) (*Assignment, error) {
	return c.DeleteAssignmentByID(a.ID)
}

// DeleteAssignmentByID will delete an assignment givent only an assignment ID.
func (c *Course) DeleteAssignmentByID(id int) (*Assignment, error) {
	resp, err := delete(c.client, fmt.Sprintf("courses/%d/assignments/%d", c.ID, id), nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	a := &Assignment{}
	return a, json.NewDecoder(resp.Body).Decode(&a)
}

// EditAssignment will edit the assignment given. Returns the new edited assignment.
func (c *Course) EditAssignment(a *Assignment) (*Assignment, error) {
	opts := assignmentOptions{*a}
	q, err := query.Values(&opts)
	if err != nil {
		return nil, err
	}

	resp, err := put(c.client, fmt.Sprintf("/courses/%d/assignments/%d", c.ID, a.ID), q)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	newas := &Assignment{}
	return newas, json.NewDecoder(resp.Body).Decode(newas)
}

// GradingType is a grading type
type GradingType string

const (
	// PassFail is the grading type for pass fail assignments
	PassFail GradingType = "pass_fail"
	// Percent is the grading type for percent graded assignments
	Percent GradingType = "percent"
	// LetterGrade is the grading type for letter grade assignments
	LetterGrade GradingType = "letter_grade"
	// GPAScale is the grading type for GPA scale assignments
	GPAScale GradingType = "gpa_scale"
	// Points is the grading type for point graded assignments
	Points GradingType = "points"
	// NotGraded is the grading type for assignments that are not graded
	NotGraded GradingType = "not_graded"
)

type assignmentOptions struct {
	Assignment `url:"assignment"`
}

// Assignment is a struct holding assignment data
type Assignment struct {
	Name        string `json:"name" url:"name,omitempty"`
	Description string `json:"description" url:"description,omitempty"`
	ID          int    `json:"id" url:"-"`

	DueAt     time.Time `json:"due_at" url:"due_at,omitempty"`
	LockAt    time.Time `json:"lock_at" url:"lock_at,omitempty"`
	UnlockAt  time.Time `json:"unlock_at" url:"unlock_at,omitempty"`
	CreatedAt time.Time `json:"created_at" url:"-"`
	UpdatedAt time.Time `json:"updated_at" url:"-"`

	Overrides              []AssignmentOverride `json:"overrides" url:"assignment_overrides,brackets,omitempty"`
	OnlyVisibleToOverrides bool                 `json:"only_visible_to_overrides" url:"only_visible_to_overrides,omitempty"`
	HasOverrides           bool                 `json:"has_overrides" url:"-"`

	AssignmentGroupID              int               `json:"assignment_group_id" url:"assignment_group_id,omitempty"`
	AllowedExtensions              []string          `json:"allowed_extensions" url:"allowed_extensions,brackets,omitempty"`
	TurnitinEnabled                bool              `json:"turnitin_enabled" url:"turnitin_enabled,omitempty"`
	VericiteEnabled                bool              `json:"vericite_enabled" url:"vericite_enabled,omitempty"`
	TurnitinSettings               *TurnitinSettings `json:"turnitin_settings" url:"turnitin_settings,omitempty"`
	GradeGroupStudentsIndividually bool              `json:"grade_group_students_individually" url:"grade_group_students_individually,omitempty"`
	ExternalToolTagAttributes      interface{}       `json:"external_tool_tag_attributes" url:"external_tool_tag_attributes,omitempty"`
	PeerReviews                    bool              `json:"peer_reviews" url:"peer_reviews,omitempty"`
	AutomaticPeerReviews           bool              `json:"automatic_peer_reviews" url:"automatic_peer_reviews,omitempty"`
	GroupCategoryID                int               `json:"group_category_id" url:"group_category_id,omitempty"`
	Position                       int               `json:"position" url:"position,omitempty"`
	IntegrationID                  string            `json:"integration_id" url:"integration_id,omitempty"`
	IntegrationData                map[string]string `json:"integration_data" url:"integration_data,omitempty"`
	NotifyOfUpdate                 bool              `json:"notify_of_update,omitempty" url:"notify_of_update,omitempty"`
	PointsPossible                 float64           `json:"points_possible" url:"points_possible,omitempty"`
	SubmissionTypes                []string          `json:"submission_types" url:"submission_types,brakets,omitempty"`
	GradingType                    GradingType       `json:"grading_type" url:"grading_type,omitempty"`
	GradingStandardID              interface{}       `json:"grading_standard_id" url:"grading_standard_id,omitempty"`
	Published                      bool              `json:"published" url:"published,omitempty"`
	SisAssignmentID                string            `json:"sis_assignment_id" url:"sis_assignment_id,omitempty"`

	PeerReviewCount            int         `json:"peer_review_count" url:"-"`
	AllDates                   interface{} `json:"all_dates" url:"-"`
	CourseID                   int         `json:"course_id" url:"-"`
	HTMLURL                    string      `json:"html_url" url:"-"`
	SubmissionsDownloadURL     string      `json:"submissions_download_url" url:"-"`
	DueDateRequired            bool        `json:"due_date_required" url:"-"`
	MaxNameLength              int         `json:"max_name_length" url:"-"`
	PeerReviewsAssignAt        time.Time   `json:"peer_reviews_assign_at" url:"-"`
	IntraGroupPeerReviews      bool        `json:"intra_group_peer_reviews" url:"-"`
	NeedsGradingCount          int         `json:"needs_grading_count" url:"-"`
	NeedsGradingCountBySection []struct {
		SectionID         string `json:"section_id" url:"-"`
		NeedsGradingCount int    `json:"needs_grading_count" url:"-"`
	} `json:"needs_grading_count_by_section" url:"-"`
	PostToSis               bool             `json:"post_to_sis" url:"-"`
	HasSubmittedSubmissions bool             `json:"has_submitted_submissions" url:"-"`
	Unpublishable           bool             `json:"unpublishable" url:"-"`
	LockedForUser           bool             `json:"locked_for_user" url:"-"`
	LockInfo                *LockInfo        `json:"lock_info" url:"-"`
	LockExplanation         string           `json:"lock_explanation" url:"-"`
	QuizID                  int              `json:"quiz_id" url:"-"`
	AnonymousSubmissions    bool             `json:"anonymous_submissions" url:"-"`
	DiscussionTopic         *DiscussionTopic `json:"discussion_topic" url:"-"`
	FreezeOnCopy            bool             `json:"freeze_on_copy" url:"-"`
	Frozen                  bool             `json:"frozen" url:"-"`
	FrozenAttributes        []string         `json:"frozen_attributes" url:"-"`
	Submission              *Submission      `json:"submission" url:"-"`
	UseRubricForGrading     bool             `json:"use_rubric_for_grading" url:"-"`
	RubricSettings          interface{}      `json:"rubric_settings" url:"-"`
	Rubric                  []RubricCriteria `json:"rubric" url:"-"`
	AssignmentVisibility    []int            `json:"assignment_visibility" url:"-"`
	PostManually            bool             `json:"post_manually" url:"-"`

	OmitFromFinalGrade              bool `json:"omit_from_final_grade" url:"omit_from_final_grade,omitempty"`
	ModeratedGrading                bool `json:"moderated_grading" url:"moderated_grading,omitempty"`
	GraderCount                     int  `json:"grader_count" url:"grader_count,omitempty"`
	FinalGraderID                   int  `json:"final_grader_id" url:"final_grader_id,omitempty"`
	GraderCommentsVisibleToGraders  bool `json:"grader_comments_visible_to_graders" url:"grader_comments_visible_to_graders,omitempty"`
	GradersAnonymousToGraders       bool `json:"graders_anonymous_to_graders" url:"graders_anonymous_to_graders,omitempty"`
	GraderNamesVisibleToFinalGrader bool `json:"grader_names_visible_to_final_grader" url:"graders_names_visible_to_final_grader,omitempty"`
	AnonymousGrading                bool `json:"anonymous_grading" url:"anonymous_grading,omitempty"`
	AllowedAttempts                 int  `json:"allowed_attempts" url:"allowed_attempts,omitempty"`

	courseCode string
	client     doer
}

// SubmitFile will submit the contents of an io.Reader as
// a file to the assignment.
//
// https://canvas.instructure.com/doc/api/submissions.html#method.submissions.create
func (a *Assignment) SubmitFile(filename string, r io.Reader, opts ...Option) (*File, error) {
	if filename == "" {
		if named, ok := r.(interface{ Name() string }); ok {
			filename = named.Name()
		}
	}
	params := fileUploadParams{
		Name:        filename,
		OnDuplicate: "rename",
		ContentType: filenameContentType(filename),
	}
	params.setOptions(opts)
	if hasstat, ok := r.(interface{ Stat() (os.FileInfo, error) }); ok {
		stat, err := hasstat.Stat()
		if err == nil {
			params.Size = int(stat.Size())
		}
	}
	endpoint := fmt.Sprintf("/courses/%d/assignments/%d/submissions/self/files", a.CourseID, a.ID)
	return uploadFile(a.client, r, endpoint, &params)
}

// SubmitOsFile is the same as SubmitFile except it takes advantage of
// the extra file data stored in an *os.File.
func (a *Assignment) SubmitOsFile(f *os.File) (*File, error) {
	return a.SubmitFile(f.Name(), f)
}

// TurnitinSettings is a settings struct for turnitin
type TurnitinSettings struct {
	OriginalityReportVisibility string `json:"originality_report_visibility"`
	SPaperCheck                 bool   `json:"s_paper_check"`
	InternetCheck               bool   `json:"internet_check"`
	JournalCheck                bool   `json:"journal_check"`
	ExcludeBiblio               bool   `json:"exclude_biblio"`
	ExcludeQuoted               bool   `json:"exclude_quoted"`
	ExcludeSmallMatchesType     string `json:"exclude_small_matches_type"`
	ExcludeSmallMatchesValue    int    `json:"exclude_small_matches_value"`
}

// RubricCriteria has the rubric information for an assignment.
type RubricCriteria struct {
	Points            float64 `json:"points"`
	ID                string  `json:"id"`
	LearningOutcomeID string  `json:"learning_outcome_id"`
	VendorGUID        string  `json:"vendor_guid"`
	Description       string  `json:"description"`
	LongDescription   string  `json:"long_description"`
	CriterionUseRange bool    `json:"criterion_use_range"`
	Ratings           []struct {
		ID              string  `json:"id"`
		Description     string  `json:"description"`
		LongDescription string  `json:"long_description"`
		Points          float64 `json:"points"`
	} `json:"ratings"`
	IgnoreForScoring bool `json:"ignore_for_scoring"`
}

// LockInfo is a struct containing assignment lock status.
type LockInfo struct {
	AssetString    string    `json:"asset_string"`
	UnlockAt       time.Time `json:"unlock_at"`
	LockAt         time.Time `json:"lock_at"`
	ContextModule  string    `json:"context_module"`
	ManuallyLocked bool      `json:"manually_locked"`
}

// AssignmentOverride is an assignment override object
type AssignmentOverride struct {
	ID              int       `json:"id" url:"-"`
	Title           string    `json:"title" url:"title"`
	StudentIds      []int     `json:"student_ids" url:"student_ids,brackets,omitempty"`
	CourseSectionID int       `json:"course_section_id" url:"course_section_id"`
	DueAt           time.Time `json:"due_at" url:"due_at,omitempty"`
	UnlockAt        time.Time `json:"unlock_at" url:"unlock_at,omitempty"`
	LockAt          time.Time `json:"lock_at" url:"lock_at,omitempty"`

	AssignmentID int       `json:"assignment_id" url:"-"`
	GroupID      int       `json:"group_id" url:"-"`
	AllDay       bool      `json:"all_day" url:"-"`
	AllDayDate   time.Time `json:"all_day_date" url:"-"`
}

// DiscussionTopics return a list of the course discussion topics.
func (c *Course) DiscussionTopics(opts ...Option) ([]*DiscussionTopic, error) {
	ch := make(chan *DiscussionTopic)
	pager := newPaginatedList(
		c.client, fmt.Sprintf("/courses/%d/discussion_topics", c.ID),
		sendDiscussionTopicFunc(ch), opts,
	)
	topics := make([]*DiscussionTopic, 0)
	errs := pager.start()
	for {
		select {
		case disc := <-ch:
			topics = append(topics, disc)
		case err := <-errs:
			return topics, err
		}
	}
}

// Activity returns a course's activity data
func (c *Course) Activity() (res interface{}, err error) {
	return res, getjson(c.client, &res, nil, "/courses/%d/analytics/activity", c.ID)
}

// Files returns a channel of all the course's files
func (c *Course) Files(opts ...Option) <-chan *File {
	return filesChannel(c.client, c.id("/courses/%d/files"), c.errorHandler, opts, nil)
}

// File will get a specific file id.
func (c *Course) File(id int, opts ...Option) (*File, error) {
	f := &File{client: c.client}
	return f, getjson(c.client, f, optEnc(opts), "courses/%d/files/%d", c.ID, id)
}

// ListFiles returns a slice of files for the course.
func (c *Course) ListFiles(opts ...Option) ([]*File, error) {
	return listFiles(c.client, c.id("courses/%d/files"), nil, opts)
}

// Folders will retrieve the course's folders.
// https://canvas.instructure.com/doc/api/files.html#method.folders.list_all_folders
func (c *Course) Folders(opts ...Option) <-chan *Folder {
	ch := make(folderChan)
	pager := c.folderspager(ch, opts)
	go handleErrs(pager.start(), ch, c.errorHandler)
	return ch
}

// Folder will the a folder from the course given a folder id.
// https://canvas.instructure.com/doc/api/files.html#method.folders.show
func (c *Course) Folder(id int, opts ...Option) (*Folder, error) {
	f := &Folder{client: c.client}
	path := fmt.Sprintf("courses/%d/folders/%d", c.ID, id)
	return f, getjson(c.client, f, optEnc(opts), path)
}

// Root will get the root folder for the course.
func (c *Course) Root(opts ...Option) (*Folder, error) {
	// TODO: change this function name to RootFolder, just Root is confusing.
	f := &Folder{client: c.client}
	path := c.id("/courses/%d/folders/root")
	return f, getjson(c.client, f, optEnc(opts), path)
}

// ListFolders returns a slice of folders for the course.
// https://canvas.instructure.com/doc/api/files.html#method.folders.list_all_folders
func (c *Course) ListFolders(opts ...Option) ([]*Folder, error) {
	return listFolders(c.client, c.id("/courses/%d/folders"), nil, opts)
}

// FolderPath will split the path and return a list containing all of the folders in the path.
func (c *Course) FolderPath(pth string) ([]*Folder, error) {
	pth = path.Join(c.id("/courses/%d/folders/by_path"), pth)
	return folderList(c.client, pth)
}

// CreateFolder will create a new folder
// https://canvas.instructure.com/doc/api/files.html#method.folders.create
func (c *Course) CreateFolder(path string, opts ...Option) (*Folder, error) {
	dir, name := filepath.Split(path)
	return createFolder(c.client, dir, name, opts, "/courses/%d/folders", c.ID)
}

// UploadFile will upload a file to the course.
// https://canvas.instructure.com/doc/api/courses.html#method.courses.create_file
func (c *Course) UploadFile(filename string, r io.Reader, opts ...Option) (*File, error) {
	p := fileUploadParams{Name: filename}
	p.setOptions(opts)
	return uploadFile(c.client, r, c.id("/courses/%d/files"), &p)
}

// SetErrorHandler will set a error handling callback that is
// used to handle errors in goroutines. The default error handler
// will simply panic.
//
// The callback should accept an error and a quit channel.
// If a value is sent on the quit channel, whatever secsion of
// code is receiving the channel will end gracefully.
func (c *Course) SetErrorHandler(f errorHandlerFunc) {
	c.errorHandler = f
}

func (c *Course) setclient(d doer) {
	c.client = d
}

// Term is a school term. One school year.
type Term struct {
	ID      int
	Name    string
	StartAt time.Time `json:"start_at"`
	EndAt   time.Time `json:"end_at"`
}

// CourseProgress is the progress through a course.
type CourseProgress struct {
	RequirementCount          int       `json:"requirement_count"`
	RequirementCompletedCount int       `json:"requirement_completed_count"`
	NextRequirementURL        string    `json:"next_requirement_url"`
	CompletedAt               time.Time `json:"completed_at"`
}

// Enrollment is an enrollment object
// https://canvas.instructure.com/doc/api/enrollments.html
type Enrollment struct {
	ID                   int    `json:"id"`
	CourseID             int    `json:"course_id"`
	CourseIntegrationID  string `json:"course_integration_id"`
	CourseSectionID      int    `json:"course_section_id"`
	SectionIntegrationID string `json:"section_integration_id"`

	EnrollmentState                string `json:"enrollment_state"`
	Role                           string `json:"role"`
	RoleID                         int    `json:"role_id"`
	Type                           string `json:"type"`
	LimitPrivilegesToCourseSection bool   `json:"limit_privileges_to_course_section"`
	UserID                         int    `json:"user_id"`
	User                           *User  `json:"user"`

	SisCourseID      string      `json:"sis_course_id"`
	SisAccountID     string      `json:"sis_account_id"`
	SisSectionID     string      `json:"sis_section_id"`
	SisUserID        string      `json:"sis_user_id"`
	SisImportID      int         `json:"sis_import_id"`
	RootAccountID    int         `json:"root_account_id"`
	AssociatedUserID interface{} `json:"associated_user_id"`

	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	StartAt           time.Time `json:"start_at"`
	EndAt             time.Time `json:"end_at"`
	LastActivityAt    time.Time `json:"last_activity_at"`
	LastAttendedAt    time.Time `json:"last_attended_at"`
	TotalActivityTime int       `json:"total_activity_time"`

	HTMLURL string `json:"html_url"`
	Grades  struct {
		HTMLURL              string  `json:"html_url"`
		CurrentScore         float64 `json:"current_score"`
		CurrentGrade         string  `json:"current_grade"`
		FinalScore           float64 `json:"final_score"`
		FinalGrade           string  `json:"final_grade"`
		UnpostedCurrentGrade string  `json:"unposted_current_grade"`
		UnpostedFinalGrade   string  `json:"unposted_final_grade"`
		UnpostedCurrentScore string  `json:"unposted_current_score"`
		UnpostedFinalScore   string  `json:"unposted_final_score"`
	} `json:"grades"`
	OverrideGrade                     string  `json:"override_grade"`
	OverrideScore                     float64 `json:"override_score"`
	UnpostedCurrentGrade              string  `json:"unposted_current_grade"`
	UnpostedFinalGrade                string  `json:"unposted_final_grade"`
	UnpostedCurrentScore              string  `json:"unposted_current_score"`
	UnpostedFinalScore                string  `json:"unposted_final_score"`
	HasGradingPeriods                 bool    `json:"has_grading_periods"`
	TotalsForAllGradingPeriodsOption  bool    `json:"totals_for_all_grading_periods_option"`
	CurrentGradingPeriodTitle         string  `json:"current_grading_period_title"`
	CurrentGradingPeriodID            int     `json:"current_grading_period_id"`
	CurrentPeriodOverrideGrade        string  `json:"current_period_override_grade"`
	CurrentPeriodOverrideScore        float64 `json:"current_period_override_score"`
	CurrentPeriodUnpostedCurrentScore float64 `json:"current_period_unposted_current_score"`
	CurrentPeriodUnpostedFinalScore   float64 `json:"current_period_unposted_final_score"`
	CurrentPeriodUnpostedCurrentGrade string  `json:"current_period_unposted_current_grade"`
	CurrentPeriodUnpostedFinalGrade   string  `json:"current_period_unposted_final_grade"`
}

// Quizzes will get all the course quizzes
func (c *Course) Quizzes(opts ...Option) ([]*Quiz, error) {
	return getQuizzes(c.client, c.ID, opts)
}

// Quiz will return a quiz given a quiz id.
func (c *Course) Quiz(id int, opts ...Option) (*Quiz, error) {
	return getQuiz(c.client, c.ID, id, opts)
}

func getQuizzes(client doer, courseID int, opts []Option) (qs []*Quiz, err error) {
	return qs, getjson(client, &qs, optEnc(opts), "courses/%d/quizzes", courseID)
}

func getQuiz(client doer, course, quiz int, opts []Option) (q *Quiz, err error) {
	return q, getjson(client, q, optEnc(opts), "courses/%d/quizzes/%d", course, quiz)
}

// Quiz is a quiz json response.
type Quiz struct {
	ID       int       `json:"id"`
	Title    string    `json:"title"`
	DueAt    time.Time `json:"due_at"`
	LockAt   time.Time `json:"lock_at"`
	UnlockAt time.Time `json:"unlock_at"`

	HTMLURL                       string          `json:"html_url"`
	MobileURL                     string          `json:"mobile_url"`
	PreviewURL                    string          `json:"preview_url"`
	Description                   string          `json:"description"`
	QuizType                      string          `json:"quiz_type"`
	AssignmentGroupID             int             `json:"assignment_group_id"`
	TimeLimit                     int             `json:"time_limit"`
	ShuffleAnswers                bool            `json:"shuffle_answers"`
	HideResults                   string          `json:"hide_results"`
	ShowCorrectAnswers            bool            `json:"show_correct_answers"`
	ShowCorrectAnswersLastAttempt bool            `json:"show_correct_answers_last_attempt"`
	ShowCorrectAnswersAt          time.Time       `json:"show_correct_answers_at"`
	HideCorrectAnswersAt          time.Time       `json:"hide_correct_answers_at"`
	OneTimeResults                bool            `json:"one_time_results"`
	ScoringPolicy                 string          `json:"scoring_policy"`
	AllowedAttempts               int             `json:"allowed_attempts"`
	OneQuestionAtATime            bool            `json:"one_question_at_a_time"`
	QuestionCount                 int             `json:"question_count"`
	PointsPossible                int             `json:"points_possible"`
	CantGoBack                    bool            `json:"cant_go_back"`
	AccessCode                    string          `json:"access_code"`
	IPFilter                      string          `json:"ip_filter"`
	Published                     bool            `json:"published"`
	Unpublishable                 bool            `json:"unpublishable"`
	LockedForUser                 bool            `json:"locked_for_user"`
	LockInfo                      interface{}     `json:"lock_info"`
	LockExplanation               string          `json:"lock_explanation"`
	SpeedgraderURL                string          `json:"speedgrader_url"`
	QuizExtensionsURL             string          `json:"quiz_extensions_url"`
	Permissions                   QuizPermissions `json:"permissions"`
	AllDates                      []string        `json:"all_dates"`
	VersionNumber                 int             `json:"version_number"`
	QuestionTypes                 []string        `json:"question_types"`
	AnonymousSubmissions          bool            `json:"anonymous_submissions"`
}

// QuizPermissions is the permissions for a quiz.
type QuizPermissions struct {
	Read           bool `json:"read"`
	Submit         bool `json:"submit"`
	Create         bool `json:"create"`
	Manage         bool `json:"manage"`
	ReadStatistics bool `json:"read_statistics"`
	ReviewGrades   bool `json:"review_grades"`
	Update         bool `json:"update"`
}

func (c *Course) filespager(ch chan *File, params []Option) *paginated {
	return newPaginatedList(
		c.client, c.id("/courses/%d/files"),
		sendFilesFunc(c.client, ch, nil),
		params,
	)
}

func (c *Course) folderspager(ch chan *Folder, params []Option) *paginated {
	return newPaginatedList(
		c.client, c.id("/courses/%d/folders"),
		sendFoldersFunc(c.client, ch, nil),
		params,
	)
}

func (c *Course) assignmentspager(ch chan *Assignment, params []Option) *paginated {
	return newPaginatedList(
		c.client, c.id("/courses/%d/assignments"),
		func(r io.Reader) error {
			asses := make([]*Assignment, 0, 10)
			err := json.NewDecoder(r).Decode(&asses)
			if err != nil {
				return err
			}
			for _, a := range asses {
				a.client = c.client
				a.courseCode = c.CourseCode
				ch <- a
			}
			return nil
		}, params,
	)
}

func (c *Course) collectUsers(path string, opts []Option) (users []*User, err error) {
	ch := make(chan *User)
	errs := newPaginatedList(
		c.client, fmt.Sprintf(path, c.ID),
		sendUserFunc(c.client, ch), opts,
	).start()
	for {
		select {
		case u := <-ch:
			users = append(users, u)
		case err := <-errs:
			return users, err
		}
	}
}

func sendFilesFunc(d doer, ch chan *File, folder *Folder) func(io.Reader) error {
	return func(r io.Reader) error {
		files := make([]*File, 0, defaultPerPage)
		err := json.NewDecoder(r).Decode(&files)
		if err != nil {
			return err
		}
		for _, f := range files {
			f.setclient(d)
			f.folder = folder
			ch <- f
		}
		return nil
	}
}

func sendFoldersFunc(d doer, ch chan *Folder, parent *Folder) sendFunc {
	return func(r io.Reader) error {
		folders := make([]*Folder, 0, defaultPerPage)
		err := json.NewDecoder(r).Decode(&folders)
		if err != nil {
			return err
		}
		for _, f := range folders {
			f.setclient(d)
			f.parent = parent
			ch <- f
		}
		return nil
	}
}

func sendUserFunc(d doer, ch chan *User) sendFunc {
	return func(r io.Reader) error {
		list := make([]*User, 0, defaultPerPage)
		err := json.NewDecoder(r).Decode(&list)
		if err != nil {
			return err
		}
		for _, u := range list {
			u.client = d
			ch <- u
		}
		return nil
	}
}

func defaultErrorHandler(err error) error {
	panic(err)
}

type assignmentChan chan *Assignment

func (ac assignmentChan) Close() {
	close(ac)
}

func (c *Course) id(s string) string {
	return fmt.Sprintf(s, c.ID)
}
