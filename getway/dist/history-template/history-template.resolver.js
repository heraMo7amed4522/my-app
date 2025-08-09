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
exports.HistoryTemplateResolver = void 0;
const graphql_1 = require("@nestjs/graphql");
const history_template_service_1 = require("./history-template.service");
const history_template_types_1 = require("./history-template.types");
let HistoryTemplateResolver = class HistoryTemplateResolver {
    historyTemplateService;
    constructor(historyTemplateService) {
        this.historyTemplateService = historyTemplateService;
    }
    async getTemplateById(id) {
        return await this.historyTemplateService.getTemplateById(id);
    }
    async getAllTemplates(page, limit, sortBy, order) {
        return await this.historyTemplateService.getAllTemplates(page, limit, sortBy, order);
    }
    async getTemplatesByEra(era, page, limit) {
        return await this.historyTemplateService.getTemplatesByEra(era, page, limit);
    }
    async getTemplatesByDynasty(dynasty, page, limit) {
        return await this.historyTemplateService.getTemplatesByDynasty(dynasty, page, limit);
    }
    async getTemplatesByPharaoh(pharaohId, page, limit) {
        return await this.historyTemplateService.getTemplatesByPharaoh(pharaohId, page, limit);
    }
    async getTemplatesByDifficulty(difficulty, page, limit) {
        return await this.historyTemplateService.getTemplatesByDifficulty(difficulty, page, limit);
    }
    async searchTemplates(query, fields, page, limit) {
        return await this.historyTemplateService.searchTemplates(query, fields || [], page, limit);
    }
    async getTemplatesByTag(tag, page, limit) {
        return await this.historyTemplateService.getTemplatesByTag(tag, page, limit);
    }
    async getRelatedTemplates(templateId, limit) {
        return await this.historyTemplateService.getRelatedTemplates(templateId, limit);
    }
    async createTemplate(input, context) {
        return await this.historyTemplateService.createTemplate(input);
    }
    async updateTemplate(id, input, context) {
        return await this.historyTemplateService.updateTemplate(id, input);
    }
    async deleteTemplate(id, context) {
        return await this.historyTemplateService.deleteTemplate(id);
    }
};
exports.HistoryTemplateResolver = HistoryTemplateResolver;
__decorate([
    (0, graphql_1.Query)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "getTemplateById", null);
__decorate([
    (0, graphql_1.Query)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('page', { type: () => graphql_1.Int, defaultValue: 1 })),
    __param(1, (0, graphql_1.Args)('limit', { type: () => graphql_1.Int, defaultValue: 10 })),
    __param(2, (0, graphql_1.Args)('sortBy', { nullable: true })),
    __param(3, (0, graphql_1.Args)('order', { nullable: true })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Number, String, String]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "getAllTemplates", null);
__decorate([
    (0, graphql_1.Query)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('era')),
    __param(1, (0, graphql_1.Args)('page', { type: () => graphql_1.Int, defaultValue: 1 })),
    __param(2, (0, graphql_1.Args)('limit', { type: () => graphql_1.Int, defaultValue: 10 })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Number, Number]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "getTemplatesByEra", null);
__decorate([
    (0, graphql_1.Query)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('dynasty', { type: () => graphql_1.Int })),
    __param(1, (0, graphql_1.Args)('page', { type: () => graphql_1.Int, defaultValue: 1 })),
    __param(2, (0, graphql_1.Args)('limit', { type: () => graphql_1.Int, defaultValue: 10 })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Number, Number, Number]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "getTemplatesByDynasty", null);
__decorate([
    (0, graphql_1.Query)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('pharaohId')),
    __param(1, (0, graphql_1.Args)('page', { type: () => graphql_1.Int, defaultValue: 1 })),
    __param(2, (0, graphql_1.Args)('limit', { type: () => graphql_1.Int, defaultValue: 10 })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Number, Number]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "getTemplatesByPharaoh", null);
__decorate([
    (0, graphql_1.Query)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('difficulty')),
    __param(1, (0, graphql_1.Args)('page', { type: () => graphql_1.Int, defaultValue: 1 })),
    __param(2, (0, graphql_1.Args)('limit', { type: () => graphql_1.Int, defaultValue: 10 })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Number, Number]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "getTemplatesByDifficulty", null);
__decorate([
    (0, graphql_1.Query)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('query')),
    __param(1, (0, graphql_1.Args)('fields', { type: () => [String], nullable: true })),
    __param(2, (0, graphql_1.Args)('page', { type: () => graphql_1.Int, defaultValue: 1 })),
    __param(3, (0, graphql_1.Args)('limit', { type: () => graphql_1.Int, defaultValue: 10 })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Array, Number, Number]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "searchTemplates", null);
__decorate([
    (0, graphql_1.Query)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('tag')),
    __param(1, (0, graphql_1.Args)('page', { type: () => graphql_1.Int, defaultValue: 1 })),
    __param(2, (0, graphql_1.Args)('limit', { type: () => graphql_1.Int, defaultValue: 10 })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Number, Number]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "getTemplatesByTag", null);
__decorate([
    (0, graphql_1.Query)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('templateId')),
    __param(1, (0, graphql_1.Args)('limit', { type: () => graphql_1.Int, defaultValue: 5 })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Number]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "getRelatedTemplates", null);
__decorate([
    (0, graphql_1.Mutation)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [history_template_types_1.CreateTemplateInput, Object]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "createTemplate", null);
__decorate([
    (0, graphql_1.Mutation)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __param(1, (0, graphql_1.Args)('input')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, history_template_types_1.UpdateTemplateInput, Object]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "updateTemplate", null);
__decorate([
    (0, graphql_1.Mutation)(() => history_template_types_1.HistoryTemplateResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], HistoryTemplateResolver.prototype, "deleteTemplate", null);
exports.HistoryTemplateResolver = HistoryTemplateResolver = __decorate([
    (0, graphql_1.Resolver)(() => history_template_types_1.HistoryTemplate),
    __metadata("design:paramtypes", [history_template_service_1.HistoryTemplateService])
], HistoryTemplateResolver);
//# sourceMappingURL=history-template.resolver.js.map