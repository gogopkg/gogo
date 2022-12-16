package container

import "container/list"

type TodoItem struct {
	Tag   string
	Data  string
	SetFn func(x any)
}

type TodoList struct {
	list.List
}

func (list *TodoList) VisitAll(fn func(item *TodoItem) (err error)) error {
	for e := list.Front(); e != nil; e = e.Next() {
		if err := fn(e.Value.(*TodoItem)); err != nil {
			return err
		}
	}
	return nil
}
