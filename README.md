# Checkout

This end point purpose to calculate final product checkout, with list SKU and Quantity as input

## Request Mutation

```graphql
mutation {
  checkout(reqs: [ { sku: "120P90", qty: 1} ]){
    products{
      sku
      name
      price
      qty
      discount
      subtotal // final price after discount
    }
    total
  }
}
```

## Sample cURL

```curl

// sample curl

curl --location --request POST 'localhost:8080/graphql' \
--header 'Content-Type: application/json' \
--data-raw '{"query":"mutation {\r\n    checkout(reqs: [{ sku: \"120P90\" qty: 3}{ sku: \"A304SD\" qty: 4}{ sku: \"43N23P\" qty: 1}{ sku: \"234234\" qty: 1}]){\r\n        products{\r\n            sku\r\n            name\r\n            price\r\n            qty\r\n            discount\r\n            subtotal\r\n        }\r\n        total\r\n    }\r\n}","variables":{}}'

```
