package main

import (
	"bufio"
	"log"
	"os"
	"regexp"
	"time"

	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

var (
	inputFile = kingpin.
		Arg("file", "log file to tail").
		Required().
		ExistingFile()
)

// HaproxyDateRegex matches the beginning of a Haproxy 1.8 log string to be able
// to grab the date.
var HaproxyDateRegex = regexp.MustCompile(`^((?P<server_host>[a-zA-Z0-9_-]+):(?P<server_port>[0-9]+))? *(?P<client_host>[^ ]+?):(?P<client_port>[0-9]+) \[(?P<request_date>[^\]]+)\]`)

// HaproxyDateLayout is a time.Parse layout
var HaproxyDateLayout = "02/Jan/2006:15:04:05.000"

func main() {
	kingpin.Parse()

	file, err := os.Open(*inputFile)

	if err != nil {
		log.Fatalf("unable to open input file: %s", err)
	}

	defer file.Close()

	reader := bufio.NewReader(file)

	logFirstLine, err := reader.ReadString('\n')

	match := HaproxyDateRegex.FindStringSubmatch(logFirstLine)

	dateStr := match[len(match)-1]

	startTime, err := time.Parse(HaproxyDateLayout, dateStr)

	// log.Fatalf("got something: %+v", t)

	for {
		line, err := reader.ReadString('\n')

		if err != nil {
			log.Fatalf("unable to read line: %s", err)
		}

		lineTime, err := time.Parse(HaproxyDateLayout, line)
		log.Fatalf("got linetime: %+v", lineTime)

		offset := lineTime.Sub(startTime)
		log.Fatalf("got something: %+v", offset)

	}

}
