package services

import (
	"github.com/hkeel/Go-API-Tech-Challenge/internal/models"
	"github.com/hkeel/Go-API-Tech-Challenge/internal/repositories"
)

type CourseService struct {
	Repo *repositories.CourseRepository
}

func (s *CourseService) GetAllCourses() ([]models.Course, error) {
	return s.Repo.GetAllCourses()
}

func (s *CourseService) GetCourse(id string) (*models.Course, error) {
	return s.Repo.GetCourse(id)
}

func (s *CourseService) UpdateCourse(id string, course models.Course) (models.Course, error) {
	return s.Repo.UpdateCourse(id, course)
}

func (s *CourseService) CreateCourse(course models.Course) (models.Course, error) {
	return s.Repo.CreateCourse(course)
}

func (s *CourseService) DeleteCourse(id string) error {
	return s.Repo.DeleteCourse(id)
}
