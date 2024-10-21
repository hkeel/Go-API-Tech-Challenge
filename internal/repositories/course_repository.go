package repositories

import (
	"database/sql"

	"github.com/hkeel/Go-API-Tech-Challenge/internal/models"
)

type CourseRepository struct {
	DB *sql.DB
}

func (r *CourseRepository) GetAllCourses() ([]models.Course, error) {
	rows, err := r.DB.Query("SELECT id, name FROM course")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []models.Course
	for rows.Next() {
		var course models.Course
		if err := rows.Scan(&course.ID, &course.Name); err != nil {
			return nil, err
		}
		courses = append(courses, course)
	}

	return courses, nil
}
