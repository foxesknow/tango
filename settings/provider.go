package settings

import "maps"

type Provider interface {
	GetSetting(name string) (value any, found bool)
}

type mapProvider struct {
	data map[string]string
}

// Creates a settings provider from a map.
// A copy of the map is made
func FromMap(data map[string]string) Provider {
	return &mapProvider{
		data: maps.Clone(data),
	}
}

func (m *mapProvider) GetSetting(name string) (value any, found bool) {
	value, found = m.data[name]
	return
}

//////////////////////////////////////////////////////////////////////////////////////////

type cons struct {
	head Provider
	tail Provider
}

// Combines 2 providers into 1.
// When looking up a setting the head is checked first. If not found then the tail is checked
func Cons(head Provider, tail Provider) Provider {
	return &cons{
		head: head,
		tail: tail,
	}
}

func (c *cons) GetSetting(name string) (value any, found bool) {
	value, found = c.head.GetSetting(name)
	if !found {
		value, found = c.tail.GetSetting(name)
	}

	return
}

//////////////////////////////////////////////////////////////////////////////////////////

type nameValue struct {
	name  string
	value string
}

// Creates a provider that has a single setting
func NameValue(name string, value string) Provider {
	return &nameValue{
		name:  name,
		value: value,
	}
}

func (nv *nameValue) GetSetting(name string) (value any, found bool) {
	if nv.name == name {
		value = nv.value
		found = true
	}

	return
}
