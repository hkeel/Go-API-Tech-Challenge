package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/hkeel/Go-API-Tech-Challenge/internal/models"
	"github.com/hkeel/Go-API-Tech-Challenge/internal/services"
)

type PersonHandler struct {
	Service *services.PersonService
}

func NewPersonHandler(service *services.PersonService) *PersonHandler {
	return &PersonHandler{Service: service}
}

func (h *PersonHandler) GetAllPeople(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	ageStr := r.URL.Query().Get("age")
	var age int
	if ageStr != "" {
		var err error
		age, err = strconv.Atoi(ageStr)
		if err != nil {
			http.Error(w, "Invalid age parameter", http.StatusBadRequest)
			return
		}
	}
	people, err := h.Service.GetPeople(name, age)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(people)
}

func (h *PersonHandler) GetPerson(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	people, err := h.Service.GetPeople(name, 0)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(people) == 0 {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}
	if len(people) > 1 {
		http.Error(w, "Multiple people found", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(people[0])
}

func (h *PersonHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	var person models.Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	person, err = h.Service.UpdatePerson(name, person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(person)
}

func (h *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person models.Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	person, err = h.Service.CreatePerson(person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(person)
}

func (h *PersonHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	err := h.Service.DeletePerson(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Person deleted successfully"})
}
