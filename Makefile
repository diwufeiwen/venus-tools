SHELL=/usr/bin/env bash

BINS:=

build:
	go build -o venus-tools ./app
BINS+=venus-tools

clean:
	rm -rf $(BINS)
.PHONY: clean
