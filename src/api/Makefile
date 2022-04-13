SHELL=/bin/bash

all: compile

compile:
	@go build -o main.o main.go

run: compile
	@./main.o

clean:
ifneq (,$(wildcard ./main.o))
	@rm ./main.o
endif