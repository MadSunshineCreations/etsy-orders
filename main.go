package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mrjones/oauth"
)

var ordersMutex = &sync.Mutex{}

type ordersReply struct {
	Orders []Order `json:"orders"`
}

type etsyConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
	ShopID         string
}

var config etsyConfig
var reply ordersReply

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Alive!")
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	ordersMutex.Lock()
	json.NewEncoder(w).Encode(reply)
	ordersMutex.Unlock()
}

func main() {
	loadConfig()
	go func() {
		for {
			loadOrders()
			time.Sleep(time.Minute * 15)
		}
	}()

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/orders", orderHandler)

	log.Printf("Listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func loadConfig() {
	consumerKey := os.Getenv("ETSY_ORDERS_CONSUMER_KEY")
	consumerSecret := os.Getenv("ETSY_ORDERS_CONSUMER_SECRET")
	accessToken := os.Getenv("ETSY_ORDERS_ACCESS_KEY")
	accessSecret := os.Getenv("ETSY_ORDERS_ACCESS_SECRET")
	shopID := os.Getenv("ETSY_ORDERS_SHOP_ID")

	if len(consumerKey) == 0 || len(consumerSecret) == 0 {
		fmt.Println("You must set the ETSY_ORDERS_CONSUMER_KEY and ETSY_ORDERS_CONSUMER_SECRET Environment Variables.")
		fmt.Println("---")
		os.Exit(1)
	}

	config.ConsumerKey = consumerKey
	config.ConsumerSecret = consumerSecret
	config.AccessSecret = accessSecret
	config.AccessToken = accessToken
	config.ShopID = shopID
}

func loadOrders() {
	log.Println("Loading Orders")
	c := oauth.NewConsumer(
		config.ConsumerKey,
		config.ConsumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://openapi.etsy.com/v2/oauth/request_token?scope=transactions_w",
			AuthorizeTokenUrl: "https://openapi.etsy.com/v2/oauth/authorize",
			AccessTokenUrl:    "https://openapi.etsy.com/v2/oauth/access_token",
		})

	// c.Debug(true)

	token := oauth.AccessToken{Token: config.AccessToken, Secret: config.AccessSecret}

	client, err := c.MakeHttpClient(&token)
	if err != nil {
		log.Fatal(err)
	}

	openOrders := GetOrders(client)
	ordersMutex.Lock()
	reply.Orders = openOrders
	ordersMutex.Unlock()
	log.Printf("Done loading %v Orders", len(openOrders))
}
