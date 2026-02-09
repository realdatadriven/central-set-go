package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	
	 paddle "github.com/PaddleHQ/paddle-go-sdk/v3"
	//"github.com/PaddleHQ/paddle-go-sdk/v2"
	//"github.com/PaddleHQ/paddle-go-sdk/v2/client"
	//"github.com/PaddleHQ/paddle-go-sdk/v2/models"
)

// PaddleConfig holds your Paddle credentials
type PaddleConfig struct {
	APIKey   string
	Env      client.Environment // client.Sandbox or client.Production
}

var paddleClient *client.Client

func (app *application) initPaddle() {
	apiKey := os.Getenv("PADDLE_API_KEY")
	if apiKey == "" {
		log.Fatal("PADDLE_API_KEY not set")
	}

	/*paddleClient = client.NewClient(
		apiKey,
		client.WithEnvironment(client.Production), // or client.Sandbox
		// client.WithHTTPClient(&http.Client{Timeout: 30 * time.Second}),
	)*/
	paddleEnv := paddle.SandboxBaseURL
	if os.Getenv("PADDLE_ENV") == "production" {
		paddleEnv := paddle.ProductionBaseURL
	}
	if paddleClient == nil {
		paddleClient, err := paddle.New(
			os.Getenv("PADDLE_API_KEY"),
			paddle.WithBaseURL(paddleEnv) // or paddle.ProductionBaseURL for accessing live API
		)
		if err != nil {
			log.Fatalf("Failed to initialize Paddle client: %v", err)
		}
	}
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

	existingProductID := data["pay_provider_product_id"].(string)

	ctx := context.Background()
	var product *models.Product
	var err error

	// 1. Try to get existing product
	if existingProductID != "" {
		product, err = paddleClient.Products.Get(ctx, existingProductID)
		if err == nil {
			log.Printf("Using existing Paddle product: %s", product.Id)
		} else {
			log.Printf("Paddle product not found: %v - will create new", err)
		}
	}

	// 2. Create or update product
	if product == nil {
		createReq := &models.CreateProductRequest{
			Name:        name,
			Description: data["description"].(string),
			CustomData: map[string]any{
				"internal_plan_id": planID,
			},
		}

		created, err := paddleClient.Products.Create(ctx, createReq)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle product: %v", err)}
		}
		product = created
	} else {
		// Optional: update if name changed
		if product.Name != name {
			updateReq := &models.UpdateProductRequest{
				Name: &name,
			}
			product, err = paddleClient.Products.Update(ctx, product.Id, updateReq)
			if err != nil {
				return Dict{"success": false, "msg": fmt.Sprintf("Failed to update Paddle product: %v", err)}
			}
		}
	}

	productID := product.Id

	// 3. Handle prices
	var monthlyPriceID, annualPriceID string

	createOrUpdatePrice := func(interval string, amount float64, existingID string) (string, error) {
		if amount <= 0 {
			return "", nil
		}

		unitPrice := float64(int64(amount * 100)) // Paddle uses minor units

		priceReq := &models.CreatePriceRequest{
			ProductId: productID,
			UnitPrice: models.Money{
				Amount:       unitPrice,
				CurrencyCode: currency,
			},
			BillingCycle: models.BillingCycle{
				Interval:  interval,
				Frequency: 1,
			},
			CustomData: map[string]any{
				"interval": interval,
			},
		}

		var price *models.Price
		var err error

		if existingID != "" {
			// Paddle doesn't have direct update for price — you usually create new versions
			// For simplicity, we'll create a new one if needed
			price, err = paddleClient.Prices.Create(ctx, priceReq)
		} else {
			price, err = paddleClient.Prices.Create(ctx, priceReq)
		}

		if err != nil {
			return "", err
		}

		return price.Id, nil
	}

	monthlyPriceID, _ = createOrUpdatePrice("month", monthlyAmount, priceData["pay_provider_monthly_id"].(string))
	annualPriceID, _ = createOrUpdatePrice("year", annualAmount, priceData["pay_provider_annual_id"].(string))

	// 4. Save back to your database
	updateData := Dict{
		"pay_provider_product_id":        productID,
		"pay_provider_product_metadata":  toDict(product),
		"pay_provider_monthly_id":        monthlyPriceID,
		"pay_provider_annual_id":         annualPriceID,
		"pay_provider_price_metadata":    toDict(Dict{"monthly": monthlyPriceID, "annual": annualPriceID}),
		"pay_provider_last_sync_at":      time.Now(),
	}

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

	existingCustomerID := data["pay_provider_customer_id"].(string)

	ctx := context.Background()
	var customer *models.Customer
	var err error

	if existingCustomerID != "" {
		customer, err = paddleClient.Customers.Get(ctx, existingCustomerID)
		if err == nil {
			log.Printf("Using existing Paddle customer: %s", customer.CustomerId)
		}
	}

	if customer == nil {
		createReq := &models.CreateCustomerRequest{
			Name:  name,
			Email: email,
			CustomData: map[string]any{
				"tenant_id": tenantID,
			},
		}

		created, err := paddleClient.Customers.Create(ctx, createReq)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle customer: %v", err)}
		}
		customer = created
	}

	customerID := customer.CustomerId

	// Save back to DB
	updateData := Dict{
		"pay_provider_customer_id":       customerID,
		"pay_provider_customer_metadata": toDict(customer),
		"pay_provider_last_sync_at":      time.Now(),
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
		"success":     true,
		"msg":         "Paddle customer synced",
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
	planID := data["plan_id"].(string)

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

	// Get price ID from plan/prices
	priceResp := app.read(Dict{
		"table":   "price",
		"filters": []Dict{{"field": "plan_id", "value": planID}},
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

	ctx := context.Background()
	var sub *models.Subscription
	var err error

	if existingSubID != "" {
		sub, err = paddleClient.Subscriptions.Get(ctx, existingSubID)
		if err == nil {
			log.Printf("Using existing Paddle subscription: %s", sub.SubscriptionId)
		}
	}

	if sub == nil {
		createReq := &models.CreateSubscriptionRequest{
			CustomerId: customerID,
			Items: []models.SubscriptionItem{
				{
					PriceId:  priceID,
					Quantity: 1,
				},
			},
			CustomData: map[string]any{
				"tenant_id": tenantID,
				"plan_id":   planID,
			},
		}

		created, err := paddleClient.Subscriptions.Create(ctx, createReq)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("Failed to create Paddle subscription: %v", err)}
		}
		sub = created
	}

	subID := sub.SubscriptionId

	// Save back to DB
	updateData := Dict{
		"pay_provider_subscription_id":       subID,
		"pay_provider_subscription_metadata": toDict(sub),
		"pay_provider_last_sync_at":          time.Now(),
		"status":                             sub.Status,
		"next_billing_at":                    sub.NextBilledAt,
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
		"success":         true,
		"msg":             "Paddle subscription created/updated",
		"subscription_id": subID,
	}
}

// PaddleSyncAllSubscriptionStatuses – batch job example
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
        subID := sub["pay_provider_subscription_id"].(string)
        if subID == "" {
            continue
        }

        paddleSub, err := paddleClient.Subscriptions.Get(ctx, subID)
        if err != nil {
            log.Printf("Failed to get Paddle sub %s: %v", subID, err)
            continue
        }

        newStatus := string(paddleSub.Status) // e.g. "active", "past_due", ...
        nextBilling := paddleSub.NextBilledAt // *time.Time

        updateParams := map[string]any{
            "data": Dict{
                "table": "subscription",
                "filters": []Dict{{"field": "pay_provider_subscription_id", "value": subID}},
                "data": Dict{
                    "status": newStatus,
                    "next_billing_at": nextBilling,
                    "pay_provider_last_sync_at": time.Now(),
                    "pay_provider_subscription_metadata": toDict(paddleSub),
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
        "msg": fmt.Sprintf("Synced status for %d subscriptions", updatedCount),
    }
}