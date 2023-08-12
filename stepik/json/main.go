package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type OkvedList []Okved

type Okved struct {
	GlobalID int `json:"global_id"`
	//	SystemObjectID string `json:"system_object_id"`
	//	SignatureDate  string `json:"signature_date"`
	//	Razdel         string `json:"Razdel"`
	//	Kod            string `json:"Kod,omitempty"`
	//	Name           string `json:"Name"`
	//	Idx            string `json:"Idx"`
	//	Nomdescr       string `json:"Nomdescr,omitempty"`
}

func main() {
	filename := "./data.json"
	data, err := ioutil.ReadFile(filename)

	okvedlist := OkvedList{}

	err = json.Unmarshal(data, &okvedlist)
	if err != nil {
		panic(err)
	}

	sum := 0
	for _, o := range okvedlist {
		sum += o.GlobalID
	}

	fmt.Println(sum)
}
