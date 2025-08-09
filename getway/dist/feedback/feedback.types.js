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
exports.UpdateFeedbackInput = exports.SubmitFeedbackInput = exports.FeedbackResponse = exports.RatingDistribution = exports.FeedbackStatistics = exports.FeedbackList = exports.Feedback = void 0;
const graphql_1 = require("@nestjs/graphql");
let Feedback = class Feedback {
    id;
    userId;
    templateId;
    rating;
    comment;
    language;
    createdAt;
};
exports.Feedback = Feedback;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Feedback.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Feedback.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Feedback.prototype, "templateId", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float),
    __metadata("design:type", Number)
], Feedback.prototype, "rating", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Feedback.prototype, "comment", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Feedback.prototype, "language", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Feedback.prototype, "createdAt", void 0);
exports.Feedback = Feedback = __decorate([
    (0, graphql_1.ObjectType)()
], Feedback);
let FeedbackList = class FeedbackList {
    feedback;
    totalCount;
    page;
    limit;
};
exports.FeedbackList = FeedbackList;
__decorate([
    (0, graphql_1.Field)(() => [Feedback]),
    __metadata("design:type", Array)
], FeedbackList.prototype, "feedback", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], FeedbackList.prototype, "totalCount", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], FeedbackList.prototype, "page", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], FeedbackList.prototype, "limit", void 0);
exports.FeedbackList = FeedbackList = __decorate([
    (0, graphql_1.ObjectType)()
], FeedbackList);
let FeedbackStatistics = class FeedbackStatistics {
    templateId;
    averageRating;
    totalRatings;
    ratingDistribution;
};
exports.FeedbackStatistics = FeedbackStatistics;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FeedbackStatistics.prototype, "templateId", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float),
    __metadata("design:type", Number)
], FeedbackStatistics.prototype, "averageRating", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], FeedbackStatistics.prototype, "totalRatings", void 0);
__decorate([
    (0, graphql_1.Field)(() => [RatingDistribution]),
    __metadata("design:type", Array)
], FeedbackStatistics.prototype, "ratingDistribution", void 0);
exports.FeedbackStatistics = FeedbackStatistics = __decorate([
    (0, graphql_1.ObjectType)()
], FeedbackStatistics);
let RatingDistribution = class RatingDistribution {
    rating;
    count;
};
exports.RatingDistribution = RatingDistribution;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], RatingDistribution.prototype, "rating", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], RatingDistribution.prototype, "count", void 0);
exports.RatingDistribution = RatingDistribution = __decorate([
    (0, graphql_1.ObjectType)()
], RatingDistribution);
let FeedbackResponse = class FeedbackResponse {
    statusCode;
    message;
    feedback;
    feedbackList;
    statistics;
    success;
};
exports.FeedbackResponse = FeedbackResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], FeedbackResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FeedbackResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => Feedback, { nullable: true }),
    __metadata("design:type", Feedback)
], FeedbackResponse.prototype, "feedback", void 0);
__decorate([
    (0, graphql_1.Field)(() => FeedbackList, { nullable: true }),
    __metadata("design:type", FeedbackList)
], FeedbackResponse.prototype, "feedbackList", void 0);
__decorate([
    (0, graphql_1.Field)(() => FeedbackStatistics, { nullable: true }),
    __metadata("design:type", FeedbackStatistics)
], FeedbackResponse.prototype, "statistics", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], FeedbackResponse.prototype, "success", void 0);
exports.FeedbackResponse = FeedbackResponse = __decorate([
    (0, graphql_1.ObjectType)()
], FeedbackResponse);
let SubmitFeedbackInput = class SubmitFeedbackInput {
    userId;
    templateId;
    rating;
    comment;
    language;
};
exports.SubmitFeedbackInput = SubmitFeedbackInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], SubmitFeedbackInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], SubmitFeedbackInput.prototype, "templateId", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float),
    __metadata("design:type", Number)
], SubmitFeedbackInput.prototype, "rating", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], SubmitFeedbackInput.prototype, "comment", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], SubmitFeedbackInput.prototype, "language", void 0);
exports.SubmitFeedbackInput = SubmitFeedbackInput = __decorate([
    (0, graphql_1.InputType)()
], SubmitFeedbackInput);
let UpdateFeedbackInput = class UpdateFeedbackInput {
    rating;
    comment;
};
exports.UpdateFeedbackInput = UpdateFeedbackInput;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float),
    __metadata("design:type", Number)
], UpdateFeedbackInput.prototype, "rating", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateFeedbackInput.prototype, "comment", void 0);
exports.UpdateFeedbackInput = UpdateFeedbackInput = __decorate([
    (0, graphql_1.InputType)()
], UpdateFeedbackInput);
//# sourceMappingURL=feedback.types.js.map