---
weight: 7050
title: "Users, Roles & Permissions"
description: "Create users, assign roles, and control application, table, column, and row access."
icon: manage_accounts
date: 2025-12-16T01:04:15+00:00
lastmod: 2026-11-14
draft: false
images: []
---

Central Set separates **authentication** (who signed in) from **authorization** (what that identity can access). Users authenticate with the configured login provider; roles and permission records in the `ADMIN` database determine which applications, menus, operations, fields, and rows they can use.

This page covers the standard ADMIN user and role model. For login, tokens, password management, and LDAP, see [Authentication & Access](../auth/). For the underlying ADMIN tables, see [Admin Schema](../admin-schema/) and the [ADMIN model](../admin-model/).

## The access model

An authenticated user has a user record and a role context. The `users.role_id` field stores the user's primary role. The `user_role` table can associate the user with additional roles. Central Set evaluates the role in the authenticated user context together with the user's role assignments when it loads permissions.

Roles are reusable permission profiles. Assigning a role to a person gives that person the access configured for the role; it does not copy permission settings into the user record.

| ADMIN table | What it controls |
| --- | --- |
| `users` | Login identity, profile, active state, language, and primary `role_id`. |
| `role` | Named roles such as `root`, `tenant`, or roles created for your organization. |
| `user_role` | Additional user-to-role assignments. |
| `role_app` | Whether a role can access an application. |
| `role_app_menu` | Whether a role can access a menu in an application. |
| `role_app_menu_table` | Table-level `create`, `read`, `update`, `delete`, and `share` permissions. |
| `column_level_access` | Per-user column permissions for `create`, `read`, and `update`. |
| `row_level_access` | Per-user permissions on individual rows. |
| `role_row_level_access` | Per-role permissions on individual rows. |

The access layers can be thought of as a path through the application:

```text
User → Role → Application → Menu → Table operations → Columns / Rows
```

Grant only the layers the user needs. For example, a table permission does not by itself grant the role access to an application or menu that the role cannot reach.

## Create a role and assign it to users

Use the **ADMIN** application to manage users and roles. The exact menu labels depend on the ADMIN model enabled for your installation; the relevant records are in the `role`, `users`, and `user_role` tables.

1. **Create a role** with a clear name and description. Prefer roles that describe a responsibility, such as `sales-reader` or `support-agent`, rather than a person's name.
2. **Configure the role's permissions** using the permissions/key action on the role. This opens the application permission editor when the ADMIN model's `AdminApps` component is enabled.
3. **Create or select a user** in the Users screen. Set the user's primary role in `users.role_id`.
4. **Add any additional role assignments** in the user's `user_role` subform when that person needs more than the primary role.
5. **Save and test** with a normal account that has the role. Check both the visible navigation and the operations the account can perform.

An assignment can be removed without deleting the user or role. ADMIN records use an `excluded` flag for soft exclusion; when changing assignments, use the ADMIN UI or your normal model-management process so the record state and audit metadata remain consistent.

## Configure role access

The role permission editor applies access from broad to specific. Configure each layer for the target application.

### 1. Application access

Grant the role access to the application. This is represented by a `role_app` record with `access` enabled. Without application access, the application is not available to that role in the application list.

### 2. Menu access

For each permitted application, grant access to the menus the role should see. These grants are stored in `role_app_menu`. A role may have access to an application but still be limited to selected menus.

### 3. Table operations

For each table exposed under a menu, set the operations the role needs in `role_app_menu_table`:

| Permission | Allows |
| --- | --- |
| `read` | Read/list records through the configured Central Set data path. |
| `create` | Create records. |
| `update` | Change existing records. |
| `delete` | Delete records. |
| `share` | Use the configured sharing permission for the table or row. |

Granting `read` or `create` also makes the table eligible to appear in the role's menu data. Keep the table operation settings narrow: a user who only needs to review data should not receive create, update, or delete access.

### Example: read-only role

For a role that should review sales orders in the `SALES` application:

1. Grant access to the `SALES` application.
2. Grant access to the menu that contains orders.
3. Enable `read` for the orders table.
4. Leave `create`, `update`, and `delete` disabled unless the role needs those operations.
5. Test that the user can read permitted orders and cannot create, change, or delete them.

The example describes the intended permission setup; use the actual app, menu, and table names registered in your ADMIN database.

## Declare role permissions in a model

For repeatable deployments, define role access in the model instead of entering every permission by hand. The ADMIN model includes a disabled [`ROLE_ACCESS` example](https://github.com/realdatadriven/central-set-go/blob/main/assets/models/admin_model.md#L2051-L2072) and a disabled [`ROLE_USERS` example](https://github.com/realdatadriven/central-set-go/blob/main/assets/models/admin_model.md#L2074-L2090). These use ETLX model modes to create or update ADMIN permission records.

### `ROLE_ACCESS`: create or update a role's access

The header identifies the role-access declaration and its ADMIN connection. The `database` value identifies the database ETLX uses when it must look up rows for RLA entries.

````markdown {linenos=table}
# ROLE_ACCESS
```yaml
name: ROLE_ACCESS
description: Role Model Example, Create Role anonymous and gives access to ADMIN/Arrow Flight ...
database: ADMIN #RLA may need to access records on other DBs
runs_as: ROLE
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
active: true
```

## ANONYMOUS
```yaml
name: anonymous
description: Anonymous Role
access:
  - ADMIN:
    - Arrow Flight:
      - {table: flight_catalog, read: true, rla: [{flight_catalog: admin, read: true, share: true}]}
      - flight_schema
      - flight_schema_table
active: true
```
````

Each named entry below the header (`# ROLE_ACCESS`) defines a role. Its `access` list nests application, menu, and table names. This example grants the `anonymous` role read access to `flight_catalog` and read access to a particular catalog row:

The declaration creates the role if it does not exist, then creates or updates its application, menu, and table permission records. For a table mapping, omitted CRUD flags default to `false`. A bare table name such as `- flight_schema` is shorthand for granting **all five** table operations: create, read, update, delete, and share. Use a mapping with explicit flags when you want narrower access.

For example, this table mapping grants read and update but does not grant create, delete, or share:

```yaml {linenos=table}
...
access:
  - SALES:
      - Orders:
          - table: orders
              create: false
              read: true
              update: true
              delete: false
              share: false
```

**Update caveat:** this project currently pins ETLX v1.200.10. In that version, when a `ROLE_ACCESS` declaration reaches a table permission record that already exists, the runner merges the stored record over the declared CRUD flags before updating it. As a result, changing these flags in the model may leave the existing table permissions unchanged. Use the ADMIN permission editor to change existing table CRUD grants, and verify the effective permissions after applying a declaration. The RLA flags follow a separate update path.

An `rla` entry identifies a row by a field and value, then assigns row-level flags for the matched record. In the example, ETLX looks up the `flight_catalog` row whose `flight_catalog` field equals `admin`, then adds or updates the role's RLA record for that row. Choose a field value that identifies the intended record unambiguously. The target menu-table must already have `requires_rla: true`; the `ROLE_ACCESS` runner checks this and fails if RLA is not enabled there. `database` must point to the database containing the row to look up, while `admin_conn` points to the ADMIN database where the role and permission records are stored.

The runner processes the records described by the declaration; it does not remove grants simply because they are omitted from a later version. Review existing grants and explicitly revoke obsolete access through ADMIN before relying on a reduced declaration.

### `ROLE_USERS`: assign users to a role

Use `ROLE_USERS` to declare which existing users should receive a role. The example can refer to each user by username or email address:

````markdown {linenos=table}
# ROLE_USERS
```yaml
name: ROLE_USERS
description: Give user access to a role
runs_as: ROLE_USERS
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
active: true
```

## ANONYMOUS
```yaml
name: anonymous
description: Anonymous Role
users: [root, admin@domain.com, anonymous.user]
active: true
```
````

The runner creates the role if needed and creates or updates `user_role` assignments for the listed users. User records must already exist in ADMIN. This declaration ensures assignments for listed users; it does not remove assignments for users omitted from the list. Remove unwanted memberships explicitly in the ADMIN UI or through a reviewed model change.

Both headers in [`admin_model.md`](https://github.com/realdatadriven/central-set-go/blob/main/assets/models/admin_model.md#L2051-L2090) currently have `active: false`, so they are examples and are disabled as checked in. Set the relevant header to `active: true` only when you intend the model runner to apply it. Review the target database, user identifiers, table shorthands, and resulting access before enabling the declarations in a deployed environment.

## Column-level permissions

Use `column_level_access` when a user needs different permissions for particular fields of a table. A record is scoped to a user, application, table, and column, with `create`, `read`, and `update` flags.

Column-level access is more specific than table CRUD access. Configure both levels deliberately: first grant the table operation the user needs, then configure field access where the application requires additional restrictions. Column rules are not a substitute for excluding sensitive fields from responses or for enforcing permissions in custom SQL code.

## Row-level access (RLA)

Row-level access narrows permissions to particular records within an otherwise accessible table. The ADMIN model supports both user-specific records in `row_level_access` and role-based records in `role_row_level_access`. These records identify the application, table, and row, and can grant `read`, `update`, `delete`, or `share` access.

To use RLA:

1. Register the target table in the application metadata.
2. Mark the table or its menu-table entry with `requires_rla: true` when row-level access is required.
3. Add user- or role-specific row permissions for the relevant row IDs and operations.
4. Test requests for a permitted row and a row with no grant using the actual UI and API paths used by the application.

The Users and Roles permission controls can open row-access management when the model's `AdminApps` component is enabled. Row rules are scoped to application and table; a row ID by itself is not a complete permission scope.

> Row-level filtering depends on the request going through Central Set's access-aware data path. If custom code connects directly to a database or runs its own SQL, that code must enforce the same restrictions itself. For Dynamic UI pages, see the [generic CRUD endpoint guidance](../dynamic-ui/#generic-crud-endpoints).

## How permissions are applied

The authenticated user context is used by the server to load the user's role IDs and evaluate application, menu, and table permissions. The ADMIN model also includes the `anonymous` role in access evaluation for anonymous contexts. A role or assignment marked excluded is not treated as an ordinary active permission record.

The UI uses these permissions to decide what to display, while Central Set's protected API and CRUD paths also perform server-side access checks. Hiding a menu item alone is not an authorization boundary. Always test the API operations the user's client can call, not only what appears in the sidebar.

The seeded `root` role is highly privileged. Reserve it for trusted administrators and use purpose-specific roles for routine work. The role IDs and seed data can differ when an installation customizes or updates the ADMIN model; do not assume a role ID is portable across environments.

## Test a permission change

After changing a role or user assignment, verify the result with a separate account that has the intended permissions:

1. Sign in as that account and confirm that only the expected applications and menus are available.
2. Try each required table operation: read, create, update, and delete.
3. Try a restricted operation directly through the application's API or CRUD client, not only through the UI.
4. If RLA is enabled, compare a row the user is allowed to access with one they are not.
5. If column permissions are configured, inspect both the displayed data and permitted writes for the restricted fields.

Repeat these checks when a user has multiple roles, since the effective access context includes their assigned roles. Validate the real requests and data path used by each application.

## Configuration and related guides

The ADMIN model defines these tables and the role/user permission controls in [`assets/models/admin_model.md`](../../../../assets/models/admin_model.md). In particular, review the `role`, `users`, `user_role`, `role_app`, `role_app_menu`, `role_app_menu_table`, `column_level_access`, `row_level_access`, and `role_row_level_access` definitions.

Feature availability can be limited by server configuration. The [configuration guide](config/) documents `ENABLE_APP`, `ENABLE_USER`, and `CORE_TABLES`; check these values if a permission table or management screen is unavailable.

See also:

- [Authentication & Access](../auth/) for login, tokens, password management, and LDAP.
- [Applications](../apps/) for application and menu structure.
- [Admin Schema](../admin-schema/) for the broader ADMIN database model.
- [Dynamic UI](../dynamic-ui/) for authenticated data access and custom UI behavior.
