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

resource "aws_s3_object" "vocabulary_filters_folder" {
  bucket = aws_s3_bucket.bucket.id
  key    = "vocabulary_filters/"
}

resource "aws_s3_object" "vocabulary_filter" {
  bucket      = aws_s3_bucket.bucket.id
  key         = "${aws_s3_object.vocabulary_filters_folder.key}vocabulary_filter.txt"
  source      = "${path.module}/transcribe/vocabulary_filter.txt"
  source_hash = filemd5("${path.module}/transcribe/vocabulary_filter.txt")
}


resource "aws_s3_object" "clm" {
  bucket = aws_s3_bucket.bucket.id
  key    = "clm/"
}


resource "aws_s3_object" "training_data" {
  bucket = aws_s3_bucket.bucket.id
  key    = "${aws_s3_object.clm.key}training_data/"
}
resource "aws_s3_object" "tune_data" {
  bucket = aws_s3_bucket.bucket.id
  key    = "${aws_s3_object.clm.key}tune_data/"
}

resource "aws_s3_object" "nintendo_switch" {
  bucket      = aws_s3_bucket.bucket.id
  key         = "${aws_s3_object.training_data.key}NintendoSwitch.txt"
  source      = "${path.module}/training_data/NintendoSwitch.txt"
  source_hash = filemd5("${path.module}/training_data/NintendoSwitch.txt")
}

resource "aws_s3_object" "ps5" {
  bucket      = aws_s3_bucket.bucket.id
  key         = "${aws_s3_object.training_data.key}PlayStation5.txt"
  source      = "${path.module}/training_data/PlayStation5.txt"
  source_hash = filemd5("${path.module}/training_data/PlayStation5.txt")
}

resource "aws_s3_object" "xbox" {
  bucket      = aws_s3_bucket.bucket.id
  key         = "${aws_s3_object.training_data.key}XboxSeries.txt"
  source      = "${path.module}/training_data/XboxSeries.txt"
  source_hash = filemd5("${path.module}/training_data/XboxSeries.txt")
}

resource "aws_s3_object" "tune_data" {
  bucket      = aws_s3_bucket.bucket.id
  key         = "${aws_s3_object.tune_data.key}tune_data.txt"
  source      = "${path.module}/tune_data/tune_data.txt"
  source_hash = filemd5("${path.module}/tune_data/tune_data.txt")
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
