BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS flight_schema (
	flight_schema_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	flight_schema VARCHAR(200) NOT NULL, 
	flight_schema_desc TEXT, 
	flight_schema_conf TEXT, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE (flight_schema)
);
CREATE TABLE IF NOT EXISTS dashboard (
	dashboard_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	dashboard VARCHAR(200), 
	dashboard_desc TEXT, 
	dashboard_conf TEXT NOT NULL, 
	"order" INTEGER, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN
);
CREATE TABLE IF NOT EXISTS dashboard_comment (
	dashboard_comment_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	dashboard_comment TEXT, 
	dashboard VARCHAR(200), 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN
);
CREATE TABLE IF NOT EXISTS etlx (
	etlx_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	etl VARCHAR(200) NOT NULL, 
	etl_desc TEXT, 
	attach_etlx_conf VARCHAR(200), 
	etlx_conf TEXT, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE (etl)
);
CREATE TABLE IF NOT EXISTS etlx_conf (
	etlx_conf_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	etlx_conf VARCHAR(200) NOT NULL, 
	etlx_conf_desc TEXT, 
	etlx_extra_conf TEXT, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE (etlx_conf)
);
CREATE TABLE IF NOT EXISTS manage_query (
	manage_query_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	manage_query VARCHAR(200) NOT NULL, 
	"database" VARCHAR(200) NOT NULL, 
	manage_query_conf TEXT, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN
);
CREATE TABLE IF NOT EXISTS notebook (
	notebook_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	notebook VARCHAR(200), 
	notebook_desc TEXT, 
	notebook_conf TEXT NOT NULL, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN
);
COMMIT;
