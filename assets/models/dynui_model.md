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
    - {label: Page, fields: [ui_id, page_key, page_title, meta_description, cache_seconds, default_page, active]}
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