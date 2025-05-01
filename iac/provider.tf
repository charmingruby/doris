terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.49.0"
    }
  }

  backend "s3" {
    bucket         = "doris-dev-state"
    key            = "state/terraform.tfstate"
    region         = "us-east-1"
    dynamodb_table = "doris-dev-lock"
  }
}
