terraform {
  required_providers {
    aws = {
      source  = "hashicorp/aws"
      version = "~> 5.84"
    }
  }

  backend "s3" {
    bucket = "terraform-medium-api-notification"
    key    = "bedrock/state.tfstate"
  }
}

provider "aws" {

}
