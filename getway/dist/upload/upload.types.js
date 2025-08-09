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
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetFileURLInput = exports.DeleteFileInput = exports.UploadFileInput = exports.GetFileURLResponse = exports.DeleteFileResponse = exports.UploadFileResponse = exports.URLSuccess = exports.DeleteSuccess = exports.UploadSuccess = exports.FileType = void 0;
const graphql_1 = require("@nestjs/graphql");
const shared_types_1 = require("../shared/shared.types");
var FileType;
(function (FileType) {
    FileType["UNKNOWN"] = "UNKNOWN";
    FileType["IMAGE"] = "IMAGE";
    FileType["VIDEO"] = "VIDEO";
})(FileType || (exports.FileType = FileType = {}));
(0, graphql_1.registerEnumType)(FileType, {
    name: 'FileType',
});
let UploadSuccess = class UploadSuccess {
    fileUrl;
    fileKey;
    fileId;
    fileSize;
    uploadedAt;
};
exports.UploadSuccess = UploadSuccess;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UploadSuccess.prototype, "fileUrl", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UploadSuccess.prototype, "fileKey", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UploadSuccess.prototype, "fileId", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], UploadSuccess.prototype, "fileSize", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UploadSuccess.prototype, "uploadedAt", void 0);
exports.UploadSuccess = UploadSuccess = __decorate([
    (0, graphql_1.ObjectType)()
], UploadSuccess);
let DeleteSuccess = class DeleteSuccess {
    fileKey;
    deletedAt;
};
exports.DeleteSuccess = DeleteSuccess;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeleteSuccess.prototype, "fileKey", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeleteSuccess.prototype, "deletedAt", void 0);
exports.DeleteSuccess = DeleteSuccess = __decorate([
    (0, graphql_1.ObjectType)()
], DeleteSuccess);
let URLSuccess = class URLSuccess {
    signedUrl;
    expiresAt;
};
exports.URLSuccess = URLSuccess;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], URLSuccess.prototype, "signedUrl", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], URLSuccess.prototype, "expiresAt", void 0);
exports.URLSuccess = URLSuccess = __decorate([
    (0, graphql_1.ObjectType)()
], URLSuccess);
let UploadFileResponse = class UploadFileResponse {
    statusCode;
    message;
    success;
    error;
};
exports.UploadFileResponse = UploadFileResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], UploadFileResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UploadFileResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => UploadSuccess, { nullable: true }),
    __metadata("design:type", UploadSuccess)
], UploadFileResponse.prototype, "success", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], UploadFileResponse.prototype, "error", void 0);
exports.UploadFileResponse = UploadFileResponse = __decorate([
    (0, graphql_1.ObjectType)()
], UploadFileResponse);
let DeleteFileResponse = class DeleteFileResponse {
    statusCode;
    message;
    success;
    error;
};
exports.DeleteFileResponse = DeleteFileResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], DeleteFileResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeleteFileResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => DeleteSuccess, { nullable: true }),
    __metadata("design:type", DeleteSuccess)
], DeleteFileResponse.prototype, "success", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], DeleteFileResponse.prototype, "error", void 0);
exports.DeleteFileResponse = DeleteFileResponse = __decorate([
    (0, graphql_1.ObjectType)()
], DeleteFileResponse);
let GetFileURLResponse = class GetFileURLResponse {
    statusCode;
    message;
    success;
    error;
};
exports.GetFileURLResponse = GetFileURLResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], GetFileURLResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetFileURLResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => URLSuccess, { nullable: true }),
    __metadata("design:type", URLSuccess)
], GetFileURLResponse.prototype, "success", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], GetFileURLResponse.prototype, "error", void 0);
exports.GetFileURLResponse = GetFileURLResponse = __decorate([
    (0, graphql_1.ObjectType)()
], GetFileURLResponse);
let UploadFileInput = class UploadFileInput {
    file;
    fileName;
    contentType;
    userId;
    fileType;
};
exports.UploadFileInput = UploadFileInput;
__decorate([
    (0, graphql_1.Field)(() => String),
    __metadata("design:type", String)
], UploadFileInput.prototype, "file", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UploadFileInput.prototype, "fileName", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UploadFileInput.prototype, "contentType", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UploadFileInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(() => FileType),
    __metadata("design:type", String)
], UploadFileInput.prototype, "fileType", void 0);
exports.UploadFileInput = UploadFileInput = __decorate([
    (0, graphql_1.InputType)()
], UploadFileInput);
let DeleteFileInput = class DeleteFileInput {
    fileKey;
    userId;
};
exports.DeleteFileInput = DeleteFileInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeleteFileInput.prototype, "fileKey", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeleteFileInput.prototype, "userId", void 0);
exports.DeleteFileInput = DeleteFileInput = __decorate([
    (0, graphql_1.InputType)()
], DeleteFileInput);
let GetFileURLInput = class GetFileURLInput {
    fileKey;
    expiresIn;
};
exports.GetFileURLInput = GetFileURLInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetFileURLInput.prototype, "fileKey", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true, defaultValue: 3600 }),
    __metadata("design:type", Number)
], GetFileURLInput.prototype, "expiresIn", void 0);
exports.GetFileURLInput = GetFileURLInput = __decorate([
    (0, graphql_1.InputType)()
], GetFileURLInput);
//# sourceMappingURL=upload.types.js.map