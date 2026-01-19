---
weight: 7010
title: "Authentication & Access"
description: "User authentication, password management, and API access in Central Set Go"
icon: lock
date: 2025-12-16T01:04:15+00:00
lastmod: 2025-12-16T01:04:15+00:00
draft: false
images: []
---

## 🔐 Authentication & Access

Central Set Go (CSGO) provides a **unified authentication system** used consistently across:

- The **Admin UI**
- The **REST API**
- External integrations and automation

This section documents **each authentication-related feature** with:

- The **UI screen**
- The **user workflow**
- The **exact API call** backing that screen

> 🧠 The Admin UI is just another API client.  
> Everything you can do in the UI can be done programmatically.

---

## Login

The login process authenticates a user and returns a **JWT token** used for all subsequent API requests.

{{< tabs tabTotal="2" >}}

{{% tab tabName="UI" %}}

### Login Screen

The login screen is the entry point to CSGO.

- Users provide their **username** and **password**
- Language and theme (light/dark) can be selected
- Optionally supports “Remember me”

**Light mode**

![CSGO Login - Light Mode](/images/auth/auth_login_light.png)

**Dark mode**

![CSGO Login - Dark Mode](/images/auth/auth_login_dark.png)

Once authenticated, the user is redirected to the main admin interface and the token is stored client-side.

{{% /tab %}}

{{% tab tabName="API" %}}

### Login API

**Endpoint**

```

POST /dyn_api/login/login

````

Authenticates a user and returns a **JWT token**.

This endpoint:
- Does **not** require authentication
- Is the **entry point** for all user sessions

---

#### Headers

```http
Content-Type: application/json
````

---

#### Request Body

```json {linenos=table}
{
  "lang": "pt",
  "data": {
    "username": "root",
    "password": "1234"
  }
}
```

**Parameters**

| Field           | Description                          |
| --------------- | ------------------------------------ |
| `lang`          | Language code (`pt`, `en`, `es`, …)  |
| `data.username` | User login name                      |
| `data.password` | User password (encrypted in transit) |

---

#### Response

```json {linenos=table}
{
  "success": true,
  "token": "<JWT_TOKEN>",
  "user": {
    "user_id": 1,
    "username": "root",
    "first_name": "Super",
    "last_name": "Admin",
    "email": "root@domain.com",
    "role_id": 1,
    "lang_id": 1
  }
}
```

The returned token must be sent in all subsequent requests:

```http
Authorization: Bearer <JWT_TOKEN>
```

{{% /tab %}}

{{< /tabs >}}

---

## Change Password

Users can change their own password from the UI or via the API.

{{< tabs tabTotal="2" >}}

{{% tab tabName="UI" %}}

### Change Password Screen

The **Change Password** screen allows a logged-in user to:

* Provide the current password
* Set a new password
* Confirm the new password

This action immediately invalidates previous credentials and updates the stored password hash.

> 🔒 Password changes always require the **current password**.

{{% /tab %}}

{{% tab tabName="API" %}}

### Change Password API

**Endpoint**

```
POST /dyn_api/login/alter_password
```

---

#### Headers

```http
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

---

#### Request Body

```json {linenos=table}
{
  "lang": "en",
  "data": {
    "old_password": "1234",
    "new_password": "new_secure_password"
  }
}
```

**Parameters**

| Field          | Description      |
| -------------- | ---------------- |
| `old_password` | Current password |
| `new_password` | New password     |

---

#### Response

```json {linenos=table}
{
  "success": true,
  "message": "Password updated successfully"
}
```

{{% /tab %}}

{{< /tabs >}}

---

## API Access & Tokens

CSGO uses **JWT-based authentication** for all protected endpoints.

* Tokens are issued at login
* Tokens encode:

  * User ID
  * Username
  * Role assignments
  * Permissions
* Token expiration is configurable

### Using the Token

All authenticated requests must include:

```http
Authorization: Bearer <JWT_TOKEN>
```

### Security Model

* UI and API share the **same RBAC rules**
* Permissions are enforced **server-side**
* Tokens cannot bypass UI restrictions

---

## LDAP Authentication (Optional)

Central Set Go supports **direct authentication against an LDAP directory** instead of the internal `users` table.

When LDAP authentication is enabled:

* The **Login UI remains exactly the same**
* The **Login API endpoint remains exactly the same**
* Credentials are validated **against LDAP**
* User records are **resolved dynamically**, not stored locally
* Roles and permissions are still managed inside CSGO

> 🔁 **UI and API do not change** — only the authentication backend does.

---

## How LDAP Authentication Works

1. A user submits credentials via:

   * Login UI **or**
   * `POST /dyn_api/login/login`
2. CSGO validates credentials against the configured LDAP server
3. If authentication succeeds:

   * A JWT token is issued
   * The user session behaves like a normal CSGO user
4. Authorization (roles, permissions) is still enforced by CSGO

---

## Enabling LDAP Authentication

LDAP authentication is enabled **entirely via environment variables**.

### `.env` Configuration

```env
# Enable / Disable LDAP authentication
USE_LDAP_AUTH=false

# LDAP connection
LDAP_URL=ldap://localhost:1389
LDAP_BIND_USER=cn=admin,dc=example,dc=com
LDAP_PASSWORD=admin
LDAP_BASE_DN=dc=example,dc=com

# TLS / certificate behavior
LDAP_SKIP_VERIFY_CERT=true

# User lookup filter
LDAP_SEARCHREQ_FILTER="(|(uid=%[1]s)(cn=%[1]s)(mail=%[1]s))"
```

### Key Variables Explained

| Variable                | Description                                       |
| ----------------------- | ------------------------------------------------- |
| `USE_LDAP_AUTH`         | Enables LDAP authentication when set to `true`    |
| `LDAP_URL`              | LDAP server URL                                   |
| `LDAP_BIND_USER`        | Bind DN used for authentication                   |
| `LDAP_PASSWORD`         | Password for the bind user                        |
| `LDAP_BASE_DN`          | Base DN for user searches                         |
| `LDAP_SKIP_VERIFY_CERT` | Skip TLS certificate verification                 |
| `LDAP_SEARCHREQ_FILTER` | User search filter (supports username, CN, email) |

---

## Login Flow with LDAP Enabled

{{< tabs tabTotal="2" >}}

{{% tab tabName="UI" %}}

### Login UI (Unchanged)

The login screen remains exactly the same:

* Username
* Password
* Language
* Theme

The user does **not** need to know whether authentication is backed by:

* Local users table
* LDAP directory

![CSGO Login - LDAP uses the same UI](/images/auth/auth_login_light.png)

{{% /tab %}}

{{% tab tabName="API" %}}

### Login API (Unchanged)

**Endpoint**

```
POST /dyn_api/login/login
```

**Request Body**

```json {linenos=table}
{
  "lang": "en",
  "data": {
    "username": "jdoe",
    "password": "secret"
  }
}
```

CSGO will:

* Resolve `jdoe` via the configured LDAP filter
* Bind and authenticate against LDAP
* Issue a JWT token on success

{{% /tab %}}

{{< /tabs >}}

---

## Authorization with LDAP Users

LDAP handles **authentication only**.

CSGO remains responsible for:

* Roles
* Permissions
* App / menu / table access
* API authorization

This allows you to:

* Centralize identity in LDAP
* Keep fine-grained access control inside CSGO
* Use the same RBAC model for:

  * Local users
  * LDAP users
  * Service tokens

---

## Tested LDAP Setup

LDAP authentication has been tested using:

* **LDAP Server**
  `osixia/openldap:1.5.0`

* **LDAP Admin UI**
  `osixia/phpldapadmin:0.9.0`

Both managed via Docker.

This setup is suitable for:

* Local development
* Testing
* Integration with enterprise LDAP-compatible directories

---

## Summary

CSGO authentication supports:

* ✅ Local users (database-backed)
* ✅ LDAP-backed authentication
* ✅ Unified UI and API login flow
* ✅ JWT-based sessions
* ✅ Centralized authorization

> 🔧 **Switch authentication backends without changing your UI or API clients.**

---

## Next

👉 **Security & Permissions**
Learn how roles, permissions, and table-level access control apply to both local and LDAP-authenticated users.
