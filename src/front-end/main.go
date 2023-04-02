package main

import (
	"bytes"
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

	"github.com/joho/godotenv"
)

const (
	PORT = 10234
)

var (
	store       = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	randomState = getToken(12)
)

// init: function that get's called on initialization
func init() {
	// loads the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type POSTREQ struct {
	Sig     string `json:"sig"`
	Ucid    string `json:"ucid"`
	Time    string `json:"time"`
	Meeting string `json:"meeting"`
}

type Template struct {
	Sig, Favicon string
}

func GOOGLE0authConfigFunc() *oauth2.Config {

	var redirectURL string

	redirectType := os.Getenv("TYPE")

	switch strings.ToLower(redirectType) {
	case "test":
		redirectURL = "http://localhost:10234/oauth2/callback"
	case "prod":
		redirectURL = "https://sig-track.com/oauth2/callback"
	default:
		redirectURL = "http://localhost:10234/oauth2/callback"
	}

	return &oauth2.Config{
		RedirectURL:  redirectURL,
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}
}

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

func CheckHandler(err error, errMsg string, w http.ResponseWriter, r *http.Request) {
	if err != nil {
		fmt.Fprintf(w, "something up with %s: %v", errMsg, err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
}

// handleLogin: does the google oauth2 authorization
func handleLogin(w http.ResponseWriter, r *http.Request) {
	EnableCors(&w)
	google0authConfig := GOOGLE0authConfigFunc()
	url := google0authConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleCallback: after successful connections
func handleCallback(w http.ResponseWriter, r *http.Request) {

	// get template setup
	tpl := template.Must(template.ParseGlob("templates/*.html"))

	// enable cors
	EnableCors(&w)
	// store := sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
	session, err := store.Get(r, "session-name")
	Check(err)

	// Get sig value from session cookie
	sig, ok := session.Values["sig"].(string)
	if !ok {
		fmt.Fprintf(w, "couldn't get session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Get meeting value from session cookie
	meeting, ok := session.Values["meeting"].(string)
	if !ok {
		fmt.Fprintf(w, "couldn't get session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var res map[string]interface{}

	if r.FormValue("state") != randomState {
		fmt.Println("state not valid")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// get token from oauth2 jazz
	google0authConfig := GOOGLE0authConfigFunc()
	token, err := google0authConfig.Exchange(context.Background(), r.FormValue("code"))
	CheckHandler(err, "getting token", w, r)

	// get user email data from google's oauth2 using token
	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	CheckHandler(err, "getting response", w, r)

	defer resp.Body.Close()

	// read content from response
	content, err := ioutil.ReadAll(resp.Body)
	Check(err)

	json.Unmarshal(content, &res)

	// convert email to string from any type
	email := fmt.Sprintf("%s", res["email"])

	// parse info into lovely data structure
	checkIn := POSTREQ{
		Sig:     sig,
		Ucid:    email[:strings.Index(email, "@")],
		Time:    time.Now().UTC().Format(time.UnixDate),
		Meeting: meeting,
	}

	fmt.Println(checkIn)

	json_data, err := json.Marshal(checkIn)
	Check(err)

	_, err = http.Post(os.Getenv("BENDPOINT"), "application/json", bytes.NewBuffer(json_data))
	Check(err)

	tpl.ExecuteTemplate(w, "done", nil)

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

	// route handling
	http.HandleFunc("/oauth2/sign_in", handleLogin)

	http.HandleFunc("/oauth2/callback", handleCallback)

	// does dynamic url routing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tpl := template.Must(template.ParseGlob("templates/*.html"))

		// enables cors
		EnableCors(&w)

		// gets the sig from the url path
		sig := r.URL.Path[1:]
		meeting := r.URL.Query().Get("meeting")
		if len(sig) == 0 {
			// root handler
			tpl.ExecuteTemplate(w, "root", sigs)
		}

		// checks if sig in the set, if so the magic starts to begin
		_, exists := validSigs[sig]
		if exists {
			// starts up a session
			session, err := store.Get(r, "session-name")
			Check(err)

			// Adds key:sig = value:sig in session
			session.Values["sig"] = sig
			session.Values["meeting"] = meeting
			// Saves session
			session.Save(r, w)

			vals := Template{
				Sig:     sig,
				Favicon: "http://jerseyctf.com/assets/img/white_hollow_acm.png",
			}
			tpl.ExecuteTemplate(w, "attendance", vals)

		} else {
			if len(sig) != 0 {
				fmt.Fprintf(w, "<html><p> <pre>%s</pre> is not a valid sig! (But you can create it if you want!)</p></html>", sig)
			}
		}
	})

	fmt.Printf("http://localhost:%s\n", port)
	log.Fatalln(http.ListenAndServe(":"+port, nil))
}
