<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD012 -->
<!-- markdownlint-disable MD047 -->
# STORE_MODEL
```yaml
name: STORE
description: Basic online store backend model with catalog, customers, orders, payments, shipping, promotions and workflow automation
runs_as: MODEL
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
create_all: checkfirst
_drop_all: checkfirst
update_table_metadata: true
active: true
cs_app:
  Dashboards:
    menu_icon: document-report
    menu_order: 1
    active: true
    menu_config: '{"label": "dashboard","tooltip": "dashboard_desc","load_items": {"table": "dashboard","tables": ["dashboard"]}}'
    tables:
      - dashboard
  Customers:
    menu_icon: users
    menu_order: 2
    active: true
    tables:
      - customer
      - customer_address
  Catalog:
    menu_icon: shopping-bag
    menu_order: 3
    active: true
    tables:
      - product_category
      - product
  Orders:
    menu_icon: shopping-cart
    menu_order: 4
    active: true
    tables:
      - order_status
      - order_header
      - order_item
      - shipment
      - order_payment
  Marketing:
    menu_icon: megaphone
    menu_order: 5
    active: true
    tables:
      - promo_code
  Automation:
    menu_icon: sparkles
    menu_order: 6
    active: true
    tables:
      - crud_action
      - action_trigger_action
  Settings:
    menu_icon: cog
    menu_order: 7
    active: true
    tables:
      - store_setting
```

## STORE_SETTING
```yaml
table: store_setting
comment: Store Settings
tooltip: Core settings for the store backend
columns:
  store_setting_id: { type: integer, pk: true, autoincrement: true, comment: "Store Setting ID" }
  store_name:       { type: varchar, len: 200, nullable: false, comment: "Store Name", form_display: true, table_display: true, form_size: 6, order: 1 }
  store_slug:       { type: varchar, len: 100, nullable: false, unique: true, comment: "Store Slug", form_display: true, table_display: true, form_size: 6, order: 2 }
  store_email:      { type: varchar, len: 200, comment: "Store Email", form_display: true, table_display: true, form_size: 6, order: 3 }
  support_email:    { type: varchar, len: 200, comment: "Support Email", form_display: true, table_display: true, form_size: 6, order: 4 }
  currency:         { type: varchar, len: 10, default: "USD", comment: "Currency", form_display: true, table_display: true, form_size: 4, order: 5 }
  tax_rate:         { type: decimal, len: 8, scale: 2, default: 0, comment: "Default Tax Rate", form_display: true, table_display: true, form_size: 4, order: 6 }
  free_shipping_threshold: { type: decimal, len: 12, scale: 2, default: 0, comment: "Free Shipping Threshold", form_display: true, table_display: true, form_size: 4, order: 7 }
  checkout_enabled: { type: boolean, default: true, comment: "Checkout Enabled", form_display: true, table_display: true, form_size: 4, order: 8 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 9 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
data:
  - {store_setting_id: 1, store_name: Sample Store, store_slug: sample-store, store_email: store@example.com, support_email: support@example.com, currency: USD, tax_rate: 0.20, free_shipping_threshold: 100.00, checkout_enabled: true, active: true, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
table_layout:
  default_order: [{field: store_setting_id, order: ASC}]
```

## CUSTOMER
```yaml
table: customer
comment: Customer
tooltip: Customers placing orders in the online store
columns:
  customer_id:      { type: integer, pk: true, autoincrement: true, comment: "Customer ID" }
  first_name:       { type: varchar, len: 100, nullable: false, comment: "First Name", form_display: true, table_display: true, form_size: 4, order: 1 }
  last_name:        { type: varchar, len: 100, nullable: false, comment: "Last Name", form_display: true, table_display: true, form_size: 4, order: 2 }
  email:            { type: varchar, len: 200, nullable: false, unique: true, comment: "Email", form_display: true, table_display: true, form_size: 4, order: 3 }
  phone:            { type: varchar, len: 50, comment: "Phone", form_display: true, table_display: true, form_size: 4, order: 4 }
  company:          { type: varchar, len: 200, comment: "Company", form_display: true, table_display: true, form_size: 4, order: 5 }
  status:           { type: varchar, len: 20, default: "Active", comment: "Status", form_display: true, table_display: true, form_size: 4, order: 6 }
  notes:            { type: text, comment: "Notes", form_display: true, table_display: true, form_long_text: true, order: 7 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 8 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
  allow_in_subform: {customer_address: true, order_header: true}
table_layout:
  default_order: [{field: created_at, order: DESC}]
```

## CUSTOMER_ADDRESS
```yaml
table: customer_address
comment: Customer Address
tooltip: Shipping and billing addresses for a customer
columns:
  customer_address_id: { type: integer, pk: true, autoincrement: true, comment: "Customer Address ID" }
  customer_id:         { type: integer, fk: "customer.customer_id", nullable: false, comment: "Customer", form_display: true, table_display: true, form_size: 4, order: 1 }
  address_type:        { type: varchar, len: 20, default: "Shipping", comment: "Address Type", form_display: true, table_display: true, form_size: 4, order: 2 }
  address_line1:       { type: varchar, len: 200, comment: "Address Line 1", form_display: true, table_display: true, form_size: 6, order: 3 }
  address_line2:       { type: varchar, len: 200, comment: "Address Line 2", form_display: true, table_display: true, form_size: 6, order: 4 }
  city:                { type: varchar, len: 100, comment: "City", form_display: true, table_display: true, form_size: 4, order: 5 }
  state:               { type: varchar, len: 100, comment: "State", form_display: true, table_display: true, form_size: 4, order: 6 }
  postal_code:         { type: varchar, len: 20, comment: "Postal Code", form_display: true, table_display: true, form_size: 4, order: 7 }
  country:             { type: varchar, len: 100, comment: "Country", form_display: true, table_display: true, form_size: 4, order: 8 }
  active:              { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 9 }
  user_id:             { type: integer, comment: "User ID" }
  app_id:              { type: integer, comment: "App ID" }
  created_at:          { type: datetime, comment: "Created at" }
  updated_at:          { type: datetime, comment: "Updated at" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 8
table_layout:
  default_order: [{field: customer_address_id, order: DESC}]
```

## PRODUCT_CATEGORY
```yaml
table: product_category
comment: Product Category
tooltip: Catalog grouping for products
columns:
  product_category_id: { type: integer, pk: true, autoincrement: true, comment: "Product Category ID" }
  category:           { type: varchar, len: 100, nullable: false, unique: true, comment: "Category", form_display: true, table_display: true, form_size: 8, order: 1 }
  category_desc:      { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 2 }
  active:             { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 4, order: 3 }
  user_id:            { type: integer, comment: "User ID" }
  app_id:             { type: integer, comment: "App ID" }
  created_at:         { type: datetime, comment: "Created at" }
  updated_at:         { type: datetime, comment: "Updated at" }
  excluded:           { type: boolean, default: false, comment: "Excluded" }
data:
  - {product_category_id: 1, category: Electronics, category_desc: Gadgets and devices, active: true, excluded: false}
  - {product_category_id: 2, category: Home, category_desc: Home and lifestyle items, active: true, excluded: false}
  - {product_category_id: 3, category: Fashion, category_desc: Accessories and apparel, active: true, excluded: false}
form_layout:
  size: 6
table_layout:
  default_order: [{field: category, order: ASC}]
```

## PRODUCT
```yaml
table: product
comment: Product
tooltip: Items sold through the online store
columns:
  product_id:           { type: integer, pk: true, autoincrement: true, comment: "Product ID" }
  sku:                  { type: varchar, len: 100, nullable: false, unique: true, comment: "SKU", form_display: true, table_display: true, form_size: 4, order: 1 }
  product_name:         { type: varchar, len: 200, nullable: false, comment: "Product Name", form_display: true, table_display: true, form_size: 8, order: 2 }
  product_desc:         { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, order: 3 }
  product_category_id:  { type: integer, fk: "product_category.product_category_id", comment: "Category", form_display: true, table_display: true, form_size: 4, order: 4 }
  price:                { type: decimal, len: 12, scale: 2, default: 0, comment: "Price", form_display: true, table_display: true, form_size: 4, order: 5 }
  sale_price:           { type: decimal, len: 12, scale: 2, default: 0, comment: "Sale Price", form_display: true, table_display: true, form_size: 4, order: 6 }
  cost_price:           { type: decimal, len: 12, scale: 2, default: 0, comment: "Cost Price", form_display: true, table_display: true, form_size: 4, order: 7 }
  stock_quantity:       { type: decimal, len: 12, scale: 2, default: 0, comment: "Stock Quantity", form_display: true, table_display: true, form_size: 4, order: 8 }
  track_inventory:      { type: boolean, default: true, comment: "Track Inventory", form_display: true, table_display: true, form_size: 4, order: 9 }
  active:               { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 10 }
  user_id:              { type: integer, comment: "User ID" }
  app_id:               { type: integer, comment: "App ID" }
  created_at:           { type: datetime, comment: "Created at" }
  updated_at:           { type: datetime, comment: "Updated at" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
data:
  - {product_id: 1, sku: LAPTOP-001, product_name: UltraBook Laptop, product_desc: Lightweight laptop for remote work, product_category_id: 1, price: 999.99, sale_price: 899.99, cost_price: 650.00, stock_quantity: 12, track_inventory: true, active: true, excluded: false}
  - {product_id: 2, sku: MOUSE-001, product_name: Wireless Mouse, product_desc: Ergonomic wireless mouse, product_category_id: 1, price: 29.99, sale_price: 24.99, cost_price: 12.00, stock_quantity: 50, track_inventory: true, active: true, excluded: false}
  - {product_id: 3, sku: LAMP-001, product_name: Desk Lamp, product_desc: Minimal desk lamp, product_category_id: 2, price: 49.99, sale_price: 39.99, cost_price: 20.00, stock_quantity: 20, track_inventory: true, active: true, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
  allow_in_subform: {order_item: true}
table_layout:
  default_order: [{field: product_name, order: ASC}]
```

## ORDER_STATUS
```yaml
table: order_status
comment: Order Status
tooltip: Lifecycle status for store orders
columns:
  order_status_id: { type: integer, pk: true, autoincrement: true, comment: "Order Status ID" }
  order_status:    { type: varchar, len: 50, nullable: false, unique: true, comment: "Status", form_display: true, table_display: true, form_size: 8, order: 1 }
  order_status_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 2 }
  active:          { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 4, order: 3 }
  user_id:         { type: integer, comment: "User ID" }
  app_id:          { type: integer, comment: "App ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
data:
  - {order_status_id: 1, order_status: Pending, order_status_desc: Order received, active: true, excluded: false}
  - {order_status_id: 2, order_status: Processing, order_status_desc: Order being prepared, active: true, excluded: false}
  - {order_status_id: 3, order_status: Shipped, order_status_desc: Order sent to carrier, active: true, excluded: false}
  - {order_status_id: 4, order_status: Completed, order_status_desc: Order delivered, active: true, excluded: false}
  - {order_status_id: 5, order_status: Canceled, order_status_desc: Order canceled, active: true, excluded: false}
form_layout:
  size: 6
table_layout:
  default_order: [{field: order_status_id, order: ASC}]
```

## ORDER_HEADER
```yaml
table: order_header
comment: Order Header
tooltip: Main order record capturing checkout and totals
columns:
  order_header_id:  { type: integer, pk: true, autoincrement: true, comment: "Order ID" }
  order_no:         { type: varchar, len: 100, nullable: false, unique: true, comment: "Order Number", form_display: true, table_display: true, form_size: 4, order: 1 }
  customer_id:      { type: integer, fk: "customer.customer_id", nullable: false, comment: "Customer", form_display: true, table_display: true, form_size: 4, order: 2 }
  order_status_id:  { type: integer, fk: "order_status.order_status_id", default: 1, comment: "Status", form_display: true, table_display: true, form_size: 4, order: 3 }
  order_date:       { type: datetime, nullable: false, comment: "Order Date", form_display: true, table_display: true, form_size: 4, order: 4 }
  currency:         { type: varchar, len: 10, default: "USD", comment: "Currency", form_display: true, table_display: true, form_size: 4, order: 5 }
  subtotal:         { type: decimal, len: 12, scale: 2, default: 0, comment: "Subtotal", form_display: true, table_display: true, form_size: 4, order: 6 }
  discount_amount:  { type: decimal, len: 12, scale: 2, default: 0, comment: "Discount", form_display: true, table_display: true, form_size: 4, order: 7 }
  tax_amount:       { type: decimal, len: 12, scale: 2, default: 0, comment: "Tax", form_display: true, table_display: true, form_size: 4, order: 8 }
  shipping_amount:  { type: decimal, len: 12, scale: 2, default: 0, comment: "Shipping", form_display: true, table_display: true, form_size: 4, order: 9 }
  total_amount:     { type: decimal, len: 12, scale: 2, default: 0, comment: "Total", form_display: true, table_display: true, form_size: 4, order: 10 }
  paid_amount:      { type: decimal, len: 12, scale: 2, default: 0, comment: "Paid", form_display: true, table_display: true, form_size: 4, order: 11 }
  balance_due:      { type: decimal, len: 12, scale: 2, default: 0, comment: "Balance Due", form_display: true, table_display: true, form_size: 4, order: 12 }
  notes:            { type: text, comment: "Notes", form_display: true, table_display: true, form_long_text: true, order: 13 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
  allow_in_subform: {order_item: true, shipment: true, order_payment: true}
table_layout:
  default_order: [{field: order_date, order: DESC}]
```

## ORDER_ITEM
```yaml
table: order_item
comment: Order Item
tooltip: Line items inside an order
columns:
  order_item_id:   { type: integer, pk: true, autoincrement: true, comment: "Order Item ID" }
  order_header_id: { type: integer, fk: "order_header.order_header_id", nullable: false, comment: "Order", form_display: true, table_display: true, form_size: 4, order: 1 }
  product_id:      { type: integer, fk: "product.product_id", nullable: false, comment: "Product", form_display: true, table_display: true, form_size: 4, order: 2 }
  qty:             { type: decimal, len: 12, scale: 2, default: 1, comment: "Quantity", form_display: true, table_display: true, form_size: 4, order: 3 }
  unit_price:      { type: decimal, len: 12, scale: 2, default: 0, comment: "Unit Price", form_display: true, table_display: true, form_size: 4, order: 4 }
  line_total:      { type: decimal, len: 12, scale: 2, default: 0, comment: "Line Total", form_display: true, table_display: true, form_size: 4, order: 5 }
  user_id:         { type: integer, comment: "User ID" }
  app_id:          { type: integer, comment: "App ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 8
table_layout:
  default_order: [{field: order_item_id, order: DESC}]
```

## SHIPMENT
```yaml
table: shipment
comment: Shipment
tooltip: Shipping details for an order
columns:
  shipment_id:      { type: integer, pk: true, autoincrement: true, comment: "Shipment ID" }
  order_header_id:  { type: integer, fk: "order_header.order_header_id", nullable: false, comment: "Order", form_display: true, table_display: true, form_size: 4, order: 1 }
  carrier:          { type: varchar, len: 100, comment: "Carrier", form_display: true, table_display: true, form_size: 4, order: 2 }
  tracking_no:      { type: varchar, len: 100, comment: "Tracking Number", form_display: true, table_display: true, form_size: 4, order: 3 }
  shipment_status:  { type: varchar, len: 20, default: "Pending", comment: "Status", form_display: true, table_display: true, form_size: 4, order: 4 }
  shipped_at:       { type: datetime, comment: "Shipped At", form_display: true, table_display: true, form_size: 6, order: 5 }
  delivered_at:     { type: datetime, comment: "Delivered At", form_display: true, table_display: true, form_size: 6, order: 6 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 8
table_layout:
  default_order: [{field: created_at, order: DESC}]
```

## PAYMENT_METHOD
```yaml
table: payment_method
comment: Payment Method
tooltip: Payment options available to shoppers
columns:
  payment_method_id:   { type: integer, pk: true, autoincrement: true, comment: "Payment Method ID" }
  payment_method:      { type: varchar, len: 100, nullable: false, unique: true, comment: "Method", form_display: true, table_display: true, form_size: 8, order: 1 }
  payment_method_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 2 }
  active:              { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 4, order: 3 }
  user_id:             { type: integer, comment: "User ID" }
  app_id:              { type: integer, comment: "App ID" }
  created_at:          { type: datetime, comment: "Created at" }
  updated_at:          { type: datetime, comment: "Updated at" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
data:
  - {payment_method_id: 1, payment_method: Card, payment_method_desc: Credit or debit card, active: true, excluded: false}
  - {payment_method_id: 2, payment_method: Bank Transfer, payment_method_desc: Bank transfer, active: true, excluded: false}
  - {payment_method_id: 3, payment_method: Cash, payment_method_desc: Cash on delivery, active: true, excluded: false}
form_layout:
  size: 6
table_layout:
  default_order: [{field: payment_method, order: ASC}]
```

## ORDER_PAYMENT
```yaml
table: order_payment
comment: Order Payment
tooltip: Payments received for an order
columns:
  order_payment_id:   { type: integer, pk: true, autoincrement: true, comment: "Order Payment ID" }
  order_header_id:    { type: integer, fk: "order_header.order_header_id", nullable: false, comment: "Order", form_display: true, table_display: true, form_size: 4, order: 1 }
  payment_method_id:  { type: integer, fk: "payment_method.payment_method_id", nullable: false, comment: "Method", form_display: true, table_display: true, form_size: 4, order: 2 }
  payment_ref:        { type: varchar, len: 100, comment: "Reference", form_display: true, table_display: true, form_size: 4, order: 3 }
  amount:             { type: decimal, len: 12, scale: 2, default: 0, comment: "Amount", form_display: true, table_display: true, form_size: 4, order: 4 }
  payment_status:     { type: varchar, len: 20, default: "Pending", comment: "Status", form_display: true, table_display: true, form_size: 4, order: 5 }
  paid_at:            { type: datetime, comment: "Paid At", form_display: true, table_display: true, form_size: 6, order: 6 }
  notes:              { type: text, comment: "Notes", form_display: true, table_display: true, form_long_text: true, order: 7 }
  user_id:            { type: integer, comment: "User ID" }
  app_id:             { type: integer, comment: "App ID" }
  created_at:         { type: datetime, comment: "Created at" }
  updated_at:         { type: datetime, comment: "Updated at" }
  excluded:           { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 8
table_layout:
  default_order: [{field: created_at, order: DESC}]
```

## PROMO_CODE
```yaml
table: promo_code
comment: Promotion Code
tooltip: Discount promotions for the store
columns:
  promo_code_id:   { type: integer, pk: true, autoincrement: true, comment: "Promo Code ID" }
  promo_code:      { type: varchar, len: 100, nullable: false, unique: true, comment: "Code", form_display: true, table_display: true, form_size: 4, order: 1 }
  promo_desc:      { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, order: 2 }
  promo_type:      { type: varchar, len: 20, default: "Percent", comment: "Type", form_display: true, table_display: true, form_size: 4, order: 3 }
  promo_value:     { type: decimal, len: 12, scale: 2, default: 0, comment: "Value", form_display: true, table_display: true, form_size: 4, order: 4 }
  valid_from:      { type: datetime, comment: "Valid From", form_display: true, table_display: true, form_size: 6, order: 5 }
  valid_to:        { type: datetime, comment: "Valid To", form_display: true, table_display: true, form_size: 6, order: 6 }
  active:          { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 7 }
  user_id:         { type: integer, comment: "User ID" }
  app_id:          { type: integer, comment: "App ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
data:
  - {promo_code_id: 1, promo_code: WELCOME10, promo_desc: 10 percent off first order, promo_type: Percent, promo_value: 10.00, valid_from: "2025-01-01 00:00:00", valid_to: "2030-12-31 23:59:59", active: true, excluded: false}
  - {promo_code_id: 2, promo_code: FREESHIP, promo_desc: Free shipping above threshold, promo_type: Fixed, promo_value: 5.00, valid_from: "2025-01-01 00:00:00", valid_to: "2030-12-31 23:59:59", active: true, excluded: false}
form_layout:
  size: 8
table_layout:
  default_order: [{field: promo_code, order: ASC}]
```

## DASHBOARD
```yaml
table: dashboard
comment: Dashboards
columns:
  dashboard_id:   { type: integer, pk: true, autoincrement: true, comment: "Dashboard ID" }
  dashboard:      { type: varchar, len: 200, comment: "Dashboard", form_display: true, table_display: true, form_size: 8, order: 1 }
  dashboard_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 4 }
  dashboard_conf: { type: text, nullable: false, comment: "Conf / Params", form_display: true, form_long_text: true, form_code: markdown, order: 5 }
  order:          { type: integer, comment: "Order", form_display: true, table_display: true, form_size: 2, order: 2 }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 3 }
  user_id:        { type: integer, comment: "User ID" }
  app_id:         { type: integer, comment: "App ID" }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
data:
  - {dashboard_id: 1, dashboard: Store Overview, dashboard_desc: Store KPIs and revenue overview, dashboard_conf: "# Store Overview\n- Orders today\n- Revenue this month\n- Low stock alerts", order: 1, active: true, excluded: false}
  - {dashboard_id: 2, dashboard: Order Funnel, dashboard_desc: Pending, processing and shipped orders, dashboard_conf: "# Order Funnel\n- Pending orders\n- Processing orders\n- Shipped orders", order: 2, active: true, excluded: false}
  - {dashboard_id: 3, dashboard: Inventory Alerts, dashboard_desc: Products with critical stock, dashboard_conf: "# Inventory Alerts\n- Stock below reorder threshold\n- Fast-moving products", order: 3, active: true, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
table_layout:
  default_order: [{field: order, order: ASC}]
table_extra_options:
  - { component: EvidenceDash, label: dashboard, intercept_r: true, size: 12 }
```

# STORE_DATA
```yaml
name: STORE_DATA
description: Seed data and CRUD action rules for the store model
database: STORE
runs_as: MODEL_DATA
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## STORE_DASHBOARD_OVERVIEW
```yaml
table: dashboard
description: Add default store overview dashboard
cond: 'WHERE dashboard_id = :dashboard_id AND excluded = false'
data:
  dashboard_id: 1
  dashboard: Store Overview
  dashboard_desc: Store KPIs and revenue overview
  dashboard_conf: |
    # Store Overview
    - Orders today
    - Revenue this month
    - Low stock alerts
  order: 1
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## STORE_DASHBOARD_FUNNEL
```yaml
table: dashboard
description: Add default order funnel dashboard
cond: 'WHERE dashboard_id = :dashboard_id AND excluded = false'
data:
  dashboard_id: 2
  dashboard: Order Funnel
  dashboard_desc: Pending, processing and shipped orders
  dashboard_conf: |
    # Order Funnel
    - Pending orders
    - Processing orders
    - Shipped orders
  order: 2
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## STORE_DASHBOARD_INVENTORY
```yaml
table: dashboard
description: Add default inventory alerts dashboard
cond: 'WHERE dashboard_id = :dashboard_id AND excluded = false'
data:
  dashboard_id: 3
  dashboard: Inventory Alerts
  dashboard_desc: Products with critical stock
  dashboard_conf: |
    # Inventory Alerts
    - Stock below reorder threshold
    - Fast-moving products
  order: 3
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## STORE_TRIGGER_ORDER_EMAIL
```yaml
table: crud_action
description: Send an email when a store order is created
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: STORE_TRIGGER_ORDER_EMAIL
  crud_action: Order Created
  action_type_id: 2
  err_msg: 'Error sending the email to {{.email}}!'
  table: order_header
  db: STORE
  active: true
  create: true
  update: false
  delete: false
  parallel: true
  email_subject: Your Order Was Created
  email_to: '{{.customer_data.email}}'
  email_template: |
    <p>Hi {{.customer_data.first_name}},</p>
    <p>Your order has been created.</p>
    <p><strong>Customer:</strong> {{.customer_data.first_name}}</p>
    <p>If you did not request this, contact support.</p>
  children:
    table: action_data
    cond: 'WHERE action_data = :action_data'
    data:
      action_data: customer_data
      action_data_desc: Customer data lookup for order email
      action_data_type_id: 3
      crud_action_id: crud_action_id()
      odata_path: "STORE/customer?$filter=customer_id eq {{.customer_id}}"
      sigle_row_obj: true
      active: true
      user_id: 1
      app_id: appId()
      created_at: Now()
      updated_at: Now()
      excluded: false
```

## STORE_TRIGGER_ORDER_STATUS_EMAIL
```yaml
table: crud_action
description: Send an email when a store order status changes
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: STORE_TRIGGER_ORDER_STATUS_EMAIL
  crud_action: Order Status Updated
  action_type_id: 2
  err_msg: 'Error sending the status email to {{.email}}!'
  table: order_header
  db: STORE
  active: true
  create: false
  update: true
  delete: false
  parallel: true
  email_subject: Your Order Status Changed
  email_to: '{{.customer_data.email}}'
  email_template: |
    <p>Hi {{.customer_data.first_name}},</p>
    <p>Your order status has been updated.</p>
    <p>If you need help, contact support.</p>
  children:
    table: action_data
    cond: 'WHERE action_data = :action_data'
    data:
      action_data: customer_data
      action_data_desc: Customer data lookup for order status email
      action_data_type_id: 3
      crud_action_id: crud_action_id()
      odata_path: "STORE/customer?$filter=customer_id eq {{.customer_id}}"
      sigle_row_obj: true
      active: true
      user_id: 1
      app_id: appId()
      created_at: Now()
      updated_at: Now()
      excluded: false
```

## STORE_TRIGGER_INVENTORY_ALERT
```yaml
table: crud_action
description: Trigger an internal API call when stock becomes low
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: STORE_TRIGGER_INVENTORY_ALERT
  crud_action: Low Stock Alert
  action_type_id: 3
  err_msg: 'Error sending the inventory alert for product {{.product_id}}!'
  table: product
  db: STORE
  active: true
  create: false
  update: true
  delete: false
  parallel: true
  api: '/api/store/inventory/alert/{{.product_id}}'
```
