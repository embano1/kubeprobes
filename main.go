package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envpath = "/env"
	healhtz = "/healthz"
	ready   = "/ready"
)

// please excuse using global vars for simplicity here ;)
var started = time.Now()
var lifefail *time.Duration
var beready *time.Duration
var version string

func main() {

	handlers := map[string]func(w http.ResponseWriter, r *http.Request){
		"/healthz":     healthzHandler,
		"/ready":       readyHandler,
		"/info":        infoHandler,
		"/started":     startedHandler,
		"/debug/pprof": noOpHandler,
		// "runtime":  runtimeHandler,
	}

	listen := flag.String("addr", ":8080", "Address and port to listen on")
	lifefail = flag.Duration("l", 60*time.Second, "When to fail the readiness check")
	beready = flag.Duration("r", 10*time.Second, "When to respond with http.StatusOK (ready)")
	flag.Parse()

	// set up root context (http.Shutdown()) and prepare to catch OS signals
	ctx := context.Background()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	defer func() {
		// cancel()
		signal.Stop(c)
	}()

	mux := http.NewServeMux()
	srv := &http.Server{
		Addr:    *listen,
		Handler: mux,
	}

	log.Printf("Starting http server (app version: %v)", version)
	log.Print("Registering handlers:")
	for h, f := range handlers {
		log.Println(h)
		// avoid registering pprof handler twice (panic)
		if h == "/debug/pprof" {
			continue
		}
		mux.HandleFunc(h, f)
	}

	go func() {
		srv.ListenAndServe()
	}()

	// catch OS sigs
	sig := <-c
	log.Printf("Got %v\n", sig)
	log.Println("Attempting graceful shutdown (closing open handlers, etc.)")
	err := srv.Shutdown(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	} else {
		fmt.Println("Done")
	}

}

func noOpHandler(w http.ResponseWriter, r *http.Request) {}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "App version: %s\n", version)

	podname, _ := os.Hostname()
	fmt.Fprintf(w, "Pod name: %s\n", podname)

	node := os.Getenv("NODE")
	fmt.Fprintf(w, "Node name: %s\n", node)

	addrs, _ := net.InterfaceAddrs()
	fmt.Fprintf(w, "IP(s): %s\n", addrs)

	/*for _, e := range os.Environ() {
		fmt.Fprintf(w, "%s\n", e)
	}*/
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	if time.Now().Sub(started).Seconds() > lifefail.Seconds() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failure"))
	} else {
		w.Write([]byte("ok"))
	}
	// time.Sleep(10 * time.Second)
}
func readyHandler(w http.ResponseWriter, r *http.Request) {
	if time.Now().Sub(started).Seconds() >= beready.Seconds() {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ready"))
	} else {
		w.WriteHeader(http.StatusServiceUnavailable)
		w.Write([]byte("failure"))
	}
}

func startedHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	data := (time.Now().Sub(started)).String()
	w.Write([]byte("Running for " + data))
}
