package main

import (
	"bufio"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

const RESULT_FILE = "result.csv"
const FILE = "redirection.csv"
const BASE_PATH = "https://integration-5ojmyuq-6rg55inh3jus4.eu.platform.sh"
const HTTP_USER = "x"
const HTTP_PASSWORD = "x"

var file *os.File
var writer csv.Writer

func main() {
	file, err := os.Create(RESULT_FILE)
	checkError("Cannot create file", err)

	writer := csv.NewWriter(file)

	csvFile, _ := os.Open(FILE)
	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.LazyQuotes = true

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}

		queryUri(writer, line)

	}

}

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func queryUri(writer *csv.Writer, record []string) {

	uri := BASE_PATH + record[0]

	client := &http.Client{}

	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth(HTTP_USER, HTTP_PASSWORD))
	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode != 200 {
		line := []string{record[0], record[1], strconv.Itoa(response.StatusCode)}
		fmt.Println(line)
		writer.Write(line)
		writer.Flush()
	}

}
