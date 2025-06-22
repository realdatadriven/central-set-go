# Central-Set-Go

## Overview

This open-source project is a dynamic, data-driven, and configuration-driven application built with Golang. Out of the box, it provides an admin app that allows users to manage multiple applications, offering built-in authentication, user management, and role-based access control at the CRUD level for each table.

The application provides:

- User and Role Management with fine-grained access control (Create, Read, Update, Delete per table).
- App Management: Create apps, add menus, associate tables to menus, define roles, and set access permissions.
- Default UI with customizable table layouts and CRUD operations.
- Hooks to override default CRUD actions with custom components.
- Report automation and ETL support powered by DuckDB and etlx.
- Built-in support for basic dashboards using `evidence.dev`.
- The UI is built using Svelte.

## Features

### Authentication & User Management

Central-Set-Go implements a robust authentication system based on JSON Web Tokens (JWT) that provides secure access control across all application endpoints. The authentication system supports multiple languages and maintains user sessions through token-based authentication, ensuring that user credentials are validated securely without storing sensitive information on the client side.

The authentication flow begins with user login through the dedicated login endpoint, which validates credentials against the configured database and returns a JWT token upon successful authentication. This token contains encoded user information including user ID, username, role assignments, and permissions, all cryptographically signed to prevent tampering. The token-based approach allows for stateless authentication, making the system highly scalable and suitable for distributed deployments.

Role-based access control is implemented at the granular level, allowing administrators to assign specific CRUD permissions for each database table to different user roles. This fine-grained permission system ensures that users can only access and modify data according to their assigned roles and responsibilities. The system supports multiple roles per user and dynamic role assignment, providing flexibility in managing complex organizational structures.

User management capabilities include comprehensive user profile management with support for multiple languages, timezone settings, and profile attachments. The system tracks user activity through creation and modification timestamps, and supports both active and inactive user states for administrative control. Password management includes secure hashing using bcrypt and support for password change operations through dedicated endpoints.

### App & Menu Management

The application management system in Central-Set-Go provides a sophisticated framework for creating and managing multiple applications within a single deployment. Each application can have its own database configuration, menu structure, and access permissions, allowing organizations to maintain separate business applications while sharing common infrastructure and user management.

Application configuration includes essential metadata such as application ID, name, description, version information, and database connection details. The system supports multiple database backends including SQLite, PostgreSQL, and MySQL, with each application able to connect to its own dedicated database or share databases with other applications as needed. This flexibility allows for both isolated and integrated application architectures depending on organizational requirements.

Menu management provides dynamic menu generation based on database table structures and user permissions. Administrators can create hierarchical menu structures that automatically reflect the underlying data model, with each menu item corresponding to specific database tables or custom views. The menu system supports role-based visibility, ensuring that users only see menu items for which they have appropriate access permissions.

The system automatically generates CRUD interfaces for each table associated with menu items, providing immediate functionality without requiring custom development. Menu items can be configured with custom labels, icons, and ordering to create intuitive user interfaces that match organizational workflows and terminology.

### Dynamic CRUD Operations

Central-Set-Go's dynamic CRUD system represents one of its most powerful features, automatically generating complete Create, Read, Update, and Delete interfaces for any database table without requiring manual coding. This system analyzes database schema information and creates appropriate user interfaces and API endpoints that respect the underlying data types, constraints, and relationships.

The Read operations support sophisticated querying capabilities including field-specific filtering, pattern matching, sorting, and pagination. Users can apply multiple filters simultaneously using various comparison operators such as equality, greater than, less than, and pattern matching for text fields. The system supports both simple and complex queries, with the ability to join multiple tables and apply aggregate functions when needed.

Create operations include comprehensive data validation based on database constraints and custom validation rules. The system automatically generates appropriate input forms with proper field types, validation messages, and user-friendly error handling. Support for file uploads and complex data types ensures that the system can handle diverse data requirements across different business domains.

Update operations maintain data integrity through optimistic locking and change tracking. The system records modification timestamps and user information for audit purposes, and supports both partial and complete record updates. Validation rules are applied consistently during updates to ensure data quality and business rule compliance.

Delete operations support both soft deletion and permanent removal, with configurable policies based on business requirements. Soft deletion maintains data integrity by marking records as excluded rather than physically removing them, allowing for data recovery and audit trail maintenance. The system includes safeguards against accidental deletion and supports bulk operations for efficient data management.

### ETL & Reporting

The Extract, Transform, Load (ETL) capabilities in Central-Set-Go provide comprehensive data integration and processing functionality powered by DuckDB, a high-performance analytical database engine. The ETL system supports multiple data sources including flat files, ODBC databases, and direct database connections, making it suitable for complex data integration scenarios.

Data extraction supports various file formats including CSV, Excel, and other structured data formats. The system can connect to external databases through ODBC connections, enabling integration with enterprise systems such as IBM iSeries, Oracle, SQL Server, and other database platforms. Each extraction operation can be configured with custom SQL queries, data validation rules, and transformation logic to ensure data quality and consistency.

The transformation engine leverages DuckDB's powerful SQL capabilities to perform complex data transformations including aggregations, joins, calculations, and data type conversions. Users can define custom transformation pipelines using SQL scripts, with support for parameterized queries and dynamic data processing based on runtime conditions.

Data validation is integrated throughout the ETL process, with configurable validation rules that can check data completeness, format compliance, and business rule adherence. The system can automatically reject invalid data, generate error reports, and provide detailed logging for troubleshooting and audit purposes.

Load operations support both full and incremental data loading strategies, with automatic duplicate detection and conflict resolution. The system maintains data lineage information and provides comprehensive logging of all ETL operations for compliance and debugging purposes.

### Dashboard Support

Central-Set-Go includes integrated dashboard capabilities that leverage the evidence.dev framework to provide real-time data visualization and reporting. The dashboard system automatically compiles configurations and connects to the underlying data sources to provide up-to-date insights and analytics.

Dashboard components can display various chart types, tables, and key performance indicators based on the data stored in the system. The integration with DuckDB enables high-performance analytical queries that can process large datasets efficiently, making the dashboard system suitable for both operational and analytical reporting requirements.

The dashboard configuration system allows users to create custom views and reports without requiring technical expertise. Dashboard layouts are responsive and support both desktop and mobile viewing, ensuring that critical business information is accessible across different devices and platforms.




## API Documentation

Central-Set-Go exposes a comprehensive RESTful API that provides programmatic access to all system functionality. The API follows consistent patterns and conventions, making it easy to integrate with external systems and develop custom applications. All API endpoints use JSON for data exchange and support internationalization through language parameters.

### API Architecture and Design Principles

The API architecture is built around a unified endpoint structure that provides consistent request and response patterns across all operations. Every API request follows a standardized format that includes language specification, operation-specific data, and application context information. This consistency simplifies client development and reduces the learning curve for developers working with the API.

Authentication is handled through JWT tokens that are included in the Authorization header of each request. The token-based authentication system ensures stateless operation and enables horizontal scaling of the application. Tokens contain encoded user information and permissions, allowing the API to make authorization decisions without additional database queries.

Error handling follows HTTP status code conventions with detailed error messages in JSON format. The API provides comprehensive error information including error codes, descriptive messages, and field-specific validation errors when applicable. This approach enables client applications to provide meaningful feedback to users and implement appropriate error recovery strategies.

### Authentication Endpoints

The authentication system provides secure access control through JWT token-based authentication. These endpoints handle user login, token validation, and password management operations.

#### User Login

**Endpoint:** `POST /dyn_api/login/login`

The login endpoint authenticates users and returns a JWT token for subsequent API requests. This endpoint does not require authentication and serves as the entry point for all user sessions.

**Headers:**
```
Content-Type: application/json
```

**Request Body:**
```json
{
    "lang": "pt",
    "data": {
        "username": "root",
        "password": "1234"
    }
}
```

**Request Parameters:**
- `lang`: Language code for response messages (pt, en, es, etc.)
- `data.username`: User's login username
- `data.password`: User's password in plain text (encrypted during transmission)

**Response:**
The endpoint returns a JWT token that must be included in the Authorization header of subsequent requests. The token contains encoded user information including user ID, username, role assignments, and permissions. Token expiration is configurable and typically set to several hours or days depending on security requirements.

**Example Response:**
```json
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

**Error Responses:**
- `401 Unauthorized`: Invalid username or password
- `403 Forbidden`: User account is inactive or excluded
- `400 Bad Request`: Missing required fields or invalid request format

#### Token Validation

**Endpoint:** `POST /dyn_api/login/chk_token`

This endpoint validates the current JWT token and returns updated user information. It's useful for checking token validity and refreshing user session data.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "lang": "en"
}
```

**Request Parameters:**
- `lang`: Language code for response messages

**Response:**
Returns current user information and token validity status. If the token is valid, the response includes complete user profile data. If the token is expired or invalid, an error response is returned.

**Example Response:**
```json
{
    "success": true,
    "valid": true,
    "user": {
        "user_id": 1,
        "username": "root",
        "first_name": "Super",
        "last_name": "Admin",
        "email": "root@domain.com",
        "role_id": 1,
        "lang_id": 1,
        "active": true,
        "created_at": "2024-08-29T15:36:59.318618Z",
        "updated_at": "2024-08-29T15:36:59.318618Z"
    }
}
```

#### Password Change

**Endpoint:** `POST /api/`

This legacy endpoint handles password change operations for authenticated users. It requires the current password for verification before allowing the password change.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "lang": "en",
    "controller": "login",
    "action": "alter_pass",
    "data": {
        "username": "root",
        "password": 12344654,
        "new_password": 1234
    }
}
```

**Request Parameters:**
- `lang`: Language code for response messages
- `controller`: Always "login" for authentication operations
- `action`: Always "alter_pass" for password change
- `data.username`: Username for the account
- `data.password`: Current password for verification
- `data.new_password`: New password to set

**Response:**
Returns success confirmation or error details if the password change fails due to invalid current password or password policy violations.

### Administrative Endpoints

Administrative endpoints provide access to application management, table configuration, and menu structure operations. These endpoints typically require elevated permissions and are used for system configuration and maintenance.

#### Application Management

**Endpoint:** `POST /dyn_api/admin/apps`

This endpoint retrieves and manages application configurations within the Central-Set-Go system. It returns a list of available applications with their configuration details and access permissions.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "lang": "en"
}
```

**Request Parameters:**
- `lang`: Language code for response messages

**Response:**
Returns an array of application configurations including application metadata, database connections, and permission settings. Each application entry includes comprehensive information needed for client applications to connect and interact with the specific application instance.

**Example Response:**
```json
{
    "success": true,
    "applications": [
        {
            "app_id": 1,
            "app": "ADMIN",
            "app_desc": "Admin",
            "version": "1.0.0",
            "db": "ADMIN",
            "created_at": "2024-08-29T15:36:59.318618Z",
            "updated_at": "2024-08-29T15:36:59.318618Z",
            "active": true
        }
    ]
}
```

#### Table Management

**Endpoint:** `POST /dyn_api/admin/tables`

The table management endpoint provides access to database table configurations and metadata. It supports querying table structures, relationships, and access permissions within specific applications.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "lang": "en",
    "db--": {
        "drivername": "sqlite",
        "dsn": "ADM"
    },
    "data": {
        "table": "app",
        "table--.": ["lang", "app"]
    },
    "app": {
        "app_id": 1,
        "app": "ADMIN",
        "db": "ADMIN"
    }
}
```

**Request Parameters:**
- `lang`: Language code for response messages
- `db.drivername`: Database driver type (sqlite, postgresql, mysql)
- `db.dsn`: Database connection string or identifier
- `data.table`: Target table name for operations
- `data.table`: Array of related tables or fields to include
- `app`: Application context information

**Response:**
Returns detailed table metadata including column definitions, data types, constraints, relationships, and access permissions. The response includes schema information needed for dynamic UI generation and data validation.

#### Menu Management

**Endpoint:** `POST /dyn_api/admin/menu`

This endpoint manages menu structures and navigation configurations for applications. It returns hierarchical menu data that reflects the user's permissions and the application's table structure.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "lang": "en",
    "app": {
        "app_id": 1,
        "app": "ADMIN",
        "app_desc": "Admin",
        "version": "1.0.0",
        "db": "ADMIN"
    }
}
```

**Request Parameters:**
- `lang`: Language code for response messages
- `app`: Complete application context including ID, name, description, version, and database

**Response:**
Returns a hierarchical menu structure with items organized by categories and access permissions. Each menu item includes navigation information, icons, labels, and associated table or function references.

**Example Response:**
```json
{
    "success": true,
    "menu": [
        {
            "menu_id": 1,
            "menu_name": "User Management",
            "menu_desc": "Manage users and roles",
            "table_name": "users",
            "icon": "users",
            "order": 1,
            "active": true,
            "permissions": {
                "create": true,
                "read": true,
                "update": true,
                "delete": true
            }
        }
    ]
}
```


### CRUD Operations

The CRUD (Create, Read, Update, Delete) endpoints form the core of Central-Set-Go's data management capabilities. These endpoints provide comprehensive data manipulation functionality with advanced filtering, validation, and transaction support.

#### Read Operations

**Endpoint:** `POST /dyn_api/crud/read`

The read endpoint provides sophisticated data querying capabilities with support for filtering, sorting, pagination, and joins. It serves as the primary interface for data retrieval across all tables in the system.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "lang": "en",
    "data": {
        "distinct--": true,
        "join--": "none|all",
        "table": "menu",
        "limit": 11,
        "offset": 0,
        "fields--": ["user_log_id", "user_id", "action"],
        "filters--": [
            {"field": "lang_id", "cond": "=", "value": 3},
            {"field": "menu_id", "cond": "=", "value": 3}
        ],
        "order_by": [
            {"field": "lang_id2", "order": "DESC"}
        ],
        "pattern--": "brasil"
    },
    "app": {
        "app_id": 1,
        "app": "ADMIN",
        "version": "1.0.0",
        "db": "ADMIN"
    }
}
```

**Request Parameters:**
- `lang`: Language code for response messages
- `data.distinct`: Boolean flag to return only distinct records
- `data.join`: Join strategy ("none", "all", or specific join configuration)
- `data.table`: Target table name for the query
- `data.limit`: Maximum number of records to return (pagination)
- `data.offset`: Number of records to skip (pagination)
- `data.fields`: Array of specific fields to return (if not specified, returns all fields)
- `data.filters`: Array of filter conditions with field, condition, and value
- `data.order_by`: Array of sorting specifications with field and order direction
- `data.pattern`: Text pattern for full-text search across applicable fields
- `app`: Application context information

**Filter Conditions:**
The filtering system supports various comparison operators:
- `=`: Exact equality match
- `!=` or `<>`: Not equal
- `>`: Greater than
- `>=`: Greater than or equal
- `<`: Less than
- `<=`: Less than or equal
- `LIKE`: Pattern matching with wildcards
- `IN`: Value in a list of options
- `BETWEEN`: Value within a range
- `IS NULL`: Field is null
- `IS NOT NULL`: Field is not null

**Response:**
Returns paginated data results with metadata including total record count, pagination information, and the requested data records. The response structure includes both the data and metadata needed for client-side pagination and display.

**Example Response:**
```json
{
    "success": true,
    "data": [
        {
            "menu_id": 1,
            "menu_name": "Users",
            "table_name": "users",
            "lang_id": 3,
            "created_at": "2024-08-29T15:36:59.318618Z",
            "updated_at": "2024-08-29T15:36:59.318618Z"
        }
    ],
    "pagination": {
        "total": 1,
        "limit": 11,
        "offset": 0,
        "has_more": false
    }
}
```

#### Create Operations

**Endpoint:** `POST /dyn_api/crud/create`

The create endpoint handles the insertion of new records into specified tables with comprehensive validation and constraint checking.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "lang": "en",
    "data": {
        "table": "lang",
        "data": {
            "lang": "es4",
            "lang_desc": "Português"
        }
    },
    "app": {
        "app_id": 1,
        "app": "ADMIN",
        "app_desc": "Admin",
        "version": "1.0.0",
        "db": "ADMIN"
    }
}
```

**Request Parameters:**
- `lang`: Language code for response messages
- `data.table`: Target table name for the insert operation
- `data.data`: Object containing the field values for the new record
- `app`: Complete application context information

**Validation:**
The create operation performs comprehensive validation including:
- Data type validation based on table schema
- Required field validation
- Unique constraint checking
- Foreign key constraint validation
- Custom business rule validation
- Field length and format validation

**Response:**
Returns the created record with any auto-generated fields (such as IDs and timestamps) populated. The response includes the complete record data as it was stored in the database.

**Example Response:**
```json
{
    "success": true,
    "data": {
        "lang_id": 4,
        "lang": "es4",
        "lang_desc": "Português",
        "created_at": "2024-10-19T14:30:00.000Z",
        "updated_at": "2024-10-19T14:30:00.000Z"
    },
    "msg": "Record created successfully"
}
```

#### Update Operations

**Endpoint:** `POST /dyn_api/crud/update`

The update endpoint modifies existing records with support for partial updates and optimistic locking to prevent concurrent modification conflicts.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "lang": "pt",
    "data": {
        "table": "lang",
        "data": {
            "lang_id": 3,
            "lang": "pt_BR",
            "lang_desc": "Português Brasil",
            "created_at": "2022-07-22T10:32:00.791024",
            "updated_at": "2022-07-22T10:33:22.185972"
        }
    },
    "app": {
        "app_id": 1,
        "app": "ADMIN",
        "app_desc": "Admin",
        "version": "1.0.0",
        "email": null,
        "db": "ADMIN"
    }
}
```

**Request Parameters:**
- `lang`: Language code for response messages
- `data.table`: Target table name for the update operation
- `data.data`: Object containing the field values to update (must include primary key)
- `app`: Complete application context information

**Update Behavior:**
- Only specified fields are updated (partial update support)
- Primary key fields are used to identify the target record
- Timestamp fields are automatically updated
- Optimistic locking prevents concurrent modification conflicts
- Validation rules are applied to all modified fields

**Response:**
Returns the updated record with current field values and updated timestamps. The response confirms the successful update and provides the current state of the record.

**Example Response:**
```json
{
    "success": true,
    "data": {
        "lang_id": 3,
        "lang": "pt_BR",
        "lang_desc": "Português Brasil",
        "created_at": "2022-07-22T10:32:00.791024",
        "updated_at": "2024-10-19T14:35:00.000Z"
    },
    "msg": "Record updated successfully"
}
```

#### Delete Operations

**Endpoint:** `POST /dyn_api/crud/delete`

The delete endpoint supports both soft deletion (marking records as excluded) and permanent deletion based on configuration and user permissions.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "lang": "en",
    "data": {
        "table": "lang",
        "data": {
            "lang_id": 3,
            "exclude": true,
            "permanently--": true
        }
    },
    "app": {
        "app_id": 1,
        "app": "ADMIN",
        "app_desc": "Admin",
        "version": "1.0.0",
        "email": null,
        "db": "ADMIN"
    }
}
```

**Request Parameters:**
- `lang`: Language code for response messages
- `data.table`: Target table name for the delete operation
- `data.data.{primary_key}`: Primary key value to identify the record
- `data.data.exclude`: Boolean flag for soft deletion
- `data.data.permanently--`: Boolean flag for permanent deletion
- `app`: Complete application context information

**Deletion Types:**
- **Soft Delete**: Sets the `excluded` flag to true, preserving data for audit and recovery
- **Permanent Delete**: Physically removes the record from the database
- **Cascade Delete**: Handles related records based on foreign key constraints

**Response:**
Confirms the deletion operation and provides information about the affected record. For soft deletes, the response includes the updated record state.

**Example Response:**
```json
{
    "success": true,
    "msg": "Record deleted successfully",
    "deleted_record": {
        "lang_id": 3,
        "excluded": true,
        "updated_at": "2024-10-19T14:40:00.000Z"
    }
}
```

#### Custom Query Operations

**Endpoint:** `POST /dyn_api/crud/query`

The query endpoint allows execution of custom SQL queries using the integrated DuckDB engine, providing advanced analytical and reporting capabilities.

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "lang": "en",
    "data": {
        "db": "DB.duckdb",
        "limit": 10,
        "offset": 0,
        "query": "select * from \"table\""
    },
    "app": {
        "app_id": 1,
        "app": "ADMIN",
        "app_desc": "Admin",
        "version": "1.0.0",
        "email": null,
        "db": "ADMIN"
    }
}
```

**Request Parameters:**
- `lang`: Language code for response messages
- `data.db`: Target database file (for DuckDB operations)
- `data.limit`: Maximum number of records to return
- `data.offset`: Number of records to skip for pagination
- `data.query`: SQL query string to execute
- `app`: Application context information

**Query Capabilities:**
- Full SQL SELECT statement support
- Join operations across multiple tables
- Aggregate functions and grouping
- Window functions for advanced analytics
- Common Table Expressions (CTEs)
- Subqueries and complex filtering

**Security Considerations:**
- Query execution is restricted to SELECT statements for security
- User permissions are validated before query execution
- Query timeout limits prevent resource exhaustion
- SQL injection protection through parameterized queries

**Response:**
Returns query results with column metadata and pagination information. The response structure adapts to the query output format.

**Example Response:**
```json
{
    "success": true,
    "data": [
        {
            "id": 1,
            "name": "Sample Record",
            "value": 100.50,
            "date": "2024-10-19"
        }
    ],
    "columns": [
        {"name": "id", "type": "INTEGER"},
        {"name": "name", "type": "VARCHAR"},
        {"name": "value", "type": "DECIMAL"},
        {"name": "date", "type": "DATE"}
    ],
    "pagination": {
        "total": 1,
        "limit": 10,
        "offset": 0,
        "has_more": false
    }
}
```


### File Management

Central-Set-Go provides comprehensive file upload and management capabilities through dedicated endpoints that support various file types and processing workflows.

#### File Upload

**Endpoint:** `POST /upload`

The upload endpoint handles multipart file uploads with support for temporary and permanent storage, file validation, and metadata extraction.

**Headers:**
```
enctype: multipart/form-data
Authorization: Bearer <JWT_TOKEN>
```

**Request Body (Form Data):**
```
lang: pt
tmp: true
path: xxxx
file: [binary file data]
```

**Request Parameters:**
- `lang`: Language code for response messages
- `tmp`: Boolean flag indicating temporary storage (true) or permanent storage (false)
- `path`: Target path or directory for file storage
- `file`: Binary file data (multipart form field)

**File Processing:**
The upload system provides several processing capabilities:
- File type validation based on MIME type and file extension
- Virus scanning and security validation
- Automatic file naming and collision resolution
- Metadata extraction for supported file types
- Thumbnail generation for image files
- File size and quota validation

**Storage Options:**
- **Temporary Storage**: Files stored temporarily for processing workflows
- **Permanent Storage**: Files stored permanently with backup and versioning
- **Cloud Storage**: Integration with cloud storage providers
- **Local Storage**: Files stored on local filesystem with appropriate permissions

**Response:**
Returns file upload confirmation with file metadata, storage location, and processing status. The response includes information needed for subsequent file operations.

**Example Response:**
```json
{
    "success": true,
    "file": {
        "file_id": "8a158fac-f01b-4ad0-a87b-0d7551f034d1",
        "original_name": "File_Name.20240829.xlsx",
        "stored_name": "8a158fac-f01b-4ad0-a87b-0d7551f034d1_20240917.csv",
        "file_size": 1048576,
        "mime_type": "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet",
        "upload_date": "2024-10-19T14:45:00.000Z",
        "temporary": true,
        "path": "/uploads/temp/"
    },
    "msg": "File uploaded successfully"
}
```

**Error Handling:**
- File size exceeding limits
- Unsupported file types
- Storage quota exceeded
- File corruption or validation failures
- Permission denied errors

### ETL Operations

The ETL (Extract, Transform, Load) system provides comprehensive data integration capabilities with support for multiple data sources, transformation pipelines, and validation rules.

#### Data Extraction

**Endpoint:** `POST /dyn_api/etl/extract`

The extraction endpoint supports multiple data source types including files, ODBC databases, and direct database connections with configurable transformation and validation rules.

#### File-Based Extraction

**Headers:**
```
Content-Type: application/json
Authorization: Bearer <JWT_TOKEN>
```

**Request Body:**
```json
{
    "lang": "en",
    "data": {
        "data": {
            "date_ref": "2024-09-17",
            "date_ref_": ["2024-08-29"],
            "file_": "file_name.20240829.xlsx",
            "file__": "8a158fac-f01b-4ad0-a87b-0d7551f034d1_20240917.csv",
            "database": "sqlite_test.duckdb",
            "name": "test",
            "save_only_temp": false,
            "destination_table": "table_name",
            "check_ref_date": false,
            "ref_date_field": "file_ref",
            "etl_rbase_input_conf": {
                "type_": "file-duckdb",
                "type": "duckdb",
                "duckdb": {
                    "pragmas_config_sql_start": [
                        "ATTACH 'database/db.duckdb' AS db"
                    ],
                    "pragmas_config_sql_end": [
                        "DETACH DB"
                    ],
                    "valid_": [
                        {
                            "sql": "SELECT * FROM \"<table>\" WHERE date_field={YYYYMMDD} LIMIT 10",
                            "rule": "throw_if_not_empty",
                            "msg": "The table (<table>) already has the data from the date YYYY/MM/DD"
                        },
                        {
                            "sql": "SELECT * FROM '<file>' WHERE date_field={YYYYMMDD} LIMIT 10",
                            "rule": "throw_if_empty",
                            "msg": "The file (<file>) has no data from the date \"YYYY/MM/DD\""
                        }
                    ],
                    "sql": "INSERT INTO main.\"<table>\" BY NAME SELECT * FROM DB.\"<table>\" WHERE \"date_field\" = YYYYMMDD"
                }
            }
        }
    },
    "app": {
        "app_id": 1,
        "app": "ADMIN",
        "app_desc": "Admin",
        "version": "1.0.0",
        "email": null,
        "db": "ADMIN"
    }
}
```

**Request Parameters:**
- `lang`: Language code for response messages
- `data.data.date_ref`: Reference date for the extraction operation
- `data.data.file_`: Source file name for extraction
- `data.data.database`: Target DuckDB database file
- `data.data.destination_table`: Target table name for loaded data
- `data.data.check_ref_date`: Boolean flag for date validation
- `data.data.etl_rbase_input_conf`: ETL configuration object with processing rules

**ETL Configuration:**
The ETL configuration supports various processing types:
- **file-duckdb**: Extract from files and load into DuckDB
- **odbc-csv-duckdb**: Extract from ODBC sources via CSV intermediate format
- **database-direct**: Direct database-to-database transfers

**Validation Rules:**
- `throw_if_empty`: Fail if query returns no results
- `throw_if_not_empty`: Fail if query returns any results
- Custom validation with SQL expressions
- Data quality checks and business rule validation

#### ODBC-Based Extraction

**Request Body:**
```json
{
    "lang": "en",
    "data": {
        "data": {
            "date_ref": "2024-09-30",
            "database": "sqlite_test.duckdb",
            "save_only_temp": true,
            "destination_table": "table_name",
            "check_ref_date": false,
            "ref_date_field": "date_field",
            "etl_rbase_input_conf": {
                "type": "odbc-csv-duckdb",
                "params": {"odbc_conn": "Driver={...ODBC Driver};System=host;Uid=@USERNAME;Pwd=@PASSWORD" },
                "query": "SELECT * FROM **",
                "duckdb": {
                    "extentions": [],
                    "valid": [
                        {
                            "sql": "SELECT DISTINCT \"date_field\" FROM READ_CSV('<filename>', HEADER = TRUE) WHERE \"date_field\" = '{YYYYMMDD}' LIMIT 10",
                            "rule": "throw_if_empty",
                            "msg": "A data na origem (<table>) é diferente de DD/MM/YYYY!"
                        }
                    ],
                    "sql": "INSERT INTO \"<table>\" BY NAME SELECT * FROM READ_CSV('<filename>', HEADER = TRUE)"
                }
            }
        }
    },
    "app": {
        "app_id": 1,
        "app": "ADMIN",
        "app_desc": "Admin",
        "version": "1.0.0",
        "email": null,
        "db": "ADMIN"
    }
}
```

**ODBC Configuration:**
- Connection string with driver specification
- Environment variable support for credentials
- Query parameterization for dynamic data extraction
- Error handling and retry logic

**Response:**
Returns extraction status, processed record counts, validation results, and any error messages. The response includes detailed information about the ETL operation for monitoring and debugging.

**Example Response:**
```json
{
    "success": true,
    "msg": "ETL extraction completed successfully"
}
```

### API Response Patterns

All API endpoints follow consistent response patterns that facilitate client development and error handling. Understanding these patterns is essential for effective API integration.

#### Success Responses

Successful API responses include a `success: true` field and relevant data. The response structure varies by endpoint but maintains consistency in status indication and error handling.

```json
{
    "success": true,
    "data": { /* endpoint-specific data */ },
    "msg": "Operation completed successfully",
    "metadata": { /* additional information */ }
}
```

#### Error Responses

Error responses include detailed error information with HTTP status codes and descriptive messages. The error structure provides sufficient information for client applications to handle errors appropriately.

```json
{
    "success": false,
    "error": {
        "code": "VALIDATION_ERROR",
        "msg": "Invalid input data",
        "details": [
            {
                "field": "username",
                "msg": "Username is required"
            }
        ]
    }
}
```

#### Pagination

Endpoints that return multiple records include pagination metadata to support efficient data browsing and loading.

```json
{
    "success": true,
    "data": [ /* array of records */ ],
    "pagination": {
        "total": 1000,
        "limit": 50,
        "offset": 0,
        "has_more": true,
        "next_offset": 50
    }
}
```

### API Security and Best Practices

Central-Set-Go implements comprehensive security measures to protect API endpoints and ensure data integrity. Understanding these security features is crucial for secure API integration.

#### Authentication Security

- JWT tokens use strong cryptographic signatures
- Token expiration prevents long-term token abuse
- Refresh token mechanism for session management
- Rate limiting prevents brute force attacks

#### Authorization Controls

- Role-based access control at the endpoint level
- Table-level permissions for CRUD operations
- Field-level access control for sensitive data
- Dynamic permission evaluation based on context

#### Data Protection

- Input validation and sanitization
- SQL injection prevention through parameterized queries
- Cross-site scripting (XSS) protection
- Data encryption for sensitive information

#### API Rate Limiting

- Request rate limiting per user and endpoint
- Burst protection for high-volume operations
- Graceful degradation under load
- Error responses for rate limit violations

#### Monitoring and Logging

- Comprehensive API access logging
- Performance monitoring and alerting
- Security event detection and response
- Audit trail for compliance requirements


## Getting Started

Central-Set-Go is designed for easy deployment and configuration across various environments. The following sections provide comprehensive guidance for setting up and running the application in development, testing, and production environments.

### Prerequisites

Before installing Central-Set-Go, ensure that your system meets the following requirements:

**Go Programming Language**: Central-Set-Go requires Go version 1.20 or higher. The application leverages modern Go features and standard library improvements introduced in recent versions. You can download Go from the official website at https://golang.org/dl/ and follow the installation instructions for your operating system.

**Database Support**: The application supports multiple database backends including SQLite, PostgreSQL, and MySQL. For development and testing, SQLite provides a lightweight option that requires no additional setup. For production deployments, PostgreSQL or MySQL offer better performance and scalability characteristics. Ensure that your chosen database system is installed and properly configured before proceeding with the application setup.

**Docker (Optional)**: While not required for basic operation, Docker provides a convenient deployment option that simplifies environment management and ensures consistent behavior across different systems. Docker installation instructions are available at https://docs.docker.com/get-docker/ for all major operating systems.

**DuckDB Dependencies**: The ETL functionality relies on DuckDB for analytical processing. The necessary DuckDB libraries are included with the application, but ensure that your system has the required runtime dependencies for optimal performance.

### Installation

The installation process involves cloning the repository, installing dependencies, and configuring the application for your specific environment.

#### Repository Setup

Begin by cloning the Central-Set-Go repository from GitHub to your local development environment:

```bash
git clone https://github.com/realdatadriven/central-set-go.git
cd central-set-go
```

This command downloads the complete source code including all necessary files, documentation, and example configurations. The repository structure includes source code, configuration templates, database migration scripts, and deployment resources.

#### Dependency Installation

Central-Set-Go uses Go modules for dependency management, which simplifies the installation process and ensures reproducible builds:

```bash
go mod tidy
```

This command downloads and installs all required dependencies as specified in the go.mod file. The process includes downloading third-party libraries for database connectivity, JWT token handling, HTTP routing, and other essential functionality. The go mod tidy command also removes any unused dependencies and updates the go.sum file with cryptographic checksums for security verification.

#### Application Startup

Once dependencies are installed, you can start the application in development mode:

```bash
go run ./cmd/api
```

This command compiles and runs the application with default settings suitable for development and testing. The application will start on the default port (typically 4444) and create necessary database tables if they don't already exist. During the first startup, the application performs initial setup including creating the admin user account and default application configuration.

For production deployments, compile the application to a binary executable:

```bash
go build -o central-set-go ./cmd/api
./central-set-go
```

The compiled binary includes all dependencies and can be deployed to production servers without requiring Go to be installed on the target system.

### Configuration

Central-Set-Go uses environment-based configuration to support different deployment scenarios and maintain security best practices. Configuration is managed through environment variables and configuration files.

#### Environment Configuration

The application includes a template configuration file `.env-exemple` that demonstrates all available configuration options. Copy this file to `.env` and modify the values according to your environment requirements:

```bash
cp .env-exemple .env
```

Edit the `.env` file to configure essential settings:

**Database Configuration:**
```
DB_DRIVER=sqlite
DB_DSN=./data/central-set-go.db
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
```

**Server Configuration:**
```
SERVER_PORT=4444
SERVER_HOST=0.0.0.0
SERVER_READ_TIMEOUT=30s
SERVER_WRITE_TIMEOUT=30s
```

**Security Configuration:**
```
JWT_SECRET=your-secret-key-here
JWT_EXPIRATION=24h
BCRYPT_COST=12
```

**File Storage Configuration:**
```
UPLOAD_PATH=./uploads
MAX_FILE_SIZE=10MB
ALLOWED_FILE_TYPES=jpg,jpeg,png,pdf,xlsx,csv
```

#### Database Configuration

Central-Set-Go supports multiple database backends with specific configuration requirements for each:

**SQLite Configuration:**
SQLite provides a lightweight, file-based database suitable for development and small-scale deployments. Configure SQLite by specifying the database file path:

```
DB_DRIVER=sqlite
DB_DSN=./data/central-set-go.db
```

**PostgreSQL Configuration:**
PostgreSQL offers robust features for production deployments including advanced indexing, replication, and performance optimization:

```
DB_DRIVER=postgres
DB_DSN=postgres://username:password@localhost:5432/central_set_go?sslmode=disable
```

**MySQL Configuration:**
MySQL provides excellent performance and widespread compatibility for enterprise deployments:

```
DB_DRIVER=mysql
DB_DSN=username:password@tcp(localhost:3306)/central_set_go?parseTime=true
```

#### Security Configuration

Security configuration includes authentication settings, encryption parameters, and access control policies:

**JWT Token Configuration:**
JSON Web Tokens provide stateless authentication with configurable expiration and signing algorithms:

```
JWT_SECRET=your-256-bit-secret-key
JWT_EXPIRATION=24h
JWT_ALGORITHM=HS256
```

**Password Security:**
Password hashing uses bcrypt with configurable cost parameters to balance security and performance:

```
BCRYPT_COST=12
PASSWORD_MIN_LENGTH=8
PASSWORD_REQUIRE_SPECIAL=true
```

**CORS Configuration:**
Cross-Origin Resource Sharing settings control browser access from different domains:

```
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://yourdomain.com
CORS_ALLOWED_METHODS=GET,POST,PUT,DELETE,OPTIONS
CORS_ALLOWED_HEADERS=Content-Type,Authorization
```

### Deployment

Central-Set-Go supports various deployment strategies including Docker containers, traditional server deployment, and cloud platform deployment.

#### Docker Deployment

Docker provides a consistent deployment environment that simplifies configuration management and ensures reproducible deployments across different environments.

**Building the Docker Image:**
```bash
docker build -t central-set-go .
```

This command creates a Docker image containing the compiled application and all necessary runtime dependencies. The Dockerfile includes multi-stage builds to minimize image size while including all required components.

**Running with Docker:**
```bash
docker run -p 4444:4444 -v $(pwd)/data:/app/data central-set-go
```

This command starts the application in a Docker container with port mapping and volume mounting for persistent data storage. The volume mount ensures that database files and uploaded content persist across container restarts.

**Docker Compose Deployment:**
For more complex deployments involving multiple services, use Docker Compose:

```yaml
version: '3.8'
services:
  central-set-go:
    build: .
    ports:
      - "4444:4444"
    volumes:
      - ./data:/app/data
      - ./uploads:/app/uploads
    environment:
      - DB_DRIVER=postgres
      - DB_DSN=postgres://user:password@postgres:5432/central_set_go
    depends_on:
      - postgres
  
  postgres:
    image: postgres:13
    environment:
      - POSTGRES_DB=central_set_go
      - POSTGRES_USER=user
      - POSTGRES_PASSWORD=password
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

#### Traditional Server Deployment

For traditional server deployments, compile the application and deploy the binary with appropriate system service configuration:

**Compilation for Production:**
```bash
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o central-set-go ./cmd/api
```

**Systemd Service Configuration:**
Create a systemd service file for automatic startup and process management:

```ini
[Unit]
Description=Central-Set-Go Application
After=network.target

[Service]
Type=simple
User=central-set-go
WorkingDirectory=/opt/central-set-go
ExecStart=/opt/central-set-go/central-set-go
Restart=always
RestartSec=5
Environment=ENV=production

[Install]
WantedBy=multi-user.target
```

#### Cloud Platform Deployment

Central-Set-Go can be deployed on various cloud platforms including AWS, Google Cloud Platform, and Azure:

**AWS Deployment:**
- Use AWS Elastic Beanstalk for simple application deployment
- Deploy on AWS ECS for containerized deployments
- Use AWS RDS for managed database services
- Configure AWS S3 for file storage and backups

**Google Cloud Platform:**
- Deploy on Google App Engine for serverless scaling
- Use Google Cloud Run for containerized deployments
- Configure Google Cloud SQL for managed databases
- Use Google Cloud Storage for file management

**Azure Deployment:**
- Deploy on Azure App Service for web applications
- Use Azure Container Instances for containerized deployments
- Configure Azure Database for managed database services
- Use Azure Blob Storage for file storage

### Binary Releases

Central-Set-Go provides prebuilt binary releases for major operating systems, eliminating the need for local compilation and Go installation. These releases are available through the GitHub releases page and include executables for Windows, Linux, and macOS.

#### Release Artifacts

Each release includes the following artifacts:

**Windows Releases:**
- `central-set-go-windows-amd64.exe`: 64-bit Windows executable
- `central-set-go-windows-386.exe`: 32-bit Windows executable

**Linux Releases:**
- `central-set-go-linux-amd64`: 64-bit Linux executable
- `central-set-go-linux-arm64`: ARM64 Linux executable
- `central-set-go-linux-386`: 32-bit Linux executable

**macOS Releases:**
- `central-set-go-darwin-amd64`: Intel Mac executable
- `central-set-go-darwin-arm64`: Apple Silicon Mac executable

#### Installation from Binary

Download the appropriate binary for your operating system and architecture:

```bash
# Linux/macOS example
wget https://github.com/realdatadriven/central-set-go/releases/latest/download/central-set-go-linux-amd64
chmod +x central-set-go-linux-amd64
./central-set-go-linux-amd64
```

Binary releases include all necessary dependencies and can be executed directly without additional setup requirements.

## Contributing

Central-Set-Go welcomes contributions from the community including bug reports, feature requests, documentation improvements, and code contributions. The project follows standard open-source contribution practices and maintains high standards for code quality and documentation.

### Development Guidelines

Contributors should follow established development practices including:

**Code Style:** Follow Go standard formatting and naming conventions using `gofmt` and `golint` tools. Maintain consistent code style across the entire codebase and include comprehensive comments for complex functionality.

**Testing:** Include unit tests for new functionality and ensure that existing tests continue to pass. The project maintains high test coverage and requires tests for all critical functionality.

**Documentation:** Update documentation for new features and API changes. Include examples and usage instructions for new functionality to help other developers understand and use the features effectively.

**Commit Messages:** Use clear, descriptive commit messages that explain the purpose and scope of changes. Follow conventional commit message format when possible to facilitate automated changelog generation.

### Contribution Process

1. **Fork the Repository:** Create a personal fork of the Central-Set-Go repository on GitHub
2. **Create Feature Branch:** Create a new branch for your feature or bug fix
3. **Implement Changes:** Make your changes following the development guidelines
4. **Test Thoroughly:** Ensure all tests pass and add new tests for your changes
5. **Submit Pull Request:** Create a pull request with a clear description of your changes
6. **Code Review:** Participate in the code review process and address feedback
7. **Merge:** Once approved, your changes will be merged into the main branch

### Issue Reporting

When reporting issues, include:
- Detailed description of the problem
- Steps to reproduce the issue
- Expected vs. actual behavior
- System information and configuration details
- Log files and error messages when applicable

## License

Central-Set-Go is released under the MIT License, which provides broad permissions for use, modification, and distribution while maintaining copyright attribution requirements. The MIT License is a permissive open-source license that allows both commercial and non-commercial use of the software.

The complete license text is available in the LICENSE file included with the source code. By using, modifying, or distributing Central-Set-Go, you agree to comply with the terms and conditions of the MIT License.

## Support and Community

Central-Set-Go maintains an active community of users and contributors who provide support, share knowledge, and collaborate on improvements. Community resources include:

**GitHub Issues:** Use the GitHub issue tracker for bug reports, feature requests, and technical questions. The development team actively monitors and responds to issues.

**Documentation:** Comprehensive documentation is available in the repository wiki and through inline code comments. Documentation includes API references, configuration guides, and deployment instructions.

**Community Forums:** Participate in community discussions through GitHub Discussions where users share experiences, ask questions, and provide mutual support.

**Professional Support:** Commercial support options are available for organizations requiring dedicated assistance, custom development, or enterprise-level support agreements.

The Central-Set-Go project is committed to maintaining a welcoming and inclusive community where all participants can contribute effectively and feel valued for their contributions.

