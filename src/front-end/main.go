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
	"strings"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/gorilla/sessions"
)

const (
	PORT = 10234
)

type POSTREQ struct {
	Sig, Ucid, Time string
}

type Template struct {
	Sig, Favicon, Redirect string
}

var (
	google0authConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:10234/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
	randomState = getToken(12)
	store       = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
)

func EnableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
}

func getToken(length int) string {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		panic(err)
	}
	return base32.StdEncoding.EncodeToString(randomBytes)[:length]
}

// Check : Func to do error checking
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL.RawQuery)
	EnableCors(&w)
	url := google0authConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	session, err := store.Get(r, "session-name")
	Check(err)

	sig, ok := session.Values["sig"].(string)
	if !ok {
		fmt.Println("couldn't get session")
		return
	}

	fmt.Fprintf(w, "%v\n", sig)

	var res map[string]interface{}

	if r.FormValue("state") != randomState {
		fmt.Println("state not valid")
		return
	}

	token, err := google0authConfig.Exchange(context.Background(), r.FormValue("code"))
	Check(err)

	// fmt.Fprintf(w, "%v\n", r.FormValue("code"))

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	Check(err)

	// fmt.Fprintf(w, "%v\n", token.AccessToken)

	content, err := ioutil.ReadAll(resp.Body)
	Check(err)

	// fmt.Fprintf(w, "%v\n", r.URL.Query().Get("email"))

	// fmt.Fprintf(w, "%s\n", string(content))

	json.Unmarshal(content, &res)
	email := fmt.Sprintf("%s", res["email"])
	// fmt.Fprintf(w, "%s", email)

	checkIn := POSTREQ{
		Sig:  sig,
		Ucid: email[:strings.Index(email, "@")],
		Time: time.Now().UTC().Format(time.UnixDate),
	}

	fmt.Println(checkIn)
	// fmt.Fprintf(w, "%v\n", checkIn)
	// json_data, err := json.Marshal(checkIn)
	// Check(err)

	// resp, err = http.Post("", "application/json", bytes.NewBuffer(json_data))
	// Check(err)

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

		EnableCors(&w)
		// gets the sig from the url path
		sig := r.URL.Path[1:]
		// checks if sig in the set, if so the magic starts to begin
		_, exists := validSigs[sig]
		if exists {
			session, err := store.Get(r, "session-name")
			Check(err)

			session.Values["sig"] = sig
			session.Save(r, w)

			tpl, err := template.ParseFiles("templates/layout.html")
			Check(err)

			vals := Template{
				Sig:      sig,
				Favicon:  "http://jerseyctf.com/assets/img/white_hollow_acm.png",
				Redirect: "/login",
			}
			tpl.ExecuteTemplate(w, "startcore", vals)

			tpl.ExecuteTemplate(w, "end", nil)

			fmt.Println(r.URL.RawQuery)
		} else {
			fmt.Fprintf(w, "<html><p>%s is not a valid sig! (But you can create it if you want!)</p></html>", sig)
		}
	})

	fmt.Printf("http://localhost:%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
