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

func (r *CourseRepository) GetCourse(id string) (*models.Course, error) {
	var course models.Course
	if err := r.DB.QueryRow("SELECT id, name FROM course WHERE id = $1", id).Scan(&course.ID, &course.Name); err != nil {
		return nil, err
	}

	return &course, nil
}

func (r *CourseRepository) UpdateCourse(id string, course models.Course) (models.Course, error) {
	if _, err := r.DB.Exec("UPDATE course SET name = $1 WHERE id = $2", course.Name, id); err != nil {
		return models.Course{}, err
	}

	return course, nil
}

func (r *CourseRepository) CreateCourse(course models.Course) (models.Course, error) {
	var id int
	err := r.DB.QueryRow("INSERT INTO course (name) VALUES ($1) RETURNING id", course.Name).Scan(&id)
	if err != nil {
		return models.Course{}, err
	}

	course.ID = id
	return course, nil
}

func (r *CourseRepository) DeleteCourse(id string) error {
	_, err := r.DB.Exec("DELETE FROM course WHERE id = $1", id)
	return err
}
