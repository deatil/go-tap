package ws

import (
    "log"
    "net/url"
    "net/http"
    "net/http/httputil"
)

type websocket struct {
    src string
    dst string
}

func New(src, dst string) *websocket {
    ws := &websocket{
        src: src,
        dst: dst,
    }

    return ws
}

func (w *websocket) Server() {
    url, err := url.Parse(w.dst)
    if err != nil {
        log.Println(err)
    }

    proxy := httputil.NewSingleHostReverseProxy(url)
    http.ListenAndServe(w.src, proxy)
}
