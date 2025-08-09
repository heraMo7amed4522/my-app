"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.QuizService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let QuizService = class QuizService {
    quizServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/quiz.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const quizProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.QUIZ_SERVICE_URL || 'localhost:50056';
        this.quizServiceClient = new quizProto.quiz.QuizService(serviceUrl, grpc.credentials.createInsecure());
    }
    async getQuizzesBySection(sectionId) {
        return new Promise((resolve, reject) => {
            this.quizServiceClient.GetQuizzesBySection({ section_id: sectionId }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getQuizById(id) {
        return new Promise((resolve, reject) => {
            this.quizServiceClient.GetQuizByID({ id }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async createQuiz(quizData) {
        return new Promise((resolve, reject) => {
            this.quizServiceClient.CreateQuiz({ quiz: quizData }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updateQuiz(id, quizData) {
        return new Promise((resolve, reject) => {
            this.quizServiceClient.UpdateQuiz({ quiz: { id, ...quizData } }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async deleteQuiz(id) {
        return new Promise((resolve, reject) => {
            this.quizServiceClient.DeleteQuiz({ id }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async submitQuizAnswer(userId, quizId, selectedOption) {
        return new Promise((resolve, reject) => {
            this.quizServiceClient.SubmitQuizAnswer({ user_id: userId, quiz_id: quizId, selected_option: selectedOption }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getUserQuizHistory(userId) {
        return new Promise((resolve, reject) => {
            this.quizServiceClient.GetUserQuizHistory({ user_id: userId }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getQuizStatistics(quizId) {
        return new Promise((resolve, reject) => {
            this.quizServiceClient.GetQuizStatistics({ quiz_id: quizId }, (error, response) => {
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
exports.QuizService = QuizService;
exports.QuizService = QuizService = __decorate([
    (0, common_1.Injectable)()
], QuizService);
//# sourceMappingURL=quiz.service.js.map