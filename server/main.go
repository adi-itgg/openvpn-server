package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	controller2 "server/controller"
	usecase2 "server/usecase"
	"time"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(Middleware)

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	usecase := usecase2.NewUsecase()
	controller := controller2.NewController(usecase)

	r.Get("/api/status", controller.Status)
	r.Post("/api/activate", controller.Activate)

	listenPort := "1195"

	fmt.Printf("Listening on port %s\n", listenPort)
	err := http.ListenAndServe(":"+listenPort, r)
	if err != nil {
		panic(err)
	}
}

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
