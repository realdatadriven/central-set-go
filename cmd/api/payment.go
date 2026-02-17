package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	// go get github.com/robfig/cron/v3
	"github.com/stripe/stripe-go/v84"
	//"github.com/stripe/stripe-go/v84/billingmeterevent"
	"github.com/stripe/stripe-go/v84/customer"
	"github.com/stripe/stripe-go/v84/invoice"
	"github.com/stripe/stripe-go/v84/invoiceitem"
	"github.com/stripe/stripe-go/v84/price"
	"github.com/stripe/stripe-go/v84/product"
	"github.com/stripe/stripe-go/v84/subscription"
	"github.com/stripe/stripe-go/v84/webhook"
)

func _init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		fmt.Println("STRIPE_SECRET_KEY not set")
	} else {
		fmt.Println("Stripe initialized with secret key", stripe.Key)
	}
}

func (app *application) PaymentInit() {
	if stripe.Key == "" {
		stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
		if stripe.Key == "" {
			fmt.Println("STRIPE_SECRET_KEY not set")
		} else {
			fmt.Println("Stripe initialized with secret key", stripe.Key)
		}
	}
}

// ------------------- DB Update Placeholders -------------------
// Replace these with your actual DB operations (e.g., SQL UPDATE or GORM Save)

func (app *application) updateProductStripeIDs(internalProductID string, stripeProductID, stripeMonthlyPriceID, stripeAnnualPriceID string) error {
	// e.g., UPDATE products SET payment_product_id = ?, monthly_price_id = ?, annual_price_id = ? WHERE id = ?
	//query := "UPDATE plan SET payment_product_id = ? WHERE plan_id = ?"
	//newDB, err := app.GetDBNameFromParams(params)
	//log.Printf("DB: Updated product %s with Stripe IDs: prod=%s, monthly=%s, annual=%s", internalProductID, stripeProductID, stripeMonthlyPriceID, stripeAnnualPriceID)
	return nil // return real error
}

func (app *application) updateTenantStripeCustomerID(tenantID, stripeCustomerID string) error {
	// e.g., UPDATE tenants SET payment_customer_id = ? WHERE tenant_id = ?
	log.Printf("DB: Updated tenant %s with payment_customer_id=%s", tenantID, stripeCustomerID)
	return nil
}

func (app *application) updateSubscriptionStripeIDs(subInternalID, stripeSubID, stripePriceID string, status string) error {
	// e.g., UPDATE subscriptions SET payment_subs_id = ?, stripe_price_id = ?, status = ? WHERE id = ?
	log.Printf("DB: Updated subscription %s: subID=%s, priceID=%s, status=%s", subInternalID, stripeSubID, stripePriceID, status)
	return nil
}

// ------------------- Core Sync/Create Functions -------------------

// SyncOrCreateProduct takes map data from your products + subscription tables
// Returns created/updated product and prices; updates DB with Stripe IDs
// SyncOrCreateProduct — idempotent version
func (app *application) SyncOrCreateProduct(params map[string]any) Dict {
	// Extract identifiers from your DB
	data := Dict{}
	if _, ok := params["data"].(Dict); ok {
		data = params["data"].(Dict)["data"].(Dict) // adjust based on your actual params structure
	}
	internalID, _ := data["plan_id"].(string)
	name, _ := data["plan"].(string)
	stripeProductID, _ := data["payment_product_id"].(string) // already saved?
	//fmt.Println(name, data["plan"], data["plan_id"], data["currency"])
	monthlyAmount, _ := data["monthly_amount"].(float64)
	annualAmount, _ := data["annual_amount"].(float64)
	currency, _ := data["currency"].(string)
	if currency == "" {
		currency = "usd"
	}
	stripePriceMonthlyID, _ := data["payment_price_monthly_id"].(string)
	stripePriceAnnualID, _ := data["payment_annual_id"].(string)
	var prod *stripe.Product
	var err error
	// 1. If we already have payment_product_id → retrieve it
	if stripeProductID != "" {
		prod, err = product.Get(stripeProductID, nil)
		if err != nil {
			// If not found → treat as not existing (maybe was deleted)
			log.Printf("Product %s not found in Stripe, will recreate: %v", stripeProductID, err)
			stripeProductID = ""
		} else {
			log.Printf("Using existing product: %s", prod.ID)
		}
	}
	// 2. Create if not found
	if stripeProductID == "" {
		sparams := &stripe.ProductParams{
			Name:     stripe.String(name),
			Metadata: map[string]string{"internal_id": internalID},
		}
		prod, err = product.New(sparams)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("failed to create product: %w", err)}
		}
		stripeProductID = prod.ID
	} else {
		prod.Name = name // Optional: update name if changed
		prod, err = product.Update(stripeProductID, &stripe.ProductParams{
			Name: stripe.String(name),
		})
		stripeProductID = prod.ID
	}
	// Now handle prices (we'll check existing prices later if needed)
	var prices []*stripe.Price
	// Helper to create or get price
	createOrGetPrice := func(id, interval string, amount float64) (*stripe.Price, error) {
		if amount <= 0 {
			return nil, nil
		}
		// In real world, you might also store stripe_price_monthly_id etc. in DB
		// For simplicity here we create new if not matching — in production you'd list prices
		sparams := &stripe.PriceParams{
			Product:  stripe.String(prod.ID),
			Currency: stripe.String(currency),
			//UnitAmount:        stripe.Int64(int64(amount)),
			UnitAmountDecimal: stripe.Float64(amount * 100),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String(interval),
			},
		}
		if id != "" {
			p, err := price.Get(id, nil)
			if err != nil {
				p, err = price.New(sparams)
				if err != nil {
					fmt.Printf("PRICE:", id, interval, amount)
					return nil, err
				}
				return p, nil
			}
			_, err = price.Update(id, sparams)
			if err != nil {
				fmt.Printf("PRICE:", id, interval, amount)
				return nil, err
			}
			p, _ = price.Get(id, nil)
			return p, nil
		} else {
			p, err := price.New(sparams)
			if err != nil {
				fmt.Printf("PRICE:", id, interval, amount)
				return nil, err
			}
			return p, nil
		}
	}
	monthly, err := createOrGetPrice(stripePriceMonthlyID, "month", monthlyAmount)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to create/get monthly price: %w", err)}
	}
	if monthly != nil {
		prices = append(prices, monthly)
	}
	annual, err := createOrGetPrice(stripePriceAnnualID, "year", annualAmount)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to create/get annual price: %w", err)}
	}
	if annual != nil {
		prices = append(prices, annual)
	}
	// Save price IDs if you store them separately
	var monthlyID, annualID string
	if len(prices) > 0 {
		monthlyID = prices[0].ID
	}
	if len(prices) > 1 {
		annualID = prices[1].ID
	}
	params["data"].(Dict)["table"] = "plan" // update params for DB save
	_data := data
	payment_product_metadata, err := json.Marshal(prod)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to marshal product metadata: %w", err)}
	}
	_data["payment_product_id"] = stripeProductID
	_data["payment_price_monthly_id"] = monthlyID
	_data["payment_price_annual_id"] = annualID
	_data["payment_product_metadata"] = string(payment_product_metadata)
	_data["payment_product_last_sync_at"] = time.Now() //.Format(time.RFC3339)
	params["data"].(Dict)["data"] = _data
	res := app.create_update(params)
	if _, ok := res["success"].(bool); !ok {
		return res
	}
	//if err := app.updateProductStripeIDs(internalID, stripeProductID, monthlyID, annualID); err != nil {
	//	log.Printf("Warning: failed to save price IDs: %v", err)
	//}
	// fmt.Print(Dict{"success": true, "msg": "Product and prices created/updated", "product": prod, "prices": prices})
	return Dict{"success": true, "msg": "Product and prices created/updated", "product": prod, "prices": prices}
}

// CreateOrSyncCustomer
func (app *application) CreateOrSyncCustomer(params map[string]any) Dict {
	data := Dict{}
	if _, ok := params["data"].(Dict); ok {
		data = params["data"].(Dict)["data"].(Dict) // adjust based on your actual params structure
	}
	tenantID, _ := data["tenant_id"].(string)
	stripeCustomerID, _ := data["payment_customer_id"].(string) // already saved?
	name, _ := data["name"].(string)
	email, _ := data["email"].(string)
	var cust *stripe.Customer
	var err error
	// If we have ID → retrieve
	if stripeCustomerID != "" {
		cust, err = customer.Get(stripeCustomerID, nil)
		if err != nil {
			log.Printf("Customer %s not found, will recreate: %v", stripeCustomerID, err)
			stripeCustomerID = ""
		} else {
			// Optional: update email/name if changed
			if cust.Email != email || cust.Name != name {
				updateParams := &stripe.CustomerParams{
					Email: stripe.String(email),
					Name:  stripe.String(name),
				}
				cust, err = customer.Update(stripeCustomerID, updateParams)
				if err != nil {
					return Dict{"success": false, "msg": fmt.Sprintf("failed to update customer: %w", err)}
				}
			}
			log.Printf("Using existing customer: %s", cust.ID)
			//return Dict{"success": true, "msg": "Customer retrieved", "customer": cust}
		}
	} else {
		// Create new
		cparams := &stripe.CustomerParams{
			Name:  stripe.String(name),
			Email: stripe.String(email),
			Metadata: map[string]string{
				"tenant_id": tenantID,
			},
		}
		cust, err = customer.New(cparams)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("failed to create customer: %w", err)}
		}
	}
	params["data"].(Dict)["table"] = "tenant" // update params for DB save
	_data := data
	payment_customer_metadata, err := json.Marshal(cust)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to marshal tenant metadata: %w", err)}
	}
	_data["payment_customer_metadata"] = string(payment_customer_metadata)
	_data["payment_customer_id"] = cust.ID
	_data["payment_customer_last_sync_at"] = time.Now() //.Format(time.RFC3339)
	params["data"].(Dict)["data"] = _data
	res := app.create_update(params)
	if _, ok := res["success"].(bool); !ok {
		return res
	}
	return Dict{"success": true, "msg": "Customer created", "customer": cust}
}

// CreateOrUpdateSubscription
// data should include tenant_id (or payment_customer_id), price_id (Stripe or internal), etc.
func (app *application) CreateOrUpdateSubscription(params map[string]any) Dict {
	data := Dict{}
	if _, ok := params["data"].(Dict); ok {
		data = params["data"].(Dict)["data"].(Dict) // adjust based on your actual params structure
	}
	//subInternalID, _ := data["subscription_id"].(string)
	params["data"].(Dict)["table"] = "plan"                                                   // ensure we know which table to update
	params["data"].(Dict)["filters"] = []Dict{{"field": "plan_id", "value": data["plan_id"]}} // for DB update
	_plan := app.read(params)                                                                 // read existing data for price updates if needed
	if _, ok := _plan["success"].(bool); !ok {
		return _plan
	}
	planData := Dict{}
	if _, ok := _plan["data"].([]Dict); !ok {
		return Dict{"success": false, "msg": "No price data found"}
	} else if len(_plan["data"].([]Dict)) == 0 {
		return Dict{"success": false, "msg": "Empty price data"}
	} else {
		planData = _plan["data"].([]Dict)[0]
	}
	stripePriceMonthlyID := planData["payment_price_monthly_id"].(string)
	stripePriceAnnualID := planData["payment_price_annual_id"].(string)
	params["data"].(Dict)["table"] = "tenant"                                                     // ensure we know which table to update
	params["data"].(Dict)["filters"] = []Dict{{"field": "tenant_id", "value": data["tenant_id"]}} // for DB update
	tenant := app.read(params)                                                                    // read existing data for price updates if needed
	if _, ok := tenant["success"].(bool); !ok {
		return tenant
	}
	tenantData := Dict{}
	if _, ok := tenant["data"].([]Dict); !ok {
		return Dict{"success": false, "msg": "No tenant data found"}
	} else if len(tenant["data"].([]Dict)) == 0 {
		return Dict{"success": false, "msg": "Empty tenant data"}
	} else {
		tenantData = tenant["data"].([]Dict)[0]
	}
	stripeCustomerID, _ := tenantData["payment_customer_id"].(string)
	stripePriceID := stripePriceMonthlyID
	if _, ok := data["recurring_interval_id"]; ok {
		if data["recurring_interval_id"] == any(2) { // annual
			stripePriceID = stripePriceAnnualID
		}
	}
	stripeSubscriptionID, _ := data["payment_subs_id"].(string) // already have?
	if stripeCustomerID == "" {
		//return nil, fmt.Errorf("payment_customer_id is required")
		return Dict{"success": false, "msg": "payment_customer_id is required"}
	}
	if stripePriceID == "" {
		return Dict{"success": false, "msg": "payment_price_id is required"}
	}
	var sub *stripe.Subscription
	var err error
	// 1. If we already have subscription ID → retrieve and possibly update
	if stripeSubscriptionID != "" {
		sub, err = subscription.Get(stripeSubscriptionID, nil)
		if err != nil {
			log.Printf("Subscription %s not found, will create new: %v", stripeSubscriptionID, err)
			stripeSubscriptionID = ""
		} else {
			// Optional: update if price changed, etc.
			if sub.Items.Data[0].Price.ID != stripePriceID {
				updateParams := &stripe.SubscriptionParams{
					Items: []*stripe.SubscriptionItemsParams{
						{
							ID:    stripe.String(sub.Items.Data[0].ID), // existing item
							Price: stripe.String(stripePriceID),
						},
					},
				}
				sub, err = subscription.Update(stripeSubscriptionID, updateParams)
				if err != nil {
					return Dict{"success": false, "msg": fmt.Sprintf("failed to update subscription: %w", err)}
				}
			}
			log.Printf("Using existing subscription: %s (status: %s)", sub.ID, sub.Status)
			//return sub, nil
		}
	} else {
		// 2. Create new subscription
		fmt.Printf("SUBS:", stripeCustomerID, stripePriceID)
		sparams := &stripe.SubscriptionParams{
			Customer: stripe.String(stripeCustomerID),
			Items: []*stripe.SubscriptionItemsParams{
				{Price: stripe.String(stripePriceID)},
			},
			// Add trial, proration, etc. from data map if needed
		}
		sub, err = subscription.New(sparams)
		if err != nil {
			return Dict{"success": false, "msg": fmt.Sprintf("failed to create subscription: %w", err)}
		}
	}
	params["data"].(Dict)["table"] = "subscription" // update params for DB save
	_data := data
	payment_subscription_metadata, err := json.Marshal(sub)
	if err != nil {
		return Dict{"success": false, "msg": fmt.Sprintf("failed to marshal subscription metadata: %w", err)}
	}
	_data["payment_subs_metadata"] = string(payment_subscription_metadata)
	_data["payment_subs_id"] = sub.ID
	_data["payment_subs_last_sync_at"] = time.Now() //.Format(time.RFC3339)
	params["data"].(Dict)["data"] = _data
	res := app.create_update(params)
	if _, ok := res["success"].(bool); !ok {
		return res
	}
	return Dict{"success": true, "msg": "Subscription created", "subscription": sub}
}

// ------------------- Webhook & Event Processing -------------------

func (app *application) processStripeEvent(evt *stripe.Event) error {
	switch evt.Type {
	case "invoice.paid":
		var inv stripe.Invoice
		if err := json.Unmarshal(evt.Data.Raw, &inv); err != nil {
			return err
		}
		// Update DB: mark paid, update next billing date, activate tenant access
		//log.Printf("Invoice paid: %s for sub %s", inv.ID, inv.Subscription.ID)

	case "invoice.payment_failed":
		var inv stripe.Invoice
		if err := json.Unmarshal(evt.Data.Raw, &inv); err != nil {
			return err
		}
		//log.Printf("Payment failed: %s (attempt %d), sub %s", inv.ID, inv.AttemptCount, inv.Subscription.ID)
		// Notify tenant, mark delayed in DB

	case "customer.subscription.updated":
		var sub stripe.Subscription
		if err := json.Unmarshal(evt.Data.Raw, &sub); err != nil {
			return err
		}
		// Update DB status (active, past_due, canceled, etc.)
		log.Printf("Sub %s updated to %s", sub.ID, sub.Status)

	// Add more cases as needed
	default:
		log.Printf("Ignored event: %s", evt.Type)
	}
	return nil
}

func (app *application) WebhookHandler(w http.ResponseWriter, r *http.Request) {
	const MaxBody = 65536
	r.Body = http.MaxBytesReader(w, r.Body, MaxBody)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Read error", http.StatusServiceUnavailable)
		return
	}

	sig := r.Header.Get("Stripe-Signature")
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")

	evt, err := webhook.ConstructEvent(payload, sig, endpointSecret)
	if err != nil {
		http.Error(w, "Signature failed", http.StatusBadRequest)
		return
	}

	if err := app.processStripeEvent(&evt); err != nil {
		log.Printf("Event processing error: %v", err)
		http.Error(w, "Processing error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ------------------- Cron Reconciliation (Regular Checks) -------------------

func (app *application) ReconcileRecentIssues() error {
	log.Println("Running Stripe reconciliation...")

	// Example: Check for past_due or unpaid subscriptions
	params := &stripe.SubscriptionListParams{
		Status: stripe.String("past_due"),
		//Limit:  stripe.Int64(100),
	}
	iter := subscription.List(params)
	for iter.Next() {
		sub := iter.Subscription()
		// Simulate event processing or direct DB update
		log.Printf("Reconcile: past_due sub %s (customer %s)", sub.ID, sub.Customer.ID)
		// You can call processStripeEvent with a fake event or directly update DB
		// e.g., updateSubscriptionStripeIDs(..., sub.ID, ..., string(sub.Status))
	}

	// Also check open invoices past due
	invParams := &stripe.InvoiceListParams{
		Status:  stripe.String("open"),
		DueDate: stripe.Int64(time.Now().Unix()),
		DueDateRange: &stripe.RangeQueryParams{
			GreaterThanOrEqual: *stripe.Int64(2), //stripe.Int64(time.Now().Unix()),

		},
		//Limit: stripe.Int64(50),
	}
	invIter := invoice.List(invParams)
	for invIter.Next() {
		inv := invIter.Invoice()
		log.Printf("Reconcile: overdue invoice %s for sub %s", inv.ID, inv.Status)
		// Handle: notify, mark in DB, etc.
	}

	log.Println("Reconciliation complete")
	return nil
}

type FuturePayment struct {
	ExpectedDate   time.Time `json:"expected_date"`
	AmountCents    int64     `json:"amount_cents"`
	AmountDisplay  string    `json:"amount_display"` // e.g. "19.99 USD"
	CustomerName   string    `json:"customer_name"`  // or tenant name if you add lookup
	SubscriptionID string    `json:"subscription_id"`
	Currency       string    `json:"currency"`
	Interval       string    `json:"interval"` // month / year
}

type FuturePayments struct {
	Payments    []FuturePayment `json:"payments"`
	TotalCount  int             `json:"total_count"`
	NextNMonths int             `json:"next_n_months"`
	AsOf        time.Time       `json:"as_of"`
}

// ListFuturePayments returns a sorted list of expected future charges for the next nMonths
// includeTrialing: count trials that will start billing soon
// You can later enhance with customer/tenant name lookup from your DB using metadata or stored payment_customer_id
func (app *application) ListFuturePayments(nMonths int, includeTrialing bool) (*FuturePayments, error) {
	if nMonths < 1 || nMonths > 36 {
		return nil, fmt.Errorf("nMonths should be 1–36")
	}

	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	// List active (and trialing) subscriptions
	params := &stripe.SubscriptionListParams{
		Status: stripe.String("all"),

		//Limit: stripe.Int64(100), // paginate in production if >100
	}
	iter := subscription.List(params)

	var payments []FuturePayment
	now := time.Now()
	endHorizon := now.AddDate(0, nMonths+1, 0) // a bit extra to catch edge cases

	for iter.Next() {
		sub := iter.Subscription()

		if sub.Status != stripe.SubscriptionStatusActive &&
			!(includeTrialing && sub.Status == stripe.SubscriptionStatusTrialing) {
			continue
		}

		if len(sub.Items.Data) == 0 || sub.Items.Data[0].Price == nil {
			continue // skip weird/empty subs
		}

		price := sub.Items.Data[0].Price
		if price.Recurring == nil {
			continue // not recurring
		}

		amount := price.UnitAmount
		if amount <= 0 {
			continue
		}

		interval := price.Recurring.Interval
		intervalCount := price.Recurring.IntervalCount
		if intervalCount == 0 {
			intervalCount = 1
		}

		nextBill := time.Unix(sub.EndedAt, 0)

		// For trialing subs: first real bill is at trial end
		if sub.Status == stripe.SubscriptionStatusTrialing {
			nextBill = time.Unix(sub.TrialEnd, 0)
		}

		for nextBill.Before(endHorizon) {
			if nextBill.After(now) {
				payments = append(payments, FuturePayment{
					ExpectedDate:   nextBill,
					AmountCents:    amount,
					AmountDisplay:  fmt.Sprintf("%.2f %s", float64(amount)/100, price.Currency),
					CustomerName:   sub.Customer.Name, // fallback; better to lookup from your DB
					SubscriptionID: sub.ID,
					Currency:       string(price.Currency),
					//Interval:       interval,
				})
			}

			// Advance to next cycle
			switch interval {
			case "month":
				nextBill = nextBill.AddDate(0, int(intervalCount), 0)
			case "year":
				nextBill = nextBill.AddDate(int(intervalCount), 0, 0)
			default:
				break // unsupported → stop projecting
			}
		}
	}

	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("subscription list error: %w", err)
	}

	// Sort by date ascending
	sort.Slice(payments, func(i, j int) bool {
		return payments[i].ExpectedDate.Before(payments[j].ExpectedDate)
	})

	return &FuturePayments{
		Payments:    payments,
		TotalCount:  len(payments),
		NextNMonths: nMonths,
		AsOf:        now,
	}, nil
}

// ReportCloudUsage - Reports consumption to Stripe Meter for variable billing
// data map example: {"payment_customer_id": "cus_...", "meter_event_name": "cloud_compute_hours", "value": 150, "timestamp": unixTime, "description": "VM usage March"}
func (app *application) ReportCloudUsage(data map[string]any) error {
	customerID, ok := data["payment_customer_id"].(string)
	if !ok || customerID == "" {
		return fmt.Errorf("payment_customer_id required")
	}
	eventName, ok := data["meter_event_name"].(string)
	if !ok || eventName == "" {
		return fmt.Errorf("meter_event_name required (e.g., 'cloud_compute_hours')")
	}
	_, ok = data["value"].(float64) // or int64, must be number value
	if !ok {
		return fmt.Errorf("value required (numeric usage amount)")
	}
	ts, _ := data["timestamp"].(int64) // optional Unix timestamp, defaults to now
	if ts == 0 {
		ts = time.Now().Unix()
	}

	/*params := &stripe.BillingMeterEventParams{
		EventName: stripe.String(eventName),
		Payload: stripe.Map{
			"payment_customer_id": customerID,
			"value":              value,
			// optional: add dimensions e.g. "region": "eu", "resource_id": "vm_123"
		},
		Timestamp: stripe.Int64(ts),
	}
	// Optional: IdempotencyKey: stripe.String("usage-" + tenantID + "-" + time.Now().Format("20060102")),
	params.SetIdempotencyKey(fmt.Sprintf("usage-%s-%d", customerID, ts))

	/_, err := billingmeterevent.New(params)
	if err != nil {
		return fmt.Errorf("failed to report meter event: %w", err)
	}

	log.Printf("Reported %f %s for customer %s", value, eventName, customerID)*/
	return nil
}

// AddVariableInvoiceItem - Adds a monthly variable charge (e.g., cloud usage)
// data map example: {"payment_customer_id": "cus_...", "payment_subs_id": "sub_...", "amount_cents": 1500, "description": "Cloud storage March: 75 GB", "currency": "usd"}
func (app *application) AddVariableInvoiceItem(data map[string]any) (*stripe.InvoiceItem, error) {
	customerID, _ := data["payment_customer_id"].(string)
	subID, _ := data["payment_subs_id"].(string) // optional: attach to sub's next invoice
	amount, _ := data["amount_cents"].(int64)
	desc, _ := data["description"].(string)
	currency, _ := data["currency"].(string)
	if currency == "" {
		currency = "usd"
	}

	if customerID == "" || amount <= 0 || desc == "" {
		return nil, fmt.Errorf("required: payment_customer_id, amount_cents >0, description")
	}

	params := &stripe.InvoiceItemParams{
		Customer:     stripe.String(customerID),
		Amount:       stripe.Int64(amount),
		Currency:     stripe.String(currency),
		Description:  stripe.String(desc),
		Subscription: stripe.String(subID), // attaches to this sub's pending invoice
		// Optional: Period: &stripe.InvoiceItemPeriodParams{Start: ..., End: ...} for display
	}
	item, err := invoiceitem.New(params)
	if err != nil {
		return nil, fmt.Errorf("failed to add invoice item: %w", err)
	}

	log.Printf("Added variable charge %d %s to customer %s: %s", amount, currency, customerID, desc)
	return item, nil
}

// AddOrUpdateResourceBilling - Attaches/updates metered prices for enabled resources
// data example:
//
//	{
//	  "payment_customer_id": "cus_...",
//	  "payment_subs_id": "sub_...",  // optional if already exists
//	  "resources": []map[string]any{
//	    {"type": "storage", "enabled": true},   // will use price_storage
//	    {"type": "compute", "enabled": true},
//	  }
//	}
func (app *application) AddOrUpdateResourceBilling(data map[string]any) (*stripe.Subscription, error) {
	custID := data["payment_customer_id"].(string)
	subID, _ := data["payment_subs_id"].(string)

	// Fetch your DB-mapped price IDs (example helper)
	priceMap := map[string]string{
		"storage": "price_abc123_storage", // your fixed price IDs
		"compute": "price_def456_compute",
		// ...
	}

	items := []*stripe.SubscriptionItemsParams{}
	for _, res := range data["resources"].([]map[string]any) {
		typ := res["type"].(string)
		enabled := res["enabled"].(bool)

		if !enabled {
			continue
		}

		priceID, ok := priceMap[typ]
		if !ok {
			return nil, fmt.Errorf("no price for resource type: %s", typ)
		}

		items = append(items, &stripe.SubscriptionItemsParams{
			Price: stripe.String(priceID),
			// Optional: Quantity: stripe.Int64(1), but for metered usually leave as-is or use metadata
		})
	}

	if len(items) == 0 {
		return nil, fmt.Errorf("no resources enabled")
	}

	var sub *stripe.Subscription
	var err error

	if subID != "" {
		// Update existing sub (add/remove items)
		params := &stripe.SubscriptionParams{
			Items: items, // This will REPLACE existing items → be careful!
			// To add without replacing: fetch current items first, merge, then update
			ProrationBehavior: stripe.String("create_prorations"),
		}
		sub, err = subscription.Update(subID, params)
	} else {
		// Create new sub with fixed plan + resources
		params := &stripe.SubscriptionParams{
			Customer: stripe.String(custID),
			Items:    items, // add your fixed monthly price here too if needed
		}
		sub, err = subscription.New(params)
	}

	if err != nil {
		return nil, err
	}

	// Save sub.ID back to your subscriptions or tenants table if new
	return sub, nil
}

/*func (app *application) main() {
	// HTTP routes
	http.HandleFunc("/webhook", WebhookHandler)
	// Add admin routes, e.g.:
	// http.HandleFunc("/sync-product", func(w, r) { ... call SyncOrCreateProduct(data) ... })
	// http.HandleFunc("/create-sub", func(w, r) { ... call CreateOrUpdateSubscription(data) ... })

	// Cron setup (example: run every hour)
	c := cron.New()
	c.AddFunc("@hourly", func() {
		if err := ReconcileRecentIssues(); err != nil {
			log.Printf("Cron error: %v", err)
		}
	})
	c.Start()

	log.Println("Server on :8080 | Cron active")
	log.Fatal(http.ListenAndServe(":8080", nil))
}*/
