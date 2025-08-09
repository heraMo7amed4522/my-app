import { ObjectType, Field, InputType, Int, registerEnumType } from '@nestjs/graphql';
import { Readable } from 'stream';
import { ErrorDetails } from '../shared/shared.types';

// Define FileUpload interface
export interface FileUpload {
  filename: string;
  mimetype: string;
  encoding: string;
  createReadStream(): Readable;
}

export enum FileType {
  UNKNOWN = 'UNKNOWN',
  IMAGE = 'IMAGE',
  VIDEO = 'VIDEO',
}

registerEnumType(FileType, {
  name: 'FileType',
});

@ObjectType()
export class UploadSuccess {
  @Field()
  fileUrl: string;

  @Field()
  fileKey: string;

  @Field()
  fileId: string;

  @Field(() => Int)
  fileSize: number;

  @Field()
  uploadedAt: string;
}

@ObjectType()
export class DeleteSuccess {
  @Field()
  fileKey: string;

  @Field()
  deletedAt: string;
}

@ObjectType()
export class URLSuccess {
  @Field()
  signedUrl: string;

  @Field()
  expiresAt: string;
}


// Remove the duplicate ErrorDetails definition (lines 58-70)
// The ErrorDetails class should be imported from shared.types.ts instead

@ObjectType()
export class UploadFileResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => UploadSuccess, { nullable: true })
  success?: UploadSuccess;

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;
}

@ObjectType()
export class DeleteFileResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => DeleteSuccess, { nullable: true })
  success?: DeleteSuccess;

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;
}

@ObjectType()
export class GetFileURLResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => URLSuccess, { nullable: true })
  success?: URLSuccess;

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;
}

@InputType()
export class UploadFileInput {
  @Field(() => String)
  file: string; // Base64 encoded file data

  @Field()
  fileName: string;

  @Field()
  contentType: string; // MIME type (e.g., 'image/jpeg', 'video/mp4')

  @Field()
  userId: string;

  @Field(() => FileType)
  fileType: FileType;
}

@InputType()
export class DeleteFileInput {
  @Field()
  fileKey: string;

  @Field()
  userId: string;
}

@InputType()
export class GetFileURLInput {
  @Field()
  fileKey: string;

  @Field(() => Int, { nullable: true, defaultValue: 3600 })
  expiresIn?: number;
}