package httpmanual

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
)

type responseLine struct {
	Version string
	Code string
	CodeExplanation string
}

func (rl responseLine) String() string {
	return fmt.Sprintf("%s %s %s", rl.Version, rl.Code, rl.CodeExplanation)
}

type responseWriter struct {
	buf *bytes.Buffer
	Conn net.Conn
	rl responseLine
	headers map[string][]string
}

func (c *responseWriter) setupBuffer() {
	if c.buf == nil {
		c.buf = &bytes.Buffer{}
	}
}

func (c *responseWriter) Write(b []byte) (n int, err error){
	c.setupBuffer()
	n, err = c.buf.Write(b)
	return

	/*
	if err != nil {
		log.Fatalf("failed to write to buffer: %v", err)
	}
	if n != len(b) {
		log.Fatalf("failed to write %d to buffer, wrote %d instead", len(b), n)
	}
	*/
}

func (c *responseWriter) WriteString(s string) {
	c.setupBuffer()

	n, err := c.buf.WriteString(s)
	if err != nil {
		log.Fatalf("failed to write to buffer: %v", err)
	}
	if n != len(s) {
		log.Fatalf("failed to write %d to buffer, wrote %d instead", len(s), n)
	}
}

func (c *responseWriter) Close() {
	if _, err := fmt.Fprintf(c.Conn, "%s\r\n", c.rl.String()); err != nil {
		log.Fatalln("failed to write response first line")
		return
	}

	c.setupBuffer()

	// Set standard headers.
	c.SetHeader("Connection", "close")
	c.SetHeader("Content-Type", "text/html; charset=UTF-8")
	c.SetHeader("Content-Length", strconv.Itoa(c.buf.Len()))
	c.SetHeader("Server", "Restlet-Framework/2.0.3")
	c.SetHeader("Vary", "Accept-Charset, Accept-Encoding, Accept-Language, Accept")

	// Write headers.
	for k, valueArray := range c.getHeaders() {
		for _, v := range valueArray {
			writeHeader(c.Conn, k, v)
		}
	}

	if _, err := fmt.Fprintf(c.Conn, "\r\n"); err != nil {
		log.Fatalf("failed to write separator to body: %v", err)
	}

	bodyLength := int64(c.buf.Len())

	n, err := io.Copy(c.Conn, c.buf)
	if err != nil {
		log.Fatalf("failed to copy body %s: %v\n", c.buf.String(), err)
	}
	if n != bodyLength {
		log.Fatalf("failed to write whole buffer %d, wrote only %d", bodyLength, n)
	}
}

func (c *responseWriter) SetResponseLine(rl responseLine) {
	c.rl = rl
}

func (c *responseWriter) SetHeader(key string, value string) {
	c.setupHeaderMap()

	c.headers[key] = []string{value}
}

func (c *responseWriter) AddHeader(key string, value string) {
	c.setupHeaderMap()

	if arr, exists := c.headers[key]; exists {
		c.headers[key] = append(arr, value)
	} else {
		c.headers[key] = []string{value}
	}
}

func (c *responseWriter) getHeaders() map[string][]string {
	c.setupHeaderMap()

	return c.headers
}

func (c *responseWriter) setupHeaderMap() {
	if c.headers == nil {
		c.headers = make(map[string][]string)
	}
}
