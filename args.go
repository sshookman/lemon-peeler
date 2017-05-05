package main

import (
    "os"
    "strings"
)

func parseArgs() (m map[string]string) {
    m = make(map[string]string)
    argLen := len(os.Args)

    for i := 1; i < argLen; i++ {
        if strings.HasPrefix(os.Args[i], "-") {

            if i+1 < argLen && !strings.HasPrefix(os.Args[i+1], "-") {
                m[os.Args[i][1:]] = os.Args[i+1]
                i++
            } else {
                m[os.Args[i][1:]] = "true"
            }
        }
    }

    return
}
