package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Paddle config (use environment variables in production)
var (
	paddleVendorID   = os.Getenv("PADDLE_VENDOR_ID")   // Your Paddle Vendor ID
	paddleApiKey     = os.Getenv("PADDLE_API_KEY")     // Paddle API Auth Key (secret)
	paddleBaseURL    = "https://api.paddle.com"        // Change to https://sandbox-api.paddle.com for testing
	paddleHttpClient = &http.Client{Timeout: 30 * time.Second}
)

// PaddleProductResponse represents Paddle product response
type PaddleProductResponse struct {
	Data struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		// ... more fields
	} `json:"data"`
	Error string `json:"error,omitempty"`
}

// PaddlePriceResponse
type PaddlePriceResponse struct {
	Data struct {
		ID        string `json:"id"`
		ProductID string `json:"product_id"`
		UnitPrice struct {
			Amount   string `json:"amount"`
			Currency string `json:"currency_code"`
		} `json:"unit_price"`
		Recurring bool   `json:"recurring"`
		Interval  string `json:"interval"` // "month", "year"
		// ...
	} `json:"data"`
}

// PaddleCustomerResponse
type PaddleCustomerResponse struct {
	Data struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
		// ...
	} `json:"data"`
}

// PaddleSubscriptionResponse
type PaddleSubscriptionResponse struct {
	Data struct {
		ID            string `json:"id"`
		Status        string `json:"status"`
		CustomerID    string `json:"customer_id"`
		PriceID       string `json:"price_id"`
		NextBillingAt string `json:"next_billed_at"`
		// ...
	} `json:"data"`
}

// Helper: make Paddle API request
func paddleRequest(method, path string, body interface{}, result interface{}) error {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return err
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	req, err := http.NewRequest(method, paddleBaseURL+path, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+paddleApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := paddleHttpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		return fmt.Errorf("paddle error %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		return json.Unmarshal(bodyBytes, result)
	}
	return nil
}

// -----------------------------------------------------------------------------
// 1. SyncOrCreateProduct (Paddle version)
// -----------------------------------------------------------------------------
func (app *application) SyncOrCreateProductPaddle(params map[string]any) Dict {
	data := Dict{}
	if d, ok := params["data"].(Dict); ok {
		if dd, ok := d["data"].(Dict); ok {
			data = dd
		}
	}

	internalID := data["plan_id"].(string)
	name := data["plan"].(string)

	// 1. Check if product already exists in Paddle
	var existingProduct PaddleProductResponse
	err := paddleRequest("GET", "/products?external_id="+internalID, nil, &existingProduct)
	if err == nil && existingProduct.Data.ID != "" {
		log.Printf("Using existing Paddle product: %s", existingProduct.Data.ID)

		// Optional: update name if changed
		updateBody := map[string]interface{}{
			"name": name,
		}
		var updated PaddleProductResponse
		paddleRequest("PATCH", "/products/"+existingProduct.Data.ID, updateBody, &updated)

		productID := updated.Data.ID

		// Now handle prices (Paddle prices are separate)
		return app.syncPaddlePrices(internalID, productID, data)
	}

	// 2. Create new product in Paddle
	createBody := map[string]interface{}{
		"name":        name,
		"description": data["description"],
		"external_id": internalID, // important: link to your internal ID
	}

	var created PaddleProductResponse
	err = paddleRequest("POST", "/products", createBody, &created)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle product: %v", err)}
	}

	productID := created.Data.ID

	// Save product ID to your DB
	app.savePaymentProductID(internalID, productID, "paddle", created)

	// 3. Create prices
	return app.syncPaddlePrices(internalID, productID, data)
}

// Helper: sync/create Paddle prices
func (app *application) syncPaddlePrices(internalID, productID string, data Dict) Dict {
	monthlyAmount := data["monthly_amount"].(float64)
	annualAmount := data["annual_amount"].(float64)
	currency := data["currency"].(string)
	if currency == "" {
		currency = "USD"
	}

	prices := []Dict{}

	// Monthly price
	if monthlyAmount > 0 {
		priceBody := map[string]interface{}{
			"product_id":   productID,
			"unit_price": map[string]interface{}{
				"amount":        fmt.Sprintf("%.2f", monthlyAmount),
				"currency_code": currency,
			},
			"recurring": true,
			"interval":  "month",
			"external_id": internalID + "-monthly",
		}

		var monthlyResp PaddlePriceResponse
		err := paddleRequest("POST", "/prices", priceBody, &monthlyResp)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to create monthly price: %v", err)}
		}
		prices = append(prices, Dict{
			"type": "monthly",
			"id":   monthlyResp.Data.ID,
		})
	}

	// Annual price
	if annualAmount > 0 {
		priceBody := map[string]interface{}{
			"product_id":   productID,
			"unit_price": map[string]interface{}{
				"amount":        fmt.Sprintf("%.2f", annualAmount),
				"currency_code": currency,
			},
			"recurring": true,
			"interval":  "year",
			"external_id": internalID + "-annual",
		}

		var annualResp PaddlePriceResponse
		err := paddleRequest("POST", "/prices", priceBody, &annualResp)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to create annual price: %v", err)}
		}
		prices = append(prices, Dict{
			"type": "annual",
			"id":   annualResp.Data.ID,
		})
	}

	// Save prices to your DB
	app.savePaymentPriceIDs(internalID, prices, "paddle")

	return Dict{
		"success": true,
		"msg":     "Paddle product and prices synced",
		"product_id": productID,
		"prices":     prices,
	}
}

// -----------------------------------------------------------------------------
// 2. CreateOrSyncCustomer (Paddle version)
// -----------------------------------------------------------------------------
func (app *application) CreateOrSyncCustomerPaddle(params map[string]any) Dict {
	data := extractData(params)

	tenantID := data["tenant_id"].(string)
	name := data["name"].(string)
	email := data["email"].(string)
	paymentProviderCustomerID := data["pay_provider_customer_id"].(string)

	var customer PaddleCustomerResponse

	// Try to find existing customer
	if paymentProviderCustomerID != "" {
		err := paddleRequest("GET", "/customers/"+paymentProviderCustomerID, nil, &customer)
		if err == nil {
			log.Printf("Using existing Paddle customer: %s", customer.Data.ID)
			// Optional: update if changed
			if customer.Data.Email != email || customer.Data.Name != name {
				updateBody := map[string]interface{}{
					"name":  name,
					"email": email,
				}
				paddleRequest("PATCH", "/customers/"+customer.Data.ID, updateBody, &customer)
			}
			return Dict{"success": true, "customer": customer.Data}
		}
	}

	// Create new customer
	createBody := map[string]interface{}{
		"name":  name,
		"email": email,
		"custom": map[string]string{
			"tenant_id": tenantID,
		},
	}

	err := paddleRequest("POST", "/customers", createBody, &customer)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle customer: %v", err)}
	}

	// Save to DB
	app.savePaymentCustomerID(tenantID, customer.Data.ID, "paddle", customer)

	return Dict{"success": true, "msg": "Paddle customer created", "customer": customer.Data}
}

// -----------------------------------------------------------------------------
// 3. CreateOrUpdateSubscription (Paddle version)
// -----------------------------------------------------------------------------
func (app *application) CreateOrUpdateSubscriptionPaddle(params map[string]any) Dict {
	data := extractData(params)

	tenantID := data["tenant_id"].(string)
	paymentPlanID := data["payment_plan_id"].(string) // your internal plan/price ID

	// Get Paddle customer ID from tenant
	customerID := app.getPaddleCustomerIDForTenant(tenantID)
	if customerID == "" {
		return Dict{"success": false, "msg": "No Paddle customer ID found for tenant"}
	}

	// Get Paddle price ID from your payment plan
	priceID := app.getPaddlePriceIDForPlan(paymentPlanID)
	if priceID == "" {
		return Dict{"success": false, "msg": "No Paddle price ID found for plan"}
	}

	// Check if subscription already exists
	existingSubID := data["pay_provider_subscription_id"].(string)
	var subscription PaddleSubscriptionResponse

	if existingSubID != "" {
		err := paddleRequest("GET", "/subscriptions/"+existingSubID, nil, &subscription)
		if err == nil {
			log.Printf("Using existing Paddle subscription: %s", subscription.Data.ID)

			// Optional: update price if changed
			if subscription.Data.PriceID != priceID {
				updateBody := map[string]interface{}{
					"items": []map[string]string{
						{"price_id": priceID},
					},
				}
				paddleRequest("PATCH", "/subscriptions/"+existingSubID, updateBody, &subscription)
			}

			app.savePaymentSubscriptionData(tenantID, subscription.Data.ID, subscription, "paddle")
			return Dict{"success": true, "msg": "Subscription updated", "subscription": subscription.Data}
		}
	}

	// Create new subscription
	createBody := map[string]interface{}{
		"customer_id": customerID,
		"items": []map[string]interface{}{
			{
				"price_id": priceID,
				"quantity": 1,
			},
		},
		// Optional: trial, currency, etc.
	}

	var newSub PaddleSubscriptionResponse
	err := paddleRequest("POST", "/subscriptions", createBody, &newSub)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle subscription: %v", err)}
	}

	// Save to your DB
	app.savePaymentSubscriptionData(tenantID, newSub.Data.ID, newSub, "paddle")

	return Dict{
		"success":      true,
		"msg":          "Paddle subscription created",
		"subscription": newSub.Data,
	}
}

// Helper functions you need to implement
func (app *application) savePaymentProductID(internalID, paddleID, provider string, resp interface{}) {
	// Save to your DB: update plan table with pay_provider_product_id, payment_provider, metadata
}

func (app *application) savePaymentPriceIDs(internalID string, prices []Dict, provider string) {
	// Update your price table with monthly/annual pay_provider_price_id
}

func (app *application) savePaymentCustomerID(tenantID, paddleCustomerID, provider string, resp interface{}) {
	// Update tenant table
}

func (app *application) savePaymentSubscriptionData(tenantID, subID string, resp interface{}, provider string) {
	// Update subscription table with pay_provider_subscription_id, metadata, status, next_billing_at, etc.
}

func (app *application) getPaddleCustomerIDForTenant(tenantID string) string {
	// Query your DB
	return ""
}

func (app *application) getPaddlePriceIDForPlan(planID string) string {
	// Query your DB
	return ""
}

func extractData(params map[string]any) Dict {
	// Your existing logic to extract data from params
	return Dict{}
}