package main

import (
	"io/ioutil"
	"net/http"

	"availability-checker/checker"
	"availability-checker/server"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Checkers []struct {
		Type             string
		URL              string `yaml:",omitempty"`
		ConnectionString string `yaml:",omitempty"`
	}
}

func main() {
	data, _ := ioutil.ReadFile("config.yaml")
	var config Config
	yaml.Unmarshal(data, &config)

	checkers := make([]checker.Checker, len(config.Checkers))
	for i, confChecker := range config.Checkers {
		switch confChecker.Type {
		case "http":
			checkers[i] = &checker.HttpChecker{URL: confChecker.URL}
			// case "sql":
			// 	checkers[i] = &checker.SqlChecker{ConnectionString: confChecker.ConnectionString}
			// case "vertica":
			// 	checkers[i] = &checker.VerticaChecker{ConnectionString: confChecker.ConnectionString}
		}
	}

	serverInstance := server.NewServer(checkers, "template.gotmpl")

	go serverInstance.StartChecking()

	http.ListenAndServe(":8080", serverInstance)
}
