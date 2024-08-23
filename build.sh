#!/bin/bash

BINARY="tmp/main"
LOG="tmp/build-errors.log"



# Define the build target
function build() {
	mkdir -p data
	go build -o $BINARY . 2>&1 | tee $LOG
}
# Define the clean target
function clean() {
	rm -rv  tmp
	rm -rv  data
}
clean; 
build; 
