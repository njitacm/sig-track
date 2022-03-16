# sig-track

## :memo: Description
Tracks attendance for sig-meetings (with little overhead)

## :hammer: How to Build
- Prereqs:
    - Be sure to add `valid-sigs.json` in `/src/front-end` with a json string slice with the name of the sig ex. [ "swe", "algo", "sec" ] 
    - Be sure to have `GOOGLE_CLIENT_ID`, `GOOGLE_CLIENT_SECRET` & `SESSION_KEY` as env variables!

```sh
# Change directory to repo 
cd sig-track

# Change into src
cd src

# Build executable with output => (main.o)
make compile
```

## :alembic: How to Use
```sh
# Run Go executable (if already built)
./main.o

# Auto run from `make` 
make
```

## :microscope: Technologies
- Languages: `go`, `html`, `css`, `sh`
- Packages: `make`

## :card_file_box: Directory Explanation
| Directory      | Explanation
| :-------:      | :-----:
| [infra](infra) | Infrastructure Configuration
| [src](src)     | Source Code for front-end and api services


## :blue_book: Technical Details

## :books: Resources
- oauth2
    - **[getting-started-with-oauth2](https://www.youtube.com/watch?v=OdyXIi6DGYw)**  <-- BEST
    - [authentication-authorization-in-oauth2:golang](https://www.youtube.com/watch?v=Vmi3trk0rCk)
    - [getting-started,code walkthrough](https://www.youtube.com/watch?v=PdpQJsR-BpE)
- URL Dyno Routing
    - [Effective Go: Writing Web Applications](https://go.dev/doc/articles/wiki/)
- Go Env Variables
    - [gobyexample: env vars](https://gobyexample.com/environment-variables)
- url parsing
    - [gobyexample: url parse](https://gobyexample.com/url-parsing)
    - [Parsing Queries](https://www.youtube.com/watch?v=cl7_ouTMFh0)
- Sessions
    - [Go Sessions](https://gowebexamples.com/sessions/)
    - [Gorilla Sessions](https://github.com/gorilla/sessions)