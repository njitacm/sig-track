package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
)

const (
	PORT = 10234
)

// var (
// 	google0authConfig = &oauth2.Config()
// )

type Template struct {
	Sig, Favicon, Redirect string
}

// Check : Func to do error checking
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// main: does routing and core logic
func main() {
	// converts const PORT to string
	port := strconv.Itoa(PORT)

	// initializes new values (especially set)
	var sigs []string
	validSigs := make(map[string]struct{})
	var exists = struct{}{}

	// read file of valid sigs
	file, err := os.ReadFile("valid-sigs.json")
	Check(err)

	// convert json file into a string slice
	err = json.Unmarshal(file, &sigs)
	Check(err)

	// adds value to set
	for _, sig := range sigs {
		validSigs[sig] = exists
	}

	// does dynamic url routing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// gets the sig from the url path
		sig := r.URL.Path[1:]
		// checks if sig in the set, if so the magic starts to begin
		_, exists := validSigs[sig]
		if exists {
			tpl, err := template.ParseFiles("templates/layout.html")
			Check(err)

			vals := Template{
				Sig:      sig,
				Favicon:  "http://jerseyctf.com/assets/img/white_hollow_acm.png",
				Redirect: "https://empty-room.xyz",
			}
			tpl.ExecuteTemplate(w, "startcore", vals)

			tpl.ExecuteTemplate(w, "end", nil)
		} else {
			fmt.Fprintf(w, "<html><p>%s is not a valid sig! (But you can create it if you want!)</p></html>", sig)
		}
	})

	fmt.Printf("http://localhost:%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
