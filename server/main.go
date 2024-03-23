package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

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

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 200*time.Millisecond)
	defer cancel()
	usdValue, err := FetchUsdValue(ctx)

	if ctx.Err() != nil {
		fmt.Printf("Error while request: %v", ctx.Err())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	if err != nil {
		fmt.Printf("Error while request: %v", err)
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	ctx, cancel = context.WithTimeout(ctx, 10*time.Millisecond)
	defer cancel()
	CreateCurrencyDb(*usdValue, *db, ctx)

	if ctx.Err() != nil {
		fmt.Printf("Error while save: %v", ctx.Err())
		res.WriteHeader(http.StatusBadRequest)
		return
	}

	res.Header().Set("Content-type", "application/json")
	json.NewEncoder(res).Encode(usdValue)
}

func FetchUsdValue(ctx context.Context) (*UsdValueDto, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accepts", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
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

func CreateCurrencyDb(UsdValueDto UsdValueDto, db gorm.DB, ctx context.Context) {
	db.WithContext(ctx).Create(&Currency{
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
