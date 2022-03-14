package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
)

const (
	PORT = 9998
)

type Ret struct {
	Ucid string `json:"ucid"`
	Sig  string `json:"sig"`
}

func root(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello World!")
	// dec := json.NewEncoder(w)
	// dec

	// enc := json.NewEncoder(w)
	// enc.Encode()

}

func main() {
	// L.QRCodeGen("ucid", "qr.png")
	port := strconv.Itoa(PORT)

	http.HandleFunc("/", root)

	fmt.Printf("http://localhost:%s\n")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
