<!DOCTYPE html>

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
APP_ID: {{pluckJson .queries.ids "app_id"}}
APP_NAME: {{pluckJson .queries.ids "app"}}

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
      <div class="stat-figure"><tc-line class="spark" static data-chart="signups_30d" data-x="day" data-y="signups"></tc-line></div>
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
FROM users
WHERE coalesce(excluded, false) = false;
```

```sql kpi_logs
WITH l AS (
  SELECT *, (lower(coalesce(res_type, '')) LIKE '%err%' OR lower(coalesce(res_type, '')) LIKE '%fail%') AS is_error
  FROM user_log
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
LEFT JOIN users u ON CAST(u.created_at AS DATE) = days.day AND coalesce(u.excluded, false) = false
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
  FROM user_log
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
FROM users u
LEFT JOIN role r ON r.role_id = u.role_id
WHERE coalesce(u.excluded, false) = false
GROUP BY 1
ORDER BY users DESC;
```

```sql top_actions
SELECT action, count(*) AS calls
FROM user_log
WHERE req_at >= now() - INTERVAL 30 DAY AND coalesce(excluded, false) = false
GROUP BY action
ORDER BY calls DESC
LIMIT 8;
```

```sql top_users
SELECT u.username,
       count(*)                                                            AS requests,
       count(*) FILTER (WHERE lower(coalesce(l.res_type, '')) LIKE '%err%'
                           OR lower(coalesce(l.res_type, '')) LIKE '%fail%') AS errors,
       strftime(max(l.req_at), '%Y-%m-%d %H:%M')                           AS last_seen
FROM user_log l
JOIN users u ON u.user_id = l.user_id
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
FROM users
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
FROM user_log l
LEFT JOIN users u ON u.user_id = l.user_id
WHERE coalesce(l.excluded, false) = false
ORDER BY l.req_at DESC
LIMIT 15;
```

<script type="module">
// ---- Data source -----------------------------------------------------------
// Leave endpoint empty to render sample data. Otherwise the page POSTs
// {name, sql} and expects a JSON array of row objects (or {data: [...]}).
const CONFIG = { endpoint: '', headers: {} };

// ---- Parse the ```sql <name> blocks from the page body --------------------
const SQL = {};
document.getElementById('queries').textContent
  .replace(/```sql\s+(\w+)\n([\s\S]*?)```/g, (_, n, s) => { SQL[n] = s.trim(); });

// ---- Sample data ------------------------------------------------------------
const rnd = (a, b) => Math.round(a + Math.random() * (b - a));
const days = n => Array.from({ length: n }, (_, i) => new Date(Date.now() - (n - 1 - i) * 864e5).toISOString().slice(0, 10));
const MOCK = {
  kpi_users: () => [{ total: 128, active: 117, inactive: 11, email_confirmed_pct: 75, two_fa_pct: 32, with_failed_logins: 5 }],
  kpi_logs: () => [{ requests_24h: 2184, errors_24h: 37, error_rate_pct: 1.7, avg_ms: 142, active_users_24h: 34 }],
  signups_30d: () => days(30).map(day => ({ day, signups: rnd(0, 9) })),
  requests_30d: () => days(30).map(day => ({ day, requests: rnd(900, 2600), errors: rnd(5, 70) })),
  users_by_role: () => [{ role: 'tenant', users: 88 }, { role: 'no-role', users: 24 }, { role: 'root', users: 3 }, { role: 'anonymous', users: 1 }],
  top_actions: () => ['read', 'update', 'create', 'login', 'delete', 'export', 'odata', 'logout'].map((action, i) => ({ action, calls: 1800 - i * 210 })),
  top_users: () => ['root', 'ana', 'joao', 'maria', 'sam'].map((username, i) => ({ username, requests: 900 - i * 140, errors: rnd(1, 20), last_seen: '2026-10-06 09:' + (10 + i * 7) })),
  failed_logins: () => [{ username: 'sam', email: 'sam@domain.com', attempts: 4, last_failed: '2026-10-06 08:12' }, { username: 'ana', email: 'ana@domain.com', attempts: 2, last_failed: '2026-10-05 17:40' }],
  recent_logs: () => Array.from({ length: 8 }, (_, i) => ({ at: '2026-10-06 09:' + (50 - i * 4), username: ['root', 'ana', 'joao'][i % 3], action: ['read', 'update', 'login'][i % 3], result: i === 3 ? 'error' : 'success', table: ['users', 'role', 'app'][i % 3], ip: '10.0.0.' + (i + 4), message: i === 3 ? 'validation USR01 failed' : 'ok' }))
};
async function run(name) {
  if (!CONFIG.endpoint) return MOCK[name]();
  const r = await fetch(CONFIG.endpoint, { method: 'POST', headers: { 'Content-Type': 'application/json', ...CONFIG.headers }, body: JSON.stringify({ name, sql: SQL[name] }) });
  if (!r.ok) throw new Error(name + ': HTTP ' + r.status);
  const j = await r.json();
  return Array.isArray(j) ? j : j.data;
}

// ---- Render -----------------------------------------------------------------
const $$ = s => document.querySelectorAll(s);
const F = { int: v => Number(v ?? 0).toLocaleString(), pct: v => Number(v ?? 0).toFixed(1) + '%' };
const PAL = ['primary', 'secondary', 'accent', 'info', 'warning', 'neutral'];
const D = {};
if (!CONFIG.endpoint) document.getElementById('mock-badge').classList.remove('hidden');

try {
  await Promise.all(Object.keys(SQL).map(async n => { D[n] = await run(n); }));
} catch (e) {
  const a = document.getElementById('alert');
  a.textContent = 'Could not load data: ' + e.message;
  a.classList.remove('hidden');
}

// stats: data-tpl="{field|int}" reads from the query named by the closest stat group
const STAT_SRC = { total: 'kpi_users', inactive: 'kpi_users', active: 'kpi_users', with_failed_logins: 'kpi_users', email_confirmed_pct: 'kpi_users', two_fa_pct: 'kpi_users' };
$$('[data-tpl]').forEach(el => {
  el.textContent = el.dataset.tpl.replace(/\{(\w+)(?:\|(\w+))?\}/g, (_, f, fm) => {
    const row = (D[STAT_SRC[f] || 'kpi_logs'] || [])[0] || {};
    return fm ? F[fm](row[f]) : (row[f] ?? '-');
  });
});

// charts
$$('[data-chart]').forEach(el => {
  const rows = D[el.dataset.chart] || [];
  el.setAttribute('values', JSON.stringify(rows.map(r => +r[el.dataset.y])));
  el.setAttribute('labels', JSON.stringify(rows.map(r => String(r[el.dataset.x]))));
});

// legend + badges
$$('[data-legend]').forEach(el => {
  el.innerHTML = '';
  (D[el.dataset.legend] || []).forEach((r, i) => {
    const li = document.createElement('li');
    li.className = 'flex items-center gap-2';
    li.innerHTML = '<span class="size-2.5 rounded-full" style="background:var(--color-' + PAL[i % 6] + ')"></span><span></span><span class="ml-auto font-medium"></span>';
    li.children[1].textContent = r[el.dataset.x];
    li.children[2].textContent = F.int(r[el.dataset.y]);
    el.appendChild(li);
  });
});
$$('[data-badges]').forEach(el => (D[el.dataset.badges] || []).forEach(r => {
  const b = document.createElement('span');
  b.className = 'badge badge-ghost badge-sm';
  b.textContent = r[el.dataset.x] + ' ' + F.int(r[el.dataset.y]);
  el.appendChild(b);
}));

// tables
$$('[data-table]').forEach(t => {
  const rows = D[t.dataset.table] || [];
  if (!rows.length) { t.innerHTML = '<tbody><tr><td class="text-base-content/60">Nothing to show yet.</td></tr></tbody>'; return; }
  const cols = Object.keys(rows[0]);
  const thead = t.createTHead().insertRow();
  cols.forEach(c => { const th = document.createElement('th'); th.textContent = c.replace(/_/g, ' '); thead.appendChild(th); });
  const tb = t.createTBody();
  rows.forEach(r => {
    const tr = tb.insertRow();
    cols.forEach(c => {
      const td = tr.insertCell();
      td.textContent = r[c] ?? '-';
      if (c === 'result') td.className = /err|fail/i.test(r[c]) ? 'text-error' : 'text-success';
    });
  });
});

// show each widget's SQL under it
$$('[data-sql-for]').forEach(card => {
  (card.dataset.sqlFor).split(',').forEach(n => {
    const d = document.createElement('details');
    d.className = 'collapse collapse-arrow border border-base-300 mt-2';
    d.innerHTML = '<summary class="collapse-title text-xs min-h-0 py-2"></summary><div class="collapse-content"><pre class="text-xs overflow-x-auto"></pre></div>';
    d.querySelector('summary').textContent = 'SQL: ' + n;
    d.querySelector('pre').textContent = SQL[n] || '';
    (card.querySelector('.card-body') || card).appendChild(d);
  });
});
</script>
</body>
</html>