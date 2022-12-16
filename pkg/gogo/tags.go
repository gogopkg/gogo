package gogo

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/pkg/errors"

	"github.com/gogopkg/gogo/pkg/internal/container"
	"github.com/gogopkg/gogo/pkg/internal/exec"
	"github.com/gogopkg/gogo/pkg/internal/template"
)

/*
- Parse "Data"
  - Render Template
  - Execute Shell
  - Decode Json
- Parse "Env"
  - Render Template
  - Execute Shell
*/

func (gg *Gogo) ProcessTags(ctx context.Context, todos *container.TodoList, data any) (err error) {
	// trim
	err = todos.VisitAll(func(item *container.TodoItem) (err error) {
		item.Tag = strings.Trim(item.Tag, "!:")
		if strings.HasPrefix(item.Tag, "str") {
			item.Tag = strings.TrimLeft(item.Tag[3:], ":")
		}
		return nil
	})
	if err != nil {
		return err
	}

	// template
	err = todos.VisitAll(func(item *container.TodoItem) (err error) {
		if strings.HasPrefix(item.Tag, "raw") {
			item.Tag = strings.TrimLeft(item.Tag[3:], ":")
		} else {
			s, err := template.Render(item.Data, data)
			if err != nil {
				return err
			}
			item.Data = s
			if item.Tag == "" {
				item.SetFn(item.Data)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// shell
	err = todos.VisitAll(func(item *container.TodoItem) (err error) {
		if strings.HasPrefix(item.Tag, "shell") {
			stdout, stderr, err := exec.Exec(ctx, nil, "sh", "-c", item.Data)
			if err != nil {
				if stdout != "" {
					log.Println("stdout:", stdout)
				}
				if stderr != "" {
					log.Println("stderr:", stderr)
				}
				return err
			}
			item.Data = stdout
			item.Tag = strings.TrimLeft(item.Tag[5:], ":")
			if item.Tag == "" {
				item.SetFn(item.Data)
			}
		}

		return nil
	})
	if err != nil {
		return err
	}

	// json
	err = todos.VisitAll(func(item *container.TodoItem) (err error) {
		if strings.HasPrefix(item.Tag, "json") {
			b := []byte(item.Data)
			var v any
			if err := json.Unmarshal(b, &v); err != nil {
				return errors.Wrap(err, "unmarshal json")
			}
			item.Tag = strings.TrimLeft(item.Tag[4:], ":")
			if item.Tag == "" {
				item.SetFn(v)
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	// error
	err = todos.VisitAll(func(item *container.TodoItem) (err error) {
		if item.Tag != "" {
			return errors.Errorf("unsupported tag: %v", item.Tag)
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}
