SHELL=/bin/bash
OUTFILE=main.o

all: run

getenvid := "os.Getenv(\"GOOGLE_CLIENT_ID\")"
clientid := "\"`printenv GOOGLE_CLIENT_ID`\""
getenvsecret := "os.Getenv(\"GOOGLE_CLIENT_ID\")"
clientsecret := `printenv GOOGLE_CLIENT_SECRET`

compile:
	@go build -o $(OUTFILE) main.go

run: compile
	@./main.o

prod: compile
	@echo $(getenvid)
	@echo $(clientid)
	@sed -i -e s/$(getenvid)/$(clientid)/g $(OUTFILE)
	@sed -i -e s/$(getenvsecret)/$(clientsecret)/g $(OUTFILE) 

clean:
ifneq (,$(wildcard ./main.o))
	@rm ./main.o
endif