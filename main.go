package main

import (
	"log"

	"github.com/gfes980615/Diana/server"
)

func main() {
	defer func() {
		log.Error("Server shutdown...")
		if err := recover(); err != nil {
			log.Errorf("error: %v", err)
		}
	}()

	ErrExit(server.Run())
}

func ErrExit(err error) {
	if err != nil {
		// log.Error(err)
	}
}
