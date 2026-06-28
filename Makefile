#! /usr/bin/make -f

# Calendar makefile

include .env

# REMOTE= rsync destination for publishing the calendar

YEAR=2026
SPREADSHEET="data/${YEAR}/2026-AppDev-V2.xlsx"

BUILD=build
CACHE=cache

all:	build

.PHONY:	build clean compare fetch prerequisites push backup

build:	
	go run wbc-calendar.go -o $(BUILD) -f $(SPREADSHEET)

pull:
	rm -rf live
	mkdir live
	rsync -a --delete $(REMOTE)/$(YEAR)/ live/

publish:
	rsync -a --delete $(BUILD)/ $(REMOTE)/$(YEAR)/

clean:
	rm -rf $(BUILD)
	mkdir $(BUILD)

backup:
	rm -rf save
	mkdir save
	rsync -a $(BUILD)/ save/

meld:
	meld "live" "${BUILD}"

