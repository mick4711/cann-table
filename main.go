// go server calls handlers for home page routes
package main

import (
	"html/template"
	"log"
	"net/http"
	"time"

	"github.com/mick4711/moh/cann"
	"github.com/mick4711/moh/fpl"
	"github.com/mick4711/moh/huxley"
)

const (
	ServerReadTimeout  = 5 * time.Second
	ServerWriteTimeout = 10 * time.Second
)

// main entry point - http server
func main() {
	mux := new(http.ServeMux)
	mux.Handle("GET /{$}", http.HandlerFunc(homeHandler))
	mux.Handle("GET /cann", http.HandlerFunc(cannHandler))
	mux.Handle("GET /huxley", http.HandlerFunc(huxleyHandler))
	mux.Handle("GET /fpl", http.HandlerFunc(fplHandler))

	srv := &http.Server{
		ReadTimeout:  ServerReadTimeout,
		WriteTimeout: ServerWriteTimeout,
		Addr:         ":8080",
		Handler:      mux,
	}

	log.Println("Listening on port 8080")
	log.Fatal(srv.ListenAndServe())
}

// log request details
func logRequest(req *http.Request) {
	if req.RequestURI == "/favicon.ico" {
		return
	}

	log.Printf("\n============ route = [%s]  ===================\n", req.RequestURI)
	log.Println("Method:", req.Method)
	log.Println("User-Agent:", req.Header["User-Agent"])
	log.Println("Cf-Ipcountry:", req.Header["Cf-Ipcountry"])
	log.Println("Cf-Connecting-Ip:", req.Header["Cf-Connecting-Ip"])
	log.Println("Sec-Ch-Ua-Platform:", req.Header["Sec-Ch-Ua-Platform"])
	log.Println("Sec-Ch-Ua:", req.Header["Sec-Ch-Ua"])
}

// displays landing page with links to other pages
func homeHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)

	// generate html output
	homeTemplate := template.Must(template.ParseFiles("HomeTemplate.html"))
	if err := homeTemplate.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

// displays Huxley's personal details
func huxleyHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)

	// generate html output
	huxley.DogStats(w, req)
}

// displays FPL league table
func fplHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)

	// get json for consumption by vercel app
	fpl.Points(w, req)
}

// fetches the standard table standings, generates and outputs the Cann table
func cannHandler(w http.ResponseWriter, req *http.Request) {
	logRequest(req)

	cann.GenerateTable(w, req)
}
