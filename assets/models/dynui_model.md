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
  ui_name:          { type: varchar, len: 200, nullable: false, comment: "Website Name", form_display: true, table_display: true, form_size: 5, order: 1 }
  ui_slug:          { type: varchar, len: 100, unique: true, nullable: false, comment: "URL Slug", form_display: true, table_display: true, form_size: 3, order: 2, form_regex_val: "^[a-z][a-z0-9-]*$", form_val_msg: "Use lowercase letters, numbers and hyphens; begin with a letter." }
  ui_desc:          { type: text, comment: "Description", form_display: true, table_display: true, form_long_text: true, order: 5 }
  default_locale:   { type: varchar, len: 10, default: "en", comment: "Default Locale", form_display: true, table_display: true, form_size: 2, form_order: 3 }
#  database:         { type: varchar, len: 50, default: "en", comment: "Database", form_display: true, table_display: true, form_size: 3, order: 4 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, form_order: 4 }
  user_id:          { type: integer, comment: "User ID" }
  app_id:           { type: integer, comment: "App ID" }
  created_at:       { type: datetime, comment: "Created At" }
  updated_at:       { type: datetime, comment: "Updated At" }
  excluded:         { type: boolean, default: false, comment: "Excluded" }
form_layout:
  tabs_steps: tabs
  form_in_popup: false
  size: 9
  allow_in_subform: {ui_page: true, ui_partial: true, ui_asset: true}
#  tabs_steps_conf: [{label: Website, fields: [ui_slug, ui_name, ui_desc, default_locale, active]}]
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
  ui_partial:        { type: varchar, len: 150, nullable: false, comment: "Partial Name", form_display: true, table_display: true, form_size: 7, order: 2, form_regex_val: "^[a-z][a-z0-9_-]*$", form_val_msg: "Use lowercase letters, numbers, underscores and hyphens." }
  ui_partial_desc:   { type: text, comment: "Partial Description", form_display: true, table_display: true, form_long_text: true, order: 4 }
  partial_template:  { type: text, nullable: false, comment: "Partial Template", form_display: true, form_long_text: true, form_code: html, order: 5}
  active:            { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, form_order: 3 }
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
    - {label: Partial, fields: [ui, ui_partial, ui_partial_desc, active]}
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
  ui_partial_data:      { type: varchar, len: 100, nullable: false, comment: "Data Name", form_display: true, table_display: true, form_size: 6, order: 3, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not begin with a number; use letters, numbers and underscores." }
  ui_partial_data_desc: { type: text, comment: "Data Description", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 7 }
  odata_path:           { type: text, nullable: false, comment: "OData Path", form_display: true, table_display: true, form_size: 8, order: 6 }
  single_row_obj:       { type: boolean, default: false, comment: "Single Row", form_display: true, table_display: true, form_size: 2, form_order: 5 }
  active:               { type: boolean, default: true, comment: "Active", table_display: true, form_display: true, form_size: 2, form_order: 4 }
  ui_partial_id:        { type: integer, fk: "ui_partial.ui_partial_id", nullable: false, comment: "Partial", form_display: true, table_display: true, form_size: 3, order: 2 }
  ui_id:                { type: integer, fk: "ui.ui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 3, order: 1 }
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
  page_title:       { type: varchar, len: 255, nullable: false, comment: "Title", form_display: true, table_display: true, form_size: 4, order: 3 }
  meta_description: { type: text, comment: "SEO Description", form_display: true, form_long_text: false, order: 6, form_size: 8 }
  page_template:    { type: text, nullable: false, comment: "Page Template", form_display: true, form_long_text: true, form_code: html, order: 10 }
  cache_seconds:    { type: integer, default: 0, comment: "Public Cache Seconds", form_display: true, table_display: true, form_size: 2, order: 7 }
  default_page:     { type: boolean, default: true, comment: "Default Page", form_display: true, table_display: true, form_size: 2, order: 8 }
  login_required:   { type: boolean, default: false, comment: "Requires Login", form_display: true, table_display: true, form_size: 2, order: 9 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, form_order: 4 }
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
  allow_in_subform: {ui_page_data: true, ui_page_partial: true}
  tabs_steps_conf:
    - {label: Page, fields: [ui, page_key, page_title, meta_description, cache_seconds, default_page, login_required, active]}
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
  ui_page_data:      { type: varchar, len: 100, nullable: false, comment: "Data Name", form_display: true, table_display: true, form_size: 4, order: 3, form_regex_val: "^[A-Za-z_][A-Za-z0-9_]*$", form_val_msg: "Must not begin with a number; use letters, numbers and underscores." }
  ui_page_data_desc: { type: text, comment: "Data Description", form_display: true, form_long_text: true, form_code: text, table_display: true, form_size: 12, order: 7 }
  odata_path:        { type: text, nullable: false, comment: "OData Path", form_display: true, table_display: true, form_size: 10, order: 5 }
  single_row_obj:    { type: boolean, default: false, comment: "Single Row Object", form_display: true, table_display: true, form_size: 3, order: 6 }
  active:            { type: boolean, default: true, comment: "Active", table_display: true, form_display: true, form_size: 2, order: 4 }
  ui_page_id:        { type: integer, fk: "ui_page.ui_page_id", nullable: false, comment: "Page", form_display: true, table_display: true, form_size: 3, order: 2 }
  ui_id:             { type: integer, fk: "ui.ui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 3, order: 1 }
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

## UI_PAGE_PARTIAL
```yaml
table: ui_page_partial
comment: UI Page Partial
tooltip: Partial to be available to the page
columns:
  ui_page_partial_id: { type: integer, pk: true, autoincrement: true, comment: "Page Partial ID" }
  ui_page_partial:    { type: varchar, len: 100, nullable: false, comment: "Page Partial", form_display: true, table_display: true, form_size: 9, order: 4 }
  active:             { type: boolean, default: true, comment: "Active", table_display: true, form_display: true, form_size: 3, order: 5 }
  ui_partial_id:      { type: integer, fk: "ui_partial.ui_partial_id", nullable: false, comment: "Partial", form_display: true, table_display: true, form_size: 4, order: 3 }
  ui_page_id:         { type: integer, fk: "ui_page.ui_page_id", nullable: false, comment: "Page", form_display: true, table_display: true, form_size: 4, order: 2 }
  ui_id:              { type: integer, fk: "ui.ui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 4, order: 1 }
  user_id:            { type: integer, comment: "User ID" }
  app_id:             { type: integer, comment: "App ID" }
  created_at:         { type: datetime, comment: "Created At" }
  updated_at:         { type: datetime, comment: "Updated At" }
  excluded:           { type: boolean, default: false, comment: "Excluded" }
form_layout:
  size: 5
table_layout:
  default_order: [{field: ui_page_partial_id, order: ASC}]
```

## UI_ASSET
```yaml
table: ui_asset
comment: UI Asset
tooltip: Store UTF-8 assets directly or binary assets as base64 text. The content is database-backed, not a filesystem path.
columns:
  ui_asset_id:      { type: integer, pk: true, autoincrement: true, comment: "Asset ID" }
  asset_path:       { type: varchar, len: 500, nullable: false, comment: "Relative Asset Path", form_display: true, table_display: true, form_size: 5, order: 2 }
  mime_type:        { type: varchar, len: 150, nullable: false, comment: "MIME Type", form_display: true, table_display: true, form_size: 3, order: 3 }
  content_encoding: { type: varchar, len: 20, default: "utf-8", nullable: false, comment: "Encoding", form_display: true, table_display: true, form_size: 3, order: 5 }
  asset_content:    { type: text, nullable: false, comment: "Content", form_display: true, form_long_text: true, form_code: text, order: 8 }
  checksum:         { type: varchar, len: 128, comment: "Content Checksum", form_display: true, table_display: true, form_size: 3, order: 6 }
  cache_seconds:    { type: integer, default: 86400, comment: "Cache Seconds", form_display: true, table_display: true, form_size: 3, order: 7 }
  active:           { type: boolean, default: true, comment: "Active", form_display: true, table_display: true, form_size: 2, order: 4 }
  ui_id:            { type: integer, fk: "ui.ui_id", nullable: false, comment: "Website", form_display: true, table_display: true, form_size: 2, order: 1 }
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

<!--
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
-->

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
  partial_template: FileContent(ui/parts/header.html)
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
  partial_template: FileContent(ui/parts/footer.html)
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## UI_PARTIAL_CRUD
```yaml
table: ui_partial
description: Add the RealDataDriven CRUD partial
cond: 'WHERE ui_id = :ui_id AND ui_partial = :ui_partial AND excluded = false'
data:
  ui_partial_id: 3
  ui_id: 1
  ui_partial: crud
  ui_partial_desc: Crud with open-source links and copyright
  partial_template: FileContent(ui/parts/general_crud.html)
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
  children:
    table: ui_partial_data
    data:
      ui_partial_data: crud_data
      ui_partial_data_desc: CRUD Data
      odata_path: |
        {{.PathParams.db}}/{{.PathParams.table}}?
          {{if .PathParams.select}}&$select={{.PathParams.select}}{{end}}
          {{if .PathParams.filter}}&$filter={{.PathParams.filter}}{{end}}
          {{if .PathParams.top}}&$top={{.PathParams.top}}{{end}}
          {{if .PathParams.skip}}&$skip={{.PathParams.skip}}{{end}}
          {{if .PathParams.schema}}&$schema={{.PathParams.schema}}{{end}}
          {{if .PathParams.format}}&$format={{.PathParams.format}}{{end}}
      ui_partial_id: ui_partial_id()
      ui_id: ui_id()
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
  page_template: FileContent(ui/landing.html)
  cache_seconds: 60
  default_page: true
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## UI_PAGE_DASH
```yaml
table: ui_page
description: Add the RealDataDriven Dashboard
cond: 'WHERE ui_id = :ui_id AND page_key = :page_key AND excluded = false'
data:
  ui_id: 1
  page_key: dashboard
  page_title: RealDataDriven - Data Engineering Consulting
  meta_description: ETL/ELT/Reverse ETL, data warehouses, lakes and lakehouses, and spreadsheet-replacing data-entry apps - built on our own open-source stack, etlx and central-set-go.
  page_template: FileContent(ui/dashboard.html)
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
  ui_id: 1
  page_key: login
  page_title: Log in - RealDataDriven
  meta_description: Log in to your RealDataDriven account
  page_template: FileContent(ui/auth/design-2-login.html)
  cache_seconds: 0
  default_page: false
  active: true
  user_id: 1
  app_id: appId()
  created_at: Now()
  updated_at: Now()
  excluded: false
```

## UI_ASSET_EX
```yaml
table: ui_asset
description: Add the ui asset example
cond: 'WHERE ui_id = :ui_id AND asset_path = :asset_path AND excluded = false'
data:
  - asset_path:       logo.svg
    mime_type:        MimeType(assets/static/img/logo.svg)
    content_encoding: utf-8
    asset_content:    FileContent(assets/static/img/logo.svg)
    active:           true
    ui_id:            1
  - asset_path:       logo.png
    mime_type:        MimeType(assets/static/img/icon.png)
    content_encoding: utf-8
    asset_content:    Base64(assets/static/img/icon.png)
    active:           true
    ui_id:            1
```
