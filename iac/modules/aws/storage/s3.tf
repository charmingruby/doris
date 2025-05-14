resource "aws_s3_bucket" "this" {
  bucket = var.bucket_name

  tags = merge(var.tags, {
    Name = var.bucket_name
  })
}

resource "aws_s3_bucket_public_access_block" "this" {
  bucket = aws_s3_bucket.this.id

  block_public_acls       = true
  block_public_policy     = true
  ignore_public_acls      = true
  restrict_public_buckets = true
}

resource "aws_s3_bucket_policy" "this" {
  bucket = aws_s3_bucket.this.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [
      {
        Sid    = "AllowOnlySpecificIAM"
        Effect = "Allow"
        Principal = {
          AWS = var.account_arn
        }
        Action   = ["s3:GetObject", "s3:PutObject"]
        Resource = "${aws_s3_bucket.this.arn}/*"
      }
    ]
  })
}
