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
exports.RefundToWalletInput = exports.GetWalletTransactionsInput = exports.DeductFromWalletInput = exports.FundWalletInput = exports.CreateWalletInput = exports.WalletBalanceResponse = exports.WalletTransactionsResponse = exports.WalletResponse = exports.WalletTransaction = exports.Wallet = void 0;
const graphql_1 = require("@nestjs/graphql");
const shared_types_1 = require("../shared/shared.types");
let Wallet = class Wallet {
    id;
    userID;
    balance;
    currency;
    createdAt;
    updatedAt;
};
exports.Wallet = Wallet;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Wallet.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Wallet.prototype, "userID", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Wallet.prototype, "balance", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Wallet.prototype, "currency", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Wallet.prototype, "createdAt", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Wallet.prototype, "updatedAt", void 0);
exports.Wallet = Wallet = __decorate([
    (0, graphql_1.ObjectType)()
], Wallet);
let WalletTransaction = class WalletTransaction {
    id;
    walletId;
    type;
    amount;
    currency;
    description;
    createdAt;
};
exports.WalletTransaction = WalletTransaction;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], WalletTransaction.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], WalletTransaction.prototype, "walletId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], WalletTransaction.prototype, "type", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], WalletTransaction.prototype, "amount", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], WalletTransaction.prototype, "currency", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], WalletTransaction.prototype, "description", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], WalletTransaction.prototype, "createdAt", void 0);
exports.WalletTransaction = WalletTransaction = __decorate([
    (0, graphql_1.ObjectType)()
], WalletTransaction);
let WalletResponse = class WalletResponse {
    statusCode;
    message;
    wallet;
    error;
};
exports.WalletResponse = WalletResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], WalletResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], WalletResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => Wallet, { nullable: true }),
    __metadata("design:type", Wallet)
], WalletResponse.prototype, "wallet", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], WalletResponse.prototype, "error", void 0);
exports.WalletResponse = WalletResponse = __decorate([
    (0, graphql_1.ObjectType)()
], WalletResponse);
let WalletTransactionsResponse = class WalletTransactionsResponse {
    statusCode;
    message;
    transactions;
    error;
};
exports.WalletTransactionsResponse = WalletTransactionsResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], WalletTransactionsResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], WalletTransactionsResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => [WalletTransaction], { nullable: true }),
    __metadata("design:type", Array)
], WalletTransactionsResponse.prototype, "transactions", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], WalletTransactionsResponse.prototype, "error", void 0);
exports.WalletTransactionsResponse = WalletTransactionsResponse = __decorate([
    (0, graphql_1.ObjectType)()
], WalletTransactionsResponse);
let WalletBalanceResponse = class WalletBalanceResponse {
    statusCode;
    message;
    balance;
    currency;
    error;
};
exports.WalletBalanceResponse = WalletBalanceResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], WalletBalanceResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], WalletBalanceResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], WalletBalanceResponse.prototype, "balance", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], WalletBalanceResponse.prototype, "currency", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], WalletBalanceResponse.prototype, "error", void 0);
exports.WalletBalanceResponse = WalletBalanceResponse = __decorate([
    (0, graphql_1.ObjectType)()
], WalletBalanceResponse);
let CreateWalletInput = class CreateWalletInput {
    userID;
    currency;
};
exports.CreateWalletInput = CreateWalletInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateWalletInput.prototype, "userID", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreateWalletInput.prototype, "currency", void 0);
exports.CreateWalletInput = CreateWalletInput = __decorate([
    (0, graphql_1.InputType)()
], CreateWalletInput);
let FundWalletInput = class FundWalletInput {
    userID;
    cardID;
    amount;
    currency;
};
exports.FundWalletInput = FundWalletInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FundWalletInput.prototype, "userID", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FundWalletInput.prototype, "cardID", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FundWalletInput.prototype, "amount", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FundWalletInput.prototype, "currency", void 0);
exports.FundWalletInput = FundWalletInput = __decorate([
    (0, graphql_1.InputType)()
], FundWalletInput);
let DeductFromWalletInput = class DeductFromWalletInput {
    userID;
    amount;
    currency;
    transactionID;
};
exports.DeductFromWalletInput = DeductFromWalletInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeductFromWalletInput.prototype, "userID", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeductFromWalletInput.prototype, "amount", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeductFromWalletInput.prototype, "currency", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeductFromWalletInput.prototype, "transactionID", void 0);
exports.DeductFromWalletInput = DeductFromWalletInput = __decorate([
    (0, graphql_1.InputType)()
], DeductFromWalletInput);
let GetWalletTransactionsInput = class GetWalletTransactionsInput {
    userID;
    limit;
    offset;
};
exports.GetWalletTransactionsInput = GetWalletTransactionsInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetWalletTransactionsInput.prototype, "userID", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], GetWalletTransactionsInput.prototype, "limit", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], GetWalletTransactionsInput.prototype, "offset", void 0);
exports.GetWalletTransactionsInput = GetWalletTransactionsInput = __decorate([
    (0, graphql_1.InputType)()
], GetWalletTransactionsInput);
let RefundToWalletInput = class RefundToWalletInput {
    userID;
    amount;
    currency;
    transactionID;
    reason;
};
exports.RefundToWalletInput = RefundToWalletInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], RefundToWalletInput.prototype, "userID", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], RefundToWalletInput.prototype, "amount", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], RefundToWalletInput.prototype, "currency", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], RefundToWalletInput.prototype, "transactionID", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], RefundToWalletInput.prototype, "reason", void 0);
exports.RefundToWalletInput = RefundToWalletInput = __decorate([
    (0, graphql_1.InputType)()
], RefundToWalletInput);
//# sourceMappingURL=wallet.types.js.map