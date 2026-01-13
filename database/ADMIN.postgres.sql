-- Adminer 5.3.0 PostgreSQL 18.0 dump

\connect "ADMIN";

DROP TABLE IF EXISTS "access_key";
DROP SEQUENCE IF EXISTS access_key_access_key_id_seq;
CREATE SEQUENCE access_key_access_key_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."access_key" (
    "access_key_id" integer DEFAULT nextval('access_key_access_key_id_seq') NOT NULL,
    "access_key_desc" character varying(200) NOT NULL,
    "access_token" text NOT NULL,
    "expires_at" timestamp,
    "active" boolean,
    "for_user_id" integer,
    "user_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "access_key_pkey" PRIMARY KEY ("access_key_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."access_key" IS 'Access Keys';

COMMENT ON COLUMN "public"."access_key"."access_key_id" IS 'Access Key ID';

COMMENT ON COLUMN "public"."access_key"."access_key_desc" IS 'Description';

COMMENT ON COLUMN "public"."access_key"."access_token" IS 'Token';

COMMENT ON COLUMN "public"."access_key"."expires_at" IS 'Expires at';

COMMENT ON COLUMN "public"."access_key"."active" IS 'Active';

COMMENT ON COLUMN "public"."access_key"."for_user_id" IS 'Created For';

COMMENT ON COLUMN "public"."access_key"."user_id" IS 'Created BY';

COMMENT ON COLUMN "public"."access_key"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."access_key"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."access_key"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "app";
DROP SEQUENCE IF EXISTS app_app_id_seq;
CREATE SEQUENCE app_app_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."app" (
    "app_id" integer DEFAULT nextval('app_app_id_seq') NOT NULL,
    "app" character varying(20) NOT NULL,
    "app_desc" text,
    "version" character varying(10) NOT NULL,
    "email" character varying(200),
    "db" character varying(200) NOT NULL,
    "attach_logo" character varying(200),
    "config" text,
    "user_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "app_pkey" PRIMARY KEY ("app_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."app" IS 'Applications';

COMMENT ON COLUMN "public"."app"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."app"."app" IS 'App Name';

COMMENT ON COLUMN "public"."app"."app_desc" IS 'Description';

COMMENT ON COLUMN "public"."app"."version" IS 'Version';

COMMENT ON COLUMN "public"."app"."email" IS 'Email';

COMMENT ON COLUMN "public"."app"."db" IS 'Database';

COMMENT ON COLUMN "public"."app"."attach_logo" IS 'Logo';

COMMENT ON COLUMN "public"."app"."config" IS 'Config';

COMMENT ON COLUMN "public"."app"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."app"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."app"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."app"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX app_app_key ON public.app USING btree (app);


DROP TABLE IF EXISTS "calendar";
DROP SEQUENCE IF EXISTS calendar_calendar_id_seq;
CREATE SEQUENCE calendar_calendar_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."calendar" (
    "calendar_id" integer DEFAULT nextval('calendar_calendar_id_seq') NOT NULL,
    "calendar" character varying(100) NOT NULL,
    "calendar_desc" text,
    "calendar_email" character varying(200),
    "calendar_color" character varying(50),
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "calendar_pkey" PRIMARY KEY ("calendar_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."calendar" IS 'Calendar';

COMMENT ON COLUMN "public"."calendar"."calendar_id" IS 'Calendar ID';

COMMENT ON COLUMN "public"."calendar"."calendar" IS 'Calendar';

COMMENT ON COLUMN "public"."calendar"."calendar_desc" IS 'Description';

COMMENT ON COLUMN "public"."calendar"."calendar_email" IS 'Email';

COMMENT ON COLUMN "public"."calendar"."calendar_color" IS 'Calendar Color';

COMMENT ON COLUMN "public"."calendar"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."calendar"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."calendar"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."calendar"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."calendar"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "column_level_access";
DROP SEQUENCE IF EXISTS column_level_access_column_level_access_id_seq;
CREATE SEQUENCE column_level_access_column_level_access_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."column_level_access" (
    "column_level_access_id" integer DEFAULT nextval('column_level_access_column_level_access_id_seq') NOT NULL,
    "column" integer NOT NULL,
    "table_id" integer,
    "table" character varying(200) NOT NULL,
    "db" character varying(200) NOT NULL,
    "user_id" integer,
    "app_id" integer,
    "create" boolean,
    "read" boolean,
    "update" boolean,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "column_level_access_pkey" PRIMARY KEY ("column_level_access_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."column_level_access" IS 'Column Level Access';

COMMENT ON COLUMN "public"."column_level_access"."column_level_access_id" IS 'Column Level Access ID';

COMMENT ON COLUMN "public"."column_level_access"."column" IS 'Column';

COMMENT ON COLUMN "public"."column_level_access"."table_id" IS 'Table ID';

COMMENT ON COLUMN "public"."column_level_access"."table" IS 'Table';

COMMENT ON COLUMN "public"."column_level_access"."db" IS 'Database';

COMMENT ON COLUMN "public"."column_level_access"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."column_level_access"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."column_level_access"."create" IS 'Create';

COMMENT ON COLUMN "public"."column_level_access"."read" IS 'Read';

COMMENT ON COLUMN "public"."column_level_access"."update" IS 'Update';

COMMENT ON COLUMN "public"."column_level_access"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."column_level_access"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."column_level_access"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "cron";
DROP SEQUENCE IF EXISTS cron_cron_id_seq;
CREATE SEQUENCE cron_cron_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."cron" (
    "cron_id" integer DEFAULT nextval('cron_cron_id_seq') NOT NULL,
    "cron" character varying(50) NOT NULL,
    "cron_desc" character varying(200) NOT NULL,
    "api" character varying(200) NOT NULL,
    "app_id" integer,
    "db" character varying(200),
    "table" character varying(50),
    "active" boolean,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "cron_pkey" PRIMARY KEY ("cron_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."cron" IS 'Jobs scheduling';

COMMENT ON COLUMN "public"."cron"."cron_id" IS 'Cron ID';

COMMENT ON COLUMN "public"."cron"."cron" IS 'Cron';

COMMENT ON COLUMN "public"."cron"."cron_desc" IS 'Decription';

COMMENT ON COLUMN "public"."cron"."api" IS 'API';

COMMENT ON COLUMN "public"."cron"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."cron"."db" IS 'Database';

COMMENT ON COLUMN "public"."cron"."table" IS 'Table';

COMMENT ON COLUMN "public"."cron"."active" IS 'Active';

COMMENT ON COLUMN "public"."cron"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."cron"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."cron"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "cron_log";
DROP SEQUENCE IF EXISTS cron_log_cron_log_id_seq;
CREATE SEQUENCE cron_log_cron_log_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."cron_log" (
    "cron_log_id" integer DEFAULT nextval('cron_log_cron_log_id_seq') NOT NULL,
    "cron_id" integer,
    "cron" character varying(50) NOT NULL,
    "cron_desc" character varying(200) NOT NULL,
    "api" character varying(200) NOT NULL,
    "start_at" timestamp,
    "end_at" timestamp,
    "success" boolean,
    "cron_msg" text NOT NULL,
    "app_id" integer,
    "db" character varying(200),
    "table" character varying(50),
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "cron_log_pkey" PRIMARY KEY ("cron_log_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."cron_log" IS 'Jobs scheduling logs';

COMMENT ON COLUMN "public"."cron_log"."cron_log_id" IS 'Cron Log ID';

COMMENT ON COLUMN "public"."cron_log"."cron_id" IS 'Cron ID';

COMMENT ON COLUMN "public"."cron_log"."cron" IS 'Cron';

COMMENT ON COLUMN "public"."cron_log"."cron_desc" IS 'Decription';

COMMENT ON COLUMN "public"."cron_log"."api" IS 'API';

COMMENT ON COLUMN "public"."cron_log"."start_at" IS 'Job Start';

COMMENT ON COLUMN "public"."cron_log"."end_at" IS 'Job End';

COMMENT ON COLUMN "public"."cron_log"."success" IS 'Success';

COMMENT ON COLUMN "public"."cron_log"."cron_msg" IS 'Message';

COMMENT ON COLUMN "public"."cron_log"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."cron_log"."db" IS 'Database';

COMMENT ON COLUMN "public"."cron_log"."table" IS 'Table';

COMMENT ON COLUMN "public"."cron_log"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."cron_log"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."cron_log"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "custom_form";
DROP SEQUENCE IF EXISTS custom_form_custom_form_id_seq;
CREATE SEQUENCE custom_form_custom_form_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."custom_form" (
    "custom_form_id" integer DEFAULT nextval('custom_form_custom_form_id_seq') NOT NULL,
    "table" character varying(200),
    "db" character varying(200),
    "config" text,
    "app_id" integer,
    "user_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "custom_form_pkey" PRIMARY KEY ("custom_form_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."custom_form" IS 'Custom Form';

COMMENT ON COLUMN "public"."custom_form"."custom_form_id" IS 'Custom Form ID';

COMMENT ON COLUMN "public"."custom_form"."table" IS 'Table';

COMMENT ON COLUMN "public"."custom_form"."db" IS 'Database';

COMMENT ON COLUMN "public"."custom_form"."config" IS 'Config';

COMMENT ON COLUMN "public"."custom_form"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."custom_form"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."custom_form"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."custom_form"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."custom_form"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "custom_table";
DROP SEQUENCE IF EXISTS custom_table_custom_table_id_seq;
CREATE SEQUENCE custom_table_custom_table_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."custom_table" (
    "custom_table_id" integer DEFAULT nextval('custom_table_custom_table_id_seq') NOT NULL,
    "table" character varying(200),
    "db" character varying(200),
    "config" text,
    "app_id" integer,
    "user_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "custom_table_pkey" PRIMARY KEY ("custom_table_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."custom_table" IS 'Custom Table';

COMMENT ON COLUMN "public"."custom_table"."custom_table_id" IS 'Custom Table ID';

COMMENT ON COLUMN "public"."custom_table"."table" IS 'Table';

COMMENT ON COLUMN "public"."custom_table"."db" IS 'Database';

COMMENT ON COLUMN "public"."custom_table"."config" IS 'Config';

COMMENT ON COLUMN "public"."custom_table"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."custom_table"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."custom_table"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."custom_table"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."custom_table"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "dashboard";
DROP SEQUENCE IF EXISTS dashboard_dashboard_id_seq;
CREATE SEQUENCE dashboard_dashboard_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."dashboard" (
    "dashboard_id" integer DEFAULT nextval('dashboard_dashboard_id_seq') NOT NULL,
    "dashboard" character varying(200),
    "dashboard_desc" text,
    "dashboard_conf" text NOT NULL,
    "order" integer,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "dashboard_pkey" PRIMARY KEY ("dashboard_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."dashboard" IS 'Dashboards';

COMMENT ON COLUMN "public"."dashboard"."dashboard_id" IS 'Dashboard ID';

COMMENT ON COLUMN "public"."dashboard"."dashboard" IS 'Dashboard';

COMMENT ON COLUMN "public"."dashboard"."dashboard_desc" IS 'Description';

COMMENT ON COLUMN "public"."dashboard"."dashboard_conf" IS 'Conf / Params';

COMMENT ON COLUMN "public"."dashboard"."order" IS 'Order';

COMMENT ON COLUMN "public"."dashboard"."active" IS 'Active';

COMMENT ON COLUMN "public"."dashboard"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."dashboard"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."dashboard"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."dashboard"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."dashboard"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "dashboard_comment";
DROP SEQUENCE IF EXISTS dashboard_comment_dashboard_comment_id_seq;
CREATE SEQUENCE dashboard_comment_dashboard_comment_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."dashboard_comment" (
    "dashboard_comment_id" integer DEFAULT nextval('dashboard_comment_dashboard_comment_id_seq') NOT NULL,
    "dashboard_comment" text,
    "dashboard" character varying(200),
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "dashboard_comment_pkey" PRIMARY KEY ("dashboard_comment_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."dashboard_comment" IS 'Dashboards Comments';

COMMENT ON COLUMN "public"."dashboard_comment"."dashboard_comment_id" IS 'Comment ID';

COMMENT ON COLUMN "public"."dashboard_comment"."dashboard_comment" IS 'Comments';

COMMENT ON COLUMN "public"."dashboard_comment"."dashboard" IS 'Dashboard';

COMMENT ON COLUMN "public"."dashboard_comment"."active" IS 'Active';

COMMENT ON COLUMN "public"."dashboard_comment"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."dashboard_comment"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."dashboard_comment"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."dashboard_comment"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."dashboard_comment"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_rb_exp_dtail";
DROP SEQUENCE IF EXISTS etl_rb_exp_dtail_etl_rb_exp_dtail_id_seq;
CREATE SEQUENCE etl_rb_exp_dtail_etl_rb_exp_dtail_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_rb_exp_dtail" (
    "etl_rb_exp_dtail_id" integer DEFAULT nextval('etl_rb_exp_dtail_etl_rb_exp_dtail_id_seq') NOT NULL,
    "etl_rb_exp_dtail" character varying(200) NOT NULL,
    "etl_rb_exp_dtail_desc" text,
    "etl_rbase_export_id" integer NOT NULL,
    "etl_report_base_id" integer NOT NULL,
    "sql_export_query" text,
    "database" character varying(200) NOT NULL,
    "dest_sheet_name" character varying(200),
    "dest_table_name" character varying(200),
    "etl_rb_exp_dtail_conf" text,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_rb_exp_dtail_pkey" PRIMARY KEY ("etl_rb_exp_dtail_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_rb_exp_dtail" IS 'Export Details';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."etl_rb_exp_dtail_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."etl_rb_exp_dtail" IS 'Export Detail';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."etl_rb_exp_dtail_desc" IS 'Export Detail Description';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."etl_rbase_export_id" IS 'Export ID';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."sql_export_query" IS 'Export SQL Query';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."database" IS 'Database';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."dest_sheet_name" IS 'Dest. Sheet Name';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."dest_table_name" IS 'Dest. Table Name';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."etl_rb_exp_dtail_conf" IS 'Config';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_rb_exp_dtail"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_rb_output_field";
DROP SEQUENCE IF EXISTS etl_rb_output_field_etl_rb_output_field_id_seq;
CREATE SEQUENCE etl_rb_output_field_etl_rb_output_field_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_rb_output_field" (
    "etl_rb_output_field_id" integer DEFAULT nextval('etl_rb_output_field_etl_rb_output_field_id_seq') NOT NULL,
    "etl_rb_output_field" character varying(200) NOT NULL,
    "etl_rb_output_field_desc" text,
    "etl_rbase_output_id" integer NOT NULL,
    "sql_select" text NOT NULL,
    "sql_from" text,
    "sql_join" text,
    "sql_where" text,
    "sql_group_by" text,
    "sql_order_by" text,
    "sql_window" text,
    "sql_having" text,
    "field_order" integer,
    "fields_used" text,
    "etl_report_base_id" integer,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_rb_output_field_pkey" PRIMARY KEY ("etl_rb_output_field_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_rb_output_field" IS 'Output Fields';

COMMENT ON COLUMN "public"."etl_rb_output_field"."etl_rb_output_field_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_rb_output_field"."etl_rb_output_field" IS 'Field';

COMMENT ON COLUMN "public"."etl_rb_output_field"."etl_rb_output_field_desc" IS 'Field Description';

COMMENT ON COLUMN "public"."etl_rb_output_field"."etl_rbase_output_id" IS 'Output ID';

COMMENT ON COLUMN "public"."etl_rb_output_field"."sql_select" IS 'SELECT';

COMMENT ON COLUMN "public"."etl_rb_output_field"."sql_from" IS 'FROM';

COMMENT ON COLUMN "public"."etl_rb_output_field"."sql_join" IS 'JOIN';

COMMENT ON COLUMN "public"."etl_rb_output_field"."sql_where" IS 'WHERE';

COMMENT ON COLUMN "public"."etl_rb_output_field"."sql_group_by" IS 'GROUP BY';

COMMENT ON COLUMN "public"."etl_rb_output_field"."sql_order_by" IS 'ORDER BY';

COMMENT ON COLUMN "public"."etl_rb_output_field"."sql_window" IS 'WINDOW';

COMMENT ON COLUMN "public"."etl_rb_output_field"."sql_having" IS 'HAVING';

COMMENT ON COLUMN "public"."etl_rb_output_field"."field_order" IS 'Field Order';

COMMENT ON COLUMN "public"."etl_rb_output_field"."fields_used" IS 'Fields Used';

COMMENT ON COLUMN "public"."etl_rb_output_field"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_rb_output_field"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_rb_output_field"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_rb_output_field"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_rb_output_field"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_rb_output_field"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_rb_output_field"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_rb_reconc_dtail";
DROP SEQUENCE IF EXISTS etl_rb_reconc_dtail_etl_rb_reconc_dtail_id_seq;
CREATE SEQUENCE etl_rb_reconc_dtail_etl_rb_reconc_dtail_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_rb_reconc_dtail" (
    "etl_rb_reconc_dtail_id" integer DEFAULT nextval('etl_rb_reconc_dtail_etl_rb_reconc_dtail_id_seq') NOT NULL,
    "etl_rb_reconc_dtail" character varying(50) NOT NULL,
    "etl_rb_reconc_dtail_desc" character varying(200) NOT NULL,
    "sql_query_val_1" text,
    "sql_query_val_2" text,
    "is_eval_formula" boolean,
    "sql_reconcilia_query" text,
    "comments" text,
    "etl_rb_reconcilia_id" integer NOT NULL,
    "etl_report_base_id" integer,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_rb_reconc_dtail_pkey" PRIMARY KEY ("etl_rb_reconc_dtail_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_rb_reconc_dtail" IS 'Data Reconciliation';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."etl_rb_reconc_dtail_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."etl_rb_reconc_dtail" IS 'Variable Name';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."etl_rb_reconc_dtail_desc" IS 'Var. Description';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."sql_query_val_1" IS 'SQL Query Valor 1';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."sql_query_val_2" IS 'SQL Query Valor 2';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."is_eval_formula" IS 'Is Eval Formula';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."sql_reconcilia_query" IS 'SQL / Formula Reconc';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."comments" IS 'Comments / Justifications';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."etl_rb_reconcilia_id" IS 'Data Reconciliation ID';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_rb_reconc_dtail"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_rb_reconcilia";
DROP SEQUENCE IF EXISTS etl_rb_reconcilia_etl_rb_reconcilia_id_seq;
CREATE SEQUENCE etl_rb_reconcilia_etl_rb_reconcilia_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_rb_reconcilia" (
    "etl_rb_reconcilia_id" integer DEFAULT nextval('etl_rb_reconcilia_etl_rb_reconcilia_id_seq') NOT NULL,
    "etl_rb_reconcilia" character varying(200) NOT NULL,
    "etl_rb_reconcilia_desc" text,
    "etl_rb_reconc_template" text,
    "database" character varying(200) NOT NULL,
    "active" boolean,
    "etl_report_base_id" integer NOT NULL,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_rb_reconcilia_pkey" PRIMARY KEY ("etl_rb_reconcilia_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_rb_reconcilia" IS 'Data Reconciliation';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."etl_rb_reconcilia_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."etl_rb_reconcilia" IS 'Data Reconciliation Rule';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."etl_rb_reconcilia_desc" IS 'Rule Description';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."etl_rb_reconc_template" IS 'Template';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."database" IS 'Database';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_rb_reconcilia"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_rbase_backup";
DROP SEQUENCE IF EXISTS etl_rbase_backup_backup_id_seq;
CREATE SEQUENCE etl_rbase_backup_backup_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_rbase_backup" (
    "backup_id" integer DEFAULT nextval('etl_rbase_backup_backup_id_seq') NOT NULL,
    "backup" character varying(200) NOT NULL,
    "backup_sql" text,
    "backup_copy_to" boolean,
    "backup_copy_path" character varying(200),
    "backup_conf" text,
    "active" boolean,
    "etl_report_base_id" integer,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_rbase_backup_pkey" PRIMARY KEY ("backup_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_rbase_backup" IS 'Backups';

COMMENT ON COLUMN "public"."etl_rbase_backup"."backup_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_rbase_backup"."backup" IS 'Backup';

COMMENT ON COLUMN "public"."etl_rbase_backup"."backup_sql" IS 'SQL';

COMMENT ON COLUMN "public"."etl_rbase_backup"."backup_copy_to" IS 'Copy';

COMMENT ON COLUMN "public"."etl_rbase_backup"."backup_copy_path" IS 'Path';

COMMENT ON COLUMN "public"."etl_rbase_backup"."backup_conf" IS 'Conf';

COMMENT ON COLUMN "public"."etl_rbase_backup"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_rbase_backup"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_rbase_backup"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_rbase_backup"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_rbase_backup"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_rbase_backup"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_rbase_backup"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_rbase_export";
DROP SEQUENCE IF EXISTS etl_rbase_export_etl_rbase_export_id_seq;
CREATE SEQUENCE etl_rbase_export_etl_rbase_export_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_rbase_export" (
    "etl_rbase_export_id" integer DEFAULT nextval('etl_rbase_export_etl_rbase_export_id_seq') NOT NULL,
    "etl_rbase_export" character varying(200) NOT NULL,
    "etl_rbase_export_desc" text,
    "etl_report_base_id" integer NOT NULL,
    "export_type_id" integer NOT NULL,
    "attach_file_template" character varying(200) NOT NULL,
    "txt_fix_format_layout" text,
    "txt_fix_format_header" text,
    "etl_rbase_export_conf" text,
    "etl_rbase_output_id" integer,
    "database" character varying(200) NOT NULL,
    "active" boolean,
    "ignore" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_rbase_export_pkey" PRIMARY KEY ("etl_rbase_export_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_rbase_export" IS 'Exports';

COMMENT ON COLUMN "public"."etl_rbase_export"."etl_rbase_export_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_rbase_export"."etl_rbase_export" IS 'Export';

COMMENT ON COLUMN "public"."etl_rbase_export"."etl_rbase_export_desc" IS 'Export Description';

COMMENT ON COLUMN "public"."etl_rbase_export"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_rbase_export"."export_type_id" IS 'Export ID';

COMMENT ON COLUMN "public"."etl_rbase_export"."attach_file_template" IS 'File name / Template';

COMMENT ON COLUMN "public"."etl_rbase_export"."txt_fix_format_layout" IS 'Text with fixed format Layout';

COMMENT ON COLUMN "public"."etl_rbase_export"."txt_fix_format_header" IS 'Text with fixed format Headers';

COMMENT ON COLUMN "public"."etl_rbase_export"."etl_rbase_export_conf" IS 'Config';

COMMENT ON COLUMN "public"."etl_rbase_export"."etl_rbase_output_id" IS 'Output ID';

COMMENT ON COLUMN "public"."etl_rbase_export"."database" IS 'Database';

COMMENT ON COLUMN "public"."etl_rbase_export"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_rbase_export"."ignore" IS 'Ignore';

COMMENT ON COLUMN "public"."etl_rbase_export"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_rbase_export"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_rbase_export"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_rbase_export"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_rbase_export"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_rbase_input";
DROP SEQUENCE IF EXISTS etl_rbase_input_etl_rbase_input_id_seq;
CREATE SEQUENCE etl_rbase_input_etl_rbase_input_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_rbase_input" (
    "etl_rbase_input_id" integer DEFAULT nextval('etl_rbase_input_etl_rbase_input_id_seq') NOT NULL,
    "etl_rbase_input" character varying(200) NOT NULL,
    "etl_rbase_input_desc" text,
    "etl_report_base_id" integer NOT NULL,
    "input_type_id" integer,
    "save_only_temp" boolean,
    "replace_existing_data" boolean,
    "check_ref_date" boolean,
    "ref_date_field" character varying(200),
    "date_format_org" character varying(200),
    "other_date_fields" character varying(200),
    "ref_id_keys" character varying(200),
    "last_update_date_field" character varying(200),
    "incremental_extract" boolean,
    "destination_table" character varying(200) NOT NULL,
    "database" character varying(200) NOT NULL,
    "allow_import" boolean,
    "multiple_sheets" boolean,
    "specific_sheets" character varying(200),
    "specific_range" character varying(200),
    "columns_to_import" character varying(200),
    "txt_fix_format_layout" character varying(200),
    "headers" character varying(200),
    "spreadsheet_forms" boolean,
    "spreadsheet_forms_map" character varying(200),
    "etl_rbase_input_conf" text,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_rbase_input_pkey" PRIMARY KEY ("etl_rbase_input_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_rbase_input" IS 'Inputs';

COMMENT ON COLUMN "public"."etl_rbase_input"."etl_rbase_input_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_rbase_input"."etl_rbase_input" IS 'Input';

COMMENT ON COLUMN "public"."etl_rbase_input"."etl_rbase_input_desc" IS 'Input Description';

COMMENT ON COLUMN "public"."etl_rbase_input"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_rbase_input"."input_type_id" IS 'Input Type ID';

COMMENT ON COLUMN "public"."etl_rbase_input"."save_only_temp" IS 'Save only temp';

COMMENT ON COLUMN "public"."etl_rbase_input"."replace_existing_data" IS 'Replace existing data';

COMMENT ON COLUMN "public"."etl_rbase_input"."check_ref_date" IS 'Check ref. date';

COMMENT ON COLUMN "public"."etl_rbase_input"."ref_date_field" IS 'Ref. date field';

COMMENT ON COLUMN "public"."etl_rbase_input"."date_format_org" IS 'Date Format in origin';

COMMENT ON COLUMN "public"."etl_rbase_input"."other_date_fields" IS 'Other Date Fields';

COMMENT ON COLUMN "public"."etl_rbase_input"."ref_id_keys" IS 'Reference / Id Keys';

COMMENT ON COLUMN "public"."etl_rbase_input"."last_update_date_field" IS 'Last update date field';

COMMENT ON COLUMN "public"."etl_rbase_input"."incremental_extract" IS 'Incremental Extraction';

COMMENT ON COLUMN "public"."etl_rbase_input"."destination_table" IS 'Destination Table';

COMMENT ON COLUMN "public"."etl_rbase_input"."database" IS 'Database';

COMMENT ON COLUMN "public"."etl_rbase_input"."allow_import" IS 'Allow Import';

COMMENT ON COLUMN "public"."etl_rbase_input"."multiple_sheets" IS 'Multiple Sheets';

COMMENT ON COLUMN "public"."etl_rbase_input"."specific_sheets" IS 'Specific Sheets';

COMMENT ON COLUMN "public"."etl_rbase_input"."specific_range" IS 'Specific Range';

COMMENT ON COLUMN "public"."etl_rbase_input"."columns_to_import" IS 'Columns to Import';

COMMENT ON COLUMN "public"."etl_rbase_input"."txt_fix_format_layout" IS 'Text with fixed format Layout';

COMMENT ON COLUMN "public"."etl_rbase_input"."headers" IS 'Headers';

COMMENT ON COLUMN "public"."etl_rbase_input"."spreadsheet_forms" IS 'Spreadsheet Forms';

COMMENT ON COLUMN "public"."etl_rbase_input"."spreadsheet_forms_map" IS 'Spreadsheet Forms Map';

COMMENT ON COLUMN "public"."etl_rbase_input"."etl_rbase_input_conf" IS 'Config';

COMMENT ON COLUMN "public"."etl_rbase_input"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_rbase_input"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_rbase_input"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_rbase_input"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_rbase_input"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_rbase_input"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_rbase_notify";
DROP SEQUENCE IF EXISTS etl_rbase_notify_notify_id_seq;
CREATE SEQUENCE etl_rbase_notify_notify_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_rbase_notify" (
    "notify_id" integer DEFAULT nextval('etl_rbase_notify_notify_id_seq') NOT NULL,
    "notify_subject" character varying(200) NOT NULL,
    "notify_body" text NOT NULL,
    "notify_to" character varying(200) NOT NULL,
    "notify_cc" character varying(200),
    "notify_attach_exports" boolean,
    "notify_copy_exports_to" boolean,
    "notify_copy_exports_path" character varying(200),
    "notify_conf" text,
    "send_email" boolean,
    "active" boolean,
    "etl_report_base_id" integer,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_rbase_notify_pkey" PRIMARY KEY ("notify_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_rbase_notify" IS 'Notify';

COMMENT ON COLUMN "public"."etl_rbase_notify"."notify_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_rbase_notify"."notify_subject" IS 'Subject';

COMMENT ON COLUMN "public"."etl_rbase_notify"."notify_body" IS 'Body';

COMMENT ON COLUMN "public"."etl_rbase_notify"."notify_to" IS 'Mail TO';

COMMENT ON COLUMN "public"."etl_rbase_notify"."notify_cc" IS 'Mail CC';

COMMENT ON COLUMN "public"."etl_rbase_notify"."notify_attach_exports" IS 'Attach Exports';

COMMENT ON COLUMN "public"."etl_rbase_notify"."notify_copy_exports_to" IS 'Copy Exports';

COMMENT ON COLUMN "public"."etl_rbase_notify"."notify_copy_exports_path" IS 'Copy Exports Path';

COMMENT ON COLUMN "public"."etl_rbase_notify"."notify_conf" IS 'Conf';

COMMENT ON COLUMN "public"."etl_rbase_notify"."send_email" IS 'Send Email';

COMMENT ON COLUMN "public"."etl_rbase_notify"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_rbase_notify"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_rbase_notify"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_rbase_notify"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_rbase_notify"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_rbase_notify"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_rbase_notify"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_rbase_output";
DROP SEQUENCE IF EXISTS etl_rbase_output_etl_rbase_output_id_seq;
CREATE SEQUENCE etl_rbase_output_etl_rbase_output_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_rbase_output" (
    "etl_rbase_output_id" integer DEFAULT nextval('etl_rbase_output_etl_rbase_output_id_seq') NOT NULL,
    "etl_rbase_output" character varying(200) NOT NULL,
    "etl_rbase_output_desc" text,
    "etl_report_base_id" integer NOT NULL,
    "output_type_id" integer,
    "date_field" character varying(200),
    "date_field_format" character varying(200),
    "destination_table" character varying(200) NOT NULL,
    "database" character varying(200) NOT NULL,
    "append_it" boolean,
    "output_order" integer,
    "etl_rbase_output_conf" text,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_rbase_output_pkey" PRIMARY KEY ("etl_rbase_output_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_rbase_output" IS 'Outputs';

COMMENT ON COLUMN "public"."etl_rbase_output"."etl_rbase_output_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_rbase_output"."etl_rbase_output" IS 'Output';

COMMENT ON COLUMN "public"."etl_rbase_output"."etl_rbase_output_desc" IS 'Output Description';

COMMENT ON COLUMN "public"."etl_rbase_output"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_rbase_output"."output_type_id" IS 'Input Type ID';

COMMENT ON COLUMN "public"."etl_rbase_output"."date_field" IS 'Date field';

COMMENT ON COLUMN "public"."etl_rbase_output"."date_field_format" IS 'Date Format';

COMMENT ON COLUMN "public"."etl_rbase_output"."destination_table" IS 'Destination Table';

COMMENT ON COLUMN "public"."etl_rbase_output"."database" IS 'Database';

COMMENT ON COLUMN "public"."etl_rbase_output"."append_it" IS 'Append';

COMMENT ON COLUMN "public"."etl_rbase_output"."output_order" IS 'Field Order';

COMMENT ON COLUMN "public"."etl_rbase_output"."etl_rbase_output_conf" IS 'Config';

COMMENT ON COLUMN "public"."etl_rbase_output"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_rbase_output"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_rbase_output"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_rbase_output"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_rbase_output"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_rbase_output"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_rbase_quality";
DROP SEQUENCE IF EXISTS etl_rbase_quality_etl_rbase_quality_id_seq;
CREATE SEQUENCE etl_rbase_quality_etl_rbase_quality_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_rbase_quality" (
    "etl_rbase_quality_id" integer DEFAULT nextval('etl_rbase_quality_etl_rbase_quality_id_seq') NOT NULL,
    "etl_rbase_quality" character varying(200) NOT NULL,
    "etl_rbase_quality_desc" text,
    "etl_report_base_id" integer NOT NULL,
    "sql_quality_check" text NOT NULL,
    "sql_quality_fix" text,
    "comments" text,
    "fields" text,
    "tables" text,
    "database" character varying(200) NOT NULL,
    "etl_rbase_quality_conf" text,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_rbase_quality_pkey" PRIMARY KEY ("etl_rbase_quality_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_rbase_quality" IS 'Data Quality';

COMMENT ON COLUMN "public"."etl_rbase_quality"."etl_rbase_quality_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_rbase_quality"."etl_rbase_quality" IS 'Data Quality Rule';

COMMENT ON COLUMN "public"."etl_rbase_quality"."etl_rbase_quality_desc" IS 'Rule Description';

COMMENT ON COLUMN "public"."etl_rbase_quality"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_rbase_quality"."sql_quality_check" IS 'SQL Quality Check';

COMMENT ON COLUMN "public"."etl_rbase_quality"."sql_quality_fix" IS 'SQL Quality Fix';

COMMENT ON COLUMN "public"."etl_rbase_quality"."comments" IS 'Comments / Justifications';

COMMENT ON COLUMN "public"."etl_rbase_quality"."fields" IS 'Fields';

COMMENT ON COLUMN "public"."etl_rbase_quality"."tables" IS 'Tables';

COMMENT ON COLUMN "public"."etl_rbase_quality"."database" IS 'Database';

COMMENT ON COLUMN "public"."etl_rbase_quality"."etl_rbase_quality_conf" IS 'Rule Description';

COMMENT ON COLUMN "public"."etl_rbase_quality"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_rbase_quality"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_rbase_quality"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_rbase_quality"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_rbase_quality"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_rbase_quality"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_rbase_script";
DROP SEQUENCE IF EXISTS etl_rbase_script_script_id_seq;
CREATE SEQUENCE etl_rbase_script_script_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_rbase_script" (
    "script_id" integer DEFAULT nextval('etl_rbase_script_script_id_seq') NOT NULL,
    "script" character varying(200) NOT NULL,
    "script_sql" text,
    "script_conf" text,
    "active" boolean,
    "etl_report_base_id" integer,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_rbase_script_pkey" PRIMARY KEY ("script_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_rbase_script" IS 'Scripts';

COMMENT ON COLUMN "public"."etl_rbase_script"."script_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_rbase_script"."script" IS 'Script';

COMMENT ON COLUMN "public"."etl_rbase_script"."script_sql" IS 'SQL';

COMMENT ON COLUMN "public"."etl_rbase_script"."script_conf" IS 'Conf';

COMMENT ON COLUMN "public"."etl_rbase_script"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_rbase_script"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_rbase_script"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_rbase_script"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_rbase_script"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_rbase_script"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_rbase_script"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "etl_report_base";
DROP SEQUENCE IF EXISTS etl_report_base_etl_report_base_id_seq;
CREATE SEQUENCE etl_report_base_etl_report_base_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_report_base" (
    "etl_report_base_id" integer DEFAULT nextval('etl_report_base_etl_report_base_id_seq') NOT NULL,
    "etl_report_base" character varying(200) NOT NULL,
    "etl_report_base_desc" text,
    "attach_etl_rbase_doc" character varying(200),
    "periodicity_id" integer,
    "database" character varying(200) NOT NULL,
    "includes_output" boolean,
    "includes_data_quality" boolean,
    "includes_data_reconci" boolean,
    "includes_exports" boolean,
    "includes_backup" boolean,
    "includes_script" boolean,
    "includes_notify" boolean,
    "etl_report_base_conf" text,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_report_base_pkey" PRIMARY KEY ("etl_report_base_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_report_base" IS 'DB | ETL | Report | Quality';

COMMENT ON COLUMN "public"."etl_report_base"."etl_report_base_id" IS 'ID';

COMMENT ON COLUMN "public"."etl_report_base"."etl_report_base" IS 'DB | ETL | Report | Quality';

COMMENT ON COLUMN "public"."etl_report_base"."etl_report_base_desc" IS 'Description';

COMMENT ON COLUMN "public"."etl_report_base"."attach_etl_rbase_doc" IS 'Documentation';

COMMENT ON COLUMN "public"."etl_report_base"."periodicity_id" IS 'Periodicity';

COMMENT ON COLUMN "public"."etl_report_base"."database" IS 'Database';

COMMENT ON COLUMN "public"."etl_report_base"."includes_output" IS 'Includes Output';

COMMENT ON COLUMN "public"."etl_report_base"."includes_data_quality" IS 'Includes Data Quality';

COMMENT ON COLUMN "public"."etl_report_base"."includes_data_reconci" IS 'Includes Data Reconc.';

COMMENT ON COLUMN "public"."etl_report_base"."includes_exports" IS 'Includes Exports';

COMMENT ON COLUMN "public"."etl_report_base"."includes_backup" IS 'Includes Backup';

COMMENT ON COLUMN "public"."etl_report_base"."includes_script" IS 'Includes Scripts';

COMMENT ON COLUMN "public"."etl_report_base"."includes_notify" IS 'Includes Notefy';

COMMENT ON COLUMN "public"."etl_report_base"."etl_report_base_conf" IS 'Config';

COMMENT ON COLUMN "public"."etl_report_base"."active" IS 'Active';

COMMENT ON COLUMN "public"."etl_report_base"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_report_base"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_report_base"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_report_base"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_report_base"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX etl_report_base_etl_report_base_key ON public.etl_report_base USING btree (etl_report_base);


DROP TABLE IF EXISTS "etl_report_base_log";
DROP SEQUENCE IF EXISTS etl_report_base_log_log_id_seq;
CREATE SEQUENCE etl_report_base_log_log_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etl_report_base_log" (
    "log_id" integer DEFAULT nextval('etl_report_base_log_log_id_seq') NOT NULL,
    "type" character varying(50) NOT NULL,
    "name" character varying(100) NOT NULL,
    "ref" date,
    "start" timestamp,
    "end" timestamp,
    "timer" character varying(10),
    "success" boolean,
    "msg" text,
    "num_rows" integer,
    "errors" integer,
    "fixes" integer,
    "fname" character varying(200),
    "html" text,
    "etl_report_base_id" integer,
    "app_id" integer,
    "user_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etl_report_base_log_pkey" PRIMARY KEY ("log_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etl_report_base_log" IS 'Logs DB | ETL | Report | Quality';

COMMENT ON COLUMN "public"."etl_report_base_log"."log_id" IS 'ID Log';

COMMENT ON COLUMN "public"."etl_report_base_log"."type" IS 'Type';

COMMENT ON COLUMN "public"."etl_report_base_log"."name" IS 'Name';

COMMENT ON COLUMN "public"."etl_report_base_log"."ref" IS 'Ref. Date';

COMMENT ON COLUMN "public"."etl_report_base_log"."start" IS 'Started At';

COMMENT ON COLUMN "public"."etl_report_base_log"."end" IS 'End At';

COMMENT ON COLUMN "public"."etl_report_base_log"."timer" IS 'Duration';

COMMENT ON COLUMN "public"."etl_report_base_log"."success" IS 'Success';

COMMENT ON COLUMN "public"."etl_report_base_log"."msg" IS 'Response';

COMMENT ON COLUMN "public"."etl_report_base_log"."num_rows" IS 'Affected Rows';

COMMENT ON COLUMN "public"."etl_report_base_log"."errors" IS 'Errors';

COMMENT ON COLUMN "public"."etl_report_base_log"."fixes" IS 'Automated Fixes';

COMMENT ON COLUMN "public"."etl_report_base_log"."fname" IS 'Generated Files';

COMMENT ON COLUMN "public"."etl_report_base_log"."html" IS 'Html';

COMMENT ON COLUMN "public"."etl_report_base_log"."etl_report_base_id" IS 'DB | ETL | Report | Quality ID';

COMMENT ON COLUMN "public"."etl_report_base_log"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etl_report_base_log"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etl_report_base_log"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etl_report_base_log"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etl_report_base_log"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "export_type";
DROP SEQUENCE IF EXISTS export_type_export_type_id_seq;
CREATE SEQUENCE export_type_export_type_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."export_type" (
    "export_type_id" integer DEFAULT nextval('export_type_export_type_id_seq') NOT NULL,
    "export_type" character varying(100) NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "export_type_pkey" PRIMARY KEY ("export_type_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."export_type" IS 'Export Type';

COMMENT ON COLUMN "public"."export_type"."export_type_id" IS 'Export Type ID';

COMMENT ON COLUMN "public"."export_type"."export_type" IS 'Export Type';

COMMENT ON COLUMN "public"."export_type"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."export_type"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."export_type"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX export_type_export_type_key ON public.export_type USING btree (export_type);


DROP TABLE IF EXISTS "input_type";
DROP SEQUENCE IF EXISTS input_type_input_type_id_seq;
CREATE SEQUENCE input_type_input_type_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."input_type" (
    "input_type_id" integer DEFAULT nextval('input_type_input_type_id_seq') NOT NULL,
    "input_type" character varying(100) NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "input_type_pkey" PRIMARY KEY ("input_type_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."input_type" IS 'Input Type';

COMMENT ON COLUMN "public"."input_type"."input_type_id" IS 'Input Type ID';

COMMENT ON COLUMN "public"."input_type"."input_type" IS 'Input Type';

COMMENT ON COLUMN "public"."input_type"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."input_type"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."input_type"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX input_type_input_type_key ON public.input_type USING btree (input_type);


DROP TABLE IF EXISTS "lang";
DROP SEQUENCE IF EXISTS lang_lang_id_seq;
CREATE SEQUENCE lang_lang_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."lang" (
    "lang_id" integer DEFAULT nextval('lang_lang_id_seq') NOT NULL,
    "lang" character varying(4) NOT NULL,
    "lang_desc" character varying(200),
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "lang_pkey" PRIMARY KEY ("lang_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."lang" IS 'Languages';

COMMENT ON COLUMN "public"."lang"."lang_id" IS 'Lang ID';

COMMENT ON COLUMN "public"."lang"."lang" IS 'Language';

COMMENT ON COLUMN "public"."lang"."lang_desc" IS 'Description';

COMMENT ON COLUMN "public"."lang"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."lang"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."lang"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX lang_lang_key ON public.lang USING btree (lang);


DROP TABLE IF EXISTS "manage_query";
DROP SEQUENCE IF EXISTS manage_query_manage_query_id_seq;
CREATE SEQUENCE manage_query_manage_query_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."manage_query" (
    "manage_query_id" integer DEFAULT nextval('manage_query_manage_query_id_seq') NOT NULL,
    "manage_query" character varying(200) NOT NULL,
    "database" character varying(200) NOT NULL,
    "manage_query_conf" text,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "manage_query_pkey" PRIMARY KEY ("manage_query_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."manage_query" IS 'Queries';

COMMENT ON COLUMN "public"."manage_query"."manage_query_id" IS 'ID';

COMMENT ON COLUMN "public"."manage_query"."manage_query" IS 'Query Desc';

COMMENT ON COLUMN "public"."manage_query"."database" IS 'Database';

COMMENT ON COLUMN "public"."manage_query"."manage_query_conf" IS 'Query Config';

COMMENT ON COLUMN "public"."manage_query"."active" IS 'Active';

COMMENT ON COLUMN "public"."manage_query"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."manage_query"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."manage_query"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."manage_query"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."manage_query"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "menu";
DROP SEQUENCE IF EXISTS menu_menu_id_seq;
CREATE SEQUENCE menu_menu_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."menu" (
    "menu_id" integer DEFAULT nextval('menu_menu_id_seq') NOT NULL,
    "menu" character varying(200) NOT NULL,
    "menu_desc" text,
    "menu_icon" character varying(20),
    "menu_order" integer,
    "menu_config" text,
    "app_id" integer NOT NULL,
    "user_id" integer NOT NULL,
    "active" boolean,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "menu_pkey" PRIMARY KEY ("menu_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."menu" IS 'Menus';

COMMENT ON COLUMN "public"."menu"."menu_id" IS 'Menu ID';

COMMENT ON COLUMN "public"."menu"."menu" IS 'Menu';

COMMENT ON COLUMN "public"."menu"."menu_desc" IS 'Description';

COMMENT ON COLUMN "public"."menu"."menu_icon" IS 'Icon';

COMMENT ON COLUMN "public"."menu"."menu_order" IS 'Order';

COMMENT ON COLUMN "public"."menu"."menu_config" IS 'Config';

COMMENT ON COLUMN "public"."menu"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."menu"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."menu"."active" IS 'Active';

COMMENT ON COLUMN "public"."menu"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."menu"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."menu"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "menu_table";
DROP SEQUENCE IF EXISTS menu_table_menu_table_id_seq;
CREATE SEQUENCE menu_table_menu_table_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."menu_table" (
    "menu_table_id" integer DEFAULT nextval('menu_table_menu_table_id_seq') NOT NULL,
    "menu_id" integer,
    "table_id" integer,
    "app_id" integer,
    "user_id" integer,
    "active" boolean,
    "requires_rla" boolean,
    "menu_table_cnf" text,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "menu_table_pkey" PRIMARY KEY ("menu_table_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."menu_table" IS 'Menu Tables';

COMMENT ON COLUMN "public"."menu_table"."menu_table_id" IS 'Menu Table ID';

COMMENT ON COLUMN "public"."menu_table"."menu_id" IS 'Menu ID';

COMMENT ON COLUMN "public"."menu_table"."table_id" IS 'Table ID';

COMMENT ON COLUMN "public"."menu_table"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."menu_table"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."menu_table"."active" IS 'Active';

COMMENT ON COLUMN "public"."menu_table"."requires_rla" IS 'Requires Row Level Access';

COMMENT ON COLUMN "public"."menu_table"."menu_table_cnf" IS 'Config';

COMMENT ON COLUMN "public"."menu_table"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."menu_table"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."menu_table"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "output_type";
DROP SEQUENCE IF EXISTS output_type_output_type_id_seq;
CREATE SEQUENCE output_type_output_type_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."output_type" (
    "output_type_id" integer DEFAULT nextval('output_type_output_type_id_seq') NOT NULL,
    "output_type" character varying(100) NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "output_type_pkey" PRIMARY KEY ("output_type_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."output_type" IS 'Output Type';

COMMENT ON COLUMN "public"."output_type"."output_type_id" IS 'Output Type ID';

COMMENT ON COLUMN "public"."output_type"."output_type" IS 'Souce Type';

COMMENT ON COLUMN "public"."output_type"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."output_type"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."output_type"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX output_type_output_type_key ON public.output_type USING btree (output_type);


DROP TABLE IF EXISTS "periodicity";
DROP SEQUENCE IF EXISTS periodicity_periodicity_id_seq;
CREATE SEQUENCE periodicity_periodicity_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."periodicity" (
    "periodicity_id" integer DEFAULT nextval('periodicity_periodicity_id_seq') NOT NULL,
    "periodicity" character varying(100) NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "periodicity_pkey" PRIMARY KEY ("periodicity_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."periodicity" IS 'Periodicity';

COMMENT ON COLUMN "public"."periodicity"."periodicity_id" IS 'Periodicity ID';

COMMENT ON COLUMN "public"."periodicity"."periodicity" IS 'Periodicity';

COMMENT ON COLUMN "public"."periodicity"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."periodicity"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."periodicity"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX periodicity_periodicity_key ON public.periodicity USING btree (periodicity);


DROP TABLE IF EXISTS "repeat_type";
DROP SEQUENCE IF EXISTS repeat_type_repeat_type_id_seq;
CREATE SEQUENCE repeat_type_repeat_type_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."repeat_type" (
    "repeat_type_id" integer DEFAULT nextval('repeat_type_repeat_type_id_seq') NOT NULL,
    "repeat_type" character varying(10) NOT NULL,
    "excluded" boolean,
    CONSTRAINT "repeat_type_pkey" PRIMARY KEY ("repeat_type_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."repeat_type" IS 'Type of repetition';

COMMENT ON COLUMN "public"."repeat_type"."repeat_type_id" IS 'Type of repetition ID';

COMMENT ON COLUMN "public"."repeat_type"."repeat_type" IS 'Type of Repetition';

COMMENT ON COLUMN "public"."repeat_type"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "role";
DROP SEQUENCE IF EXISTS role_role_id_seq;
CREATE SEQUENCE role_role_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."role" (
    "role_id" integer DEFAULT nextval('role_role_id_seq') NOT NULL,
    "role" character varying(20) NOT NULL,
    "role_desc" text,
    "config" text,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "role_pkey" PRIMARY KEY ("role_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."role" IS 'Roles';

COMMENT ON COLUMN "public"."role"."role_id" IS 'Role ID';

COMMENT ON COLUMN "public"."role"."role" IS 'Role';

COMMENT ON COLUMN "public"."role"."role_desc" IS 'Description';

COMMENT ON COLUMN "public"."role"."config" IS 'Config';

COMMENT ON COLUMN "public"."role"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."role"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."role"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX role_role_key ON public.role USING btree (role);


DROP TABLE IF EXISTS "role_app";
DROP SEQUENCE IF EXISTS role_app_role_app_id_seq;
CREATE SEQUENCE role_app_role_app_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."role_app" (
    "role_app_id" integer DEFAULT nextval('role_app_role_app_id_seq') NOT NULL,
    "role_id" integer,
    "app_id" integer,
    "access" boolean,
    "user_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "role_app_pkey" PRIMARY KEY ("role_app_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."role_app" IS 'Role Apps';

COMMENT ON COLUMN "public"."role_app"."role_app_id" IS 'Role App ID';

COMMENT ON COLUMN "public"."role_app"."role_id" IS 'Role ID';

COMMENT ON COLUMN "public"."role_app"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."role_app"."access" IS 'Access';

COMMENT ON COLUMN "public"."role_app"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."role_app"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."role_app"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."role_app"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "role_app_menu";
DROP SEQUENCE IF EXISTS role_app_menu_role_app_menu_id_seq;
CREATE SEQUENCE role_app_menu_role_app_menu_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."role_app_menu" (
    "role_app_menu_id" integer DEFAULT nextval('role_app_menu_role_app_menu_id_seq') NOT NULL,
    "role_id" integer,
    "app_id" integer,
    "menu_id" integer,
    "access" boolean,
    "user_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "role_app_menu_pkey" PRIMARY KEY ("role_app_menu_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."role_app_menu" IS 'Role App Menus';

COMMENT ON COLUMN "public"."role_app_menu"."role_app_menu_id" IS 'Role App Menu ID';

COMMENT ON COLUMN "public"."role_app_menu"."role_id" IS 'Role ID';

COMMENT ON COLUMN "public"."role_app_menu"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."role_app_menu"."menu_id" IS 'Menu ID';

COMMENT ON COLUMN "public"."role_app_menu"."access" IS 'Access';

COMMENT ON COLUMN "public"."role_app_menu"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."role_app_menu"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."role_app_menu"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."role_app_menu"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "role_app_menu_table";
DROP SEQUENCE IF EXISTS role_app_menu_table_role_app_menu_table_id_seq;
CREATE SEQUENCE role_app_menu_table_role_app_menu_table_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."role_app_menu_table" (
    "role_app_menu_table_id" integer DEFAULT nextval('role_app_menu_table_role_app_menu_table_id_seq') NOT NULL,
    "role_id" integer,
    "app_id" integer,
    "menu_id" integer,
    "table_id" integer,
    "create" boolean,
    "read" boolean,
    "update" boolean,
    "delete" boolean,
    "user_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "role_app_menu_table_pkey" PRIMARY KEY ("role_app_menu_table_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."role_app_menu_table" IS 'Role App Menu Tables';

COMMENT ON COLUMN "public"."role_app_menu_table"."role_app_menu_table_id" IS 'Role App Menu Table ID';

COMMENT ON COLUMN "public"."role_app_menu_table"."role_id" IS 'Role ID';

COMMENT ON COLUMN "public"."role_app_menu_table"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."role_app_menu_table"."menu_id" IS 'Menu ID';

COMMENT ON COLUMN "public"."role_app_menu_table"."table_id" IS 'Table ID';

COMMENT ON COLUMN "public"."role_app_menu_table"."create" IS 'Create';

COMMENT ON COLUMN "public"."role_app_menu_table"."read" IS 'Read';

COMMENT ON COLUMN "public"."role_app_menu_table"."update" IS 'Update';

COMMENT ON COLUMN "public"."role_app_menu_table"."delete" IS 'Delete';

COMMENT ON COLUMN "public"."role_app_menu_table"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."role_app_menu_table"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."role_app_menu_table"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."role_app_menu_table"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "role_row_level_access";
DROP SEQUENCE IF EXISTS role_row_level_access_role_row_level_access_id_seq;
CREATE SEQUENCE role_row_level_access_role_row_level_access_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."role_row_level_access" (
    "role_row_level_access_id" integer DEFAULT nextval('role_row_level_access_role_row_level_access_id_seq') NOT NULL,
    "role_id" integer,
    "row_id" integer NOT NULL,
    "table_id" integer,
    "table" character varying(200) NOT NULL,
    "db" character varying(200) NOT NULL,
    "user_id" integer,
    "app_id" integer,
    "read" boolean,
    "update" boolean,
    "delete" boolean,
    "share" boolean,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "role_row_level_access_pkey" PRIMARY KEY ("role_row_level_access_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."role_row_level_access" IS 'Role Row Level Access';

COMMENT ON COLUMN "public"."role_row_level_access"."role_row_level_access_id" IS 'Role Row Level Access ID';

COMMENT ON COLUMN "public"."role_row_level_access"."role_id" IS 'Role ID';

COMMENT ON COLUMN "public"."role_row_level_access"."row_id" IS 'Row ID';

COMMENT ON COLUMN "public"."role_row_level_access"."table_id" IS 'Table ID';

COMMENT ON COLUMN "public"."role_row_level_access"."table" IS 'Table';

COMMENT ON COLUMN "public"."role_row_level_access"."db" IS 'Database';

COMMENT ON COLUMN "public"."role_row_level_access"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."role_row_level_access"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."role_row_level_access"."read" IS 'Read';

COMMENT ON COLUMN "public"."role_row_level_access"."update" IS 'Update';

COMMENT ON COLUMN "public"."role_row_level_access"."delete" IS 'Delete';

COMMENT ON COLUMN "public"."role_row_level_access"."share" IS 'Share';

COMMENT ON COLUMN "public"."role_row_level_access"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."role_row_level_access"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."role_row_level_access"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "row_level_access";
DROP SEQUENCE IF EXISTS row_level_access_row_level_access_id_seq;
CREATE SEQUENCE row_level_access_row_level_access_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."row_level_access" (
    "row_level_access_id" integer DEFAULT nextval('row_level_access_row_level_access_id_seq') NOT NULL,
    "row_id" integer NOT NULL,
    "table_id" integer,
    "table" character varying(200) NOT NULL,
    "db" character varying(200) NOT NULL,
    "user_id" integer,
    "app_id" integer,
    "read" boolean,
    "update" boolean,
    "delete" boolean,
    "share" boolean,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "row_level_access_pkey" PRIMARY KEY ("row_level_access_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."row_level_access" IS 'Row Level Access';

COMMENT ON COLUMN "public"."row_level_access"."row_level_access_id" IS 'Row Level Access ID';

COMMENT ON COLUMN "public"."row_level_access"."row_id" IS 'Row ID';

COMMENT ON COLUMN "public"."row_level_access"."table_id" IS 'Table ID';

COMMENT ON COLUMN "public"."row_level_access"."table" IS 'Table';

COMMENT ON COLUMN "public"."row_level_access"."db" IS 'Database';

COMMENT ON COLUMN "public"."row_level_access"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."row_level_access"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."row_level_access"."read" IS 'Read';

COMMENT ON COLUMN "public"."row_level_access"."update" IS 'Update';

COMMENT ON COLUMN "public"."row_level_access"."delete" IS 'Delete';

COMMENT ON COLUMN "public"."row_level_access"."share" IS 'Share';

COMMENT ON COLUMN "public"."row_level_access"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."row_level_access"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."row_level_access"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "source_type";
DROP SEQUENCE IF EXISTS source_type_source_type_id_seq;
CREATE SEQUENCE source_type_source_type_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."source_type" (
    "source_type_id" integer DEFAULT nextval('source_type_source_type_id_seq') NOT NULL,
    "source_type" character varying(100) NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "source_type_pkey" PRIMARY KEY ("source_type_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."source_type" IS 'Source Type';

COMMENT ON COLUMN "public"."source_type"."source_type_id" IS 'Source Type ID';

COMMENT ON COLUMN "public"."source_type"."source_type" IS 'Souce Type';

COMMENT ON COLUMN "public"."source_type"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."source_type"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."source_type"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX source_type_source_type_key ON public.source_type USING btree (source_type);


DROP TABLE IF EXISTS "table";
DROP SEQUENCE IF EXISTS table_table_id_seq;
CREATE SEQUENCE table_table_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."table" (
    "table_id" integer DEFAULT nextval('table_table_id_seq') NOT NULL,
    "table" character varying(50) NOT NULL,
    "table_desc" character varying(200),
    "db" character varying(50),
    "requires_rla" boolean,
    "user_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "table_pkey" PRIMARY KEY ("table_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."table" IS 'Tables';

COMMENT ON COLUMN "public"."table"."table_id" IS 'Table ID';

COMMENT ON COLUMN "public"."table"."table" IS 'Table';

COMMENT ON COLUMN "public"."table"."table_desc" IS 'Description';

COMMENT ON COLUMN "public"."table"."db" IS 'Database';

COMMENT ON COLUMN "public"."table"."requires_rla" IS 'Requires Row Level Access';

COMMENT ON COLUMN "public"."table"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."table"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."table"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."table"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "table_schema";
DROP SEQUENCE IF EXISTS table_schema_table_schema_id_seq;
CREATE SEQUENCE table_schema_table_schema_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."table_schema" (
    "table_schema_id" integer DEFAULT nextval('table_schema_table_schema_id_seq') NOT NULL,
    "db" character varying(200) NOT NULL,
    "table" character varying(200) NOT NULL,
    "field" character varying(200) NOT NULL,
    "type" character varying(200) NOT NULL,
    "comment" character varying(200),
    "pk" boolean,
    "autoincrement" boolean,
    "nullable" boolean,
    "computed" boolean,
    "default" boolean,
    "fk" boolean,
    "referred_table" character varying(200),
    "referred_column" character varying(200),
    "user_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "table_schema_pkey" PRIMARY KEY ("table_schema_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."table_schema" IS 'Table Schema';

COMMENT ON COLUMN "public"."table_schema"."table_schema_id" IS 'Table field ID';

COMMENT ON COLUMN "public"."table_schema"."db" IS 'Database';

COMMENT ON COLUMN "public"."table_schema"."table" IS 'Table';

COMMENT ON COLUMN "public"."table_schema"."field" IS 'Field';

COMMENT ON COLUMN "public"."table_schema"."type" IS 'Type';

COMMENT ON COLUMN "public"."table_schema"."comment" IS 'Comment';

COMMENT ON COLUMN "public"."table_schema"."pk" IS 'Primary Key';

COMMENT ON COLUMN "public"."table_schema"."autoincrement" IS 'Auto Increment';

COMMENT ON COLUMN "public"."table_schema"."nullable" IS 'Nullable';

COMMENT ON COLUMN "public"."table_schema"."computed" IS 'Nullable';

COMMENT ON COLUMN "public"."table_schema"."default" IS 'Default';

COMMENT ON COLUMN "public"."table_schema"."fk" IS 'Foreign  Key';

COMMENT ON COLUMN "public"."table_schema"."referred_table" IS 'Ref. Table.';

COMMENT ON COLUMN "public"."table_schema"."referred_column" IS 'Ref. Column';

COMMENT ON COLUMN "public"."table_schema"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."table_schema"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."table_schema"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."table_schema"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "task";
DROP SEQUENCE IF EXISTS task_task_id_seq;
CREATE SEQUENCE task_task_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."task" (
    "task_id" integer DEFAULT nextval('task_task_id_seq') NOT NULL,
    "task" character varying(100) NOT NULL,
    "task_desc" text,
    "starts_at" timestamp NOT NULL,
    "ends_at" timestamp,
    "calendar_id" integer,
    "calendar_color" character varying(50),
    "calendar_email" character varying(200),
    "task_status_id" integer,
    "task_status" character varying(50),
    "attach_task" character varying(200),
    "repeat" boolean,
    "repeat_type_id" integer,
    "days_of_week" character varying(200),
    "repeat_start_date" date,
    "repeat_start_time" time without time zone,
    "repeat_end_date" date,
    "repeat_end_time" time without time zone,
    "attendees" text,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "task_pkey" PRIMARY KEY ("task_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."task" IS 'Tasks';

COMMENT ON COLUMN "public"."task"."task_id" IS 'Task Id';

COMMENT ON COLUMN "public"."task"."task" IS 'task';

COMMENT ON COLUMN "public"."task"."task_desc" IS 'Description';

COMMENT ON COLUMN "public"."task"."starts_at" IS 'Starts at';

COMMENT ON COLUMN "public"."task"."ends_at" IS 'Ends at';

COMMENT ON COLUMN "public"."task"."calendar_id" IS 'Calendar ID';

COMMENT ON COLUMN "public"."task"."calendar_color" IS 'Calendar Color';

COMMENT ON COLUMN "public"."task"."calendar_email" IS 'Email';

COMMENT ON COLUMN "public"."task"."task_status_id" IS 'Status ID';

COMMENT ON COLUMN "public"."task"."task_status" IS 'Status';

COMMENT ON COLUMN "public"."task"."attach_task" IS 'Attachment';

COMMENT ON COLUMN "public"."task"."repeat" IS 'Repeat';

COMMENT ON COLUMN "public"."task"."repeat_type_id" IS 'Repeat Type ID';

COMMENT ON COLUMN "public"."task"."days_of_week" IS 'Days of week';

COMMENT ON COLUMN "public"."task"."repeat_start_date" IS 'Repeat start date';

COMMENT ON COLUMN "public"."task"."repeat_start_time" IS 'Repeat start time';

COMMENT ON COLUMN "public"."task"."repeat_end_date" IS 'Repeat end date';

COMMENT ON COLUMN "public"."task"."repeat_end_time" IS 'Repeat end time';

COMMENT ON COLUMN "public"."task"."attendees" IS 'Attendees';

COMMENT ON COLUMN "public"."task"."active" IS 'Active';

COMMENT ON COLUMN "public"."task"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."task"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."task"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."task"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."task"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "task_status";
DROP SEQUENCE IF EXISTS task_status_task_status_id_seq;
CREATE SEQUENCE task_status_task_status_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."task_status" (
    "task_status_id" integer DEFAULT nextval('task_status_task_status_id_seq') NOT NULL,
    "task_status" character varying(10) NOT NULL,
    "excluded" boolean,
    CONSTRAINT "task_status_pkey" PRIMARY KEY ("task_status_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."task_status" IS 'Status';

COMMENT ON COLUMN "public"."task_status"."task_status_id" IS 'Status ID';

COMMENT ON COLUMN "public"."task_status"."task_status" IS 'Status';

COMMENT ON COLUMN "public"."task_status"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "task_track";
DROP SEQUENCE IF EXISTS task_track_task_track_id_seq;
CREATE SEQUENCE task_track_task_track_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."task_track" (
    "task_track_id" integer DEFAULT nextval('task_track_task_track_id_seq') NOT NULL,
    "task_track" character varying(100) NOT NULL,
    "task_track_desc" text,
    "task_track_date" timestamp NOT NULL,
    "task_status_id" integer,
    "attach_task_track" character varying(200),
    "task_id" integer,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "task_track_pkey" PRIMARY KEY ("task_track_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."task_track" IS 'Status Updates';

COMMENT ON COLUMN "public"."task_track"."task_track_id" IS 'Status Update ID';

COMMENT ON COLUMN "public"."task_track"."task_track" IS 'Status Update';

COMMENT ON COLUMN "public"."task_track"."task_track_desc" IS 'Description';

COMMENT ON COLUMN "public"."task_track"."task_track_date" IS 'Date';

COMMENT ON COLUMN "public"."task_track"."task_status_id" IS 'Status ID';

COMMENT ON COLUMN "public"."task_track"."attach_task_track" IS 'Attachment';

COMMENT ON COLUMN "public"."task_track"."task_id" IS 'Task ID';

COMMENT ON COLUMN "public"."task_track"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."task_track"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."task_track"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."task_track"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."task_track"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "translate_table";
DROP SEQUENCE IF EXISTS translate_table_transl_tbl_id_seq;
CREATE SEQUENCE translate_table_transl_tbl_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."translate_table" (
    "transl_tbl_id" integer DEFAULT nextval('translate_table_transl_tbl_id_seq') NOT NULL,
    "table_org_desc" character varying(200) NOT NULL,
    "table_transl_desc" character varying(200) NOT NULL,
    "table" character varying(200) NOT NULL,
    "db" character varying(200) NOT NULL,
    "lang" character varying(5) NOT NULL,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "translate_table_pkey" PRIMARY KEY ("transl_tbl_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."translate_table" IS 'Translate Table';

COMMENT ON COLUMN "public"."translate_table"."transl_tbl_id" IS 'Translate Table ID';

COMMENT ON COLUMN "public"."translate_table"."table_org_desc" IS 'Table Org. Desc';

COMMENT ON COLUMN "public"."translate_table"."table_transl_desc" IS 'Table Transl. Desc';

COMMENT ON COLUMN "public"."translate_table"."table" IS 'Table';

COMMENT ON COLUMN "public"."translate_table"."db" IS 'Database';

COMMENT ON COLUMN "public"."translate_table"."lang" IS 'Lang';

COMMENT ON COLUMN "public"."translate_table"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."translate_table"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."translate_table"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."translate_table"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."translate_table"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "translate_table_field";
DROP SEQUENCE IF EXISTS translate_table_field_transl_tbl_field_id_seq;
CREATE SEQUENCE translate_table_field_transl_tbl_field_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."translate_table_field" (
    "transl_tbl_field_id" integer DEFAULT nextval('translate_table_field_transl_tbl_field_id_seq') NOT NULL,
    "field_org_desc" character varying(200) NOT NULL,
    "field_transl_desc" character varying(200) NOT NULL,
    "field" character varying(200) NOT NULL,
    "table" character varying(200) NOT NULL,
    "db" character varying(200) NOT NULL,
    "lang" character varying(5) NOT NULL,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "translate_table_field_pkey" PRIMARY KEY ("transl_tbl_field_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."translate_table_field" IS 'Translate Table Fields';

COMMENT ON COLUMN "public"."translate_table_field"."transl_tbl_field_id" IS 'Translate Table Field ID';

COMMENT ON COLUMN "public"."translate_table_field"."field_org_desc" IS 'Field Org. Desc';

COMMENT ON COLUMN "public"."translate_table_field"."field_transl_desc" IS 'Field Transl. Desc';

COMMENT ON COLUMN "public"."translate_table_field"."field" IS 'Field';

COMMENT ON COLUMN "public"."translate_table_field"."table" IS 'Table';

COMMENT ON COLUMN "public"."translate_table_field"."db" IS 'Database';

COMMENT ON COLUMN "public"."translate_table_field"."lang" IS 'Lang';

COMMENT ON COLUMN "public"."translate_table_field"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."translate_table_field"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."translate_table_field"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."translate_table_field"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."translate_table_field"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "user_log";
DROP SEQUENCE IF EXISTS user_log_user_log_id_seq;
CREATE SEQUENCE user_log_user_log_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."user_log" (
    "user_log_id" integer DEFAULT nextval('user_log_user_log_id_seq') NOT NULL,
    "user_id" integer,
    "action" character varying(200) NOT NULL,
    "req_ip" character varying(200),
    "req_at" timestamp,
    "req_data" text,
    "res_at" timestamp,
    "res_type" character varying(200),
    "res_msg" character varying(500),
    "res_data" text,
    "table" character varying(200),
    "db" character varying(200),
    "row_id" integer,
    "app_id" integer,
    "new_data" text,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "user_log_pkey" PRIMARY KEY ("user_log_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."user_log" IS 'User Logs';

COMMENT ON COLUMN "public"."user_log"."user_log_id" IS 'User Log ID';

COMMENT ON COLUMN "public"."user_log"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."user_log"."action" IS 'Action';

COMMENT ON COLUMN "public"."user_log"."req_ip" IS 'Request IP';

COMMENT ON COLUMN "public"."user_log"."req_at" IS 'Request at';

COMMENT ON COLUMN "public"."user_log"."req_data" IS 'Request Data';

COMMENT ON COLUMN "public"."user_log"."res_at" IS 'Response at';

COMMENT ON COLUMN "public"."user_log"."res_type" IS 'Response Type';

COMMENT ON COLUMN "public"."user_log"."res_msg" IS 'Response Message';

COMMENT ON COLUMN "public"."user_log"."res_data" IS 'Request Data';

COMMENT ON COLUMN "public"."user_log"."table" IS 'Table';

COMMENT ON COLUMN "public"."user_log"."db" IS 'Database';

COMMENT ON COLUMN "public"."user_log"."row_id" IS 'Database';

COMMENT ON COLUMN "public"."user_log"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."user_log"."new_data" IS 'New Data';

COMMENT ON COLUMN "public"."user_log"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."user_log"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."user_log"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "user_role";
DROP SEQUENCE IF EXISTS user_role_user_role_id_seq;
CREATE SEQUENCE user_role_user_role_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."user_role" (
    "user_role_id" integer DEFAULT nextval('user_role_user_role_id_seq') NOT NULL,
    "user_id" integer,
    "role_id" integer,
    "active" boolean,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "user_role_pkey" PRIMARY KEY ("user_role_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."user_role" IS 'User Roles';

COMMENT ON COLUMN "public"."user_role"."user_role_id" IS 'User Role ID';

COMMENT ON COLUMN "public"."user_role"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."user_role"."role_id" IS 'Role ID';

COMMENT ON COLUMN "public"."user_role"."active" IS 'Active';

COMMENT ON COLUMN "public"."user_role"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."user_role"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."user_role"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "users";
DROP SEQUENCE IF EXISTS users_user_id_seq;
CREATE SEQUENCE users_user_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."users" (
    "user_id" integer DEFAULT nextval('users_user_id_seq') NOT NULL,
    "username" character varying(50) NOT NULL,
    "first_name" character varying(50) NOT NULL,
    "last_name" character varying(50) NOT NULL,
    "email" character varying(50),
    "phone" character varying(50),
    "password" character varying(200) NOT NULL,
    "role_id" integer,
    "lang_id" integer,
    "timezone" character varying(50),
    "attach_profile_pic" character varying(200),
    "active" boolean,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "users_pkey" PRIMARY KEY ("user_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."users" IS 'Users';

COMMENT ON COLUMN "public"."users"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."users"."username" IS 'Username';

COMMENT ON COLUMN "public"."users"."first_name" IS 'Fisrt Name';

COMMENT ON COLUMN "public"."users"."last_name" IS 'Last Name';

COMMENT ON COLUMN "public"."users"."email" IS 'Email';

COMMENT ON COLUMN "public"."users"."phone" IS 'Phone';

COMMENT ON COLUMN "public"."users"."password" IS 'Password';

COMMENT ON COLUMN "public"."users"."role_id" IS 'Default Role ID';

COMMENT ON COLUMN "public"."users"."lang_id" IS 'Lang ID';

COMMENT ON COLUMN "public"."users"."timezone" IS 'Timezone';

COMMENT ON COLUMN "public"."users"."attach_profile_pic" IS 'Profile Picture';

COMMENT ON COLUMN "public"."users"."active" IS 'Active';

COMMENT ON COLUMN "public"."users"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."users"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."users"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX users_username_key ON public.users USING btree (username);

CREATE UNIQUE INDEX users_email_key ON public.users USING btree (email);

CREATE UNIQUE INDEX users_phone_key ON public.users USING btree (phone);


ALTER TABLE ONLY "public"."access_key" ADD CONSTRAINT "access_key_for_user_id_fkey" FOREIGN KEY (for_user_id) REFERENCES users(user_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."access_key" ADD CONSTRAINT "access_key_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."app" ADD CONSTRAINT "app_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."calendar" ADD CONSTRAINT "calendar_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."calendar" ADD CONSTRAINT "calendar_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."column_level_access" ADD CONSTRAINT "column_level_access_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."column_level_access" ADD CONSTRAINT "column_level_access_table_id_fkey" FOREIGN KEY (table_id) REFERENCES "table"(table_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."column_level_access" ADD CONSTRAINT "column_level_access_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."cron" ADD CONSTRAINT "cron_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."cron_log" ADD CONSTRAINT "cron_log_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."cron_log" ADD CONSTRAINT "cron_log_cron_id_fkey" FOREIGN KEY (cron_id) REFERENCES cron(cron_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."custom_form" ADD CONSTRAINT "custom_form_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."custom_form" ADD CONSTRAINT "custom_form_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."custom_table" ADD CONSTRAINT "custom_table_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."custom_table" ADD CONSTRAINT "custom_table_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."dashboard" ADD CONSTRAINT "dashboard_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."dashboard" ADD CONSTRAINT "dashboard_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."dashboard_comment" ADD CONSTRAINT "dashboard_comment_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."dashboard_comment" ADD CONSTRAINT "dashboard_comment_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_rb_exp_dtail" ADD CONSTRAINT "etl_rb_exp_dtail_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rb_exp_dtail" ADD CONSTRAINT "etl_rb_exp_dtail_etl_rbase_export_id_fkey" FOREIGN KEY (etl_rbase_export_id) REFERENCES etl_rbase_export(etl_rbase_export_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rb_exp_dtail" ADD CONSTRAINT "etl_rb_exp_dtail_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rb_exp_dtail" ADD CONSTRAINT "etl_rb_exp_dtail_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_rb_output_field" ADD CONSTRAINT "etl_rb_output_field_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rb_output_field" ADD CONSTRAINT "etl_rb_output_field_etl_rbase_output_id_fkey" FOREIGN KEY (etl_rbase_output_id) REFERENCES etl_rbase_output(etl_rbase_output_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rb_output_field" ADD CONSTRAINT "etl_rb_output_field_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rb_output_field" ADD CONSTRAINT "etl_rb_output_field_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_rb_reconc_dtail" ADD CONSTRAINT "etl_rb_reconc_dtail_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rb_reconc_dtail" ADD CONSTRAINT "etl_rb_reconc_dtail_etl_rb_reconcilia_id_fkey" FOREIGN KEY (etl_rb_reconcilia_id) REFERENCES etl_rb_reconcilia(etl_rb_reconcilia_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rb_reconc_dtail" ADD CONSTRAINT "etl_rb_reconc_dtail_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rb_reconc_dtail" ADD CONSTRAINT "etl_rb_reconc_dtail_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_rb_reconcilia" ADD CONSTRAINT "etl_rb_reconcilia_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rb_reconcilia" ADD CONSTRAINT "etl_rb_reconcilia_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rb_reconcilia" ADD CONSTRAINT "etl_rb_reconcilia_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_rbase_backup" ADD CONSTRAINT "etl_rbase_backup_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_backup" ADD CONSTRAINT "etl_rbase_backup_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_backup" ADD CONSTRAINT "etl_rbase_backup_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_rbase_export" ADD CONSTRAINT "etl_rbase_export_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_export" ADD CONSTRAINT "etl_rbase_export_etl_rbase_output_id_fkey" FOREIGN KEY (etl_rbase_output_id) REFERENCES etl_rbase_output(etl_rbase_output_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_export" ADD CONSTRAINT "etl_rbase_export_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_export" ADD CONSTRAINT "etl_rbase_export_export_type_id_fkey" FOREIGN KEY (export_type_id) REFERENCES export_type(export_type_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_export" ADD CONSTRAINT "etl_rbase_export_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_rbase_input" ADD CONSTRAINT "etl_rbase_input_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_input" ADD CONSTRAINT "etl_rbase_input_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_input" ADD CONSTRAINT "etl_rbase_input_input_type_id_fkey" FOREIGN KEY (input_type_id) REFERENCES input_type(input_type_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_input" ADD CONSTRAINT "etl_rbase_input_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_rbase_notify" ADD CONSTRAINT "etl_rbase_notify_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_notify" ADD CONSTRAINT "etl_rbase_notify_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_notify" ADD CONSTRAINT "etl_rbase_notify_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_rbase_output" ADD CONSTRAINT "etl_rbase_output_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_output" ADD CONSTRAINT "etl_rbase_output_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_output" ADD CONSTRAINT "etl_rbase_output_output_type_id_fkey" FOREIGN KEY (output_type_id) REFERENCES output_type(output_type_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_output" ADD CONSTRAINT "etl_rbase_output_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_rbase_quality" ADD CONSTRAINT "etl_rbase_quality_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_quality" ADD CONSTRAINT "etl_rbase_quality_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_quality" ADD CONSTRAINT "etl_rbase_quality_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_rbase_script" ADD CONSTRAINT "etl_rbase_script_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_script" ADD CONSTRAINT "etl_rbase_script_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_rbase_script" ADD CONSTRAINT "etl_rbase_script_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_report_base" ADD CONSTRAINT "etl_report_base_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_report_base" ADD CONSTRAINT "etl_report_base_periodicity_id_fkey" FOREIGN KEY (periodicity_id) REFERENCES periodicity(periodicity_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_report_base" ADD CONSTRAINT "etl_report_base_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."etl_report_base_log" ADD CONSTRAINT "etl_report_base_log_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_report_base_log" ADD CONSTRAINT "etl_report_base_log_etl_report_base_id_fkey" FOREIGN KEY (etl_report_base_id) REFERENCES etl_report_base(etl_report_base_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."etl_report_base_log" ADD CONSTRAINT "etl_report_base_log_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."manage_query" ADD CONSTRAINT "manage_query_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."manage_query" ADD CONSTRAINT "manage_query_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."menu" ADD CONSTRAINT "menu_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."menu" ADD CONSTRAINT "menu_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."menu_table" ADD CONSTRAINT "menu_table_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."menu_table" ADD CONSTRAINT "menu_table_menu_id_fkey" FOREIGN KEY (menu_id) REFERENCES menu(menu_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."menu_table" ADD CONSTRAINT "menu_table_table_id_fkey" FOREIGN KEY (table_id) REFERENCES "table"(table_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."menu_table" ADD CONSTRAINT "menu_table_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."role_app" ADD CONSTRAINT "role_app_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_app" ADD CONSTRAINT "role_app_role_id_fkey" FOREIGN KEY (role_id) REFERENCES role(role_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_app" ADD CONSTRAINT "role_app_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."role_app_menu" ADD CONSTRAINT "role_app_menu_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_app_menu" ADD CONSTRAINT "role_app_menu_menu_id_fkey" FOREIGN KEY (menu_id) REFERENCES menu(menu_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_app_menu" ADD CONSTRAINT "role_app_menu_role_id_fkey" FOREIGN KEY (role_id) REFERENCES role(role_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_app_menu" ADD CONSTRAINT "role_app_menu_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."role_app_menu_table" ADD CONSTRAINT "role_app_menu_table_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_app_menu_table" ADD CONSTRAINT "role_app_menu_table_menu_id_fkey" FOREIGN KEY (menu_id) REFERENCES menu(menu_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_app_menu_table" ADD CONSTRAINT "role_app_menu_table_role_id_fkey" FOREIGN KEY (role_id) REFERENCES role(role_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_app_menu_table" ADD CONSTRAINT "role_app_menu_table_table_id_fkey" FOREIGN KEY (table_id) REFERENCES "table"(table_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_app_menu_table" ADD CONSTRAINT "role_app_menu_table_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."role_row_level_access" ADD CONSTRAINT "role_row_level_access_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_row_level_access" ADD CONSTRAINT "role_row_level_access_role_id_fkey" FOREIGN KEY (role_id) REFERENCES role(role_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_row_level_access" ADD CONSTRAINT "role_row_level_access_table_id_fkey" FOREIGN KEY (table_id) REFERENCES "table"(table_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."role_row_level_access" ADD CONSTRAINT "role_row_level_access_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."row_level_access" ADD CONSTRAINT "row_level_access_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."row_level_access" ADD CONSTRAINT "row_level_access_table_id_fkey" FOREIGN KEY (table_id) REFERENCES "table"(table_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."row_level_access" ADD CONSTRAINT "row_level_access_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."table" ADD CONSTRAINT "table_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."table_schema" ADD CONSTRAINT "table_schema_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."task" ADD CONSTRAINT "task_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."task" ADD CONSTRAINT "task_calendar_id_fkey" FOREIGN KEY (calendar_id) REFERENCES calendar(calendar_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."task" ADD CONSTRAINT "task_repeat_type_id_fkey" FOREIGN KEY (repeat_type_id) REFERENCES repeat_type(repeat_type_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."task" ADD CONSTRAINT "task_task_status_id_fkey" FOREIGN KEY (task_status_id) REFERENCES task_status(task_status_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."task" ADD CONSTRAINT "task_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."task_track" ADD CONSTRAINT "task_track_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."task_track" ADD CONSTRAINT "task_track_task_id_fkey" FOREIGN KEY (task_id) REFERENCES task(task_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."task_track" ADD CONSTRAINT "task_track_task_status_id_fkey" FOREIGN KEY (task_status_id) REFERENCES task_status(task_status_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."task_track" ADD CONSTRAINT "task_track_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."translate_table" ADD CONSTRAINT "translate_table_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."translate_table" ADD CONSTRAINT "translate_table_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."translate_table_field" ADD CONSTRAINT "translate_table_field_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."translate_table_field" ADD CONSTRAINT "translate_table_field_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."user_log" ADD CONSTRAINT "user_log_app_id_fkey" FOREIGN KEY (app_id) REFERENCES app(app_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."user_log" ADD CONSTRAINT "user_log_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."user_role" ADD CONSTRAINT "user_role_role_id_fkey" FOREIGN KEY (role_id) REFERENCES role(role_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."user_role" ADD CONSTRAINT "user_role_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(user_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."users" ADD CONSTRAINT "users_lang_id_fkey" FOREIGN KEY (lang_id) REFERENCES lang(lang_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."users" ADD CONSTRAINT "users_role_id_fkey" FOREIGN KEY (role_id) REFERENCES role(role_id) NOT DEFERRABLE;

-- 2026-01-13 13:14:09 UTC
