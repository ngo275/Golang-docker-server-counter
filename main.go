// main.go
package main

import (
    "fmt"
    "log"
    "net/http"

    "github.com/garyburd/redigo/redis"
)

func main() {
    redi, err := redis.Dial("tcp", "redis:6379")
    if err != nil {
        log.Fatal(err)
    }
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        res, err := redi.Do("incr", "counter")
        if err != nil {
            w.WriteHeader(500)
            w.Write([]byte(err.Error()))
            return
        }

        if res, ok := res.(int64); ok {
            w.Write([]byte(fmt.Sprintf("counter: %d", res)))
        } else {
            w.WriteHeader(500)
            w.Write([]byte("unexpected value"))
        }
    })
    log.Fatal(http.ListenAndServe(":5000", nil))
}
