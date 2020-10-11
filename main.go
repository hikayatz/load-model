package main

import (
	"encoding/json"
	"fmt"
	. "github.com/zayinul/load-model/helpers"
	"net/http"
)

type Person struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
	Age     int64  `json:"age"`
}

// sample code
func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/person", handlePostPerson)

	fmt.Println("webserver start at http://localhost:3000")
	_ = http.ListenAndServe(":3000", nil)
}

func handlePostPerson(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(HeaderContentType, MIMEApplicationJSON)
	if request.Method != "POST" {
		http.Error(writer, "method allow only post", http.StatusBadRequest)
		return
	}
	var person = &Person{}
	err := LoadModel(person, request, "json")
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}
	jsonByte, err := json.Marshal(&person)
	_, _ = writer.Write(jsonByte)

}

func handleIndex(writer http.ResponseWriter, request *http.Request) {
	writer.Header().Set(HeaderContentType, MIMEApplicationJSON)
	_, _ = writer.Write([]byte(" -- welcome--"))
}
