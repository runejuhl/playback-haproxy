package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"time"
)

var inputFile string
var speed = 1.0
var err error
var version = "0.1.0"

// HaproxyDateRegex matches the beginning of a Haproxy 1.8 log string to be able
// to grab the date.
var HaproxyDateRegex = regexp.MustCompile(`^((?P<server_host>[a-zA-Z0-9_-]+):(?P<server_port>[0-9]+))? *(?P<client_host>[^ ]+?):(?P<client_port>[0-9]+) \[(?P<request_date>[^\]]+)\]`)

var dateStringIndex = len(HaproxyDateRegex.SubexpNames()) - 1

// HaproxyDateLayout is a time.Parse layout
var HaproxyDateLayout = "02/Jan/2006:15:04:05.000"

// LineToDate takes a Haproxy 1.8 log line and returns the timestamp
func LineToDate(s string) (t time.Time, err error) {
	match := HaproxyDateRegex.FindStringSubmatch(s)
	dateStr := match[dateStringIndex]

	t, err = time.Parse(HaproxyDateLayout, dateStr)

	return
}

func usage() {
	fmt.Printf(
		"Play back haproxy log files in realtime using timing information from logs\n"+
			"\n"+
			"Usage:\n"+
			"  %s <input-file> [speed]\n"+
			"\n"+
			"Examples:\n"+
			"  %s /var/log/haproxy.log\n"+
			"  %s /var/log/haproxy.log 4.2\n"+
			"\n"+
			"https://gitlab.com/runejuhl/playback-haproxy\n"+
			"https://github.com/runejuhl/playback-haproxy\n"+
			"\n"+
			"playback-haproxy v%s copyright (C) 2019 Rune Juhl Jacobsen\n"+
			"This program comes with ABSOLUTELY NO WARRANTY.\n"+
			"This is free software, and you are welcome to redistribute it \n"+
			"under certain conditions; see source repository for details.\n",
		os.Args[0], os.Args[0], os.Args[0], version)
	os.Exit(0)
}

func main() {

	if len(os.Args) < 2 || len(os.Args) > 3 {
		usage()
	}

	inputFile = os.Args[1]

	if inputFile == "-v" ||
		inputFile == "--version" ||
		inputFile == "-h" ||
		inputFile == "--help" {
		usage()
	}

	if len(os.Args) > 2 {
		speed, err = strconv.ParseFloat(os.Args[2], 64)
		if err != nil {
			log.Fatalf("invalid speed '%s': %s", os.Args[2], err)
		}
	}

	if speed <= 0 {
		log.Fatalf("speed cannot be <= 0")
	}

	speedMultiplier := (time.Duration(speed * 1e6))

	file, err := os.Open(inputFile)

	if err != nil {
		log.Fatalf("unable to open input file: %s", err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)
	currentTime := time.Date(2999, 0, 0, 0, 0, 0, 0, time.UTC)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalf("unable to read line: %s", err)
		}

		lineTime, err := LineToDate(line)

		offset := lineTime.Sub(currentTime)

		if offset > 0 {
			sleep := offset * speedMultiplier / 1e6
			log.Printf(
				"current time %s, next time %s, sleeping for: %+v",
				currentTime, lineTime,
				sleep)
			time.Sleep(sleep)
		}

		fmt.Print(line)

		currentTime = lineTime
	}

}
