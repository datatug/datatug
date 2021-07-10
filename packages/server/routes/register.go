package routes

import (
	"github.com/datatug/datatug/packages/server/endpoints"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

type wrapper = func(f http.HandlerFunc) http.HandlerFunc

type Mode = int

const (
	WriteOnlyHandlers Mode = iota
	AllHandlers
)

// RegisterDatatugHandlers registers datatug HTTP handlers
func RegisterDatatugHandlers(
	path string,
	router *httprouter.Router,
	mode Mode,
	wrap wrapper,
	handler endpoints.Handler,
) {
	log.Println("Registering DataTug handlers on path:", path)
	endpoints.SetHandler(handler)
	registerRoutes(path, router, wrap, mode == WriteOnlyHandlers)
}