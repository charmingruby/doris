variable "sender_email" {
  type        = string
  description = "Test sender email to be verified"
}

variable "recipient_emails" {
  type        = list(string)
  description = "Test recipient emails to be verified"
}

variable "tags" {
  type        = map(string)
  description = "The tags to apply to the resources"
}
