package http

import (
    "fmt"
    "io"
    "net"
    "net/http"
    "strings"
)

type Server struct {}

func (p *Server) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
    fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

    transport := http.DefaultTransport

    outReq := new(http.Request)
    *outReq = *req

    if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
        if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
            clientIP = strings.Join(prior, ", ") + ", " + clientIP
        }

        outReq.Header.Set("X-Forwarded-For", clientIP)
    }

    res, err := transport.RoundTrip(outReq)
    if err != nil {
        rw.WriteHeader(http.StatusBadGateway)
        return
    }

    for key, value := range res.Header {
        for _, v := range value {
            rw.Header().Add(key, v)
        }
    }

    rw.WriteHeader(res.StatusCode)
    io.Copy(rw, res.Body)
    res.Body.Close()
}
