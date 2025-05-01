output "remote_state_bucket_name" {
  value = module.remote_state.bucket_name
}

output "remote_state_dynamo_table_name" {
  value = module.remote_state.dynamo_table_name
}
