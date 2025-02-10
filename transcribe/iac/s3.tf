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

resource "aws_s3_bucket_notification" "bucket" {
  bucket = aws_s3_bucket.bucket.id
  lambda_function {
    filter_prefix       = "images/"
    events              = ["s3:ObjectCreated:*"]
    lambda_function_arn = data.aws_lambda_function.transcribe.arn
  }
}

resource "aws_lambda_permission" "eventbridge" {
  action        = "lambda:InvokeFunction"
  function_name = data.transcribe.function_name
  source_arn    = aws_s3_bucket.bucket.arn
  principal     = "s3.amazonaws.com"
}

data "aws_lambda_function" "transcribe" {
  function_name = "transcribe"
}

