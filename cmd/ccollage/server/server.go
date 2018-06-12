// Package server is a HTTP server that reads the URL and handles it accordingly
package server

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/shurcooL/httpfs/html/vfstemplate"
	"github.com/valentine/ccollage/cmd/ccollage/svg"
	"github.com/valentine/ccollage/internal/client/github"
)

var serverPort int

func ghHandler(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Path
	q := r.URL.Query()

	su := strings.Split(u, "/")

	var buf bytes.Buffer

	switch su[len(su)-1] {
	case "contributors.svg":
		// making sure that github.com is in the URL
		if strings.ToLower(su[1]) == "github.com" {
			var owner = su[2]
			var repo = su[3]

			c, err := github.GetAllContributors(owner, repo)
			if err != nil {
				buf = parseTemplate(fmt.Sprintln(err))
				w.WriteHeader(http.StatusBadGateway) // 502
				w.Write(buf.Bytes())
				return
			}

			width, padding, maxWidth := parseQueries(q)

			buf = svg.BuildCollage(c, width, padding, maxWidth)
			w.Header().Set("Content-Type", "image/svg+xml")
		} else {
			buf = parseTemplate("Please provide a valid URL.")
			w.WriteHeader(http.StatusBadRequest) // 400
		}
	default:
		buf = parseTemplate("Please add <code>contributors.svg</code> to the end of the URL.")
		w.WriteHeader(http.StatusNotExtended) // 510
		// log.Printf("ERROR: last part of the path needs to be a filename.")
	}
	w.Write(buf.Bytes())
}

// TO-DO
func glHandler(w http.ResponseWriter, r *http.Request) {
}

func readmeHandler(w http.ResponseWriter, r *http.Request) {
	var buf bytes.Buffer

	u := r.URL.Path
	su := strings.Split(u, "/")

	if len(su) > 2 {
		buf = parseTemplate("Please provide a valid URL.")
	} else {
		buf = parseTemplate()
	}

	w.Write(buf.Bytes())
}

func rateLimitHandler(w http.ResponseWriter, r *http.Request) {
	rateTotal, rateLeft, resetTime, err := github.GetRateLimit()
	if err != nil {
		buf := parseTemplate(fmt.Sprintln(err))
		w.WriteHeader(http.StatusBadGateway) // 502
		w.Write(buf.Bytes())
		return
	}

	rt := time.Unix(int64(resetTime), 0)
	diff := rt.Sub(time.Now())
	difft := diff.Truncate(time.Second)

	buf := parseTemplate(fmt.Sprintf("API limit remaining: %d/%d; %+v before reset", rateLeft, rateTotal, difft))
	w.Write(buf.Bytes())
}

func parseTemplate(flashMsg ...string) (output bytes.Buffer) {
	type Messages struct {
		FlashMessage template.HTML
	}

	// initialise the slice if no values were passed into the function
	if flashMsg == nil {
		flashMsg = append(flashMsg, "")
	}

	msg := Messages{template.HTML(flashMsg[0])}

	tmpl := template.Must(vfstemplate.ParseFiles(templates, nil, "readme.html"))

	err := tmpl.Execute(&output, msg)
	if err != nil {
		log.Printf("ERROR: HTML template could not be executed:\n%+v", err)
	}

	return output
}

func parseQueries(q url.Values) (width int, padding int, maxWidth int) {
	var err error
	width, padding, maxWidth = 80, 5, 800

	for k, v := range q {
		switch k {
		case "size", "s":
			width, err = strconv.Atoi(v[0])
			if err != nil {
				log.Printf("ERROR:%v\n", err)
			}
		case "padding", "p":
			padding, err = strconv.Atoi(v[0])
			if err != nil {
				log.Printf("ERROR:%v\n", err)
			}
		case "width", "w":
			maxWidth, err = strconv.Atoi(v[0])
			if err != nil {
				log.Printf("ERROR:%v\n", err)
			}
		}
	}
	return width, padding, maxWidth
}

func init() {
	flag.IntVar(&serverPort, "port", 8080, "Port to listen on.")
}

// Serve starts the web server
func Serve() {
	flag.Parse()

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", serverPort),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	http.HandleFunc("/github.com/", ghHandler)
	http.HandleFunc("/ratelimit", rateLimitHandler)
	//	http.HandleFunc("/gitlab.com/", glHandler)
	http.HandleFunc("/", readmeHandler)
	s.ListenAndServe()
}
