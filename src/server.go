package main

import (
	"fmt"
	"net"
	"os"

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
	// len == 0 means root is requested
	case len(req.Target) == 0:
		conn.Write(lib.ToResponse(req, 200, nil))
	case req.Target[0] == "echo":
		status, body := handleEcho(req)
		conn.Write(lib.ToResponse(req, status, body))
	case req.Target[0] == "files":
		status, body := handleFiles(req)
		conn.Write(lib.ToResponse(req, status, body))
	case req.Target[0] == "user-agent":
		status, body := handleUserAgent(req)
		conn.Write(lib.ToResponse(req, status, body))
	default:
		conn.Write(lib.ToResponse(req, 404, nil))
	}
}

type status = lib.Status
type body = lib.Body

func handleEcho(req lib.Request) (status, body) {
	if len(req.Target) < 2 {
		return 400, nil
	}
	yell := req.Target[1]
	return 200, yell
}

func handleFiles(req lib.Request) (status, body) {
	if len(req.Target) < 2 {
		return 400, nil
	}
	fileName := req.Target[1]
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
