# sig-track

## :memo: Description
Tracks attendance for sig-meetings (with little overhead)

Try it out at: http://sig-track.com/

## :hammer: How to Build
- Prereqs:
    - Create `valid-sigs.json` file in `src/front-end` with a json string slice with the name of the sig 
        - ex. `[ "swe", "algo", "sec" ]` 
    - Create an `.env` file in `src/front-end` with the following values:
        - `GOOGLE_CLIENT_ID`
        - `GOOGLE_CLIENT_SECRET`
        - `SESSION_KEY`
        - `BENDPOINT`
        - `TYPE`
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
- Packages: `make`, `nginx`

## :card_file_box: Directory Explanation
| Directory          | Explanation
| :-------:          | :-----:
| [configs](configs) | Infrastructure Configuration
| [src](src)         | Source Code for front-end and api services


## :blue_book: Technical Details
- Check out the [Wiki](https://github.com/njitacm/sig-track/wiki)