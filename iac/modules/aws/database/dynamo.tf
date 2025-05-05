resource "aws_dynamodb_table" "notifications" {
  name         = "Notifications"
  billing_mode = "PAY_PER_REQUEST"

  hash_key  = "PK"
  range_key = "SK"

  attribute {
    name = "PK"
    type = "S"
  }

  attribute {
    name = "SK"
    type = "S"
  }

  attribute {
    name = "timestamp"
    type = "N"
  }

  attribute {
    name = "correlationId"
    type = "S"
  }

  global_secondary_index {
    name            = "CorrelationIndex"
    hash_key        = "correlationId"
    range_key       = "timestamp"
    projection_type = "ALL"
  }
}
