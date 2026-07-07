<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD012 -->
<!-- markdownlint-disable MD047 -->
# ERP_MODEL
```yaml
name: ERP
description: Basic ERP model for a small business with customers, suppliers, inventory, sales, purchases, payments and accounting
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
  Inventory:
    menu_icon: boxes-2
    menu_order: 2
    active: true
    tables:
      - item_type
      - item
      - inventory_movement
  Sales:
    menu_icon: shopping-cart
    menu_order: 3
    active: true
    tables:
      - invoice
      - invoice_line
  Purchases:
    menu_icon: clipboard-list
    menu_order: 4
    active: true
    tables:
      - invoice
      - invoice_line
  Journals:
    menu_icon: book
    menu_order: 5
    active: true
    tables:
      - journal_entry
      - journal_entry_line
  Payments:
    menu_icon: credit-card
    menu_order: 6
    active: true
    tables:
      - payment
      - payment_method
      - account
  Parties:
    menu_icon: users
    menu_order: 7
    active: true
    tables:
      - party_type
      - party
```

## PARTY_TYPE
```yaml
table: party_type
comment: Party Type
tooltip: Types of business partners such as customers and suppliers
columns:
  party_type_id:   { type: integer, pk: true, autoincrement: true, comment: "Party Type ID" }
  party_type:      { type: varchar, len: 100, nullable: false, unique: true, comment: "Type", form_display: true, table_display: true, form_size: 6, order: 1 }
  party_type_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 2 }
  active:          { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 3 }
  user_id:         { type: integer, comment: "User ID" }
  app_id:          { type: integer, comment: "App ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
data:
  - {party_type_id: 1, party_type: Customer, party_type_desc: Customer / Client, active: true, excluded: false}
  - {party_type_id: 2, party_type: Supplier, party_type_desc: Supplier / Vendor, active: true, excluded: false}
  - {party_type_id: 3, party_type: Internal, party_type_desc: Internal / Staff, active: true, excluded: false}
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 6
table_layout:
  default_order: [{field: party_type_id, order: ASC}]
```

## PARTY
```yaml
table: party
comment: Party / Entity
tooltip: Customers, suppliers and other business entities
columns:
  party_id:          { type: integer, pk: true, autoincrement: true, comment: "Party ID" }
  party:             { type: varchar, len: 200, nullable: false, comment: "Name", form_display: true, table_display: true, form_size: 8, order: 1 }
  party_desc:        { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, order: 4 }
  party_type_id:     { type: integer, fk: "party_type.party_type_id", comment: "Type", form_display: true, table_display: true, form_size: 4, order: 2 }
  tax_id:            { type: varchar, len: 100, comment: "Tax / Registration Number", form_display: true, table_display: true, form_size: 4, order: 3 }
  email:             { type: varchar, len: 200, comment: "Email", form_display: true, table_display: true, form_size: 6, order: 5 }
  phone:             { type: varchar, len: 50, comment: "Phone", form_display: true, table_display: true, form_size: 6, order: 6 }
  address:           { type: text, comment: "Address", form_display: true, table_display: true, form_long_text: true, order: 7 }
  payment_terms_days: { type: integer, default: 30, comment: "Payment Terms (Days)", form_display: true, table_display: true, form_size: 3, order: 8 }
  active:            { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 9 }
  user_id:           { type: integer, comment: "User ID" }
  app_id:            { type: integer, comment: "App ID" }
  created_at:        { type: datetime, comment: "Created at" }
  updated_at:        { type: datetime, comment: "Updated at" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
table_layout:
  default_order: [{field: party, order: ASC}]
```

## ITEM_TYPE
```yaml
table: item_type
comment: Item Type
tooltip: Goods or services sold or bought by the business
columns:
  item_type_id:   { type: integer, pk: true, autoincrement: true, comment: "Item Type ID" }
  item_type:      { type: varchar, len: 100, nullable: false, unique: true, comment: "Type", form_display: true, table_display: true, form_size: 6, order: 1 }
  item_type_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 2 }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 3 }
  user_id:        { type: integer, comment: "User ID" }
  app_id:         { type: integer, comment: "App ID" }
  created_at:     { type: datetime, comment: "Created at" }
  updated_at:     { type: datetime, comment: "Updated at" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
data:
  - {item_type_id: 1, item_type: Goods, item_type_desc: Physical product with inventory tracking, active: true, excluded: false}
  - {item_type_id: 2, item_type: Service, item_type_desc: Service with no stock movement, active: true, excluded: false}
form_layout:
  size: 6
```

## ITEM
```yaml
table: item
comment: Item / Product / Service
tooltip: Goods and services offered by the business
columns:
  item_id:            { type: integer, pk: true, autoincrement: true, comment: "Item ID" }
  item:               { type: varchar, len: 200, nullable: false, comment: "Name", form_display: true, table_display: true, form_size: 8, order: 1 }
  item_desc:          { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, order: 4 }
  item_type_id:       { type: integer, fk: "item_type.item_type_id", comment: "Type", form_display: true, table_display: true, form_size: 4, order: 2 }
  sku:                { type: varchar, len: 100, comment: "SKU / Code", form_display: true, table_display: true, form_size: 4, order: 3 }
  sell_price:         { type: decimal, len: 12, scale: 2, default: 0, comment: "Sell Price", form_display: true, table_display: true, form_size: 4, order: 5 }
  cost_price:         { type: decimal, len: 12, scale: 2, default: 0, comment: "Cost Price", form_display: true, table_display: true, form_size: 4, order: 6 }
  stock_on_hand:      { type: decimal, len: 12, scale: 2, default: 0, comment: "Stock On Hand", form_display: true, table_display: true, form_size: 4, order: 7 }
  track_inventory:    { type: boolean, default: false, comment: "Track Inventory", form_display: true, table_display: true, form_size: 4, order: 8 }
  active:             { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 9 }
  user_id:            { type: integer, comment: "User ID" }
  app_id:             { type: integer, comment: "App ID" }
  created_at:         { type: datetime, comment: "Created at" }
  updated_at:         { type: datetime, comment: "Updated at" }
  excluded:           { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 8
table_layout:
  default_order: [{field: item, order: ASC}]
```

## INVENTORY_MOVEMENT
```yaml
table: inventory_movement
comment: Inventory Movement
tooltip: Stock movement from purchases, sales, adjustments and losses
columns:
  inventory_movement_id: { type: integer, pk: true, autoincrement: true, comment: "Movement ID" }
  item_id:               { type: integer, fk: "item.item_id", nullable: false, comment: "Item", form_display: true, table_display: true, form_size: 4, order: 1 }
  movement_type:         { type: varchar, len: 20, nullable: false, comment: "Type", form_display: true, table_display: true, form_size: 4, order: 2 }
  qty:                   { type: decimal, len: 12, scale: 2, nullable: false, comment: "Quantity", form_display: true, table_display: true, form_size: 4, order: 3 }
  unit_cost:             { type: decimal, len: 12, scale: 2, default: 0, comment: "Unit Cost", form_display: true, table_display: true, form_size: 4, order: 4 }
  reference_table:       { type: varchar, len: 100, comment: "Reference Table", form_display: true, table_display: true, form_size: 6, order: 5 }
  reference_id:          { type: integer, comment: "Reference ID", form_display: true, table_display: true, form_size: 6, order: 6 }
  movement_date:         { type: datetime, comment: "Date", form_display: true, table_display: true, form_size: 6, order: 7 }
  user_id:               { type: integer, comment: "User ID" }
  app_id:                { type: integer, comment: "App ID" }
  created_at:            { type: datetime, comment: "Created at" }
  updated_at:            { type: datetime, comment: "Updated at" }
  excluded:              { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 8
table_layout:
  default_order: [{field: movement_date, order: DESC}]
```

## INVOICE
```yaml
table: invoice
comment: Invoice / Document
tooltip: Sales and purchase invoices
columns:
  invoice_id:         { type: integer, pk: true, autoincrement: true, comment: "Invoice ID" }
  invoice_no:         { type: varchar, len: 100, nullable: false, unique: true, comment: "Document Number", form_display: true, table_display: true, form_size: 4, order: 1 }
  invoice_type:       { type: varchar, len: 20, nullable: false, comment: "Sale / Purchase", form_display: true, table_display: true, form_size: 4, order: 2 }
  invoice_date:       { type: datetime, nullable: false, comment: "Date", form_display: true, table_display: true, form_size: 4, order: 3 }
  party_id:           { type: integer, fk: "party.party_id", nullable: false, comment: "Customer / Supplier", form_display: true, table_display: true, form_size: 6, order: 4 }
  reference_no:       { type: varchar, len: 100, comment: "Reference", form_display: true, table_display: true, form_size: 6, order: 5 }
  status:             { type: varchar, len: 20, default: "Draft", comment: "Status", form_display: true, table_display: true, form_size: 3, order: 6 }
  subtotal:           { type: decimal, len: 12, scale: 2, default: 0, comment: "Subtotal", form_display: true, table_display: true, form_size: 3, order: 7 }
  discount_amount:    { type: decimal, len: 12, scale: 2, default: 0, comment: "Discount", form_display: true, table_display: true, form_size: 3, order: 8 }
  tax_amount:         { type: decimal, len: 12, scale: 2, default: 0, comment: "Tax", form_display: true, table_display: true, form_size: 3, order: 9 }
  total:              { type: decimal, len: 12, scale: 2, default: 0, comment: "Total", form_display: true, table_display: true, form_size: 3, order: 10 }
  paid_amount:        { type: decimal, len: 12, scale: 2, default: 0, comment: "Paid", form_display: true, table_display: true, form_size: 3, order: 11 }
  balance_due:        { type: decimal, len: 12, scale: 2, default: 0, comment: "Balance Due", form_display: true, table_display: true, form_size: 3, order: 12 }
  notes:              { type: text, comment: "Notes", form_display: true, table_display: true, form_long_text: true, order: 13 }
  user_id:            { type: integer, comment: "User ID" }
  app_id:             { type: integer, comment: "App ID" }
  created_at:         { type: datetime, comment: "Created at" }
  updated_at:         { type: datetime, comment: "Updated at" }
  excluded:           { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 10
  allow_in_subform: {invoice_line: true}
table_layout:
  default_order: [{field: invoice_date, order: DESC}]
```

## INVOICE_LINE
```yaml
table: invoice_line
comment: Invoice Line
tooltip: Lines inside sales and purchase documents
columns:
  invoice_line_id: { type: integer, pk: true, autoincrement: true, comment: "Invoice Line ID" }
  invoice_id:      { type: integer, fk: "invoice.invoice_id", nullable: false, comment: "Invoice", form_display: true, table_display: true, form_size: 4, order: 1 }
  line_no:         { type: integer, default: 1, comment: "Line Number", form_display: true, table_display: true, form_size: 2, order: 2 }
  item_id:         { type: integer, fk: "item.item_id", nullable: false, comment: "Item", form_display: true, table_display: true, form_size: 4, order: 3 }
  description:     { type: varchar, len: 200, comment: "Description", form_display: true, table_display: true, form_size: 6, order: 4 }
  qty:             { type: decimal, len: 12, scale: 2, default: 1, comment: "Quantity", form_display: true, table_display: true, form_size: 3, order: 5 }
  unit_price:      { type: decimal, len: 12, scale: 2, default: 0, comment: "Unit Price", form_display: true, table_display: true, form_size: 3, order: 6 }
  line_total:      { type: decimal, len: 12, scale: 2, default: 0, comment: "Line Total", form_display: true, table_display: true, form_size: 3, order: 7 }
  cost_amount:     { type: decimal, len: 12, scale: 2, default: 0, comment: "Cost Amount", form_display: true, table_display: true, form_size: 3, order: 8 }
  user_id:         { type: integer, comment: "User ID" }
  app_id:          { type: integer, comment: "App ID" }
  created_at:      { type: datetime, comment: "Created at" }
  updated_at:      { type: datetime, comment: "Updated at" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 8
table_layout:
  default_order: [{field: line_no, order: ASC}]
```

## PAYMENT_METHOD
```yaml
table: payment_method
comment: Payment Method
tooltip: Cash, bank transfer, card or other payment methods
columns:
  payment_method_id:   { type: integer, pk: true, autoincrement: true, comment: "Payment Method ID" }
  payment_method:      { type: varchar, len: 100, nullable: false, unique: true, comment: "Method", form_display: true, table_display: true, form_size: 6, order: 1 }
  payment_method_desc: { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 2 }
  active:              { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 3 }
  user_id:             { type: integer, comment: "User ID" }
  app_id:              { type: integer, comment: "App ID" }
  created_at:          { type: datetime, comment: "Created at" }
  updated_at:          { type: datetime, comment: "Updated at" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
data:
  - {payment_method_id: 1, payment_method: Cash, payment_method_desc: Local cash, active: true, excluded: false}
  - {payment_method_id: 2, payment_method: Bank, payment_method_desc: Bank account transfer, active: true, excluded: false}
  - {payment_method_id: 3, payment_method: Card, payment_method_desc: Card payment, active: true, excluded: false}
form_layout:
  size: 6
```

## ACCOUNT
```yaml
table: account
comment: Accounting Account
tooltip: Cash, bank, receivables, payables, revenue and expense accounts
columns:
  account_id:        { type: integer, pk: true, autoincrement: true, comment: "Account ID" }
  account_code:      { type: varchar, len: 50, nullable: false, unique: true, comment: "Code", form_display: true, table_display: true, form_size: 4, order: 1 }
  account:           { type: varchar, len: 200, nullable: false, comment: "Account Name", form_display: true, table_display: true, form_size: 8, order: 2 }
  account_type:      { type: varchar, len: 50, nullable: false, comment: "Account Type", form_display: true, table_display: true, form_size: 4, order: 3 }
  payment_method_id: { type: integer, fk: "payment_method.payment_method_id", comment: "Payment Method", form_display: true, table_display: true, form_size: 4, order: 4 }
  opening_balance:   { type: decimal, len: 12, scale: 2, default: 0, comment: "Opening Balance", form_display: true, table_display: true, form_size: 4, order: 5 }
  current_balance:   { type: decimal, len: 12, scale: 2, default: 0, comment: "Current Balance", form_display: true, table_display: true, form_size: 4, order: 6 }
  active:            { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 3, order: 7 }
  user_id:           { type: integer, comment: "User ID" }
  app_id:            { type: integer, comment: "App ID" }
  created_at:        { type: datetime, comment: "Created at" }
  updated_at:        { type: datetime, comment: "Updated at" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
data:
  - {account_code: CASH, account: Local Cash, account_type: Asset, payment_method_id: 1, opening_balance: 0, current_balance: 0, active: true, excluded: false}
  - {account_code: BANK, account: Bank Account, account_type: Asset, payment_method_id: 2, opening_balance: 0, current_balance: 0, active: true, excluded: false}
  - {account_code: AR, account: Accounts Receivable, account_type: Asset, payment_method_id: 2, opening_balance: 0, current_balance: 0, active: true, excluded: false}
  - {account_code: AP, account: Accounts Payable, account_type: Liability, payment_method_id: 2, opening_balance: 0, current_balance: 0, active: true, excluded: false}
  - {account_code: REV, account: Sales Revenue, account_type: Revenue, payment_method_id: 2, opening_balance: 0, current_balance: 0, active: true, excluded: false}
  - {account_code: COGS, account: Cost of Goods Sold, account_type: Expense, payment_method_id: 2, opening_balance: 0, current_balance: 0, active: true, excluded: false}
form_layout:
  size: 8
table_layout:
  default_order: [{field: account_code, order: ASC}]
```

## JOURNAL_ENTRY
```yaml
table: journal_entry
comment: Journal Entry
tooltip: Double-entry accounting document for invoices and payments
columns:
  journal_entry_id: { type: integer, pk: true, autoincrement: true, comment: "Journal Entry ID" }
  entry_no:         { type: varchar, len: 100, nullable: false, unique: true, comment: "Entry Number", form_display: true, table_display: true, form_size: 4, order: 1 }
  entry_date:       { type: datetime, nullable: false, comment: "Date", form_display: true, table_display: true, form_size: 4, order: 2 }
  entry_type:       { type: varchar, len: 50, nullable: false, comment: "Type", form_display: true, table_display: true, form_size: 4, order: 3 }
  description:      { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 4 }
  source_table:     { type: varchar, len: 100, comment: "Source Table", form_display: true, table_display: true, form_size: 6, order: 5 }
  source_id:        { type: integer, comment: "Source ID", form_display: true, table_display: true, form_size: 6, order: 6 }
  status:           { type: varchar, len: 20, default: "Posted", comment: "Status", form_display: true, table_display: true, form_size: 3, order: 7 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created at" }
  updated_at:       { type: datetime, comment: "Updated at" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 8
  allow_in_subform: {journal_entry_line: true}
table_layout:
  default_order: [{field: entry_date, order: DESC}]
```

## JOURNAL_ENTRY_LINE
```yaml
table: journal_entry_line
comment: Journal Entry Line
tooltip: Debit and credit lines of a journal entry
columns:
  journal_entry_line_id: { type: integer, pk: true, autoincrement: true, comment: "Journal Entry Line ID" }
  journal_entry_id:      { type: integer, fk: "journal_entry.journal_entry_id", nullable: false, comment: "Journal Entry", form_display: true, table_display: true, form_size: 4, order: 1 }
  account_id:            { type: integer, fk: "account.account_id", nullable: false, comment: "Account", form_display: true, table_display: true, form_size: 4, order: 2 }
  debit:                 { type: decimal, len: 12, scale: 2, default: 0, comment: "Debit", form_display: true, table_display: true, form_size: 3, order: 3 }
  credit:                { type: decimal, len: 12, scale: 2, default: 0, comment: "Credit", form_display: true, table_display: true, form_size: 3, order: 4 }
  description:           { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, form_size: 12, order: 5 }
  user_id:               { type: integer, comment: "User ID" }
  app_id:                { type: integer, comment: "App ID" }
  created_at:            { type: datetime, comment: "Created at" }
  updated_at:            { type: datetime, comment: "Updated at" }
  excluded:              { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 8
table_layout:
  default_order: [{field: journal_entry_line_id, order: DESC}]
```

## PAYMENT
```yaml
table: payment
comment: Payment
tooltip: Money movement for invoices and cash/bank operations
columns:
  payment_id:         { type: integer, pk: true, autoincrement: true, comment: "Payment ID" }
  payment_no:         { type: varchar, len: 100, nullable: false, unique: true, comment: "Payment Number", form_display: true, table_display: true, form_size: 4, order: 1 }
  payment_date:       { type: datetime, nullable: false, comment: "Date", form_display: true, table_display: true, form_size: 4, order: 2 }
  party_id:           { type: integer, fk: "party.party_id", nullable: false, comment: "Client / Supplier", form_display: true, table_display: true, form_size: 4, order: 3 }
  invoice_id:         { type: integer, fk: "invoice.invoice_id", comment: "Invoice", form_display: true, table_display: true, form_size: 4, order: 4 }
  account_id:         { type: integer, fk: "account.account_id", nullable: false, comment: "Cash / Bank Account", form_display: true, table_display: true, form_size: 4, order: 5 }
  payment_method_id:  { type: integer, fk: "payment_method.payment_method_id", nullable: false, comment: "Method", form_display: true, table_display: true, form_size: 4, order: 6 }
  amount:             { type: decimal, len: 12, scale: 2, default: 0, comment: "Amount", form_display: true, table_display: true, form_size: 3, order: 7 }
  direction:          { type: varchar, len: 10, nullable: false, comment: "IN / OUT", form_display: true, table_display: true, form_size: 3, order: 8 }
  status:             { type: varchar, len: 20, default: "Posted", comment: "Status", form_display: true, table_display: true, form_size: 3, order: 9 }
  notes:              { type: text, comment: "Notes", form_display: true, table_display: true, form_long_text: true, order: 10 }
  user_id:            { type: integer, comment: "User ID" }
  app_id:             { type: integer, comment: "App ID" }
  created_at:         { type: datetime, comment: "Created at" }
  updated_at:         { type: datetime, comment: "Updated at" }
  excluded:           { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 8
table_layout:
  default_order: [{field: payment_date, order: DESC}]
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
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
table_layout:
  default_order: [{field: order, order: ASC}]
table_extra_options:
  - { component: EvidenceDash, label: dashboard, intercept_r: true, size: 12 }
```

# DATA
```yaml
name: ERP_DATA
description: Seed data and CRUD action rules for the ERP model
database: ADMIN
runs_as: MODEL_DATA
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## ERP_CRUD_RECALC_TOTAL
```yaml
table: crud_action
description: Recalculate invoice totals when invoice lines are created, updated or deleted
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: ERP_RECALC_INVOICE_TOTAL
  crud_action: Recalculate Invoice Total
  action_type_id: 1
  err_msg: 'Error recalculating invoice totals for invoice {{.invoice_id}}'
  table: invoice_line
  db: ADMIN
  active: true
  create: true
  update: true
  delete: true
  parallel: false
  sql: |
    UPDATE invoice
    SET subtotal = COALESCE((SELECT SUM(line_total) FROM invoice_line WHERE invoice_id = :invoice_id), 0),
        total = COALESCE((SELECT SUM(line_total) FROM invoice_line WHERE invoice_id = :invoice_id), 0),
        balance_due = COALESCE((SELECT SUM(line_total) FROM invoice_line WHERE invoice_id = :invoice_id), 0) - COALESCE(paid_amount, 0),
        updated_at = NOW()
    WHERE invoice_id = :invoice_id
```

## ERP_CRUD_STOCK_OUT_SALE
```yaml
table: crud_action
description: Decrease inventory when a sale invoice line is created
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: ERP_STOCK_OUT_ON_SALE
  crud_action: Decrease stock on sale
  action_type_id: 1
  err_msg: 'Error updating inventory for invoice line {{.invoice_line_id}}'
  table: invoice_line
  db: ADMIN
  active: true
  create: true
  update: false
  delete: false
  parallel: false
  sql: |
    INSERT INTO inventory_movement (item_id, movement_type, qty, unit_cost, reference_table, reference_id, movement_date, created_at, updated_at, excluded)
    VALUES (:item_id, 'OUT', :qty, :unit_price, 'invoice_line', :invoice_line_id, NOW(), NOW(), NOW(), FALSE)
```

## ERP_CRUD_STOCK_IN_PURCHASE
```yaml
table: crud_action
description: Increase inventory when a purchase invoice line is created
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: ERP_STOCK_IN_ON_PURCHASE
  crud_action: Increase stock on purchase
  action_type_id: 1
  err_msg: 'Error updating inventory for purchase invoice line {{.invoice_line_id}}'
  table: invoice_line
  db: ADMIN
  active: true
  create: true
  update: false
  delete: false
  parallel: false
  sql: |
    INSERT INTO inventory_movement (item_id, movement_type, qty, unit_cost, reference_table, reference_id, movement_date, created_at, updated_at, excluded)
    VALUES (:item_id, 'IN', :qty, :unit_price, 'invoice_line', :invoice_line_id, NOW(), NOW(), NOW(), FALSE)
```

## ERP_CRUD_POST_SALES_ACCOUNTING
```yaml
table: crud_action
description: Create accounting entries when a sales invoice is created
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: ERP_POST_SALES_ACCOUNTING
  crud_action: Post sales accounting
  action_type_id: 1
  err_msg: 'Error posting accounting for invoice {{.invoice_id}}'
  table: invoice
  db: ADMIN
  active: true
  create: true
  update: false
  delete: false
  parallel: false
  sql: |
    INSERT INTO journal_entry (entry_no, entry_date, entry_type, description, source_table, source_id, status, created_at, updated_at, excluded)
    VALUES ('INV-' || :invoice_id, NOW(), 'SALE', 'Sales invoice posting', 'invoice', :invoice_id, 'Posted', NOW(), NOW(), FALSE)
```

## ERP_CRUD_POST_PURCHASES_ACCOUNTING
```yaml
table: crud_action
description: Create accounting entries when a purchase invoice is created
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: ERP_POST_PURCHASES_ACCOUNTING
  crud_action: Post purchases accounting
  action_type_id: 1
  err_msg: 'Error posting accounting for purchase invoice {{.invoice_id}}'
  table: invoice
  db: ADMIN
  active: true
  create: true
  update: false
  delete: false
  parallel: false
  sql: |
    INSERT INTO journal_entry (entry_no, entry_date, entry_type, description, source_table, source_id, status, created_at, updated_at, excluded)
    VALUES ('PUR-' || :invoice_id, NOW(), 'PURCHASE', 'Purchase invoice posting', 'invoice', :invoice_id, 'Posted', NOW(), NOW(), FALSE)
```

## ERP_CRUD_POST_PAYMENT
```yaml
table: crud_action
description: Record money movement when a payment is created
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: ERP_POST_PAYMENT
  crud_action: Record payment movement
  action_type_id: 1
  err_msg: 'Error recording payment {{.payment_id}}'
  table: payment
  db: ADMIN
  active: true
  create: true
  update: false
  delete: false
  parallel: false
  sql: |
    INSERT INTO journal_entry (entry_no, entry_date, entry_type, description, source_table, source_id, status, created_at, updated_at, excluded)
    VALUES ('PAY-' || :payment_id, NOW(), 'PAYMENT', 'Payment movement', 'payment', :payment_id, 'Posted', NOW(), NOW(), FALSE)
```

## ERP_CRUD_ACTION_CHAIN
```yaml
table: crud_action
description: Example of wiring invoice creation to inventory and accounting rules
cond: 'WHERE crud_action_code = :crud_action_code'
data:
  crud_action_code: ERP_INVOICE_ACTION_CHAIN
  crud_action: Invoice workflow chain
  action_type_id: 6
  err_msg: 'Error chaining invoice actions for {{.invoice_id}}'
  table: invoice
  db: ADMIN
  active: true
  create: true
  update: false
  delete: false
  parallel: false
  etlx_md_template: |
    1. Recalculate invoice totals
    2. Create inventory movement for sale or purchase
    3. Create accounting journal entries
    4. Record cash or bank payment movement
  children:
    table: action_trigger_action
    cond: 'WHERE action_trigger_action = :action_trigger_action'
    data:
      crud_action_id: crud_action_id()
      action_trigger_action: ERP_TRIGGER_RECALC_TOTAL
      action_trigger_action_desc: Trigger invoice total recalculation after invoice line changes
      action_trigger_code: ERP_RECALC_INVOICE_TOTAL
      trigger_order: 1
      active: true
```

This model is intentionally simple and practical for a small business. It gives you a solid starting point for:
- parties and entities such as customers and suppliers
- products and services
- sales and purchase invoices
- inventory updates when items are sold or bought
- cash and bank movement through payments
- accounting entries with debit and credit logic via journal entries

A good next step is to add a small set of default accounts and then enable the CRUD actions progressively for sales, purchases and payments.
