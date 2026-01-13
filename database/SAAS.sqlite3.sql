BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS currency (
	currency_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	currency VARCHAR(3) NOT NULL, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE (currency)
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
CREATE TABLE IF NOT EXISTS deployment (
	deployment_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	deployment VARCHAR(200), 
	deployment_desc TEXT, 
	product_id INTEGER NOT NULL, 
	provider_id INTEGER NOT NULL, 
	terraform_template TEXT NOT NULL, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	FOREIGN KEY(product_id) REFERENCES product (product_id), 
	FOREIGN KEY(provider_id) REFERENCES provider (provider_id)
);
CREATE TABLE IF NOT EXISTS env (
	env_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	env_name VARCHAR(200) NOT NULL, 
	env_value TEXT NOT NULL, 
	tenant_id INTEGER NOT NULL, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE (env_name), 
	FOREIGN KEY(tenant_id) REFERENCES tenant (tenant_id)
);
CREATE TABLE IF NOT EXISTS interval (
	interval_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	interval VARCHAR(100) NOT NULL, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE (interval)
);
CREATE TABLE IF NOT EXISTS payment_plan (
	payment_plan_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	plan_id INTEGER NOT NULL, 
	deployment_id INTEGER NOT NULL, 
	product_id INTEGER NOT NULL, 
	price FLOAT NOT NULL, 
	currency_id INTEGER NOT NULL, 
	interval_id INTEGER NOT NULL, 
	stripe_price_id VARCHAR(255), 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	FOREIGN KEY(plan_id) REFERENCES "plan" (plan_id), 
	FOREIGN KEY(deployment_id) REFERENCES deployment (deployment_id), 
	FOREIGN KEY(product_id) REFERENCES product (product_id), 
	FOREIGN KEY(currency_id) REFERENCES currency (currency_id), 
	FOREIGN KEY(interval_id) REFERENCES interval (interval_id)
);
CREATE TABLE IF NOT EXISTS "plan" (
	plan_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	"plan" VARCHAR(3) NOT NULL, 
	price FLOAT, 
	deployment_id INTEGER NOT NULL, 
	product_id INTEGER NOT NULL, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE ("plan"), 
	FOREIGN KEY(deployment_id) REFERENCES deployment (deployment_id), 
	FOREIGN KEY(product_id) REFERENCES product (product_id)
);
CREATE TABLE IF NOT EXISTS product (
	product_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	product VARCHAR(255) NOT NULL, 
	description TEXT, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN
);
CREATE TABLE IF NOT EXISTS provider (
	provider_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	provider VARCHAR(3) NOT NULL, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE (provider)
);
CREATE TABLE IF NOT EXISTS subscription (
	subscription_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	tenant_id INTEGER NOT NULL, 
	plan_id INTEGER NOT NULL, 
	deployment_id INTEGER NOT NULL, 
	payment_plan_id INTEGER NOT NULL, 
	product_id INTEGER NOT NULL, 
	terraform_outputs TEXT, 
	tf_public_ip VARCHAR(100), 
	tf_public_dns VARCHAR(255), 
	terraform_state TEXT, 
	terraform_lock TEXT, 
	tf_err_msg TEXT, 
	deployed BOOLEAN, 
	active BOOLEAN, 
	stripe_subscription_id VARCHAR(255), 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	FOREIGN KEY(tenant_id) REFERENCES tenant (tenant_id), 
	FOREIGN KEY(plan_id) REFERENCES "plan" (plan_id), 
	FOREIGN KEY(deployment_id) REFERENCES deployment (deployment_id), 
	FOREIGN KEY(payment_plan_id) REFERENCES payment_plan (payment_plan_id), 
	FOREIGN KEY(product_id) REFERENCES product (product_id), 
	UNIQUE (stripe_subscription_id)
);
CREATE TABLE IF NOT EXISTS tenant (
	tenant_id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
	tenant VARCHAR(200) NOT NULL, 
	email VARCHAR(200) NOT NULL, 
	password VARCHAR(200), 
	phone VARCHAR(200), 
	address VARCHAR(200), 
	currency_id INTEGER, 
	active BOOLEAN, 
	user_id INTEGER, 
	app_id INTEGER, 
	created_at DATETIME, 
	updated_at DATETIME, 
	excluded BOOLEAN, 
	UNIQUE (email), 
	FOREIGN KEY(currency_id) REFERENCES currency (currency_id)
);
INSERT INTO "currency" VALUES (1,'USD','2026-01-11 12:00:33.405715','2026-01-11 12:00:33.405715',0);
INSERT INTO "currency" VALUES (2,'EUR','2026-01-11 12:00:33.405715','2026-01-11 12:00:33.405715',0);
INSERT INTO "deployment" VALUES (1,'FREE','Free Central Set Plan',1,1,'# ==========================================
# TERRAFORM CONFIGURATION
# ==========================================

terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = ">= 4.60, < 7.0"
    }
  }
  required_version = ">= 1.3"
}

# ==========================================
# VARIABLES
# ==========================================

variable "aws_region" {
  description = "AWS region"
  type        = string
  default     = "us-east-1" # us-east-1
}

variable "instance_type" {
  description = "EC2 instance type"
  type        = string
  default     = "t3.micro"
}

# ==========================================
# PROVIDER CONFIGURATION
# ==========================================

provider "aws" {
  region = var.aws_region
}

# ==========================================
# DATA SOURCES
# ==========================================

# Latest Amazon Linux 2 AMI
data "aws_ami" "amazon_linux2" {
  most_recent = true
  owners      = ["amazon"]

  filter {
    name   = "name"
    values = ["amzn2-ami-hvm-*-x86_64-gp2"]
  }
}

# ==========================================
# NETWORKING RESOURCES
# ==========================================

# Simple VPC
resource "aws_vpc" "this" {
  cidr_block           = "10.0.0.0/16"
  enable_dns_hostnames = true

  tags = {
    Name = "central-set-vpc"
  }
}

# Public subnet
resource "aws_subnet" "public" {
  vpc_id                  = aws_vpc.this.id
  cidr_block              = "10.0.1.0/24"
  availability_zone       = "us-east-1a"
  map_public_ip_on_launch = true

  tags = {
    Name = "central-set-public-subnet"
  }
}

# Internet Gateway
resource "aws_internet_gateway" "igw" {
  vpc_id = aws_vpc.this.id

  tags = {
    Name = "central-set-igw"
  }
}

# Route table for public subnet
resource "aws_route_table" "public" {
  vpc_id = aws_vpc.this.id

  route {
    cidr_block = "0.0.0.0/0"
    gateway_id = aws_internet_gateway.igw.id
  }

  tags = {
    Name = "central-set-public-rt"
  }
}

# Route table association
resource "aws_route_table_association" "ra" {
  subnet_id      = aws_subnet.public.id
  route_table_id = aws_route_table.public.id
}

# ==========================================
# SECURITY GROUPS
# ==========================================

# Security group allowing HTTP, HTTPS, and SSH
resource "aws_security_group" "http" {
  name        = "allow_http"
  description = "Allow HTTP, HTTPS, and SSH traffic"
  vpc_id      = aws_vpc.this.id

  # Allow HTTP from anywhere
  ingress {
    description = "HTTP"
    from_port   = 80
    to_port     = 80
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow HTTPS from anywhere
  ingress {
    description = "HTTPS"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow SSH access
  ingress {
    description = "SSH"
    from_port   = 22
    to_port     = 22
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  # Allow all outbound traffic
  egress {
    description = "All outbound traffic"
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }

  tags = {
    Name = "central-set-security-group"
  }
}

# ==========================================
# KEY PAIR
# ==========================================

# SSH Key Pair
# Note: Generate the key pair first with:
# ssh-keygen -t rsa -b 4096 -f ~/.ssh/centralset-key
#resource "aws_key_pair" "deployer" {
#  key_name   = "centralset-key"
#  public_key = file("~/.ssh/centralset-key.pub")

#  tags = {
#    Name = "central-set-key-pair"
#  }
#}

# ==========================================
# COMPUTE RESOURCES
# ==========================================

# EC2 Instance
resource "aws_instance" "app" {
  ami                         = data.aws_ami.amazon_linux2.id
  instance_type               = var.instance_type
  subnet_id                   = aws_subnet.public.id
  associate_public_ip_address = true
  vpc_security_group_ids      = [aws_security_group.http.id]
  key_name                    = aws_key_pair.deployer.key_name

  user_data = <<-EOF
              #!/bin/bash
              yum update -y
              amazon-linux-extras install docker -y || yum install -y docker
              systemctl enable --now docker
              docker run docker.io/realdatadriven/central-set-go:latest --init
              docker run docker.io/realdatadriven/central-set-go:latest --init --dbname ETLX
              docker run --rm -d -p 80:4444 docker.io/realdatadriven/central-set-go:latest
              EOF

  tags = {
    Name = "central-set-free-demo"
  }
}

# ==========================================
# OUTPUTS
# ==========================================

output "public_ip" {
  description = "Public IP address of the EC2 instance"
  value       = aws_instance.app.public_ip
}

output "public_dns" {
  description = "Public DNS name of the EC2 instance"
  value       = aws_instance.app.public_dns
}

output "url" {
  description = "Public URL to reach the app"
  value       = "http://${aws_instance.app.public_ip}"
}

output "ssh_command" {
  description = "SSH command to connect to the instance"
  value       = "ssh -i ~/.ssh/centralset-key ec2-user@${aws_instance.app.public_ip}"
}

# ==========================================
# ADDITIONAL INFORMATION
# ==========================================

# To use this configuration:
# 1. Generate SSH key pair:
#    ssh-keygen -t rsa -b 4096 -f ~/.ssh/centralset-key
# 
# 2. Initialize Terraform:
#    terraform init
# 
# 3. Plan the deployment:
#    terraform plan
# 
# 4. Apply the configuration:
#    terraform apply
# 
# 5. Connect to the instance:
#    ssh -i ~/.ssh/centralset-key ec2-user@<public_ip>
',1,1,3,'2026-01-11T12:00:33.405715Z','2026-01-13 09:17:28.072430622-01:00',0);
INSERT INTO "interval" VALUES (1,'Daily','2026-01-11 12:00:33.405715','2026-01-11 12:00:33.405715',0);
INSERT INTO "interval" VALUES (2,'Monthly','2026-01-11 12:00:33.405715','2026-01-11 12:00:33.405715',0);
INSERT INTO "interval" VALUES (3,'Yearly','2026-01-11 12:00:33.405715','2026-01-11 12:00:33.405715',0);
INSERT INTO "payment_plan" VALUES (1,1,1,1,0.0,1,2,NULL,1,1,3,'2026-01-11 12:00:33.405715','2026-01-11 12:00:33.405715',0);
INSERT INTO "plan" VALUES (1,'Free',0.0,1,1,'2026-01-11 12:00:33.405715','2026-01-11 12:00:33.405715',0);
INSERT INTO "product" VALUES (1,'CS','Central Set',1,1,3,'2026-01-11 12:00:33.405715','2026-01-11 12:00:33.405715',0);
INSERT INTO "provider" VALUES (1,'AWS','2026-01-11 12:00:33.405715','2026-01-11 12:00:33.405715',0);
INSERT INTO "provider" VALUES (2,'GCP','2026-01-11 12:00:33.405715','2026-01-11 12:00:33.405715',0);
INSERT INTO "provider" VALUES (3,'Azure','2026-01-11 12:00:33.405715','2026-01-11 12:00:33.405715',0);
INSERT INTO "subscription" VALUES (1,1,1,1,1,1,'null','','','{
  "version": 4,
  "terraform_version": "1.14.3",
  "serial": 1,
  "lineage": "89305648-88a8-80d0-7994-1a608896946b",
  "outputs": {},
  "resources": [
    {
      "mode": "data",
      "type": "aws_ami",
      "name": "amazon_linux2",
      "provider": "provider[\"registry.terraform.io/hashicorp/aws\"]",
      "instances": [
        {
          "schema_version": 0,
          "attributes": {
            "allow_unsafe_filter": null,
            "architecture": "x86_64",
            "arn": "arn:aws:ec2:us-east-1::image/ami-0fcb14c72c80bdef2",
            "block_device_mappings": [
              {
                "device_name": "/dev/xvda",
                "ebs": {
                  "delete_on_termination": "true",
                  "encrypted": "false",
                  "iops": "0",
                  "snapshot_id": "snap-00ea161610c9fd4bd",
                  "throughput": "0",
                  "volume_initialization_rate": "0",
                  "volume_size": "8",
                  "volume_type": "gp2"
                },
                "no_device": "",
                "virtual_name": ""
              }
            ],
            "boot_mode": "",
            "creation_date": "2026-01-02T18:47:06.000Z",
            "deprecation_time": "2026-04-02T00:00:00.000Z",
            "description": "Amazon Linux 2 AMI 2.0.20260105.1 x86_64 HVM gp2",
            "ena_support": true,
            "executable_users": null,
            "filter": [
              {
                "name": "name",
                "values": [
                  "amzn2-ami-hvm-*-x86_64-gp2"
                ]
              }
            ],
            "hypervisor": "xen",
            "id": "ami-0fcb14c72c80bdef2",
            "image_id": "ami-0fcb14c72c80bdef2",
            "image_location": "amazon/amzn2-ami-hvm-2.0.20260105.1-x86_64-gp2",
            "image_owner_alias": "amazon",
            "image_type": "machine",
            "imds_support": "",
            "include_deprecated": false,
            "kernel_id": "",
            "last_launched_time": "",
            "most_recent": true,
            "name": "amzn2-ami-hvm-2.0.20260105.1-x86_64-gp2",
            "name_regex": null,
            "owner_id": "137112412989",
            "owners": [
              "amazon"
            ],
            "platform": "",
            "platform_details": "Linux/UNIX",
            "product_codes": [],
            "public": true,
            "ramdisk_id": "",
            "region": "us-east-1",
            "root_device_name": "/dev/xvda",
            "root_device_type": "ebs",
            "root_snapshot_id": "snap-00ea161610c9fd4bd",
            "sriov_net_support": "simple",
            "state": "available",
            "state_reason": {
              "code": "UNSET",
              "message": "UNSET"
            },
            "tags": {},
            "timeouts": null,
            "tpm_support": "",
            "uefi_data": null,
            "usage_operation": "RunInstances",
            "virtualization_type": "hvm"
          },
          "sensitive_attributes": [],
          "identity_schema_version": 0
        }
      ]
    }
  ],
  "check_results": null
}
','# This file is maintained automatically by "terraform init".
# Manual edits may be lost in future updates.

provider "registry.terraform.io/hashicorp/aws" {
  version     = "6.28.0"
  constraints = ">= 4.60.0, < 7.0.0"
  hashes = [
    "h1:wzZdGs0FFmNqIgPyo9tKnGKJ37BGNSgwRrEXayL29+0=",
    "zh:0ba0d5eb6e0c6a933eb2befe3cdbf22b58fbc0337bf138f95bf0e8bb6e6df93e",
    "zh:23eacdd4e6db32cf0ff2ce189461bdbb62e46513978d33c5de4decc4670870ec",
    "zh:307b06a15fc00a8e6fd243abde2cbe5112e9d40371542665b91bec1018dd6e3c",
    "zh:37a02d5b45a9d050b9642c9e2e268297254192280df72f6e46641daca52e40ec",
    "zh:3da866639f07d92e734557d673092719c33ede80f4276c835bf7f231a669aa33",
    "zh:480060b0ba310d0f6b6a14d60b276698cb103c48fd2f7e2802ae47c963995ec6",
    "zh:57796453455c20db80d9168edbf125bf6180e1aae869de1546a2be58e4e405ec",
    "zh:69139cba772d4df8de87598d8d8a2b1b4b254866db046c061dccc79edb14e6b9",
    "zh:7312763259b859ff911c5452ca8bdf7d0be6231c5ea0de2df8f09d51770900ac",
    "zh:8d2d6f4015d3c155d7eb53e36f019a729aefb46ebfe13f3a637327d3a1402ecc",
    "zh:94ce589275c77308e6253f607de96919b840c2dd36c44aa798f693c9dd81af42",
    "zh:9b12af85486a96aedd8d7984b0ff811a4b42e3d88dad1a3fb4c0b580d04fa425",
    "zh:adaceec6a1bf4f5df1e12bd72cf52b72087c72efed078aef636f8988325b1a8b",
    "zh:d37be1ce187d94fd9df7b13a717c219964cd835c946243f096c6b230cdfd7e92",
    "zh:fe6205b5ca2ff36e68395cb8d3ae10a3728f405cdbcd46b206a515e1ebcf17a1",
  ]
}
','deploy: exit status 1

Error: creating EC2 VPC: operation error EC2: CreateVpc, https response error StatusCode: 400, RequestID: d3577c53-3809-48eb-a4a2-02dbbee0419d, api error VpcLimitExceeded: The maximum number of VPCs has been reached.

  with aws_vpc.this,
  on main.tf line 59, in resource "aws_vpc" "this":
  59: resource "aws_vpc" "this" {

',1,1,NULL,1,3,'2026-01-12T22:18:36.907978Z','2026-01-13 08:42:36.250822221-01:00',0);
INSERT INTO "tenant" VALUES (1,'Tenant 1','tenant1@cs.io','$2a$12$DfApyEJ1yVaVgiUe9G0L6elnsRFhWoFYJHd3jEgJkfZHUqmOea7ei',NULL,NULL,1,1,1,3,'2026-01-11T12:00:33.405715Z','2026-01-11 12:19:40.652971511-01:00',0);
COMMIT;
