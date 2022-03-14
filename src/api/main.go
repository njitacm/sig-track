package main

import (
	"fmt"
	"net/http"
	"strconv"

	L "github.com/njitacm/sig-track/src/api/util"
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
	L.QRCodeGen("ucid", "qr.png")

	port := strconv.Itoa(PORT)
	http.HandleFunc("/", root)
	err := http.ListenAndServe(":"+port, nil)
	L.Check(err)
}
