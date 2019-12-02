package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {

	jsonFile, err := os.Open("1575245198.json")

	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	batchRSP := result["batchrsp"].(map[string]interface{})

	for key, value := range batchRSP {
		fmt.Println(key, value.(string))
	}
}
