# Etsy Orders

This is a simple service that access the etsy API so that I can get a list of current open orders to be displayed in a user interface alongside my current running 3D Printers.

It runs a webservice with two endpoints.

`/` - Returns "Alive"

`/orders` - Returns a list of orders for the given shop. The list is refreshed every 15 minutes and kept cached. This is done to not bash the Etsy APIs if the api is polled or used frequently.

The JSON reply is very customized to my needs, even searching for specific variations and flattening these into keys in the JOSN structure. This could be made generic, but I have intentionally simplified for my needs.

## Configuration

In order for the application to run there must be 4 environment variables set.

The set of consumer key/secret is used to authenticate agains the Etsy API, the set of Access Key/Secret is used to identify your Etsy Application and used to grant access to the specific shop this will get data from.

The Consumer key/secret are what's given to you buy etsy when you sign up as a developer.

ETSY_ORDERS_CONSUMER_KEY=<consumer_key>

ETSY_ORDERS_CONSUMER_SECRET=<consumer_secret>

ETSY_ORDERS_ACCESS_KEY=<access_key>

ETSY_ORDERS_ACCESS_SECRET=<access_secret>

ETSY_ORDERS_SHOP_ID=<shop ID>

### Auth Access Token

The easiest way to to do this using the etsy gem and just save the perm token for deployment.

These are what you add to the ETSY_ORDERS_ACCESS_KEY and ETSY_ORDERS_ACCESS_SECRET environment variables.

https://github.com/kytrinyx/etsy

## Reply

```json
{
  "orders": [
    {
      "name": "Joy-con Basic Grip Plus",
      "primary_color": "6 Pink",
      "secondary_color": "",
      "slot_amount": "",
      "quantity": 1,
      "message_from_seller": "",
      "gift_message": "",
      "url": "https://www.etsy.com/transaction/123",
      "days_from_due_date": 7,
      "image_url": "https://i.etsystatic.com/image.jpg"
    }
  ]
}
```
