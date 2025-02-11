resource "aws_lambda_function" "transcribe" {
  function_name    = "transcribe"
  runtime          = "python3.13"
  handler          = "main.handler"
  filename         = data.archive_file.file.output_path
  source_code_hash = data.archive_file.file.output_base64sha256
  role             = aws_iam_role.role.arn

  environment {
    variables = {
      JOB_ROLE_ARN = aws_iam_role.job_role.arn
    }
  }
}

resource "aws_iam_role" "role" {
  name               = "transcribe-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy" "policies" {
  role   = aws_iam_role.role.name
  policy = data.aws_iam_policy_document.policies.json
}

resource "aws_iam_role" "job_role" {
  name               = "transcribe-job-role"
  assume_role_policy = data.aws_iam_policy_document.assume_job_role.json
}

resource "aws_iam_role_policy" "job_policies" {
  role   = aws_iam_role.role.name
  policy = data.aws_iam_policy_document.job_policies.json
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

data "aws_iam_policy_document" "assume_job_role" {
  statement {
    effect = "Allow"

    principals {
      type        = "Service"
      identifiers = ["transcribe.amazonaws.com"]
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

    actions = ["iam:PassRole"]

    resources = [aws_iam_role.job_role.arn]
  }

  statement {
    effect = "Allow"

    actions = ["transcribe:StartTranscriptionJob"]

    resources = ["arn:aws:transcribe:*:*:transcription-job/*"]
  }
}



data "aws_iam_policy_document" "job_policies" {
  statement {
    effect = "Allow"

    actions = ["s3:GetObject"]

    resources = ["${data.aws_s3_object.audio.arn}*"]
  }
  statement {
    effect = "Allow"

    actions = ["s3:PutObject"]

    resources = ["${data.aws_s3_object.transcription.arn}*"]
  }
}

data "archive_file" "file" {
  source_dir  = "${path.root}/../src"
  output_path = "lambda_payload.zip"
  type        = "zip"
}
data "aws_s3_bucket" "bucket" {
  bucket = "aws-ai-transcribe"
}

data "aws_s3_object" "audio" {
  bucket = data.aws_s3_bucket.bucket.id
  key    = "audio/"
}

data "aws_s3_object" "transcription" {
  bucket = data.aws_s3_bucket.bucket.id
  key    = "transcription/"
}
