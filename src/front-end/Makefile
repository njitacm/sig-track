SHELL=/bin/bash
BUILDFILE=main.go
OUTFILE=main.o

all: run

getenvid := "os.Getenv(\"GOOGLE_CLIENT_ID\")"
clientid := "\"`printenv GOOGLE_CLIENT_ID`\""
getenvsecret := "os.Getenv(\"GOOGLE_CLIENT_SECRET\")"
clientsecret := "\"`printenv GOOGLE_CLIENT_SECRET`\""

compile:
	@go build -o $(OUTFILE) $(BUILDFILE)

run: compile
	@./main.o

prod: compile
# Replace CLIENT_ID with Hardcoded String
	@sed -i -e s/$(getenvid)/$(clientid)/g $(BUILDFILE)
# Replace CLIENT_SECRET with Hardcoded String
	@sed -i -e s/$(getenvsecret)/$(clientsecret)/g $(BUILDFILE)

clean:
ifneq (,$(wildcard ./$(OUTFILE)))
	@rm ./$(OUTFILE)
endif