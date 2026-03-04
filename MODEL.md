# ADMMIN_MODEL
```yaml
name: ADMIN
description: CS ADMIN Model
runs_as: MODEL
conn: 'sqlite3:database/ADMIN.db'
active: true
```

## LANGS
```yaml
table: lang
comment: Languages
columns:
  lang_id:     { type: integer, pk: true, autoincrement: true, comment: "Lang ID" }
  lang:        { type: varchar(4), unique: true, nullable: false, comment: "Language" }
  lang_desc:   { type: varchar(200), comment: "Description" }
  created_at:  { type: datetime, comment: "Created at" }
  updated_at:  { type: datetime, comment: "Updated at" }
  excluded:    { type: boolean, default: false, comment: "Excluded" }
data:
  - {lang_id: 1, lang: en, lang_desc: English, excluded: false}
```