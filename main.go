package main

import (
	"fmt"
	"github.com/itzg/go-flagsfiller"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var config struct {
	Bind     string `usage:"The [host:port] to bind, but using port flag is preferred"`
	Port     int    `default:"8080" usage:"The port to bind"`
	Response struct {
		Status           int    `default:"200" usage:"When set, specifies the status code to use in responses"`
		FixedBody        string `usage:"When set, specifies a fixed body to write to the response"`
		FixedContentType string `default:"text/plain" usage:"When FixedBody is set, specifies the content type to set"`
	}
}

func main() {
	err := flagsfiller.Parse(&config, flagsfiller.WithEnv(""))
	if err != nil {
		log.Fatal(err)
	}

	bind := config.Bind
	if bind == "" {
		bind = fmt.Sprintf(":%d", config.Port)
	}

	httpServer := &http.Server{
		Addr:    bind,
		Handler: &debugHandler{},
	}

	log.Printf("Ready for connections at %s", bind)
	log.Fatal(httpServer.ListenAndServe())
}

type debugHandler struct{}

func (h *debugHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	if config.Response.FixedBody != "" {
		resp.Header().Set("Content-Type", config.Response.FixedContentType)
		resp.WriteHeader(config.Response.Status)
		_, err := resp.Write([]byte(config.Response.FixedBody))
		if err != nil {
			log.Printf("Failed to write response body: %s", err)
		}
		return
	}

	resp.Header().Set("Content-Type", "text/html")
	resp.WriteHeader(config.Response.Status)

	fmt.Fprintln(resp, `<html>

<head>
<style>
th {
  text-align: left;
}
</style>
</head>

<body>`)

	log.Printf("INF Handling %s %s", req.Method, req.URL.String())

	startSection(resp, "Request")
	writeField(resp, "Method", req.Method)
	writeField(resp, "URL", req.URL.String())
	writeField(resp, "Host", req.Host)
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
