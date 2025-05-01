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
