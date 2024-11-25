package converter

import "fmt"

type Jmap = map[string]any

type Jmaps = []Jmap

// JmapsFromAny extract a Jmaps from an "any" value which is actually a Jmap
// or a slice in which elements are also "any" values with a dynamic Jmap value.
func JmapsFromAny(v any) (Jmaps, error) {
	slice, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("items expected as []any but got %T: %#v", v, v)
	}
	items := make(Jmaps, len(slice))
	for i, item := range slice {
		items[i], ok = item.(Jmap)
		if !ok {
			return nil, fmt.Errorf("item [%d] expected as Jmap but got %T: %#v", i, item, item)
		}
	}
	return items, nil
}

func JmapFromAny(v any) (Jmap, error) {
	shortcut, ok := v.(Jmap)
	if ok {
		return shortcut, nil
	}
	slice, err := JmapsFromAny(v)
	if err != nil {
		return nil, err
	}
	if len(slice) != 1 {
		return nil, fmt.Errorf("items expected as convertible to Jmaps of len 1 but got %T: %#v", v, v)
	}
	return slice[0], nil
}
