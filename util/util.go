package util

import (
	"bufio"
	"fmt"
	"os"
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

/*
 * Read the contents of a file into a string
 */
func ReadFile(current_path string, target_file string) (string, error) {

	filename := filepath.Join(current_path, target_file)

	fileBytes, err := os.ReadFile(filename)

	if err != nil {
		fmt.Printf("ERROR: can't open file for reading: %s -> %s", filename, err.Error())
		return "", err
	}

	return string(fileBytes), nil

}

/*
 * Strings all whiespace from a string
 */
func RemoveAllWhitespace(str string) string {
	return strings.ReplaceAll(str, " ", "")
}
