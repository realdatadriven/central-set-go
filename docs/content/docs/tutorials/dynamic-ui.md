---
weight: 7065
title: "Dynamic UI Applications"
description: "Build and deploy database-backed websites with Go templates, OData data bindings, reusable partials, and Central Set authentication."
icon: web
date: 2026-10-18
lastmod: 2026-10-18
draft: false
images: []
---

{{% alert context="warning" text="Dynamic UI templates are powerful server-side configuration. Only trusted administrators should be allowed to edit them. OData-backed data uses Central Set's read and access-control path; embedded SQL blocks are a separate execution path and do not automatically inherit OData permissions or row-level access." /%}}

## What is a Dynamic UI?

A Dynamic UI is a database-backed website served by Central Set Go. Its site, pages, reusable partials, data bindings, and static assets are configured as records rather than compiled into the main frontend. Once the records are active, the site is available at `/ui/{ui_slug}` without rebuilding or redeploying the application.

This is a different tool from the opinionated ADMIN application:

| ADMIN application | Dynamic UI |
| --- | --- |
| Uses Central Set's configured menus, tables, forms, and datatables | Uses pages and HTML authored with Go's `html/template` syntax |
| Follows the standard admin layout and navigation | You choose the layout, navigation, and frontend libraries |
| Best for back-office data management | Best for focused portals, dashboards, public pages, and custom workflows |

Both can work with the same registered applications and data-access rules. A Dynamic UI is not a new database or a bypass around Central Set's authorization: its OData bindings resolve through the application's existing read path.

## The configuration model

The UI model in `assets/models/dynui_model.md` defines these main records:

* **`ui`** — the website's name, URL slug, default locale, and active state.
* **`ui_page`** — a page key, Go HTML template, optional response template, cache duration, default-page flag, and login requirement.
* **`ui_page_data`** — a name and OData path whose result is added to the page's template data.
* **`ui_partial`** — a reusable named Go template, such as a header, footer, or CRUD component.
* **`ui_page_partial`** — associates reusable partials with a page.
* **`ui_partial_data`** — data bindings available to a partial.
* **`ui_asset`** — database-backed static files served for a UI.

Use the UI model from the admin application to create and maintain these records. Activate the UI, page, binding, and any associated partials; the server reads the configuration at request time. No frontend build is required.

The request flow is implemented in `cmd/api/dyn_ui.go`. The UI routes in `cmd/api/routes.go` send page requests to `serve_ui_page`, which calls `RenderUIPage`; that function finds the active UI and page, resolves active data bindings, parses the page and its associated partials, and renders the response.

## Build a page backed by OData

Suppose the registered `SALES` application exposes an `orders` table. Create a UI with slug `sales-portal`, then add a page with key `orders`. Associate a data row with that page:

```text
ui_page_data: orders
odata_path:   SALES/orders?$top=25&$orderby=created_at desc
active:       true
```

The path is an OData path in the form `{app-or-db}/{table}` followed by optional OData query options. During rendering, `ui_page_data` is resolved through `ODataRead`, which maps the request to the Central Set read logic. The returned result is available to the Go template under the configured data name:

```html
<h1>Recent orders</h1>
<ul>
  {{range .orders.data}}
    <li>
      Order {{.order_id}} — {{.status}}
    </li>
  {{else}}
    <li>No orders found.</li>
  {{end}}
</ul>
```

Open the page at:

```text
/ui/sales-portal/orders
```

Bindings can use URL query parameters through `.PathParams`. For example, the model's store example filters a store record by the page's `id` query parameter:

```text
odata_path: store/store?$filter=store_id eq {{.PathParams.id}}
```

The page is then available at a URL such as `/ui/store/home?id=42`. Treat values from `.PathParams` as user input: constrain or validate them before interpolating them into an OData filter, and do not let request parameters select arbitrary databases, tables, or query structure.

By default, a binding exposes the OData result envelope, including its `data` field. Set `single_row_obj: true` when the template should receive one row object in `.data` rather than a collection; the implementation uses the first returned row, or an empty object if there are no rows. The model's store and dashboard records demonstrate these single-row bindings.

Page templates use Go's `html/template`, with the Sprig function map and Central Set's `T` and `pluckJson` helpers. The `user`, `UI`, `Page`, `Path`, `Host`, `PathParams`, and `config` values are also available to templates. HTML is contextually escaped by the HTML template engine; avoid bypassing that protection for user-controlled content.

## Reuse page components and add frontend assets

Create a named `ui_partial` such as `header`, then associate it to a page through `ui_page_partial`. The partial is parsed with the page's template set, so the page can invoke it as a named template:

```html
{{template "header" .}}

<main class="container mx-auto p-6">
  <h1>{{.Page.page_title}}</h1>
  {{range .orders.data}}
    <p>{{.order_id}} — {{.status}}</p>
  {{end}}
</main>
```

Partials can have their own `ui_partial_data` OData bindings. Active page partials and their bindings are loaded for rendering, allowing shared layouts and components to use the same template data context.

Add stylesheets, scripts, or other files as `ui_asset` rows. Dynamic UI assets are served below `/ui/{ui_slug}/static/{asset...}` (and the compatible `/asset/` route). This makes it practical to load a CDN-hosted Tailwind CSS or daisyUI stylesheet and use HTMX for progressive interactions; keep third-party asset URLs and versions explicit in the page template. For example:

```html
<link href="https://cdn.jsdelivr.net/npm/daisyui@latest/daisyui.css" rel="stylesheet">
<script src="https://unpkg.com/htmx.org"></script>
```

For production, pin library versions and apply your deployment's content-security and asset policies.

The model examples are useful starting points:

* `UI_PARTIAL_HEADER`, `UI_PARTIAL_FOOTER`, `UI_PARTIAL_CRUD`, and `UI_PARTIAL_CRUD_FORM` show reusable partial definitions and parameterized data bindings.
* `STORE_UI_PAGE_HOME` shows a page with a path-parameter-driven OData binding and associated partials.
* `SURVEY` and `DASHBOARD2` show other page and data-binding configurations.
* The seeded login, signup, recovery, reset, alter-password, and two-factor pages show response templates and auth-oriented pages.

## Optional SQL blocks

`RenderUIPage` also extracts fenced SQL blocks from the page template, renders their Go template expressions, executes them in an in-memory DuckDB connection, and makes named results available under `.queries`. A named block can look like this:

````markdown
```sql
INSTALL erpl_web FROM community;
LOAD erpl_web;
```

```sql
CREATE SECRET api_auth (
    TYPE http_bearer,
    TOKEN '{{.user.token}}',
    SCOPE '{{.Host}}'
);
```

```sql
ATTACH IF NOT EXISTS '{{.Host}}/odata/SALES' AS sales (TYPE ODATA);
```

```sql recent_totals
SELECT status, count(*) AS order_count
FROM sales.orders
GROUP BY status
```
````

The result can then be rendered as `{{range .queries.recent_totals}}...{{end}}`. The UI model's `DASHBOARD2` example is a reference for a dashboard-backed page.

**Important security boundary:** When querying C7-managed data from a Dynamic UI template, attach the database through the C7 OData endpoint rather than connecting directly to the underlying database. A direct database attachment does not enforce the CRUD privileges or row-level access rules configured in C7.

To enforce these permissions, use the `erpl_web` extension, authenticate with the current user's token, and attach the relevant C7 OData application:

```sql
INSTALL erpl_web FROM community;
LOAD erpl_web;
```

```sql
CREATE SECRET api_auth (
    TYPE http_bearer,
    TOKEN '{{.user.token}}',
    SCOPE '{{.Host}}'
);
```

```sql
ATTACH IF NOT EXISTS '{{.Host}}/odata/SALES' AS sales (TYPE ODATA);
```

Queries against this OData-attached database are evaluated through C7's authorization layer, enforcing the applicable CRUD privileges and row-level access rules for the authenticated user.

This approach is especially useful for dashboards and analytical queries. For example, a dashboard can aggregate sales, calculate department-level totals, or display filtered metrics while respecting the current user's department scope and other configured access restrictions. Aggregations and derived results must use the authorized OData source so that users cannot infer information from records they are not permitted to access.

**Rule of thumb:** Use direct database connections for data that does not require C7 authorization. For protected C7 application data, use the authenticated OData attachment so that queries, including dashboard filters and aggregations, respect the permissions defined in C7.

## Authentication and account flows

Dynamic UIs can use the configured Central Set session/token identity. Set `login_required` on pages that should reject anonymous access; otherwise the page may render with the anonymous user context. The UI route set includes handlers for:

* Login, signup, two-factor validation, and logout
* Password recovery, password reset, and password change
* Email confirmation
* OAuth provider login and callback

The snippets below follow the seeded `UI_PAGE_LOGIN`, `UI_PAGE_SIGNUP`, `UI_PAGE_RECOVER_PASSWORD`, and `UI_PAGE_RESET_PASSWORD` examples in `assets/models/dynui_model.md`. Replace `portal` with your UI's slug and make sure the corresponding page records are active.

### Login

The seeded login page is stored with `page_key: login`, uses a Go HTML page template, and defines a response template for displaying authentication errors:

```yaml
table: ui_page
data:
  page_key: login
  page_title: Log in
  page_template: FileContent(ui/auth/login.html)
  response_tmpl: |
    <div role="alert"
        class="alert {{ if .success }}alert-success{{ else }}alert-error{{ end }}">
        <span>{{ .msg }}</span>
    </div>
  cache_seconds: 0
  default_page: false
  active: true
```

A minimal login form posts the field names expected by the handler:

```html
<form method="post" action="/ui/portal/login">
  <label>
    Email
    <input type="email" name="email" autocomplete="username" required>
  </label>
  <label>
    Password
    <input type="password" name="password" autocomplete="current-password" required>
  </label>
  <button type="submit">Log in</button>
</form>
<a href="/ui/portal/password-recovery">Forgot your password?</a>
<a href="/ui/portal/signup">Create an account</a>
```

Successful login creates the UI session and redirects to `/ui/portal` (or sends an `HX-Redirect` for an HTMX request). If two-factor authentication is enabled, the handler sends the user to the two-factor page instead.

### Signup / registration

The model seeds a `signup` page in the same way and uses the configured dynamic-login table and field names for registration. The example deployment defaults are `users`, `user_id`, `username`, `email`, `password`, and `active`; check the `DYN_LOGIN_*` settings in your environment before using these field names.

```yaml
table: ui_page
data:
  page_key: signup
  page_title: Sign Up
  page_template: FileContent(ui/auth/signup.html)
  response_tmpl: |
    <div role="alert"
        class="alert {{ if .success }}alert-success{{ else }}alert-error{{ end }}">
        <span>{{ .msg }}</span>
    </div>
  cache_seconds: 0
  default_page: false
  active: true
```

The signup handler accepts the configured username, email, and password form fields:

```html
<form method="post" action="/ui/portal/signup">
  <label>
    Username
    <input type="text" name="username" autocomplete="username" required>
  </label>
  <label>
    Email
    <input type="email" name="email" autocomplete="email" required>
  </label>
  <label>
    Password
    <input type="password" name="password" autocomplete="new-password" required>
  </label>
  <button type="submit">Create account</button>
</form>
```

**Handler behavior:** the seeded model includes `response_tmpl` on the signup page, but the current signup route uses a generic action response and does not render that page-specific template. Do not rely on this `response_tmpl` for signup feedback unless the handler behavior changes.

### Password recovery and reset

The recovery page in the model is `page_key: password-recovery`. Its response template displays the result returned by the recovery action:

```yaml
table: ui_page
data:
  page_key: password-recovery
  page_title: Recover Password
  page_template: FileContent(ui/auth/password-recovery.html)
  response_tmpl: |
    {{if .success}}
      <div class="alert alert-soft alert-success mb-4">
        <span>Recovery email sent successfully</span>
      </div>
    {{else}}
      <div role="alert"
          class="alert {{ if .success }}alert-success{{ else }}alert-error{{ end }}">
          <span>{{ .msg }}</span>
      </div>
    {{end}}
  cache_seconds: 0
  default_page: false
  active: true
```

Post the email and include `ui` and `page` in the query string. The handler uses those values to find the page's `response_tmpl`. `recover_path` is optional; it tells the recovery action where the emailed reset link should lead:

```html
<form method="post"
      action="/ui/portal/recover_pass?ui=portal&amp;page=password-recovery">
  <label>
    Email
    <input type="email" name="email" autocomplete="email" required>
  </label>
  <input type="hidden" name="recover_path" value="/ui/portal/reset-password">
  <button type="submit">Send recovery email</button>
</form>
```

The seeded reset page uses `page_key: reset-password` and a response template for reset errors:

```yaml
table: ui_page
data:
  page_key: reset-password
  page_title: Reset Password
  page_template: FileContent(ui/auth/password-reset.html)
  response_tmpl: |
    <div role="alert"
        class="alert {{ if .success }}alert-success{{ else }}alert-error{{ end }}">
        <span>{{ .msg }}</span>
    </div>
  cache_seconds: 0
  default_page: false
  active: true
```

The emailed link includes a short-lived `token` query parameter; preserve it when submitting the new password:

```html
<form method="post"
      action="/ui/portal/reset_pass?ui=portal&amp;page=reset-password&amp;token={{.PathParams.token}}">
  <label>
    New password
    <input type="password" name="new_password" autocomplete="new-password" required>
  </label>
  <label>
    Confirm new password
    <input type="password" name="confirm_password" autocomplete="new-password" required>
  </label>
  <button type="submit">Reset password</button>
</form>
```

The reset handler requires matching passwords and enforces the password rules in the authentication code. On success it logs out the current session. Recovery email delivery and the reset URL depend on the server's email and frontend URL configuration; verify those settings before publishing the flow.

Login and account actions are available below `/ui/{ui_slug}`: `POST /login`, `POST /signup`, `POST /recover_pass`, and `POST /reset_pass`. The model also seeds alter-password and two-factor pages. OAuth uses `/ui/{ui_slug}/oauth/{provider}/login` and `/callback`. Configure the underlying Central Set user/authentication settings and providers as required; adding page templates alone does not configure identity providers or mail delivery.

## Capture and log requests centrally

Requests handled by the Dynamic UI and generic CRUD endpoints are logged in the `user_log` table in the `ADMIN` database (`ADMIN/user_log`).

## Generic CRUD endpoints

Dynamic UI pages are free to use their own HTML and interaction patterns. Separately, Central Set exposes a generic REST-style CRUD API for a registered application's tables:

```text
GET    /crud/{db}/{table}
GET    /crud/{db}/{table}/{id}
POST   /crud/{db}/{table}
PUT    /crud/{db}/{table}/{id}
PATCH  /crud/{db}/{table}/{id}
DELETE /crud/{db}/{table}/{id}
```

For example, a page might use HTMX to submit a form to `POST /crud/SALES/orders`, or fetch a row from `GET /crud/SALES/orders/42`. The `{db}` segment is resolved to a registered Central Set application; it is not an arbitrary connection string. Reads follow the OData-backed read path, and writes use Central Set's create/update flow. JSON requests receive JSON responses; HTML form requests can use a configured page `response_tmpl`.

This API is generic and not tied to the admin UI's menu or form layout, but it is still governed by the server's authentication, application/table permissions, and row-level access configuration. Grant only the access the intended users need, and verify read, create, update, and delete separately. The ADMIN interface is opinionated; a Dynamic UI lets you build another frontend over the same governed application data.

## Recommended workflow

1. Create or select the registered Central Set application and tables; configure the intended user, role, and row-level permissions.
2. In the UI model, create an active UI with a stable slug.
3. Add a default page and any additional pages, using Go HTML templates.
4. Add `ui_page_data` rows for OData results, then use those names in the templates.
5. Extract shared layout into active partials and associate them through `ui_page_partial`.
6. Add static assets or reference pinned CDN assets; use HTMX or another lightweight frontend library where useful.
7. Test the pages and CRUD endpoints as anonymous, permitted, and restricted users. Confirm row-level filtering as well as each write operation.
8. Deploy through the admin UI/model configuration, then activate the records. No application frontend rebuild is needed.

Dynamic UIs work particularly well for lightweight portals, AI-generated page drafts that an administrator reviews, custom dashboards, and small CRUD views. Keep the template and SQL authoring privileges restricted, favor OData for protected data, and treat every published page as application code.

## Related tutorials

* [Applications](../apps/) — how apps map to databases and configure their admin interface.
* [OData v4 API](../odata-v4/) — query options, token requirements, and data access.
* [Authentication](../auth/) — Central Set authentication configuration.
* [Analytics & Dashboards](./dashboards/overview/) — dashboard architecture and examples.
