# Go HTTP Server

This repo contains the program I wrote for the *[Build Your Own HTTP server Challenge](https://app.codecrafters.io/courses/http-server/overview)* on [codecrafters.io](https://codecrafters.io). I've been eyeing their challenges for a bit so when I heard that this one would be free for the month, I couldn't help but use the opportunity to get more comfortable writing [Go](https://go.dev/). If you think this is bad/messy code now, you should have seen what it looked like while I was completing the challenge!

Jokes aside, this was so much fun to do and their platform felt natural and easy to use. Each step was nicely digestible so I could focus on figuring out  how to implement my ideas while maintaining a sense of momentum towards the overall goal. And being able to test and progress by pushing git commits to their remote meant I could spend most of my time in my own editor and terminal rathan than on a webpage.

— Jason, July 2024

### From their overview

[HTTP](https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol) is the
protocol that powers the web. In this challenge, you'll build a HTTP/1.1 server
that is capable of serving multiple clients.

Along the way you'll learn about TCP servers,
[HTTP request syntax](https://www.w3.org/Protocols/rfc2616/rfc2616-sec5.html),
and more.

# Usage
- Run the server with:
```bash
go run src/server.go
```
- Once running you can send it some requsts and you *should* get the right responses:
```bash
curl -v -H "Accept-Encoding: gzip" http://localhost:4221/echo/hello | gzip -dc
```

## Some things to try
- Request to `/` receives response 200 OK
- Request to unknown path receives response 404 Not Found
- Request to `/echo/:str` receives `str` in response body
- Request to `/user-agent` receives the User-Agent request header back in response body
- Can handle concurrent connections
- Returns a file's contents (specified by the `--directory` arg)
- Writes the request body to a file (specified by the `--directory` arg)
- Handles gzip encoding the response body
