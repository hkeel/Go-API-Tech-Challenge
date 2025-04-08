package repositories

import (
	"database/sql"

	"github.com/hkeel/Go-API-Tech-Challenge/internal/models"
)

type PersonRepository struct {
	DB *sql.DB
}

func (r *PersonRepository) GetPeople(name string, age int) ([]models.Person, error) {
	query := "SELECT id, first_name, last_name, type, age FROM person WHERE 1=1"
	args := []interface{}{}

	if name != "" {
		query += " AND (first_name ILIKE $1 OR last_name ILIKE $1)"
		args = append(args, "%"+name+"%")
	}

	if age > 0 {
		query += " AND age = $2"
		args = append(args, age)
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var people []models.Person
	for rows.Next() {
		var person models.Person
		if err := rows.Scan(&person.ID, &person.FirstName, &person.LastName, &person.Type, &person.Age); err != nil {
			return nil, err
		}
		// Fetch courses for the person
		courses, err := r.GetCoursesForPerson(person.ID)
		if err != nil {
			return nil, err
		}
		person.Courses = courses
		people = append(people, person)
	}

	return people, nil
}

func (r *PersonRepository) UpdatePerson(name string, person models.Person) (models.Person, error) {
	if _, err := r.DB.Exec("UPDATE person SET first_name = $1, last_name = $2, type = $3, age = $4 WHERE first_name = $5", person.FirstName, person.LastName, person.Type, person.Age, name); err != nil {
		return models.Person{}, err
	}

	return person, nil
}

func (r *PersonRepository) CreatePerson(person models.Person) (models.Person, error) {
	var id int
	err := r.DB.QueryRow("INSERT INTO person (first_name, last_name, type, age) VALUES ($1, $2, $3, $4) RETURNING id", person.FirstName, person.LastName, person.Type, person.Age).Scan(&id)
	if err != nil {
		return models.Person{}, err
	}

	person.ID = id
	return person, nil
}

func (r *PersonRepository) AddPersonToCourse(personID, courseID int) error {
	_, err := r.DB.Exec("INSERT INTO person_course (person_id, course_id) VALUES ($1, $2)", personID, courseID)
	return err
}

func (r *PersonRepository) DeletePersonFromCourses(personID int) error {
	_, err := r.DB.Exec("DELETE FROM person_course WHERE person_id = $1", personID)
	return err
}

func (r *PersonRepository) GetCoursesForPerson(personID int) ([]int, error) {
	rows, err := r.DB.Query("SELECT course_id FROM person_course WHERE person_id = $1", personID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var courses []int
	for rows.Next() {
		var courseID int
		if err := rows.Scan(&courseID); err != nil {
			return nil, err
		}
		courses = append(courses, courseID)
	}

	return courses, nil
}

func (r *PersonRepository) DeletePerson(name string) error {
	_, err := r.DB.Exec("DELETE FROM person WHERE first_name = $1", name)
	if err != nil {
		return err
	}

	return nil
}
