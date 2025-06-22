CREATE TABLE "account_emailaddress"(
    "id" INTEGER NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "verified" BOOLEAN NOT NULL,
    "primary" BOOLEAN NOT NULL,
    "user_id" INTEGER NOT NULL
);
ALTER TABLE
    "account_emailaddress" ADD PRIMARY KEY("id");
ALTER TABLE
    "account_emailaddress" ADD CONSTRAINT "account_emailaddress_email_unique" UNIQUE("email");
CREATE INDEX "account_emailaddress_user_id_index" ON
    "account_emailaddress"("user_id");
CREATE TABLE "account_emailconfirmation"(
    "id" INTEGER NOT NULL,
    "created" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "sent" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "key" VARCHAR(255) NOT NULL,
    "email_address_id" INTEGER NOT NULL
);
ALTER TABLE
    "account_emailconfirmation" ADD PRIMARY KEY("id");
ALTER TABLE
    "account_emailconfirmation" ADD CONSTRAINT "account_emailconfirmation_key_unique" UNIQUE("key");
CREATE INDEX "account_emailconfirmation_email_address_id_index" ON
    "account_emailconfirmation"("email_address_id");
CREATE TABLE "auth_group"(
    "id" INTEGER NOT NULL,
    "name" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "auth_group" ADD PRIMARY KEY("id");
ALTER TABLE
    "auth_group" ADD CONSTRAINT "auth_group_name_unique" UNIQUE("name");
CREATE TABLE "auth_group_permissions"(
    "id" INTEGER NOT NULL,
    "group_id" INTEGER NOT NULL,
    "permission_id" INTEGER NOT NULL
);
ALTER TABLE
    "auth_group_permissions" ADD CONSTRAINT "auth_group_permissions_group_id_permission_id_unique" UNIQUE("group_id", "permission_id");
ALTER TABLE
    "auth_group_permissions" ADD PRIMARY KEY("id");
CREATE INDEX "auth_group_permissions_group_id_index" ON
    "auth_group_permissions"("group_id");
CREATE INDEX "auth_group_permissions_permission_id_index" ON
    "auth_group_permissions"("permission_id");
CREATE TABLE "auth_permission"(
    "id" INTEGER NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "content_type_id" INTEGER NOT NULL,
    "codename" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "auth_permission" ADD CONSTRAINT "auth_permission_content_type_id_codename_unique" UNIQUE("content_type_id", "codename");
ALTER TABLE
    "auth_permission" ADD PRIMARY KEY("id");
CREATE INDEX "auth_permission_content_type_id_index" ON
    "auth_permission"("content_type_id");
CREATE TABLE "django_admin_log"(
    "id" INTEGER NOT NULL,
    "action_time" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "object_id" TEXT NULL,
    "object_repr" VARCHAR(255) NOT NULL,
    "action_flag" SMALLINT NOT NULL,
    "change_message" TEXT NOT NULL,
    "content_type_id" INTEGER NULL,
    "user_id" INTEGER NOT NULL
);
ALTER TABLE
    "django_admin_log" ADD PRIMARY KEY("id");
CREATE INDEX "django_admin_log_content_type_id_index" ON
    "django_admin_log"("content_type_id");
CREATE INDEX "django_admin_log_user_id_index" ON
    "django_admin_log"("user_id");
CREATE TABLE "django_content_type"(
    "id" INTEGER NOT NULL,
    "app_label" VARCHAR(255) NOT NULL,
    "model" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "django_content_type" ADD CONSTRAINT "django_content_type_app_label_model_unique" UNIQUE("app_label", "model");
ALTER TABLE
    "django_content_type" ADD PRIMARY KEY("id");
CREATE TABLE "django_migrations"(
    "id" INTEGER NOT NULL,
    "app" VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "applied" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "django_migrations" ADD PRIMARY KEY("id");
CREATE TABLE "django_session"(
    "session_key" VARCHAR(255) NOT NULL,
    "session_data" TEXT NOT NULL,
    "expire_date" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "django_session" ADD PRIMARY KEY("session_key");
CREATE INDEX "django_session_expire_date_index" ON
    "django_session"("expire_date");
CREATE TABLE "django_site"(
    "id" INTEGER NOT NULL,
    "domain" VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "django_site" ADD PRIMARY KEY("id");
ALTER TABLE
    "django_site" ADD CONSTRAINT "django_site_domain_unique" UNIQUE("domain");
CREATE TABLE "djstripe_customer"(
    "djstripe_id" BIGINT NOT NULL,
    "id" VARCHAR(255) NOT NULL,
    "livemode" BOOLEAN NOT NULL,
    "created" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "metadata" TEXT NULL,
    "description" TEXT NULL,
    "djstripe_created" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "djstripe_updated" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "balance" BIGINT NOT NULL,
    "business_vat_id" VARCHAR(255) NOT NULL,
    "currency" VARCHAR(255) NOT NULL,
    "delinquent" BOOLEAN NOT NULL,
    "coupon_start" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "coupon_end" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "email" TEXT NOT NULL,
    "shipping" TEXT NULL,
    "date_purged" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "coupon_id" BIGINT NULL,
    "default_source_id" VARCHAR(255) NULL,
    "subscriber_id" INTEGER NULL,
    "address" TEXT NULL,
    "invoice_prefix" VARCHAR(255) NOT NULL,
    "invoice_settings" TEXT NULL,
    "name" TEXT NOT NULL,
    "phone" TEXT NOT NULL,
    "preferred_locales" TEXT NULL,
    "tax_exempt" VARCHAR(255) NOT NULL,
    "default_payment_method_id" BIGINT NULL
);
ALTER TABLE
    "djstripe_customer" ADD CONSTRAINT "djstripe_customer_subscriber_id_livemode_unique" UNIQUE("subscriber_id", "livemode");
ALTER TABLE
    "djstripe_customer" ADD PRIMARY KEY("djstripe_id");
ALTER TABLE
    "djstripe_customer" ADD CONSTRAINT "djstripe_customer_id_unique" UNIQUE("id");
CREATE INDEX "djstripe_customer_coupon_id_index" ON
    "djstripe_customer"("coupon_id");
CREATE INDEX "djstripe_customer_default_source_id_index" ON
    "djstripe_customer"("default_source_id");
CREATE INDEX "djstripe_customer_subscriber_id_index" ON
    "djstripe_customer"("subscriber_id");
CREATE INDEX "djstripe_customer_default_payment_method_id_index" ON
    "djstripe_customer"("default_payment_method_id");
CREATE TABLE "djstripe_plan"(
    "djstripe_id" BIGINT NOT NULL,
    "id" VARCHAR(255) NOT NULL,
    "livemode" BOOLEAN NULL,
    "created" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "metadata" TEXT NULL,
    "description" TEXT NULL,
    "djstripe_created" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "djstripe_updated" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "aggregate_usage" VARCHAR(255) NOT NULL,
    "amount" DECIMAL(8, 2) NULL,
    "billing_scheme" VARCHAR(255) NOT NULL,
    "currency" VARCHAR(255) NOT NULL,
    "interval" VARCHAR(255) NOT NULL,
    "interval_count" INTEGER NULL,
    "nickname" TEXT NOT NULL,
    "tiers" TEXT NULL,
    "tiers_mode" VARCHAR(255) NULL,
    "transform_usage" TEXT NULL,
    "trial_period_days" INTEGER NULL,
    "usage_type" VARCHAR(255) NOT NULL,
    "name" TEXT NULL,
    "statement_descriptor" VARCHAR(255) NULL,
    "product_id" BIGINT NULL,
    "active" BOOLEAN NOT NULL
);
ALTER TABLE
    "djstripe_plan" ADD PRIMARY KEY("djstripe_id");
ALTER TABLE
    "djstripe_plan" ADD CONSTRAINT "djstripe_plan_id_unique" UNIQUE("id");
CREATE INDEX "djstripe_plan_product_id_index" ON
    "djstripe_plan"("product_id");
CREATE TABLE "djstripe_product"(
    "djstripe_id" BIGINT NOT NULL,
    "id" VARCHAR(255) NOT NULL,
    "livemode" BOOLEAN NOT NULL,
    "created" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "metadata" TEXT NULL,
    "description" TEXT NULL,
    "djstripe_created" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "djstripe_updated" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "name" TEXT NOT NULL,
    "type" VARCHAR(255) NOT NULL,
    "active" BOOLEAN NOT NULL,
    "attributes" TEXT NULL,
    "caption" TEXT NOT NULL,
    "deactivate_on" TEXT NULL,
    "images" TEXT NULL,
    "package_dimensions" TEXT NULL,
    "shippable" BOOLEAN NOT NULL,
    "url" VARCHAR(255) NULL,
    "statement_descriptor" VARCHAR(255) NOT NULL,
    "unit_label" VARCHAR(255) NOT NULL
);
ALTER TABLE
    "djstripe_product" ADD PRIMARY KEY("djstripe_id");
ALTER TABLE
    "djstripe_product" ADD CONSTRAINT "djstripe_product_id_unique" UNIQUE("id");
CREATE TABLE "djstripe_subscription"(
    "djstripe_id" BIGINT NOT NULL,
    "id" VARCHAR(255) NOT NULL,
    "livemode" BOOLEAN NULL,
    "created" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "metadata" TEXT NOT NULL,
    "description" TEXT NOT NULL,
    "djstripe_created" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "djstripe_updated" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "application_fee_percent" DECIMAL(8, 2) NULL,
    "billing" VARCHAR(255) NOT NULL,
    "billing_cycle_anchor" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "cancel_at_period_end" BOOLEAN NOT NULL,
    "canceled_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "current_period_end" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "current_period_start" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "days_until_due" INTEGER NULL,
    "ended_at" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "quantity" INTEGER NULL,
    "start" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "status" VARCHAR(255) NOT NULL,
    "tax_percent" DECIMAL(8, 2) NULL,
    "trial_end" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "trial_start" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "customer_id" BIGINT NOT NULL,
    "plan_id" BIGINT NULL,
    "pending_setup_intent_id" BIGINT NULL
);
ALTER TABLE
    "djstripe_subscription" ADD PRIMARY KEY("djstripe_id");
ALTER TABLE
    "djstripe_subscription" ADD CONSTRAINT "djstripe_subscription_id_unique" UNIQUE("id");
CREATE INDEX "djstripe_subscription_customer_id_index" ON
    "djstripe_subscription"("customer_id");
CREATE INDEX "djstripe_subscription_plan_id_index" ON
    "djstripe_subscription"("plan_id");
CREATE INDEX "djstripe_subscription_pending_setup_intent_id_index" ON
    "djstripe_subscription"("pending_setup_intent_id");
CREATE TABLE "pegasus_examples_employee"(
    "id" INTEGER NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "department" VARCHAR(255) NOT NULL,
    "salary" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL
);
ALTER TABLE
    "pegasus_examples_employee" ADD PRIMARY KEY("id");
CREATE INDEX "pegasus_examples_employee_user_id_index" ON
    "pegasus_examples_employee"("user_id");
CREATE TABLE "pegasus_examples_payment"(
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "id" UUID NOT NULL,
    "charge_id" VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "amount" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL
);
ALTER TABLE
    "pegasus_examples_payment" ADD PRIMARY KEY("created_at");
CREATE INDEX "pegasus_examples_payment_user_id_index" ON
    "pegasus_examples_payment"("user_id");
CREATE TABLE "teams_invitation"(
    "id" UUID NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "role" VARCHAR(255) NOT NULL,
    "is_accepted" BOOLEAN NOT NULL,
    "invited_by_id" INTEGER NOT NULL,
    "team_id" INTEGER NOT NULL,
    "accepted_by_id" INTEGER NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "teams_invitation" ADD CONSTRAINT "teams_invitation_team_id_email_unique" UNIQUE("team_id", "email");
ALTER TABLE
    "teams_invitation" ADD PRIMARY KEY("id");
CREATE INDEX "teams_invitation_invited_by_id_index" ON
    "teams_invitation"("invited_by_id");
CREATE INDEX "teams_invitation_team_id_index" ON
    "teams_invitation"("team_id");
CREATE INDEX "teams_invitation_accepted_by_id_index" ON
    "teams_invitation"("accepted_by_id");
CREATE TABLE "teams_membership"(
    "id" INTEGER NOT NULL,
    "role" VARCHAR(255) NOT NULL,
    "team_id" INTEGER NOT NULL,
    "user_id" INTEGER NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL
);
ALTER TABLE
    "teams_membership" ADD PRIMARY KEY("id");
CREATE INDEX "teams_membership_team_id_index" ON
    "teams_membership"("team_id");
CREATE INDEX "teams_membership_user_id_index" ON
    "teams_membership"("user_id");
CREATE TABLE "teams_team"(
    "id" INTEGER NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    "slug" VARCHAR(255) NOT NULL,
    "created_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "updated_at" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "subscription_id" BIGINT NULL
);
ALTER TABLE
    "teams_team" ADD PRIMARY KEY("id");
ALTER TABLE
    "teams_team" ADD CONSTRAINT "teams_team_slug_unique" UNIQUE("slug");
CREATE INDEX "teams_team_subscription_id_index" ON
    "teams_team"("subscription_id");
CREATE TABLE "users_customuser"(
    "id" INTEGER NOT NULL,
    "password" VARCHAR(255) NOT NULL,
    "last_login" TIMESTAMP(0) WITHOUT TIME ZONE NULL,
    "is_superuser" BOOLEAN NOT NULL,
    "username" VARCHAR(255) NOT NULL,
    "first_name" VARCHAR(255) NOT NULL,
    "last_name" VARCHAR(255) NOT NULL,
    "email" VARCHAR(255) NOT NULL,
    "is_staff" BOOLEAN NOT NULL,
    "is_active" BOOLEAN NOT NULL,
    "date_joined" TIMESTAMP(0) WITHOUT TIME ZONE NOT NULL,
    "avatar" VARCHAR(255) NULL,
    "customer_id" BIGINT NULL,
    "subscription_id" BIGINT NULL
);
ALTER TABLE
    "users_customuser" ADD PRIMARY KEY("id");
ALTER TABLE
    "users_customuser" ADD CONSTRAINT "users_customuser_username_unique" UNIQUE("username");
CREATE INDEX "users_customuser_customer_id_index" ON
    "users_customuser"("customer_id");
CREATE INDEX "users_customuser_subscription_id_index" ON
    "users_customuser"("subscription_id");
CREATE TABLE "users_customuser_groups"(
    "id" INTEGER NOT NULL,
    "customuser_id" INTEGER NOT NULL,
    "group_id" INTEGER NOT NULL
);
ALTER TABLE
    "users_customuser_groups" ADD CONSTRAINT "users_customuser_groups_customuser_id_group_id_unique" UNIQUE("customuser_id", "group_id");
ALTER TABLE
    "users_customuser_groups" ADD PRIMARY KEY("id");
CREATE INDEX "users_customuser_groups_customuser_id_index" ON
    "users_customuser_groups"("customuser_id");
CREATE INDEX "users_customuser_groups_group_id_index" ON
    "users_customuser_groups"("group_id");
CREATE TABLE "users_customuser_user_permissions"(
    "id" INTEGER NOT NULL,
    "customuser_id" INTEGER NOT NULL,
    "permission_id" INTEGER NOT NULL
);
ALTER TABLE
    "users_customuser_user_permissions" ADD CONSTRAINT "users_customuser_user_permissions_customuser_id_permission_id_unique" UNIQUE("customuser_id", "permission_id");
ALTER TABLE
    "users_customuser_user_permissions" ADD PRIMARY KEY("id");
CREATE INDEX "users_customuser_user_permissions_customuser_id_index" ON
    "users_customuser_user_permissions"("customuser_id");
CREATE INDEX "users_customuser_user_permissions_permission_id_index" ON
    "users_customuser_user_permissions"("permission_id");
ALTER TABLE
    "auth_group_permissions" ADD CONSTRAINT "auth_group_permissions_permission_id_foreign" FOREIGN KEY("permission_id") REFERENCES "auth_permission"("id");
ALTER TABLE
    "teams_invitation" ADD CONSTRAINT "teams_invitation_team_id_foreign" FOREIGN KEY("team_id") REFERENCES "teams_membership"("id");
ALTER TABLE
    "account_emailconfirmation" ADD CONSTRAINT "account_emailconfirmation_email_address_id_foreign" FOREIGN KEY("email_address_id") REFERENCES "account_emailaddress"("id");
ALTER TABLE
    "teams_invitation" ADD CONSTRAINT "teams_invitation_invited_by_id_foreign" FOREIGN KEY("invited_by_id") REFERENCES "users_customuser"("id");
ALTER TABLE
    "account_emailaddress" ADD CONSTRAINT "account_emailaddress_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users_customuser"("id");
ALTER TABLE
    "auth_permission" ADD CONSTRAINT "auth_permission_content_type_id_foreign" FOREIGN KEY("content_type_id") REFERENCES "django_content_type"("id");
ALTER TABLE
    "teams_invitation" ADD CONSTRAINT "teams_invitation_accepted_by_id_foreign" FOREIGN KEY("accepted_by_id") REFERENCES "users_customuser"("id");
ALTER TABLE
    "pegasus_examples_payment" ADD CONSTRAINT "pegasus_examples_payment_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users_customuser"("id");
ALTER TABLE
    "users_customuser_groups" ADD CONSTRAINT "users_customuser_groups_group_id_foreign" FOREIGN KEY("group_id") REFERENCES "auth_group"("id");
ALTER TABLE
    "djstripe_subscription" ADD CONSTRAINT "djstripe_subscription_plan_id_foreign" FOREIGN KEY("plan_id") REFERENCES "djstripe_plan"("djstripe_id");
ALTER TABLE
    "teams_membership" ADD CONSTRAINT "teams_membership_team_id_foreign" FOREIGN KEY("team_id") REFERENCES "teams_team"("id");
ALTER TABLE
    "teams_membership" ADD CONSTRAINT "teams_membership_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users_customuser"("id");
ALTER TABLE
    "auth_group_permissions" ADD CONSTRAINT "auth_group_permissions_group_id_foreign" FOREIGN KEY("group_id") REFERENCES "auth_group"("id");
ALTER TABLE
    "django_admin_log" ADD CONSTRAINT "django_admin_log_content_type_id_foreign" FOREIGN KEY("content_type_id") REFERENCES "django_content_type"("id");
ALTER TABLE
    "djstripe_customer" ADD CONSTRAINT "djstripe_customer_subscriber_id_foreign" FOREIGN KEY("subscriber_id") REFERENCES "users_customuser"("id");
ALTER TABLE
    "users_customuser" ADD CONSTRAINT "users_customuser_subscription_id_foreign" FOREIGN KEY("subscription_id") REFERENCES "djstripe_subscription"("djstripe_id");
ALTER TABLE
    "users_customuser_user_permissions" ADD CONSTRAINT "users_customuser_user_permissions_customuser_id_foreign" FOREIGN KEY("customuser_id") REFERENCES "users_customuser"("id");
ALTER TABLE
    "teams_team" ADD CONSTRAINT "teams_team_subscription_id_foreign" FOREIGN KEY("subscription_id") REFERENCES "djstripe_subscription"("djstripe_id");
ALTER TABLE
    "users_customuser" ADD CONSTRAINT "users_customuser_customer_id_foreign" FOREIGN KEY("customer_id") REFERENCES "djstripe_customer"("djstripe_id");
ALTER TABLE
    "djstripe_subscription" ADD CONSTRAINT "djstripe_subscription_customer_id_foreign" FOREIGN KEY("customer_id") REFERENCES "djstripe_customer"("djstripe_id");
ALTER TABLE
    "pegasus_examples_employee" ADD CONSTRAINT "pegasus_examples_employee_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users_customuser"("id");
ALTER TABLE
    "django_admin_log" ADD CONSTRAINT "django_admin_log_user_id_foreign" FOREIGN KEY("user_id") REFERENCES "users_customuser"("id");
ALTER TABLE
    "djstripe_plan" ADD CONSTRAINT "djstripe_plan_product_id_foreign" FOREIGN KEY("product_id") REFERENCES "djstripe_product"("djstripe_id");
ALTER TABLE
    "users_customuser_groups" ADD CONSTRAINT "users_customuser_groups_customuser_id_foreign" FOREIGN KEY("customuser_id") REFERENCES "users_customuser"("id");
ALTER TABLE
    "users_customuser_user_permissions" ADD CONSTRAINT "users_customuser_user_permissions_permission_id_foreign" FOREIGN KEY("permission_id") REFERENCES "auth_permission"("id");