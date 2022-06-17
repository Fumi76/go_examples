package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"io/ioutil"
	//"regexp"
    //"html/template"

    //"github.com/microcosm-cc/bluemonday"
    blackfriday "github.com/russross/blackfriday/v2"
    "github.com/Depado/bfchroma" 
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

    content, err := ioutil.ReadFile("template.txt")

    if err != nil {
        log.Fatal(err)
    }

    s := string(content)
	//fmt.Println(s)

/*
    s = "This is some sample code.\n\n```go\n" +
	`func main() {
	fmt.Println("Hi")
}
` + "```"
*/
    
    
    br := blackfriday.NewHTMLRenderer(blackfriday.HTMLRendererParameters{
        Flags: blackfriday.HrefTargetBlank,
    })
    
    renderer := bfchroma.NewRenderer(bfchroma.Extend(br),
    bfchroma.Style("pygments"))
    
    output := blackfriday.Run([]byte(s), 
    blackfriday.WithExtensions(blackfriday.HardLineBreak + blackfriday.Autolink + blackfriday.FencedCode + blackfriday.Tables), 
    blackfriday.WithRenderer(renderer),
    //blackfriday.WithRenderer(bfchroma.NewRenderer()),
    )
    
    
    /*
    output := blackfriday.Run([]byte(s), 
    blackfriday.WithRenderer(bfchroma.NewRenderer(bfchroma.Style("pygments"))))
    */
    
    // html := bluemonday.UGCPolicy().SanitizeBytes(output)

    /*
    p := bluemonday.UGCPolicy()
    p.AllowAttrs("class").Matching(regexp.MustCompile("^language-[a-zA-Z0-9]+$")).OnElements("code")
    html := p.SanitizeBytes(output)
    */
    
    //fmt.Println(string(html))
    fmt.Fprintf(w, "<html><head><link rel=\"stylesheet\" type=\"text/css\" href=\"css/style.css\"/></style></head><body>")
	fmt.Fprintf(w, string(output))
	fmt.Fprintf(w, "</body></html>")
}

func main() {

	time.Local = time.FixedZone("Asia/Tokyo", 9*60*60)
	time.LoadLocation("Asia/Tokyo")

	l := log.New(os.Stdout, "", 0)
	l.SetPrefix(time.Now().Format("2006-01-02 15:04:05") + " ")

	http.HandleFunc("/", dump)
	http.HandleFunc("/dump", dump)

	dir, _ := os.Getwd()
	fmt.Println("current path :" + dir)

    // URIの/css/込みでファイルシステムを探すので取り除く
	http.Handle(
		"/css/",
		http.StripPrefix(
			"/css/",
			http.FileServer(http.Dir("./css")),
		),
	)
	
	l.Println("http_server starts ...")
	http.ListenAndServe(":8080", nil)
}
