import {
  MediaFormat,
  StartTranscriptionJobCommand,
  TranscribeClient,
  type StartTranscriptionJobRequest,
} from '@aws-sdk/client-transcribe';
import type { S3Event } from 'aws-lambda';

const JOB_ROLE_ARN = process.env.JOB_ROLE_ARN;
const OUTPUT_KEY = process.env.OUTPUT_KEY; // transcription/
const transcribeClient = new TranscribeClient({});

export const handler = async (event: S3Event) => {
  for (let record of event.Records) {
    const bucket = record.s3.bucket.name;
    const key = record.s3.object.key;

    const fileInput = `s3://${bucket}/${key}`;
    const mediaFormat = fileInput.split('.').at(-1);

    if (
      !mediaFormat ||
      !Object.values(MediaFormat).includes(mediaFormat as MediaFormat)
    ) {
      console.warn('No media format for this file');
      return;
    }

    const jobName = key.replace('/', '_') + Date.now();

    const jobRequest: StartTranscriptionJobRequest = {
      TranscriptionJobName: jobName,
      Media: { MediaFileUri: fileInput },
      MediaFormat: mediaFormat as MediaFormat,
      LanguageCode: 'en-US',
      OutputBucketName: bucket,
      OutputKey: `${OUTPUT_KEY}${jobName}.json`,
      JobExecutionSettings: {
        DataAccessRoleArn: JOB_ROLE_ARN,
      },
    };

    const job = new StartTranscriptionJobCommand(jobRequest);

    try {
      const response = await transcribeClient.send(job);
      console.log(
        'Finished job %s. Data %s',
        jobName,
        response.TranscriptionJob
      );
    } catch (error: any) {
      console.error(
        "Couldn't start transcription job %s. Error: %s",
        jobName,
        error
      );
      throw error;
    }
  }
};
