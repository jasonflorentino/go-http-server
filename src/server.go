package main

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/jasonflorentino/go-http-server/src/lib"
)

const PORT int = 4221
const BUF_SIZE int = 1024

var FILE_PATH string

func main() {
	args := lib.ToArgsMap(os.Args[1:])
	fmt.Println("Args:", args)
	FILE_PATH = args["--directory"]

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", PORT))
	if err != nil {
		fmt.Println("Error: Failed to bind to port", PORT, err.Error())
		os.Exit(1)
	}
	fmt.Println("Listening on", PORT)

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			os.Exit(1)
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Handling connection")
	readBuf := make([]byte, BUF_SIZE)
	n, err := conn.Read(readBuf)
	if err != nil {
		fmt.Println("Error reading from connection:", err.Error())
		conn.Write(lib.ToResponse(lib.Request{}, 500, nil))
		return
	}
	fmt.Println("Read", n, "bytes")

	// Important to convert only the section of the
	// byte slice that has been read into. Otherwise
	// downstream usage of the "request body" will
	// have excess length from the unused slice area
	req, err := lib.ToRequest(string(readBuf[:n]))
	if err != nil {
		fmt.Println("Error creating request:", err.Error())
		conn.Write(lib.ToResponse(req, 500, nil))
		return
	}

	switch {
	case req.Target == "/":
		conn.Write(lib.ToResponse(req, 200, nil))
	case strings.HasPrefix(req.Target, "/echo/"):
		status, body := handleEcho(req)
		conn.Write(lib.ToResponse(req, status, body))
	case strings.HasPrefix(req.Target, "/files"):
		status, body := handleFiles(req)
		conn.Write(lib.ToResponse(req, status, body))
	case strings.HasPrefix(req.Target, "/user-agent"):
		status, body := handleUserAgent(req)
		conn.Write(lib.ToResponse(req, status, body))
	default:
		conn.Write(lib.ToResponse(req, 404, nil))
	}
}

type status = lib.Status
type body = lib.Body

func handleEcho(req lib.Request) (status, body) {
	// TODO: Split out path parts
	yell := req.Target[len("/echo/"):]
	return 200, yell
}

func handleFiles(req lib.Request) (status, body) {
	fileName := req.Target[len("/files/"):]
	switch req.Method {
	case "GET":
		dat, err := os.ReadFile(fmt.Sprintf("%s%s", FILE_PATH, fileName))
		if err != nil {
			fmt.Println("Error reading file:", err.Error())
			return 404, nil
		}
		return 200, dat
	case "POST":
		err := os.WriteFile(fmt.Sprintf("%s%s", FILE_PATH, fileName), []byte(req.Body), 0644)
		if err != nil {
			fmt.Println("Error writing file:", err.Error())
			return 404, nil
		}
		return 201, nil
	default:
		return 400, nil
	}
}

func handleUserAgent(req lib.Request) (status, body) {
	userAgent := req.Headers["User-Agent"]
	if userAgent == "" {
		return 400, nil
	}
	return 200, userAgent
}
