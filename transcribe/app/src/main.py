import boto3
import logging
import time
import os
from urllib.parse import unquote_plus
from botocore.exceptions import ClientError

logger = logging.getLogger()
logger.setLevel("INFO")

transcribe_client = boto3.client("transcribe")

def handler(event, context):
  for record in event["Records"]:
    bucket = record['s3']['bucket']['name']
    key = unquote_plus(record['s3']['object']['key'])
    fileInput = f"s3://{bucket}/{key}"

    job_name = key + time.time()
    
    media_format = key.split(".")[-1]

    try:
        job_args = {
            "TranscriptionJobName": job_name,
            "Media": {"MediaFileUri": fileInput},
            "MediaFormat": media_format,
            "LanguageCode": "en-US",
            "OutputBucketName": bucket,
            "OutputKey": f"transcription/{job_name}.txt"
        }

        response = transcribe_client.start_transcription_job(**job_args)
        job = response["TranscriptionJob"]
        logger.info("Started transcription job %s.", job_name)
    except ClientError:
        logger.exception("Couldn't start transcription job %s.", job_name)
        raise
    else:
        return job