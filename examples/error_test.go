package examples

import (
	"errors"
	"fmt"
	"testing"
)

type ErrorHandler func(message string) error

func function1(message string) error {
	return errors.New("function1 error")
}

func function2(message string) error {
	return errors.New("function2 error")
}

func (fn ErrorHandler) Execute(message string) {
	if err := fn(message); err != nil {
		fmt.Println("Error occured: ", err)
	}
}

type Application struct {
	handlers map[string]ErrorHandler
}

func (app *Application) Handle(path string, fn ErrorHandler) {
	if len(app.handlers) == 0 {
		app.handlers = make(map[string]ErrorHandler)
	}
	app.handlers[path] = fn
}

func (app *Application) Send(path, message string) error {
	if len(app.handlers) == 0 || app.handlers[path] == nil {
		return errors.New("Handler not found for given path")
	}
	app.handlers[path].Execute(message)
	return nil
}

func TestErrorHandler(t *testing.T) {
	app := &Application{}
	app.Handle("/func1", ErrorHandler(function1))
	app.Handle("/func2", ErrorHandler(function2))

	app.Send("/func2", "send me inventory")
}
