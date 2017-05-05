package main

import (
    "io"
    "log"
    "net/http"
)

func sendGetRequest(url string) (doc io.Reader) {
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }

    return resp.Body
}
