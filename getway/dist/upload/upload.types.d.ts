import { Readable } from 'stream';
import { ErrorDetails } from '../shared/shared.types';
export interface FileUpload {
    filename: string;
    mimetype: string;
    encoding: string;
    createReadStream(): Readable;
}
export declare enum FileType {
    UNKNOWN = "UNKNOWN",
    IMAGE = "IMAGE",
    VIDEO = "VIDEO"
}
export declare class UploadSuccess {
    fileUrl: string;
    fileKey: string;
    fileId: string;
    fileSize: number;
    uploadedAt: string;
}
export declare class DeleteSuccess {
    fileKey: string;
    deletedAt: string;
}
export declare class URLSuccess {
    signedUrl: string;
    expiresAt: string;
}
export declare class UploadFileResponse {
    statusCode: number;
    message: string;
    success?: UploadSuccess;
    error?: ErrorDetails;
}
export declare class DeleteFileResponse {
    statusCode: number;
    message: string;
    success?: DeleteSuccess;
    error?: ErrorDetails;
}
export declare class GetFileURLResponse {
    statusCode: number;
    message: string;
    success?: URLSuccess;
    error?: ErrorDetails;
}
export declare class UploadFileInput {
    file: string;
    fileName: string;
    contentType: string;
    userId: string;
    fileType: FileType;
}
export declare class DeleteFileInput {
    fileKey: string;
    userId: string;
}
export declare class GetFileURLInput {
    fileKey: string;
    expiresIn?: number;
}
