package controller

import (
	"reflect"

	"github.com/gfes980615/Diana/injection"

	"github.com/gin-gonic/gin"
)

var ControllerDictionary []reflect.Value

func InitController() (*gin.Engine, error) {
	if err := InitControllerDictionary(); err != nil {
		return nil, err
	}

	return InitControllerImpl()
}

func InitControllerDictionary() error {
	controller, err := injection.ParseSuffix("Controller")
	if err != nil {
		return err
	}

	for _, c := range controller {
		ControllerDictionary = append(ControllerDictionary, c.MethodByName("SetupRouter"))
	}

	return nil
}

func InitControllerImpl() (*gin.Engine, error) {
	router := gin.New()

	for _, control := range ControllerDictionary {
		params := make([]reflect.Value, 1)
		params[0] = reflect.ValueOf(router)
		control.Call(params)
	}

	return router, nil
}
