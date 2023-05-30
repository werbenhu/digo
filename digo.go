// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package digo

import "errors"

var (
	singletons = make(map[string]any)   // Map to store singleton objects by their IDs.
	groups     = make(map[string][]any) // Map to store groups of objects by their group IDs.
)

// RegisterSingleton registers a singleton object with the provided ID.
func RegisterSingleton(id string, object any) {
	singletons[id] = object
}

// RegisterMember registers a member object with the provided group ID.
func RegisterMember(groupId string, object any) {
	members, ok := groups[groupId]
	if !ok {
		members = make([]any, 0)
	}
	members = append(members, object)
	groups[groupId] = members
}

// Members returns the group of objects associated with the provided group ID.
// It returns an error if the group does not exist.
func Members(name string) ([]any, error) {
	group, ok := groups[name]
	if !ok {
		return nil, errors.New("group not found")
	}
	return group, nil
}

// Provide returns the singleton object associated with the provided ID.
// It returns an error if the object does not exist.
func Provide(id string) (any, error) {
	p, ok := singletons[id]
	if ok {
		return p, nil
	}
	return nil, errors.New("object not found")
}
