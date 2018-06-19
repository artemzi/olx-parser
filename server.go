package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/artemzi/olx-parser/version"
	simplejson "github.com/bitly/go-simplejson"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

// DEFAULTPORT returns default port number
const DEFAULTPORT = "8080"

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `Ok`)
}

func info(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json := simplejson.New()
	json.Set("version", version.RELEASE)
	json.Set("commit", version.COMMIT)
	json.Set("repo", version.REPO)

	payload, err := json.MarshalJSON()
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(payload)
}

func root(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, fmt.Sprintf("PARSER v%s\n", version.RELEASE))
}

func notFound(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	io.WriteString(w, `Not Found`)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info(fmt.Sprintf("%s: %s", r.Method, r.RequestURI))

		next.ServeHTTP(w, r)
	})
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Warn(fmt.Sprintf("Not Found (%s) %s", r.Method, r.RequestURI))

		f(w, r)
	}
}

func init() {
	switch version.STAGE {
	case "dev":
		log.SetLevel(log.DebugLevel)
	case "prod":
		f, err := os.OpenFile(fmt.Sprintf("./storage/logs/%s.log", time.Now().Local().Format("2006-01-02")),
			os.O_APPEND|os.O_WRONLY|os.O_CREATE,
			0755)
		if err != nil {
			log.Error(err)
		}

		log.SetFormatter(&log.JSONFormatter{})
		log.SetOutput(f)
		log.SetLevel(log.InfoLevel)
	}
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout",
		time.Second*15,
		"the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	port := os.Getenv("SERVICE_PORT")
	if len(port) == 0 {
		port = DEFAULTPORT
	}

	r := mux.NewRouter()
	r.HandleFunc("/info", info).Methods("GET")
	r.HandleFunc("/healthz", healthz).Methods("GET")
	r.HandleFunc("/", root).Methods("GET")
	r.NotFoundHandler = http.HandlerFunc(logging(notFound))

	r.Use(loggingMiddleware)

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	go func() {
		log.Debug(fmt.Sprintf("Server started on: http://0.0.0.0:%s", port))
		if err := srv.ListenAndServe(); err != nil {
			log.Error(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Info("shutting down")
	os.Exit(0)
}
