package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func main() {
	// metadata := new(spotlight.MetaData)
	// metadata.ImportAll()

	// panic("ABORTED")

	data, err := readJSON("image.json")
	if err != nil {
		panic(err)
	}

	var result map[string]interface{}
	json.Unmarshal(data, &result)

	batchRSP := result["batchrsp"].(map[string]interface{})
	items := batchRSP["items"].([]interface{})

	for idx, obj := range items {
		fmt.Printf("====================\nItem %d:\n", idx)
		itemObj := obj.(map[string]interface{})
		itemStr := itemObj["item"]
		itemBytes := []byte(itemStr.(string))

		var itemMap map[string]interface{}
		json.Unmarshal(itemBytes, &itemMap)

		adObj := itemMap["ad"].(map[string]interface{})
		fmt.Printf("'ad' subMap:\n\n")
		for key, val := range adObj {
			fmt.Printf("%s: %v\n\n", key, val)
		}
	}
}

func readJSON(fileName string) ([]byte, error) {
	data, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	if !json.Valid(data) {
		return nil, fmt.Errorf("Specified file does not contain valid JSON")
	}

	return data, nil
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func readFile(filePath string) []byte {
	data, err := ioutil.ReadFile(filePath)
	checkErr(err)
	return data
}

// func main() {
// 	path := "1575245198.json"
// 	data := readFile(path)

// 	var f interface{}

// 	err := json.Unmarshal(data, &f)
// 	checkErr(err)

// 	m := f.(map[string]interface{})

// 	for k, v := range m {
// 		switch vv := v.(type) {
// 		case string:
// 			fmt.Println(k, "is string", vv)
// 		case int:
// 			fmt.Println(k, "is int", vv)
// 		case []interface{}:
// 			fmt.Println(k, "is an array:")
// 			for i, u := range vv {
// 				fmt.Println(i, u)
// 			}
// 		default:
// 			fmt.Println(k, "is of a type I don't know how to handle")
// 		}
// 	}
// }
