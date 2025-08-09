"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.UserProgressService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let UserProgressService = class UserProgressService {
    userProgressServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/user_progress.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const userProgressProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.USER_PROGRESS_SERVICE_URL || 'localhost:50060';
        this.userProgressServiceClient = new userProgressProto.user_progress.UserProgressService(serviceUrl, grpc.credentials.createInsecure());
    }
    async getUserProgress(userId, templateId) {
        return new Promise((resolve, reject) => {
            this.userProgressServiceClient.GetUserProgress({ user_id: userId, template_id: templateId || '' }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updateProgress(progressData) {
        return new Promise((resolve, reject) => {
            this.userProgressServiceClient.UpdateProgress(progressData, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getCompletedTemplates(userId, page, limit) {
        return new Promise((resolve, reject) => {
            this.userProgressServiceClient.GetCompletedTemplates({ user_id: userId, page, limit }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getUserAchievements(userId) {
        return new Promise((resolve, reject) => {
            this.userProgressServiceClient.GetUserAchievements({ user_id: userId }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async unlockAchievement(userId, achievementId) {
        return new Promise((resolve, reject) => {
            this.userProgressServiceClient.UnlockAchievement({ user_id: userId, achievement_id: achievementId }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getLearningPath(userId) {
        return new Promise((resolve, reject) => {
            this.userProgressServiceClient.GetLearningPath({ user_id: userId }, (error, response) => {
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
exports.UserProgressService = UserProgressService;
exports.UserProgressService = UserProgressService = __decorate([
    (0, common_1.Injectable)()
], UserProgressService);
//# sourceMappingURL=user-progress.service.js.map