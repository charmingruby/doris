output "verified_emails" {
  value = concat(
    [aws_ses_email_identity.sender.email],
    [for recipient in aws_ses_email_identity.recipients : recipient.email]
  )
}
