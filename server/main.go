package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Usdbrl struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type UsdValue struct {
	Usdbrl Usdbrl `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/", HandleFetchUsdValue)
	http.ListenAndServe(":8080", nil)
}

func HandleFetchUsdValue(res http.ResponseWriter, req *http.Request) {

	if req.URL.Path != "/" {
		res.WriteHeader(http.StatusNotFound)
		return
	}

	usdValue, err := FetchUsdValue()

	if err != nil {
		fmt.Println(err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(usdValue)
}

func FetchUsdValue() (*UsdValue, error) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var data UsdValue
	err = json.Unmarshal(res, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}
