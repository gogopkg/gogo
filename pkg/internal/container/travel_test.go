package container

import (
	"log"
	"testing"

	"gopkg.in/yaml.v3"
)

// override
var b []byte = []byte(`
---
int: 200
float: 3.14
bool: true
str: "string"
none: none
slice:
  - 
  - 200
  -
  - 
  -
  -
  - 700
`)

func TestTravel(t *testing.T) {

	var n yaml.Node
	if err := yaml.Unmarshal(b, &n); err != nil {
		t.Fatal(err)
	}

	// origin
	var v = map[string]any{
		"foo":   "bar",
		"int":   100,
		"bool":  false,
		"slice": []any{10, 20, 30, 40, 50},
		"str":   false,
	}
	todo, err := TravelTree(&v, &n, []string{"!!todo"})
	if err != nil {
		t.Fatal(err)
	}

	log.Println(v, todo)
}
