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
exports.SubmitQuizAnswerInput = exports.UpdateQuizInput = exports.QuizOptionInput = exports.CreateQuizInput = exports.QuizServiceResponse = exports.QuizStatistics = exports.QuizList = exports.QuizResponse = exports.QuizOption = exports.Quiz = void 0;
const graphql_1 = require("@nestjs/graphql");
let Quiz = class Quiz {
    id;
    sectionId;
    question;
    options;
    correctAnswer;
    explanation;
    difficulty;
    createdAt;
};
exports.Quiz = Quiz;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Quiz.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Quiz.prototype, "sectionId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Quiz.prototype, "question", void 0);
__decorate([
    (0, graphql_1.Field)(() => [QuizOption]),
    __metadata("design:type", Array)
], Quiz.prototype, "options", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Quiz.prototype, "correctAnswer", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Quiz.prototype, "explanation", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Quiz.prototype, "difficulty", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Quiz.prototype, "createdAt", void 0);
exports.Quiz = Quiz = __decorate([
    (0, graphql_1.ObjectType)()
], Quiz);
let QuizOption = class QuizOption {
    key;
    value;
};
exports.QuizOption = QuizOption;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], QuizOption.prototype, "key", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], QuizOption.prototype, "value", void 0);
exports.QuizOption = QuizOption = __decorate([
    (0, graphql_1.ObjectType)()
], QuizOption);
let QuizResponse = class QuizResponse {
    id;
    userId;
    quizId;
    selectedOption;
    isCorrect;
    answeredAt;
};
exports.QuizResponse = QuizResponse;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], QuizResponse.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], QuizResponse.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], QuizResponse.prototype, "quizId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], QuizResponse.prototype, "selectedOption", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], QuizResponse.prototype, "isCorrect", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], QuizResponse.prototype, "answeredAt", void 0);
exports.QuizResponse = QuizResponse = __decorate([
    (0, graphql_1.ObjectType)()
], QuizResponse);
let QuizList = class QuizList {
    quizzes;
    totalCount;
    page;
    limit;
};
exports.QuizList = QuizList;
__decorate([
    (0, graphql_1.Field)(() => [Quiz]),
    __metadata("design:type", Array)
], QuizList.prototype, "quizzes", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], QuizList.prototype, "totalCount", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], QuizList.prototype, "page", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], QuizList.prototype, "limit", void 0);
exports.QuizList = QuizList = __decorate([
    (0, graphql_1.ObjectType)()
], QuizList);
let QuizStatistics = class QuizStatistics {
    totalAttempts;
    correctAnswers;
    uniqueUsers;
    successRate;
};
exports.QuizStatistics = QuizStatistics;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], QuizStatistics.prototype, "totalAttempts", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], QuizStatistics.prototype, "correctAnswers", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], QuizStatistics.prototype, "uniqueUsers", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], QuizStatistics.prototype, "successRate", void 0);
exports.QuizStatistics = QuizStatistics = __decorate([
    (0, graphql_1.ObjectType)()
], QuizStatistics);
let QuizServiceResponse = class QuizServiceResponse {
    statusCode;
    message;
    quiz;
    quizzes;
    quizResponses;
    statistics;
    success;
};
exports.QuizServiceResponse = QuizServiceResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], QuizServiceResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], QuizServiceResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => Quiz, { nullable: true }),
    __metadata("design:type", Quiz)
], QuizServiceResponse.prototype, "quiz", void 0);
__decorate([
    (0, graphql_1.Field)(() => QuizList, { nullable: true }),
    __metadata("design:type", QuizList)
], QuizServiceResponse.prototype, "quizzes", void 0);
__decorate([
    (0, graphql_1.Field)(() => [QuizResponse], { nullable: true }),
    __metadata("design:type", Array)
], QuizServiceResponse.prototype, "quizResponses", void 0);
__decorate([
    (0, graphql_1.Field)(() => QuizStatistics, { nullable: true }),
    __metadata("design:type", QuizStatistics)
], QuizServiceResponse.prototype, "statistics", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], QuizServiceResponse.prototype, "success", void 0);
exports.QuizServiceResponse = QuizServiceResponse = __decorate([
    (0, graphql_1.ObjectType)()
], QuizServiceResponse);
let CreateQuizInput = class CreateQuizInput {
    sectionId;
    question;
    options;
    correctAnswer;
    explanation;
    difficulty;
};
exports.CreateQuizInput = CreateQuizInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateQuizInput.prototype, "sectionId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateQuizInput.prototype, "question", void 0);
__decorate([
    (0, graphql_1.Field)(() => [QuizOptionInput]),
    __metadata("design:type", Array)
], CreateQuizInput.prototype, "options", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateQuizInput.prototype, "correctAnswer", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreateQuizInput.prototype, "explanation", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateQuizInput.prototype, "difficulty", void 0);
exports.CreateQuizInput = CreateQuizInput = __decorate([
    (0, graphql_1.InputType)()
], CreateQuizInput);
let QuizOptionInput = class QuizOptionInput {
    key;
    value;
};
exports.QuizOptionInput = QuizOptionInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], QuizOptionInput.prototype, "key", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], QuizOptionInput.prototype, "value", void 0);
exports.QuizOptionInput = QuizOptionInput = __decorate([
    (0, graphql_1.InputType)()
], QuizOptionInput);
let UpdateQuizInput = class UpdateQuizInput {
    sectionId;
    question;
    options;
    correctAnswer;
    explanation;
    difficulty;
};
exports.UpdateQuizInput = UpdateQuizInput;
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateQuizInput.prototype, "sectionId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateQuizInput.prototype, "question", void 0);
__decorate([
    (0, graphql_1.Field)(() => [QuizOptionInput], { nullable: true }),
    __metadata("design:type", Array)
], UpdateQuizInput.prototype, "options", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateQuizInput.prototype, "correctAnswer", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateQuizInput.prototype, "explanation", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateQuizInput.prototype, "difficulty", void 0);
exports.UpdateQuizInput = UpdateQuizInput = __decorate([
    (0, graphql_1.InputType)()
], UpdateQuizInput);
let SubmitQuizAnswerInput = class SubmitQuizAnswerInput {
    userId;
    quizId;
    selectedOption;
};
exports.SubmitQuizAnswerInput = SubmitQuizAnswerInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], SubmitQuizAnswerInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], SubmitQuizAnswerInput.prototype, "quizId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], SubmitQuizAnswerInput.prototype, "selectedOption", void 0);
exports.SubmitQuizAnswerInput = SubmitQuizAnswerInput = __decorate([
    (0, graphql_1.InputType)()
], SubmitQuizAnswerInput);
//# sourceMappingURL=quiz.types.js.map