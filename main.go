package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/graphql-go/graphql"
)

// CheckoutRequest checkout request payload struct
type CheckoutRequest struct {
	SKU string `json:"sku"`
	Qty int    `json:"qty"`
}

// CheckoutResponse checkout response
type CheckoutResponse struct {
	Products []Product `json:"products"`
	Total    float64   `json:"total"`
}

// Product struct store product information
type Product struct {
	SKU      string  `json:"sku"`
	Name     string  `json:"name"`
	Price    float64 `json:"price"`
	Qty      int     `json:"qty"`
	Discount float64 `json:"discount"`
	Subtotal float64 `json:"subtotal"`
}

// Promo struct store Promo information
type Promo struct {
	PromoID            string  `json:"promo_id,omitempty"`
	SKU                string  `json:"sku,omitempty"`
	MinPurchase        int     `json:"min_purchase,omitempty"`
	FreeItemSKU        string  `json:"free_item_sku,omitempty"`
	FreeItemQty        int     `json:"free_item_qty,omitempty"`
	DiscountPercentage float64 `json:"discount_percentage,omitempty"`
	PromoQty           int     `json:"promo_qty,omitempty"`
}

var (
	ProductList map[string]Product
	PromoList   map[string]Promo
	Schema      graphql.Schema
)

func init() {
	ProductList = initProductList()
	PromoList = initPromoList()
}

func main() {
	schema, err := appSchema()
	if err != nil {
		panic(err)
	}
	Schema = schema

	http.HandleFunc("/graphql", AppHandler)

	fmt.Println("Running on port 8080")
	http.ListenAndServe(":8080", nil)
}

// initProductList init product list for API
func initProductList() map[string]Product {
	products := make(map[string]Product)

	file, err := os.Open("data/products.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var list []Product
	err = json.Unmarshal(jsonBytes, &list)
	if err != nil {
		log.Println("Error Unmarshal")
		panic(err)
	}

	for _, product := range list {
		products[product.SKU] = product
	}

	fmt.Println("Success init products")
	return products
}

// initPromoList to init promo list data
func initPromoList() map[string]Promo {
	promos := make(map[string]Promo)

	file, err := os.Open("data/promo.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	jsonBytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	var list []Promo
	err = json.Unmarshal(jsonBytes, &list)
	if err != nil {
		log.Println("Error Unmarshal")
		panic(err)
	}

	for _, promo := range list {
		promos[promo.SKU] = promo
	}

	fmt.Println("Success init promos")
	return promos
}

// initiate app schema for the API
func appSchema() (graphql.Schema, error) {
	checkoutType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name: "checkoutType",
		Fields: graphql.InputObjectConfigFieldMap{
			"sku": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"qty": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.Int),
			},
		},
	})

	productType := graphql.NewObject(graphql.ObjectConfig{
		Name: "productType",
		Fields: graphql.Fields{
			"sku": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"qty": &graphql.Field{
				Type: graphql.Int,
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
			"discount": &graphql.Field{
				Type: graphql.Float,
			},
			"subtotal": &graphql.Field{
				Type: graphql.Float,
			},
		},
	})

	checkoutResponse := graphql.NewObject(graphql.ObjectConfig{
		Name: "checkoutResponse",
		Fields: graphql.Fields{
			"products": &graphql.Field{
				Type: graphql.NewList(productType),
			},
			"total": &graphql.Field{
				Type: graphql.Float,
			},
		},
	})

	mutation := graphql.NewObject(graphql.ObjectConfig{
		Name: "RootMutation",
		Fields: graphql.Fields{
			"checkout": &graphql.Field{
				Type:        checkoutResponse,
				Description: "Mutation to proceed checkout request",
				Args: graphql.FieldConfigArgument{
					"reqs": &graphql.ArgumentConfig{
						Type: graphql.NewList(checkoutType),
					},
				},
				Resolve: RootCheckoutResolver,
			},
		},
	})

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    &graphql.Object{},
		Mutation: mutation,
	})
}
