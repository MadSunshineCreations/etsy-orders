package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
)

//Orders - list of orders currently open in shop
type Orders struct {
	Count   int `json:"count"`
	Results []struct {
		ReceiptID         int    `json:"receipt_id"`
		ReceiptType       int    `json:"receipt_type"`
		OrderID           int    `json:"order_id"`
		CanRefund         bool   `json:"can_refund"`
		MessageFromSeller string `json:"message_from_seller"`
		MessageFromBuyer  string `json:"message_from_buyer"`
		WasPaid           bool   `json:"was_paid"`
		GiftMessage       string `json:"gift_message"`
		IsOverdue         bool   `json:"is_overdue"`
		DaysFromDueDate   int    `json:"days_from_due_date"`
		IsDead            bool   `json:"is_dead"`
	} `json:"results"`
	Pagination struct {
		EffectiveLimit  int `json:"effective_limit"`
		EffectiveOffset int `json:"effective_offset"`
		NextOffset      int `json:"next_offset"`
		EffectivePage   int `json:"effective_page"`
		NextPage        int `json:"next_page"`
	} `json:"pagination"`
}

//Receipts - list of receipts for an order. Basically the order details
type Receipts struct {
	Count   int `json:"count"`
	Results []struct {
		TransactionID  int64  `json:"transaction_id"`
		Title          string `json:"title"`
		URL            string `json:"url"`
		Quantity       int    `json:"quantity"`
		ListingID      int    `json:"listing_id"`
		ImageListingID int    `json:"image_listing_id"`
		Variations     []struct {
			PropertyID     int    `json:"property_id"`
			ValueID        int64  `json:"value_id"`
			FormattedName  string `json:"formatted_name"`
			FormattedValue string `json:"formatted_value"`
		} `json:"variations"`
	} `json:"results"`
	Params struct {
		ReceiptID string `json:"receipt_id"`
		Limit     int    `json:"limit"`
		Offset    int    `json:"offset"`
		Page      int    `json:"page"`
	} `json:"params"`
	Type       string `json:"type"`
	Pagination struct {
		EffectiveLimit  int `json:"effective_limit"`
		EffectiveOffset int `json:"effective_offset"`
		NextOffset      int `json:"next_offset"`
		EffectivePage   int `json:"effective_page"`
		NextPage        int `json:"next_page"`
	} `json:"pagination"`
}

//An Order is consolidation of an Order and Receipt from Etsy
type Order struct {
	Name              string `json:"name"`
	PrimaryColor      string `json:"primary_color"`
	SecondaryColor    string `json:"secondary_color"`
	SlotAmount        string `json:"slot_amount"`
	Quantity          int    `json:"quantity"`
	MessageFromSeller string `json:"message_from_seller"`
	GiftMessage       string `json:"gift_message"`
	URL               string `json:"url"`
	DaysFromDueDate   int    `json:"days_from_due_date"`
	ImageURL          string `json:"image_url"`
}

//Images is a list of URL for a listing Image
type Images struct {
	Count   int `json:"count"`
	Results []struct {
		SmallURL string `json:"url_75x75"`
	}
}

// GetOrders expects a http client setup with auth ready to go.
func GetOrders(client *http.Client) []Order {
	response, err := client.Get(
		fmt.Sprintf("https://openapi.etsy.com/v2/shops/%v/receipts?was_paid=true&was_shipped=false", config.ShopID))
	if err != nil {
		log.Fatal(err)
	}
	if response.StatusCode == 400 {
		// Probably rate limited
		return make([]Order, 0)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	orders := Orders{}
	if err := json.Unmarshal(bytes, &orders); err != nil {
		panic(err)
	}
	var openOrders []Order

	for i := 0; i < len(orders.Results); i++ {
		receipts := getTransactions(client, orders.Results[i].ReceiptID)

		for j := 0; j < receipts.Count; j++ {
			primaryColor := ""
			secondaryColor := ""
			slotAmount := ""
			for k := 0; k < len(receipts.Results[j].Variations); k++ {
				switch variation := receipts.Results[j].Variations[k].FormattedName; variation {
				case "Complete Comfort Grip Color", "Color":
					primaryColor = receipts.Results[j].Variations[k].FormattedValue
				case "Secondary color":
					secondaryColor = receipts.Results[j].Variations[k].FormattedValue
				case "Slot Amount":
					slotAmount = receipts.Results[j].Variations[k].FormattedValue
				}
			}

			imageURL := getListingImage(client, receipts.Results[j].ListingID, receipts.Results[j].ImageListingID)

			//Make a new order for each
			order := Order{
				receipts.Results[j].Title,
				primaryColor,
				secondaryColor,
				slotAmount,
				receipts.Results[j].Quantity,
				orders.Results[i].MessageFromBuyer,
				orders.Results[i].GiftMessage,
				receipts.Results[j].URL,
				orders.Results[i].DaysFromDueDate,
				imageURL,
			}
			openOrders = append(openOrders, order)
		}
	}
	sort.Slice(openOrders, func(i, j int) bool {
		return openOrders[i].DaysFromDueDate < openOrders[j].DaysFromDueDate
	})
	return openOrders
}

func getTransactions(client *http.Client, recieptID int) Receipts {
	response, err := client.Get(
		fmt.Sprintf("https://openapi.etsy.com/v2/receipts/%v/transactions", recieptID))
	if err != nil {
		log.Fatal(err)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	receipts := Receipts{}
	if err := json.Unmarshal(bytes, &receipts); err != nil {
		panic(err)
	}
	return receipts
}

func getListingImage(client *http.Client, listingID int, listingImageID int) string {
	response, err := client.Get(
		fmt.Sprintf("https://openapi.etsy.com/v2/listings/%v/images/%v", listingID, listingImageID))
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(response.Body)
	images := Images{}
	if err := json.Unmarshal(bytes, &images); err != nil {
		panic(err)
	}
	if images.Count > 0 {
		return images.Results[0].SmallURL
	}
	return ""
}
