package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func dump(w http.ResponseWriter, req *http.Request) {
	/*
		dump, err := httputil.DumpRequest(req, true)

		if err != nil {
			http.Error(w, fmt.Sprint(err), http.StatusInternalServerError)
			return
		}

		log.Printf("Dump:\n[%q]\n", dump)
	*/
	fmt.Printf("Method %v\n", req.Method)
	fmt.Printf("Host   %v\n", req.Host)
	fmt.Printf("URI    %v\n", req.URL.RequestURI())

	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Printf("Header %v: %v\n", name, h)
		}
	}

	for key, values := range req.URL.Query() {
		for _, value := range values {
			fmt.Printf("Query %v=%v\n", key, value)
		}
	}

	req.ParseForm()

	for k, values := range req.Form {
		for _, value := range values {
			fmt.Printf("Form %v=%v\n", k, value)
		}
	}

	fmt.Println("-----")

	fmt.Fprintf(w, "OK from http_server")
}

func main() {

	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
	time.LoadLocation("Asia/Tokyo")

	l := log.New(os.Stdout, "", 0)
	l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " ")

	http.HandleFunc("/", dump)
	http.HandleFunc("/dump", dump)

	l.Println("http_server starts ...")
	http.ListenAndServe(":8080", nil)
}
