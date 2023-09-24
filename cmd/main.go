package main

import "deferredMessage/internal"

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {

	app, err := internal.NewApp("./.env")
	if err != nil {
		panic(err)
	}
	err = app.Run()
	if err != nil {
		panic(err)
	}
}
