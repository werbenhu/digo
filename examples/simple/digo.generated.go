
//
// This file is generated by digogen. Run 'digogen' to regenerate.
//
// You can install this tool by running `go install github.com/werbenhu/digo/digogen`.
// For more details, please refer to https://github.com/werbenhu/digo. 
//
package main

import "github.com/werbenhu/digo"

// init_main_db_url registers the singleton object with ID main.db.url into the DI object manager
// Now you can retrieve the singleton object by using `obj, err := di.Provide("main.db.url")`.
// The obj obtained from the above code is of type `any`.
// You will need to forcefully cast the obj to its corresponding actual object type.
func init_main_db_url() {
	main_db_url_obj := NewDbUrl()
	digo.RegisterSingleton("main.db.url", main_db_url_obj)
}

// init_main_db registers the singleton object with ID main.db into the DI object manager
// Now you can retrieve the singleton object by using `obj, err := di.Provide("main.db")`.
// The obj obtained from the above code is of type `any`.
// You will need to forcefully cast the obj to its corresponding actual object type.
func init_main_db() {
	url_obj, err := digo.Provide("main.db.url")
	if err != nil {
		panic(err)
	}
	url := url_obj.(string)
	main_db_obj := NewDb(url)
	digo.RegisterSingleton("main.db", main_db_obj)
}

// init_main_redis registers the singleton object with ID main.redis into the DI object manager
// Now you can retrieve the singleton object by using `obj, err := di.Provide("main.redis")`.
// The obj obtained from the above code is of type `any`.
// You will need to forcefully cast the obj to its corresponding actual object type.
func init_main_redis() {
	main_redis_obj := NewRedis()
	digo.RegisterSingleton("main.redis", main_redis_obj)
}

// init_main_app registers the singleton object with ID main.app into the DI object manager
// Now you can retrieve the singleton object by using `obj, err := di.Provide("main.app")`.
// The obj obtained from the above code is of type `any`.
// You will need to forcefully cast the obj to its corresponding actual object type.
func init_main_app() {
	db_obj, err := digo.Provide("main.db")
	if err != nil {
		panic(err)
	}
	db := db_obj.(*Db)
	redis_obj, err := digo.Provide("main.redis")
	if err != nil {
		panic(err)
	}
	redis := redis_obj.(*Redis)
	main_app_obj := NewApp(db, redis)
	digo.RegisterSingleton("main.app", main_app_obj)
}

// init registers all providers in the current package into the DI object manager.
func init() {
	init_main_db_url()
	init_main_db()
	init_main_redis()
	init_main_app()
}
