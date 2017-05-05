package main

import (
    "os"
    "io"
    "io/ioutil"
    "bytes"
    "log"
    "fmt"
    "strings"
    "strconv"
    "net/http"
    "golang.org/x/net/html"
)

func main() {
    path, search, output, extension, dive := readArguments()
    document := getDocument(path)

    processDocument(document, path, search, output, extension, dive)
}

func getDocument(path string) (document io.Reader){
    if strings.HasPrefix(path, "FILE:") {
        file, err := ioutil.ReadFile(path[5:])
        if err != nil {
            log.Fatal(err)
        }

        document = bytes.NewReader(file)
    } else {
        resp, err := http.Get(path)
        if err != nil {
            log.Fatal(err)
        }

        document = resp.Body
    }

    return
}

func processDocument(document io.Reader, url string, search []string, output, extension string, dive int) {
    tokenizer := html.NewTokenizer(document)
    for {

        nextToken := tokenizer.Next()
        switch {
            case nextToken == html.ErrorToken:
                return

            case nextToken == html.StartTagToken:
                token := tokenizer.Token()

                isAnchor := token.Data == "a"
                if isAnchor {
                    for _, a := range token.Attr {
                        if a.Key == "href" && strings.Contains(a.Val, search[0]) {
                            if (dive > 0) {
                                doc := sendGetRequest(a.Val)
                                processDocument(doc, a.Val, search[1:], output, extension, dive-1)
                            } else {
                                if output == "download" {
                                    downloadFile(url, a.Val, extension)
                                } else {
                                    fmt.Println(url + a.Val)
                                }
                                break
                            }
                        }
                    }
                }
        }
    }
}

func sendGetRequest(url string) (doc io.Reader) {
    resp, err := http.Get(url)
    if err != nil {
        log.Fatal(err)
    }

    return resp.Body
}

func readArguments() (path string, search []string, output, extension string, dive int) {
    if len(os.Args) < 2 {
        fmt.Println("Please Provide a path (URL or FILE)")
        log.Fatal()
    }

    path = os.Args[1]
    if len(os.Args) > 2 {
        search = strings.Split(os.Args[2], ",")
    }
    if len(os.Args) > 3 {
        output = os.Args[3]
    }
    if len(os.Args) > 4 {
        extension = os.Args[4]
    }
    if len(os.Args) > 5 {
        var err error
        dive, err = strconv.Atoi(os.Args[5])
        if err != nil {
            log.Fatal("Please Provide and integer for the dive count")
        }
    }

    return
}

func downloadFile(url, downloadUrl, extension string) {

    if strings.HasPrefix(downloadUrl, "/") {
        downloadUrl = url + downloadUrl
    }

    fmt.Print("Downloading ", downloadUrl)
    fmt.Print("...")

    response, e := http.Get(downloadUrl)
    if e != nil {
        log.Fatal(e)
    }

    defer response.Body.Close()


    file, err := os.Create(url[strings.LastIndex(url, "/")+1:] + extension)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print("...")

    _, err = io.Copy(file, response.Body)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Print("...")

    file.Close()
    fmt.Println("Success!")
}
