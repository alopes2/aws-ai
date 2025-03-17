resource "aws_transcribe_vocabulary" "vocabulary" {
  vocabulary_name     = "example"
  vocabulary_file_uri = "s3://${aws_s3_object.vocabulary.bucket}/${aws_s3_object.vocabulary.key}"
  language_code       = "en-US"
}

resource "aws_transcribe_vocabulary_filter" "filter" {
  vocabulary_filter_name     = "example"
  vocabulary_filter_file_uri = "s3://${aws_s3_object.vocabulary_filter.bucket}/${aws_s3_object.vocabulary_filter.key}"
  language_code              = "en-US"
}

resource "aws_transcribe_vocabulary_filter" "inline" {
  vocabulary_filter_name = "inline_example"
  language_code          = "en-US"
  words                  = ["content", "profane"]
}

resource "aws_transcribe_language_model" "model" {
  model_name      = "example"
  base_model_name = "Narrowband" # If you want a Wideband (audio sample rates over 16,000 Hz) or Narrowband (audio sample rates under 16,000 Hz) base model
  language_code   = "en-US"
  input_data_config {
    s3_uri               = "s3://${aws_s3_object.training_data.bucket}/${aws_s3_object.training_data.key}"
    tuning_data_s3_uri   = "s3://${aws_s3_object.tune_data.bucket}/${aws_s3_object.tune_data.key}"
    data_access_role_arn = aws_iam_role.transcribe_clm.arn
  }
}

resource "aws_iam_role" "transcribe_clm" {
  name               = "transcribe_clm"
  assume_role_policy = data.aws_iam_policy_document.transcribe_assume_role.json
}

resource "aws_iam_role_policy" "test_policy" {
  name = "transcribe_clm"
  role = aws_iam_role.transcribe_clm.id

  policy = data.aws_iam_policy_document.transcribe_s3.json
}
