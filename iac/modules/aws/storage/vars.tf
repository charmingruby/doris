variable "bucket_name" {
  description = "The name of the S3 bucket"
  type        = string
}

variable "account_arn" {
  description = "The ARN of the account"
  type        = string
}

variable "tags" {
  type        = map(string)
  description = "The tags to apply to the resources"
}
