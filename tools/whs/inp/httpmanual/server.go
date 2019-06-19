package httpmanual

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strconv"
	"strings"
)

type firstLine struct {
	Type string
	Path string
	Version string
}

func parseFirstLine(line string) (*firstLine, error) {
	parts := strings.Split(line, " ")
	if len(parts) != 3 {
		return nil, fmt.Errorf("first line malformed: %#v", line)
	}

	return &firstLine{
		Type:    parts[0],
		Path:    parts[1],
		Version: parts[2],
	}, nil
}

func parseHeaders(scanner *bufio.Scanner) map[string]string {
	headers := map[string]string{}

	for scanner.Scan() {
		line := scanner.Text()
		if line == ""{
			fmt.Println("found end of line")
			break
		}

		parts := strings.SplitN(line, ": ", 2)
		//fmt.Printf("line:%#v\n", line)
		//fmt.Println("parts:", parts)

		headers[parts[0]] = parts[1]
	}

	return headers
}

func writeHeader(conn net.Conn, key string, value string) {
	fmt.Fprintf(conn, fmt.Sprintf("%s: %s\r\n", key, value))
}

const endOfLine = "\r\n"
var defaultBody = []byte("<html><head><title>hello world!</title></head><body><h1>Yay</h1></body></html>")

type HandlerFunc = func(request *Request, resp *responseWriter)

type Request struct {
	Method string
	Path string
	Headers map[string]string
	Body string
}

type entry struct {
	path string
	handler HandlerFunc
}

var handlers = []entry{
	{path: "/", handler: indexHandler},
	{path: "/index.html", handler: indexHandler},
	{path: ".*", handler: notFoundHandler},
}

func handleConnection(conn net.Conn) {
	//io.Copy(conn, conn)

	// Shutdown connection following HTTP/1.1 serial.
	defer func() {
		if err := conn.Close(); err != nil {
			fmt.Printf("failed to close connection to %v\n: %v", conn.RemoteAddr().String(), err)
		}
		fmt.Printf("Closed connection to %v\n", conn.RemoteAddr().String())
	}()

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		fmt.Println("no first line")
		return
	}

	firstLine, err := parseFirstLine(scanner.Text())
	if err != nil {
		fmt.Println("failed to parse first line:", err)
		return
	}

	fmt.Printf("first line: %#v\n", firstLine)

	// Parse Headers.
	headers := parseHeaders(scanner)

	// Print headers.
	fmt.Println()
	fmt.Println("Headers:")
	for k, v := range headers {
		fmt.Printf("%s: %s\n", k, v)
	}
	fmt.Println()

	req := &Request{
		Method:  firstLine.Type,
		Path:    firstLine.Path,
		Headers: headers,
	}

	contentLengthString, hasContentLength := headers["Content-Length"]
	if hasContentLength {
		fmt.Println("content length:", contentLengthString)

		contentLength, err := strconv.Atoi(contentLengthString)
		log.Fatalf("failed to convert content length %s to int: %v", contentLengthString, err)

		var buf bytes.Buffer
		gesWritten := int64(0)
		for gesWritten < int64(contentLength) {
			written, err := io.Copy(&buf, conn)
			gesWritten += written

			if err != nil {
				if err != io.EOF {
					log.Fatalf("failed to read request body: %v", err)
				}
				if gesWritten != int64(contentLength) {
					log.Fatalf("expected to read %d from body but read %d instead", contentLength, gesWritten)
				}
			}
		}

		req.Body = buf.String()
	}

	fmt.Println("Scanning errors:")
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

	for _, e := range handlers {
		path := e.path
		handlerFunc := e.handler
		pattern := fmt.Sprintf("^%s$", path)
		matched, err := regexp.MatchString(pattern, req.Path)
		if err != nil {
			log.Fatalf("failed to parse path regex: %v\n", err)
		}
		if matched {
			resp := &responseWriter{
				Conn: conn,
			}
			defer resp.Close()

			fmt.Println("found handler")
			handlerFunc(req, resp)
			break
		}
	}

	/*
	if value, hasAcceptEncoding := headers["Accept-Encoding"]; hasAcceptEncoding && value == "gzip" {
		writeHeader(conn, "Content-Encoding", "gzip")

		fmt.Fprintf(conn, "\r\n")

		gzw := gzip.NewWriter(conn)
		gzw.Write(defaultBody)
		gzw.Flush()
	} else {
		fmt.Fprintf(conn, "\r\n")

		fmt.Fprintf(conn, string(defaultBody))
	}
	*/




}