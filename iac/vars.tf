variable "project" {
  description = "The project name"
  type        = string
}

variable "environment" {
  description = "The project name"
  type        = string
  default     = "dev"
}

variable "managed_by" {
  description = "The project name"
  type        = string
  default     = "Terraform"
}

variable "created_at" {
  description = "The project name"
  type        = string
}

variable "sender_email" {
  type        = string
  description = "Test sender email to be verified"
}

variable "recipient_emails" {
  type        = list(string)
  description = "Test recipient emails to be verified"
}
