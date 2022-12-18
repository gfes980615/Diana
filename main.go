package main

import (
	"fmt"

	"github.com/gfes980615/Diana/injection"
	"github.com/gfes980615/Diana/injection/controller"

	_ "github.com/gfes980615/Diana/service"
	_ "github.com/gfes980615/Diana/transport/http/controller"
)

// func main() {
// 	defer func() {
// 		log.Error("Server shutdown...")
// 		if err := recover(); err != nil {
// 			log.Errorf("error: %v", err)
// 		}
// 	}()

// 	ErrExit(server.Run())
// }

func main() {
	if err := injection.InitInject(); err != nil {
		fmt.Println(err)
		return
	}

	controller.InitController()

}

func ErrExit(err error) {
	if err != nil {
		// log.Error(err)
	}
}
