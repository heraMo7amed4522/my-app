import { OnModuleInit } from '@nestjs/common';
import { UploadFileInput, DeleteFileInput, GetFileURLInput } from './upload.types';
export declare class UploadService implements OnModuleInit {
    private uploadServiceClient;
    onModuleInit(): Promise<void>;
    uploadFile(input: UploadFileInput): Promise<any>;
    deleteFile(input: DeleteFileInput): Promise<any>;
    getFileURL(input: GetFileURLInput): Promise<any>;
    private mapFileType;
}
