"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.HistoryTemplateService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let HistoryTemplateService = class HistoryTemplateService {
    historyTemplateServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/history_template.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const historyTemplateProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.HISTORY_TEMPLATE_SERVICE_URL || 'localhost:50055';
        this.historyTemplateServiceClient = new historyTemplateProto.history_templates.HistoryTemplateService(serviceUrl, grpc.credentials.createInsecure());
    }
    async getTemplateById(id) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.GetTemplateByID({ id }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getAllTemplates(page, limit, sortBy, order) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.GetAllTemplates({ page, limit, sortBy, order }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getTemplatesByEra(era, page, limit) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.GetTemplatesByEra({ era, page, limit }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getTemplatesByDynasty(dynasty, page, limit) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.GetTemplatesByDynasty({ dynasty, page, limit }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getTemplatesByPharaoh(pharaohId, page, limit) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.GetTemplatesByPharaoh({ pharaoh_id: pharaohId, page, limit }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getTemplatesByDifficulty(difficulty, page, limit) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.GetTemplatesByDifficulty({ difficulty, page, limit }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async searchTemplates(query, fields, page, limit) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.SearchTemplates({ query, fields, page, limit }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async createTemplate(templateData) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.CreateTemplate({ template: templateData }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updateTemplate(id, templateData) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.UpdateTemplate({ id, template: templateData }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async deleteTemplate(id) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.DeleteTemplate({ id }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getTemplatesByTag(tag, page, limit) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.GetTemplatesByTag({ tag, page, limit }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getRelatedTemplates(templateId, limit) {
        return new Promise((resolve, reject) => {
            this.historyTemplateServiceClient.GetRelatedTemplates({ template_id: templateId, limit }, (error, response) => {
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
exports.HistoryTemplateService = HistoryTemplateService;
exports.HistoryTemplateService = HistoryTemplateService = __decorate([
    (0, common_1.Injectable)()
], HistoryTemplateService);
//# sourceMappingURL=history-template.service.js.map