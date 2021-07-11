package interfaces

type Store interface {
	// course related
	CreateCourse(courses []map[string]string) error
	GetCourseList(query string, page, offset int) ([]map[string]string, error)
}
