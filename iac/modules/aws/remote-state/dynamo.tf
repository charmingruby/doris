resource "aws_dynamodb_table" "this" {
  name         = "${var.tags.project}-${var.tags.environment}-lock"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  lifecycle {
    prevent_destroy = true
  }

  attribute {
    name = "LockID"
    type = "S"
  }

  tags = merge(var.tags, {
    Name = "${var.tags.project}-${var.tags.environment}-lock"
  })
}
