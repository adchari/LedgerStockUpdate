package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"
)

func main() {
	ledgerFile := flag.String("f", "ledger.ledger", "Ledger File")
	priceDbFile := flag.String("p", "prices.db", "Price Database File")
	flag.Parse()

	commodities := GetCommodities(*ledgerFile)

	pricedb, err := os.OpenFile(*priceDbFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Price database file access failed with %s\n", err)
	}
	defer pricedb.Close()

	for _, c := range commodities {
		priceString, err := GetPriceString(c)
		if err != nil {
			continue
		}
		pricedb.WriteString("P " + GetTimeString() + " " + c + " " + priceString + "\n")
	}
}

func GetPriceString(ticker string) (string, error) {
	resp, err := http.Get("https://api.worldtradingdata.com/api/v1/stock?symbol=" + ticker + "&api_token=demo")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var f interface{}
	json.Unmarshal(body, &f)
	m := f.(map[string]interface{})
	dataInterface := m["data"]
	arr := dataInterface.([]interface{})
	elem := arr[0]
	ma := elem.(map[string]interface{})
	priceInterface := ma["price"]
	price := priceInterface.(string)
	return "$" + price, nil
}

func GetTimeString() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetCommodities(ledger string) []string {
	cmd := exec.Command("ledger", "-f", ledger, "commodities")
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
