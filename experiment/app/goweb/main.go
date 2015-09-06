package main

import (
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"
	"unsafe"
)

func bigBytes() *[]byte {
	s := make([]byte, 250000000)
	return &s
}

func hello(w http.ResponseWriter, r *http.Request) {
	a := bigBytes()
	fmt.Printf("memory type:%T address: %p  size: %d\n", a, &a, unsafe.Sizeof(a))
	b := bigBytes()
	fmt.Printf("memory type:%T address: %p  size: %d\n", b, &b, unsafe.Sizeof(b))
	time.Sleep(3000 * time.Millisecond)
	io.WriteString(w, "Hello world!")
	debug.FreeOSMemory()
}

func main() {
	http.HandleFunc("/", hello)
	http.ListenAndServe(":8000", nil)
}
