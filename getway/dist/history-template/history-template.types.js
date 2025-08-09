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
exports.UpdateTemplateInput = exports.CreateTemplateInput = exports.HistoryTemplateResponse = exports.TemplateList = exports.HistoryTemplate = exports.TemplateSection = void 0;
const graphql_1 = require("@nestjs/graphql");
let TemplateSection = class TemplateSection {
    id;
    templateId;
    title;
    subtitle;
    contentType;
    content;
    metadata;
    orderIndex;
    optional;
    createdAt;
};
exports.TemplateSection = TemplateSection;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TemplateSection.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TemplateSection.prototype, "templateId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TemplateSection.prototype, "title", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], TemplateSection.prototype, "subtitle", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TemplateSection.prototype, "contentType", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TemplateSection.prototype, "content", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], TemplateSection.prototype, "metadata", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TemplateSection.prototype, "orderIndex", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], TemplateSection.prototype, "optional", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TemplateSection.prototype, "createdAt", void 0);
exports.TemplateSection = TemplateSection = __decorate([
    (0, graphql_1.ObjectType)()
], TemplateSection);
let HistoryTemplate = class HistoryTemplate {
    id;
    title;
    description;
    era;
    dynasty;
    pharaohId;
    difficulty;
    thumbnailUrl;
    language;
    isActive;
    version;
    publishedAt;
    createdAt;
    updatedAt;
    sections;
    tags;
};
exports.HistoryTemplate = HistoryTemplate;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], HistoryTemplate.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], HistoryTemplate.prototype, "title", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], HistoryTemplate.prototype, "description", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], HistoryTemplate.prototype, "era", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], HistoryTemplate.prototype, "dynasty", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], HistoryTemplate.prototype, "pharaohId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], HistoryTemplate.prototype, "difficulty", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], HistoryTemplate.prototype, "thumbnailUrl", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], HistoryTemplate.prototype, "language", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], HistoryTemplate.prototype, "isActive", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], HistoryTemplate.prototype, "version", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], HistoryTemplate.prototype, "publishedAt", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], HistoryTemplate.prototype, "createdAt", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], HistoryTemplate.prototype, "updatedAt", void 0);
__decorate([
    (0, graphql_1.Field)(() => [TemplateSection], { nullable: true }),
    __metadata("design:type", Array)
], HistoryTemplate.prototype, "sections", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], HistoryTemplate.prototype, "tags", void 0);
exports.HistoryTemplate = HistoryTemplate = __decorate([
    (0, graphql_1.ObjectType)()
], HistoryTemplate);
let TemplateList = class TemplateList {
    templates;
    totalCount;
    page;
    limit;
};
exports.TemplateList = TemplateList;
__decorate([
    (0, graphql_1.Field)(() => [HistoryTemplate]),
    __metadata("design:type", Array)
], TemplateList.prototype, "templates", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TemplateList.prototype, "totalCount", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TemplateList.prototype, "page", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TemplateList.prototype, "limit", void 0);
exports.TemplateList = TemplateList = __decorate([
    (0, graphql_1.ObjectType)()
], TemplateList);
let HistoryTemplateResponse = class HistoryTemplateResponse {
    statusCode;
    message;
    template;
    templates;
    success;
};
exports.HistoryTemplateResponse = HistoryTemplateResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], HistoryTemplateResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], HistoryTemplateResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => HistoryTemplate, { nullable: true }),
    __metadata("design:type", HistoryTemplate)
], HistoryTemplateResponse.prototype, "template", void 0);
__decorate([
    (0, graphql_1.Field)(() => TemplateList, { nullable: true }),
    __metadata("design:type", TemplateList)
], HistoryTemplateResponse.prototype, "templates", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], HistoryTemplateResponse.prototype, "success", void 0);
exports.HistoryTemplateResponse = HistoryTemplateResponse = __decorate([
    (0, graphql_1.ObjectType)()
], HistoryTemplateResponse);
let CreateTemplateInput = class CreateTemplateInput {
    title;
    description;
    era;
    dynasty;
    pharaohId;
    difficulty;
    thumbnailUrl;
    language;
    isActive;
    tags;
};
exports.CreateTemplateInput = CreateTemplateInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTemplateInput.prototype, "title", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTemplateInput.prototype, "description", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTemplateInput.prototype, "era", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], CreateTemplateInput.prototype, "dynasty", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTemplateInput.prototype, "pharaohId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTemplateInput.prototype, "difficulty", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreateTemplateInput.prototype, "thumbnailUrl", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTemplateInput.prototype, "language", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], CreateTemplateInput.prototype, "isActive", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], CreateTemplateInput.prototype, "tags", void 0);
exports.CreateTemplateInput = CreateTemplateInput = __decorate([
    (0, graphql_1.InputType)()
], CreateTemplateInput);
let UpdateTemplateInput = class UpdateTemplateInput {
    title;
    description;
    era;
    dynasty;
    pharaohId;
    difficulty;
    thumbnailUrl;
    language;
    isActive;
    tags;
};
exports.UpdateTemplateInput = UpdateTemplateInput;
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateTemplateInput.prototype, "title", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateTemplateInput.prototype, "description", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateTemplateInput.prototype, "era", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], UpdateTemplateInput.prototype, "dynasty", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateTemplateInput.prototype, "pharaohId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateTemplateInput.prototype, "difficulty", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateTemplateInput.prototype, "thumbnailUrl", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateTemplateInput.prototype, "language", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], UpdateTemplateInput.prototype, "isActive", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], UpdateTemplateInput.prototype, "tags", void 0);
exports.UpdateTemplateInput = UpdateTemplateInput = __decorate([
    (0, graphql_1.InputType)()
], UpdateTemplateInput);
//# sourceMappingURL=history-template.types.js.map