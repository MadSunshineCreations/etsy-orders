package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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
		TransactionID int64  `json:"transaction_id"`
		Title         string `json:"title"`
		URL           string `json:"url"`
		Variations    []struct {
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
	Name              string
	PrimaryColor      string
	SecondaryColor    string
	SlotAmount        string
	MessageFromSeller string
	GiftMessage       string
	URL               string
}

// GetOrders expects a http client setup with auth ready to go.
func GetOrders(client *http.Client) []Order {
	response, err := client.Get(
		"https://openapi.etsy.com/v2/shops/15212421/receipts?was_paid=true&was_shipped=false")
	if err != nil {
		log.Fatal(err)
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
				case "Color":
					primaryColor = receipts.Results[j].Variations[k].FormattedValue
				case "Secondary color":
					secondaryColor = receipts.Results[j].Variations[k].FormattedValue
				case "Slot Amount":
					slotAmount = receipts.Results[j].Variations[k].FormattedValue
				}
			}

			//Make a new order for each
			order := Order{
				receipts.Results[j].Title,
				primaryColor,
				secondaryColor,
				slotAmount,
				orders.Results[i].MessageFromBuyer,
				orders.Results[i].GiftMessage,
				receipts.Results[j].URL,
			}
			openOrders = append(openOrders, order)
		}
	}
	return openOrders
}

func getTransactions(client *http.Client, recieptID int) Receipts {
	response, err := client.Get(
		"https://openapi.etsy.com/v2/receipts/" + fmt.Sprint(recieptID) + "/transactions")
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
