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
