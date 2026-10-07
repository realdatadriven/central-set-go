
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
ATTACH IF NOT EXISTS '{{.Host}}/odata/ADMIN' AS admin (TYPE odata);
```

```sql apps
SELECT count(*) AS total, SUM(CASE WHEN excluded = 1 THEN 1 ELSE 0 END) AS excluded
FROM admin.app;
```

```sql ids
SELECT *
FROM admin.app;
```

<!--APP_ID: {{pluckJson .queries.ids "app_id"}}
APP_NAME: {{pluckJson .queries.ids "app"}}-->

<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="utf-8">
<meta name="viewport" content="width=device-width, initial-scale=1">
<title>Admin overview</title>
<link href="https://cdn.jsdelivr.net/npm/daisyui@5" rel="stylesheet" type="text/css">
<link href="https://cdn.jsdelivr.net/npm/daisyui@5/themes.css" rel="stylesheet" type="text/css">
<script src="https://cdn.jsdelivr.net/npm/@tailwindcss/browser@4"></script>
<script type="module" src="https://cdn.jsdelivr.net/npm/@weblogin/trendchart-elements/+esm"></script>
<style>
  tc-line, tc-bar, tc-pie, tc-stack {
    --shape-color: var(--color-primary);
    --area-color: var(--color-primary);
    --tooltip-background: var(--color-neutral);
    --tooltip-font-color: var(--color-neutral-content);
    --shape-color-1: var(--color-primary);  --shape-color-2: var(--color-secondary);
    --shape-color-3: var(--color-accent);   --shape-color-4: var(--color-info);
    --shape-color-5: var(--color-warning);  --shape-color-6: var(--color-neutral);
  }
  tc-line { --area-opacity: .14; width: 100%; height: 180px; }
  tc-bar  { width: 100%; height: 180px; }
  tc-pie  { width: 160px; height: 160px; }
  tc-line.spark { width: 6rem; height: 2.5rem; }
  .err    { --shape-color: var(--color-error); --area-color: var(--color-error); }
</style>
</head>
<body class="bg-base-200 min-h-screen">

<div class="navbar bg-base-100 shadow-sm px-4 lg:px-8">
  <div class="flex-1">
    <span class="text-lg font-semibold">Admin overview</span>
    <span id="mock-badge" class="badge badge-warning badge-sm ml-3 hidden">sample data</span>
  </div>
  <label class="flex items-center gap-2 text-sm">Dark
    <input type="checkbox" value="dim" class="toggle toggle-sm theme-controller">
  </label>
</div>

<main class="max-w-7xl mx-auto p-4 lg:p-8 space-y-6">
  <div id="alert" role="alert" class="alert alert-error hidden"></div>

  <!-- Users -->
  <div class="stats stats-vertical sm:stats-horizontal shadow w-full bg-base-100">
    <div class="stat">
      <div class="stat-figure">
        <tc-line class="spark" static data-chart="signups_30d" data-x="day" data-y="signups"></tc-line>
      </div>
      <div class="stat-title">Users</div>
      <div class="stat-value" data-tpl="{total|int}">-</div>
      <div class="stat-desc" data-tpl="{inactive|int} inactive">&nbsp;</div>
    </div>
    <div class="stat">
      <div class="stat-title">Active users</div>
      <div class="stat-value text-success" data-tpl="{active|int}">-</div>
      <div class="stat-desc" data-tpl="{with_failed_logins|int} with failed logins">&nbsp;</div>
    </div>
    <div class="stat">
      <div class="stat-title">Email confirmed</div>
      <div class="stat-value" data-tpl="{email_confirmed_pct}%">-</div>
      <div class="stat-desc">of all users</div>
    </div>
    <div class="stat">
      <div class="stat-title">Two-factor enabled</div>
      <div class="stat-value" data-tpl="{two_fa_pct}%">-</div>
      <div class="stat-desc">of all users</div>
    </div>
  </div>

  <!-- User log, last 24h -->
  <div class="stats stats-vertical sm:stats-horizontal shadow w-full bg-base-100">
    <div class="stat">
      <div class="stat-figure"><tc-line class="spark" static data-chart="requests_30d" data-x="day" data-y="requests"></tc-line></div>
      <div class="stat-title">Requests, last 24h</div>
      <div class="stat-value" data-tpl="{requests_24h|int}">-</div>
      <div class="stat-desc" data-tpl="{active_users_24h|int} distinct users">&nbsp;</div>
    </div>
    <div class="stat">
      <div class="stat-figure"><tc-line class="spark err" static data-chart="requests_30d" data-x="day" data-y="errors"></tc-line></div>
      <div class="stat-title">Error rate</div>
      <div class="stat-value text-error" data-tpl="{error_rate_pct}%">-</div>
      <div class="stat-desc" data-tpl="{errors_24h|int} failed requests">&nbsp;</div>
    </div>
    <div class="stat">
      <div class="stat-title">Avg response time</div>
      <div class="stat-value" data-tpl="{avg_ms|int} ms">-</div>
      <div class="stat-desc">request to response</div>
    </div>
  </div>

  <div class="grid gap-6 lg:grid-cols-3">
    <div class="card bg-base-100 shadow-sm lg:col-span-2" data-sql-for="signups_30d">
      <div class="card-body"><h2 class="card-title text-base">New users, last 30 days</h2>
        <tc-line data-chart="signups_30d" data-x="day" data-y="signups" tooltip-format="@L: @V new"></tc-line>
      </div>
    </div>
    <div class="card bg-base-100 shadow-sm" data-sql-for="users_by_role">
      <div class="card-body"><h2 class="card-title text-base">Users by role</h2>
        <div class="flex items-center gap-4">
          <tc-pie donut="28" data-chart="users_by_role" data-x="role" data-y="users"></tc-pie>
          <ul class="text-sm space-y-1 flex-1" data-legend="users_by_role" data-x="role" data-y="users"></ul>
        </div>
      </div>
    </div>
  </div>

  <div class="grid gap-6 lg:grid-cols-3">
    <div class="card bg-base-100 shadow-sm lg:col-span-2" data-sql-for="requests_30d">
      <div class="card-body"><h2 class="card-title text-base">Requests per day</h2>
        <tc-bar data-chart="requests_30d" data-x="day" data-y="requests" tooltip-format="@L: @V requests"></tc-bar>
        <h3 class="text-sm font-medium mt-2">Failed requests</h3>
        <tc-line class="err" style="height:70px" data-chart="requests_30d" data-x="day" data-y="errors" tooltip-format="@L: @V failed"></tc-line>
      </div>
    </div>
    <div class="card bg-base-100 shadow-sm" data-sql-for="top_actions">
      <div class="card-body"><h2 class="card-title text-base">Top actions, last 30 days</h2>
        <tc-bar horizontal gap="4" style="height:220px" data-chart="top_actions" data-x="action" data-y="calls" tooltip-format="@L: @V"></tc-bar>
        <div class="flex flex-wrap gap-1" data-badges="top_actions" data-x="action" data-y="calls"></div>
      </div>
    </div>
  </div>

  <div class="grid gap-6 lg:grid-cols-2">
    <div class="card bg-base-100 shadow-sm" data-sql-for="top_users">
      <div class="card-body"><h2 class="card-title text-base">Most active users, last 7 days</h2>
        <div class="overflow-x-auto"><table class="table table-sm table-zebra" data-table="top_users"></table></div>
      </div>
    </div>
    <div class="card bg-base-100 shadow-sm" data-sql-for="failed_logins">
      <div class="card-body"><h2 class="card-title text-base">Users with failed logins</h2>
        <div class="overflow-x-auto"><table class="table table-sm table-zebra" data-table="failed_logins"></table></div>
      </div>
    </div>
  </div>

  <div class="card bg-base-100 shadow-sm" data-sql-for="recent_logs">
    <div class="card-body"><h2 class="card-title text-base">Recent activity</h2>
      <div class="overflow-x-auto"><table class="table table-sm table-zebra" data-table="recent_logs"></table></div>
    </div>
  </div>
</main>

<!--
  SQL for each widget. Written for DuckDB (the central-set-go data app attaches the ADMIN db as `adm`;
  prefix tables with `adm.` if you don't `USE adm`). Block name = key in the dashboard.
  Assumption: user_log.res_type holds a value containing "err" or "fail" for failed requests - adjust is_error to your values.
-->
```sql kpi_users
SELECT count(*)                                              AS total,
       count(*) FILTER (WHERE active)                        AS active,
       count(*) FILTER (WHERE NOT coalesce(active, false))   AS inactive,
       round(100.0 * count(*) FILTER (WHERE email_confirmed) / nullif(count(*), 0), 1)  AS email_confirmed_pct,
       round(100.0 * count(*) FILTER (WHERE enable_2f_auth)  / nullif(count(*), 0), 1)  AS two_fa_pct,
       count(*) FILTER (WHERE coalesce(failed_login_attmpt, 0) > 0)                     AS with_failed_logins
FROM adim.users
WHERE coalesce(excluded, false) = false;
```

```sql kpi_logs
WITH l AS (
  SELECT *, (lower(coalesce(res_type, '')) LIKE '%err%' OR lower(coalesce(res_type, '')) LIKE '%fail%') AS is_error
  FROM admin.user_log
  WHERE req_at >= now() - INTERVAL 1 DAY AND coalesce(excluded, false) = false
)
SELECT count(*)                                                                   AS requests_24h,
       count(*) FILTER (WHERE is_error)                                           AS errors_24h,
       round(100.0 * count(*) FILTER (WHERE is_error) / nullif(count(*), 0), 1)   AS error_rate_pct,
       round(avg(date_diff('millisecond', req_at, res_at)), 0)                    AS avg_ms,
       count(DISTINCT user_id)                                                    AS active_users_24h
FROM l;
```

```sql signups_30d
WITH days AS (
  SELECT CAST(d AS DATE) AS day
  FROM generate_series(current_date - INTERVAL 29 DAY, current_date, INTERVAL 1 DAY) AS t(d)
)
SELECT strftime(days.day, '%Y-%m-%d') AS day, count(u.user_id) AS signups
FROM days
LEFT JOIN admin.users u ON CAST(u.created_at AS DATE) = days.day AND coalesce(u.excluded, false) = false
GROUP BY days.day
ORDER BY days.day;
```

```sql requests_30d
WITH days AS (
  SELECT CAST(d AS DATE) AS day
  FROM generate_series(current_date - INTERVAL 29 DAY, current_date, INTERVAL 1 DAY) AS t(d)
), l AS (
  SELECT CAST(req_at AS DATE) AS day,
         (lower(coalesce(res_type, '')) LIKE '%err%' OR lower(coalesce(res_type, '')) LIKE '%fail%') AS is_error
  FROM admin.user_log
  WHERE coalesce(excluded, false) = false
)
SELECT strftime(days.day, '%Y-%m-%d') AS day,
       count(l.day)                       AS requests,
       count(*) FILTER (WHERE l.is_error) AS errors
FROM days
LEFT JOIN l ON l.day = days.day
GROUP BY days.day
ORDER BY days.day;
```

```sql users_by_role
SELECT coalesce(r.role, 'no-role') AS role, count(*) AS users
FROM admin.users u
LEFT JOIN admin.role r ON r.role_id = u.role_id
WHERE coalesce(u.excluded, false) = false
GROUP BY 1
ORDER BY users DESC;
```

```sql top_actions
SELECT action, count(*) AS calls
FROM admin.user_log
WHERE req_at >= now() - INTERVAL 30 DAY AND coalesce(excluded, false) = false
GROUP BY action
ORDER BY calls DESC
LIMIT 8;
```

```sql top_users
SELECT u.username,
       count(*) AS requests,
       count(*) FILTER (WHERE lower(coalesce(l.res_type, '')) LIKE '%err%'
                           OR lower(coalesce(l.res_type, '')) LIKE '%fail%') AS errors,
       strftime(max(l.req_at), '%Y-%m-%d %H:%M') AS last_seen
FROM admin.user_log l
JOIN admin.users u ON u.user_id = l.user_id
WHERE l.req_at >= now() - INTERVAL 7 DAY AND coalesce(l.excluded, false) = false
GROUP BY u.username
ORDER BY requests DESC
LIMIT 8;
```

```sql failed_logins
SELECT username,
       email,
       failed_login_attmpt AS attempts,
       strftime(last_failed_login, '%Y-%m-%d %H:%M') AS last_failed
FROM admin.users
WHERE coalesce(failed_login_attmpt, 0) > 0 AND coalesce(excluded, false) = false
ORDER BY failed_login_attmpt DESC
LIMIT 8;
```

```sql recent_logs
SELECT strftime(l.req_at, '%Y-%m-%d %H:%M:%S') AS at,
       u.username,
       l.action,
       l.res_type AS result,
       l."table",
       l.req_ip AS ip,
       l.res_msg AS message
FROM admin.user_log l
LEFT JOIN admin.users u ON u.user_id = l.user_id
WHERE coalesce(l.excluded, false) = false
ORDER BY l.req_at DESC
LIMIT 15;
```

</body>
</html>