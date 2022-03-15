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
	switch r.Method {
	case "POST":
		fmt.Fprintf(w, "bruh")
	default:
		fmt.Fprintf(w, "POST plz!")
	}
}

func main() {
	// L.QRCodeGen("ucid", "qr.png")
	port := strconv.Itoa(PORT)

	http.HandleFunc("/", root)

	fmt.Printf("http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
