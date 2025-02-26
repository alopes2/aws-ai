resource "aws_lambda_function" "transcribe" {
  function_name    = "transcribe"
  runtime          = "python3.13"
  handler          = "main.handler"
  filename         = data.archive_file.file.output_path
  source_code_hash = data.archive_file.file.output_base64sha256
  role             = aws_iam_role.role.arn

  environment {
    variables = {
      JOB_ROLE_ARN = "${aws_iam_role.job_role.arn}"
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
  role   = aws_iam_role.job_role.name
  policy = data.aws_iam_policy_document.job_policies.json
}

data "archive_file" "file" {
  source_dir  = "${path.root}/../app/src"
  output_path = "lambda_payload.zip"
  type        = "zip"
}

