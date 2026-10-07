
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
FROM admin.users
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
<!--APP_ID: {{pluckJson .queries.ids "app_id"}}
APP_NAME: {{pluckJson .queries.ids "app"}}-->

{{- $q := .queries -}}
{{- $pal := list "primary" "secondary" "accent" "info" "warning" "neutral" -}}
{{- $u := dict -}}{{- if kindIs "slice" $q.kpi_users -}}{{- with first $q.kpi_users -}}{{- $u = . -}}{{- end -}}{{- end -}}
{{- $l := dict -}}{{- if kindIs "slice" $q.kpi_logs -}}{{- with first $q.kpi_logs -}}{{- $l = . -}}{{- end -}}{{- end -}}
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
  <div class="flex-1"><span class="text-lg font-semibold">Admin overview</span></div>
  <label class="flex items-center gap-2 text-sm">Dark
    <input type="checkbox" value="dim" class="toggle toggle-sm theme-controller">
  </label>
</div>

<main class="max-w-7xl mx-auto p-4 lg:p-8 space-y-6">

  {{- /* one alert per failed query */ -}}
  {{range $name, $v := $q}}{{if kindIs "string" $v}}
  <div role="alert" class="alert alert-error"><span><b>{{$name}}</b>: {{$v}}</span></div>
  {{end}}{{end}}

  <!-- Users -->
  <div class="stats stats-vertical sm:stats-horizontal shadow w-full bg-base-100">
    <div class="stat">
      <div class="stat-figure"><tc-line class="spark" static values='{{pluckJson $q.signups_30d "signups"}}' labels='{{pluckJson $q.signups_30d "day"}}'></tc-line></div>
      <div class="stat-title">Users</div>
      <div class="stat-value">{{$u.total | default 0}}</div>
      <div class="stat-desc">{{$u.inactive | default 0}} inactive</div>
    </div>
    <div class="stat">
      <div class="stat-title">Active users</div>
      <div class="stat-value text-success">{{$u.active | default 0}}</div>
      <div class="stat-desc">{{$u.with_failed_logins | default 0}} with failed logins</div>
    </div>
    <div class="stat">
      <div class="stat-title">Email confirmed</div>
      <div class="stat-value">{{printf "%.1f" (float64 ($u.email_confirmed_pct | default 0))}}%</div>
      <div class="stat-desc">of all users</div>
    </div>
    <div class="stat">
      <div class="stat-title">Two-factor enabled</div>
      <div class="stat-value">{{printf "%.1f" (float64 ($u.two_fa_pct | default 0))}}%</div>
      <div class="stat-desc">of all users</div>
    </div>
  </div>

  <!-- User log, last 24h -->
  <div class="stats stats-vertical sm:stats-horizontal shadow w-full bg-base-100">
    <div class="stat">
      <div class="stat-figure"><tc-line class="spark" static values='{{pluckJson $q.requests_30d "requests"}}' labels='{{pluckJson $q.requests_30d "day"}}'></tc-line></div>
      <div class="stat-title">Requests, last 24h</div>
      <div class="stat-value">{{$l.requests_24h | default 0}}</div>
      <div class="stat-desc">{{$l.active_users_24h | default 0}} distinct users</div>
    </div>
    <div class="stat">
      <div class="stat-figure"><tc-line class="spark err" static values='{{pluckJson $q.requests_30d "errors"}}' labels='{{pluckJson $q.requests_30d "day"}}'></tc-line></div>
      <div class="stat-title">Error rate</div>
      <div class="stat-value text-error">{{printf "%.1f" (float64 ($l.error_rate_pct | default 0))}}%</div>
      <div class="stat-desc">{{$l.errors_24h | default 0}} failed requests</div>
    </div>
    <div class="stat">
      <div class="stat-title">Avg response time</div>
      <div class="stat-value">{{$l.avg_ms | default 0 | int64}} ms</div>
      <div class="stat-desc">request to response</div>
    </div>
  </div>

  <div class="grid gap-6 lg:grid-cols-3">
    <div class="card bg-base-100 shadow-sm lg:col-span-2">
      <div class="card-body"><h2 class="card-title text-base">New users, last 30 days</h2>
        <tc-line values='{{pluckJson $q.signups_30d "signups"}}' labels='{{pluckJson $q.signups_30d "day"}}' tooltip-format="@L: @V new"></tc-line>
      </div>
    </div>
    <div class="card bg-base-100 shadow-sm">
      <div class="card-body"><h2 class="card-title text-base">Users by role</h2>
        <div class="flex items-center gap-4">
          <tc-pie donut="28" values='{{pluckJson $q.users_by_role "users"}}' labels='{{pluckJson $q.users_by_role "role"}}'></tc-pie>
          <ul class="text-sm space-y-1 flex-1">
            {{- if kindIs "slice" $q.users_by_role}}{{range $i, $r := $q.users_by_role}}
            <li class="flex items-center gap-2">
              <span class="size-2.5 rounded-full" style="background: var(--color-{{index $pal (mod $i 6)}})"></span>
              <span>{{$r.role}}</span><span class="ml-auto font-medium">{{$r.users}}</span>
            </li>
            {{- end}}{{end}}
          </ul>
        </div>
      </div>
    </div>
  </div>

  <div class="grid gap-6 lg:grid-cols-3">
    <div class="card bg-base-100 shadow-sm lg:col-span-2">
      <div class="card-body"><h2 class="card-title text-base">Requests per day</h2>
        <tc-bar values='{{pluckJson $q.requests_30d "requests"}}' labels='{{pluckJson $q.requests_30d "day"}}' tooltip-format="@L: @V requests"></tc-bar>
        <h3 class="text-sm font-medium mt-2">Failed requests</h3>
        <tc-line class="err" style="height:70px" values='{{pluckJson $q.requests_30d "errors"}}' labels='{{pluckJson $q.requests_30d "day"}}' tooltip-format="@L: @V failed"></tc-line>
      </div>
    </div>
    <div class="card bg-base-100 shadow-sm">
      <div class="card-body"><h2 class="card-title text-base">Top actions, last 30 days</h2>
        <tc-bar horizontal gap="4" style="height:220px" values='{{pluckJson $q.top_actions "calls"}}' labels='{{pluckJson $q.top_actions "action"}}' tooltip-format="@L: @V"></tc-bar>
        <div class="flex flex-wrap gap-1">
          {{- if kindIs "slice" $q.top_actions}}{{range $q.top_actions}}
          <span class="badge badge-ghost badge-sm">{{.action}} {{.calls}}</span>
          {{- end}}{{end}}
        </div>
      </div>
    </div>
  </div>

  <div class="grid gap-6 lg:grid-cols-2">
    <div class="card bg-base-100 shadow-sm">
      <div class="card-body"><h2 class="card-title text-base">Most active users, last 7 days</h2>
        <div class="overflow-x-auto"><table class="table table-sm table-zebra">
          <thead><tr><th>username</th><th>requests</th><th>errors</th><th>last seen</th></tr></thead>
          <tbody>
            {{- if kindIs "slice" $q.top_users}}{{range $q.top_users}}
            <tr><td>{{.username}}</td><td>{{.requests}}</td><td>{{.errors}}</td><td>{{default "-" .last_seen}}</td></tr>
            {{- else}}<tr><td colspan="4" class="text-base-content/60">Nothing to show yet.</td></tr>{{end}}{{end}}
          </tbody>
        </table></div>
      </div>
    </div>
    <div class="card bg-base-100 shadow-sm">
      <div class="card-body"><h2 class="card-title text-base">Users with failed logins</h2>
        <div class="overflow-x-auto"><table class="table table-sm table-zebra">
          <thead><tr><th>username</th><th>email</th><th>attempts</th><th>last failed</th></tr></thead>
          <tbody>
            {{- if kindIs "slice" $q.failed_logins}}{{range $q.failed_logins}}
            <tr><td>{{.username}}</td><td>{{default "-" .email}}</td><td>{{.attempts}}</td><td>{{default "-" .last_failed}}</td></tr>
            {{- else}}<tr><td colspan="4" class="text-base-content/60">Nothing to show yet.</td></tr>{{end}}{{end}}
          </tbody>
        </table></div>
      </div>
    </div>
  </div>

  <div class="card bg-base-100 shadow-sm">
    <div class="card-body"><h2 class="card-title text-base">Recent activity</h2>
      <div class="overflow-x-auto"><table class="table table-sm table-zebra">
        <thead><tr><th>at</th><th>username</th><th>action</th><th>result</th><th>table</th><th>ip</th><th>message</th></tr></thead>
        <tbody>
          {{- if kindIs "slice" $q.recent_logs}}{{range $q.recent_logs}}
          <tr>
            <td>{{default "-" .at}}</td><td>{{default "-" .username}}</td><td>{{.action}}</td>
            <td class="{{if regexMatch "(?i)err|fail" (toString .result)}}text-error{{else}}text-success{{end}}">{{default "-" .result}}</td>
            <td>{{default "-" .table}}</td><td>{{default "-" .ip}}</td><td>{{default "-" .message}}</td>
          </tr>
          {{- else}}<tr><td colspan="7" class="text-base-content/60">Nothing to show yet.</td></tr>{{end}}{{end}}
        </tbody>
      </table></div>
    </div>
  </div>
</main>
</body>
</html>


