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
exports.FindCardsByStatusInput = exports.GetCardsByUserInput = exports.StripePaymentInput = exports.SearchCardInput = exports.UpdateCardStatusInput = exports.UpdateCardInput = exports.CreateCardInput = exports.PaymentResponse = exports.CardsResponse = exports.CardResponse = exports.Card = void 0;
const graphql_1 = require("@nestjs/graphql");
const shared_types_1 = require("../shared/shared.types");
let Card = class Card {
    id;
    userId;
    cardNumber;
    cardHolderName;
    expirationDate;
    cvv;
    status;
    cardType;
    balance;
    createAt;
    updateAt;
};
exports.Card = Card;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Card.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Card.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Card.prototype, "cardNumber", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Card.prototype, "cardHolderName", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Card.prototype, "expirationDate", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Card.prototype, "cvv", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Card.prototype, "status", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Card.prototype, "cardType", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Card.prototype, "balance", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Card.prototype, "createAt", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Card.prototype, "updateAt", void 0);
exports.Card = Card = __decorate([
    (0, graphql_1.ObjectType)()
], Card);
let CardResponse = class CardResponse {
    statusCode;
    message;
    card;
    error;
};
exports.CardResponse = CardResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], CardResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CardResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => Card, { nullable: true }),
    __metadata("design:type", Card)
], CardResponse.prototype, "card", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], CardResponse.prototype, "error", void 0);
exports.CardResponse = CardResponse = __decorate([
    (0, graphql_1.ObjectType)()
], CardResponse);
let CardsResponse = class CardsResponse {
    statusCode;
    message;
    cards;
    error;
    totalCount;
    currentPage;
    totalPages;
};
exports.CardsResponse = CardsResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], CardsResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CardsResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => [Card], { nullable: true }),
    __metadata("design:type", Array)
], CardsResponse.prototype, "cards", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], CardsResponse.prototype, "error", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], CardsResponse.prototype, "totalCount", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], CardsResponse.prototype, "currentPage", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], CardsResponse.prototype, "totalPages", void 0);
exports.CardsResponse = CardsResponse = __decorate([
    (0, graphql_1.ObjectType)()
], CardsResponse);
let PaymentResponse = class PaymentResponse {
    statusCode;
    message;
    paymentIntentId;
    clientSecret;
    error;
};
exports.PaymentResponse = PaymentResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], PaymentResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], PaymentResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], PaymentResponse.prototype, "paymentIntentId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], PaymentResponse.prototype, "clientSecret", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], PaymentResponse.prototype, "error", void 0);
exports.PaymentResponse = PaymentResponse = __decorate([
    (0, graphql_1.ObjectType)()
], PaymentResponse);
let CreateCardInput = class CreateCardInput {
    userId;
    cardNumber;
    cardHolderName;
    expirationDate;
    cvv;
    cardType;
    balance;
};
exports.CreateCardInput = CreateCardInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateCardInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateCardInput.prototype, "cardNumber", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateCardInput.prototype, "cardHolderName", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateCardInput.prototype, "expirationDate", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateCardInput.prototype, "cvv", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreateCardInput.prototype, "cardType", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreateCardInput.prototype, "balance", void 0);
exports.CreateCardInput = CreateCardInput = __decorate([
    (0, graphql_1.InputType)()
], CreateCardInput);
let UpdateCardInput = class UpdateCardInput {
    id;
    cardHolderName;
    expirationDate;
    cvv;
    cardType;
    balance;
};
exports.UpdateCardInput = UpdateCardInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateCardInput.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateCardInput.prototype, "cardHolderName", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateCardInput.prototype, "expirationDate", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateCardInput.prototype, "cvv", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateCardInput.prototype, "cardType", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateCardInput.prototype, "balance", void 0);
exports.UpdateCardInput = UpdateCardInput = __decorate([
    (0, graphql_1.InputType)()
], UpdateCardInput);
let UpdateCardStatusInput = class UpdateCardStatusInput {
    id;
    status;
};
exports.UpdateCardStatusInput = UpdateCardStatusInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateCardStatusInput.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateCardStatusInput.prototype, "status", void 0);
exports.UpdateCardStatusInput = UpdateCardStatusInput = __decorate([
    (0, graphql_1.InputType)()
], UpdateCardStatusInput);
let SearchCardInput = class SearchCardInput {
    userId;
    cardHolderName;
    cardType;
    status;
    page;
    limit;
};
exports.SearchCardInput = SearchCardInput;
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], SearchCardInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], SearchCardInput.prototype, "cardHolderName", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], SearchCardInput.prototype, "cardType", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], SearchCardInput.prototype, "status", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], SearchCardInput.prototype, "page", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], SearchCardInput.prototype, "limit", void 0);
exports.SearchCardInput = SearchCardInput = __decorate([
    (0, graphql_1.InputType)()
], SearchCardInput);
let StripePaymentInput = class StripePaymentInput {
    cardId;
    amount;
    currency;
    description;
};
exports.StripePaymentInput = StripePaymentInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], StripePaymentInput.prototype, "cardId", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], StripePaymentInput.prototype, "amount", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], StripePaymentInput.prototype, "currency", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], StripePaymentInput.prototype, "description", void 0);
exports.StripePaymentInput = StripePaymentInput = __decorate([
    (0, graphql_1.InputType)()
], StripePaymentInput);
let GetCardsByUserInput = class GetCardsByUserInput {
    userId;
    page;
    limit;
};
exports.GetCardsByUserInput = GetCardsByUserInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetCardsByUserInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], GetCardsByUserInput.prototype, "page", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], GetCardsByUserInput.prototype, "limit", void 0);
exports.GetCardsByUserInput = GetCardsByUserInput = __decorate([
    (0, graphql_1.InputType)()
], GetCardsByUserInput);
let FindCardsByStatusInput = class FindCardsByStatusInput {
    status;
    page;
    limit;
};
exports.FindCardsByStatusInput = FindCardsByStatusInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FindCardsByStatusInput.prototype, "status", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], FindCardsByStatusInput.prototype, "page", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], FindCardsByStatusInput.prototype, "limit", void 0);
exports.FindCardsByStatusInput = FindCardsByStatusInput = __decorate([
    (0, graphql_1.InputType)()
], FindCardsByStatusInput);
//# sourceMappingURL=card.types.js.map