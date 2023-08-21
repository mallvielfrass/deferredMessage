package main

import "deferredMessage/internal"

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
