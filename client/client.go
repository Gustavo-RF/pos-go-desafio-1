package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type UsdValueResponseDto struct {
	Bid string `json:"bid"`
}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 300*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "http://localhost:8080/cotacao", nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Accepts", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	res, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var data UsdValueResponseDto
	err = json.Unmarshal(res, &data)

	if err != nil {
		panic(err)
	}

	fmt.Printf("Valor atual do câmbio: R$ %v\n", data.Bid)

	SaveFile(data)
}

func SaveFile(data UsdValueResponseDto) error {
	fileOpened, err := os.OpenFile("cotacao.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if os.IsNotExist(err) {
		fileOpened, err = os.Create("cotacao.txt")
		if err != nil {
			panic(err)
		}
	}

	defer fileOpened.Close()

	_, err = fileOpened.WriteString("Dolar: " + data.Bid + "\n")
	if err != nil {
		panic(err)
	}
	return nil
}
