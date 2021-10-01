package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/wspirrat/trashkv/core"
	"golang.org/x/sync/syncmap"
)

var (
	db syncmap.Map
)

func main() {
	port := os.Getenv("PORT")

	db = syncmap.Map{}

	http.HandleFunc("/connect", connect)
	http.HandleFunc("/save", compare_and_save)

	http.ListenAndServe(":"+port, nil)
}

func connect(w http.ResponseWriter, r *http.Request) {
	dataMap := make(map[string]interface{})
	db.Range(func(k interface{}, v interface{}) bool {
		dataMap[k.(string)] = v
		return true
	})

	j, err := json.Marshal(&dataMap)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Fprint(w, string(j))
}

func compare_and_save(w http.ResponseWriter, r *http.Request) {
	var request map[string]interface{}
	var newdb syncmap.Map

	// Try to decode the request body into the struct. If there is an error,
	// respond to the client with the error message and a 400 status code.
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		fmt.Println(err)
	}

	// check if request is not nil
	if len(request) > 0 {
		for key, value := range request {
			newdb.Store(key, value)
		}

		db = newdb
	}
}

// file functions
func NewKeySaveToFile(key string, value interface{}) error {
	var backupArray []string

	// string value is encoded to base64
	// because many text messages can contain commas
	// and text is encoded to prevent wrong string splits
	if reflect.TypeOf(value).String() == "string" {
		valueInString := fmt.Sprintf("%v", value)
		value = toBase64(valueInString)
		// the same with []string
	} else if reflect.TypeOf(value).String() == "[]string" {
		switch reflect.TypeOf(value).Kind() {
		case reflect.Slice:
			s := reflect.ValueOf(value)

			for i := 0; i < s.Len(); i++ {
				backupArray = append(backupArray, toBase64(s.Index(i).String()))
				value = backupArray
			}
		}
	}

	if err := SaveHardData(key, value); err != nil {
		return err
	}

	return nil
}

func SaveHardData(key string, value interface{}) error {
	db_line := fmt.Sprintf("%s,%v,%s;", key, value, reflect.TypeOf(value))
	// if DatabasePath exist
	if _, err := os.Stat(core.DatabasePath); err == nil {
		input, err := ioutil.ReadFile(core.DatabasePath)
		if err != nil {
			return err
		}

		lines := strings.Split(string(input), "\n")

		var isFind bool
		for i, line := range lines {
			split := strings.Split(line, ",")
			if split[0] == key {
				lines[i] = db_line
				isFind = true
			}
		}
		if !isFind {
			lines = append(lines, db_line)
		}

		output := strings.Join(lines, "\n")
		err = ioutil.WriteFile(core.DatabasePath, []byte(output), 0644)
		if err != nil {
			return err
		}
	} else {
		err = ioutil.WriteFile(core.DatabasePath, []byte(db_line), 0644)
		if err != nil {
			return err
		}
	}

	return nil
}

func readHardData(PrivateKey string) (syncmap.Map, error) {
	var res syncmap.Map

	if _, err := os.Stat(core.DatabasePath); err == nil {
		input, err := ioutil.ReadFile(core.DatabasePath)
		if err != nil {
			return syncmap.Map{}, err
		}

		splittedInput := strings.Split(string(input), ";")
		for _, value := range splittedInput {
			insideData := strings.Split(value, ",")
			if len(insideData) > 2 {
				converted, err := ConvertDataTo(insideData[2], insideData[1])
				if err != nil {
					return syncmap.Map{}, err
				}

				// I am trim spacing keys here, because of empty space line in db file
				// between keys
				res.Store(strings.TrimSpace(insideData[0]), converted)
			}
		}
	}

	return res, nil
}

func delInFile(key string) {
	input, err := ioutil.ReadFile(core.DatabasePath)
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		split := strings.Split(line, ",")
		if split[0] == key {
			lines[i] = ""
		}
	}
}

func ConvertDataTo(what, value string) (interface{}, error) {
	what = strings.TrimSpace(what)

	if what == "string" {
		return base64ToString(value), nil
	} else if what == "int" {
		intVar, err := strconv.Atoi(value)
		return intVar, err
	} else if what == "byte" {
		return []byte(value), nil
	} else if what == "[]string" {
		value = strings.Trim(value, "[ ]")
		split := strings.Split(value, " ")
		for index, val := range split {
			split[index] = base64ToString(strings.TrimSpace(val))
		}
		return split, nil
	} else if what == "bool" {
		boolType, err := strconv.ParseBool(value)
		return boolType, err
	}

	return value, nil
}

// utils
func toBase64(value string) string {
	return base64.StdEncoding.EncodeToString(
		[]byte(value),
	)
}

func base64ToString(value string) string {
	res, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		fmt.Println(err)
		return "[error]"
	}
	return string(res)
}
