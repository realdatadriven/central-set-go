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

	"github.com/stripe/stripe-go/v84" // still needed if you want to keep both

	// Paddle doesn't have an official Go SDK (as of 2025), so we use HTTP
	// You can also use community SDKs like github.com/PaddleHQ/paddle-sdk-go (if available)
)

// PaddleConfig holds your Paddle credentials
type PaddleConfig struct {
	VendorID     string // or Team ID
	APIKey       string // Paddle API key (from Developer Tools → Authentication)
	WebhookSecret string // for webhook verification
	BaseURL      string // "https://api.paddle.com" (production) or "https://sandbox-api.paddle.com"
}

var paddleConfig PaddleConfig

func init() {
	paddleConfig = PaddleConfig{
		VendorID:     os.Getenv("PADDLE_VENDOR_ID"),
		APIKey:       os.Getenv("PADDLE_API_KEY"),
		WebhookSecret: os.Getenv("PADDLE_WEBHOOK_SECRET"),
		BaseURL:      "https://api.paddle.com", // change to sandbox for testing
	}
	if paddleConfig.APIKey == "" {
		log.Println("Warning: PADDLE_API_KEY not set")
	}
}

// paddleRequest is a helper to make authenticated requests to Paddle API
func paddleRequest(method, path string, body any) (map[string]any, error) {
	url := paddleConfig.BaseURL + path

	var bodyBytes []byte
	var err error
	if body != nil {
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+paddleConfig.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("paddle error %d: %s", resp.StatusCode, string(respBody))
	}

	var result map[string]any
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ====================================================================
// Paddle: Sync or Create Product
// ====================================================================
func (app *application) PaddleSyncOrCreateProduct(params map[string]any) Dict {
	data := extractData(params)
	if data == nil {
		return Dict{"success": false, "msg": "invalid params structure"}
	}

	planID, _ := data["plan_id"].(string)
	name, _ := data["plan"].(string)

	// Read existing prices from your DB
	params["data"] = Dict{"table": "price", "filters": []Dict{{"field": "plan_id", "value": planID}}}
	pricesResp := app.read(params)
	if !pricesResp["success"].(bool) {
		return pricesResp
	}

	priceDataList := pricesResp["data"].([]Dict)
	if len(priceDataList) == 0 {
		return Dict{"success": false, "msg": "No price data found for plan"}
	}
	priceData := priceDataList[0]

	monthlyAmount := priceData["monthly_amount"].(float64)
	annualAmount := priceData["annual_amount"].(float64)
	currency := priceData["currency"].(string)
	if currency == "" {
		currency = "USD"
	}

	existingProductID := priceData["pay_provider_product_id"].(string)

	var productID string
	var productData map[string]any

	// 1. Try to get existing product
	if existingProductID != "" {
		resp, err := paddleRequest("GET", "/products/"+existingProductID, nil)
		if err == nil {
			productData = resp["data"].(map[string]any)
			productID = productData["id"].(string)
			log.Printf("Using existing Paddle product: %s", productID)
		} else {
			log.Printf("Paddle product not found: %v - will create new", err)
		}
	}

	// 2. Create or update product
	if productID == "" {
		createBody := map[string]any{
			"name":        name,
			"description": data["description"],
			"custom_data": map[string]any{
				"internal_plan_id": planID,
			},
		}

		resp, err := paddleRequest("POST", "/products", createBody)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle product: %v", err)}
		}

		productData = resp["data"].(map[string]any)
		productID = productData["id"].(string)
	}

	// 3. Handle prices (Paddle calls them prices, similar to Stripe)
	var monthlyPriceID, annualPriceID string

	createOrUpdatePrice := func(interval string, amount float64, existingID string) (string, error) {
		if amount <= 0 {
			return "", nil
		}

		priceBody := map[string]any{
			"product_id":  productID,
			"unit_price": map[string]any{
				"amount": amount * 100, // Paddle uses minor units (cents)
				"currency_code": currency,
			},
			"billing_cycle": map[string]any{
				"interval": interval,
				"frequency": 1,
			},
			"custom_data": map[string]any{
				"interval": interval,
			},
		}

		var resp map[string]any
		var err error

		if existingID != "" {
			resp, err = paddleRequest("PATCH", "/prices/"+existingID, priceBody)
		} else {
			resp, err = paddleRequest("POST", "/prices", priceBody)
		}

		if err != nil {
			return "", err
		}

		return resp["data"].(map[string]any)["id"].(string), nil
	}

	monthlyPriceID, _ = createOrUpdatePrice("month", monthlyAmount, priceData["pay_provider_monthly_id"].(string))
	annualPriceID, _ = createOrUpdatePrice("year", annualAmount, priceData["pay_provider_annual_id"].(string))

	// 4. Save back to your database
	updateData := Dict{
		"pay_provider_product_id": productID,
		"pay_provider_product_metadata": toJSON(productData),
		"pay_provider_monthly_id": monthlyPriceID,
		"pay_provider_annual_id": annualPriceID,
		"pay_provider_price_metadata": toJSON(Dict{
			"monthly": monthlyPriceID,
			"annual":  annualPriceID,
		}),
		"pay_provider_last_sync_at": time.Now(),
	}

	params["data"] = Dict{
		"table": "plan",
		"data":  updateData,
		"filters": []Dict{{"field": "plan_id", "value": planID}},
	}

	updateResult := app.create_update(params)
	if !updateResult["success"].(bool) {
		return updateResult
	}

	return Dict{
		"success": true,
		"msg":     "Paddle product and prices synced",
		"product_id": productID,
		"monthly_price_id": monthlyPriceID,
		"annual_price_id": annualPriceID,
	}
}

// ====================================================================
// Paddle: Create or Sync Customer
// ====================================================================
func (app *application) PaddleCreateOrSyncCustomer(params map[string]any) Dict {
	data := extractData(params)
	if data == nil {
		return Dict{"success": false, "msg": "invalid params"}
	}

	tenantID, _ := data["tenant_id"].(string)
	name, _ := data["name"].(string)
	email, _ := data["email"].(string)

	existingCustomerID := data["pay_provider_customer_id"].(string)

	var customerID string
	var customerData map[string]any

	if existingCustomerID != "" {
		resp, err := paddleRequest("GET", "/customers/"+existingCustomerID, nil)
		if err == nil {
			customerData = resp["data"].(map[string]any)
			customerID = customerData["customer_id"].(string)
			log.Printf("Using existing Paddle customer: %s", customerID)
		}
	}

	if customerID == "" {
		createBody := map[string]any{
			"name":  name,
			"email": email,
			"custom_data": map[string]any{
				"tenant_id": tenantID,
			},
		}

		resp, err := paddleRequest("POST", "/customers", createBody)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle customer: %v", err)}
		}

		customerData = resp["data"].(map[string]any)
		customerID = customerData["customer_id"].(string)
	}

	// Save back to DB
	updateData := Dict{
		"pay_provider_customer_id": customerID,
		"pay_provider_customer_metadata": toJSON(customerData),
		"pay_provider_last_sync_at": time.Now(),
	}

	params["data"] = Dict{
		"table":   "tenant",
		"data":    updateData,
		"filters": []Dict{{"field": "tenant_id", "value": tenantID}},
	}

	updateResult := app.create_update(params)
	if !updateResult["success"].(bool) {
		return updateResult
	}

	return Dict{
		"success": true,
		"msg":     "Paddle customer synced",
		"customer_id": customerID,
	}
}

// ====================================================================
// Paddle: Create or Update Subscription
// ====================================================================
func (app *application) PaddleCreateOrUpdateSubscription(params map[string]any) Dict {
	data := extractData(params)
	if data == nil {
		return Dict{"success": false, "msg": "invalid params"}
	}

	tenantID := data["tenant_id"].(string)
	priceID := data["price_id"].(string)

	// Get tenant's Paddle customer ID
	tenantResp := app.read(Dict{
		"table":   "tenant",
		"filters": []Dict{{"field": "tenant_id", "value": tenantID}},
	})
	if !tenantResp["success"].(bool) {
		return tenantResp
	}
	tenantData := tenantResp["data"].([]Dict)[0]
	customerID := tenantData["pay_provider_customer_id"].(string)

	// Get price IDs from plan/prices
	priceResp := app.read(Dict{
		"table":   "price",
		"filters": []Dict{{"field": "price_id", "value": priceID}},
	})
	if !priceResp["success"].(bool) {
		return priceResp
	}
	priceData := priceResp["data"].([]Dict)[0]

	priceID := priceData["pay_provider_price_id"].(string)
	if priceID == "" {
		return Dict{"success": false, "msg": "No Paddle price ID found for plan"}
	}

	existingSubID := data["pay_provider_subscription_id"].(string)

	var subID string
	var subData map[string]any

	if existingSubID != "" {
		resp, err := paddleRequest("GET", "/subscriptions/"+existingSubID, nil)
		if err == nil {
			subData = resp["data"].(map[string]any)
			subID = subData["subscription_id"].(string)
			log.Printf("Using existing Paddle subscription: %s", subID)
		}
	}

	if subID == "" {
		createBody := map[string]any{
			"customer_id": customerID,
			"items": []map[string]any{
				{
					"price_id": priceID,
					"quantity": 1,
				},
			},
			"custom_data": map[string]any{
				"tenant_id": tenantID,
				"plan_id":   planID,
			},
		}

		resp, err := paddleRequest("POST", "/subscriptions", createBody)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle subscription: %v", err)}
		}

		subData = resp["data"].(map[string]any)
		subID = subData["subscription_id"].(string)
	}

	// Save back to DB
	updateData := Dict{
		"pay_provider_subscription_id": subID,
		"pay_provider_subscription_metadata": toJSON(subData),
		"pay_provider_last_sync_at": time.Now(),
		"status": subData["status"], // e.g. active, past_due, canceled
		"next_billing_at": subData["next_billed_at"],
	}

	params["data"] = Dict{
		"table":   "subscription",
		"data":    updateData,
		"filters": []Dict{{"field": "tenant_id", "value": tenantID}},
	}

	updateResult := app.create_update(params)
	if !updateResult["success"].(bool) {
		return updateResult
	}

	return Dict{
		"success": true,
		"msg":     "Paddle subscription created/updated",
		"subscription_id": subID,
	}
}

// Helper functions
func extractData(params map[string]any) Dict {
	if d, ok := params["data"].(Dict); ok {
		if inner, ok := d["data"].(Dict); ok {
			return inner
		}
	}
	return nil
}

func toJSON(v any) string {
	b, _ := json.Marshal(v)
	return string(b)
}