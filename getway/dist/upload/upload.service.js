"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.UploadService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let UploadService = class UploadService {
    uploadServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/upload.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const uploadProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.UPLOAD_SERVICE_URL || 'localhost:50061';
        this.uploadServiceClient = new uploadProto.upload.UploadService(serviceUrl, grpc.credentials.createInsecure());
    }
    async uploadFile(input) {
        try {
            const { file, fileName, contentType, userId, fileType } = input;
            if (!file) {
                throw new common_1.BadRequestException('File data is required');
            }
            if (!fileName) {
                throw new common_1.BadRequestException('File name is required');
            }
            if (!contentType) {
                throw new common_1.BadRequestException('Content type is required');
            }
            if (!userId) {
                throw new common_1.BadRequestException('User ID is required');
            }
            let fileBuffer;
            try {
                const base64Data = file.includes(',') ? file.split(',')[1] : file;
                fileBuffer = Buffer.from(base64Data, 'base64');
            }
            catch (error) {
                throw new common_1.BadRequestException('Invalid file data format');
            }
            return new Promise((resolve, reject) => {
                this.uploadServiceClient.UploadFile({
                    file_data: fileBuffer,
                    file_name: fileName,
                    content_type: contentType,
                    user_id: userId,
                    file_type: this.mapFileType(fileType),
                }, (error, response) => {
                    if (error) {
                        console.error('Upload error:', error);
                        reject(new common_1.BadRequestException(`Upload failed: ${error.message}`));
                    }
                    else {
                        resolve(response);
                    }
                });
            });
        }
        catch (error) {
            console.error('Upload service error:', error);
            if (error instanceof common_1.BadRequestException) {
                throw error;
            }
            throw new common_1.BadRequestException('Upload failed');
        }
    }
    async deleteFile(input) {
        try {
            const { fileKey, userId } = input;
            if (!fileKey) {
                throw new common_1.BadRequestException('File key is required');
            }
            if (!userId) {
                throw new common_1.BadRequestException('User ID is required');
            }
            return new Promise((resolve, reject) => {
                this.uploadServiceClient.DeleteFile({
                    file_key: fileKey,
                    user_id: userId,
                }, (error, response) => {
                    if (error) {
                        console.error('Delete error:', error);
                        reject(new common_1.BadRequestException(`Delete failed: ${error.message}`));
                    }
                    else {
                        resolve(response);
                    }
                });
            });
        }
        catch (error) {
            console.error('Delete service error:', error);
            if (error instanceof common_1.BadRequestException) {
                throw error;
            }
            throw new common_1.BadRequestException('Delete failed');
        }
    }
    async getFileURL(input) {
        try {
            const { fileKey, expiresIn } = input;
            if (!fileKey) {
                throw new common_1.BadRequestException('File key is required');
            }
            return new Promise((resolve, reject) => {
                this.uploadServiceClient.GetFileURL({
                    file_key: fileKey,
                    expires_in: expiresIn || 3600,
                }, (error, response) => {
                    if (error) {
                        console.error('Get URL error:', error);
                        reject(new common_1.BadRequestException(`Get URL failed: ${error.message}`));
                    }
                    else {
                        resolve(response);
                    }
                });
            });
        }
        catch (error) {
            console.error('Get URL service error:', error);
            if (error instanceof common_1.BadRequestException) {
                throw error;
            }
            throw new common_1.BadRequestException('Get URL failed');
        }
    }
    mapFileType(fileType) {
        switch (fileType) {
            case 'IMAGE':
                return 1;
            case 'VIDEO':
                return 2;
            default:
                return 0;
        }
    }
};
exports.UploadService = UploadService;
exports.UploadService = UploadService = __decorate([
    (0, common_1.Injectable)()
], UploadService);
//# sourceMappingURL=upload.service.js.map