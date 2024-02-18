package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	L "github.com/njitacm/sig-track/src/api/util"
)

func init() {
	// loads the .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

const (
	PORT     = 10233
	FILENAME = "attendeeList.json"
)

type POSTREQ struct {
	Sig     string `json:"sig"`
	Ucid    string `json:"ucid"`
	Time    string `json:"time"`
	Meeting string `json:"meeting"`
}

func handleGen(w http.ResponseWriter, r *http.Request) {
	/*
		ex:
		handleGen:
		http://localhost:10233/gen?sig=swe&meeting=7
	*/
	var fendpoint string
	redirectType := os.Getenv("TYPE")

	// enable cors
	L.EnableCors(&w)

	switch strings.ToLower(redirectType) {
	case "test":
		fendpoint = "http://localhost:10234"
	case "prod":
		fendpoint = "https://sig-track.com"
	default:
		fendpoint = "http://localhost:10234"
	}

	switch r.Method {
	case "GET":
		q := r.URL.Query()
		sig := q.Get("sig")
		meeting := q.Get("meeting")
		if len(sig) == 0 {
			fmt.Fprintf(w, "error in query, must have `sig={sig-name}`")
		}
		// fmt.Println(fmt.Sprintf("%s/%s?meeting=%s", fendpoint, sig, meeting))
		image := L.QRCodeGen(fmt.Sprintf("%s/%s?meeting=%s", fendpoint, sig, meeting))
		w.Write(image)
	case "POST":
		var res map[string]string
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&res)
		L.Check(err)
		sig := res["sig"]
		meeting := res["meeting"]
		image := L.QRCodeGen(fmt.Sprintf("%s/%s?meeting=%s", fendpoint, sig, meeting))
		w.Write(image)
	default:
		fmt.Fprintf(w, "No support yet!")
	}
}

func main() {
	port := strconv.Itoa(PORT)

	//http.HandleFunc("/", handleList)
	http.HandleFunc("/gen", handleGen)
	//http.HandleFunc("/list", handleList)
	//http.HandleFunc("/stats", handleList)

	fmt.Printf("http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
