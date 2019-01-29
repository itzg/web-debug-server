package main

import (
	"flag"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	bind := flag.String("bind", ":8080", "Host:port to bind")
	flag.Parse()

	httpServer := &http.Server{
		Addr:    *bind,
		Handler: &debugHandler{},
	}

	log.Printf("Ready for connections at %s", *bind)
	log.Fatal(httpServer.ListenAndServe())
}

type debugHandler struct{}

func (h *debugHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	resp.Header().Set("Content-Type", "text/html")

	fmt.Fprintln(resp, `<html>

<head>
<style>
th {
  text-align: left;
}
</style>
</head>

<body>`)

	startSection(resp, "Request")
	writeField(resp, "Method", req.Method)
	writeField(resp, "URL", req.URL.String())
	writeField(resp, "Remote address", req.RemoteAddr)
	hostname, err := os.Hostname()
	if err == nil {
		writeField(resp, "Server hostname", hostname)
	}
	endSection(resp)

	startSection(resp, "Headers")
	for key, values := range req.Header {
		for _, value := range values {
			writeField(resp, key, value)
		}
	}
	endSection(resp)

	startSection(resp, "Content")
	bodyBytes, err := ioutil.ReadAll(req.Body)
	if err != nil {
		writeField(resp, "ERROR", err.Error())
	} else {
		writeField(resp, "", html.EscapeString(string(bodyBytes)))
	}
	endSection(resp)

	fmt.Fprintln(resp, "</body></html>")
}

func endSection(resp http.ResponseWriter) {
	fmt.Fprintln(resp, "</table></div>")
}

func startSection(resp http.ResponseWriter, section string) {
	fmt.Fprintf(resp, "<div><h1>%s</h1>\n", section)
	fmt.Fprintln(resp, "<table>")
}

func writeField(resp http.ResponseWriter, field, value string) {
	fmt.Fprintf(resp, "<tr><th>%s</th><td>%s</td></tr>\n", field, value)
}
