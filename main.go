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
    Count int64
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

func createDomFromTemplate (w http.ResponseWriter, res int64) {
    tmpl, err := template.ParseFiles("view.html")
    if err != nil {
        panic(err)
    }

    page := Page{Count: res}
    err = tmpl.Execute(w, page)
    if err != nil {
        panic(err)
    }
}