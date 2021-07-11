package database

import (
	"fmt"
	"strings"
)

type CourseStore struct {
	ID       int64
	Email    string
	Name     string
	Password string
}

func (s *Store) CreateCourse(courses []map[string]string) error {
	key := []string{}
	val := []string{}
	for i, course := range courses {
		k := []string{
			fmt.Sprintf("$%d", i*3+1),
			fmt.Sprintf("$%d", i*3+2),
			fmt.Sprintf("$%d", i*3+3),
		}
		key = append(key, fmt.Sprintf("(%s)", strings.Join(k, ",")))

		val = append(val, []string{
			course["title"],
			course["description"],
			course["author"],
		}...)
	}
	sql := fmt.Sprintf("INSERT INTO courses (name, description, author) VALUES %s ON CONFLICT (name, author) DO NOTHING", strings.Join(key, ","))
	_, err := s.db.Exec(s.ctx, sql, Strings(val).ToInterfaceSlice()...)

	if err != nil {
		s.logger.Error("Error in inserting courses:", err)
	}

	return err
}

func (s *Store) GetCourseList(query string, limit, offset int) ([]map[string]string, error) {
	sql := `SELECT name, description, author
				FROM courses
				WHERE name ILIKE $1 OR author ILIKE $1 OR description ILIKE $1
				LIMIT $2 OFFSET $3`

	query = "%" + query + "%"
	fmt.Println(query)
	rows, err := s.db.Query(s.ctx, sql, query, limit, offset)
	if err != nil {
		s.logger.Error("Error in fetching course list:", err, sql)
		return nil, err
	}

	list := []map[string]string{}
	for rows.Next() {
		var name, description, author string
		err := rows.Scan(&name, &description, &author)
		if err != nil {
			s.logger.Error("Error in scanning course list:", err)
			continue
		}
		list = append(list, map[string]string{
			"name":        name,
			"description": description,
			"author":      author,
		})
	}

	return list, nil
}

type Strings []string

func (ss Strings) ToInterfaceSlice() []interface{} {
	iface := make([]interface{}, len(ss))
	for i := range ss {
		iface[i] = ss[i]
	}
	return iface
}
