# sig-track

## :memo: Description
Tracks attendance for sig-meetings (with little overhead)

## :hammer: How to Build
- Prereqs:
    - Be sure to add `valid-sigs.json` with a json string slice with the name of the sig ex. [ "swe", "algo", "sec" ] 

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