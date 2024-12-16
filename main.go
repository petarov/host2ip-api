package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

const (
	VERSION           = "1.0"
	HEART             = "\u2764"
	DEFAULT_PORT      = 7029
	CONTENT_TYPE_JSON = "application/json"
)

var (
	// Server listen address
	listenAddress string
	// Server listen port
	listenPort int
	// API key to prevent public access
	apiKey string
)

func init() {
	flag.StringVar(&listenAddress, "addr", "[::1]", "Server listen address")
	flag.IntVar(&listenPort, "port", DEFAULT_PORT, "Server listen port")
	flag.StringVar(&apiKey, "apikey", "", "API key to prevent public access")
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	routes := map[string]string{
		"/":             "Shows this",
		"/lookup/:name": "Fetches a list of IP addresses for the specified FQDN or hostname",
	}
	resp, _ := json.Marshal(routes)
	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.Write(resp)
}

func writeError(w http.ResponseWriter, status int, err error) {
	resp, _ := json.Marshal(map[string]interface{}{
		"code":   status,
		"status": http.StatusText(status),
		"err":    fmt.Sprintf("%v", err),
	})
	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.Write(resp)
}

func handleLookup(w http.ResponseWriter, r *http.Request) {
	key := strings.TrimSpace(r.URL.Query().Get("key"))
	if len(apiKey) > 0 && apiKey != key {
		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, "Unauthorized")
		return
	}

	name := strings.TrimPrefix(strings.TrimPrefix(r.URL.Path, "/lookup"), "/")
	if len(name) == 0 {
		writeError(w, http.StatusBadRequest, fmt.Errorf("hostname not specified"))
		return
	}

	ips, err := net.LookupIP(name)
	if err != nil {
		writeError(w, http.StatusNotFound, err)
		return
	}

	arr := make([]string, 0)

	for _, ip := range ips {
		arr = append(arr, ip.String())
	}

	resp, _ := json.Marshal(map[string]interface{}{
		"addresses": arr,
	})
	w.Header().Set("Content-Type", CONTENT_TYPE_JSON)
	w.Write(resp)
}

func createRouter() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/", handleIndex)
	router.HandleFunc("/lookup", handleLookup)
	router.HandleFunc("/lookup/", handleLookup)
	return router
}

func main() {
	fmt.Printf("%s host2ip API service v%s %s\n\n", HEART, VERSION, HEART)
	flag.Parse()

	fmt.Printf("Serving at %s and port %d ...\n", listenAddress, listenPort)

	apiKey = strings.TrimSpace(apiKey)
	if len(apiKey) == 0 {
		fmt.Println("*** NO API KEY: The server is open to public queries")
	}

	background := context.Background()
	ctx, cancel := signal.NotifyContext(background, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		if err := http.ListenAndServe(fmt.Sprintf("%s:%d", listenAddress, listenPort), createRouter()); err != nil {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		fmt.Printf("server error: %v\n", err)
	case <-ctx.Done():
		fmt.Printf("\n%s Goodbye.\n", HEART)
	}
}
