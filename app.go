package main

import "log"

import "github.com/peaberberian/OscarGoGo/config"
import "github.com/peaberberian/OscarGoGo/routes"

func main() {
	log.Printf("starting application")
	start()
}

func start() {
	var conf, readErr = config.GetConfig()
	if readErr != nil {
		panic(readErr)
	}
	routes.StartServer(conf)
}
