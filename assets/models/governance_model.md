# GOVERNANCE_MODEL
```yaml
name: GOVERNANCE
description: Metadata and data governance model covering domains, glossary terms, data sources, assets, fields, classifications, and quality rules
runs_as: MODEL
admin_conn: '@DB_DRIVER_NAME:@DB_DSN'
create_all: checkfirst
_drop_all: checkfirst
update_table_metadata: true
active: true
cs_app:
  Governance:
    menu_icon: document-check
    menu_order: 1
    active: true
    tables:
      - stakeholders
      - business_units
      - domains
      - glossary_terms
      - data_sources
      - asset_schemas
      - data_assets
      - asset_fields
  Classifications:
    menu_icon: tag
    menu_order: 2
    active: true
    tables:
      - tag_categories
      - tags
      - asset_term_mappings
      - field_term_mappings
      - asset_tag_mappings
      - field_tag_mappings
  Quality:
    menu_icon: chart-bar
    menu_order: 3
    active: true
    tables:
      - data_quality_rules
      - data_quality_executions
```

## STAKEHOLDERS
```yaml
table: stakeholders
comment: Stakeholder
tooltip: Central repository of individuals who act as technical owners, business stewards, or data producers.
columns:
  stakeholder_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the stakeholder.", form_display: true, table_display: true, order: 1 }
  email: { type: varchar, len: 255, nullable: false, unique: true, comment: "Email", tooltip: "Primary corporate email address of the stakeholder.", form_display: true, table_display: true, order: 2 }
  full_name: { type: varchar, len: 150, nullable: false, comment: "Name", tooltip: "Full name of the stakeholder.", form_display: true, table_display: true, order: 3 }
  title: { type: varchar, len: 100, comment: "Title", tooltip: "Job title or role description of the stakeholder.", form_display: true, table_display: true, order: 4 }
  slack_handle: { type: varchar, len: 50, comment: "Slack", tooltip: "Internal Slack username or handle for rapid communication.", form_display: true, table_display: true, order: 5 }
  created_at: { type: datetime, comment: "Created", tooltip: "Timestamp when the stakeholder record was created.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Updated", tooltip: "Timestamp when the stakeholder record was last updated.", form_display: false, table_display: true }
```

## BUSINESS_UNITS
```yaml
table: business_units
comment: Business Unit
tooltip: High-level organizational divisions such as Finance, Risk, Marketing, or Product.
columns:
  business_unit_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the business unit.", form_display: true, table_display: true, order: 1 }
  name: { type: varchar, len: 100, nullable: false, unique: true, comment: "Name", tooltip: "Name of the business unit.", form_display: true, table_display: true, order: 2 }
  description: { type: text, comment: "Description", tooltip: "Detailed description of the business unit scope and responsibilities.", form_display: true, table_display: true, order: 3 }
  created_at: { type: datetime, comment: "Created", tooltip: "Timestamp when the business unit was recorded.", form_display: false, table_display: true }
```

## DOMAINS
```yaml
table: domains
comment: Domain
tooltip: Data Mesh domains representing distinct business bounded contexts such as a Customer Domain or Revenue Domain.
columns:
  domain_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the domain.", form_display: true, table_display: true, order: 1 }
  business_unit_id: { type: integer, fk: "business_units.business_unit_id", comment: "Business Unit", tooltip: "Foreign key referencing the parent business unit.", form_display: true, table_display: true, order: 2 }
  name: { type: varchar, len: 100, nullable: false, unique: true, comment: "Name", tooltip: "Name of the domain.", form_display: true, table_display: true, order: 3 }
  description: { type: text, comment: "Description", tooltip: "Detailed description of the data domain boundary.", form_display: true, table_display: true, order: 4 }
  domain_lead_id: { type: integer, fk: "stakeholders.stakeholder_id", comment: "Lead", tooltip: "Foreign key referencing the stakeholder acting as domain lead.", form_display: true, table_display: true, order: 5 }
  created_at: { type: datetime, comment: "Created", tooltip: "Timestamp when the domain was created.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Updated", tooltip: "Timestamp when the domain was last updated.", form_display: false, table_display: true }
```

## TAG_CATEGORIES
```yaml
table: tag_categories
comment: Tag Category
tooltip: Taxonomy groupings for tags and classifications such as Sensitivity, GDPR, or Data Tier.
columns:
  category_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the tag category.", form_display: true, table_display: true, order: 1 }
  name: { type: varchar, len: 50, nullable: false, unique: true, comment: "Name", tooltip: "Name of the tag category such as Sensitivity or Data Tier.", form_display: true, table_display: true, order: 2 }
  description: { type: text, comment: "Description", tooltip: "Description of what this category classifies.", form_display: true, table_display: true, order: 3 }
```

## TAGS
```yaml
table: tags
comment: Tag
tooltip: Individual tags or classification labels applied to data assets, tables, or fields.
columns:
  tag_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the tag.", form_display: true, table_display: true, order: 1 }
  category_id: { type: integer, fk: "tag_categories.category_id", comment: "Category", tooltip: "Foreign key referencing the parent tag category.", form_display: true, table_display: true, order: 2 }
  name: { type: varchar, len: 100, nullable: false, unique: true, comment: "Name", tooltip: "Name of the tag such as PII, Confidential, or Gold-Tier.", form_display: true, table_display: true, order: 3 }
  description: { type: text, comment: "Description", tooltip: "Definition and criteria for applying this tag.", form_display: true, table_display: true, order: 4 }
  created_at: { type: datetime, comment: "Created", tooltip: "Timestamp when the tag was created.", form_display: false, table_display: true }
```

## GLOSSARY_TERMS
```yaml
table: glossary_terms
comment: Glossary Term
tooltip: Enterprise business glossary establishing authoritative definitions for metrics and concepts.
columns:
  term_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the glossary term.", form_display: true, table_display: true, order: 1 }
  domain_id: { type: integer, fk: "domains.domain_id", comment: "Domain", tooltip: "Foreign key referencing the domain responsible for the term.", form_display: true, table_display: true, order: 2 }
  term_name: { type: varchar, len: 150, nullable: false, unique: true, comment: "Name", tooltip: "Standardized business term name such as Monthly Active User or Net Revenue.", form_display: true, table_display: true, order: 3 }
  definition: { type: text, nullable: false, comment: "Definition", tooltip: "Authoritative business definition of the term.", form_display: true, table_display: true, order: 4 }
  business_rules: { type: text, comment: "Rules", tooltip: "Calculation logic, constraints, or qualitative rules defining the term.", form_display: true, table_display: true, order: 5 }
  status: { type: varchar, len: 30, default: "Draft", comment: "Status", tooltip: "Lifecycle status of the term such as Draft, Approved, or Deprecated.", form_display: true, table_display: true, order: 6 }
  steward_id: { type: integer, fk: "stakeholders.stakeholder_id", comment: "Steward", tooltip: "Foreign key referencing the business steward responsible for the term.", form_display: true, table_display: true, order: 7 }
  created_at: { type: datetime, comment: "Created", tooltip: "Timestamp when the term was created.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Updated", tooltip: "Timestamp when the term was last updated.", form_display: false, table_display: true }
```

## DATA_SOURCES
```yaml
table: data_sources
comment: Data Source
tooltip: Physical or cloud data storage systems containing datasets.
columns:
  source_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the data source.", form_display: true, table_display: true, order: 1 }
  name: { type: varchar, len: 100, nullable: false, unique: true, comment: "Name", tooltip: "Name of the data source instance such as Production PostgreSQL or Snowflake DW.", form_display: true, table_display: true, order: 2 }
  source_type: { type: varchar, len: 50, nullable: false, comment: "Type", tooltip: "Technology platform type such as PostgreSQL, Snowflake, S3, BigQuery, or Kafka.", form_display: true, table_display: true, order: 3 }
  connection_uri: { type: text, comment: "URI", tooltip: "Connection endpoint or URI with credentials scrubbed.", form_display: true, table_display: true, order: 4 }
  domain_id: { type: integer, fk: "domains.domain_id", comment: "Domain", tooltip: "Foreign key referencing the domain owning this source.", form_display: true, table_display: true, order: 5 }
  technical_owner_id: { type: integer, fk: "stakeholders.stakeholder_id", comment: "Owner", tooltip: "Foreign key referencing the technical owner responsible for infrastructure.", form_display: true, table_display: true, order: 6 }
  created_at: { type: datetime, comment: "Created", tooltip: "Timestamp when the source was registered.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Updated", tooltip: "Timestamp when the source metadata was last updated.", form_display: false, table_display: true }
```

## ASSET_SCHEMAS
```yaml
table: asset_schemas
comment: Asset Schema
tooltip: Database schema or dataset grouping within a data source.
columns:
  schema_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the schema.", form_display: true, table_display: true, order: 1 }
  source_id: { type: integer, fk: "data_sources.source_id", comment: "Source", tooltip: "Foreign key referencing the parent data source.", form_display: true, table_display: true, order: 2 }
  name: { type: varchar, len: 150, nullable: false, comment: "Name", tooltip: "Name of the schema or namespace such as public or finance_mart.", form_display: true, table_display: true, order: 3 }
  description: { type: text, comment: "Description", tooltip: "Description of the data contents within this schema.", form_display: true, table_display: true, order: 4 }
  created_at: { type: datetime, comment: "Created", tooltip: "Timestamp when the schema was recorded.", form_display: false, table_display: true }
```

## DATA_ASSETS
```yaml
table: data_assets
comment: Data Asset
tooltip: Tables, views, or streams representing discrete data entities.
columns:
  asset_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the data asset.", form_display: true, table_display: true, order: 1 }
  schema_id: { type: integer, fk: "asset_schemas.schema_id", comment: "Schema", tooltip: "Foreign key referencing the parent schema.", form_display: true, table_display: true, order: 2 }
  name: { type: varchar, len: 200, nullable: false, comment: "Name", tooltip: "Name of the table, view, or stream such as dim_customer or fct_monthly_revenue.", form_display: true, table_display: true, order: 3 }
  asset_type: { type: varchar, len: 50, default: "Table", comment: "Type", tooltip: "Type of asset such as Table, View, Stream, API, or File.", form_display: true, table_display: true, order: 4 }
  description: { type: text, comment: "Description", tooltip: "Business and technical description of the asset.", form_display: true, table_display: true, order: 5 }
  business_owner_id: { type: integer, fk: "stakeholders.stakeholder_id", comment: "Business Owner", tooltip: "Foreign key referencing the business owner responsible for asset data quality and meaning.", form_display: true, table_display: true, order: 6 }
  technical_owner_id: { type: integer, fk: "stakeholders.stakeholder_id", comment: "Tech Owner", tooltip: "Foreign key referencing the technical owner responsible for the ETL or pipeline.", form_display: true, table_display: true, order: 7 }
  row_count: { type: integer, comment: "Rows", tooltip: "Latest observed total row count.", form_display: true, table_display: true, order: 8 }
  bytes_size: { type: integer, comment: "Size", tooltip: "Latest observed storage size in bytes.", form_display: true, table_display: true, order: 9 }
  created_at: { type: datetime, comment: "Created", tooltip: "Timestamp when the asset was registered.", form_display: false, table_display: true }
  updated_at: { type: datetime, comment: "Updated", tooltip: "Timestamp when the asset metadata was last refreshed.", form_display: false, table_display: true }
```

## ASSET_FIELDS
```yaml
table: asset_fields
comment: Asset Field
tooltip: Columns or fields within a data asset including typing, nullability, and relational PK or FK constraints.
columns:
  field_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the field or column.", form_display: true, table_display: true, order: 1 }
  asset_id: { type: integer, fk: "data_assets.asset_id", comment: "Asset", tooltip: "Foreign key referencing the parent data asset or table.", form_display: true, table_display: true, order: 2 }
  name: { type: varchar, len: 150, nullable: false, comment: "Name", tooltip: "Name of the column or field such as customer_id or email_address.", form_display: true, table_display: true, order: 3 }
  data_type: { type: varchar, len: 100, nullable: false, comment: "Type", tooltip: "SQL data type of the column such as VARCHAR(255), TIMESTAMP, INTEGER, or NUMERIC(18,2).", form_display: true, table_display: true, order: 4 }
  max_length: { type: integer, comment: "Length", tooltip: "Maximum character or byte length if applicable.", form_display: true, table_display: true, order: 5 }
  is_nullable: { type: boolean, default: true, comment: "Nullable", tooltip: "Flag indicating whether the column allows NULL values.", form_display: true, table_display: true, order: 6 }
  is_primary_key: { type: boolean, default: false, comment: "PK", tooltip: "Flag indicating whether this field is part of the asset primary key.", form_display: true, table_display: true, order: 7 }
  is_foreign_key: { type: boolean, default: false, comment: "FK", tooltip: "Flag indicating whether this field acts as a foreign key.", form_display: true, table_display: true, order: 8 }
  foreign_key_target_field_id: { type: integer, fk: "asset_fields.field_id", comment: "Target", tooltip: "Self-referencing foreign key linking this column directly to its target primary key field.", form_display: true, table_display: true, order: 9 }
  description: { type: text, comment: "Description", tooltip: "Column-level documentation explaining the meaning of the data stored.", form_display: true, table_display: true, order: 10 }
  created_at: { type: datetime, comment: "Created", tooltip: "Timestamp when the field metadata was recorded.", form_display: false, table_display: true }
```

## ASSET_TERM_MAPPINGS
```yaml
table: asset_term_mappings
comment: Asset Term Mapping
tooltip: Associates data assets with business glossary terms they embody.
columns:
  asset_id: { type: integer, fk: "data_assets.asset_id", comment: "Asset", tooltip: "Foreign key referencing the data asset.", form_display: true, table_display: true, order: 1 }
  term_id: { type: integer, fk: "glossary_terms.term_id", comment: "Term", tooltip: "Foreign key referencing the glossary term.", form_display: true, table_display: true, order: 2 }
```

## FIELD_TERM_MAPPINGS
```yaml
table: field_term_mappings
comment: Field Term Mapping
tooltip: Associates specific columns or fields with glossary terms.
columns:
  field_id: { type: integer, fk: "asset_fields.field_id", comment: "Field", tooltip: "Foreign key referencing the asset field or column.", form_display: true, table_display: true, order: 1 }
  term_id: { type: integer, fk: "glossary_terms.term_id", comment: "Term", tooltip: "Foreign key referencing the glossary term.", form_display: true, table_display: true, order: 2 }
```

## ASSET_TAG_MAPPINGS
```yaml
table: asset_tag_mappings
comment: Asset Tag Mapping
tooltip: Applies tags and classifications to data assets.
columns:
  asset_id: { type: integer, fk: "data_assets.asset_id", comment: "Asset", tooltip: "Foreign key referencing the data asset.", form_display: true, table_display: true, order: 1 }
  tag_id: { type: integer, fk: "tags.tag_id", comment: "Tag", tooltip: "Foreign key referencing the tag.", form_display: true, table_display: true, order: 2 }
```

## FIELD_TAG_MAPPINGS
```yaml
table: field_tag_mappings
comment: Field Tag Mapping
tooltip: Applies granular classifications at the individual column level.
columns:
  field_id: { type: integer, fk: "asset_fields.field_id", comment: "Field", tooltip: "Foreign key referencing the asset field or column.", form_display: true, table_display: true, order: 1 }
  tag_id: { type: integer, fk: "tags.tag_id", comment: "Tag", tooltip: "Foreign key referencing the tag.", form_display: true, table_display: true, order: 2 }
```

## DATA_QUALITY_RULES
```yaml
table: data_quality_rules
comment: Data Quality Rule
tooltip: Automated data quality checks and assertions defined for assets or fields.
columns:
  rule_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the data quality rule.", form_display: true, table_display: true, order: 1 }
  asset_id: { type: integer, fk: "data_assets.asset_id", comment: "Asset", tooltip: "Foreign key referencing the target data asset being validated.", form_display: true, table_display: true, order: 2 }
  field_id: { type: integer, fk: "asset_fields.field_id", comment: "Field", tooltip: "Optional foreign key referencing a specific column if the rule is column-level.", form_display: true, table_display: true, order: 3 }
  rule_type: { type: varchar, len: 50, nullable: false, comment: "Type", tooltip: "Type of assertion such as NotNull, Unique, Range, Freshness, or Expression.", form_display: true, table_display: true, order: 4 }
  rule_expression: { type: text, nullable: false, comment: "Expression", tooltip: "Executable condition or expression such as NULL_COUNT == 0 or VALUE >= 0.", form_display: true, table_display: true, order: 5 }
  severity: { type: varchar, len: 20, default: "Warning", comment: "Severity", tooltip: "Alert severity level if the rule fails.", form_display: true, table_display: true, order: 6 }
  created_at: { type: datetime, comment: "Created", tooltip: "Timestamp when the rule was created.", form_display: false, table_display: true }
```

## DATA_QUALITY_EXECUTIONS
```yaml
table: data_quality_executions
comment: Data Quality Execution
tooltip: Execution history and logs for data quality rules.
columns:
  execution_id: { type: integer, pk: true, autoincrement: true, comment: "ID", tooltip: "Unique auto-incrementing identifier for the execution record.", form_display: true, table_display: true, order: 1 }
  rule_id: { type: integer, fk: "data_quality_rules.rule_id", comment: "Rule", tooltip: "Foreign key referencing the executed data quality rule.", form_display: true, table_display: true, order: 2 }
  status: { type: varchar, len: 30, nullable: false, comment: "Status", tooltip: "Execution outcome status such as Passed, Failed, or Error.", form_display: true, table_display: true, order: 3 }
  failed_records_count: { type: integer, default: 0, comment: "Failed", tooltip: "Number of records failing the assertion condition.", form_display: true, table_display: true, order: 4 }
  total_records_checked: { type: integer, comment: "Checked", tooltip: "Total number of records scanned during execution.", form_display: true, table_display: true, order: 5 }
  executed_at: { type: datetime, comment: "Executed", tooltip: "Timestamp when the quality check was executed.", form_display: false, table_display: true }
  error_message: { type: text, comment: "Error", tooltip: "Error details if the execution failed due to system exception.", form_display: true, table_display: true, order: 6 }
```
