package http

import (
    "errors"
    "net/http"
)

type pxy struct {
    address string
}

func New(address string) (*pxy, error) {
    if address == "" {
        return nil, errors.New("need address")
    }

    h := &pxy{
        address: address,
    }

    return h, nil
}

func (this *pxy) Server() {
    http.Handle("/", &Server{})
    http.ListenAndServe(this.address, nil)
}
