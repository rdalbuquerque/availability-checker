package server

import (
	"log"
	"net/http"
	"sync"
	"text/template"
	"time"

	"availability-checker/pkg/checker"
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
	var wg sync.WaitGroup
	resultsCh := make(chan checker.CheckResult)

	for _, c := range s.checkers {
		wg.Add(1)
		go func(c checker.Checker) {
			defer wg.Done()
			success, err := c.Check()
			if err != nil {
				log.Printf("Error while checking %s: %s\n", c.Name(), err)
			}
			resultsCh <- checker.CheckResult{Name: c.Name(), Status: success, LastChecked: time.Now(), IsFixable: c.IsFixable()}
		}(c)
	}

	go func() {
		wg.Wait()
		close(resultsCh)
	}()

	for r := range resultsCh {
		s.results = append(s.results, r)
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		s.mu.Lock()
		defer s.mu.Unlock()
		err := s.template.Execute(w, s.results)
		if err != nil {
			log.Printf("Error while executing template: %s\n", err)
		}
	case "/fix":
		s.fixChecker(w, r)

	default:
		http.NotFound(w, r)
	}
}

func (s *Server) fixChecker(w http.ResponseWriter, r *http.Request) {
	checkerName := r.URL.Query().Get("checker")
	if checkerName == "" {
		http.Error(w, "Missing checker parameter", http.StatusBadRequest)
		return
	}

	var check checker.Checker
	for _, c := range s.checkers {
		if c.Name() == checkerName {
			check = c
			break
		}
	}
	if check == nil {
		http.Error(w, "Invalid checker", http.StatusBadRequest)
		return
	}

	if !check.IsFixable() {
		http.Error(w, "Checker is not fixable", http.StatusBadRequest)
		return
	}

	err := check.Fix()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
