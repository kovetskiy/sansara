package main

import (
	"io"
	"log"
	"net/http"
	"sync"
)

type Controller struct {
	request  *http.Request
	response http.ResponseWriter
}

func (controller *Controller) Serve() {
	io.WriteString(controller.response, controller.request.URL.String())
}

type Handler struct {
	pool *sync.Pool
}

func (handler *Handler) init() {
	handler.pool = &sync.Pool{}
	handler.pool.New = func() interface{} {
		return &Controller{}
	}

	for i := 0; i < 50; i++ {
		handler.deposit()
	}
}

func (handler *Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	controller := handler.acquire()

	controller.request = request
	controller.response = response

	controller.Serve()

}

func (handler *Handler) acquire() *Controller {
	controller := handler.pool.Get().(*Controller)
	go handler.deposit()
	return controller
}

func (handler *Handler) deposit() {
	handler.pool.Put(&Controller{})
}

func main() {
	handler := &Handler{}
	handler.init()

	server := http.Server{
		Addr:    "localhost:8080",
		Handler: handler,
	}

	log.Fatal(server.ListenAndServe())
}
