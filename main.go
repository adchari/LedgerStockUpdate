package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Quote struct {
	Data struct {
		Price string `json:"05. price"`
	} `json:"Global Quote"`
}

func main() {
	apiToken := flag.String("a", "demo", "Alpha Vantage API Token")
	ledgerBinary := flag.String("b", "ledger", "Ledger Binary")
	ledgerFile := flag.String("f", "ledger.ledger", "Ledger File")
	priceDbFile := flag.String("p", "prices.db", "Price Database File")
	flag.Parse()

	commodities := GetCommodities(*ledgerFile, *ledgerBinary)

	pricedb, err := os.OpenFile(*priceDbFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Price database file access failed with %s\n", err)
	}
	defer pricedb.Close()

	start := time.Now()
	for i, c := range commodities {
		if i%5 == 0 {
			elapsed := time.Now().Sub(start)
			if elapsed < time.Minute {
				time.Sleep(time.Minute - elapsed)
			}
			start = time.Now()
		}

		priceString, err := GetPriceString(c, *apiToken)
		if err != nil {
			log.Println("Skipped " + c)
			continue
		}
		pricedb.WriteString("P " + GetTimeString() + " " + c + " " + priceString[:len(priceString)-2] + "\n")
	}
	log.Println("Stock price update complete")
}

func GetPriceString(ticker string, apiToken string) (string, error) {
	resp, err := http.Get("https://www.alphavantage.co/query?function=GLOBAL_QUOTE&symbol=" + ticker + "&apikey=" + apiToken)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var f Quote
	err = json.Unmarshal(body, &f)
	if err != nil {
		return "", err
	}

	if f.Data.Price == "" {
		return "", errors.New("Conversion Error")
	}
	return "$" + f.Data.Price, nil
}

func GetTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetCommodities(ledger string, binary string) []string {
	cmd := exec.Command(binary, "-f", ledger, "commodities")
	out, err := cmd.Output()
	if err != nil {
		log.Fatalf("Ledger file commodity report failed with %s\n", err)
	}
	a := strings.Split(string(out), "\n")
	sliceOut := a[:len(a)-1]

	commodities := make([]string, 0)
	for _, e := range sliceOut {
		if IsTicker(e) {
			commodities = append(commodities, e)
		}
	}
	return commodities
}

func IsTicker(s string) bool {
	for _, e := range s {
		if (e < 'A' || e > 'Z') && (e < '0' || e > '9') {
			return false
		}
	}
	return true
}
