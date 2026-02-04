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

func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		//log.Fatal("STRIPE_SECRET_KEY not set")
	}
}

// ------------------- DB Update Placeholders -------------------
// Replace these with your actual DB operations (e.g., SQL UPDATE or GORM Save)

func updateProductStripeIDs(internalProductID string, stripeProductID, stripeMonthlyPriceID, stripeAnnualPriceID string) error {
	// e.g., UPDATE products SET stripe_product_id = ?, monthly_price_id = ?, annual_price_id = ? WHERE id = ?
	log.Printf("DB: Updated product %s with Stripe IDs: prod=%s, monthly=%s, annual=%s", internalProductID, stripeProductID, stripeMonthlyPriceID, stripeAnnualPriceID)
	return nil // return real error
}

func updateTenantStripeCustomerID(tenantID, stripeCustomerID string) error {
	// e.g., UPDATE tenants SET stripe_customer_id = ? WHERE tenant_id = ?
	log.Printf("DB: Updated tenant %s with stripe_customer_id=%s", tenantID, stripeCustomerID)
	return nil
}

func updateSubscriptionStripeIDs(subInternalID, stripeSubID, stripePriceID string, status string) error {
	// e.g., UPDATE subscriptions SET stripe_subscription_id = ?, stripe_price_id = ?, status = ? WHERE id = ?
	log.Printf("DB: Updated subscription %s: subID=%s, priceID=%s, status=%s", subInternalID, stripeSubID, stripePriceID, status)
	return nil
}

// ------------------- Core Sync/Create Functions -------------------

// SyncOrCreateProduct takes map data from your products + subscription tables
// Returns created/updated product and prices; updates DB with Stripe IDs
// SyncOrCreateProduct — idempotent version
func SyncOrCreateProduct(data map[string]any) (*stripe.Product, []*stripe.Price, error) {
	// Extract identifiers from your DB
	internalID, _ := data["product_id"].(string)
	stripeProductID, _ := data["stripe_product_id"].(string) // already saved?

	name, _ := data["name"].(string)
	monthlyAmount, _ := data["monthly_amount"].(int64)
	annualAmount, _ := data["annual_amount"].(int64)
	currency, _ := data["currency"].(string)
	if currency == "" {
		currency = "usd"
	}

	var prod *stripe.Product
	var err error

	// 1. If we already have stripe_product_id → retrieve it
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
		params := &stripe.ProductParams{
			Name:     stripe.String(name),
			Metadata: map[string]string{"internal_id": internalID},
		}
		prod, err = product.New(params)
		if err != nil {
			return nil, nil, fmt.Errorf("failed to create product: %w", err)
		}
		stripeProductID = prod.ID

		// Immediately save to DB
		if err := updateProductStripeIDs(internalID, stripeProductID, "", ""); err != nil {
			log.Printf("Warning: failed to save product ID to DB: %v", err)
		}
	}

	// Now handle prices (we'll check existing prices later if needed)
	var prices []*stripe.Price

	// Helper to create or get price
	createOrGetPrice := func(interval string, amount int64) (*stripe.Price, error) {
		if amount <= 0 {
			return nil, nil
		}

		// In real world, you might also store stripe_price_monthly_id etc. in DB
		// For simplicity here we create new if not matching — in production you'd list prices

		params := &stripe.PriceParams{
			Product:    stripe.String(prod.ID),
			Currency:   stripe.String(currency),
			UnitAmount: stripe.Int64(amount),
			Recurring: &stripe.PriceRecurringParams{
				Interval: stripe.String(interval),
			},
		}
		p, err := price.New(params)
		if err != nil {
			return nil, err
		}
		return p, nil
	}

	monthly, err := createOrGetPrice("month", monthlyAmount)
	if err != nil {
		return nil, nil, err
	}
	if monthly != nil {
		prices = append(prices, monthly)
	}

	annual, err := createOrGetPrice("year", annualAmount)
	if err != nil {
		return nil, nil, err
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
	if err := updateProductStripeIDs(internalID, stripeProductID, monthlyID, annualID); err != nil {
		log.Printf("Warning: failed to save price IDs: %v", err)
	}

	return prod, prices, nil
}

// CreateOrSyncCustomer
func CreateOrSyncCustomer(data map[string]any) (*stripe.Customer, error) {
	tenantID, _ := data["tenant_id"].(string)
	stripeCustomerID, _ := data["stripe_customer_id"].(string) // already saved?
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
					return nil, err
				}
			}
			log.Printf("Using existing customer: %s", cust.ID)
			return cust, nil
		}
	}

	// Create new
	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
		Metadata: map[string]string{
			"tenant_id": tenantID,
		},
	}
	cust, err = customer.New(params)
	if err != nil {
		return nil, err
	}

	// Save to DB
	if err := updateTenantStripeCustomerID(tenantID, cust.ID); err != nil {
		log.Printf("Warning: failed to save customer ID: %v", err)
	}

	return cust, nil
}

// CreateOrUpdateSubscription
// data should include tenant_id (or stripe_customer_id), price_id (Stripe or internal), etc.
func CreateOrUpdateSubscription(data map[string]any) (*stripe.Subscription, error) {
	//tenantID, _ := data["tenant_id"].(string)
	stripeCustomerID, _ := data["stripe_customer_id"].(string)
	stripePriceID, _ := data["stripe_price_id"].(string)
	stripeSubscriptionID, _ := data["stripe_subscription_id"].(string) // already have?
	subInternalID, _ := data["subscription_id"].(string)

	if stripeCustomerID == "" {
		return nil, fmt.Errorf("stripe_customer_id is required")
	}
	if stripePriceID == "" {
		return nil, fmt.Errorf("stripe_price_id is required")
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
					return nil, err
				}
			}
			log.Printf("Using existing subscription: %s (status: %s)", sub.ID, sub.Status)
			return sub, nil
		}
	}

	// 2. Create new subscription
	params := &stripe.SubscriptionParams{
		Customer: stripe.String(stripeCustomerID),
		Items: []*stripe.SubscriptionItemsParams{
			{Price: stripe.String(stripePriceID)},
		},
		// Add trial, proration, etc. from data map if needed
	}
	sub, err = subscription.New(params)
	if err != nil {
		return nil, err
	}

	// 3. Save to DB
	if err := updateSubscriptionStripeIDs(subInternalID, sub.ID, stripePriceID, string(sub.Status)); err != nil {
		log.Printf("Warning: failed to save subscription ID: %v", err)
	}

	return sub, nil
}

// ------------------- Webhook & Event Processing -------------------

func processStripeEvent(evt *stripe.Event) error {
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

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := processStripeEvent(&evt); err != nil {
		log.Printf("Event processing error: %v", err)
		http.Error(w, "Processing error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// ------------------- Cron Reconciliation (Regular Checks) -------------------

func ReconcileRecentIssues() error {
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
// You can later enhance with customer/tenant name lookup from your DB using metadata or stored stripe_customer_id
func ListFuturePayments(nMonths int, includeTrialing bool) (*FuturePayments, error) {
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
// data map example: {"stripe_customer_id": "cus_...", "meter_event_name": "cloud_compute_hours", "value": 150, "timestamp": unixTime, "description": "VM usage March"}
func ReportCloudUsage(data map[string]any) error {
	customerID, ok := data["stripe_customer_id"].(string)
	if !ok || customerID == "" {
		return fmt.Errorf("stripe_customer_id required")
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
			"stripe_customer_id": customerID,
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
// data map example: {"stripe_customer_id": "cus_...", "stripe_subscription_id": "sub_...", "amount_cents": 1500, "description": "Cloud storage March: 75 GB", "currency": "usd"}
func AddVariableInvoiceItem(data map[string]any) (*stripe.InvoiceItem, error) {
	customerID, _ := data["stripe_customer_id"].(string)
	subID, _ := data["stripe_subscription_id"].(string) // optional: attach to sub's next invoice
	amount, _ := data["amount_cents"].(int64)
	desc, _ := data["description"].(string)
	currency, _ := data["currency"].(string)
	if currency == "" {
		currency = "usd"
	}

	if customerID == "" || amount <= 0 || desc == "" {
		return nil, fmt.Errorf("required: stripe_customer_id, amount_cents >0, description")
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
//	  "stripe_customer_id": "cus_...",
//	  "stripe_subscription_id": "sub_...",  // optional if already exists
//	  "resources": []map[string]any{
//	    {"type": "storage", "enabled": true},   // will use price_storage
//	    {"type": "compute", "enabled": true},
//	  }
//	}
func AddOrUpdateResourceBilling(data map[string]any) (*stripe.Subscription, error) {
	custID := data["stripe_customer_id"].(string)
	subID, _ := data["stripe_subscription_id"].(string)

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

/*func main() {
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
