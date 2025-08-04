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
exports.TransactionResolver = void 0;
const graphql_1 = require("@nestjs/graphql");
const transaction_service_1 = require("./transaction.service");
const transaction_types_1 = require("./transaction.types");
let TransactionResolver = class TransactionResolver {
    transactionService;
    constructor(transactionService) {
        this.transactionService = transactionService;
    }
    async getTransactionByID(id, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.transactionService.getTransactionByID(id, token);
    }
    async getTransactionsByUserID(userID, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.transactionService.getTransactionByUserID(userID, token);
    }
    async getTransactionsByCardID(cardID, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.transactionService.getTransactionByCardID(cardID, token);
    }
    async getTransactionsByStatus(status, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.transactionService.getTransactionByStatus(status, token);
    }
    async getTransactionsByDate(date, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.transactionService.getTransactionByDate(date, token);
    }
};
exports.TransactionResolver = TransactionResolver;
__decorate([
    (0, graphql_1.Query)(() => transaction_types_1.TransactionResponse),
    __param(0, (0, graphql_1.Args)('id')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], TransactionResolver.prototype, "getTransactionByID", null);
__decorate([
    (0, graphql_1.Query)(() => transaction_types_1.TransactionsResponse),
    __param(0, (0, graphql_1.Args)('userID')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], TransactionResolver.prototype, "getTransactionsByUserID", null);
__decorate([
    (0, graphql_1.Query)(() => transaction_types_1.TransactionsResponse),
    __param(0, (0, graphql_1.Args)('cardID')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], TransactionResolver.prototype, "getTransactionsByCardID", null);
__decorate([
    (0, graphql_1.Query)(() => transaction_types_1.TransactionsResponse),
    __param(0, (0, graphql_1.Args)('status')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], TransactionResolver.prototype, "getTransactionsByStatus", null);
__decorate([
    (0, graphql_1.Query)(() => transaction_types_1.TransactionsResponse),
    __param(0, (0, graphql_1.Args)('date')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], TransactionResolver.prototype, "getTransactionsByDate", null);
exports.TransactionResolver = TransactionResolver = __decorate([
    (0, graphql_1.Resolver)(() => transaction_types_1.Transaction),
    __metadata("design:paramtypes", [transaction_service_1.TransactionService])
], TransactionResolver);
//# sourceMappingURL=transaction.resolver.js.map