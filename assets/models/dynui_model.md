<!-- markdownlint-disable MD022 -->
<!-- markdownlint-disable MD025 -->
<!-- markdownlint-disable MD031 -->
<!-- markdownlint-disable MD012 -->
<!-- markdownlint-disable MD047 -->

# UI_MODEL
```yaml
name: UI
description: Database-backed template websites, pages, partials, assets and C7 OData to be avaliable for the templates
runs_as: MODEL
database: UI
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
create_all: checkfirst
update_table_metadata: true
active: true
cs_app:
  UI:
    menu_icon: paint-brush
    menu_order: 8
    active: true
    tables:
      - ui
      - {table: ui_partial, active: false}
      - {table: ui_partial_data, active: false}
      - {table: ui_page, active: false}
      - {table: ui_page_data, active: false}
      - {table: ui_asset, active: false}
```

## UI
```yaml
table: ui
comment: Dynamic Website
tooltip: A website served below /ui/{page}
columns:
  ui_id:            { type: integer, pk: true, autoincrement: true, comment: "UI ID" }
  ui_slug:          { type: varchar, len: 100, unique: true, nullable: false, comment: "URL Slug", form_display: true, table_display: true, form_size: 4, order: 1, form_regex_val: "^[a-z][a-z0-9-]*$", form_val_msg: "Use lowercase letters, numbers and hyphens; begin with a letter." }
  ui_name:          { type: varchar, len: 200, nullable: false, comment: "Website Name", form_display: true, table_display: true, form_size: 5, order: 2 }
  ui_desc:          { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, order: 3 }
  default_locale:   { type: varchar, len: 10, default: "en", comment: "Default Locale", form_display: true, table_display: true, form_size: 3, order: 4 }
#  database:         { type: varchar, len: 50, default: "en", comment: "Database", form_display: true, table_display: true, form_size: 3, order: 4 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 6 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created At" }
  updated_at:       { type: datetime, comment: "Updated At" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  allow_in_subform: {ui_route: true, ui_page: true, ui_partial: true, ui_asset: true}
  tabs_steps_conf:
    - {label: Website, fields: [ui_slug, ui_name, ui_desc, default_locale, active]}
table_layout:
  default_order: [{field: ui_id, order: DESC}]
```

## UI_PARTIAL
```yaml
table: ui_partial
comment: UI Partial
tooltip: Reusable Go HTML templates. Every active partial for the selected UI is parsed with the page template before rendering.
columns:
  ui_partial_id:     { type: integer, pk: true, autoincrement: true, comment: "Partial ID" }
  ui_partial:        { type: varchar, len: 150, nullable: false, comment: "Partial Name", form_display: true, table_display: true, form_size: 4, order: 2, form_regex_val: "^[a-z][a-z0-9_-]*$", form_val_msg: "Use lowercase letters, numbers, underscores and hyphens." }
  ui_partial_desc:   { type: text, comment: "Partial Description", form_display: true, table_display: true, form_long_text: true, order: 3 }
  partial_template:  { type: text, nullable: false, comment: "Partial Template", form_display: true, form_long_text: true, form_code: html, order: 4 }
  active:            { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 5 }
  ui_id:             { type: integer, fk: "ui.ui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 3, order: 1 }
  user_id:           { type: integer, comment: "User ID" }
  app_id:            { type: integer, comment: "App ID" }
  created_at:        { type: datetime, comment: "Created At" }
  updated_at:        { type: datetime, comment: "Updated At" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  allow_in_subform: {ui_partial_data: true}
  tabs_steps_conf:
    - {label: Partial, fields: [ui_id, ui_partial, ui_partial_desc, active]}
    - {label: Template, fields: [partial_template]}
table_layout:
  default_order: [{field: ui_partial_id, order: ASC}]
```

## UI_PARTIAL_DATA
```yaml
table: ui_partial_data
comment: UI Partial Data
tooltip: Named data pulled from C7 OData path and made available to its partial template.
columns:
  ui_partial_data_id:   { type: integer, pk: true, autoincrement: true, comment: "Partial Data ID" }
  ui_partial_data:      { type: varchar, len: 100, nullable: false, comment: "Data Name", form_display: true, table_display: true, form_size: 4, order: 2, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not begin with a number; use letters, numbers and underscores." }
  ui_partial_data_desc: { type: text, comment: "Data Description", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 4 }
  odata_path:           { type: text, nullable: false, comment: "OData Path", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 5 }
  sigle_row_obj:        { type: boolean, default: false, comment: "Single Row Object", form_display: true, table_display: true, form_size: 3, order: 6 }
  active:               { type: boolean, default: true, comment: "Active", table_display: true, form_display: true, form_size: 2, order: 3 }
  ui_partial_id:        { type: integer, fk: "ui_partial.ui_partial_id", nullable: false, comment: "Partial", form_display: true, table_display: true, form_size: 4, order: 1 }
  user_id:              { type: integer, comment: "User ID" }
  app_id:               { type: integer, comment: "App ID" }
  created_at:           { type: datetime, comment: "Created At" }
  updated_at:           { type: datetime, comment: "Updated At" }
  excluded:             { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 6
table_layout:
  default_order: [{field: ui_partial_data_id, order: ASC}]
```

## UI_PAGE
```yaml
table: ui_page
comment: UI Page
tooltip: A page template parsed with the active partials for its website.
columns:
  ui_page_id:       { type: integer, pk: true, autoincrement: true, comment: "Page ID" }
  page_key:         { type: varchar, len: 150, nullable: false, comment: "Page Key", form_display: true, table_display: true, form_size: 3, order: 2, form_regex_val: "^[a-z][a-z0-9_-]*$", form_val_msg: "Use lowercase letters, numbers, underscores and hyphens." }
  page_title:       { type: varchar, len: 255, nullable: false, comment: "Title", form_display: true, table_display: true, form_size: 6, order: 3 }
  meta_description: { type: text, comment: "SEO Description", form_display: true, form_long_text: true, order: 4 }
  page_template:    { type: text, nullable: false, comment: "Page Template", form_display: true, form_long_text: true, form_code: html, order: 5 }
  cache_seconds:    { type: integer, default: 0, comment: "Public Cache Seconds", form_display: true, table_display: true, form_size: 3, order: 6 }
  default_page:     { type: boolean, default: true, comment: "Default Page", form_display: true, table_display: true, form_size: 2, order: 7 }
  requires_login:   { type: boolean, default: false, comment: "Requires Login", form_display: true, table_display: true, form_size: 2, order: 7 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 8 }
  ui_id:            { type: integer, fk: "ui.ui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 3, order: 1 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created At" }
  updated_at:       { type: datetime, comment: "Updated At" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  allow_in_subform: {ui_page_data: true}
  tabs_steps_conf:
    - {label: Page, fields: [ui_id, page_key, page_title, meta_description, cache_seconds, default_page, requires_login, active]}
    - {label: Template, fields: [page_template]}
table_layout:
  default_order: [{field: ui_page_id, order: DESC}]
```

## UI_PAGE_DATA
```yaml
table: ui_page_data
comment: UI Page Data
tooltip: Named data pulled from C7 OData path and made available to its page template.
columns:
  ui_page_data_id:   { type: integer, pk: true, autoincrement: true, comment: "Page Data ID" }
  ui_page_data:      { type: varchar, len: 100, nullable: false, comment: "Data Name", form_display: true, table_display: true, form_size: 4, order: 2, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not begin with a number; use letters, numbers and underscores." }
  ui_page_data_desc: { type: text, comment: "Data Description", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 4 }
  odata_path:        { type: text, nullable: false, comment: "OData Path", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 5 }
  sigle_row_obj:     { type: boolean, default: false, comment: "Single Row Object", form_display: true, table_display: true, form_size: 3, order: 6 }
  active:            { type: boolean, default: true, comment: "Active", table_display: true, form_display: true, form_size: 2, order: 3 }
  ui_page_id:        { type: integer, fk: "ui_page.ui_page_id", nullable: false, comment: "Page", form_display: true, table_display: true, form_size: 4, order: 1 }
  user_id:           { type: integer, comment: "User ID" }
  app_id:            { type: integer, comment: "App ID" }
  created_at:        { type: datetime, comment: "Created At" }
  updated_at:        { type: datetime, comment: "Updated At" }
  excluded:          { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 6
table_layout:
  default_order: [{field: ui_page_data_id, order: ASC}]
```

## UI_ASSET
```yaml
table: ui_asset
comment: UI Asset
tooltip: Store UTF-8 assets directly or binary assets as base64 text. The content is database-backed, not a filesystem path.
columns:
  ui_asset_id:      { type: integer, pk: true, autoincrement: true, comment: "Asset ID" }
  ui_id:            { type: integer, fk: "ui.ui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 3, order: 1 }
  asset_path:       { type: varchar, len: 500, nullable: false, comment: "Relative Asset Path", form_display: true, table_display: true, form_size: 5, order: 2 }
  mime_type:        { type: varchar, len: 150, nullable: false, comment: "MIME Type", form_display: true, table_display: true, form_size: 3, order: 3 }
  content_encoding: { type: varchar, len: 20, default: "utf-8", nullable: false, comment: "Encoding", form_display: true, table_display: true, form_size: 2, order: 4 }
  asset_content:    { type: text, nullable: false, comment: "Content", form_display: true, form_long_text: true, form_code: text, order: 5 }
  checksum:         { type: varchar, len: 128, comment: "Content Checksum", form_display: true, table_display: true, form_size: 3, order: 6 }
  cache_seconds:    { type: integer, default: 86400, comment: "Cache Seconds", form_display: true, table_display: true, form_size: 2, order: 7 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 8 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created At" }
  updated_at:       { type: datetime, comment: "Updated At" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  size: 9
  tabs_steps_conf:
    - {label: Asset, fields: [ui_id, asset_path, mime_type, content_encoding, checksum, cache_seconds, active]}
    - {label: Content, fields: [asset_content]}
table_layout:
  default_order: [{field: asset_path, order: ASC}]
```

For the `products` page-data row, `odata_path` can be
`STORE/product?$filter=active eq true&$orderby=product_name asc`. The handler
should render this stored Go template only after it has safely interpolated
route values, such as `{{ .Route.PathParams.product_id }}`, and should never
accept C7 OData path or template name directly from the browser.

# UI_DATA
```yaml
name: UI_DATA
description: Minimal single-page example with header/footer partials
database: UI
runs_as: MODEL_DATA
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## UI_DEMO
```yaml
table: ui
description: Add the demo website
cond: 'WHERE ui_slug = :ui_slug AND excluded = false'
data:
  ui_id: 1
  ui_slug: demo
  ui_name: Demo Site
  ui_desc: Minimal example used to test one page with header/footer partials
  default_locale: en
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## UI_PARTIAL_HEADER
```yaml
table: ui_partial
description: Add the demo header partial
cond: 'WHERE ui_id = :ui_id AND ui_partial = :ui_partial AND excluded = false'
data:
  ui_partial_id: 1
  ui_id: 1
  ui_partial: header
  ui_partial_desc: Site header and primary navigation
  partial_template: |
    <header style="padding:1rem;background:#222;color:#fff;">
      <strong>{{.UI.ui_name}}</strong>
      <nav style="float:right;">
        <a href="/ui/demo" style="color:#fff;margin-left:1rem;">Home</a>
        <a href="/ui/demo/about" style="color:#fff;margin-left:1rem;">About</a>
      </nav>
    </header>
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## UI_PARTIAL_FOOTER
```yaml
table: ui_partial
description: Add the demo footer partial
cond: 'WHERE ui_id = :ui_id AND ui_partial = :ui_partial AND excluded = false'
data:
  ui_partial_id: 2
  ui_id: 1
  ui_partial: footer
  ui_partial_desc: Site footer
  partial_template: |
    <footer style="padding:1rem;background:#eee;margin-top:2rem;color:#555;">
      &copy; {{.UI.ui_name}} &middot; built with central-set-go
    </footer>
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## UI_PAGE_HOME
```yaml
table: ui_page
description: Add the demo home page (default page for the demo ui)
cond: 'WHERE ui_id = :ui_id AND page_key = :page_key AND excluded = false'
data:
  ui_page_id: 1
  ui_id: 1
  page_key: home
  page_title: Welcome
  meta_description: Minimal single-page demo
  page_template: |
    {{template "header" .}}
    <main style="padding:1rem;">
      <h1>{{.Page.page_title}}</h1>
      <p>{{.Page.meta_description}}</p>
      <p>This page is served straight out of the database — the header and
      footer above are separate <code>ui_partial</code> rows parsed together
      with this <code>ui_page.page_template</code>.</p>
    </main>
    {{template "footer" .}}
  cache_seconds: 0
  default_page: true
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

# UI_DATA
```yaml
name: UI_DATA
description: RealDataDriven consulting + SaaS landing page (DaisyUI + htmx)
database: UI
runs_as: MODEL_DATA
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
```

## UI_RDD
```yaml
table: ui
description: Add the RealDataDriven consulting website
cond: 'WHERE ui_slug = :ui_slug AND excluded = false'
data:
  ui_id: 1
  ui_slug: rdd
  ui_name: RealDataDriven
  ui_desc: Data engineering consulting - ETL/ELT/Reverse ETL, warehouses, lakehouses, and open-source-backed SaaS
  default_locale: en
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## UI_PARTIAL_HEADER
```yaml
table: ui_partial
description: Add the RealDataDriven header/navbar partial
cond: 'WHERE ui_id = :ui_id AND ui_partial = :ui_partial AND excluded = false'
data:
  ui_partial_id: 1
  ui_id: 1
  ui_partial: header
  ui_partial_desc: Sticky navbar with brand and section links
  partial_template: |
    <div class="navbar bg-base-100 shadow-sm sticky top-0 z-50 px-4 lg:px-12">
      <div class="flex-1">
        <a href="/ui/rdd" class="btn btn-ghost text-xl">
          <span class="text-primary font-bold">Real</span>DataDriven
        </a>
      </div>
      <div class="hidden md:flex flex-none gap-1">
        <a href="#services" class="btn btn-ghost btn-sm">Services</a>
        <a href="#opensource" class="btn btn-ghost btn-sm">Open Source</a>
        <a href="#saas" class="btn btn-ghost btn-sm">SaaS</a>
        <a href="#contact" class="btn btn-ghost btn-sm">Contact</a>
        <a href="https://github.com/realdatadriven" target="_blank" rel="noopener" class="btn btn-ghost btn-sm">GitHub</a>
        <a href="/ui/rdd/login" class="btn btn-primary btn-sm ml-2">Login</a>
      </div>
      <div class="flex-none md:hidden">
        <div class="dropdown dropdown-end">
          <div tabindex="0" role="button" class="btn btn-ghost">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" /></svg>
          </div>
          <ul tabindex="0" class="menu menu-sm dropdown-content bg-base-100 rounded-box shadow mt-3 w-44 p-2 z-50">
            <li><a href="#services">Services</a></li>
            <li><a href="#opensource">Open Source</a></li>
            <li><a href="#saas">SaaS</a></li>
            <li><a href="#contact">Contact</a></li>
            <li><a href="/ui/rdd/login">Login</a></li>
          </ul>
        </div>
      </div>
    </div>
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## UI_PARTIAL_FOOTER
```yaml
table: ui_partial
description: Add the RealDataDriven footer partial
cond: 'WHERE ui_id = :ui_id AND ui_partial = :ui_partial AND excluded = false'
data:
  ui_partial_id: 2
  ui_id: 1
  ui_partial: footer
  ui_partial_desc: Footer with open-source links and copyright
  partial_template: |
    <footer class="footer footer-center bg-base-200 text-base-content p-10 mt-12">
      <nav class="grid grid-flow-col gap-4">
        <a href="#services" class="link link-hover">Services</a>
        <a href="#opensource" class="link link-hover">Open Source</a>
        <a href="#saas" class="link link-hover">SaaS</a>
        <a href="#contact" class="link link-hover">Contact</a>
      </nav>
      <nav class="grid grid-flow-col gap-4">
        <a href="https://github.com/realdatadriven/etlx" target="_blank" rel="noopener" class="link link-hover">etlx</a>
        <a href="https://github.com/realdatadriven/central-set-go" target="_blank" rel="noopener" class="link link-hover">central-set-go</a>
      </nav>
      <aside>
        <p>&copy; {{.UI.ui_name}} &middot; Data engineering consulting, powered by our own open-source stack.</p>
      </aside>
    </footer>
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## UI_PAGE_HOME
```yaml
table: ui_page
description: Add the RealDataDriven landing page (default page)
cond: 'WHERE ui_id = :ui_id AND page_key = :page_key AND excluded = false'
data:
  ui_page_id: 1
  ui_id: 1
  page_key: home
  page_title: RealDataDriven - Data Engineering Consulting
  meta_description: ETL/ELT/Reverse ETL, data warehouses, lakes and lakehouses, and spreadsheet-replacing data-entry apps - built on our own open-source stack, etlx and central-set-go.
  page_template: |
    <!DOCTYPE html>
    <html lang="en" data-theme="corporate">
    <head>
      <meta charset="UTF-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <title>{{.Page.page_title}}</title>
      <meta name="description" content="{{.Page.meta_description}}" />
      <!-- Tailwind + DaisyUI via CDN (pin/replace with your own build for production) -->
      <script src="https://cdn.tailwindcss.com"></script>
      <link href="https://cdn.jsdelivr.net/npm/daisyui@4.12.14/dist/full.min.css" rel="stylesheet" type="text/css" />
      <!-- htmx -->
      <script src="https://unpkg.com/htmx.org@1.9.12"></script>
    </head>
    <body class="min-h-screen bg-base-100">
      {{template "header" .}}

      <!-- HERO -->
      <section class="hero min-h-[70vh] bg-gradient-to-br from-base-200 to-base-100">
        <div class="hero-content text-center py-20">
          <div class="max-w-2xl">
            <h1 class="text-4xl md:text-5xl font-bold">Data engineering, simplified.</h1>
            <p class="py-6 text-lg opacity-80">
              10+ years designing and running data pipelines, warehouses, lakes and lakehouses,
              on-prem and in the cloud. ETL, ELT, Reverse ETL - and simple data-entry apps that
              finally replace the spreadsheets everyone's afraid to touch.
            </p>
            <div class="flex flex-wrap gap-3 justify-center">
              <a href="#saas" class="btn btn-primary">Try the free demo</a>
              <a href="#contact" class="btn btn-outline">Get in touch</a>
            </div>
          </div>
        </div>
      </section>

      <!-- SERVICES -->
      <section id="services" class="py-20 px-4 lg:px-12 max-w-6xl mx-auto">
        <h2 class="text-3xl font-bold text-center mb-2">What we do</h2>
        <p class="text-center opacity-70 mb-12">Data engineering end to end, not just tools.</p>
        <div class="grid md:grid-cols-3 gap-6">
          <div class="card bg-base-100 shadow border border-base-200">
            <div class="card-body">
              <h3 class="card-title">ETL / ELT / Reverse ETL</h3>
              <p>Pipelines that move and transform your data reliably, in whichever direction your business needs it.</p>
            </div>
          </div>
          <div class="card bg-base-100 shadow border border-base-200">
            <div class="card-body">
              <h3 class="card-title">Warehouses, lakes &amp; lakehouses</h3>
              <p>Database, data warehouse, data lake and lakehouse design - sized to what you actually need, not just what's trendy.</p>
            </div>
          </div>
          <div class="card bg-base-100 shadow border border-base-200">
            <div class="card-body">
              <h3 class="card-title">Data-entry apps</h3>
              <p>Lightweight, database-backed apps that replace the spreadsheets your team is quietly running the business on.</p>
            </div>
          </div>
          <div class="card bg-base-100 shadow border border-base-200 md:col-span-3">
            <div class="card-body">
              <h3 class="card-title">10+ years, local and cloud</h3>
              <p>On-premises or cloud, we've deployed and operated these systems on both for over a decade - we know the trade-offs, not just the marketing.</p>
            </div>
          </div>
        </div>
      </section>

      <!-- OPEN SOURCE -->
      <section id="opensource" class="py-20 px-4 lg:px-12 bg-base-200">
        <div class="max-w-6xl mx-auto">
          <h2 class="text-3xl font-bold text-center mb-2">Built on our own open source</h2>
          <p class="text-center opacity-70 mb-12">Everything we deliver is backed by tools we build and maintain in the open.</p>
          <div class="grid md:grid-cols-2 gap-6">
            <div class="card bg-base-100 shadow">
              <div class="card-body">
                <h3 class="card-title">etlx</h3>
                <p>Our open-source ETL/ELT/Reverse ETL engine - the pipeline core behind every project we run.</p>
                <div class="card-actions justify-end">
                  <a href="https://github.com/realdatadriven/etlx" target="_blank" rel="noopener" class="btn btn-sm btn-outline">View on GitHub</a>
                </div>
              </div>
            </div>
            <div class="card bg-base-100 shadow">
              <div class="card-body">
                <h3 class="card-title">central-set-go</h3>
                <p>Our database-backed application platform: data-entry apps, admin, APIs and now dynamic websites, deployable anywhere.</p>
                <div class="card-actions justify-end">
                  <a href="https://github.com/realdatadriven/central-set-go" target="_blank" rel="noopener" class="btn btn-sm btn-outline">View on GitHub</a>
                </div>
              </div>
            </div>
          </div>
          <p class="text-center opacity-70 mt-8">
            central-set-go + etlx is where it shines the most: a full platform for building and running
            production data pipelines - but our consulting doesn't stop there. Any data integration
            challenge is fair game.
          </p>
        </div>
      </section>

      <!-- SAAS -->
      <section id="saas" class="py-20 px-4 lg:px-12 max-w-6xl mx-auto">
        <h2 class="text-3xl font-bold text-center mb-2">Get a Central-Set instance running</h2>
        <p class="text-center opacity-70 mb-12">From an instant free test to a fully managed, negotiated deployment.</p>
        <div class="grid md:grid-cols-3 gap-6 items-stretch">

          <!-- Free demo: self-service, one click -->
          <div class="card bg-base-100 shadow border border-base-200 flex flex-col">
            <div class="card-body flex-1 flex flex-col">
              <div class="badge badge-primary mb-2">Free</div>
              <h3 class="card-title">Free demo</h3>
              <p class="flex-1">
                Instant, one-click Central-Set instance in our cloud, for a quick test drive.
                No talking to sales - just click and it's up.
              </p>
              <form hx-post="/api/saas/free-demo"
                    hx-target="#saas-result"
                    hx-swap="innerHTML"
                    hx-indicator="#free-demo-spinner"
                    class="mt-4">
                <button type="submit" class="btn btn-primary w-full">
                  Deploy free demo
                  <span id="free-demo-spinner" class="htmx-indicator loading loading-spinner loading-sm"></span>
                </button>
              </form>
            </div>
          </div>

          <!-- Guided demo: still self-service, more resources, time-boxed -->
          <div class="card bg-base-100 shadow border border-base-200 flex flex-col">
            <div class="card-body flex-1 flex flex-col">
              <div class="badge badge-secondary mb-2">Guided</div>
              <h3 class="card-title">Guided demo</h3>
              <p class="flex-1">
                More headroom to really try it out: up to 150 GB NVMe disk, 10 GB memory and
                30 TB bandwidth, for a limited time.
              </p>
              <form hx-post="/api/saas/guided-demo"
                    hx-target="#saas-result"
                    hx-swap="innerHTML"
                    hx-indicator="#guided-demo-spinner"
                    class="mt-4 space-y-3">
                <select name="duration_days" class="select select-bordered select-sm w-full">
                  <option value="3">3 days</option>
                  <option value="5">5 days</option>
                  <option value="15">15 days</option>
                </select>
                <button type="submit" class="btn btn-secondary w-full">
                  Deploy guided demo
                  <span id="guided-demo-spinner" class="htmx-indicator loading loading-spinner loading-sm"></span>
                </button>
              </form>
            </div>
          </div>

          <!-- Personalised setup: requires contact, not self-service -->
          <div class="card bg-base-100 shadow border border-base-200 flex flex-col">
            <div class="card-body flex-1 flex flex-col">
              <div class="badge badge-accent mb-2">Custom</div>
              <h3 class="card-title">Personalised setup</h3>
              <p class="flex-1">
                Deployed on your own servers or in the cloud, sized and configured for you.
                And since our consulting goes well beyond Central-Set and etlx, this is also
                where any broader data-integration need starts.
              </p>
              <a href="#contact" class="btn btn-accent btn-outline w-full mt-4">Get in touch</a>
            </div>
          </div>
        </div>

        <div id="saas-result" class="max-w-xl mx-auto mt-8"></div>
      </section>

      <!-- CONTACT -->
      <section id="contact" class="py-20 px-4 lg:px-12 bg-base-200">
        <div class="max-w-xl mx-auto">
          <h2 class="text-3xl font-bold text-center mb-2">Get in touch</h2>
          <p class="text-center opacity-70 mb-8">Tell us about your data - we'll tell you how we'd approach it.</p>
          <form hx-post="/api/contact" hx-target="#contact-result" hx-swap="innerHTML" class="card bg-base-100 shadow p-6 space-y-4">
            <input type="text" name="name" placeholder="Name" required class="input input-bordered w-full" />
            <input type="email" name="email" placeholder="Email" required class="input input-bordered w-full" />
            <textarea name="message" placeholder="What are you trying to solve?" required class="textarea textarea-bordered w-full h-28"></textarea>
            <button type="submit" class="btn btn-primary w-full">Send</button>
          </form>
          <div id="contact-result" class="mt-4"></div>
        </div>
      </section>

      {{template "footer" .}}
    </body>
    </html>
  cache_seconds: 60
  default_page: true
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## UI_PAGE_LOGIN
```yaml
table: ui_page
description: Add the login page (page_key "login" -> GET /ui/rdd/login)
cond: 'WHERE ui_id = :ui_id AND page_key = :page_key AND excluded = false'
data:
  ui_page_id: 2
  ui_id: 1
  page_key: login
  page_title: Log in - RealDataDriven
  meta_description: Log in to your RealDataDriven account
  page_template: |
    <!DOCTYPE html>
    <html lang="en" data-theme="corporate">
    <head>
      <meta charset="UTF-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <title>{{.Page.page_title}}</title>
      <meta name="description" content="{{.Page.meta_description}}" />
      <script src="https://cdn.tailwindcss.com"></script>
      <link href="https://cdn.jsdelivr.net/npm/daisyui@4.12.14/dist/full.min.css" rel="stylesheet" type="text/css" />
      <script src="https://unpkg.com/htmx.org@1.9.12"></script>
      <script>
        // By default htmx only swaps in 2xx/3xx responses; a 4xx/5xx just
        // fires htmx:responseError and leaves the DOM untouched. Our /login
        // failure response is a normal HTML fragment (a DaisyUI alert) meant
        // to be swapped in, not just logged - so tell htmx to treat 4xx/5xx
        // as swappable content too.
        document.addEventListener('DOMContentLoaded', function () {
          htmx.config.responseHandling = [
            { code: "204", swap: false },
            { code: "[23]..", swap: true },
            { code: "[45]..", swap: true, error: false },
            { code: "...", swap: true }
          ];
        });
      </script>
    </head>
    <body class="min-h-screen bg-base-100">
      {{template "header" .}}

      <section class="hero min-h-[70vh]">
        <div class="hero-content w-full max-w-sm">
          <div class="card w-full bg-base-100 shadow border border-base-200">
            <div class="card-body">
              <h1 class="card-title text-2xl mb-2">Log in</h1>

              <div id="login-error" class="mb-4"></div>

              <form hx-post="/ui/rdd/login"
                    hx-target="#login-error"
                    hx-swap="innerHTML"
                    hx-indicator="#login-spinner"
                    class="space-y-4">
                <input type="email" name="email" placeholder="Email" required class="input input-bordered w-full" />
                <input type="password" name="password" placeholder="Password" required class="input input-bordered w-full" />
                <button type="submit" class="btn btn-primary w-full">
                  Log in
                  <span id="login-spinner" class="htmx-indicator loading loading-spinner loading-sm"></span>
                </button>
              </form>
            </div>
          </div>
        </div>
      </section>

      {{template "footer" .}}
    </body>
    </html>
  cache_seconds: 0
  default_page: false
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

# --- backend endpoints referenced above, not part of the ui/ui_page model --
# The three htmx forms above POST to:
#   /api/saas/free-demo
#   /api/saas/guided-demo
#   /api/contact
# These perform side effects (provisioning, sending an email) rather than
# rendering a stored template, so they're regular application handlers you
# write yourself (e.g. app.serve_saas_free_demo), not RenderUIPage /
# RenderUIPartial calls - those two stay read-only / template-only. Each
# handler just needs to respond with an HTML fragment (a DaisyUI alert works
# well) for hx-target="#saas-result" / "#contact-result" to swap in, e.g.:
#
#   <div class="alert alert-success">
#     <span>Your free demo is being deployed - check your email shortly.</span>
#   </div>
#
#   <div class="alert alert-success">
#     <span>Your free demo is being deployed - check your email shortly.</span>
#   </div>
#
# --- POST /ui/rdd/login, referenced by the login page above -----------------
# Implemented as app.serve_ui_login (a plain handler, not RenderUIPage): it
# reads the htmx form post (application/x-www-form-urlencoded, not JSON),
# delegates the credential check to the existing app.dynamic_login, and on
# success sets an HttpOnly "session" cookie plus an HX-Redirect header (no
# JSON in the response). On failure it returns a small HTML alert fragment
# with a 4xx status, which the login page's htmx.config.responseHandling
# tweak allows to be swapped into #login-error.