resource "aws_s3_bucket" "bucket" {
  bucket = "aws-ai-transcribe"
}

resource "aws_s3_object" "media" {
  bucket = aws_s3_bucket.bucket.id
  key    = "media/"
}

resource "aws_s3_object" "transcription" {
  bucket = aws_s3_bucket.bucket.id
  key    = "transcription/"
}

resource "aws_s3_object" "vocabulary_folder" {
  bucket = aws_s3_bucket.bucket.id
  key    = "vocabularies/"
}

resource "aws_s3_object" "vocabulary" {
  bucket      = aws_s3_bucket.bucket.id
  key         = "${aws_s3_object.vocabulary_folder.key}vocabulary.txt"
  source      = "${path.module}/transcribe/vocabulary.txt"
  source_hash = filemd5("${path.module}/transcribe/vocabulary.txt")
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
    filter_prefix       = aws_s3_object.media.key
    events              = ["s3:ObjectCreated:*"]
    lambda_function_arn = aws_lambda_function.transcribe.arn
  }

  depends_on = [aws_lambda_permission.allow_bucket]
}
