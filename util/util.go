package util

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"raygun/log"
	"strings"
)

func ReadFile(filename string) (*string, error) {
	// Read the file into a []byte slice
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Convert the []byte slice to a string
	str := string(data)

	return &str, nil
}

func SplitHeaderAndBody(data string) (*string, *string, error) {

	var header string
	var body string = ""

	log.Debug("SplitHeaderAndBody: processing: %s", data)

	scanner := bufio.NewScanner(strings.NewReader(data))

	if !scanner.Scan() {
		return nil, nil, errors.New("SplitHeaderAndbody: invalid header/body pair:" + data)
	}

	header = scanner.Text()

	header = strings.TrimSpace(header)

	// we want to keep
	for scanner.Scan() {
		body += scanner.Text() + "\n"
	}

	return &header, &body, nil
}

func Chomp(data string) string {

	data = strings.Trim(data, " \n\t\r")

	return data
}

func Listify(data string) []string {

	list := make([]string, 0)

	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {

		text := strings.TrimSpace(scanner.Text())

		if len(text) > 0 {
			list = append(list, strings.TrimSpace(scanner.Text()))
		}
	}

	return list

}

func GetFileExtension(entity string) string {

	return filepath.Ext(entity)

}
