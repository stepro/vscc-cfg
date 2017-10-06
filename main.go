package main

import (
	"io/ioutil"
	"log"
	"os"
	"text/template"

	"gopkg.in/yaml.v2"
)

// Document represents a configuration file
type Document struct {
	Kind string
}

// FunctionProperties represents configuration for a function
type FunctionProperties struct {
}

// FunctionDocument represents a configuration file for a function
type FunctionDocument struct {
	Name           string
	Version        string
	Labels         map[string]string
	Properties     FunctionProperties
	Configurations map[string]FunctionProperties
}

// JobProperties represents configuration for a job
type JobProperties struct {
}

// JobDocument represents a configuration file for a job
type JobDocument struct {
	Name           string
	Version        string
	Labels         map[string]string
	Properties     JobProperties
	Configurations map[string]JobProperties
}

// ServiceProperties represents configuration for a service
type ServiceProperties struct {
	Build struct {
		Args       map[string]string
		Context    string
		Dockerfile string
		Labels     map[string]string
		Target     string
	}
	Command    []string
	Entrypoint []string
	Env        map[string]string
	Imports    map[string]string
	Init       *bool
	Ports      map[string]string
	Public     *bool
	PublicPort *int `yaml:"publicPort"`
	References []string
	Sync       *bool
	SyncTarget string `yaml:"syncTarget"`
	Tasks      map[string]struct {
		Command     []string
		Env         map[string]string
		Interactive *bool
		TTY         *bool
		User        string
		Workdir     string
	}
	TTY   *bool
	User  string
	Watch map[string]struct {
		Action string // ignore, sync or rebuild
	}
	Workdir string
}

// ServiceDocument represents a configuration file for a service
type ServiceDocument struct {
	Name           string
	Version        string
	Labels         map[string]string
	Properties     ServiceProperties
	Configurations map[string]ServiceProperties
}

func main() {
	tmpl, err := template.New("template").Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	bytes, _ := ioutil.ReadFile(os.Args[2])

	kind := "service"
	var doc Document
	err = yaml.Unmarshal(bytes, &doc)
	if err != nil {
		log.Fatal(err)
	}
	if doc.Kind != "" {
		kind = doc.Kind
	}

	if kind == "function" {
		var function FunctionDocument
		err = yaml.UnmarshalStrict(bytes, &function)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(os.Stdout, function)
	} else if kind == "job" {
		var job JobDocument
		err = yaml.UnmarshalStrict(bytes, &job)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(os.Stdout, job)
	} else if kind == "service" {
		var service ServiceDocument
		err = yaml.UnmarshalStrict(bytes, &service)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(os.Stdout, service)
	} else {
		log.Fatal("Unknown kind '" + kind + "'")
	}
	if err != nil {
		log.Fatal(err)
	}
}
