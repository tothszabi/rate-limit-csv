package main

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	path := os.Args[1]
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	re := regexp.MustCompile(".* -- hostname \"(.*)\"")

	reader := csv.NewReader(f)

	// Ignore csv header
	if _, err := reader.Read(); err != nil {
		log.Fatal(err)
	}

	hostnames := make(map[string]struct{})

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if len(record) != 4 {
			log.Fatal("Expected 4 fields, got ", len(record))
		}

		matches := re.FindStringSubmatch(record[3])
		if len(matches) != 2 {
			log.Fatal("Expected 2 fields, got ", len(matches))
		}

		hostnames[matches[1]] = struct{}{}
	}

	var buffer bytes.Buffer
	buffer.WriteString("(")

	for hostname := range hostnames {
		buffer.WriteString(fmt.Sprintf("'%s', ", hostname))
	}

	query := strings.TrimSuffix(buffer.String(), ", ")
	query += ")"

	fmt.Println(query)
}
