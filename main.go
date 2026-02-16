package main

import (
	"log"
	"todoer/config"
	"todoer/server"
)

func main() {
	/* Error handler */
	defer func() {
		if recovered := recover(); recovered != nil {
			log.Println("Crashed:", recovered)
		}
		log.Println("Bye")
	}()
	config.Load()
	log.Printf("Starting server on port %s", config.Port)
	server.Create(config.Port)
}
