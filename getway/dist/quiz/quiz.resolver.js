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
exports.QuizResolver = void 0;
const graphql_1 = require("@nestjs/graphql");
const quiz_service_1 = require("./quiz.service");
const quiz_types_1 = require("./quiz.types");
let QuizResolver = class QuizResolver {
    quizService;
    constructor(quizService) {
        this.quizService = quizService;
    }
    async getQuizzesBySection(sectionId) {
        return await this.quizService.getQuizzesBySection(sectionId);
    }
    async getQuizById(id) {
        return await this.quizService.getQuizById(id);
    }
    async getUserQuizHistory(userId) {
        return await this.quizService.getUserQuizHistory(userId);
    }
    async getQuizStatistics(quizId) {
        return await this.quizService.getQuizStatistics(quizId);
    }
    async createQuiz(input, context) {
        const optionsMap = {};
        input.options.forEach(option => {
            optionsMap[option.key] = option.value;
        });
        const quizData = {
            section_id: input.sectionId,
            question: input.question,
            options: optionsMap,
            correct_answer: input.correctAnswer,
            explanation: input.explanation,
            difficulty: input.difficulty,
        };
        return await this.quizService.createQuiz(quizData);
    }
    async updateQuiz(id, input, context) {
        const updateData = {
            section_id: input.sectionId,
            question: input.question,
            correct_answer: input.correctAnswer,
            explanation: input.explanation,
            difficulty: input.difficulty,
        };
        if (input.options) {
            const optionsMap = {};
            input.options.forEach(option => {
                optionsMap[option.key] = option.value;
            });
            updateData.options = optionsMap;
        }
        return await this.quizService.updateQuiz(id, updateData);
    }
    async deleteQuiz(id, context) {
        return await this.quizService.deleteQuiz(id);
    }
    async submitQuizAnswer(input, context) {
        return await this.quizService.submitQuizAnswer(input.userId, input.quizId, input.selectedOption);
    }
};
exports.QuizResolver = QuizResolver;
__decorate([
    (0, graphql_1.Query)(() => quiz_types_1.QuizServiceResponse),
    __param(0, (0, graphql_1.Args)('sectionId')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], QuizResolver.prototype, "getQuizzesBySection", null);
__decorate([
    (0, graphql_1.Query)(() => quiz_types_1.QuizServiceResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], QuizResolver.prototype, "getQuizById", null);
__decorate([
    (0, graphql_1.Query)(() => quiz_types_1.QuizServiceResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], QuizResolver.prototype, "getUserQuizHistory", null);
__decorate([
    (0, graphql_1.Query)(() => quiz_types_1.QuizServiceResponse),
    __param(0, (0, graphql_1.Args)('quizId')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], QuizResolver.prototype, "getQuizStatistics", null);
__decorate([
    (0, graphql_1.Mutation)(() => quiz_types_1.QuizServiceResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [quiz_types_1.CreateQuizInput, Object]),
    __metadata("design:returntype", Promise)
], QuizResolver.prototype, "createQuiz", null);
__decorate([
    (0, graphql_1.Mutation)(() => quiz_types_1.QuizServiceResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __param(1, (0, graphql_1.Args)('input')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, quiz_types_1.UpdateQuizInput, Object]),
    __metadata("design:returntype", Promise)
], QuizResolver.prototype, "updateQuiz", null);
__decorate([
    (0, graphql_1.Mutation)(() => quiz_types_1.QuizServiceResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], QuizResolver.prototype, "deleteQuiz", null);
__decorate([
    (0, graphql_1.Mutation)(() => quiz_types_1.QuizServiceResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [quiz_types_1.SubmitQuizAnswerInput, Object]),
    __metadata("design:returntype", Promise)
], QuizResolver.prototype, "submitQuizAnswer", null);
exports.QuizResolver = QuizResolver = __decorate([
    (0, graphql_1.Resolver)(() => quiz_types_1.Quiz),
    __metadata("design:paramtypes", [quiz_service_1.QuizService])
], QuizResolver);
//# sourceMappingURL=quiz.resolver.js.map