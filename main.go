package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/mrjones/oauth"
)

var ordersMutex = &sync.Mutex{}
var orders = []Order{}

type ordersReply struct {
	Orders []Order `json:"orders"`
}

type etsyConfig struct {
	ConsumerKey    string
	ConsumerSecret string
	AccessToken    string
	AccessSecret   string
}

var config etsyConfig

//Usage How to user
func Usage() {
	fmt.Println("Usage:")
	fmt.Print("go run etsy-orders")
	fmt.Print("  --consumerkey <consumerkey>")
	fmt.Println("  --consumersecret <consumersecret>")
	fmt.Println("")
	fmt.Println("In order to get your consumerkey and consumersecret, you must register an 'app' etsy.com:")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Alive!")
}

func orderHandler(w http.ResponseWriter, r *http.Request) {
	ordersMutex.Lock()
	json.NewEncoder(w).Encode(orders)
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

	fmt.Printf("Listening on port 8081")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func loadConfig() {
	var consumerKey *string = flag.String(
		"consumerkey",
		"",
		"Consumer Key from Etsy")

	var consumerSecret *string = flag.String(
		"consumersecret",
		"",
		"Consumer Secret Etsy")

	var accessToken *string = flag.String(
		"accesstoken",
		"",
		"Access Token for Etsy App")

	var accessSecret *string = flag.String(
		"accesssecret",
		"",
		"Access Secret Etsy App")

	flag.Parse()

	if len(*consumerKey) == 0 || len(*consumerSecret) == 0 {
		fmt.Println("You must set the --consumerkey and --consumersecret flags.")
		fmt.Println("---")
		Usage()
		os.Exit(1)
	}

	config.ConsumerKey = *consumerKey
	config.ConsumerSecret = *consumerSecret
	config.AccessSecret = *accessSecret
	config.AccessToken = *accessToken
}

func loadOrders() {
	c := oauth.NewConsumer(
		config.ConsumerKey,
		config.ConsumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://openapi.etsy.com/v2/oauth/request_token?scope=transactions_w",
			AuthorizeTokenUrl: "https://openapi.etsy.com/v2/oauth/authorize",
			AccessTokenUrl:    "https://openapi.etsy.com/v2/oauth/access_token",
		})

	c.Debug(true)

	token := oauth.AccessToken{Token: config.AccessToken, Secret: config.AccessSecret}

	client, err := c.MakeHttpClient(&token)
	if err != nil {
		log.Fatal(err)
	}

	openOrders := GetOrders(client)
	ordersMutex.Lock()
	orders = openOrders
	ordersMutex.Unlock()
}
