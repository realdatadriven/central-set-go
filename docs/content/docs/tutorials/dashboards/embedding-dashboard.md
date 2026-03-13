---
weight: 7085
title: "Embedding Dashboards"
description: "How to embed Central Set dashboards inside external applications using an iframe and JWT authentication."
icon: dashboard
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

{{% alert context="warning" text="**Caution** — Central Set is **production-ready and actively used**, but this documentation is still **under active development**. Large parts of the docs are **auto-generated**, evolving alongside the platform, and some sections may be incomplete or change frequently."  /%}}

Central Set dashboards can be **embedded directly into external applications** using a simple **iframe integration**.

This allows dashboards built in Central Set to be reused inside:

- internal portals
- client applications
- admin panels
- SaaS products
- monitoring tools

The embedding system uses **secure cross-window communication (`postMessage`)** combined with **JWT authentication**.

---

# 🧩 Basic Embedding Example

An external application can embed a Central Set dashboard using an iframe.

Example:

```html {linenos=table}
<!-- other-site.html -->

<iframe 
  id="my-frame"
  src="http://localhost:4444" 
  width="800" 
  height="600"
  sandbox="allow-scripts allow-same-origin"
/>

<script>
try {        
    const iframe = document.getElementById('my-frame');        

    iframe.onload = () => {
        const targetOrigin = 'http://localhost:4444'; // MUST match iframe origin        
        iframe.contentWindow.postMessage({
            embedded_dashboard: true,
            app: {app_id: 2, app: 'ETLX', db: 'ETLX'},
            table: 'dashboard',
            dashboard_id: 1,
            token: '<JWT_TOKEN>',
            theme: 'light'
        }, targetOrigin);
    };
} catch (e) {
    console.error("Error setting up message listener:", e);
}
</script>
````

Once the iframe loads, the parent application sends a **configuration message** to the Central Set frontend.

The dashboard is then loaded **inside the embedded frame**.

---

# 📦 Message Parameters

The `postMessage` payload defines what dashboard to load and how it should behave.

Example payload:

```json {linenos=table}
{
  "embedded_dashboard": true,
  "app": {
    "app_id": 2,
    "app": "ETLX",
    "db": "ETLX"
  },
  "table": "dashboard",
  "dashboard_id": 1,
  "token": "<JWT_TOKEN>",
  "theme": "light"
}
```

Parameter reference:

| Field                | Description                           |
| -------------------- | ------------------------------------- |
| `embedded_dashboard` | Enables embedded dashboard mode       |
| `app`                | Application context                   |
| `table`              | Table containing dashboards           |
| `dashboard_id`       | Dashboard identifier                  |
| `token`              | Central Set JWT authentication token  |
| `theme`              | Optional UI theme (`light` or `dark`) |

---

# 🧠 Application Context

The `app` object defines **which Central Set application database the dashboard belongs to**.

Example:

```json {linenos=table}
{
  "app_id": 2,
  "app": "ETLX",
  "db": "ETLX"
}
```

This ensures:

* correct metadata loading
* correct table resolution
* proper access control

---

# 🔐 Authentication

Embedding dashboards requires a **valid Central Set JWT token**.

Example:

```json {linenos=table}
"token": "<JWT_TOKEN>"
```

This token must belong to a **user with access to the requested resources**.

All normal Central Set security layers still apply:

* role permissions
* table access
* row level access
* column level access

The embedded dashboard therefore respects **the same security rules as the API**.

---

# 🔑 Token Sources

There are two ways to obtain a valid token.

## 1. Central Set Authentication

If the embedding application **already authenticates against a Central Set backend**, the same JWT token can be reused.

Example flow:

```
User login → CS backend
           → JWT issued
           → used for API + embedded dashboards
```

This is the **recommended approach**.

---

## 2. Access Keys

If the external application **does not authenticate directly with Central Set**, a token can be generated in the Admin UI.

Location:

```
Admin → Access Keys
```

Access keys create **long-lived tokens** that can be used by:

* external services
* integrations
* embedded dashboards

---

# ⚠️ Security Considerations

Tokens **should never be embedded directly in client-side applications**.

❌ Avoid this:

```
Hardcoding JWT tokens inside frontend code
```

This exposes the token to:

* browser inspection
* network interception
* client-side extraction

---

## ✅ Recommended Pattern

Serve tokens from a **secure backend**.

Example architecture:

```
User Browser
     │
     ▼
Your Application Backend
     │
     │ obtains / stores JWT
     ▼
Central Set API
```

The backend then injects the token into the page or returns it via a **secure API endpoint**.

Example flow:

```
1. User opens your app
2. Backend authenticates with Central Set
3. Backend returns a temporary JWT
4. Frontend embeds dashboard using that token
```

This ensures:

* tokens are **not exposed**
* authentication remains **controlled**
* access can be **revoked centrally**

---

# 🎨 Theme Support

The embed configuration optionally supports a theme parameter.

Example:

```json
"theme": "light"
```

Possible values:

```
light
dark
```

If omitted, the default theme is:

```
light
```

---

# 🚀 Use Cases

Embedding dashboards is useful for:

* customer portals
* SaaS analytics
* operational monitoring
* internal reporting tools
* ETLX pipeline monitoring

Because the dashboard runs inside Central Set, you benefit from:

* dynamic dashboards
* SQL data sources
* CRUD API data sources
* ETLX pipelines
* built-in security layers

---

# 🧠 Summary

Central Set dashboards can be embedded inside any application using:

* an **iframe**
* a **postMessage configuration**
* a **valid JWT token**

This enables Central Set to act as a **dashboard engine for external systems**, while maintaining:

* centralized security
* dynamic data access
* reusable dashboards
* unified data governance.
