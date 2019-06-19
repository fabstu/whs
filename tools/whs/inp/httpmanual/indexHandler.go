package httpmanual

import (
	"encoding/hex"
	"fmt"
	"html/template"
	"log"
	"time"

	"golang.org/x/crypto/sha3"
)

func Error404(req *Request, resp *responseWriter) {
	const errorContent = "<html><title>Error 404</title><body><h1>Path {{.Path}} does not exist for method {{.Method}}.</body></html>"
	var errorTemplate = template.Must(template.New("").Parse(errorContent))

	if err := errorTemplate.Execute(resp, req); err != nil {
		message := fmt.Sprintf("failed to show error due to template error: %v\n", err)
		log.Println(message)
		Error("501", message, resp)
		return
	}

	Error("404", "page not found", resp)
}

func Error(statusCode string, statusMessage string, resp *responseWriter) {
	resp.SetResponseLine(responseLine{
		Version:         "HTTP/1.1",
		Code:            statusCode,
		CodeExplanation: statusMessage,
	})
}

var content = `<html><head></head><body><h1>Willkommen</h1><p><a href="http://localhost:8081/index.html">Back here</a> </p></body></html>`
var etag string
var lastModified time.Time
var expires time.Time
var location *time.Location

func indexHandler(req *Request, resp *responseWriter) {
	if req.Method != "GET" {
		Error404(req, resp)
		return
	}

	if etag == "" {
		hash := sha3.Sum256([]byte(content))
		etag = hex.EncodeToString(hash[:])

		// Last-Modified: <day-name>, <day> <month> <year> <hour>:<minute>:<second> GMT
		lastModified = time.Now().UTC()
		expires = lastModified.Add(time.Second * 60)
		loc, err := time.LoadLocation("Europe/London")
		if err != nil {
			message := fmt.Sprintf("failed to parse location: %v", err)
			log.Println(message)
			Error("501", message, resp)
			return

		}
		location = loc
	}

	if reqETag, hasEtag := req.Headers["If-None-Match"]; hasEtag {
		if reqETag == etag {
			fmt.Println("304 due to ETag")
			Error("304", "you have it already", resp)
			return
		}
		fmt.Println("wrong req etag")
	}
	if reqModifiedSince, hasModifiedSince := req.Headers["If-Modified-Since"]; hasModifiedSince {
		modifiedSinceTime, err := time.Parse(time.RFC1123, reqModifiedSince)
		if err != nil {
			Error("501", fmt.Sprintf("failed to parse modified since %s: %v", reqModifiedSince, err), resp)
			return
		}

		if lastModified.Before(modifiedSinceTime) {
			fmt.Println("304 due to If-Modified-Since")
			Error("304", "you have it already", resp)
			return
		}
		fmt.Println("Modifier after", modifiedSinceTime.Format(time.RFC3339))
	}

	fmt.Println("sending fresh")

	resp.SetResponseLine(responseLine{
		Version:         "HTTP/1.1",
		Code:            "200",
		CodeExplanation: "OK",
	})

	resp.SetHeader("ETag", etag)
	resp.SetHeader("Last-Modified", lastModified.In(location).Format(time.RFC1123))

	// Date: <day-name>, <day> <month> <year> <hour>:<minute>:<second> GMT
	resp.SetHeader("Date", time.Now().In(location).Format(time.RFC1123))

	// Expires: Wed, 21 Oct 2015 07:28:00 GMT
	resp.SetHeader("Expires", expires.In(location).Format(time.RFC1123))

	//resp.SetHeader("Content-Location", "http://localhost:8081")

	//resp.AddHeader("Cache-Control", "public, max-age=86400")

	resp.WriteString(content)
}

