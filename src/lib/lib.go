package lib

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Request struct {
	Method  string
	Target  []string
	Version string
	Headers map[string]string
	Body    string
}

func ToRequest(req string) (Request, error) {
	parts := strings.Split(req, "\r\n\r\n")
	requestAndHeaders := strings.Split(parts[0], "\r\n")
	reqLine := strings.Fields(requestAndHeaders[0])
	headersSlice := requestAndHeaders[1:]
	headers, err := toHeaders(headersSlice)
	if err != nil {
		return Request{}, err
	}
	fmt.Println("> Request:", strings.Join(reqLine, ", "))
	fmt.Println("> Headers:", headers)
	fmt.Println("> Body:", parts[1])
	reqObj := Request{
		Method:  reqLine[0],
		Target:  toPaths(reqLine[1]),
		Version: reqLine[2],
		Headers: headers,
		Body:    parts[1],
	}
	// Check body matches Content-Length header
	body := Text(reqObj.Body)
	if body.Length() > 0 {
		contentLength, err := strconv.Atoi(headers["Content-Length"])
		if err != nil {
			return reqObj, fmt.Errorf("Error parsing Content-Length header: %w", err)
		}
		if body.Length() != contentLength {
			return reqObj, errors.New(fmt.Sprintf(
				"Body length (%d) doesn't match Content-Length (%d) for body: %v",
				body.Length(),
				contentLength,
				body,
			))
		}
	}
	return reqObj, nil
}

func toPaths(path string) []string {
	paths := make([]string, 0)
	for _, p := range strings.Split(path, "/") {
		if p != "" {
			paths = append(paths, p)
		}
	}
	return paths
}

type Status int
type Body interface{}

func toHeaders(hSlice []string) (map[string]string, error) {
	headers := make(map[string]string)
	for _, h := range hSlice {
		hParts := strings.Split(h, ": ")
		if len(hParts) != 2 {
			return headers, errors.New("Expected split headers string to have 2 parts")
		}
		headers[hParts[0]] = hParts[1]
	}
	return headers, nil
}

func encodeGzip(body Body) ([]byte, error) {
	var toWrite []byte
	switch v := body.(type) {
	case string:
		toWrite = []byte(v)
	case []byte:
		toWrite = v
	}
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)
	_, err := zw.Write(toWrite)
	if err != nil {
		return nil, err
	}
	if err := zw.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func ToResponse(req Request, status Status, body Body) []byte {
	statusLine := fmt.Sprintf("HTTP/1.1 %s", getStatusMsg(status))
	// TODO: Use a map for headers instead of a slice
	headers := make([]string, 0)
	resBody := make([]byte, 0)
	// Handle encoding response body
	if req.Headers["Accept-Encoding"] != "" {
		encodings := strings.Split(req.Headers["Accept-Encoding"], ", ")
		for _, encoding := range encodings {
			switch encoding {
			case "gzip":
				headers = append(headers, "Content-Encoding: gzip")
				compressed, err := encodeGzip(body)
				if err != nil {
					fmt.Println("Error compressing respones body:", err.Error())
					statusLine = fmt.Sprintf("HTTP/1.1 %s", getStatusMsg(500))
					break
				}
				body = compressed
				break
			}
		}
	}
	// Append body if we have one
	if body != nil {
		switch v := body.(type) {
		case string:
			resBody = append(resBody, []byte(v)...)
			headers = append(headers,
				"Content-Type: text/plain",
				fmt.Sprintf("Content-Length: %d", len(v)),
			)
		case []byte:
			resBody = append(resBody, v...)
			headers = append(headers,
				"Content-Type: application/octet-stream",
				fmt.Sprintf("Content-Length: %d", len(v)),
			)
		default:
			fmt.Println(errors.New(fmt.Sprintf("Unknown body type: %T", v)))
		}
	}
	// Double \r\n in middle to mark end of last header
	// and mark end of the headers section
	res := []byte(fmt.Sprintf("%s\r\n%s\r\n\r\n", statusLine, strings.Join(headers, "\r\n")))
	res = append(res, resBody...)
	fmt.Printf("< Response:\n%s\n", res)
	return res
}

func ToArgsMap(args []string) map[string]string {
	argsMap := make(map[string]string)
	for i := 0; i < len(args); i += 2 {
		argsMap[args[i]] = args[i+1]
	}
	return argsMap
}

func getStatusMsg(status Status) string {
	switch status {
	case 200:
		return "200 OK"
	case 201:
		return "201 Created"
	case 400:
		return "400 Bad Request"
	case 404:
		return "404 Not Found"
	case 500:
		return "500 Internal Server Error"
	default:
		return "500 Internal Server Error"
	}
}

type Text string

func (t Text) Length() int {
	l := 0
	for range t {
		l++
	}
	return l
}
