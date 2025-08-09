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
exports.UnlockAchievementInput = exports.UpdateProgressInput = exports.UserProgressResponse = exports.UserAchievementList = exports.UserProgressList = exports.LearningPath = exports.UserAchievement = exports.Achievement = exports.UserProgress = void 0;
const graphql_1 = require("@nestjs/graphql");
let UserProgress = class UserProgress {
    id;
    userId;
    templateId;
    sectionId;
    progress;
    completed;
    lastViewed;
    createdAt;
};
exports.UserProgress = UserProgress;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserProgress.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserProgress.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserProgress.prototype, "templateId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserProgress.prototype, "sectionId", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float),
    __metadata("design:type", Number)
], UserProgress.prototype, "progress", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], UserProgress.prototype, "completed", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserProgress.prototype, "lastViewed", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserProgress.prototype, "createdAt", void 0);
exports.UserProgress = UserProgress = __decorate([
    (0, graphql_1.ObjectType)()
], UserProgress);
let Achievement = class Achievement {
    id;
    code;
    name;
    description;
    badgeUrl;
    createdAt;
};
exports.Achievement = Achievement;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Achievement.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Achievement.prototype, "code", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Achievement.prototype, "name", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Achievement.prototype, "description", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Achievement.prototype, "badgeUrl", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Achievement.prototype, "createdAt", void 0);
exports.Achievement = Achievement = __decorate([
    (0, graphql_1.ObjectType)()
], Achievement);
let UserAchievement = class UserAchievement {
    id;
    userId;
    achievementId;
    achievement;
    unlockedAt;
};
exports.UserAchievement = UserAchievement;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserAchievement.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserAchievement.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserAchievement.prototype, "achievementId", void 0);
__decorate([
    (0, graphql_1.Field)(() => Achievement),
    __metadata("design:type", Achievement)
], UserAchievement.prototype, "achievement", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserAchievement.prototype, "unlockedAt", void 0);
exports.UserAchievement = UserAchievement = __decorate([
    (0, graphql_1.ObjectType)()
], UserAchievement);
let LearningPath = class LearningPath {
    userId;
    recommendedTemplateIds;
    overallProgress;
    completedTemplates;
    totalTemplates;
};
exports.LearningPath = LearningPath;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], LearningPath.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String]),
    __metadata("design:type", Array)
], LearningPath.prototype, "recommendedTemplateIds", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float),
    __metadata("design:type", Number)
], LearningPath.prototype, "overallProgress", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], LearningPath.prototype, "completedTemplates", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], LearningPath.prototype, "totalTemplates", void 0);
exports.LearningPath = LearningPath = __decorate([
    (0, graphql_1.ObjectType)()
], LearningPath);
let UserProgressList = class UserProgressList {
    progress;
    totalCount;
};
exports.UserProgressList = UserProgressList;
__decorate([
    (0, graphql_1.Field)(() => [UserProgress]),
    __metadata("design:type", Array)
], UserProgressList.prototype, "progress", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], UserProgressList.prototype, "totalCount", void 0);
exports.UserProgressList = UserProgressList = __decorate([
    (0, graphql_1.ObjectType)()
], UserProgressList);
let UserAchievementList = class UserAchievementList {
    achievements;
    totalCount;
};
exports.UserAchievementList = UserAchievementList;
__decorate([
    (0, graphql_1.Field)(() => [UserAchievement]),
    __metadata("design:type", Array)
], UserAchievementList.prototype, "achievements", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], UserAchievementList.prototype, "totalCount", void 0);
exports.UserAchievementList = UserAchievementList = __decorate([
    (0, graphql_1.ObjectType)()
], UserAchievementList);
let UserProgressResponse = class UserProgressResponse {
    statusCode;
    message;
    progress;
    progressList;
    achievement;
    achievements;
    learningPath;
    success;
};
exports.UserProgressResponse = UserProgressResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], UserProgressResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserProgressResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => UserProgress, { nullable: true }),
    __metadata("design:type", UserProgress)
], UserProgressResponse.prototype, "progress", void 0);
__decorate([
    (0, graphql_1.Field)(() => UserProgressList, { nullable: true }),
    __metadata("design:type", UserProgressList)
], UserProgressResponse.prototype, "progressList", void 0);
__decorate([
    (0, graphql_1.Field)(() => UserAchievement, { nullable: true }),
    __metadata("design:type", UserAchievement)
], UserProgressResponse.prototype, "achievement", void 0);
__decorate([
    (0, graphql_1.Field)(() => UserAchievementList, { nullable: true }),
    __metadata("design:type", UserAchievementList)
], UserProgressResponse.prototype, "achievements", void 0);
__decorate([
    (0, graphql_1.Field)(() => LearningPath, { nullable: true }),
    __metadata("design:type", LearningPath)
], UserProgressResponse.prototype, "learningPath", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], UserProgressResponse.prototype, "success", void 0);
exports.UserProgressResponse = UserProgressResponse = __decorate([
    (0, graphql_1.ObjectType)()
], UserProgressResponse);
let UpdateProgressInput = class UpdateProgressInput {
    userId;
    templateId;
    sectionId;
    progress;
    completed;
};
exports.UpdateProgressInput = UpdateProgressInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateProgressInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateProgressInput.prototype, "templateId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateProgressInput.prototype, "sectionId", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float),
    __metadata("design:type", Number)
], UpdateProgressInput.prototype, "progress", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], UpdateProgressInput.prototype, "completed", void 0);
exports.UpdateProgressInput = UpdateProgressInput = __decorate([
    (0, graphql_1.InputType)()
], UpdateProgressInput);
let UnlockAchievementInput = class UnlockAchievementInput {
    userId;
    achievementId;
};
exports.UnlockAchievementInput = UnlockAchievementInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UnlockAchievementInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UnlockAchievementInput.prototype, "achievementId", void 0);
exports.UnlockAchievementInput = UnlockAchievementInput = __decorate([
    (0, graphql_1.InputType)()
], UnlockAchievementInput);
//# sourceMappingURL=user-progress.types.js.map