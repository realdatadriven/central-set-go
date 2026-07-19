<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD012 -->
<!-- markdownlint-disable MD047 -->
# DYNUI_MODEL

This model stores a template-driven website definition in the database. It is
intended to be used together with a business model such as `STORE`: DYNUI owns
the site, routes, renderable sources and assets; STORE owns products and orders.

```yaml
name: DYNUI
description: Database-backed template websites, pages, partials, routes, assets and OData data providers
runs_as: MODEL
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
create_all: checkfirst
_drop_all: checkfirst
update_table_metadata: true
active: true
cs_app:
  Dynamic UI:
    menu_icon: paint-brush
    menu_order: 8
    active: true
    tables:
      - dynui
      - dynui_route
      - dynui_page
      - dynui_partial
      - dynui_page_data
      - dynui_navigation
      - dynui_asset
      - {table: dynui_setting, active: false}
```

## DYNUI
```yaml
table: dynui
comment: Dynamic Website
tooltip: A website served below /dynui/{ui_slug}
columns:
  dynui_id:          { type: integer, pk: true, autoincrement: true, comment: "Dynamic UI ID" }
  ui_slug:           { type: varchar, len: 100, unique: true, nullable: false, comment: "URL Slug", form_display: true, table_display: true, form_size: 4, order: 1, form_regex_val: "^[a-z][a-z0-9-]*$", form_val_msg: "Use lowercase letters, numbers and hyphens; begin with a letter." }
  ui_name:           { type: varchar, len: 200, nullable: false, comment: "Website Name", form_display: true, table_display: true, form_size: 5, order: 2 }
  ui_desc:           { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, order: 3 }
  default_locale:    { type: varchar, len: 10, default: "en", comment: "Default Locale", form_display: true, table_display: true, form_size: 3, order: 4 }
  default_page_id:   { type: integer, comment: "Default Page ID", form_display: true, table_display: true, form_size: 3, order: 5 }
  active:            { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 6 }
  user_id:           { type: integer, comment: "User ID" }
  app_id:            { type: integer, comment: "App ID" }
  created_at:        { type: datetime, comment: "Created At" }
  updated_at:        { type: datetime, comment: "Updated At" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  allow_in_subform: {dynui_route: true, dynui_page: true, dynui_partial: true, dynui_asset: true, dynui_navigation: true, dynui_setting: true}
  tabs_steps_conf:
    - {label: Website, fields: [ui_slug, ui_name, ui_desc, default_locale, default_page_id, active]}
table_layout:
  default_order: [{field: ui_slug, order: ASC}]
```

## DYNUI_PARTIAL
```yaml
table: dynui_partial
comment: Dynamic UI Partial
tooltip: Reusable Go HTML templates. Every active partial for the selected UI is parsed with the page template before rendering.
columns:
  dynui_partial_id:   { type: integer, pk: true, autoincrement: true, comment: "Partial ID" }
  dynui_id:           { type: integer, fk: "dynui.dynui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 3, order: 1 }
  dynui_partial:      { type: varchar, len: 150, nullable: false, comment: "Partial Name", form_display: true, table_display: true, form_size: 4, order: 2, form_regex_val: "^[a-z][a-z0-9_-]*$", form_val_msg: "Use lowercase letters, numbers, underscores and hyphens." }
  dynui_partial_desc: { type: text, comment: "Partial Description", form_display: true, table_display: true, form_long_text: true, order: 3 }
  partial_template:   { type: text, nullable: false, comment: "Partial Template", form_display: true, form_long_text: true, form_code: html, order: 4 }
  active:             { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 5 }
  user_id:            { type: integer, comment: "User ID" }
  app_id:             { type: integer, comment: "App ID" }
  created_at:         { type: datetime, comment: "Created At" }
  updated_at:         { type: datetime, comment: "Updated At" }
  excluded:           { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  tabs_steps_conf:
    - {label: Partial, fields: [dynui_id, dynui_partial, dynui_partial_desc, active]}
    - {label: Template, fields: [partial_template]}
table_layout:
  default_order: [{field: dynui_partial, order: ASC}]
```

## DYNUI_PAGE
```yaml
table: dynui_page
comment: Dynamic UI Page
tooltip: A page template parsed with the active partials for its website.
columns:
  dynui_page_id:       { type: integer, pk: true, autoincrement: true, comment: "Page ID" }
  dynui_id:            { type: integer, fk: "dynui.dynui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 3, order: 1 }
  page_key:            { type: varchar, len: 150, nullable: false, comment: "Page Key", form_display: true, table_display: true, form_size: 3, order: 2, form_regex_val: "^[a-z][a-z0-9_-]*$", form_val_msg: "Use lowercase letters, numbers, underscores and hyphens." }
  page_title:          { type: varchar, len: 255, nullable: false, comment: "Title", form_display: true, table_display: true, form_size: 6, order: 3 }
  meta_description:    { type: text, comment: "SEO Description", form_display: true, form_long_text: true, order: 4 }
  page_template:       { type: text, nullable: false, comment: "Page Template", form_display: true, form_long_text: true, form_code: html, order: 5 }
  cache_seconds:       { type: integer, default: 0, comment: "Public Cache Seconds", form_display: true, table_display: true, form_size: 3, order: 6 }
  active:              { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 7 }
  user_id:             { type: integer, comment: "User ID" }
  app_id:              { type: integer, comment: "App ID" }
  created_at:          { type: datetime, comment: "Created At" }
  updated_at:          { type: datetime, comment: "Updated At" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  allow_in_subform: {dynui_page_data: true}
  tabs_steps_conf:
    - {label: Page, fields: [dynui_id, page_key, page_title, meta_description, cache_seconds, active]}
    - {label: Template, fields: [page_template]}
table_layout:
  default_order: [{field: page_key, order: ASC}]
```

## DYNUI_ROUTE
```yaml
table: dynui_route
comment: Dynamic UI Route
tooltip: Relative route matched after /dynui/{ui_slug}; use :name for a path value, for example /products/:product_id.
columns:
  dynui_route_id: { type: integer, pk: true, autoincrement: true, comment: "Route ID" }
  dynui_id:       { type: integer, fk: "dynui.dynui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 3, order: 1 }
  route_key:      { type: varchar, len: 150, unique: true, nullable: false, comment: "Route Key", form_display: true, table_display: true, form_size: 3, order: 2 }
  route_path:     { type: varchar, len: 500, nullable: false, comment: "Relative Path", form_display: true, table_display: true, form_size: 5, order: 3 }
  http_method:    { type: varchar, len: 10, default: "GET", nullable: false, comment: "HTTP Method", form_display: true, table_display: true, form_size: 2, order: 4 }
  dynui_page_id:  { type: integer, fk: "dynui_page.dynui_page_id", nullable: false, comment: "Page", form_display: true, table_display: true, form_size: 3, order: 5 }
  requires_auth:  { type: boolean, default: false, comment: "Requires Authentication", form_display: true, table_display: true, form_size: 2, order: 6 }
  sort_order:     { type: integer, default: 100, comment: "Match Order", form_display: true, table_display: true, form_size: 2, order: 7 }
  active:         { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 8 }
  user_id:        { type: integer, comment: "User ID" }
  app_id:         { type: integer, comment: "App ID" }
  created_at:     { type: datetime, comment: "Created At" }
  updated_at:     { type: datetime, comment: "Updated At" }
  excluded:       { type: boolean, default: false, comment: "Excluded" }
form_layout: {size: 9, form_in_popup: false}
table_layout:
  default_order: [{field: sort_order, order: ASC}]
```

## DYNUI_PAGE_DATA
```yaml
table: dynui_page_data
comment: Dynamic UI Page Data
tooltip: Named data pulled from an OData path and made available to its page template.
columns:
  dynui_page_data_id:   { type: integer, pk: true, autoincrement: true, comment: "Page Data ID" }
  dynui_page_data:      { type: varchar, len: 100, nullable: false, comment: "Data Name", form_display: true, table_display: true, form_size: 4, order: 2, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not begin with a number; use letters, numbers and underscores." }
  dynui_page_data_desc: { type: text, comment: "Data Description", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 4 }
  dynui_page_id:        { type: integer, fk: "dynui_page.dynui_page_id", nullable: false, comment: "Page", form_display: true, table_display: true, form_size: 4, order: 1 }
  odata_path:           { type: text, nullable: false, comment: "OData Path", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 5 }
  sigle_row_obj:        { type: boolean, default: false, comment: "Single Row Object", form_display: true, table_display: true, form_size: 3, order: 6 }
  active:               { type: boolean, default: true, comment: "Active", table_display: true, form_display: true, form_size: 2, order: 3 }
  user_id:              { type: integer, comment: "User ID" }
  app_id:               { type: integer, comment: "App ID" }
  created_at:           { type: datetime, comment: "Created At" }
  updated_at:           { type: datetime, comment: "Updated At" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 9
table_layout:
  default_order: [{field: dynui_page_data, order: ASC}]
```

## DYNUI_NAVIGATION
```yaml
table: dynui_navigation
comment: Dynamic UI Navigation Item
columns:
  dynui_navigation_id: { type: integer, pk: true, autoincrement: true, comment: "Navigation ID" }
  dynui_id:            { type: integer, fk: "dynui.dynui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 3, order: 1 }
  parent_id:           { type: integer, fk: "dynui_navigation.dynui_navigation_id", comment: "Parent Item", form_display: true, table_display: true, form_size: 3, order: 2 }
  label:               { type: varchar, len: 150, nullable: false, comment: "Label", form_display: true, table_display: true, form_size: 3, order: 3 }
  href:                { type: varchar, len: 500, nullable: false, comment: "Link", form_display: true, table_display: true, form_size: 5, order: 4 }
  sort_order:          { type: integer, default: 100, comment: "Sort Order", form_display: true, table_display: true, form_size: 2, order: 5 }
  active:              { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 6 }
  user_id:             { type: integer, comment: "User ID" }
  app_id:              { type: integer, comment: "App ID" }
  created_at:          { type: datetime, comment: "Created At" }
  updated_at:          { type: datetime, comment: "Updated At" }
  excluded:            { type: boolean, default: false, comment: "Excluded" }
form_layout: {size: 8}
table_layout:
  default_order: [{field: sort_order, order: ASC}]
```

## DYNUI_ASSET
```yaml
table: dynui_asset
comment: Dynamic UI Asset
tooltip: Store UTF-8 assets directly or binary assets as base64 text. The content is database-backed, not a filesystem path.
columns:
  dynui_asset_id:  { type: integer, pk: true, autoincrement: true, comment: "Asset ID" }
  dynui_id:        { type: integer, fk: "dynui.dynui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 3, order: 1 }
  asset_path:      { type: varchar, len: 500, nullable: false, comment: "Relative Asset Path", form_display: true, table_display: true, form_size: 5, order: 2 }
  mime_type:       { type: varchar, len: 150, nullable: false, comment: "MIME Type", form_display: true, table_display: true, form_size: 3, order: 3 }
  content_encoding: { type: varchar, len: 20, default: "utf-8", nullable: false, comment: "Encoding", form_display: true, table_display: true, form_size: 2, order: 4 }
  asset_content:   { type: text, nullable: false, comment: "Content", form_display: true, form_long_text: true, form_code: text, order: 5 }
  checksum:        { type: varchar, len: 128, comment: "Content Checksum", form_display: true, table_display: true, form_size: 3, order: 6 }
  cache_seconds:   { type: integer, default: 86400, comment: "Cache Seconds", form_display: true, table_display: true, form_size: 2, order: 7 }
  active:          { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 8 }
  user_id:         { type: integer, comment: "User ID" }
  app_id:          { type: integer, comment: "App ID" }
  created_at:      { type: datetime, comment: "Created At" }
  updated_at:      { type: datetime, comment: "Updated At" }
  excluded:        { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  size: 9
  tabs_steps_conf:
    - {label: Asset, fields: [dynui_id, asset_path, mime_type, content_encoding, checksum, cache_seconds, active]}
    - {label: Content, fields: [asset_content]}
table_layout:
  default_order: [{field: asset_path, order: ASC}]
```

## DYNUI_SETTING
```yaml
table: dynui_setting
comment: Dynamic UI Setting
tooltip: Site-wide data exposed to every template as .Site.Settings.
columns:
  dynui_setting_id: { type: integer, pk: true, autoincrement: true, comment: "Setting ID" }
  dynui_id:         { type: integer, fk: "dynui.dynui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 3, order: 1 }
  setting_key:      { type: varchar, len: 150, nullable: false, comment: "Key", form_display: true, table_display: true, form_size: 3, order: 2 }
  setting_value:    { type: text, comment: "Value", form_display: true, form_long_text: true, form_code: json, order: 3 }
  is_public:        { type: boolean, default: true, comment: "Expose to Public Templates", form_display: true, table_display: true, form_size: 2, order: 4 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 5 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created At" }
  updated_at:       { type: datetime, comment: "Updated At" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout: {size: 8}
table_layout:
  default_order: [{field: setting_key, order: ASC}]
```

# Renderer contract

The later `GET /dynui/{ui_slug}/{path...}` handler should build one render
context and parse the page template plus every active `dynui_partial` in a
single `html/template.Template` set. Recommended context keys are:

```text
.Site       # dynui row, public dynui_setting values and navigation tree
.Page       # selected dynui_page row
.Route      # route row, PathParams and Query
.Data       # each dynui_page_data result under its dynui_page_data name
.User       # authenticated public-safe identity, only when applicable
.RequestID  # safe tracing value
```

A partial defines a named template and a page calls it by name, for example:

```gotemplate
{{/* dynui_partial = header */}}
{{ define "header" }}<header><a href="{{ .Site.BaseURL }}">{{ .Site.Name }}</a></header>{{ end }}
```

```gotemplate
{{ template "header" . }}
{{ range .Data.products }}<article>{{ .product_name }}</article>{{ end }}
```

For the `products` page-data row, `odata_path` can be
`STORE/product?$filter=active eq true&$orderby=product_name asc`. The handler
should render this stored Go template only after it has safely interpolated
route values, such as `{{ .Route.PathParams.product_id }}`, and should never
accept an OData path or template name directly from the browser.

# DYNUI_DATA
```yaml
name: DYNUI_DATA
description: Starter storefront UI backed by the STORE model
database: DYNUI
runs_as: MODEL_DATA
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## DYNUI_STORE
```yaml
table: dynui
description: Add the sample store website
cond: 'WHERE ui_slug = :ui_slug AND excluded = false'
data:
  dynui_id: 1
  ui_slug: store
  ui_name: Sample Store
  ui_desc: Starter template-driven storefront using STORE as its backend
  default_locale: en
  default_page_id: 1
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_PARTIAL_HEADER
```yaml
table: dynui_partial
description: Add the store header partial
cond: 'WHERE dynui_id = :dynui_id AND dynui_partial = :dynui_partial AND excluded = false'
data:
  dynui_partial_id: 1
  dynui_id: 1
  dynui_partial: header
  dynui_partial_desc: Site header and primary navigation
  partial_template: FileContent(dynui/store/partials/header.html)
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_PARTIAL_FOOTER
```yaml
table: dynui_partial
description: Add the store footer partial
cond: 'WHERE dynui_id = :dynui_id AND dynui_partial = :dynui_partial AND excluded = false'
data:
  dynui_partial_id: 2
  dynui_id: 1
  dynui_partial: footer
  dynui_partial_desc: Site footer
  partial_template: FileContent(dynui/store/partials/footer.html)
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_PAGE_HOME
```yaml
table: dynui_page
description: Add the store home page
cond: 'WHERE dynui_id = :dynui_id AND page_key = :page_key AND excluded = false'
data:
  dynui_page_id: 1
  dynui_id: 1
  page_key: home
  page_title: Sample Store
  meta_description: Welcome to the sample store
  page_template: FileContent(dynui/store/pages/home.html)
  cache_seconds: 60
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_PAGE_PRODUCTS
```yaml
table: dynui_page
description: Add the store products page
cond: 'WHERE dynui_id = :dynui_id AND page_key = :page_key AND excluded = false'
data:
  dynui_page_id: 2
  dynui_id: 1
  page_key: products
  page_title: Products
  meta_description: Browse available products
  page_template: FileContent(dynui/store/pages/products.html)
  cache_seconds: 60
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_PAGE_PRODUCT
```yaml
table: dynui_page
description: Add the store product detail page
cond: 'WHERE dynui_id = :dynui_id AND page_key = :page_key AND excluded = false'
data:
  dynui_page_id: 3
  dynui_id: 1
  page_key: product
  page_title: Product
  meta_description: Product details
  page_template: FileContent(dynui/store/pages/product.html)
  cache_seconds: 60
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_ROUTE_HOME
```yaml
table: dynui_route
description: Add the store home route
cond: 'WHERE route_key = :route_key AND excluded = false'
data:
  dynui_route_id: 1
  dynui_id: 1
  route_key: store-home
  route_path: /
  http_method: GET
  dynui_page_id: 1
  requires_auth: false
  sort_order: 1
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_ROUTE_PRODUCTS
```yaml
table: dynui_route
description: Add the store product-list route
cond: 'WHERE route_key = :route_key AND excluded = false'
data:
  dynui_route_id: 2
  dynui_id: 1
  route_key: store-products
  route_path: /products
  http_method: GET
  dynui_page_id: 2
  requires_auth: false
  sort_order: 2
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_ROUTE_PRODUCT
```yaml
table: dynui_route
description: Add the store product-detail route
cond: 'WHERE route_key = :route_key AND excluded = false'
data:
  dynui_route_id: 3
  dynui_id: 1
  route_key: store-product
  route_path: /products/:product_id
  http_method: GET
  dynui_page_id: 3
  requires_auth: false
  sort_order: 3
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_PAGE_DATA_FEATURED_PRODUCTS
```yaml
table: dynui_page_data
description: Load featured products for the home page from STORE
cond: 'WHERE dynui_page_id = :dynui_page_id AND dynui_page_data = :dynui_page_data AND excluded = false'
data:
  dynui_page_data_id: 1
  dynui_page_data: products
  dynui_page_data_desc: Active products shown on the home page
  dynui_page_id: 1
  odata_path: 'STORE/product?$filter=active eq true&$top=6&$orderby=product_name asc'
  sigle_row_obj: false
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_PAGE_DATA_PRODUCTS
```yaml
table: dynui_page_data
description: Load products for the product-list page from STORE
cond: 'WHERE dynui_page_id = :dynui_page_id AND dynui_page_data = :dynui_page_data AND excluded = false'
data:
  dynui_page_data_id: 2
  dynui_page_data: products
  dynui_page_data_desc: Active products shown on the products page
  dynui_page_id: 2
  odata_path: 'STORE/product?$filter=active eq true&$top=48&$orderby=product_name asc'
  sigle_row_obj: false
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_PAGE_DATA_PRODUCT
```yaml
table: dynui_page_data
description: Load one product using the product_id path parameter
cond: 'WHERE dynui_page_id = :dynui_page_id AND dynui_page_data = :dynui_page_data AND excluded = false'
data:
  dynui_page_data_id: 3
  dynui_page_data: product
  dynui_page_data_desc: Product selected from the route
  dynui_page_id: 3
  odata_path: 'STORE/product?$filter=product_id eq {{.Route.PathParams.product_id}}'
  sigle_row_obj: true
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## DYNUI_NAVIGATION_STORE
```yaml
table: dynui_navigation
description: Add the store navigation links
cond: 'WHERE dynui_id = :dynui_id AND label = :label AND excluded = false'
data:
  dynui_navigation_id: 1
  dynui_id: 1
  label: Home
  href: /dynui/store
  sort_order: 1
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
children:
  table: dynui_navigation
  cond: 'WHERE dynui_id = :dynui_id AND label = :label AND excluded = false'
  data:
    dynui_navigation_id: 2
    dynui_id: 1
    label: Products
    href: /dynui/store/products
    sort_order: 2
    active: true
    user_id: 1
    app_id: appId()
    created_at: Now()
    updated_at: Now()
    excluded: false
```

## DYNUI_ASSET_STORE_CSS
```yaml
table: dynui_asset
description: Add the store stylesheet as a database asset
cond: 'WHERE dynui_id = :dynui_id AND asset_path = :asset_path AND excluded = false'
data:
  dynui_asset_id: 1
  dynui_id: 1
  asset_path: assets/store.css
  mime_type: text/css; charset=utf-8
  content_encoding: utf-8
  asset_content: FileContent(dynui/store/assets/store.css)
  cache_seconds: 86400
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```
