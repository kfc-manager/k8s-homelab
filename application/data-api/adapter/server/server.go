package server

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/kfc-manager/k8s-homelab/application/data-api/service/store"
)

type server struct {
	port   string
	store  store.Service
	router *http.ServeMux
}

func New(port string, store store.Service) *server {
	s := &server{
		port:   port,
		router: http.NewServeMux(),
		store:  store,
	}
	s.router.HandleFunc("/{hash}", s.handler)
	return s
}

func (s *server) Listen() error {
	return http.ListenAndServe(fmt.Sprintf(":%s", s.port), s.router)
}

func (s *server) handler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		s.handleGet(w, r)
		return
	}

	if r.Method == "POST" {
		s.handlePost(w, r)
		return
	}

	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
}

func (s *server) handleGet(w http.ResponseWriter, r *http.Request) {
	b, err := s.store.Get(r.PathValue("hash"))
	if err != nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(b)
}

func (s *server) handlePost(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	err = s.store.Set(r.PathValue("hash"), b)
	if err != nil {
		if err.Error() == "hash mismatch" {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		log.Println(err.Error())
		http.Error(w, "Internal Server", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprintf(w, "Created")
}
