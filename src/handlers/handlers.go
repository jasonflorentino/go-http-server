package handlers

import (
	"fmt"
	"os"

	"github.com/jasonflorentino/go-http-server/src/config"
	"github.com/jasonflorentino/go-http-server/src/lib"
)

type status = lib.Status
type body = lib.Body

func HandleEcho(req lib.Request) (status, body) {
	if len(req.Target) < 2 {
		return 400, nil
	}
	yell := req.Target[1]
	return 200, yell
}

func HandleFiles(req lib.Request) (status, body) {
	if len(req.Target) < 2 {
		return 400, nil
	}
	fileName := req.Target[1]
	switch req.Method {
	case "GET":
		dat, err := os.ReadFile(fmt.Sprintf("%s%s", config.FILE_PATH, fileName))
		if err != nil {
			fmt.Println("Error reading file:", err.Error())
			return 404, nil
		}
		return 200, dat
	case "POST":
		err := os.WriteFile(fmt.Sprintf("%s%s", config.FILE_PATH, fileName), []byte(req.Body), 0644)
		if err != nil {
			fmt.Println("Error writing file:", err.Error())
			return 404, nil
		}
		return 201, nil
	default:
		return 400, nil
	}
}

func HandleUserAgent(req lib.Request) (status, body) {
	userAgent := req.Headers["User-Agent"]
	if userAgent == "" {
		return 400, nil
	}
	return 200, userAgent
}
