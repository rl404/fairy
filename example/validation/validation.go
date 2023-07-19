package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/rl404/fairy/validation"
)

type user struct {
	Name    string `validate:"required" mod:"trim"`
	Age     int    `validate:"gt=10"`
	Magic   string `mod:"magic"`
	Skill   string `mod:"skill=expert"`
	Country string `validate:"contain=konoha" mod:"lcase"`
}

func main() {
	// Init validator.
	v := validation.New(true)

	// Register custom modifier.
	v.RegisterModifier("magic", func(in string, _ ...string) string {
		return in + " magic"
	})
	v.RegisterModifier("skill", func(in string, param ...string) string {
		return fmt.Sprintf("%s (%s)", in, param[0])
	})

	// Register custom validator.
	v.RegisterValidator("contain", func(value interface{}, param ...string) bool {
		return strings.Contains(value.(string), param[0])
	})

	// Register custom error message handler.
	v.RegisterValidatorError("gt", func(field string, param ...string) error {
		return fmt.Errorf("field %s must be greater than %s", field, param[0])
	})
	v.RegisterValidatorError("contain", func(field string, param ...string) error {
		return fmt.Errorf("field %s must contain %s", field, param[0])
	})

	// Sample 'dirty' data.
	naruto := user{
		Name:    "  Naruto ",
		Age:     15,
		Magic:   "ninja",
		Skill:   "jump",
		Country: "Konohagakure",
	}

	// Validate struct fields.
	if err := v.Validate(&naruto); err != nil {
		panic(err)
	}

	j, _ := json.MarshalIndent(naruto, "", "  ")
	fmt.Println(string(j))
}
