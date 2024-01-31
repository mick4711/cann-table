// go server calls handlers for home page routes
package main

import (
	"html/template"
	"log"
	"net/http"

	"github.com/mick4711/moh/cann"
	"github.com/mick4711/moh/fpl"
	"github.com/mick4711/moh/huxley"
)

// main entry point - http server
func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/cann", cannHandler)
	http.HandleFunc("/huxley", huxleyHandler)
	http.HandleFunc("/fpl", fplHandler)

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// log request details
func logRequest(req *http.Request) {
	log.Println("User-Agent:", req.Header["User-Agent"])
	log.Println("Cf-Ipcountry:", req.Header["Cf-Ipcountry"])
	log.Println("Cf-Connecting-Ip:", req.Header["Cf-Connecting-Ip"])
	log.Println("Sec-Ch-Ua-Platform:", req.Header["Sec-Ch-Ua-Platform"])
	log.Println("Sec-Ch-Ua:", req.Header["Sec-Ch-Ua"])
}

// displays landing page with links to other pages
func homeHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Request on /")
	logRequest(req)

	// generate html output
	homeTemplate := template.Must(template.ParseFiles("HomeTemplate.html"))
	if err := homeTemplate.Execute(w, nil); err != nil {
		log.Fatal(err)
	}
}

// displays Huxley's personal details
func huxleyHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Request on /huxley")
	logRequest(req)

	// generate html output
	huxley.DogStats(w, req)
}

// displays FPL league table
func fplHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Request on /fpl")
	logRequest(req)

	// get json for consumption by vercel app
	fpl.FplPoints(w, req)
}

// fetches the standard table standings, generates and outputs the Cann table
func cannHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("Request on /cann")
	logRequest(req)

	cann.GenerateTable(w, req)
}
