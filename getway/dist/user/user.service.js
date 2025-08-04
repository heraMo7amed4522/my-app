"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.UserService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let UserService = class UserService {
    userServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/user.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const userProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.USER_SERVICE_URL || 'localhost:50051';
        this.userServiceClient = new userProto.users.UserService(serviceUrl, grpc.credentials.createInsecure());
    }
    async getUserByEmail(email) {
        return new Promise((resolve, reject) => {
            this.userServiceClient.GetUserByEmail({ email }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async createNewUser(userData) {
        return new Promise((resolve, reject) => {
            this.userServiceClient.CreateNewUser(userData, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async loginUser(email, password) {
        return new Promise((resolve, reject) => {
            this.userServiceClient.LoginUser({ email, password }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    if (response.tokens && response.tokens.accessToken) {
                        response.accessToken = response.tokens.accessToken;
                    }
                    resolve(response);
                }
            });
        });
    }
    async forgetPassword(email, type) {
        return new Promise((resolve, reject) => {
            this.userServiceClient.ForgetPassword({ email, type }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async verifyCode(email, verifyCode) {
        return new Promise((resolve, reject) => {
            this.userServiceClient.VerifyCode({ email, verifyCode }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async resetPassword(email, password) {
        return new Promise((resolve, reject) => {
            this.userServiceClient.ResetPassword({ email, password }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updateUserData(userData, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.userServiceClient.UpdateUserData(userData, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async deleteUserData(userId, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.userServiceClient.DeleteUserData({ id: userId }, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async validateToken(token) {
        return new Promise((resolve, reject) => {
            this.userServiceClient.ValidateToken({ token }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
};
exports.UserService = UserService;
exports.UserService = UserService = __decorate([
    (0, common_1.Injectable)()
], UserService);
//# sourceMappingURL=user.service.js.map