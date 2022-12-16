package gogo

import (
	"context"
	"encoding/json"
	"os"
	"testing"
)

func TestYAML(t *testing.T) {
	b := `
---
data:
  path: "."
  files: !!shell "ls {{.path}}"
  
env:
  - FIRST=foo
  - LAST=bar

tasks:
  - name:
    deps:
      - task:  task_name
        extra:
    data:
    env:
    skip: "none | duplicate | once"
    steps:
      - task: task_name
        extra:
      - shell:
    range: 
`

	ctx := context.Background()

	gg := NewGogo()
	if err := gg.LoadGlobal(ctx, []byte(b)); err != nil {
		t.Fatal(err)
	}

	data := gg.GetData()
	env := gg.GetEnv()

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	enc.Encode(data)
	enc.Encode(env)
}
