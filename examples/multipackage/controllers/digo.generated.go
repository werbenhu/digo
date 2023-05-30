
//
// This file is generated by digen. Run 'digen' to regenerate.
//
// You can install this tool by running `go install github.com/werbenhu/di/digen`.
// For more details, please refer to https://github.com/werbenhu/di. 
//
package controllers

import "github.com/werbenhu/digo"

// init_main_role_name registers the singleton object with ID main.role.name into the DI object manager
// Now you can retrieve the singleton object by using `obj, err := di.Provide("main.role.name")`.
// The obj obtained from the above code is of type `any`.
// You will need to forcefully cast the obj to its corresponding actual object type.
func init_main_role_name() {
	main_role_name_obj := NewRoleName()
	digo.RegisterSingleton("main.role.name", main_role_name_obj)
}

// init_main_user_name registers the singleton object with ID main.user.name into the DI object manager
// Now you can retrieve the singleton object by using `obj, err := di.Provide("main.user.name")`.
// The obj obtained from the above code is of type `any`.
// You will need to forcefully cast the obj to its corresponding actual object type.
func init_main_user_name() {
	main_user_name_obj := NewUserName()
	digo.RegisterSingleton("main.user.name", main_user_name_obj)
}

// Add a member object to group: group.controllers
// Now you can retrieve the group's member objects by using `objs, err := di.Members("group.controllers")`.
// The objs obtained from the above code are of type `[]any`.
// You will need to forcefully cast the objs to their corresponding actual object types.
func group_group_controllers_NewRoleController() {
	name_obj, err := digo.Provide("main.role.name")
	if err != nil {
		panic(err)
	}
	name := name_obj.(string)
	member := NewRoleController(name)
	digo.RegisterMember("group.controllers", member)
}

// Add a member object to group: group.controllers
// Now you can retrieve the group's member objects by using `objs, err := di.Members("group.controllers")`.
// The objs obtained from the above code are of type `[]any`.
// You will need to forcefully cast the objs to their corresponding actual object types.
func group_group_controllers_NewUserController() {
	name_obj, err := digo.Provide("main.user.name")
	if err != nil {
		panic(err)
	}
	name := name_obj.(string)
	member := NewUserController(name)
	digo.RegisterMember("group.controllers", member)
}

// init registers all providers in the current package into the DI object manager.
func init() {
	init_main_role_name()
	init_main_user_name()
	group_group_controllers_NewRoleController()
	group_group_controllers_NewUserController()
}