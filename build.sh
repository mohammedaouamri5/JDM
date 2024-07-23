#!/bin/bash

BINARY="tmp/main"
LOG="tmp/build-errors.log"

# Define the build target
function build() {
	go build -o $BINARY . 2>&1 | tee $LOG
}
# Define the clean target
function clean() {
	if [ -f $BINARY ]; then rm $BINARY  ; fi
}
clean; 
build; 
