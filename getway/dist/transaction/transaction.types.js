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
exports.CreateTransactionInput = exports.TransactionsResponse = exports.TransactionResponse = exports.Transaction = void 0;
const graphql_1 = require("@nestjs/graphql");
let Transaction = class Transaction {
    id;
    userID;
    cardID;
    merchantID;
    merchantName;
    cardNumber;
    merchantCategory;
    amount;
    currency;
    status;
    createAt;
    updateAt;
};
exports.Transaction = Transaction;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "userID", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "cardID", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "merchantID", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "merchantName", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "cardNumber", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "merchantCategory", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "amount", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "currency", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "status", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "createAt", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Transaction.prototype, "updateAt", void 0);
exports.Transaction = Transaction = __decorate([
    (0, graphql_1.ObjectType)()
], Transaction);
let TransactionResponse = class TransactionResponse {
    statusCode;
    message;
    transaction;
};
exports.TransactionResponse = TransactionResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TransactionResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TransactionResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => Transaction, { nullable: true }),
    __metadata("design:type", Transaction)
], TransactionResponse.prototype, "transaction", void 0);
exports.TransactionResponse = TransactionResponse = __decorate([
    (0, graphql_1.ObjectType)()
], TransactionResponse);
let TransactionsResponse = class TransactionsResponse {
    statusCode;
    message;
    transactions;
};
exports.TransactionsResponse = TransactionsResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TransactionsResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TransactionsResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => [Transaction], { nullable: true }),
    __metadata("design:type", Array)
], TransactionsResponse.prototype, "transactions", void 0);
exports.TransactionsResponse = TransactionsResponse = __decorate([
    (0, graphql_1.ObjectType)()
], TransactionsResponse);
let CreateTransactionInput = class CreateTransactionInput {
    userID;
    cardID;
    merchantID;
    merchantName;
    cardNumber;
    merchantCategory;
    amount;
    currency;
};
exports.CreateTransactionInput = CreateTransactionInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTransactionInput.prototype, "userID", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTransactionInput.prototype, "cardID", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTransactionInput.prototype, "merchantID", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTransactionInput.prototype, "merchantName", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTransactionInput.prototype, "cardNumber", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTransactionInput.prototype, "merchantCategory", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTransactionInput.prototype, "amount", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateTransactionInput.prototype, "currency", void 0);
exports.CreateTransactionInput = CreateTransactionInput = __decorate([
    (0, graphql_1.InputType)()
], CreateTransactionInput);
//# sourceMappingURL=transaction.types.js.map