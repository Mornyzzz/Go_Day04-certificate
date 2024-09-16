package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	_ "log"
	"net/http"
	"strconv"
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
	fmt.Println("handle request")
	err := json.NewDecoder(r.Body).Decode(&order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // 400
		return
	}

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
		response := map[string]interface{}{
			"thanks": "Thank you!",
			"change": order.Money - sumOrder,
		}
		w.WriteHeader(http.StatusCreated) // 201
		json.NewEncoder(w).Encode(response)
		return
	}
}

func getServer() *http.Server {
	server := &http.Server{
		Addr: "localhost:3333",
	}
	return server
}

func main() {
	server := getServer()
	http.HandleFunc("/buy_candy", buyCandy)
	fmt.Println("Starting server on port :3333...")

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
