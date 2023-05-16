package server

import (
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"availability-checker/checker"
)

type Server struct {
	checkers []checker.Checker
	results  []checker.CheckResult
	mu       sync.Mutex
	template *template.Template
}

func NewServer(checkers []checker.Checker, templateFile string) *Server {
	tmpl := template.Must(template.ParseFiles(templateFile))

	return &Server{
		checkers: checkers,
		template: tmpl,
	}
}

func (s *Server) StartChecking() {
	for {
		s.checkAll()
		time.Sleep(30 * time.Second)
	}
}

func (s *Server) checkAll() {
	results := make([]checker.CheckResult, len(s.checkers))

	for i, c := range s.checkers {
		success, err := c.Check()
		if err != nil {
			log.Printf("Error while checking %s: %s\n", c.Name(), err)
		}
		results[i] = checker.CheckResult{Name: c.Name(), Status: success}
	}

	s.mu.Lock()
	s.results = results
	s.mu.Unlock()
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.template.Execute(w, s.results)
}
