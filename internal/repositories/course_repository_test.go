package repositories

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/hkeel/Go-API-Tech-Challenge/internal/models"
	"github.com/stretchr/testify/assert"
)

func setupMockCourseRepo(t *testing.T) (*sqlmock.Sqlmock, *CourseRepository) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	repo := &CourseRepository{DB: db}
	return &mock, repo
}

func TestGetAllCourses(t *testing.T) {
	mock, repo := setupMockCourseRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectQuery("SELECT id, name FROM course").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Math").
			AddRow(2, "Science"))

	courses, err := repo.GetAllCourses()
	assert.NoError(t, err)
	assert.Len(t, courses, 2)
	assert.Equal(t, 1, courses[0].ID)
	assert.Equal(t, "Math", courses[0].Name)
	assert.Equal(t, 2, courses[1].ID)
	assert.Equal(t, "Science", courses[1].Name)
}

func TestGetAllCourses_ReturnsError(t *testing.T) {
	mock, repo := setupMockCourseRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectQuery("SELECT id, name FROM course").
		WillReturnError(assert.AnError)

	_, err := repo.GetAllCourses()
	assert.Error(t, err)
}

func TestGetCourse(t *testing.T) {
	mock, repo := setupMockCourseRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectQuery(`SELECT id, name FROM course WHERE id = \$1`).
		WithArgs("1").
		WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Math"))
	course, err := repo.GetCourse("1")
	assert.Equal(t, 1, course.ID)
	assert.NoError(t, err)
}

func TestGetCourse_ReturnsError(t *testing.T) {
	mock, repo := setupMockCourseRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectQuery("SELECT id, name FROM course").
		WithArgs("1").
		WillReturnError(assert.AnError)

	_, err := repo.GetCourse("1")
	assert.Error(t, err)
}

func TestUpdateCourse(t *testing.T) {
	mock, repo := setupMockCourseRepo(t)
	defer teardownMock(mock, t)

	course := models.Course{ID: 1, Name: "Math"}

	(*mock).ExpectExec(`UPDATE course SET name = \$1 WHERE id = \$2`).
		WithArgs("Math", "1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	course, err := repo.UpdateCourse("1", course)
	assert.NoError(t, err)
	assert.Equal(t, "Math", course.Name)
}

func TestUpdateCourse_ReturnsError(t *testing.T) {
	mock, repo := setupMockCourseRepo(t)
	defer teardownMock(mock, t)

	course := models.Course{ID: 1, Name: "Math"}

	(*mock).ExpectExec(`UPDATE course SET name = \$1 WHERE id = \$2`).
		WithArgs("Math", "1").
		WillReturnError(assert.AnError)

	_, err := repo.UpdateCourse("1", course)
	assert.Error(t, err)
}

func TestCreateCourse(t *testing.T) {
	mock, repo := setupMockCourseRepo(t)
	defer teardownMock(mock, t)

	course := models.Course{Name: "Math"}

	(*mock).ExpectQuery(`INSERT INTO course \(name\) VALUES \(\$1\) RETURNING id`).
		WithArgs("Math").
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))

	course, err := repo.CreateCourse(course)
	assert.NoError(t, err)
	assert.Equal(t, 1, course.ID)
}

func TestCreateCourse_ReturnsError(t *testing.T) {
	mock, repo := setupMockCourseRepo(t)
	defer teardownMock(mock, t)

	course := models.Course{Name: "Math"}

	(*mock).ExpectQuery(`INSERT INTO course \(name\) VALUES \(\$1\) RETURNING id`).
		WithArgs("Math").
		WillReturnError(assert.AnError)

	_, err := repo.CreateCourse(course)
	assert.Error(t, err)
}

func TestDeleteCourse(t *testing.T) {
	mock, repo := setupMockCourseRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectExec(`DELETE FROM course WHERE id = \$1`).
		WithArgs("1").
		WillReturnResult(sqlmock.NewResult(1, 1))

	err := repo.DeleteCourse("1")
	assert.NoError(t, err)
}

func TestDeleteCourse_ReturnsError(t *testing.T) {
	mock, repo := setupMockCourseRepo(t)
	defer teardownMock(mock, t)

	(*mock).ExpectExec(`DELETE FROM course WHERE id = \$1`).
		WithArgs("1").
		WillReturnError(assert.AnError)

	err := repo.DeleteCourse("1")
	assert.Error(t, err)
}
