resource "aws_transcribe_vocabulary" "vocabulary" {
  vocabulary_name     = "example"
  vocabulary_file_uri = "s3://${aws_s3_object.vocabulary.bucket}/${aws_s3_object.vocabulary.key}"
  language_code       = "en-US"
}
