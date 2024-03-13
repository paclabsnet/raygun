package util

import (
	"bufio"
	"path/filepath"
	"sort"
	"strings"
)

/*
 *  sorts the keys of a map so we always get them in the same order
 *  from run to run.  This helps with validation and debugging
 *  and operational cleanliness
 */
func SortMapKeys(data map[string]interface{}) []string {
	keys := make([]string, len(data))
	i := 0
	for k := range data {
		keys[i] = k
		i++
	}

	sort.Strings(keys)

	return keys
}

// func ReadFile(filename string) (*string, error) {
// 	// Read the file into a []byte slice
// 	data, err := os.ReadFile(filename)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Convert the []byte slice to a string
// 	str := string(data)

// 	return &str, nil
// }

// func SplitHeaderAndBody(data string) (*string, *string, error) {

// 	var header string
// 	var body string = ""

// 	log.Debug("SplitHeaderAndBody: processing: %s", data)

// 	scanner := bufio.NewScanner(strings.NewReader(data))

// 	if !scanner.Scan() {
// 		return nil, nil, errors.New("SplitHeaderAndbody: invalid header/body pair:" + data)
// 	}

// 	header = scanner.Text()

// 	header = strings.TrimSpace(header)

// 	// we want to keep
// 	for scanner.Scan() {
// 		body += scanner.Text() + "\n"
// 	}

// 	return &header, &body, nil
// }

// func Chomp(data string) string {

// 	data = strings.Trim(data, " \n\t\r")

// 	return data
// }

/*
 * Take a string that has newlines, and convert each line into a separate array element
 * in a list of strings
 */
func Listify(data string) []string {

	list := make([]string, 0)

	scanner := bufio.NewScanner(strings.NewReader(data))

	for scanner.Scan() {

		text := strings.TrimSpace(scanner.Text())

		if len(text) > 0 {
			list = append(list, text)
		}
	}

	return list

}

func GetFileExtension(entity string) string {

	return filepath.Ext(entity)

}

/*
 * returns true if the object is a string
 */
func IsString(obj interface{}) bool {

	_, ok := obj.(string)

	return ok
}

/*
 * returns true if the object is a map, and it can be three
 * different types of map - a map of itnerfaces or a map of strings
 */
func IsMap(obj interface{}) bool {

	_, ok := obj.(map[string]interface{})

	if ok {
		return true
	}

	_, ok = obj.(map[string]string)

	return ok
}

/*
 * returns true if the object is an array of interfaces
 */
func IsArray(obj interface{}) bool {
	_, ok := obj.([]interface{})

	return ok
}