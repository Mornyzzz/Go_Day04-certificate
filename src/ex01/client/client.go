package main

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

type CandyResponse struct {
	Change int    `json:"change"`
	Thanks string `json:"thanks"`
}

func getClient() *http.Client {
	data, err := os.ReadFile("../cert/minica.pem")
	if err != nil {
		fmt.Println(err)
	}
	cp, _ := x509.SystemCertPool()
	cp.AppendCertsFromPEM(data)

	cert, err := tls.LoadX509KeyPair("../cert/client/cert.pem", "../cert/client/key.pem")
	if err != nil {
		fmt.Println(err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      cp,
	}
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: config,
		},
	}
	return client
}

func main() {
	kFlag := flag.String("k", "", "candy type")
	cFlag := flag.Int("c", -1, "candy count")
	mFlag := flag.Int("m", -1, "money")
	flag.Parse()
	if *kFlag == "" || *cFlag == -1 || *mFlag == -1 {
		fmt.Println("no valid data")
		return
	}
	client := getClient()
	data := map[string]interface{}{
		"money":      *mFlag,
		"candyType":  *kFlag,
		"candyCount": *cFlag,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}
	resp, err := client.Post("https://localhost:3333/buy_candy", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	if resp.StatusCode != 201 {
		fmt.Print(string(body))
	} else {
		var response CandyResponse
		err = json.Unmarshal(body, &response)
		if err != nil {
			fmt.Println("Ошибка при декодировании JSON:", err)
			return
		}
		fmt.Println("Thank you! Your change is " + strconv.Itoa(response.Change))
	}
}
