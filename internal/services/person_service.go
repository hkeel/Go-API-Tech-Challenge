package services

import (
	"github.com/hkeel/Go-API-Tech-Challenge/internal/models"
	"github.com/hkeel/Go-API-Tech-Challenge/internal/repositories"
)

type PersonService struct {
	Repo *repositories.PersonRepository
}

func (s *PersonService) GetPeople(name string, age int) ([]models.Person, error) {
	return s.Repo.GetPeople(name, age)
}

func (s *PersonService) UpdatePerson(name string, person models.Person) (models.Person, error) {
	person, err := s.Repo.UpdatePerson(name, person)
	if err != nil {
		return person, err
	}
	err = s.Repo.DeletePersonFromCourses(person.ID)
	if err != nil {
		return person, err
	}
	for _, courseID := range person.Courses {
		err = s.Repo.AddPersonToCourse(person.ID, courseID)
		if err != nil {
			return person, err
		}
	}
	return person, nil
}

func (s *PersonService) CreatePerson(person models.Person) (models.Person, error) {
	person, err := s.Repo.CreatePerson(person)
	if err != nil {
		return person, err
	}
	for _, courseID := range person.Courses {
		err = s.Repo.AddPersonToCourse(person.ID, courseID)
		if err != nil {
			return person, err
		}
	}
	return person, nil
}

func (s *PersonService) DeletePerson(name string) error {
	// Get the person ID
	person, err := s.Repo.GetPeople(name, 0)
	if err != nil {
		return err
	}
	// DELETE the person from courses
	err = s.Repo.DeletePersonFromCourses(person[0].ID)
	if err != nil {
		return err
	}
	// Then DELETE the person
	err = s.Repo.DeletePerson(name)
	if err != nil {
		return err
	}
	return nil
}
