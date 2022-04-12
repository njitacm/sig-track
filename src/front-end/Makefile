SHELL=/bin/bash
BUILDFILE=main.go
OUTFILE=main.o

all: compile

compile:
	@go build -o $(OUTFILE) $(BUILDFILE)

run: compile
	@./main.o

clean:
ifneq (,$(wildcard ./$(OUTFILE)))
	@rm ./$(OUTFILE)
endif