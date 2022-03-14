package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
)

const (
	PORT = 10234
)

type Vals struct {
	Sig, Favicon, Redirect string
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func root(w http.ResponseWriter, r *http.Request) {
	tpl, err := template.ParseFiles("templates/layout.html")
	Check(err)

	vals := Vals{
		Sig:      "swe",
		Favicon:  "http://jerseyctf.com/assets/img/white_hollow_acm.png",
		Redirect: "https://empty-room.xyz",
	}
	tpl.ExecuteTemplate(w, "startcore", vals)

	tpl.ExecuteTemplate(w, "end", nil)
}

func main() {
	port := strconv.Itoa(PORT)

	http.HandleFunc("/", root)

	fmt.Printf("http://localhost:%s\n", port)
	err := http.ListenAndServe(":"+port, nil)
	Check(err)

}
