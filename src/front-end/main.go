package main

import (
	"context"
	"crypto/rand"
	"encoding/base32"
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const (
	PORT = 10234
)

var (
	google0authConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:10234/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	randomState = getToken(12)
)

func getToken(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

type Template struct {
	Sig, Favicon string
}

// Check : Func to do error checking
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	url := google0authConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	if r.FormValue("state") != randomState {
		fmt.Println("state not valid")
		return
	}

	token, err := google0authConfig.Exchange(context.Background(), r.FormValue("code"))
	Check(err)

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	Check(err)

	content, err := ioutil.ReadAll(resp.Body)
	Check(err)

	fmt.Fprintf(w, "%s", string(content))

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

	http.HandleFunc("/login", handleLogin)

	http.HandleFunc("/callback", handleCallback)

	// does dynamic url routing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// gets the sig from the url path
		sig := r.URL.Path[1:]
		// checks if sig in the set, if so the magic starts to begin
		_, exists := validSigs[sig]
		if exists {
			switch r.Method {
			case "POST":
				fmt.Fprintf(w, "do something")
			default:
				tpl, err := template.ParseFiles("templates/layout.html")
				Check(err)

				vals := Template{
					Sig:     sig,
					Favicon: "http://jerseyctf.com/assets/img/white_hollow_acm.png",
				}
				tpl.ExecuteTemplate(w, "startcore", vals)

				tpl.ExecuteTemplate(w, "end", nil)
			}
		} else {
			fmt.Fprintf(w, "<html><p>%s is not a valid sig! (But you can create it if you want!)</p></html>", sig)
		}
	})

	fmt.Printf("http://localhost:%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
