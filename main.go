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
    args := parseArgs()
    file := args["f"]
    url := args["u"]
    search := strings.Split(args["s"], ",")
    download := args["d"] != ""
    ext := args["e"]
    level, err := strconv.Atoi(args["l"])
    if err != nil {
        level = 1
    }

    if (file == "" && url == "") || (file != "" && url != ""){
        fmt.Println("Please provide a file (-f) or a url (-u) but not both")
        os.Exit(1)
    }

    var path string
    if (file != "") {
        path = file
    } else {
        path = url
    }

    document := getDocument(file, url)
    processDocument(document, path, search, download, ext, level)
}

func getDocument(file, url string) (document io.Reader){
    if (file != "") {
        fileData, err := ioutil.ReadFile(file)
        if err != nil {
            log.Fatal(err)
        }

        document = bytes.NewReader(fileData)
    } else {
        document = sendGetRequest(url)
    }

    return
}

func processDocument(document io.Reader, url string, search []string, download bool, extension string, level int) {
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
                            if (level > 1) {
                                doc := sendGetRequest(a.Val)
                                if (len(search) > 1) {
                                    search = search[1:]
                                }
                                processDocument(doc, a.Val, search, download, extension, level-1)
                            } else {
                                if download {
                                    downloadFile(url, a.Val, extension)
                                } else {
                                    fmt.Println(a.Val)
                                }
                                break
                            }
                        }
                    }
                }
        }
    }
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
