output "remote_state_bucket_name" {
  value = module.remote_state.bucket_name
}

output "remote_state_dynamo_table_name" {
  value = module.remote_state.dynamo_table_name
}

output "verified_emails" {
  value = module.email.verified_emails
}
