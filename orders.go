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
		ReceiptID                int           `json:"receipt_id"`
		ReceiptType              int           `json:"receipt_type"`
		OrderID                  int           `json:"order_id"`
		SellerUserID             int           `json:"seller_user_id"`
		BuyerUserID              int           `json:"buyer_user_id"`
		CreationTsz              int           `json:"creation_tsz"`
		CanRefund                bool          `json:"can_refund"`
		LastModifiedTsz          int           `json:"last_modified_tsz"`
		Name                     string        `json:"name"`
		FirstLine                string        `json:"first_line"`
		SecondLine               string        `json:"second_line"`
		City                     string        `json:"city"`
		State                    string        `json:"state"`
		Zip                      string        `json:"zip"`
		FormattedAddress         string        `json:"formatted_address"`
		CountryID                int           `json:"country_id"`
		PaymentMethod            string        `json:"payment_method"`
		PaymentEmail             string        `json:"payment_email"`
		MessageFromSeller        interface{}   `json:"message_from_seller"`
		MessageFromBuyer         interface{}   `json:"message_from_buyer"`
		WasPaid                  bool          `json:"was_paid"`
		TotalTaxCost             string        `json:"total_tax_cost"`
		TotalVatCost             string        `json:"total_vat_cost"`
		TotalPrice               string        `json:"total_price"`
		TotalShippingCost        string        `json:"total_shipping_cost"`
		CurrencyCode             string        `json:"currency_code"`
		MessageFromPayment       interface{}   `json:"message_from_payment"`
		WasShipped               bool          `json:"was_shipped"`
		BuyerEmail               string        `json:"buyer_email"`
		SellerEmail              string        `json:"seller_email"`
		IsGift                   bool          `json:"is_gift"`
		NeedsGiftWrap            bool          `json:"needs_gift_wrap"`
		GiftMessage              string        `json:"gift_message"`
		DiscountAmt              string        `json:"discount_amt"`
		Subtotal                 string        `json:"subtotal"`
		Grandtotal               string        `json:"grandtotal"`
		AdjustedGrandtotal       string        `json:"adjusted_grandtotal"`
		BuyerAdjustedGrandtotal  string        `json:"buyer_adjusted_grandtotal"`
		Shipments                []interface{} `json:"shipments"`
		ShippedDate              int           `json:"shipped_date"`
		IsOverdue                bool          `json:"is_overdue"`
		DaysFromDueDate          int           `json:"days_from_due_date"`
		TransparentPriceMessage  string        `json:"transparent_price_message"`
		ShowChannelBadge         bool          `json:"show_channel_badge"`
		ChannelBadgeSuffixString string        `json:"channel_badge_suffix_string"`
		IsDead                   bool          `json:"is_dead"`
	} `json:"results"`
	Pagination struct {
		EffectiveLimit  int         `json:"effective_limit"`
		EffectiveOffset int         `json:"effective_offset"`
		NextOffset      interface{} `json:"next_offset"`
		EffectivePage   int         `json:"effective_page"`
		NextPage        interface{} `json:"next_page"`
	} `json:"pagination"`
}

// GroupedOrders is the list
type GroupedOrders struct {
	Count  int `json:"count"`
	ByItem struct {
		name   string
		colors []string
	} `json:"orders"`
	ByColor struct {
		name  string
		count int
	}
}

// GetOrders expects a http client setup with auth ready to go.
func GetOrders(client *http.Client) {
	response, err := client.Get(
		"https://openapi.etsy.com/v2/shops/15212421/receipts?was_paid=true&was_shipped=false")
	if err != nil {
		log.Fatal(err)
	}

	bytes, err := ioutil.ReadAll(response.Body)
	// fmt.Println(string(bits))

	orders := Orders{}
	if err := json.Unmarshal(bytes, &orders); err != nil {
		panic(err)
	}
	fmt.Println(orders)

	for i := 0; i < len(orders.Results); i++ {
		getTransactions(client, orders.Results[i].ReceiptID)
	}
}

func getTransactions(client *http.Client, recieptID int) {
	response, err := client.Get(
		"https://openapi.etsy.com/v2/receipts/" + fmt.Sprint(recieptID) + "/transactions")
	if err != nil {
		log.Fatal(err)
	}

	// Need to get all the transactions (Items). Ulimately need to gruop this by property ID and Value
	// Property is like "Color, Secondary Color", value is the actual color. Blue, Red etc..

	bytes, err := ioutil.ReadAll(response.Body)

	fmt.Println(string(bytes))

}
