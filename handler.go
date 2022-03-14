package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
)

type postData struct {
	Query     string                 `json:"query"`
	Operation string                 `json:"operation"`
	Variables map[string]interface{} `json:"variables"`
}

var (
	errInvalidCasting = errors.New("invalid type cast")
)

// AppHandler handler for this API
func AppHandler(w http.ResponseWriter, req *http.Request) {
	var p postData
	if err := json.NewDecoder(req.Body).Decode(&p); err != nil {
		w.WriteHeader(400)
		return
	}
	result := graphql.Do(graphql.Params{
		Context:        req.Context(),
		Schema:         Schema,
		RequestString:  p.Query,
		VariableValues: p.Variables,
		OperationName:  p.Operation,
	})
	if err := json.NewEncoder(w).Encode(result); err != nil {
		fmt.Printf("could not write result to response: %s", err)
	}
}

// RootCheckoutResolver root resolver function for GQL checkout mutation
func RootCheckoutResolver(params graphql.ResolveParams) (interface{}, error) {
	args, valid := params.Args["reqs"].([]interface{})
	if !valid {
		return nil, errInvalidCasting
	}
	var checkoutReq []CheckoutRequest
	for _, arg := range args {
		// casting to request type
		req, valid := arg.(map[string]interface{})
		if !valid {
			return nil, errInvalidCasting
		}

		sku := req["sku"].(string)
		qty := req["qty"].(int)

		checkoutReq = append(checkoutReq, CheckoutRequest{
			SKU: sku,
			Qty: qty,
		})
	}

	fmt.Println("reqs", checkoutReq)
	resp, err := Checkout(checkoutReq)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
