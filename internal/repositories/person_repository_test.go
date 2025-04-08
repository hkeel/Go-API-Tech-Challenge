package repositories

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hkeel/Go-API-Tech-Challenge/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupMockPersonRepo(t *testing.T) (*sqlmock.Sqlmock, *PersonRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repo := &PersonRepository{DB: db}
	return &mock, repo
}

func teardownMock(mock *sqlmock.Sqlmock, t *testing.T) {
	if err := (*mock).ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestCreatePerson(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	person := models.Person{
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
	}

	(*mock).ExpectQuery("INSERT INTO person").
		WithArgs(person.FirstName, person.LastName, person.Type, person.Age).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	createdPerson, err := repo.CreatePerson(person)
	assert.NoError(t, err)
	assert.Equal(t, 1, createdPerson.ID)
	assert.Equal(t, "John", createdPerson.FirstName)
	assert.Equal(t, "Doe", createdPerson.LastName)
	assert.Equal(t, "student", createdPerson.Type)
	assert.Equal(t, 25, createdPerson.Age)
}

func TestCreatePerson_ReturnsError(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	person := models.Person{
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
	}

	(*mock).ExpectQuery("INSERT INTO person").
		WithArgs(person.FirstName, person.LastName, person.Type, person.Age).
		WillReturnError(assert.AnError)

	createdPerson, err := repo.CreatePerson(person)
	assert.Error(t, err)
	assert.Equal(t, models.Person{}, createdPerson)
}

func TestGetPeople(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name", "type", "age"}).
		AddRow(1, "John", "Doe", "student", 25).
		AddRow(2, "Jane", "Doe", "teacher", 30)

	(*mock).ExpectQuery(`SELECT id, first_name, last_name, type, age FROM person WHERE 1=1 AND \(first_name ILIKE \$1 OR last_name ILIKE \$1\)`).
		WithArgs("%John%").
		WillReturnRows(rows)

	// Mock expectation for the person_course query for person ID 1
	courseRows1 := sqlmock.NewRows([]string{"course_id"}).
		AddRow(101).
		AddRow(102)
	(*mock).ExpectQuery(`SELECT course_id FROM person_course WHERE person_id = \$1`).
		WithArgs(1).
		WillReturnRows(courseRows1)

	// Mock expectation for the person_course query for person ID 2
	courseRows2 := sqlmock.NewRows([]string{"course_id"}).
		AddRow(201).
		AddRow(202)
	(*mock).ExpectQuery(`SELECT course_id FROM person_course WHERE person_id = \$1`).
		WithArgs(2).
		WillReturnRows(courseRows2)

	people, err := repo.GetPeople("John", 0)

	assert.NoError(t, err)
	assert.Len(t, people, 2)
	assert.Equal(t, "John", people[0].FirstName)
	assert.Equal(t, "Jane", people[1].FirstName)
}

func TestGetPeople_ReturnsError(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectQuery(`SELECT id, first_name, last_name, type, age FROM person WHERE 1=1 AND \(first_name ILIKE \$1 OR last_name ILIKE \$1\)`).
		WithArgs("%John%").
		WillReturnError(assert.AnError)

	people, err := repo.GetPeople("John", 0)
	assert.Error(t, err)
	assert.Nil(t, people)
}

func TestUpdatePerson(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	person := models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
	}

	(*mock).ExpectExec(`UPDATE person SET first_name = \$1, last_name = \$2, type = \$3, age = \$4 WHERE first_name = \$5`).
		WithArgs(person.FirstName, person.LastName, person.Type, person.Age, "John").
		WillReturnResult(sqlmock.NewResult(1, 1))

	updatedPerson, err := repo.UpdatePerson("John", person)
	assert.NoError(t, err)
	assert.Equal(t, 1, updatedPerson.ID)
}

func TestUpdatePerson_ReturnsError(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	person := models.Person{
		ID:        1,
		FirstName: "John",
		LastName:  "Doe",
		Type:      "student",
		Age:       25,
	}

	(*mock).ExpectExec(`UPDATE person SET first_name = \$1, last_name = \$2, type = \$3, age = \$4 WHERE first_name = \$5`).
		WithArgs(person.FirstName, person.LastName, person.Type, person.Age, "John").
		WillReturnError(assert.AnError)

	updatedPerson, err := repo.UpdatePerson("John", person)
	assert.Error(t, err)
	assert.Equal(t, models.Person{}, updatedPerson)
}

func TestAddPersonToCourse(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectExec(`INSERT INTO person_course \(person_id, course_id\) VALUES \(\$1, \$2\)`).
		WithArgs(1, 101).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.AddPersonToCourse(1, 101)
	assert.NoError(t, err)
}

func TestAddPersonToCourse_ReturnsError(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectExec(`INSERT INTO person_course \(person_id, course_id\) VALUES \(\$1, \$2\)`).
		WithArgs(1, 101).
		WillReturnError(assert.AnError)

	err := repo.AddPersonToCourse(1, 101)
	assert.Error(t, err)
}

func TestDeletePersonFromCourses(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectExec(`DELETE FROM person_course WHERE person_id = \$1`).
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeletePersonFromCourses(1)
	assert.NoError(t, err)
}

func TestDeletePersonFromCourses_ReturnsError(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectExec(`DELETE FROM person_course WHERE person_id = \$1`).
		WithArgs(1).
		WillReturnError(assert.AnError)

	err := repo.DeletePersonFromCourses(1)
	assert.Error(t, err)
}

func TestGetCoursesForPerson(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectQuery(`SELECT course_id FROM person_course WHERE person_id = \$1`).
		WithArgs(1).
		WillReturnRows(sqlmock.NewRows([]string{"course_id"}).AddRow(101).AddRow(102))

	courses, err := repo.GetCoursesForPerson(1)
	assert.NoError(t, err)
	assert.Len(t, courses, 2)
	assert.Equal(t, 101, courses[0])
	assert.Equal(t, 102, courses[1])
}

func TestGetCoursesForPerson_ReturnsError(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectQuery(`SELECT course_id FROM person_course WHERE person_id = \$1`).
		WithArgs(1).
		WillReturnError(assert.AnError)

	courses, err := repo.GetCoursesForPerson(1)
	assert.Error(t, err)
	assert.Nil(t, courses)
}

func TestDeletePerson(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectExec(`DELETE FROM person WHERE first_name = \$1`).
		WithArgs("John").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeletePerson("John")
	assert.NoError(t, err)
}

func TestDeletePerson_ReturnsError(t *testing.T) {
	mock, repo := setupMockPersonRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectExec(`DELETE FROM person WHERE first_name = \$1`).
		WithArgs("John").
		WillReturnError(assert.AnError)

	err := repo.DeletePerson("John")
	assert.Error(t, err)
}
