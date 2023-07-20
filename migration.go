package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
)

type Questionnaire struct {
	Categories []Category `yaml:"categories"`
}

type Category struct {
	Title     string     `yaml:"title"`
	Order     uint       `yaml:"order"`
	Questions []Question `yaml:"questions"`
}

// Question represents a question in a category of a Questionnaire.
type Question struct {
	// Name is a short internal identifier that is not displayed.
	Name  string `yaml:"name"`
	Order uint   `yaml:"order"`
	// Question is the text of the question.
	Question string `yaml:"question"`
	// Description is displayed as the question's tooltip in Konveyor.
	Description string   `yaml:"description"`
	Options     []Option `yaml:"options"`
}

// Option represents a single answer to a multiple-choice question.
type Option struct {
	// Option is the display text for the option.
	Option string `yaml:"option"`
	Order  uint   `yaml:"order"`
	// Risk must be one-of UNKNOWN, RED, AMBER, GREEN.
	Risk string `yaml:"risk"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Generates a sql migration from a questionnaire yaml and outputs it to stdout.")
		fmt.Printf("Usage: %s infile > outfile\n", os.Args[0])
		os.Exit(1)
	}

	in, err := os.Open(os.Args[1])
	if err != nil {
		panic(err)
	}

	questionnaire := Questionnaire{}
	decoder := yaml.NewDecoder(in)
	err = decoder.Decode(&questionnaire)
	if err != nil {
		panic(err)
	}

	// mark the old questions deleted while leaving them for previous assessments.
	fmt.Println("UPDATE category SET deleted = true;")
	fmt.Println("UPDATE question SET deleted = true;")
	fmt.Println("UPDATE single_option SET deleted = true;")

	// insert new categories, questions, and options.
	for _, cat := range questionnaire.Categories {
		fmt.Printf("INSERT INTO category (id, category_order, name, questionnaire_id) VALUES (nextval('hibernate_sequence'), %d, $$%s$$, 1);\n", cat.Order, cat.Title)
		for _, q := range cat.Questions {
			fmt.Printf("INSERT INTO question (id, question_order, type, name, question_text, description, category_id) VALUES (nextval('hibernate_sequence'), %d, 'SINGLE', $$%s$$, $$%s$$, $$%s$$, (SELECT max(id) FROM category));\n", q.Order, q.Name, q.Question, q.Description)
			for _, o := range q.Options {
				fmt.Printf("INSERT INTO single_option (id, singleoption_order, option, risk, question_id) VALUES (nextval('hibernate_sequence'), %d, $$%s$$, '%s', (SELECT max(id) FROM question));\n", o.Order, o.Option, o.Risk)
			}
		}
	}
}
