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
exports.WalletResolver = void 0;
const graphql_1 = require("@nestjs/graphql");
const wallet_service_1 = require("./wallet.service");
const wallet_types_1 = require("./wallet.types");
let WalletResolver = class WalletResolver {
    walletService;
    constructor(walletService) {
        this.walletService = walletService;
    }
    async getWalletByUserID(userID, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.walletService.getWalletByUserID(userID, token);
    }
    async createWallet(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.walletService.createWallet(input, token);
    }
    async fundWallet(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.walletService.fundWallet(input, token);
    }
    async deductFromWallet(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.walletService.deductFromWallet(input, token);
    }
    async getWalletBalance(userID, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.walletService.getWalletBalance(userID, token);
    }
    async getWalletTransactions(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.walletService.getWalletTransactions(input, token);
    }
    async refundToWallet(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.walletService.refundToWallet(input, token);
    }
};
exports.WalletResolver = WalletResolver;
__decorate([
    (0, graphql_1.Query)(() => wallet_types_1.WalletResponse),
    __param(0, (0, graphql_1.Args)('userID')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], WalletResolver.prototype, "getWalletByUserID", null);
__decorate([
    (0, graphql_1.Mutation)(() => wallet_types_1.WalletResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [wallet_types_1.CreateWalletInput, Object]),
    __metadata("design:returntype", Promise)
], WalletResolver.prototype, "createWallet", null);
__decorate([
    (0, graphql_1.Mutation)(() => wallet_types_1.WalletResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [wallet_types_1.FundWalletInput, Object]),
    __metadata("design:returntype", Promise)
], WalletResolver.prototype, "fundWallet", null);
__decorate([
    (0, graphql_1.Mutation)(() => wallet_types_1.WalletResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [wallet_types_1.DeductFromWalletInput, Object]),
    __metadata("design:returntype", Promise)
], WalletResolver.prototype, "deductFromWallet", null);
__decorate([
    (0, graphql_1.Query)(() => wallet_types_1.WalletBalanceResponse),
    __param(0, (0, graphql_1.Args)('userID')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], WalletResolver.prototype, "getWalletBalance", null);
__decorate([
    (0, graphql_1.Query)(() => wallet_types_1.WalletTransactionsResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [wallet_types_1.GetWalletTransactionsInput, Object]),
    __metadata("design:returntype", Promise)
], WalletResolver.prototype, "getWalletTransactions", null);
__decorate([
    (0, graphql_1.Mutation)(() => wallet_types_1.WalletResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [wallet_types_1.RefundToWalletInput, Object]),
    __metadata("design:returntype", Promise)
], WalletResolver.prototype, "refundToWallet", null);
exports.WalletResolver = WalletResolver = __decorate([
    (0, graphql_1.Resolver)(() => wallet_types_1.Wallet),
    __metadata("design:paramtypes", [wallet_service_1.WalletService])
], WalletResolver);
//# sourceMappingURL=wallet.resolver.js.map