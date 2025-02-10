resource "aws_lambda_function" "transcribe" {
  function_name    = "transcribe"
  runtime          = "python3.13"
  handler          = "main.handler"
  filename         = data.archive_file.file.output_path
  source_code_hash = data.archive_file.file.output_base64sha256
  role             = aws_iam_role.role.arn
}

resource "aws_iam_role" "role" {
  name               = "transcribe-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy" "policies" {
  role   = aws_iam_role.role.name
  policy = data.aws_iam_policy_document.policies.json
}

data "aws_iam_policy_document" "assume_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["lambda.amazonaws.com"]
    }

    actions = ["sts:AssumeRole"]
  }
}
data "aws_iam_policy_document" "policies" {
  statement {
    effect = "Allow"

    actions = [
      "logs:CreateLogGroup",
      "logs:CreateLogStream",
      "logs:PutLogEvents"
    ]

    resources = ["arn:aws:logs:*:*:*"]
  }

  statement {
    effect = "Allow"

    actions = ["transcribe:StartTranscriptionJob"]

    resources = ["arn:aws:transcribe:*:*:transcription-job/*"]
  }
}

data "archive_file" "file" {
  source_dir  = "${path.root}/../src"
  output_path = "lambda_payload.zip"
  type        = "zip"
}
