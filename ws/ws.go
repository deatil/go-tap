package ws

import (
    "log"
    "errors"
    "net/url"
    "net/http"
    "net/http/httputil"
)

type websocket struct {
    addr  string
    proxy string
}

func New(addr, proxy string) (*websocket, error) {
    if addr == "" {
        return nil, errors.New("need addr flag")
    }
    if proxy == "" {
        return nil, errors.New("need proxy flag")
    }

    ws := &websocket{
        addr:  addr,
        proxy: proxy,
    }

    return ws, nil
}

func (w *websocket) Server() {
    url, err := url.Parse(w.proxy)
    if err != nil {
        log.Println(err)
    }

    proxy := httputil.NewSingleHostReverseProxy(url)
    http.ListenAndServe(w.addr, proxy)
}
