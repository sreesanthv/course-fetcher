package services

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/buger/jsonparser"
	"github.com/sirupsen/logrus"
	"github.com/sreesanthv/course-fetcher/interfaces"
)

type CourseService struct {
	Logger *logrus.Logger
	Store  interfaces.Store
}

func NewCourseService(log *logrus.Logger, store interfaces.Store) *CourseService {
	return &CourseService{
		Logger: log,
		Store:  store,
	}
}

type parsedCourse struct {
	name        string
	description string
	author      string
}

func (s *CourseService) Fetch(query string, limit int) error {
	page := 1
	courseIndex := map[string]bool{}
	for {
		courses, err := s.parseCourses(query, page, courseIndex)
		if err != nil {
			return err
		}
		count := len(courseIndex)

		if count > limit {
			trim := count - limit
			courses = courses[:len(courses)-trim]
			break
		}

		if len(courses) > 0 {
			s.Store.CreateCourse(courses)
		}
		s.Logger.Infof("Parsed total %d courses from website", count)

		if page > 5 {
			break
		}

		page++
	}

	return nil
}

// parse courses from website
func (s *CourseService) parseCourses(query string, page int, courseIndex map[string]bool) ([]map[string]string, error) {
	url := fmt.Sprintf("https://www.coursera.org/search?query=%s&page=%d", query, page)
	res, err := http.Get(url)
	if err != nil {
		s.Logger.Error("Error connecting coursera", err)
		return nil, err
	} else {
		s.Logger.Info("Fetching courses from: ", url)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		s.Logger.Error("Error connecting coursera")
		return nil, fmt.Errorf("Invalid request")
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		s.Logger.Error("Error reading coursera HTML body")
	}

	courses := []map[string]string{}

	data := doc.Find("#__NEXT_DATA__").Text()
	jsonparser.ArrayEach([]byte(data), func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		hits, _, _, err := jsonparser.Get(value, "content", "hits")
		if err == nil {
			jsonparser.ArrayEach(hits, func(course []byte, dataType jsonparser.ValueType, offset int, err error) {
				title, err := jsonparser.GetString(course, "name")
				description, err := jsonparser.GetString(course, "_snippetResult", "description", "value")
				author, _, _, err := jsonparser.Get(course, "partners")
				if err == nil {
					auth := []string{}
					json.Unmarshal(author, &auth)
					courses = append(courses, map[string]string{
						"title":       title,
						"description": description,
						"author":      strings.Join(auth, ","),
					})
					courseIndex[fmt.Sprintf("%s-%s", title, author)] = true
				}
			})
		}

	}, "props", "pageProps", "resultsState")

	return courses, nil
}

func (s *CourseService) GetList(query string, limi, offset int) ([]map[string]string, error) {
	return s.Store.GetCourseList(query, limi, offset)
}
