resource "aws_lambda_function" "connection" {
  function_name    = "websocket-connection"
  runtime          = "provided.al2023"
  handler          = "bootstrap"
  architectures    = ["arm64"]
  filename         = data.archive_file.file.output_path
  source_code_hash = data.archive_file.file.output_base64sha256
  role             = aws_iam_role.connection_role.arn

  environment {
    variables = {
      TABLE_NAME = "${aws_dynamodb_table.connections.name}"
    }
  }
}

resource "aws_iam_role" "connection_role" {
  name               = "websocket-connection-lambda-role"
  assume_role_policy = data.aws_iam_policy_document.assume_role.json
}

resource "aws_iam_role_policy" "connection_policies" {
  role   = aws_iam_role.role.name
  policy = data.aws_iam_policy_document.connection_lambda_policies.json
}
