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
