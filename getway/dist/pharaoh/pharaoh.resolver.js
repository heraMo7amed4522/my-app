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
exports.PharaohResolver = void 0;
const graphql_1 = require("@nestjs/graphql");
const pharaoh_service_1 = require("./pharaoh.service");
const pharaoh_types_1 = require("./pharaoh.types");
let PharaohResolver = class PharaohResolver {
    pharaohService;
    constructor(pharaohService) {
        this.pharaohService = pharaohService;
    }
    async getPharaohById(id) {
        return await this.pharaohService.getPharaohById(id);
    }
    async getAllPharaohs(input = {}) {
        return await this.pharaohService.getAllPharaohs(input);
    }
    async getPharaohsByDynasty(input) {
        return await this.pharaohService.getPharaohsByDynasty(input);
    }
    async getPharaohsByPeriod(input) {
        return await this.pharaohService.getPharaohsByPeriod(input);
    }
    async searchPharaohs(input) {
        return await this.pharaohService.searchPharaohs(input);
    }
    async getPharaohsByRarity(input) {
        return await this.pharaohService.getPharaohsByRarity(input);
    }
    async createPharaoh(input) {
        return await this.pharaohService.createPharaoh(input);
    }
    async updatePharaoh(input) {
        return await this.pharaohService.updatePharaoh(input);
    }
    async deletePharaoh(id) {
        return await this.pharaohService.deletePharaoh(id);
    }
    async updatePharaohRating(input) {
        return await this.pharaohService.updatePharaohRating(input);
    }
};
exports.PharaohResolver = PharaohResolver;
__decorate([
    (0, graphql_1.Query)(() => pharaoh_types_1.PharaohResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], PharaohResolver.prototype, "getPharaohById", null);
__decorate([
    (0, graphql_1.Query)(() => pharaoh_types_1.PharaohsResponse),
    __param(0, (0, graphql_1.Args)('input', { nullable: true })),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [pharaoh_types_1.GetAllPharaohsInput]),
    __metadata("design:returntype", Promise)
], PharaohResolver.prototype, "getAllPharaohs", null);
__decorate([
    (0, graphql_1.Query)(() => pharaoh_types_1.PharaohsResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [pharaoh_types_1.GetPharaohsByDynastyInput]),
    __metadata("design:returntype", Promise)
], PharaohResolver.prototype, "getPharaohsByDynasty", null);
__decorate([
    (0, graphql_1.Query)(() => pharaoh_types_1.PharaohsResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [pharaoh_types_1.GetPharaohsByPeriodInput]),
    __metadata("design:returntype", Promise)
], PharaohResolver.prototype, "getPharaohsByPeriod", null);
__decorate([
    (0, graphql_1.Query)(() => pharaoh_types_1.PharaohsResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [pharaoh_types_1.SearchPharaohsInput]),
    __metadata("design:returntype", Promise)
], PharaohResolver.prototype, "searchPharaohs", null);
__decorate([
    (0, graphql_1.Query)(() => pharaoh_types_1.PharaohsResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [pharaoh_types_1.GetPharaohsByRarityInput]),
    __metadata("design:returntype", Promise)
], PharaohResolver.prototype, "getPharaohsByRarity", null);
__decorate([
    (0, graphql_1.Mutation)(() => pharaoh_types_1.PharaohResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [pharaoh_types_1.CreatePharaohInput]),
    __metadata("design:returntype", Promise)
], PharaohResolver.prototype, "createPharaoh", null);
__decorate([
    (0, graphql_1.Mutation)(() => pharaoh_types_1.PharaohResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [pharaoh_types_1.UpdatePharaohInput]),
    __metadata("design:returntype", Promise)
], PharaohResolver.prototype, "updatePharaoh", null);
__decorate([
    (0, graphql_1.Mutation)(() => pharaoh_types_1.DeletePharaohResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String]),
    __metadata("design:returntype", Promise)
], PharaohResolver.prototype, "deletePharaoh", null);
__decorate([
    (0, graphql_1.Mutation)(() => pharaoh_types_1.PharaohResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [pharaoh_types_1.UpdatePharaohRatingInput]),
    __metadata("design:returntype", Promise)
], PharaohResolver.prototype, "updatePharaohRating", null);
exports.PharaohResolver = PharaohResolver = __decorate([
    (0, graphql_1.Resolver)(() => pharaoh_types_1.Pharaoh),
    __metadata("design:paramtypes", [pharaoh_service_1.PharaohService])
], PharaohResolver);
//# sourceMappingURL=pharaoh.resolver.js.map