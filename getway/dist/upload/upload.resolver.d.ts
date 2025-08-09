import { UploadService } from './upload.service';
import { UploadFileResponse, DeleteFileResponse, GetFileURLResponse, UploadFileInput, DeleteFileInput, GetFileURLInput } from './upload.types';
export declare class UploadResolver {
    private readonly uploadService;
    constructor(uploadService: UploadService);
    uploadFile(input: UploadFileInput): Promise<UploadFileResponse>;
    deleteFile(input: DeleteFileInput): Promise<DeleteFileResponse>;
    getFileURL(input: GetFileURLInput): Promise<GetFileURLResponse>;
}
