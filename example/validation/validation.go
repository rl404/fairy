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
	Country string `validate:"contain=konoha" mod:"lcase"`
}

func main() {
	// Init validator.
	v := validation.New(true)

	// Register custom modifier.
	v.RegisterModifier("magic", func(in string) string {
		return in + " magic"
	})

	// Register custom validator and error message handler.
	v.RegisterValidator("contain",
		func(value interface{}, param ...string) bool {
			return strings.Contains(value.(string), param[0])
		},
		func(field string, param ...string) error {
			return fmt.Errorf("field %s must contain %s", field, param[0])
		},
	)

	// Sample 'dirty' data.
	naruto := user{
		Name:    "  Naruto ",
		Age:     15,
		Magic:   "ninja",
		Country: "Konohagakure",
	}

	// Validate struct fields.
	if err := v.Validate(&naruto); err != nil {
		panic(err)
	}

	j, _ := json.MarshalIndent(naruto, "", "  ")
	fmt.Println(string(j))
}
