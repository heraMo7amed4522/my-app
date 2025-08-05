import { Resolver, Mutation, Query, Args } from '@nestjs/graphql';
import { UploadService } from './upload.service';
import {
  UploadFileResponse,
  DeleteFileResponse,
  GetFileURLResponse,
  UploadFileInput,
  DeleteFileInput,
  GetFileURLInput,
} from './upload.types';

@Resolver()
export class UploadResolver {
  constructor(private readonly uploadService: UploadService) {}

  @Mutation(() => UploadFileResponse)
  async uploadFile(
    @Args('input') input: UploadFileInput,
  ): Promise<UploadFileResponse> {
    try {
      return await this.uploadService.uploadFile(input);
    } catch (error) {
      console.error('Upload resolver error:', error);
      return {
        statusCode: 500,
        message: error.message || 'Upload failed',
        error: {
          code: 500,
          message: error.message || 'Internal server error',
          details: [error.message || 'Upload failed'],
          timestamp: new Date().toISOString(),
        },
      };
    }
  }

  @Mutation(() => DeleteFileResponse)
  async deleteFile(
    @Args('input') input: DeleteFileInput,
  ): Promise<DeleteFileResponse> {
    try {
      return await this.uploadService.deleteFile(input);
    } catch (error) {
      console.error('Delete resolver error:', error);
      return {
        statusCode: 500,
        message: error.message || 'Delete failed',
        error: {
          code: 500,
          message: error.message || 'Internal server error',
          details: [error.message || 'Delete failed'],
          timestamp: new Date().toISOString(),
        },
      };
    }
  }

  @Query(() => GetFileURLResponse)
  async getFileURL(
    @Args('input') input: GetFileURLInput,
  ): Promise<GetFileURLResponse> {
    try {
      return await this.uploadService.getFileURL(input);
    } catch (error) {
      console.error('Get URL resolver error:', error);
      return {
        statusCode: 500,
        message: error.message || 'Get URL failed',
        error: {
          code: 500,
          message: error.message || 'Internal server error',
          details: [error.message || 'Get URL failed'],
          timestamp: new Date().toISOString(),
        },
      };
    }
  }
}