package container

import (
	"reflect"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v3"
)

func valueGetFn(v reflect.Value) func() any {
	return func() any {
		if v.IsValid() {
			return v.Interface()
		}
		return nil
	}
}

func sliceGetFn(arr []any, idx int) func() any {
	return func() any {
		if len(arr) > idx {
			return arr[idx]
		}
		return nil
	}
}

func mapGetFn(m map[string]any, k string) func() any {
	return func() any {
		v, _ := m[k]
		return v
	}
}

func valueSetFn(v reflect.Value) func(x any) {
	return func(x any) {
		if x == nil {
			return
		}
		v.Set(reflect.ValueOf(x))
	}
}

func sliceSetFn(arr []any, idx int) func(x any) {
	return func(x any) {
		if x == nil {
			return
		}
		arr[idx] = x
	}
}

func mapSetFn(m map[string]any, k string) func(x any) {
	return func(x any) {
		if x == nil {
			return
		}
		m[k] = x
	}
}

func TravelTree(out any, n *yaml.Node, tags []string) (list *TodoList, err error) {
	v := reflect.ValueOf(out).Elem()
	list = new(TodoList)

	m := map[string]struct{}{}
	for _, tag := range tags {
		m[tag] = struct{}{}
	}

	err = travelTree(list, m, n, valueGetFn(v), valueSetFn(v))
	if err != nil {
		return nil, err
	}
	return list, nil
}

func travelTree(list *TodoList, tags map[string]struct{}, n *yaml.Node, getFn func() any, setFn func(x any)) (err error) {
	ori := getFn()

	switch n.Kind {
	case yaml.DocumentNode:
		if len(n.Content) == 0 {
			return errors.Errorf("the number os documents is zero")
		}
		// skip document node
		if err := travelTree(list, tags, n.Content[0], getFn, setFn); err != nil {
			return err
		}
	case yaml.SequenceNode:
		var arr []any
		if ori, ok := ori.([]any); ok {
			if len(ori) < len(n.Content) {
				arr = make([]any, len(n.Content))
				copy(arr, ori)
				setFn(arr)
			} else {
				arr = ori
			}
		} else {
			arr = make([]any, len(n.Content))
			setFn(arr)
		}

		for i := 0; i < len(n.Content); i++ {
			vn := n.Content[i]
			if err := travelTree(list, tags, vn, sliceGetFn(arr, i), sliceSetFn(arr, i)); err != nil {
				return err
			}
		}

	case yaml.MappingNode:
		var m map[string]any
		if ori, ok := ori.(map[string]any); ok {
			m = ori
		} else {
			m = make(map[string]any)
			setFn(m)
		}
		for i := 0; i < len(n.Content); i += 2 {
			k := n.Content[i].Value
			vn := n.Content[i+1]
			if err := travelTree(list, tags, vn, mapGetFn(m, k), mapSetFn(m, k)); err != nil {
				return err
			}
		}

	case yaml.ScalarNode:
		if _, ok := tags[n.Tag]; ok {
			list.PushBack(&TodoItem{Tag: n.Tag, Data: n.Value, SetFn: setFn})
		} else {
			switch n.Tag {
			case "!!null", "!!str", "!!bool", "!!int", "!!float":
				var v any
				if err := n.Decode(&v); err != nil {
					return errors.Wrapf(err, "decode yaml scalar")
				}
				setFn(v)
			default:
				return errors.Errorf("unsupported tag: %v", n.Tag)
			}
		}

	default:
		panic(errors.Errorf("unkown yaml kind:%v", n.Kind))
	}

	return nil
}
