resource "aws_ses_email_identity" "sender" {
  email = var.sender_email
}

resource "aws_ses_email_identity" "recipients" {
  count = length(var.recipient_emails)
  email = var.recipient_emails[count.index]
}
