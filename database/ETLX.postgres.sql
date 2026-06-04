-- Adminer 5.3.0 PostgreSQL 18.0 dump

\connect "ETLX";

DROP TABLE IF EXISTS "flight_schema";
DROP SEQUENCE IF EXISTS flight_schema_flight_schema_id_seq;
CREATE SEQUENCE flight_schema_flight_schema_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."flight_schema" (
    "flight_schema_id" integer DEFAULT nextval('flight_schema_flight_schema_id_seq') NOT NULL,
    "flight_schema" character varying(200) NOT NULL,
    "flight_schema_desc" text,
    "flight_schema_conf" text,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "flight_schema_pkey" PRIMARY KEY ("flight_schema_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."flight_schema" IS 'Expose Arrow Flight';

COMMENT ON COLUMN "public"."flight_schema"."flight_schema_id" IS 'ID';

COMMENT ON COLUMN "public"."flight_schema"."flight_schema" IS 'Name';

COMMENT ON COLUMN "public"."flight_schema"."flight_schema_desc" IS 'Description';

COMMENT ON COLUMN "public"."flight_schema"."flight_schema_conf" IS 'Config Text';

COMMENT ON COLUMN "public"."flight_schema"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."flight_schema"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."flight_schema"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."flight_schema"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."flight_schema"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX flight_schema_flight_schema_key ON public.flight_schema USING btree (flight_schema);


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


DROP TABLE IF EXISTS "etlx";
DROP SEQUENCE IF EXISTS etlx_etlx_id_seq;
CREATE SEQUENCE etlx_etlx_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etlx" (
    "etlx_id" integer DEFAULT nextval('etlx_etlx_id_seq') NOT NULL,
    "etl" character varying(200) NOT NULL,
    "etl_desc" text,
    "attach_etlx_conf" character varying(200),
    "etlx_conf" text,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etlx_pkey" PRIMARY KEY ("etlx_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etlx" IS 'ETLX';

COMMENT ON COLUMN "public"."etlx"."etlx_id" IS 'ID';

COMMENT ON COLUMN "public"."etlx"."etl" IS 'Name';

COMMENT ON COLUMN "public"."etlx"."etl_desc" IS 'Description';

COMMENT ON COLUMN "public"."etlx"."attach_etlx_conf" IS 'Config File';

COMMENT ON COLUMN "public"."etlx"."etlx_conf" IS 'Config Text';

COMMENT ON COLUMN "public"."etlx"."active" IS 'Active';

COMMENT ON COLUMN "public"."etlx"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etlx"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etlx"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etlx"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etlx"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX etlx_etl_key ON public.etlx USING btree (etl);


DROP TABLE IF EXISTS "etlx_conf";
DROP SEQUENCE IF EXISTS etlx_conf_etlx_conf_id_seq;
CREATE SEQUENCE etlx_conf_etlx_conf_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."etlx_conf" (
    "etlx_conf_id" integer DEFAULT nextval('etlx_conf_etlx_conf_id_seq') NOT NULL,
    "etlx_conf" character varying(200) NOT NULL,
    "etlx_conf_desc" text,
    "etlx_extra_conf" text,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "etlx_conf_pkey" PRIMARY KEY ("etlx_conf_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."etlx_conf" IS 'ETLX Extra Cofig';

COMMENT ON COLUMN "public"."etlx_conf"."etlx_conf_id" IS 'ID';

COMMENT ON COLUMN "public"."etlx_conf"."etlx_conf" IS 'Name';

COMMENT ON COLUMN "public"."etlx_conf"."etlx_conf_desc" IS 'Description';

COMMENT ON COLUMN "public"."etlx_conf"."etlx_extra_conf" IS 'Config Text';

COMMENT ON COLUMN "public"."etlx_conf"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."etlx_conf"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."etlx_conf"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."etlx_conf"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."etlx_conf"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX etlx_conf_etlx_conf_key ON public.etlx_conf USING btree (etlx_conf);


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


DROP TABLE IF EXISTS "notebook";
DROP SEQUENCE IF EXISTS notebook_notebook_id_seq;
CREATE SEQUENCE notebook_notebook_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."notebook" (
    "notebook_id" integer DEFAULT nextval('notebook_notebook_id_seq') NOT NULL,
    "notebook" character varying(200),
    "notebook_desc" text,
    "notebook_conf" text NOT NULL,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "notebook_pkey" PRIMARY KEY ("notebook_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."notebook" IS 'Notebooks';

COMMENT ON COLUMN "public"."notebook"."notebook_id" IS 'Notebook ID';

COMMENT ON COLUMN "public"."notebook"."notebook" IS 'Name';

COMMENT ON COLUMN "public"."notebook"."notebook_desc" IS 'Description';

COMMENT ON COLUMN "public"."notebook"."notebook_conf" IS 'Conf / Params';

COMMENT ON COLUMN "public"."notebook"."active" IS 'Active';

COMMENT ON COLUMN "public"."notebook"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."notebook"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."notebook"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."notebook"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."notebook"."excluded" IS 'Excluded';


-- 2026-01-13 13:14:34 UTC
