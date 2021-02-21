package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mrjones/oauth"
)

//Usage How to user
func Usage() {
	fmt.Println("Usage:")
	fmt.Print("go run etsy-orders")
	fmt.Print("  --consumerkey <consumerkey>")
	fmt.Println("  --consumersecret <consumersecret>")
	fmt.Println("")
	fmt.Println("In order to get your consumerkey and consumersecret, you must register an 'app' etsy.com:")
}

func main() {
	var consumerKey *string = flag.String(
		"consumerkey",
		"",
		"Consumer Key from Etsy")

	var consumerSecret *string = flag.String(
		"consumersecret",
		"",
		"Consumer Secret Etsy")

	flag.Parse()

	if len(*consumerKey) == 0 || len(*consumerSecret) == 0 {
		fmt.Println("You must set the --consumerkey and --consumersecret flags.")
		fmt.Println("---")
		Usage()
		os.Exit(1)
	}

	c := oauth.NewConsumer(
		*consumerKey,
		*consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://openapi.etsy.com/v2/oauth/request_token?scope=transactions_w",
			AuthorizeTokenUrl: "https://openapi.etsy.com/v2/oauth/authorize",
			AccessTokenUrl:    "https://openapi.etsy.com/v2/oauth/access_token",
		})

	c.Debug(true)

	token := oauth.AccessToken{Token: "<access_token>", Secret: "<secret>"}

	client, err := c.MakeHttpClient(&token)
	if err != nil {
		log.Fatal(err)
	}

	GetOrders(client)
	// defer response.Body.Close()

	// bits, err := ioutil.ReadAll(response.Body)
	// fmt.Println("The newest item in your home timeline is: " + string(bits))

	// if *postUpdate {
	// 	status := fmt.Sprintf("Test post via the API using Go (http://golang.org/) at %s", time.Now().String())

	// 	response, err = client.PostForm(
	// 		"https://api.twitter.com/1.1/statuses/update.json",
	// 		url.Values{"status": []string{status}})

	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}

	// 	log.Printf("%v\n", response)
	// }
}
