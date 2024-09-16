package main

/*
#cgo CFLAGS: -I.
#include <stdlib.h>
#include <string.h>
#include <stdio.h>
char *ask_cow(char phrase[]);
*/
import "C"
import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"unsafe"
)

type CandyOrder struct {
	Money      int    `json:"money"`
	CandyType  string `json:"candyType"`
	CandyCount int    `json:"candyCount"`
}

func countMoney(order CandyOrder) (int, error) {
	if order.CandyCount < 0 {
		return 0, errors.New("candy count < 0")
	}
	switch order.CandyType {
	case "CE":
		return order.CandyCount * 10, nil
	case "AA":
		return order.CandyCount * 15, nil
	case "NT":
		return order.CandyCount * 17, nil
	case "DE":
		return order.CandyCount * 21, nil
	case "YR":
		return order.CandyCount * 23, nil
	default:
		return 0, errors.New("non-existent candy")
	}
}

func buyCandy(w http.ResponseWriter, r *http.Request) {
	var order CandyOrder
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}
	fmt.Println("handle request")

	sumOrder, err := countMoney(order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

	if sumOrder > order.Money {
		http.Error(w, "You need "+strconv.Itoa(sumOrder-order.Money)+" more money!", http.StatusPaymentRequired) // 402
		return
	}

	if sumOrder <= order.Money {

		cPhrase := C.CString(" Thank you!")
		defer C.free(unsafe.Pointer(cPhrase))

		result := C.ask_cow(cPhrase)
		goResult := C.GoString(result)

		response := map[string]interface{}{
			"thanks": goResult,
			"change": order.Money - sumOrder,
		}
		w.WriteHeader(http.StatusCreated) // 201
		json.NewEncoder(w).Encode(response)

		return
	}
}

func getServer() *http.Server {
	data, err := os.ReadFile("../cert/minica.pem")
	if err != nil {
		fmt.Println(err)
	}
	cp, _ := x509.SystemCertPool()
	cp.AppendCertsFromPEM(data)

	config := &tls.Config{
		ClientCAs:  cp,
		ClientAuth: tls.RequireAndVerifyClientCert,
	}

	server := &http.Server{
		Addr:      "localhost:3333",
		TLSConfig: config,
	}
	return server
}

func main() {
	server := getServer()
	http.HandleFunc("/buy_candy", buyCandy)
	fmt.Println("Starting server on port :3333...")

	err := server.ListenAndServeTLS("../cert/localhost/cert.pem", "../cert/localhost/key.pem")
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
