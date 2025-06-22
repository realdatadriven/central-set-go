BEGIN TRANSACTION;
DROP TABLE IF EXISTS "etl_rb_exp_dtail";
CREATE TABLE IF NOT EXISTS "etl_rb_exp_dtail" (
	"etl_rb_exp_dtail_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"etl_rb_exp_dtail"	VARCHAR(200) NOT NULL,
	"etl_rb_exp_dtail_desc"	TEXT,
	"etl_rbase_export_id"	INTEGER NOT NULL,
	"etl_report_base_id"	INTEGER NOT NULL,
	"sql_export_query"	TEXT,
	"database"	VARCHAR(200) NOT NULL,
	"dest_sheet_name"	VARCHAR(200),
	"dest_table_name"	VARCHAR(200),
	"etl_rb_exp_dtail_conf"	TEXT,
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id"),
	FOREIGN KEY("etl_rbase_export_id") REFERENCES "etl_rbase_export"("etl_rbase_export_id")
);
DROP TABLE IF EXISTS "etl_rbase_export";
CREATE TABLE IF NOT EXISTS "etl_rbase_export" (
	"etl_rbase_export_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"etl_rbase_export"	VARCHAR(200) NOT NULL,
	"etl_rbase_export_desc"	TEXT,
	"etl_report_base_id"	INTEGER NOT NULL,
	"export_type_id"	INTEGER NOT NULL,
	"attach_file_template"	VARCHAR(200) NOT NULL,
	"txt_fix_format_layout"	TEXT,
	"txt_fix_format_header"	TEXT,
	"etl_rbase_export_conf"	TEXT,
	"etl_rbase_output_id"	INTEGER,
	"database"	VARCHAR(200) NOT NULL,
	"active"	BOOLEAN,
	"ignore"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("etl_rbase_output_id") REFERENCES "etl_rbase_output"("etl_rbase_output_id"),
	FOREIGN KEY("export_type_id") REFERENCES "export_type"("export_type_id"),
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id")
);
DROP TABLE IF EXISTS "etl_rb_reconc_dtail";
CREATE TABLE IF NOT EXISTS "etl_rb_reconc_dtail" (
	"etl_rb_reconc_dtail_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"etl_rb_reconc_dtail"	VARCHAR(50) NOT NULL,
	"etl_rb_reconc_dtail_desc"	VARCHAR(200) NOT NULL,
	"sql_query_val_1"	TEXT,
	"sql_query_val_2"	TEXT,
	"is_eval_formula"	BOOLEAN,
	"sql_reconcilia_query"	TEXT,
	"comments"	TEXT,
	"etl_rb_reconcilia_id"	INTEGER NOT NULL,
	"etl_report_base_id"	INTEGER,
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id"),
	FOREIGN KEY("etl_rb_reconcilia_id") REFERENCES "etl_rb_reconcilia"("etl_rb_reconcilia_id")
);
DROP TABLE IF EXISTS "etl_rb_output_field";
CREATE TABLE IF NOT EXISTS "etl_rb_output_field" (
	"etl_rb_output_field_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"etl_rb_output_field"	VARCHAR(200) NOT NULL,
	"etl_rb_output_field_desc"	TEXT,
	"etl_rbase_output_id"	INTEGER NOT NULL,
	"sql_select"	TEXT NOT NULL,
	"sql_from"	TEXT,
	"sql_join"	TEXT,
	"sql_where"	TEXT,
	"sql_group_by"	TEXT,
	"sql_order_by"	TEXT,
	"sql_window"	TEXT,
	"sql_having"	TEXT,
	"field_order"	INTEGER,
	"fields_used"	TEXT,
	"etl_report_base_id"	INTEGER,
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id"),
	FOREIGN KEY("etl_rbase_output_id") REFERENCES "etl_rbase_output"("etl_rbase_output_id")
);
DROP TABLE IF EXISTS "etl_rbase_script";
CREATE TABLE IF NOT EXISTS "etl_rbase_script" (
	"script_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"script"	VARCHAR(200) NOT NULL,
	"script_sql"	TEXT,
	"script_conf"	TEXT,
	"active"	BOOLEAN,
	"etl_report_base_id"	INTEGER,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id")
);
DROP TABLE IF EXISTS "etl_rbase_backup";
CREATE TABLE IF NOT EXISTS "etl_rbase_backup" (
	"backup_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"backup"	VARCHAR(200) NOT NULL,
	"backup_sql"	TEXT,
	"backup_copy_to"	BOOLEAN,
	"backup_copy_path"	VARCHAR(200),
	"backup_conf"	TEXT,
	"active"	BOOLEAN,
	"etl_report_base_id"	INTEGER,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id")
);
DROP TABLE IF EXISTS "etl_rbase_notify";
CREATE TABLE IF NOT EXISTS "etl_rbase_notify" (
	"notify_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"notify_subject"	VARCHAR(200) NOT NULL,
	"notify_body"	TEXT NOT NULL,
	"notify_to"	VARCHAR(200) NOT NULL,
	"notify_cc"	VARCHAR(200),
	"notify_attach_exports"	BOOLEAN,
	"notify_copy_exports_to"	BOOLEAN,
	"notify_copy_exports_path"	VARCHAR(200),
	"notify_conf"	TEXT,
	"send_email"	BOOLEAN,
	"active"	BOOLEAN,
	"etl_report_base_id"	INTEGER,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id")
);
DROP TABLE IF EXISTS "etl_report_base_log";
CREATE TABLE IF NOT EXISTS "etl_report_base_log" (
	"log_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"type"	VARCHAR(50) NOT NULL,
	"name"	VARCHAR(100) NOT NULL,
	"ref"	DATE,
	"start"	DATETIME,
	"end"	DATETIME,
	"timer"	VARCHAR(10),
	"success"	BOOLEAN,
	"msg"	TEXT,
	"num_rows"	INTEGER,
	"errors"	INTEGER,
	"fixes"	INTEGER,
	"fname"	VARCHAR(200),
	"html"	TEXT,
	"etl_report_base_id"	INTEGER,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id")
);
DROP TABLE IF EXISTS "etl_rb_reconcilia";
CREATE TABLE IF NOT EXISTS "etl_rb_reconcilia" (
	"etl_rb_reconcilia_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"etl_rb_reconcilia"	VARCHAR(200) NOT NULL,
	"etl_rb_reconcilia_desc"	TEXT,
	"etl_rb_reconc_template"	TEXT,
	"database"	VARCHAR(200) NOT NULL,
	"active"	BOOLEAN,
	"etl_report_base_id"	INTEGER NOT NULL,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id")
);
DROP TABLE IF EXISTS "etl_rbase_quality";
CREATE TABLE IF NOT EXISTS "etl_rbase_quality" (
	"etl_rbase_quality_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"etl_rbase_quality"	VARCHAR(200) NOT NULL,
	"etl_rbase_quality_desc"	TEXT,
	"etl_report_base_id"	INTEGER NOT NULL,
	"sql_quality_check"	TEXT NOT NULL,
	"sql_quality_fix"	TEXT,
	"comments"	TEXT,
	"fields"	TEXT,
	"tables"	TEXT,
	"database"	VARCHAR(200) NOT NULL,
	"etl_rbase_quality_conf"	TEXT,
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id")
);
DROP TABLE IF EXISTS "etl_rbase_output";
CREATE TABLE IF NOT EXISTS "etl_rbase_output" (
	"etl_rbase_output_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"etl_rbase_output"	VARCHAR(200) NOT NULL,
	"etl_rbase_output_desc"	TEXT,
	"etl_report_base_id"	INTEGER NOT NULL,
	"output_type_id"	INTEGER,
	"date_field"	VARCHAR(200),
	"date_field_format"	VARCHAR(200),
	"destination_table"	VARCHAR(200) NOT NULL,
	"database"	VARCHAR(200) NOT NULL,
	"append_it"	BOOLEAN,
	"output_order"	INTEGER,
	"etl_rbase_output_conf"	TEXT,
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("output_type_id") REFERENCES "output_type"("output_type_id"),
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id")
);
DROP TABLE IF EXISTS "etl_rbase_input";
CREATE TABLE IF NOT EXISTS "etl_rbase_input" (
	"etl_rbase_input_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"etl_rbase_input"	VARCHAR(200) NOT NULL,
	"etl_rbase_input_desc"	TEXT,
	"etl_report_base_id"	INTEGER NOT NULL,
	"input_type_id"	INTEGER,
	"save_only_temp"	BOOLEAN,
	"replace_existing_data"	BOOLEAN,
	"check_ref_date"	BOOLEAN,
	"ref_date_field"	VARCHAR(200),
	"date_format_org"	VARCHAR(200),
	"other_date_fields"	VARCHAR(200),
	"ref_id_keys"	VARCHAR(200),
	"last_update_date_field"	VARCHAR(200),
	"incremental_extract"	BOOLEAN,
	"destination_table"	VARCHAR(200) NOT NULL,
	"database"	VARCHAR(200) NOT NULL,
	"allow_import"	BOOLEAN,
	"multiple_sheets"	BOOLEAN,
	"specific_sheets"	VARCHAR(200),
	"specific_range"	VARCHAR(200),
	"columns_to_import"	VARCHAR(200),
	"txt_fix_format_layout"	VARCHAR(200),
	"headers"	VARCHAR(200),
	"spreadsheet_forms"	BOOLEAN,
	"spreadsheet_forms_map"	VARCHAR(200),
	"etl_rbase_input_conf"	TEXT,
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("input_type_id") REFERENCES "input_type"("input_type_id"),
	FOREIGN KEY("etl_report_base_id") REFERENCES "etl_report_base"("etl_report_base_id")
);
DROP TABLE IF EXISTS "etl_report_base";
CREATE TABLE IF NOT EXISTS "etl_report_base" (
	"etl_report_base_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"etl_report_base"	VARCHAR(200) NOT NULL UNIQUE,
	"etl_report_base_desc"	TEXT,
	"attach_etl_rbase_doc"	VARCHAR(200),
	"periodicity_id"	INTEGER,
	"database"	VARCHAR(200) NOT NULL,
	"includes_output"	BOOLEAN,
	"includes_data_quality"	BOOLEAN,
	"includes_data_reconci"	BOOLEAN,
	"includes_exports"	BOOLEAN,
	"includes_backup"	BOOLEAN,
	"includes_notify"	BOOLEAN,
	"etl_report_base_conf"	TEXT,
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN,
	FOREIGN KEY("periodicity_id") REFERENCES "periodicity"("periodicity_id")
);
DROP TABLE IF EXISTS "dashboard_comment";
CREATE TABLE IF NOT EXISTS "dashboard_comment" (
	"dashboard_comment_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"dashboard_comment"	TEXT,
	"dashboard"	VARCHAR(200),
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
DROP TABLE IF EXISTS "dashboard";
CREATE TABLE IF NOT EXISTS "dashboard" (
	"dashboard_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"dashboard"	VARCHAR(200),
	"dashboard_desc"	TEXT,
	"dashboard_conf"	TEXT NOT NULL,
	"order"	INTEGER,
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
DROP TABLE IF EXISTS "manage_query";
CREATE TABLE IF NOT EXISTS "manage_query" (
	"manage_query_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"manage_query"	VARCHAR(200) NOT NULL,
	"database"	VARCHAR(200) NOT NULL,
	"manage_query_conf"	TEXT,
	"active"	BOOLEAN,
	"user_id"	INTEGER,
	"app_id"	INTEGER,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
DROP TABLE IF EXISTS "export_type";
CREATE TABLE IF NOT EXISTS "export_type" (
	"export_type_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"export_type"	VARCHAR(100) NOT NULL UNIQUE,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
DROP TABLE IF EXISTS "output_type";
CREATE TABLE IF NOT EXISTS "output_type" (
	"output_type_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"output_type"	VARCHAR(100) NOT NULL UNIQUE,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
DROP TABLE IF EXISTS "source_type";
CREATE TABLE IF NOT EXISTS "source_type" (
	"source_type_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"source_type"	VARCHAR(100) NOT NULL UNIQUE,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
DROP TABLE IF EXISTS "input_type";
CREATE TABLE IF NOT EXISTS "input_type" (
	"input_type_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"input_type"	VARCHAR(100) NOT NULL UNIQUE,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
DROP TABLE IF EXISTS "periodicity";
CREATE TABLE IF NOT EXISTS "periodicity" (
	"periodicity_id"	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
	"periodicity"	VARCHAR(100) NOT NULL UNIQUE,
	"created_at"	DATETIME,
	"updated_at"	DATETIME,
	"excluded"	BOOLEAN
);
INSERT INTO "export_type" ("export_type_id","export_type","created_at","updated_at","excluded") VALUES (1,'File','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0),
 (2,'Template','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0);
INSERT INTO "output_type" ("output_type_id","output_type","created_at","updated_at","excluded") VALUES (1,'Table','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0),
 (2,'View','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0);
INSERT INTO "source_type" ("source_type_id","source_type","created_at","updated_at","excluded") VALUES (1,'File','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0),
 (2,'Database','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0),
 (3,'FTP','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0),
 (4,'eMail','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0),
 (5,'FileSystem','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0);
INSERT INTO "input_type" ("input_type_id","input_type","created_at","updated_at","excluded") VALUES (1,'Aux','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0),
 (2,'Main','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0);
INSERT INTO "periodicity" ("periodicity_id","periodicity","created_at","updated_at","excluded") VALUES (1,'Daily','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0),
 (2,'Monthly','2025-02-27 11:45:28.906263','2025-02-27 11:45:28.906263',0);
COMMIT;
