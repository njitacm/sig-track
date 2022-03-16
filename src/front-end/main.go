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
	PORT      = 10234
	BENDPOINT = "ec2-3-21-33-128.us-east-2.compute.amazonaws.com"
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
	store       = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))
)

type POSTREQ struct {
	Sig, Ucid, Time string
}

type Template struct {
	Sig, Favicon string
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
	url := google0authConfig.AuthCodeURL(randomState)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// handleCallback: after successful connections
func handleCallback(w http.ResponseWriter, r *http.Request) {

	// get template setup
	tpl := template.Must(template.ParseGlob("templates/*.html"))

	// enable cors
	EnableCors(&w)
	session, err := store.Get(r, "session-name")
	Check(err)

	// Get sig value from session cookie
	sig, ok := session.Values["sig"].(string)
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

	token, err := google0authConfig.Exchange(context.Background(), r.FormValue("code"))
	CheckHandler(err, "getting token", w, r)

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	CheckHandler(err, "getting token", w, r)

	defer resp.Body.Close()

	content, err := ioutil.ReadAll(resp.Body)
	Check(err)

	json.Unmarshal(content, &res)

	// convert email to string from any type
	email := fmt.Sprintf("%s", res["email"])

	checkIn := POSTREQ{
		Sig:  sig,
		Ucid: email[:strings.Index(email, "@")],
		Time: time.Now().UTC().Format(time.UnixDate),
	}

	fmt.Println(checkIn)

	// json_data, err := json.Marshal(checkIn)
	// Check(err)

	// resp, err = http.Post(BENDPOINT, "application/json", bytes.NewBuffer(json_data))
	// Check(err)

	fmt.Println()

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

	http.HandleFunc("/login", handleLogin)

	http.HandleFunc("/callback", handleCallback)

	// does dynamic url routing
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		tpl := template.Must(template.ParseGlob("templates/*.html"))

		// enables cors
		EnableCors(&w)

		// gets the sig from the url path
		sig := r.URL.Path[1:]
		if len(sig) == 0 {
			tpl.ExecuteTemplate(w, "done", nil)
		}

		// checks if sig in the set, if so the magic starts to begin
		_, exists := validSigs[sig]
		if exists {
			// starts up a session
			session, err := store.Get(r, "session-name")
			Check(err)

			// Adds key:sig = value:sig in session
			session.Values["sig"] = sig
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
