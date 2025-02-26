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

resource "aws_lambda_permission" "allow_bucket" {
  action        = "lambda:InvokeFunction"
  function_name = aws_lambda_function.transcribe.arn
  source_arn    = aws_s3_bucket.bucket.arn
  principal     = "s3.amazonaws.com"
}

resource "aws_s3_bucket_notification" "bucket" {
  bucket = aws_s3_bucket.bucket.id
  lambda_function {
    filter_prefix       = aws_s3_object.audio.key
    events              = ["s3:ObjectCreated:*"]
    lambda_function_arn = aws_lambda_function.transcribe.arn
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}