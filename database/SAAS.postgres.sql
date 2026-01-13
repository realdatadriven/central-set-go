-- Adminer 5.3.0 PostgreSQL 18.0 dump

\connect "SAAS";

DROP TABLE IF EXISTS "currency";
DROP SEQUENCE IF EXISTS currency_currency_id_seq;
CREATE SEQUENCE currency_currency_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."currency" (
    "currency_id" integer DEFAULT nextval('currency_currency_id_seq') NOT NULL,
    "currency" character varying(3) NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "currency_pkey" PRIMARY KEY ("currency_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."currency" IS 'Currency';

COMMENT ON COLUMN "public"."currency"."currency_id" IS 'Currency ID';

COMMENT ON COLUMN "public"."currency"."currency" IS 'Currency';

COMMENT ON COLUMN "public"."currency"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."currency"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."currency"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX currency_currency_key ON public.currency USING btree (currency);


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


DROP TABLE IF EXISTS "deployment";
DROP SEQUENCE IF EXISTS deployment_deployment_id_seq;
CREATE SEQUENCE deployment_deployment_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."deployment" (
    "deployment_id" integer DEFAULT nextval('deployment_deployment_id_seq') NOT NULL,
    "deployment" character varying(200),
    "deployment_desc" text,
    "product_id" integer NOT NULL,
    "provider_id" integer NOT NULL,
    "terraform_template" text NOT NULL,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "deployment_pkey" PRIMARY KEY ("deployment_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."deployment" IS 'Deployments';

COMMENT ON COLUMN "public"."deployment"."deployment_id" IS 'Deployment ID';

COMMENT ON COLUMN "public"."deployment"."deployment" IS 'Deployment';

COMMENT ON COLUMN "public"."deployment"."deployment_desc" IS 'Deployment Description';

COMMENT ON COLUMN "public"."deployment"."product_id" IS 'Product ID';

COMMENT ON COLUMN "public"."deployment"."provider_id" IS 'Provider ID';

COMMENT ON COLUMN "public"."deployment"."terraform_template" IS 'Terraform Script';

COMMENT ON COLUMN "public"."deployment"."active" IS 'Active';

COMMENT ON COLUMN "public"."deployment"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."deployment"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."deployment"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."deployment"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."deployment"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "env";
DROP SEQUENCE IF EXISTS env_env_id_seq;
CREATE SEQUENCE env_env_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."env" (
    "env_id" integer DEFAULT nextval('env_env_id_seq') NOT NULL,
    "env_name" character varying(200) NOT NULL,
    "env_value" text NOT NULL,
    "tenant_id" integer NOT NULL,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "env_pkey" PRIMARY KEY ("env_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."env" IS 'Envariomental Variables';

COMMENT ON COLUMN "public"."env"."env_id" IS 'env ID';

COMMENT ON COLUMN "public"."env"."env_name" IS 'Env Name';

COMMENT ON COLUMN "public"."env"."env_value" IS 'Env Value';

COMMENT ON COLUMN "public"."env"."tenant_id" IS 'Tenant ID';

COMMENT ON COLUMN "public"."env"."active" IS 'Active';

COMMENT ON COLUMN "public"."env"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."env"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."env"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."env"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."env"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX env_env_name_key ON public.env USING btree (env_name);


DROP TABLE IF EXISTS "interval";
DROP SEQUENCE IF EXISTS interval_interval_id_seq;
CREATE SEQUENCE interval_interval_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."interval" (
    "interval_id" integer DEFAULT nextval('interval_interval_id_seq') NOT NULL,
    "interval" character varying(100) NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "interval_pkey" PRIMARY KEY ("interval_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."interval" IS 'Intervals';

COMMENT ON COLUMN "public"."interval"."interval_id" IS 'Interval ID';

COMMENT ON COLUMN "public"."interval"."interval" IS 'Interval';

COMMENT ON COLUMN "public"."interval"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."interval"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."interval"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX interval_interval_key ON public."interval" USING btree ("interval");


DROP TABLE IF EXISTS "payment_plan";
DROP SEQUENCE IF EXISTS payment_plan_payment_plan_id_seq;
CREATE SEQUENCE payment_plan_payment_plan_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."payment_plan" (
    "payment_plan_id" integer DEFAULT nextval('payment_plan_payment_plan_id_seq') NOT NULL,
    "plan_id" integer NOT NULL,
    "deployment_id" integer NOT NULL,
    "product_id" integer NOT NULL,
    "price" real NOT NULL,
    "currency_id" integer NOT NULL,
    "interval_id" integer NOT NULL,
    "stripe_price_id" character varying(255),
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "payment_plan_pkey" PRIMARY KEY ("payment_plan_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."payment_plan" IS 'Payment Plans';

COMMENT ON COLUMN "public"."payment_plan"."payment_plan_id" IS 'Payment Plan ID';

COMMENT ON COLUMN "public"."payment_plan"."plan_id" IS 'Plan ID';

COMMENT ON COLUMN "public"."payment_plan"."deployment_id" IS 'Deployment ID';

COMMENT ON COLUMN "public"."payment_plan"."product_id" IS 'Product ID';

COMMENT ON COLUMN "public"."payment_plan"."price" IS 'Plan Price';

COMMENT ON COLUMN "public"."payment_plan"."currency_id" IS 'Currency ID';

COMMENT ON COLUMN "public"."payment_plan"."interval_id" IS 'Billing Interval';

COMMENT ON COLUMN "public"."payment_plan"."stripe_price_id" IS 'Stripe Price ID';

COMMENT ON COLUMN "public"."payment_plan"."active" IS 'Active';

COMMENT ON COLUMN "public"."payment_plan"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."payment_plan"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."payment_plan"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."payment_plan"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."payment_plan"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "plan";
DROP SEQUENCE IF EXISTS plan_plan_id_seq;
CREATE SEQUENCE plan_plan_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."plan" (
    "plan_id" integer DEFAULT nextval('plan_plan_id_seq') NOT NULL,
    "plan" character varying(200) NOT NULL,
    "price" real,
    "deployment_id" integer NOT NULL,
    "product_id" integer NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "plan_pkey" PRIMARY KEY ("plan_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."plan" IS 'Plan';

COMMENT ON COLUMN "public"."plan"."plan_id" IS 'Plan ID';

COMMENT ON COLUMN "public"."plan"."plan" IS 'Plan';

COMMENT ON COLUMN "public"."plan"."price" IS 'Plan Price';

COMMENT ON COLUMN "public"."plan"."deployment_id" IS 'Deployment ID';

COMMENT ON COLUMN "public"."plan"."product_id" IS 'Product ID';

COMMENT ON COLUMN "public"."plan"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."plan"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."plan"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX plan_plan_key ON public.plan USING btree (plan);


DROP TABLE IF EXISTS "product";
DROP SEQUENCE IF EXISTS product_product_id_seq;
CREATE SEQUENCE product_product_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."product" (
    "product_id" integer DEFAULT nextval('product_product_id_seq') NOT NULL,
    "product" character varying(255) NOT NULL,
    "description" text,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "product_pkey" PRIMARY KEY ("product_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."product" IS 'Products';

COMMENT ON COLUMN "public"."product"."product_id" IS 'Product ID';

COMMENT ON COLUMN "public"."product"."product" IS 'Product Name';

COMMENT ON COLUMN "public"."product"."description" IS 'Product Description';

COMMENT ON COLUMN "public"."product"."active" IS 'Active';

COMMENT ON COLUMN "public"."product"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."product"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."product"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."product"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."product"."excluded" IS 'Excluded';


DROP TABLE IF EXISTS "provider";
DROP SEQUENCE IF EXISTS provider_provider_id_seq;
CREATE SEQUENCE provider_provider_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."provider" (
    "provider_id" integer DEFAULT nextval('provider_provider_id_seq') NOT NULL,
    "provider" character varying(200) NOT NULL,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "provider_pkey" PRIMARY KEY ("provider_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."provider" IS 'Cloud Provider';

COMMENT ON COLUMN "public"."provider"."provider_id" IS 'Provider ID';

COMMENT ON COLUMN "public"."provider"."provider" IS 'Provider';

COMMENT ON COLUMN "public"."provider"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."provider"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."provider"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX provider_provider_key ON public.provider USING btree (provider);


DROP TABLE IF EXISTS "subscription";
DROP SEQUENCE IF EXISTS subscription_subscription_id_seq;
CREATE SEQUENCE subscription_subscription_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."subscription" (
    "subscription_id" integer DEFAULT nextval('subscription_subscription_id_seq') NOT NULL,
    "tenant_id" integer NOT NULL,
    "plan_id" integer NOT NULL,
    "deployment_id" integer NOT NULL,
    "payment_plan_id" integer NOT NULL,
    "product_id" integer NOT NULL,
    "terraform_outputs" text,
    "tf_public_ip" character varying(200),
    "tf_public_dns" character varying(255),
    "terraform_state" text,
    "terraform_lock" text,
    "tf_err_msg" text,
    "deployed" boolean,
    "active" boolean,
    "stripe_subscription_id" character varying(255),
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "subscription_pkey" PRIMARY KEY ("subscription_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."subscription" IS 'Subscriptions';

COMMENT ON COLUMN "public"."subscription"."subscription_id" IS 'Subscription ID';

COMMENT ON COLUMN "public"."subscription"."tenant_id" IS 'Tenant ID';

COMMENT ON COLUMN "public"."subscription"."plan_id" IS 'Plan ID';

COMMENT ON COLUMN "public"."subscription"."deployment_id" IS 'Deployment ID';

COMMENT ON COLUMN "public"."subscription"."payment_plan_id" IS 'Payment Plan ID';

COMMENT ON COLUMN "public"."subscription"."product_id" IS 'Product ID';

COMMENT ON COLUMN "public"."subscription"."terraform_outputs" IS 'Terraform Outputs';

COMMENT ON COLUMN "public"."subscription"."tf_public_ip" IS 'Terraform Public IP';

COMMENT ON COLUMN "public"."subscription"."tf_public_dns" IS 'Terraform Public DNS';

COMMENT ON COLUMN "public"."subscription"."terraform_state" IS 'Terraform State';

COMMENT ON COLUMN "public"."subscription"."terraform_lock" IS 'Terraform Lock';

COMMENT ON COLUMN "public"."subscription"."tf_err_msg" IS 'Terraform error msg';

COMMENT ON COLUMN "public"."subscription"."deployed" IS 'Deployed';

COMMENT ON COLUMN "public"."subscription"."active" IS 'Active';

COMMENT ON COLUMN "public"."subscription"."stripe_subscription_id" IS 'Stripe Subscription ID';

COMMENT ON COLUMN "public"."subscription"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."subscription"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."subscription"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."subscription"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."subscription"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX subscription_stripe_subscription_id_key ON public.subscription USING btree (stripe_subscription_id);


DROP TABLE IF EXISTS "tenant";
DROP SEQUENCE IF EXISTS tenant_tenant_id_seq;
CREATE SEQUENCE tenant_tenant_id_seq INCREMENT 1 MINVALUE 1 MAXVALUE 2147483647 CACHE 1;

CREATE TABLE "public"."tenant" (
    "tenant_id" integer DEFAULT nextval('tenant_tenant_id_seq') NOT NULL,
    "tenant" character varying(200) NOT NULL,
    "email" character varying(200) NOT NULL,
    "password" character varying(200),
    "phone" character varying(200),
    "address" character varying(200),
    "currency_id" integer,
    "active" boolean,
    "user_id" integer,
    "app_id" integer,
    "created_at" timestamp,
    "updated_at" timestamp,
    "excluded" boolean,
    CONSTRAINT "tenant_pkey" PRIMARY KEY ("tenant_id")
)
WITH (oids = false);

COMMENT ON TABLE "public"."tenant" IS 'Tenants';

COMMENT ON COLUMN "public"."tenant"."tenant_id" IS 'ID';

COMMENT ON COLUMN "public"."tenant"."tenant" IS 'Tenant Name';

COMMENT ON COLUMN "public"."tenant"."email" IS 'Email';

COMMENT ON COLUMN "public"."tenant"."password" IS 'Password';

COMMENT ON COLUMN "public"."tenant"."phone" IS 'Phone';

COMMENT ON COLUMN "public"."tenant"."address" IS 'Adress';

COMMENT ON COLUMN "public"."tenant"."currency_id" IS 'Currency ID';

COMMENT ON COLUMN "public"."tenant"."active" IS 'Active';

COMMENT ON COLUMN "public"."tenant"."user_id" IS 'User ID';

COMMENT ON COLUMN "public"."tenant"."app_id" IS 'App ID';

COMMENT ON COLUMN "public"."tenant"."created_at" IS 'Created at';

COMMENT ON COLUMN "public"."tenant"."updated_at" IS 'Updated at';

COMMENT ON COLUMN "public"."tenant"."excluded" IS 'Excluded';

CREATE UNIQUE INDEX tenant_email_key ON public.tenant USING btree (email);


ALTER TABLE ONLY "public"."deployment" ADD CONSTRAINT "deployment_product_id_fkey" FOREIGN KEY (product_id) REFERENCES product(product_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."deployment" ADD CONSTRAINT "deployment_provider_id_fkey" FOREIGN KEY (provider_id) REFERENCES provider(provider_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."env" ADD CONSTRAINT "env_tenant_id_fkey" FOREIGN KEY (tenant_id) REFERENCES tenant(tenant_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."payment_plan" ADD CONSTRAINT "payment_plan_currency_id_fkey" FOREIGN KEY (currency_id) REFERENCES currency(currency_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."payment_plan" ADD CONSTRAINT "payment_plan_deployment_id_fkey" FOREIGN KEY (deployment_id) REFERENCES deployment(deployment_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."payment_plan" ADD CONSTRAINT "payment_plan_interval_id_fkey" FOREIGN KEY (interval_id) REFERENCES "interval"(interval_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."payment_plan" ADD CONSTRAINT "payment_plan_plan_id_fkey" FOREIGN KEY (plan_id) REFERENCES plan(plan_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."payment_plan" ADD CONSTRAINT "payment_plan_product_id_fkey" FOREIGN KEY (product_id) REFERENCES product(product_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."plan" ADD CONSTRAINT "plan_deployment_id_fkey" FOREIGN KEY (deployment_id) REFERENCES deployment(deployment_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."plan" ADD CONSTRAINT "plan_product_id_fkey" FOREIGN KEY (product_id) REFERENCES product(product_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."subscription" ADD CONSTRAINT "subscription_deployment_id_fkey" FOREIGN KEY (deployment_id) REFERENCES deployment(deployment_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."subscription" ADD CONSTRAINT "subscription_payment_plan_id_fkey" FOREIGN KEY (payment_plan_id) REFERENCES payment_plan(payment_plan_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."subscription" ADD CONSTRAINT "subscription_plan_id_fkey" FOREIGN KEY (plan_id) REFERENCES plan(plan_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."subscription" ADD CONSTRAINT "subscription_product_id_fkey" FOREIGN KEY (product_id) REFERENCES product(product_id) NOT DEFERRABLE;
ALTER TABLE ONLY "public"."subscription" ADD CONSTRAINT "subscription_tenant_id_fkey" FOREIGN KEY (tenant_id) REFERENCES tenant(tenant_id) NOT DEFERRABLE;

ALTER TABLE ONLY "public"."tenant" ADD CONSTRAINT "tenant_currency_id_fkey" FOREIGN KEY (currency_id) REFERENCES currency(currency_id) NOT DEFERRABLE;

-- 2026-01-13 13:15:29 UTC
