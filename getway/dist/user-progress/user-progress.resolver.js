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
exports.UserProgressResolver = void 0;
const graphql_1 = require("@nestjs/graphql");
const user_progress_service_1 = require("./user-progress.service");
const user_progress_types_1 = require("./user-progress.types");
let UserProgressResolver = class UserProgressResolver {
    userProgressService;
    constructor(userProgressService) {
        this.userProgressService = userProgressService;
    }
    async getUserProgress(userId, templateId) {
        try {
            const response = await this.userProgressService.getUserProgress(userId, templateId);
            if (response.result?.progress) {
                const progressList = {
                    progress: response.result.progress.progress.map((p) => ({
                        id: p.id,
                        userId: p.user_id,
                        templateId: p.template_id,
                        sectionId: p.section_id,
                        progress: p.progress,
                        completed: p.completed,
                        lastViewed: p.last_viewed?.seconds
                            ? new Date(p.last_viewed.seconds * 1000).toISOString()
                            : new Date().toISOString(),
                        createdAt: p.created_at?.seconds
                            ? new Date(p.created_at.seconds * 1000).toISOString()
                            : new Date().toISOString(),
                    })),
                    totalCount: response.result.progress.total_count,
                };
                return {
                    statusCode: response.status_code,
                    message: response.message,
                    progressList,
                };
            }
            else {
                return {
                    statusCode: response.status_code,
                    message: response.message,
                };
            }
        }
        catch (error) {
            return {
                statusCode: 500,
                message: 'Failed to retrieve user progress',
                success: false,
            };
        }
    }
    async updateProgress(input, context) {
        try {
            const response = await this.userProgressService.updateProgress({
                user_id: input.userId,
                template_id: input.templateId,
                section_id: input.sectionId,
                progress: input.progress,
                completed: input.completed,
            });
            if (response.result?.progress) {
                return {
                    statusCode: response.status_code,
                    message: response.message,
                    progress: {
                        id: response.result.progress.id,
                        userId: response.result.progress.user_id,
                        templateId: response.result.progress.template_id,
                        sectionId: response.result.progress.section_id,
                        progress: response.result.progress.progress,
                        completed: response.result.progress.completed,
                        lastViewed: response.result.progress.last_viewed?.seconds
                            ? new Date(response.result.progress.last_viewed.seconds * 1000).toISOString()
                            : new Date().toISOString(),
                        createdAt: response.result.progress.created_at?.seconds
                            ? new Date(response.result.progress.created_at.seconds * 1000).toISOString()
                            : new Date().toISOString(),
                    },
                };
            }
            else {
                return {
                    statusCode: response.status_code,
                    message: response.message,
                };
            }
        }
        catch (error) {
            return {
                statusCode: 500,
                message: 'Failed to update progress',
                success: false,
            };
        }
    }
    async getCompletedTemplates(userId, page, limit) {
        try {
            const response = await this.userProgressService.getCompletedTemplates(userId, page, limit);
            if (response.result?.completed) {
                const progressList = {
                    progress: response.result.completed.progress.map((p) => ({
                        id: p.id,
                        userId: p.user_id,
                        templateId: p.template_id,
                        sectionId: p.section_id,
                        progress: p.progress,
                        completed: p.completed,
                        lastViewed: p.last_viewed?.seconds
                            ? new Date(p.last_viewed.seconds * 1000).toISOString()
                            : new Date().toISOString(),
                        createdAt: p.created_at?.seconds
                            ? new Date(p.created_at.seconds * 1000).toISOString()
                            : new Date().toISOString(),
                    })),
                    totalCount: response.result.completed.total_count,
                };
                return {
                    statusCode: response.status_code,
                    message: response.message,
                    progressList,
                };
            }
            else {
                return {
                    statusCode: response.status_code,
                    message: response.message,
                };
            }
        }
        catch (error) {
            return {
                statusCode: 500,
                message: 'Failed to retrieve completed templates',
                success: false,
            };
        }
    }
    async getUserAchievements(userId) {
        try {
            const response = await this.userProgressService.getUserAchievements(userId);
            if (response.result?.achievements) {
                const achievements = {
                    achievements: response.result.achievements.achievements.map((a) => ({
                        id: a.id,
                        userId: a.user_id,
                        achievementId: a.achievement_id,
                        achievement: {
                            id: a.achievement.id,
                            code: a.achievement.code,
                            name: a.achievement.name,
                            description: a.achievement.description,
                            badgeUrl: a.achievement.badge_url,
                            createdAt: a.achievement.created_at?.seconds
                                ? new Date(a.achievement.created_at.seconds * 1000).toISOString()
                                : new Date().toISOString(),
                        },
                        unlockedAt: a.unlocked_at?.seconds
                            ? new Date(a.unlocked_at.seconds * 1000).toISOString()
                            : new Date().toISOString(),
                    })),
                    totalCount: response.result.achievements.total_count,
                };
                return {
                    statusCode: response.status_code,
                    message: response.message,
                    achievements,
                };
            }
            else {
                return {
                    statusCode: response.status_code,
                    message: response.message,
                };
            }
        }
        catch (error) {
            return {
                statusCode: 500,
                message: 'Failed to retrieve user achievements',
                success: false,
            };
        }
    }
    async unlockAchievement(input, context) {
        try {
            const response = await this.userProgressService.unlockAchievement(input.userId, input.achievementId);
            if (response.result?.achievement) {
                return {
                    statusCode: response.status_code,
                    message: response.message,
                    achievement: {
                        id: response.result.achievement.id,
                        userId: response.result.achievement.user_id,
                        achievementId: response.result.achievement.achievement_id,
                        achievement: {
                            id: response.result.achievement.achievement.id,
                            code: response.result.achievement.achievement.code,
                            name: response.result.achievement.achievement.name,
                            description: response.result.achievement.achievement.description,
                            badgeUrl: response.result.achievement.achievement.badge_url,
                            createdAt: response.result.achievement.achievement.created_at?.seconds
                                ? new Date(response.result.achievement.achievement.created_at.seconds * 1000).toISOString()
                                : new Date().toISOString(),
                        },
                        unlockedAt: response.result.achievement.unlocked_at?.seconds
                            ? new Date(response.result.achievement.unlocked_at.seconds * 1000).toISOString()
                            : new Date().toISOString(),
                    },
                };
            }
            else {
                return {
                    statusCode: response.status_code,
                    message: response.message,
                };
            }
        }
        catch (error) {
            return {
                statusCode: 500,
                message: 'Failed to unlock achievement',
                success: false,
            };
        }
    }
    async getLearningPath(userId) {
        try {
            const response = await this.userProgressService.getLearningPath(userId);
            if (response.result?.path) {
                const learningPath = {
                    userId: response.result.path.user_id,
                    recommendedTemplateIds: response.result.path.recommended_template_ids,
                    overallProgress: response.result.path.overall_progress,
                    completedTemplates: response.result.path.completed_templates,
                    totalTemplates: response.result.path.total_templates,
                };
                return {
                    statusCode: response.status_code,
                    message: response.message,
                    learningPath,
                };
            }
            else {
                return {
                    statusCode: response.status_code,
                    message: response.message,
                };
            }
        }
        catch (error) {
            return {
                statusCode: 500,
                message: 'Failed to generate learning path',
                success: false,
            };
        }
    }
};
exports.UserProgressResolver = UserProgressResolver;
__decorate([
    (0, graphql_1.Query)(() => user_progress_types_1.UserProgressResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __param(1, (0, graphql_1.Args)('templateId', { nullable: true })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String]),
    __metadata("design:returntype", Promise)
], UserProgressResolver.prototype, "getUserProgress", null);
__decorate([
    (0, graphql_1.Mutation)(() => user_progress_types_1.UserProgressResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [user_progress_types_1.UpdateProgressInput, Object]),
    __metadata("design:returntype", Promise)
], UserProgressResolver.prototype, "updateProgress", null);
__decorate([
    (0, graphql_1.Query)(() => user_progress_types_1.UserProgressResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __param(1, (0, graphql_1.Args)('page', { type: () => graphql_1.Int, defaultValue: 1 })),
    __param(2, (0, graphql_1.Args)('limit', { type: () => graphql_1.Int, defaultValue: 10 })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Number, Number]),
    __metadata("design:returntype", Promise)
], UserProgressResolver.prototype, "getCompletedTemplates", null);
__decorate([
    (0, graphql_1.Query)(() => user_progress_types_1.UserProgressResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], UserProgressResolver.prototype, "getUserAchievements", null);
__decorate([
    (0, graphql_1.Mutation)(() => user_progress_types_1.UserProgressResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [user_progress_types_1.UnlockAchievementInput, Object]),
    __metadata("design:returntype", Promise)
], UserProgressResolver.prototype, "unlockAchievement", null);
__decorate([
    (0, graphql_1.Query)(() => user_progress_types_1.UserProgressResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], UserProgressResolver.prototype, "getLearningPath", null);
exports.UserProgressResolver = UserProgressResolver = __decorate([
    (0, graphql_1.Resolver)(() => user_progress_types_1.UserProgress),
    __metadata("design:paramtypes", [user_progress_service_1.UserProgressService])
], UserProgressResolver);
//# sourceMappingURL=user-progress.resolver.js.map