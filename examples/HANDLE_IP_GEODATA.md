# PARSE_ACC_JSON_RESP
```yaml
name: PARSE_ACC_JSON_RESP
runs_as: SCRIPTS
description: Handle JSON response from C7 ACtion example of integration with ETLX
active: true
```

## JSON_RESP
```yaml
name: JSON_RESP
description: Handle JSON response from C7 ACtion example of integration with ETLX
connection: "duckdb:"
script_sql:
  - LOAD SQLITE
  - ATTACH 'database/ADMIN.db' (TYPE SQLITE)
  - json_response_parse
  - DETACH ADMIN
active: true
```

```sql
-- json_response_parse
CREATE OR REPLACE TABLE test_api_response_parse AS
WITH _json AS (SELECT '{{.api_response}}'::JSON AS j)
SELECT j->>'ip'      AS ip,
    j->>'city'       AS city,
    j->>'region'     AS region,
    j->>'country'    AS country,
    j->>'loc'        AS loc,
    j->>'org'        AS org,
    j->>'timezone'   AS timezone,
    j->>'readme'     AS readme
FROM _json;
```