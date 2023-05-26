package digo

import "errors"

var (
	singletons = make(map[string]any)
	groups     = make(map[string][]any)
)

type CreatorFunc func(...any) any

func RegisterSingleton(id string, object any) {
	singletons[id] = object
}

func RegisterMember(groupId string, object any) {
	members, ok := groups[groupId]
	if !ok {
		members = make([]any, 0)
	}
	members = append(members, object)
	groups[groupId] = members
}

func Members(name string) ([]any, error) {
	group, ok := groups[name]
	if !ok {
		return nil, errors.New("")
	}
	return group, nil
}

func Provide(id string) (any, error) {
	p, ok := singletons[id]
	if ok {
		return p, nil
	}
	return nil, errors.New("")
}
