resource "aws_s3_bucket" "bucket" {
  bucket = "aws-ai-transcribe"
}

resource "aws_s3_object" "audio" {
  bucket = aws_s3_bucket.bucket.id
  key    = "audio/"
}

resource "aws_s3_object" "transcription" {
  bucket = aws_s3_bucket.bucket.id
  key    = "transcription/"
}

resource "aws_s3_bucket_policy" "policy" {
  bucket = aws_s3_bucket.bucket.id
  policy = data.aws_iam_policy_document.bucket_policy.json
}

data "aws_iam_policy_document" "bucket_policy" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["transcribe.amazonaws.com"]
      # identifiers = ["transcribe.streaming.amazonaws.com"] // This is permissions for streaming
    }

    actions = ["s3:GetObject"]

    resources = ["${aws_s3_object.audio.arn}*"]
  }
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["transcribe.amazonaws.com"]
    }

    actions = ["s3:PutObject"]

    resources = ["${aws_s3_object.transcription.arn}*"]
  }
}

resource "aws_lambda_permission" "allow_bucket" {
  action        = "lambda:InvokeFunction"
  function_name = data.aws_lambda_function.transcribe.arn
  source_arn    = aws_s3_bucket.bucket.arn
  principal     = "s3.amazonaws.com"
}

resource "aws_s3_bucket_notification" "bucket" {
  bucket = aws_s3_bucket.bucket.id
  lambda_function {
    filter_prefix       = aws_s3_object.audio.key
    events              = ["s3:ObjectCreated:*"]
    lambda_function_arn = data.aws_lambda_function.transcribe.arn
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}

data "aws_lambda_function" "transcribe" {
  function_name = "transcribe"
}

