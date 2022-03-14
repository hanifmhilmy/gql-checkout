package main

import (
	"errors"
)

// Checkout checkout function
func Checkout(reqs []CheckoutRequest) (result CheckoutResponse, err error) {
	freeItem := make(map[string]int)
	result.Products = make([]Product, 0, len(reqs))

	for i := range reqs {
		req := reqs[i]

		product, err := getProductDetails(req.SKU)
		if err != nil {
			return result, err
		}

		if req.Qty > product.Qty {
			return result, errors.New("stock not enough")
		}

		if qty, ok := freeItem[req.SKU]; ok {
			if qty >= req.Qty {
				req.Qty = qty

				product.Price = 0
			} else {
				req.Qty = req.Qty - qty
			}

			delete(freeItem, req.SKU)
		}

		subtotal := product.Price * float64(req.Qty)
		// check promo
		promo := getPromo(req.SKU)
		if req.Qty >= promo.MinPurchase {
			if promo.FreeItemSKU != "" {
				if promo.FreeItemSKU == req.SKU {
					disc := float64(promo.FreeItemQty) * product.Price
					subtotal = subtotal - disc
					product.Discount = disc

					if promo.FreeItemQty >= req.Qty {
						req.Qty = promo.FreeItemQty

						product.Price = 0
						subtotal = 0
					}
				} else {
					freeItem[promo.FreeItemSKU] = promo.FreeItemQty
				}
			} else if promo.DiscountPercentage > 0 {
				disc := (subtotal * (promo.DiscountPercentage / 100))
				subtotal = subtotal - disc
				product.Discount = disc
			}
		}

		// productList
		product.Qty = req.Qty
		product.Subtotal = subtotal
		result.Products = append(result.Products, product)

		// calculate total
		result.Total = result.Total + subtotal
	}

	// check if there is leftover free item, append to list
	for key, qty := range freeItem {
		product, _ := getProductDetails(key)

		product.Price = 0
		product.Qty = qty

		result.Products = append(result.Products, product)
	}

	return result, nil
}

func getPromo(sku string) Promo {
	if data, ok := PromoList[sku]; ok {
		return data
	}
	return Promo{}
}

func getProductDetails(sku string) (Product, error) {
	if data, ok := ProductList[sku]; ok {
		return data, nil
	}
	return Product{}, errors.New("product not found")
}
