// SPDX-License-Identifier: MIT
// SPDX-FileCopyrightText: 2023 werbenhu
// SPDX-FileContributor: werbenhu

package digo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRegisterSingleton(t *testing.T) {
	id := "singleton"
	obj := "singleton object"

	RegisterSingleton(id, obj)

	// Check if the singleton object is registered successfully
	result, err := Provide(id)
	assert.NoError(t, err)
	assert.Equal(t, obj, result)
}

func TestRegisterMember(t *testing.T) {
	groupID := "group"
	obj1 := "object 1"
	obj2 := "object 2"

	RegisterMember(groupID, obj1)
	RegisterMember(groupID, obj2)

	// Check if the members are registered successfully
	members, err := Members(groupID)
	assert.NoError(t, err)
	assert.Equal(t, []any{obj1, obj2}, members)
}

func TestMembers_NonexistentGroup(t *testing.T) {
	groupID := "nonexistent"

	// Check if the error is returned for a nonexistent group
	_, err := Members(groupID)
	assert.Error(t, err)
	assert.EqualError(t, err, "group not found")
}

func TestProvide(t *testing.T) {
	id := "singleton"
	obj := "singleton object"

	RegisterSingleton(id, obj)

	// Check if the singleton object is provided successfully
	result, err := Provide(id)
	assert.NoError(t, err)
	assert.Equal(t, obj, result)
}

func TestProvide_NonexistentObject(t *testing.T) {
	id := "nonexistent"

	// Check if the error is returned for a nonexistent object
	_, err := Provide(id)
	assert.Error(t, err)
	assert.EqualError(t, err, "object not found")
}
