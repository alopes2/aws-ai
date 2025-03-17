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
  model_name = "example"

  //NarrowBand: Use this option for audio with a sample rate of less than 16,000 Hz. This model type is typically used for telephone conversations recorded at 8,000 Hz.
  // WideBand: Use this option for audio with a sample rate greater than or equal to 16,000 Hz.
  base_model_name = "NarrowBand"

  language_code = "en-US"
  input_data_config {
    s3_uri               = "s3://${aws_s3_object.training_data.bucket}/${aws_s3_object.training_data.key}"
    tuning_data_s3_uri   = "s3://${aws_s3_object.tune_data.bucket}/${aws_s3_object.tune_data.key}"
    data_access_role_arn = aws_iam_role.transcribe_clm.arn
  }

  depends_on = [aws_iam_role_policy.transcribe_clm_policy]
}

resource "aws_iam_role" "transcribe_clm" {
  name               = "transcribe_clm"
  assume_role_policy = data.aws_iam_policy_document.transcribe_assume_role.json
}

resource "aws_iam_role_policy" "transcribe_clm_policy" {
  name = "transcribe_clm"
  role = aws_iam_role.transcribe_clm.id

  policy = data.aws_iam_policy_document.transcribe_s3.json
}
