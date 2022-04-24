package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	L "github.com/njitacm/sig-track/src/api/util"
)

const (
	PORT = 10233
	// FENDPOINT = "http://localhost:10234" // local
	FENDPOINT = "https://sig-track.xyz" // production
	FILENAME  = "attendeeList.json"
)

type POSTREQ struct {
	Sig     string `json:"sig"`
	Ucid    string `json:"ucid"`
	Time    string `json:"time"`
	Meeting string `json:"meeting"`
}

func handleList(w http.ResponseWriter, r *http.Request) {

	// create file if it does not exist
	if _, err := os.Stat(FILENAME); err != nil {
		_, err = os.Create(FILENAME)
		L.Check(err)
	}

	var attendeeList []POSTREQ

	file, err := os.Open(FILENAME)
	L.Check(err)

	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.Decode(&attendeeList)

	switch r.Method {
	case "GET":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "    ")
		enc.Encode(attendeeList)
	case "POST":
		var getPost POSTREQ

		// decode json from request body of POST request
		err := json.NewDecoder(r.Body).Decode(&getPost)
		L.Check(err)

		defer r.Body.Close()

		// append POST request to attendeeList
		attendeeList = append(attendeeList, POSTREQ{
			Sig:     getPost.Sig,
			Ucid:    getPost.Ucid,
			Time:    getPost.Time,
			Meeting: getPost.Meeting,
		})

		// convert attendeeList to []byte to write to attendeeList.json
		data, err := json.Marshal(attendeeList)
		L.Check(err)

		// code to overwrite file
		f, err := os.OpenFile(FILENAME, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		L.Check(err)

		// write to file (being overwritten)
		f.Write(data)

	default:
		fmt.Fprintf(w, "Error!")
	}
}

func handleGen(w http.ResponseWriter, r *http.Request) {
	/*
		ex:
		handleGen:
		http://localhost:10233/gen?sig=swe&meeting=7
	*/

	// enable cors
	L.EnableCors(&w)

	switch r.Method {
	case "GET":
		q := r.URL.Query()
		sig := q.Get("sig")
		meeting := q.Get("meeting")
		if len(sig) == 0 {
			fmt.Fprintf(w, "error in query, must have `sig={sig-name}`")
		}
		// fmt.Println(fmt.Sprintf("%s/%s?meeting=%s", FENDPOINT, sig, meeting))
		image := L.QRCodeGen(fmt.Sprintf("%s/%s?meeting=%s", FENDPOINT, sig, meeting))
		w.Write(image)
	case "POST":
		var res map[string]string
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&res)
		L.Check(err)
		sig := res["sig"]
		meeting := res["meeting"]
		image := L.QRCodeGen(fmt.Sprintf("%s/%s?meeting=%s", FENDPOINT, sig, meeting))
		w.Write(image)
	default:
		fmt.Fprintf(w, "No support yet!")
	}
}

func main() {
	port := strconv.Itoa(PORT)

	http.HandleFunc("/", handleList)
	http.HandleFunc("/gen", handleGen)
	http.HandleFunc("/list", handleList)
	http.HandleFunc("/stats", handleList)

	fmt.Printf("http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
