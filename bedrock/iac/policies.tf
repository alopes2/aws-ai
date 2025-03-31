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

    actions = [
      # "bedrock:InvokeModel",
      "bedrock:InvokeModelWithResponseStream"
    ]

    resources = [
      data.aws_bedrock_foundation_model.model.model_arn,
      "arn:aws:bedrock:*::foundation-model/${data.aws_bedrock_foundation_model.claude.model_id}",
      data.aws_bedrock_inference_profile.claude.inference_profile_arn
    ]
  }
}

data "aws_iam_policy_document" "connection_lambda_policies" {
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

    actions = [
      "dynamodb:PutItem",
      "dynamodb:GetItem",
    ]

    resources = [
      aws_dynamodb_table.connections.arn
    ]
  }
}
