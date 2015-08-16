package main

import (
	"io"
	"net/http"
	"fmt"
	"unsafe"
)

func bigBytes() *[]byte {
        s := make([]byte, 250000000)
        return &s
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
	a := bigBytes()
        fmt.Printf("memory type:%T address: %p  size: %d\n", a, &a, unsafe.Sizeof(a))
        b := bigBytes()
        fmt.Printf("memory type:%T address: %p  size: %d\n", b, &b, unsafe.Sizeof(b))
        c := bigBytes()
        fmt.Printf("memory type:%T address: %p  size: %d\n", c, &c, unsafe.Sizeof(c))
        d := bigBytes()
        fmt.Printf("memory type:%T address: %p  size: %d\n", d, &d, unsafe.Sizeof(d))
        e := bigBytes()
        fmt.Printf("memory type:%T address: %p  size: %d\n", e, &e, unsafe.Sizeof(e))
        f := bigBytes()
        fmt.Printf("memory type:%T address: %p  size: %d\n", f, &f, unsafe.Sizeof(f))
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}

