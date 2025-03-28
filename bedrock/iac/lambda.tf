resource "aws_lambda_function" "bedrock" {
  function_name    = "bedrock"
  runtime          = "provided.al2023"
  handler          = "bootstrap"
  architectures    = ["arm64"]
  filename         = data.archive_file.file.output_path
  source_code_hash = data.archive_file.file.output_base64sha256
  role             = aws_iam_role.role.arn
  timeout          = 30

  environment {
    variables = {
      MODEL_ID = "${data.aws_bedrock_inference_profile.claude.inference_profile_id}"
    }
  }
}

resource "aws_iam_role" "role" {
  name               = "bedrock-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy" "policies" {
  role   = aws_iam_role.role.name
  policy = data.aws_iam_policy_document.policies.json
}

data "archive_file" "file" {
  source_file = "${path.module}/init_code/bootstrap"
  output_path = "lambda_payload.zip"
  type        = "zip"
}
