package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type UsdbrlDto struct {
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

type Currency struct {
	CurrencyCompare string `json:"currency_compare"`
	CurrencyUser    string `json:"currency_user"`
	HighValue       string `json:"high_value"`
	LowValue        string `json:"low_value"`
	Variation       string `json:"variation"`
	Percentage      string `json:"percentage"`
	BuyValue        string `json:"buy_value"`
	SellValue       string `json:"sell_value"`
	CreateDate      string `json:"create_date"`
	gorm.Model
}

type UsdValueDto struct {
	Usdbrl UsdbrlDto `json:"USDBRL"`
}

func main() {

	db, err := gorm.Open(sqlite.Open("currency.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Migrate the schema
	db.AutoMigrate(&Currency{})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		HandleFetchUsdValue(w, r, db)
	})
	http.ListenAndServe(":8080", nil)
}

func HandleFetchUsdValue(res http.ResponseWriter, req *http.Request, db *gorm.DB) {

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

	CreateCurrencyDb(*usdValue, *db)

	// if err != nil {
	// 	fmt.Println(err)
	// 	res.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(usdValue)
}

func FetchUsdValue() (*UsdValueDto, error) {
	req, err := http.Get("https://economia.awesomeapi.com.br/json/last/USD-BRL")
	if err != nil {
		return nil, err
	}
	defer req.Body.Close()

	res, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var data UsdValueDto
	err = json.Unmarshal(res, &data)

	if err != nil {
		return nil, err
	}

	return &data, nil
}

func CreateCurrencyDb(UsdValueDto UsdValueDto, db gorm.DB) {
	db.Create(&Currency{
		CurrencyCompare: UsdValueDto.Usdbrl.Code,
		CurrencyUser:    UsdValueDto.Usdbrl.Codein,
		HighValue:       UsdValueDto.Usdbrl.High,
		LowValue:        UsdValueDto.Usdbrl.Low,
		Variation:       UsdValueDto.Usdbrl.VarBid,
		Percentage:      UsdValueDto.Usdbrl.PctChange,
		BuyValue:        UsdValueDto.Usdbrl.Bid,
		SellValue:       UsdValueDto.Usdbrl.Ask,
		CreateDate:      UsdValueDto.Usdbrl.CreateDate,
	})
}
