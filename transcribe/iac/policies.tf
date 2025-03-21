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

    resources = [
      "arn:aws:transcribe:*:*:transcription-job/*",
      aws_transcribe_vocabulary.vocabulary.arn,
      aws_transcribe_vocabulary_filter.filter.arn,
      aws_transcribe_language_model.model.arn
    ]
  }
}

data "aws_iam_policy_document" "job_policies" {
  statement {
    effect = "Allow"

    actions = ["s3:GetObject"]

    resources = ["${aws_s3_object.media.arn}*"]
  }
  statement {
    effect = "Allow"

    actions = ["s3:PutObject"]

    resources = ["${aws_s3_object.transcription.arn}*"]
  }
}

data "aws_iam_policy_document" "transcribe_assume_role" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["transcribe.amazonaws.com"]
    }
  }
}

data "aws_iam_policy_document" "transcribe_s3" {
  statement {
    effect = "Allow"

    actions = [
      "s3:GetObject",
    ]

    resources = [
      "${aws_s3_object.clm.arn}*",
    ]
  }
  statement {
    effect = "Allow"

    actions = [
      "s3:ListBucket"
    ]

    resources = [
      "${aws_s3_bucket.bucket.arn}",
    ]
  }
}
