package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	L "github.com/njitacm/sig-track/src/api/util"
)

const (
	PORT      = 10233
	BENDPOINT = "ec2-3-21-33-128.us-east-2.compute.amazonaws.com"
	FENDPOINT = "http://localhost:10234"
)

type POSTREQ struct {
	Sig, Ucid, Time string
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":

	case "POST":
		fmt.Fprintf(w, "bruh")
	default:
		fmt.Fprintf(w, "Error!")
	}
}

func handleGen(w http.ResponseWriter, r *http.Request) {
	/*
		ex:
		handleGen:
		http://localhost:9998/gen?sig=swe
	*/

	// enable cors
	L.EnableCors(&w)

	switch r.Method {
	case "GET":
		q := r.URL.Query()
		sig := q.Get("sig")
		if len(sig) == 0 {
			fmt.Fprintf(w, "error in query, must have `sig={sig-name}`")
		}
		image := L.QRCodeGen(fmt.Sprintf("%s/%s", FENDPOINT, sig))
		w.Write(image)
	case "POST":
		var res map[string]string
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&res)
		L.Check(err)
		sig := res["sig"]
		image := L.QRCodeGen(fmt.Sprintf("%s/%s", FENDPOINT, sig))
		w.Write(image)
	default:
		fmt.Fprintf(w, "No support yet!")
	}
}

func main() {
	// L.QRCodeGen("ucid", "qr.png")
	port := strconv.Itoa(PORT)

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/gen", handleGen)
	// http.HandleFunc("/list", handleList)
	// http.HandleFunc("/add", handleAdd)

	fmt.Printf("http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
