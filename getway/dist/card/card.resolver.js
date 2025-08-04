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
exports.CardResolver = void 0;
const graphql_1 = require("@nestjs/graphql");
const card_service_1 = require("./card.service");
const card_types_1 = require("./card.types");
let CardResolver = class CardResolver {
    cardService;
    constructor(cardService) {
        this.cardService = cardService;
    }
    async getCardById(id, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.cardService.getCardById(id, token);
    }
    async createCard(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.cardService.createCard(input, token);
    }
    async updateCard(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.cardService.updateCard(input, token);
    }
    async deleteCard(id, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.cardService.deleteCard(id, token);
    }
    async updateCardStatus(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.cardService.updateCardStatus(input, token);
    }
    async getCardsByUser(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.cardService.getCardsByUser(input, token);
    }
    async searchCards(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.cardService.searchCards(input, token);
    }
    async findCardsByStatus(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.cardService.findCardsByStatus(input, token);
    }
    async processStripePayment(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.cardService.processStripePayment(input, token);
    }
};
exports.CardResolver = CardResolver;
__decorate([
    (0, graphql_1.Query)(() => card_types_1.CardResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], CardResolver.prototype, "getCardById", null);
__decorate([
    (0, graphql_1.Mutation)(() => card_types_1.CardResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [card_types_1.CreateCardInput, Object]),
    __metadata("design:returntype", Promise)
], CardResolver.prototype, "createCard", null);
__decorate([
    (0, graphql_1.Mutation)(() => card_types_1.CardResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [card_types_1.UpdateCardInput, Object]),
    __metadata("design:returntype", Promise)
], CardResolver.prototype, "updateCard", null);
__decorate([
    (0, graphql_1.Mutation)(() => card_types_1.CardResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], CardResolver.prototype, "deleteCard", null);
__decorate([
    (0, graphql_1.Mutation)(() => card_types_1.CardResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [card_types_1.UpdateCardStatusInput, Object]),
    __metadata("design:returntype", Promise)
], CardResolver.prototype, "updateCardStatus", null);
__decorate([
    (0, graphql_1.Query)(() => card_types_1.CardsResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [card_types_1.GetCardsByUserInput, Object]),
    __metadata("design:returntype", Promise)
], CardResolver.prototype, "getCardsByUser", null);
__decorate([
    (0, graphql_1.Query)(() => card_types_1.CardsResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [card_types_1.SearchCardInput, Object]),
    __metadata("design:returntype", Promise)
], CardResolver.prototype, "searchCards", null);
__decorate([
    (0, graphql_1.Query)(() => card_types_1.CardsResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [card_types_1.FindCardsByStatusInput, Object]),
    __metadata("design:returntype", Promise)
], CardResolver.prototype, "findCardsByStatus", null);
__decorate([
    (0, graphql_1.Mutation)(() => card_types_1.PaymentResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [card_types_1.StripePaymentInput, Object]),
    __metadata("design:returntype", Promise)
], CardResolver.prototype, "processStripePayment", null);
exports.CardResolver = CardResolver = __decorate([
    (0, graphql_1.Resolver)(() => card_types_1.Card),
    __metadata("design:paramtypes", [card_service_1.CardService])
], CardResolver);
//# sourceMappingURL=card.resolver.js.map