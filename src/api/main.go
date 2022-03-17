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
	PORT      = 10233
	BENDPOINT = "ec2-3-21-33-128.us-east-2.compute.amazonaws.com"
	FENDPOINT = "http://localhost:10234"
	FILENAME  = "attendeeList.json"
)

type POSTREQ struct {
	Sig  string `json:"sig"`
	Ucid string `json:"ucid"`
	Time string `json:"time"`
}

func handleRoot(w http.ResponseWriter, r *http.Request) {

	var attendeeList []POSTREQ

	file, err := os.Open(FILENAME)
	L.Check(err)

	defer file.Close()

	decoder := json.NewDecoder(file)
	decoder.Decode(&attendeeList)

	switch r.Method {
	case "GET":
		// fmt.Fprintf(w, "%v", attendeeList)

		for i := range attendeeList {
			fmt.Fprintf(w, "%v\n", attendeeList[i])
		}
	case "POST":

		var getPost POSTREQ

		err := json.NewDecoder(r.Body).Decode(&getPost)
		L.Check(err)

		defer r.Body.Close()

		attendeeList = append(attendeeList, POSTREQ{
			Sig:  getPost.Sig,
			Ucid: getPost.Ucid,
			Time: getPost.Time,
		})

		// convert attendeeList to []byte to write to attendeeList.json
		data, err := json.Marshal(attendeeList)
		L.Check(err)

		f, err := os.OpenFile(FILENAME, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		L.Check(err)

		f.Write(data)

		// fmt.Fprintf(w, "%s", string(data))

		// encoder := json.NewEncoder(file)
		// encoder.Encode((&attendeeList))

		// err := ioutil.WriteFile(FILENAME, []byte(attendeeList))
		// L.Check(err)

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
	port := strconv.Itoa(PORT)

	http.HandleFunc("/", handleRoot)
	http.HandleFunc("/gen", handleGen)

	fmt.Printf("http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
