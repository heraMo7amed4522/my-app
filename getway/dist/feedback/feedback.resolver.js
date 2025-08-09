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
exports.FeedbackResolver = void 0;
const graphql_1 = require("@nestjs/graphql");
const feedback_service_1 = require("./feedback.service");
const feedback_types_1 = require("./feedback.types");
let FeedbackResolver = class FeedbackResolver {
    feedbackService;
    constructor(feedbackService) {
        this.feedbackService = feedbackService;
    }
    async submitFeedback(input, context) {
        try {
            const response = await this.feedbackService.submitFeedback({
                user_id: input.userId,
                template_id: input.templateId,
                rating: input.rating,
                comment: input.comment,
                language: input.language || 'en',
            });
            if (response.result?.feedback) {
                return {
                    statusCode: response.status_code,
                    message: response.message,
                    feedback: {
                        id: response.result.feedback.id,
                        userId: response.result.feedback.user_id,
                        templateId: response.result.feedback.template_id,
                        rating: response.result.feedback.rating,
                        comment: response.result.feedback.comment,
                        language: response.result.feedback.language,
                        createdAt: response.result.feedback.created_at?.seconds
                            ? new Date(response.result.feedback.created_at.seconds * 1000).toISOString()
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
                message: 'Failed to submit feedback',
            };
        }
    }
    async getTemplateFeedback(templateId, page, limit) {
        try {
            const response = await this.feedbackService.getTemplateFeedback(templateId, page, limit);
            if (response.result?.feedback) {
                const feedbackList = {
                    feedback: response.result.feedback.feedback.map((f) => ({
                        id: f.id,
                        userId: f.user_id,
                        templateId: f.template_id,
                        rating: f.rating,
                        comment: f.comment,
                        language: f.language,
                        createdAt: f.created_at?.seconds
                            ? new Date(f.created_at.seconds * 1000).toISOString()
                            : new Date().toISOString(),
                    })),
                    totalCount: response.result.feedback.total_count,
                    page: response.result.feedback.page,
                    limit: response.result.feedback.limit,
                };
                return {
                    statusCode: response.status_code,
                    message: response.message,
                    feedbackList,
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
                message: 'Failed to get template feedback',
            };
        }
    }
    async getFeedbackStatistics(templateId) {
        try {
            const response = await this.feedbackService.getFeedbackStatistics(templateId);
            if (response.result?.statistics) {
                const statistics = {
                    templateId: response.result.statistics.template_id,
                    averageRating: response.result.statistics.average_rating,
                    totalRatings: response.result.statistics.total_ratings,
                    ratingDistribution: Object.entries(response.result.statistics.rating_distribution || {}).map(([rating, count]) => ({
                        rating,
                        count: count,
                    })),
                };
                return {
                    statusCode: response.status_code,
                    message: response.message,
                    statistics,
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
                message: 'Failed to get feedback statistics',
            };
        }
    }
    async updateFeedback(id, input, context) {
        try {
            const response = await this.feedbackService.updateFeedback(id, {
                rating: input.rating,
                comment: input.comment,
            });
            if (response.result?.feedback) {
                return {
                    statusCode: response.status_code,
                    message: response.message,
                    feedback: {
                        id: response.result.feedback.id,
                        userId: response.result.feedback.user_id,
                        templateId: response.result.feedback.template_id,
                        rating: response.result.feedback.rating,
                        comment: response.result.feedback.comment,
                        language: response.result.feedback.language,
                        createdAt: response.result.feedback.created_at?.seconds
                            ? new Date(response.result.feedback.created_at.seconds * 1000).toISOString()
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
                message: 'Failed to update feedback',
            };
        }
    }
    async deleteFeedback(id, context) {
        try {
            const response = await this.feedbackService.deleteFeedback(id);
            return {
                statusCode: response.status_code,
                message: response.message,
                success: response.result?.success || false,
            };
        }
        catch (error) {
            return {
                statusCode: 500,
                message: 'Failed to delete feedback',
                success: false,
            };
        }
    }
};
exports.FeedbackResolver = FeedbackResolver;
__decorate([
    (0, graphql_1.Mutation)(() => feedback_types_1.FeedbackResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [feedback_types_1.SubmitFeedbackInput, Object]),
    __metadata("design:returntype", Promise)
], FeedbackResolver.prototype, "submitFeedback", null);
__decorate([
    (0, graphql_1.Query)(() => feedback_types_1.FeedbackResponse),
    __param(0, (0, graphql_1.Args)('templateId')),
    __param(1, (0, graphql_1.Args)('page', { type: () => graphql_1.Int, defaultValue: 1 })),
    __param(2, (0, graphql_1.Args)('limit', { type: () => graphql_1.Int, defaultValue: 10 })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Number, Number]),
    __metadata("design:returntype", Promise)
], FeedbackResolver.prototype, "getTemplateFeedback", null);
__decorate([
    (0, graphql_1.Query)(() => feedback_types_1.FeedbackResponse),
    __param(0, (0, graphql_1.Args)('templateId')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], FeedbackResolver.prototype, "getFeedbackStatistics", null);
__decorate([
    (0, graphql_1.Mutation)(() => feedback_types_1.FeedbackResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __param(1, (0, graphql_1.Args)('input')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, feedback_types_1.UpdateFeedbackInput, Object]),
    __metadata("design:returntype", Promise)
], FeedbackResolver.prototype, "updateFeedback", null);
__decorate([
    (0, graphql_1.Mutation)(() => feedback_types_1.FeedbackResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], FeedbackResolver.prototype, "deleteFeedback", null);
exports.FeedbackResolver = FeedbackResolver = __decorate([
    (0, graphql_1.Resolver)(() => feedback_types_1.Feedback),
    __metadata("design:paramtypes", [feedback_service_1.FeedbackService])
], FeedbackResolver);
//# sourceMappingURL=feedback.resolver.js.map