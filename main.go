// main.go
package main

import (
    "fmt"
    "log"
    "net/http"
    "html/template"

    "github.com/garyburd/redigo/redis"
)

func main() {
    http.HandleFunc("/", viewHandler)

    log.Fatal(http.ListenAndServe(":5000", nil))
}

type Page struct {
    Color string
}

func viewHandler (w http.ResponseWriter, r *http.Request) {
    redi, err := redis.Dial("tcp", "redis:6379")
    if err != nil {
        log.Fatal(err)
    }

    res, err := redi.Do("incr", "counter")
    if err != nil {
        w.WriteHeader(500)
        w.Write([]byte(err.Error()))

        return
    }

    if res, ok := res.(int64); ok {
        fmt.Print(res)
        //w.Write([]byte(fmt.Sprintf("counter: %d", res)))
        createDomFromTemplate(w, res)
    } else {
        w.WriteHeader(500)
        w.Write([]byte("unexpected value"))
    }
}

func intToHexStr(num int64) (rgb string) {
    // Todo: calculate rgb. This is a makeshift.
    color := fmt.Sprintf("%x", num)
    rgb = color + color + color

    return rgb
}

func createDomFromTemplate (w http.ResponseWriter, res int64) {
    tmpl, err := template.ParseFiles("template/view.html")
    if err != nil {
        panic(err)
    }

    rgb := intToHexStr(res)
    fmt.Print(rgb)
    page := Page{Color: rgb}
    err = tmpl.Execute(w, page)
    if err != nil {
        panic(err)
    }
}