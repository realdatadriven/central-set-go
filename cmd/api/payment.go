package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/robfig/cron/v3" // go get github.com/robfig/cron/v3
	"github.com/stripe/stripe-go/v84"
	"github.com/stripe/stripe-go/v84/customer"
	"github.com/stripe/stripe-go/v84/event"
	"github.com/stripe/stripe-go/v84/invoice"
	"github.com/stripe/stripe-go/v84/price"
	"github.com/stripe/stripe-go/v84/product"
	"github.com/stripe/stripe-go/v84/subscription"
	"github.com/stripe/stripe-go/v84/webhook"
)

func init() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	if stripe.Key == "" {
		log.Fatal("STRIPE_SECRET_KEY not set")
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
func SyncOrCreateProduct(data map[string]any) (*stripe.Product, []*stripe.Price, error) {
	internalID, _ := data["product_id"].(string)
	name, _ := data["name"].(string)
	monthlyAmount, _ := data["monthly_amount"].(int64) // in cents
	annualAmount, _ := data["annual_amount"].(int64)
	currency, _ := data["currency"].(string)
	if currency == "" {
		currency = "usd"
	}

	// Create or retrieve product (Stripe doesn't have upsert, so check existence via metadata or your DB)
	prodParams := &stripe.ProductParams{
		Name:     stripe.String(name),
		Metadata: stripe.StringMap{"internal_id": internalID},
	}
	prod, err := product.New(prodParams)
	if err != nil {
		return nil, nil, err
	}

	var prices []*stripe.Price

	// Monthly price
	if monthlyAmount > 0 {
		priceParams := &stripe.PriceParams{
			Product:    stripe.String(prod.ID),
			Currency:   stripe.String(currency),
			UnitAmount: stripe.Int64(monthlyAmount),
			Recurring:  &stripe.PriceRecurringParams{Interval: stripe.String("month")},
		}
		p, err := price.New(priceParams)
		if err != nil {
			return nil, nil, err
		}
		prices = append(prices, p)
	}

	// Annual price
	if annualAmount > 0 {
		priceParams := &stripe.PriceParams{
			Product:    stripe.String(prod.ID),
			Currency:   stripe.String(currency),
			UnitAmount: stripe.Int64(annualAmount),
			Recurring:  &stripe.PriceRecurringParams{Interval: stripe.String("year")},
		}
		p, err := price.New(priceParams)
		if err != nil {
			return nil, nil, err
		}
		prices = append(prices, p)
	}

	// Update your DB with Stripe IDs
	var monthlyID, annualID string
	if len(prices) > 0 {
		monthlyID = prices[0].ID
	}
	if len(prices) > 1 {
		annualID = prices[1].ID
	}
	if err := updateProductStripeIDs(internalID, prod.ID, monthlyID, annualID); err != nil {
		log.Printf("Warning: DB update failed: %v", err)
	}

	return prod, prices, nil
}

// CreateOrSyncCustomer
func CreateOrSyncCustomer(data map[string]any) (*stripe.Customer, error) {
	tenantID, _ := data["tenant_id"].(string)
	name, _ := data["name"].(string)
	email, _ := data["email"].(string)

	params := &stripe.CustomerParams{
		Name:  stripe.String(name),
		Email: stripe.String(email),
		Metadata: stripe.StringMap{
			"tenant_id": tenantID,
		},
	}
	cust, err := customer.New(params)
	if err != nil {
		return nil, err
	}

	// Update DB
	if err := updateTenantStripeCustomerID(tenantID, cust.ID); err != nil {
		log.Printf("Warning: DB update failed: %v", err)
	}

	return cust, nil
}

// CreateOrUpdateSubscription
// data should include tenant_id (or stripe_customer_id), price_id (Stripe or internal), etc.
func CreateOrUpdateSubscription(data map[string]any) (*stripe.Subscription, error) {
	stripeCustID, _ := data["stripe_customer_id"].(string)
	stripePriceID, _ := data["stripe_price_id"].(string)
	subInternalID, _ := data["subscription_id"].(string) // your DB ID
	// Optional: trial_days, quantity, etc.

	if stripeCustID == "" {
		return nil, fmt.Errorf("stripe_customer_id required")
	}
	if stripePriceID == "" {
		return nil, fmt.Errorf("stripe_price_id required")
	}

	params := &stripe.SubscriptionParams{
		Customer: stripe.String(stripeCustID),
		Items: []*stripe.SubscriptionItemsParams{
			{Price: stripe.String(stripePriceID)},
		},
		// Add more from data map as needed
	}
	sub, err := subscription.New(params)
	if err != nil {
		return nil, err
	}

	// Update DB with Stripe subscription ID and price
	if err := updateSubscriptionStripeIDs(subInternalID, sub.ID, stripePriceID, string(sub.Status)); err != nil {
		log.Printf("Warning: DB update failed: %v", err)
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
		log.Printf("Invoice paid: %s for sub %s", inv.ID, inv.Subscription.ID)

	case "invoice.payment_failed":
		var inv stripe.Invoice
		if err := json.Unmarshal(evt.Data.Raw, &inv); err != nil {
			return err
		}
		log.Printf("Payment failed: %s (attempt %d), sub %s", inv.ID, inv.AttemptCount, inv.Subscription.ID)
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
		Limit:  stripe.Int64(100),
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
		Status: stripe.String("open"),
		DueDate: &stripe.InvoiceListDueDateParams{
			Lt: stripe.Int64(time.Now().Unix()),
		},
		Limit: stripe.Int64(50),
	}
	invIter := invoice.List(invParams)
	for invIter.Next() {
		inv := invIter.Invoice()
		log.Printf("Reconcile: overdue invoice %s for sub %s", inv.ID, inv.Subscription.ID)
		// Handle: notify, mark in DB, etc.
	}

	log.Println("Reconciliation complete")
	return nil
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