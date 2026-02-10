package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	paddle "github.com/PaddleHQ/paddle-go-sdk/v4"
	//"github.com/PaddleHQ/paddle-go-sdk/v2"
	//"github.com/PaddleHQ/paddle-go-sdk/v2/client"
	//"github.com/PaddleHQ/paddle-go-sdk/v2/models"
)

func (app *application) initPaddle() (*paddle.SDK, error) {
	paddleEnv := paddle.SandboxBaseURL
	if os.Getenv("PADDLE_ENV") == "production" {
		paddleEnv = paddle.ProductionBaseURL
	}
	cliente, err := paddle.New(
		os.Getenv("PADDLE_API_KEY"),
		paddle.WithBaseURL(paddleEnv), // or paddle.ProductionBaseURL for accessing live API
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Paddle client: %w", err)
	}
	return cliente, nil
}

// Helper to convert any struct to Dict (map[string]any)
func toDict(v any) Dict {
	b, _ := json.Marshal(v)
	var m Dict
	json.Unmarshal(b, &m)
	return m
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

	monthlyAmount := data["monthly_amount"].(float64)
	annualAmount := data["annual_amount"].(float64)
	currency := data["currency"].(string)
	if currency == "" {
		currency = "USD"
	}

	existingProductID := data["payment_provider_product_id"].(string)

	ctx := context.Background()
	var product *paddle.Product
	var err error
	client, err := app.initPaddle()
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("Failed to initialize Paddle client: %v", err)}
	}
	// 1. Try to get existing product
	if existingProductID != "" {
		product, err = client.GetProduct(ctx, &paddle.GetProductRequest{ProductID: existingProductID})
		if err == nil {
			log.Printf("Using existing Paddle product: %s", product.ID)
		} else {
			log.Printf("Paddle product not found: %v - will create new", err)
		}
	}

	// 2. Create or update product
	if product == nil {
		description, _ := data["description"].(string)
		createReq := &paddle.CreateProductRequest{
			Name:        name,
			Description: &description,
			CustomData: map[string]any{
				"internal_plan_id": planID,
			},
		}

		created, err := client.CreateProduct(ctx, createReq)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle product: %v", err)}
		}
		product = created
	} else {
		// Optional: update if name changed
		if product.Name != name {
			updateReq := &paddle.UpdateProductRequest{
				ProductID: product.ID,
				Name:      paddle.NewPatchField(name),
			}
			product, err = client.UpdateProduct(ctx, updateReq)
			if err != nil {
				return Dict{"success": false, "msg": fmt.Sprintf("Failed to update Paddle product: %v", err)}
			}
		}
	}

	productID := product.ID

	// 3. Handle prices
	var monthlyPriceID, annualPriceID string

	createOrUpdatePrice := func(interval string, amount float64, existingID string) (string, error) {
		if amount <= 0 {
			return "", nil
		}

		unitPrice := fmt.Sprintf("%.0f", amount*100) // Paddle uses minor units as string

		var billingInterval *paddle.Duration
		if interval == "month" {
			billingInterval = &paddle.Duration{Interval: paddle.IntervalMonth, Frequency: 1}
		} else if interval == "year" {
			billingInterval = &paddle.Duration{Interval: paddle.IntervalYear, Frequency: 1}
		}

		var price *paddle.Price
		var err error

		if existingID != "" {
			priceReq := &paddle.UpdatePriceRequest{
				PriceID: existingID,
				UnitPrice: paddle.NewPatchField(paddle.Money{
					Amount:       unitPrice,
					CurrencyCode: paddle.CurrencyCode(currency),
				}),
				BillingCycle: paddle.NewPatchField(billingInterval),
				CustomData: paddle.NewPatchField(paddle.CustomData{
					"interval": interval,
				}),
			}
			// Paddle doesn't have direct update for price — you usually create new versions
			// For simplicity, we'll create a new one if needed
			price, err = client.UpdatePrice(ctx, priceReq)
		} else {
			priceReq := &paddle.CreatePriceRequest{
				ProductID: productID,
				UnitPrice: paddle.Money{
					Amount:       unitPrice,
					CurrencyCode: paddle.CurrencyCode(currency),
				},
				BillingCycle: billingInterval,
				CustomData: map[string]any{
					"interval": interval,
				},
			}
			price, err = client.CreatePrice(ctx, priceReq)
		}

		if err != nil {
			return "", err
		}

		return price.ID, nil
	}

	monthlyPriceID, _ = createOrUpdatePrice("month", monthlyAmount, data["payment_price_monthly_id"].(string))
	annualPriceID, _ = createOrUpdatePrice("year", annualAmount, data["payment_price_annual_id"].(string))

	// 4. Save back to your database
	updateData := data
	updateData["payment_provider_product_id"] = productID
	updateData["payment_product_metadata"] = toDict(product)
	updateData["payment_price_monthly_id"] = monthlyPriceID
	updateData["payment_price_annual_id"] = annualPriceID
	//updateData["payment_provider_price_metadata"] = toDict(Dict{"monthly": monthlyPriceID, "annual": annualPriceID})
	updateData["payment_provider_last_sync_at"] = time.Now()

	params["data"] = Dict{
		"table":   "plan",
		"data":    updateData,
		"filters": []Dict{{"field": "plan_id", "value": planID}},
	}

	updateResult := app.create_update(params)
	if !updateResult["success"].(bool) {
		return updateResult
	}

	return Dict{
		"success":          true,
		"msg":              "Paddle product and prices synced",
		"product_id":       productID,
		"monthly_price_id": monthlyPriceID,
		"annual_price_id":  annualPriceID,
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

	existingCustomerID := data["payment_provider_customer_id"].(string)

	ctx := context.Background()
	var customer *paddle.Customer
	var err error

	client, err := app.initPaddle()
	if existingCustomerID != "" {
		getCustomerReq := &paddle.GetCustomerRequest{CustomerID: existingCustomerID}
		customer, err = client.GetCustomer(ctx, getCustomerReq)
		if err == nil {
			log.Printf("Using existing Paddle customer: %s", customer.ID)
		}
	}

	if customer == nil {
		createReq := &paddle.CreateCustomerRequest{
			Name:  &name,
			Email: email,
			CustomData: map[string]any{
				"tenant_id": tenantID,
			},
		}

		created, err := client.CreateCustomer(ctx, createReq)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle customer: %v", err)}
		}
		customer = created
	}

	customerID := customer.ID

	// Save back to DB
	updateData := data
	updateData["payment_customer_id"] = customerID
	updateData["payment_customer_metadata"] = toDict(customer)
	updateData["payment_customer_last_sync_at"] = time.Now()

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
		"success":     true,
		"msg":         "Paddle customer synced",
		"customer_id": customerID,
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

/*func (app *application) Merge[M ~map[K]V, K comparable, V any](dst M, srcs ...M) {
	for _, src := range srcs {
		for k, v := range src {
			dst[k] = v
		}
	}
}*/

/*/ ====================================================================
// Paddle: Create or Update Subscription
// ====================================================================
func (app *application) PaddleCreateOrUpdateSubscription(params map[string]any) Dict {
	data := extractData(params)
	if data == nil {
		return Dict{"success": false, "msg": "invalid params"}
	}
	tenantID := data["tenant_id"].(string)
	planID := data["price_id"].(string)
	// Get tenant's Paddle customer ID
	params["data"] = Dict{
		"table":   "tenant",
		"filters": []Dict{{"field": "tenant_id", "value": tenantID}},
	}
	tenantResp := app.read(params)
	if !tenantResp["success"].(bool) {
		return tenantResp
	}
	tenantData := tenantResp["data"].([]Dict)[0]
	customerID := tenantData["payment_provider_customer_id"].(string)

	// Get price ID from plan/prices
	params["data"] = Dict{
		"table":   "plan",
		"filters": []Dict{{"field": "plan_id", "value": planID}},
	}
	planResp := app.read(params)
	if !planResp["success"].(bool) {
		return planResp
	}
	planData := planResp["data"].([]Dict)[0]
	priceID := planData["payment_provider_price_id"].(string)
	if priceID == "" {
		return Dict{"success": false, "msg": "No Paddle price ID found for plan"}
	}
	existingSubID := data["payment_provider_subscription_id"].(string)
	ctx := context.Background()
	var sub *paddle.Subscription
	var err error
	client, err := app.initPaddle()
	if existingSubID != "" {
		getSubReq := &paddle.GetSubscriptionRequest{SubscriptionID: existingSubID}
		sub, err = client.GetSubscription(ctx, getSubReq)
		if err == nil {
			log.Printf("Using existing Paddle subscription: %s", sub.ID)
		}
	}
	// client
	if sub == nil {
		createReq := &paddle.CreateSubscriptionRequest{
			CustomerID: customerID,
			Items: []paddle.SubscriptionItemCreateRequest{
				{
					PriceID:  priceID,
					Quantity: 1,
				},
			},
			CustomData: map[string]any{
				"tenant_id": tenantID,
				"plan_id":   planID,
			},
		}

		created, err := client.Create(ctx, createReq)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle subscription: %v", err)}
		}
		sub = created
	} else {
		// Optional: update subscription if needed (e.g. change price)
		updateReq := &paddle.UpdateSubscriptionRequest{
			SubscriptionID: sub.ID,
			Items: []paddle.SubscriptionItemUpdateRequest{
				{
					PriceID:  priceID,
					Quantity: 1,
				},
			},
		}
		sub, err = client.Update(ctx, updateReq)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to update Paddle subscription: %v", err)}
		}
	}

	subID := sub.ID

	// Save back to DB
	updateData := data
	updateData["payment_provider_subscription_id"] = subID
	updateData["payment_provider_subscription_metadata"] = toDict(sub)
	updateData["payment_provider_last_sync_at"] = time.Now()
	updateData["status"] = sub.Status
	updateData["next_billing_at"] = sub.NextBilledAt

	params["data"] = Dict{
		"table": "subscription",
		"data":  updateData,
		//"filters": []Dict{{"field": "tenant_id", "value": tenantID}},
	}

	updateResult := app.create_update(params)
	if !updateResult["success"].(bool) {
		return updateResult
	}

	return Dict{
		"success":         true,
		"msg":             "Paddle subscription created/updated",
		"subscription_id": subID,
	}
}
*/
/*/ PaddleSyncAllSubscriptionStatuses – batch job example
func (app *application) PaddleSyncAllSubscriptionStatuses() Dict {
	// Get all active/known subscriptions from your DB
	params := map[string]any{
		"data": Dict{
			"table": "subscription",
			"filters": []Dict{
				{"field": "payment_provider", "value": "paddle"},
				{"field": "excluded", "value": false},
			},
			"limit": 100, // paginate if many tenants
		},
	}

	subsResp := app.read(params)
	if !subsResp["success"].(bool) {
		return subsResp
	}

	subscriptions := subsResp["data"].([]Dict)
	if len(subscriptions) == 0 {
		return Dict{"success": true, "msg": "No Paddle subscriptions to sync"}
	}

	ctx := context.Background()
	updatedCount := 0

	for _, sub := range subscriptions {
		subID := sub["payment_provider_subscription_id"].(string)
		if subID == "" {
			continue
		}

		paddleSub, err := client.Subscriptions.Get(ctx, subID)
		if err != nil {
			log.Printf("Failed to get Paddle sub %s: %v", subID, err)
			continue
		}

		newStatus := string(paddleSub.Status) // e.g. "active", "past_due", ...
		nextBilling := paddleSub.NextBilledAt // *time.Time

		updateParams := map[string]any{
			"data": Dict{
				"table":   "subscription",
				"filters": []Dict{{"field": "payment_provider_subscription_id", "value": subID}},
				"data": Dict{
					"status":                             newStatus,
					"next_billing_at":                    nextBilling,
					"payment_provider_last_sync_at":          time.Now(),
					"payment_provider_subscription_metadata": toDict(paddleSub),
				},
			},
		}

		result := app.create_update(updateParams)
		if result["success"].(bool) {
			updatedCount++
			log.Printf("Synced Paddle sub %s → status: %s", subID, newStatus)
		} else {
			log.Printf("Failed to update sub %s: %v", subID, result["msg"])
		}
	}

	return Dict{
		"success": true,
		"msg":     fmt.Sprintf("Synced status for %d subscriptions", updatedCount),
	}
}*/
