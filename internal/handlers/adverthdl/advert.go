package adverthdl

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"simple-search-es/internal/services/advertsvc"
)

type Handler struct {
	service advertsvc.Service
}

func New(service advertsvc.Service) Handler {
	return Handler{
		service: service,
	}
}

// Create - POST /api/adverts
func (h Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req advertsvc.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}

	res, err := h.service.Create(r.Context(), req)
	if err != nil {
		switch err {
		case errors.New("conflict"):
			w.WriteHeader(http.StatusConflict)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	bdy, _ := json.Marshal(res)
	_, _ = w.Write(bdy)
}

// SearchOne - GET /api/adverts/{id}
func (h Handler) SearchOne(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	res, err := h.service.SearchOne(r.Context(), advertsvc.SearchOneRequest{ID: id})
	if err != nil {
		switch err {
		case errors.New("not found"):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	bdy, _ := json.Marshal(res)
	_, _ = w.Write(bdy)
}

// Search - GET /api/adverts
func (h Handler) Search(w http.ResponseWriter, r *http.Request) {
	res, err := h.service.Search(r.Context(), advertsvc.SearchRequest{Values: r.URL.Query()})
	if err != nil {
		switch err {
		case errors.New("not found"):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	bdy, _ := json.Marshal(res)
	_, _ = w.Write(bdy)
}

// Delete - DELETE /api/adverts/{id}
func (h Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := h.service.Delete(r.Context(), advertsvc.DeleteRequest{ID: id}); err != nil {
		switch err {
		case errors.New("not found"):
			w.WriteHeader(http.StatusNotFound)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}