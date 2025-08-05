import { Injectable, OnModuleInit, BadRequestException } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';
import { UploadFileInput, DeleteFileInput, GetFileURLInput } from './upload.types';

@Injectable()
export class UploadService implements OnModuleInit {
  private uploadServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/upload.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const uploadProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.UPLOAD_SERVICE_URL || 'localhost:50061';
    
    this.uploadServiceClient = new uploadProto.upload.UploadService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  async uploadFile(input: UploadFileInput): Promise<any> {
    try {
      const { file, fileName, contentType, userId, fileType } = input;
      
      // Validate input
      if (!file) {
        throw new BadRequestException('File data is required');
      }
      if (!fileName) {
        throw new BadRequestException('File name is required');
      }
      if (!contentType) {
        throw new BadRequestException('Content type is required');
      }
      if (!userId) {
        throw new BadRequestException('User ID is required');
      }

      // Convert base64 to buffer
      let fileBuffer: Buffer;
      try {
        // Remove data URL prefix if present (e.g., "data:image/jpeg;base64,")
        const base64Data = file.includes(',') ? file.split(',')[1] : file;
        fileBuffer = Buffer.from(base64Data, 'base64');
      } catch (error) {
        throw new BadRequestException('Invalid file data format');
      }

      return new Promise((resolve, reject) => {
        this.uploadServiceClient.UploadFile(
          {
            file_data: fileBuffer,
            file_name: fileName,
            content_type: contentType,
            user_id: userId,
            file_type: this.mapFileType(fileType),
          },
          (error: any, response: any) => {
            if (error) {
              console.error('Upload error:', error);
              reject(new BadRequestException(`Upload failed: ${error.message}`));
            } else {
              resolve(response);
            }
          },
        );
      });
    } catch (error) {
      console.error('Upload service error:', error);
      if (error instanceof BadRequestException) {
        throw error;
      }
      throw new BadRequestException('Upload failed');
    }
  }

  async deleteFile(input: DeleteFileInput): Promise<any> {
    try {
      const { fileKey, userId } = input;

      if (!fileKey) {
        throw new BadRequestException('File key is required');
      }
      if (!userId) {
        throw new BadRequestException('User ID is required');
      }

      return new Promise((resolve, reject) => {
        this.uploadServiceClient.DeleteFile(
          {
            file_key: fileKey,
            user_id: userId,
          },
          (error: any, response: any) => {
            if (error) {
              console.error('Delete error:', error);
              reject(new BadRequestException(`Delete failed: ${error.message}`));
            } else {
              resolve(response);
            }
          },
        );
      });
    } catch (error) {
      console.error('Delete service error:', error);
      if (error instanceof BadRequestException) {
        throw error;
      }
      throw new BadRequestException('Delete failed');
    }
  }

  async getFileURL(input: GetFileURLInput): Promise<any> {
    try {
      const { fileKey, expiresIn } = input;

      if (!fileKey) {
        throw new BadRequestException('File key is required');
      }

      return new Promise((resolve, reject) => {
        this.uploadServiceClient.GetFileURL(
          {
            file_key: fileKey,
            expires_in: expiresIn || 3600,
          },
          (error: any, response: any) => {
            if (error) {
              console.error('Get URL error:', error);
              reject(new BadRequestException(`Get URL failed: ${error.message}`));
            } else {
              resolve(response);
            }
          },
        );
      });
    } catch (error) {
      console.error('Get URL service error:', error);
      if (error instanceof BadRequestException) {
        throw error;
      }
      throw new BadRequestException('Get URL failed');
    }
  }

  private mapFileType(fileType: string): number {
    switch (fileType) {
      case 'IMAGE':
        return 1;
      case 'VIDEO':
        return 2;
      default:
        return 0; // UNKNOWN
    }
  }
}