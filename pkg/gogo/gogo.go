package gogo

import (
	"context"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"

	"github.com/gogopkg/gogo/pkg/internal/container"
)

var DataTags = []string{
	"!!raw",        // nothing
	"!!str",        // render
	"!!json",       // render > json
	"!!shell",      // render > shell
	"!!shell:json", // render > shell > json
}

type Dep struct {
	Task  string    `yaml:"task"`
	Extra yaml.Node `yaml:"data"`
}

type Step struct {
	// Task trigger
	Task  string    `yaml:"task"`
	Extra yaml.Node `yaml:"extra"`
	// Shell trigger
	Shell string `yaml:"shell"`
}

type Task struct {
	Name  string    `yaml:"name"`
	Deps  []Dep     `yaml:"deps"`
	Data  yaml.Node `yaml:"data"`
	Env   []string  `yaml:"env"`
	Skip  string    `yaml:"skip"`
	Steps []*Step   `yaml:"steps"`
	Range string    `yaml:"range"`

	// internal
	data  any
	extra any
}

type Gogo struct {
	Data  yaml.Node `yaml:"data"`
	Env   []string  `yaml:"env"`
	Tasks []*Task   `yaml:"tasks"`

	// internal
	data  any
	tasks map[string]*Task
}

func NewGogo() *Gogo {
	return &Gogo{
		tasks: make(map[string]*Task),
	}
}

func (gg *Gogo) LoadGlobal(ctx context.Context, b []byte) (err error) {
	if err := yaml.Unmarshal(b, &gg); err != nil {
		return errors.Wrap(err, "unmarshal yaml")
	}

	var data any
	todo, err := container.TravelTree(&data, &gg.Data, DataTags)
	if err != nil {
		return errors.Wrap(err, "travel data tree")
	}

	if err := gg.ProcessTags(ctx, todo, data); err != nil {
		return errors.Wrap(err, "process tags")
	}
	gg.data = data

	for _, task := range gg.Tasks {
		if _, ok := gg.tasks[task.Name]; ok {
			return errors.Errorf("task [%v] is duplicated", task.Name)
		}
		gg.tasks[task.Name] = task
	}

	return nil
}

func (gg *Gogo) GetData() (data any) {
	if gg.data == nil {
		panic("data is not yet loaded")
	}
	return gg.data
}

func (gg *Gogo) GetEnv() []string {
	return gg.Env
}

func (gg *Gogo) LoadTask(ctx context.Context, name string) (err error) {
	t, _ := gg.tasks[name]
	if t == nil {
		return nil
	}

	data := CloneData(gg.GetData())
	todo, err := container.TravelTree(&data, &t.Data, DataTags)
	if err != nil {
		return errors.Wrapf(err, "travel Task[%v] data tree", t.Name)
	}

	if err := gg.ProcessTags(ctx, todo, data); err != nil {
		return errors.Wrapf(err, "process Task[%v] tags", t.Name)
	}

	return nil
}

func CloneData(in any) any {
	// TODO
	return in
}
