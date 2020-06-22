package main

import (
	"log"

	"github.com/gfes980615/Diana/apis"
)

func main() {

	defer func() {
		if rc := recover(); rc != nil {
			log.Printf("panic:\n%v\n", rc)
		}
	}()

	apis.MainApis() // gin

}
