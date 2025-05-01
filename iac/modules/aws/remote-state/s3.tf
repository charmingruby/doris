resource "aws_s3_bucket" "this" {
  bucket        = "${var.tags.project}-${var.tags.environment}-state"
  force_destroy = true

  lifecycle {
    prevent_destroy = true
  }

  tags = merge(var.tags, {
    Name = "${var.tags.project}-${var.tags.environment}-state"
  })
}

resource "aws_s3_bucket_versioning" "this" {
  bucket = aws_s3_bucket.this.id

  lifecycle {
    prevent_destroy = true
  }

  versioning_configuration {
    status = "Enabled"
  }
}
