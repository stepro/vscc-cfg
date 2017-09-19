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
	DocVersion string `yaml:"docVersion"`
	Kind       string
}

// FunctionPropertiesV1 represents configuration for a function
type FunctionPropertiesV1 struct {
}

// FunctionDocumentV1 represents a configuration file for a function
type FunctionDocumentV1 struct {
	DocVersion     string `yaml:"docVersion"`
	Kind           string
	Name           string
	Version        string
	Language       string
	Labels         map[string]string
	Properties     FunctionPropertiesV1
	Configurations map[string]FunctionPropertiesV1
}

// JobPropertiesV1 represents configuration for a job
type JobPropertiesV1 struct {
}

// JobDocumentV1 represents a configuration file for a job
type JobDocumentV1 struct {
	DocVersion     string `yaml:"docVersion"`
	Kind           string
	Name           string
	Version        string
	Language       string
	Labels         map[string]string
	Properties     JobPropertiesV1
	Configurations map[string]JobPropertiesV1
}

// ServicePropertiesV1 represents configuration for a service
type ServicePropertiesV1 struct {
	Build struct {
		Args       map[string]string
		Context    string
		Dockerfile string
		Labels     map[string]string
		Tags       []string
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
		Action  string // restart or exec
		Command []string
	}
	Workdir string
}

// ServiceDocumentV1 represents a configuration file for a service
type ServiceDocumentV1 struct {
	DocVersion     string `yaml:"docVersion"`
	Kind           string
	Name           string
	Version        string
	Language       string
	Labels         map[string]string
	Properties     ServicePropertiesV1
	Configurations map[string]ServicePropertiesV1
}

func main() {
	tmpl, err := template.New("template").Parse(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	fileName := "vscc.yaml"
	if len(os.Args) > 2 {
		fileName = os.Args[2]
	}
	bytes, _ := ioutil.ReadFile(fileName)

	docVersion := "1.0.0"
	kind := "service"
	var doc Document
	err = yaml.Unmarshal(bytes, &doc)
	if err != nil {
		log.Fatal(err)
	}
	if doc.DocVersion != "" {
		docVersion = doc.DocVersion
	}
	if doc.Kind != "" {
		kind = doc.Kind
	}

	if docVersion != "1.0.0" {
		log.Fatal("Unknown document version '" + docVersion + "'")
	}

	if kind == "function" {
		var function FunctionDocumentV1
		err = yaml.UnmarshalStrict(bytes, &function)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(os.Stdout, function)
	} else if kind == "job" {
		var job JobDocumentV1
		err = yaml.UnmarshalStrict(bytes, &job)
		if err != nil {
			log.Fatal(err)
		}
		err = tmpl.Execute(os.Stdout, job)
	} else if kind == "service" {
		var service ServiceDocumentV1
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
