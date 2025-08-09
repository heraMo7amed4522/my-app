"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.FeedbackService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let FeedbackService = class FeedbackService {
    feedbackServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/feedback.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const feedbackProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.FEEDBACK_SERVICE_URL || 'localhost:50059';
        this.feedbackServiceClient = new feedbackProto.feedback.FeedbackService(serviceUrl, grpc.credentials.createInsecure());
    }
    async submitFeedback(feedbackData) {
        return new Promise((resolve, reject) => {
            this.feedbackServiceClient.SubmitFeedback(feedbackData, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getTemplateFeedback(templateId, page, limit) {
        return new Promise((resolve, reject) => {
            this.feedbackServiceClient.GetTemplateFeedback({ template_id: templateId, page, limit }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getFeedbackStatistics(templateId) {
        return new Promise((resolve, reject) => {
            this.feedbackServiceClient.GetFeedbackStatistics({ template_id: templateId }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updateFeedback(id, feedbackData) {
        return new Promise((resolve, reject) => {
            this.feedbackServiceClient.UpdateFeedback({ id, ...feedbackData }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async deleteFeedback(id) {
        return new Promise((resolve, reject) => {
            this.feedbackServiceClient.DeleteFeedback({ id }, (error, response) => {
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
exports.FeedbackService = FeedbackService;
exports.FeedbackService = FeedbackService = __decorate([
    (0, common_1.Injectable)()
], FeedbackService);
//# sourceMappingURL=feedback.service.js.map