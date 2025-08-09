"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.UploadResolver = void 0;
const graphql_1 = require("@nestjs/graphql");
const upload_service_1 = require("./upload.service");
const upload_types_1 = require("./upload.types");
let UploadResolver = class UploadResolver {
    uploadService;
    constructor(uploadService) {
        this.uploadService = uploadService;
    }
    async uploadFile(input) {
        try {
            return await this.uploadService.uploadFile(input);
        }
        catch (error) {
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
    async deleteFile(input) {
        try {
            return await this.uploadService.deleteFile(input);
        }
        catch (error) {
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
    async getFileURL(input) {
        try {
            return await this.uploadService.getFileURL(input);
        }
        catch (error) {
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
};
exports.UploadResolver = UploadResolver;
__decorate([
    (0, graphql_1.Mutation)(() => upload_types_1.UploadFileResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [upload_types_1.UploadFileInput]),
    __metadata("design:returntype", Promise)
], UploadResolver.prototype, "uploadFile", null);
__decorate([
    (0, graphql_1.Mutation)(() => upload_types_1.DeleteFileResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [upload_types_1.DeleteFileInput]),
    __metadata("design:returntype", Promise)
], UploadResolver.prototype, "deleteFile", null);
__decorate([
    (0, graphql_1.Query)(() => upload_types_1.GetFileURLResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [upload_types_1.GetFileURLInput]),
    __metadata("design:returntype", Promise)
], UploadResolver.prototype, "getFileURL", null);
exports.UploadResolver = UploadResolver = __decorate([
    (0, graphql_1.Resolver)(),
    __metadata("design:paramtypes", [upload_service_1.UploadService])
], UploadResolver);
//# sourceMappingURL=upload.resolver.js.map