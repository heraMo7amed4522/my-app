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
exports.UpdateUserInput = exports.CreateUserInput = exports.TokenValidationResponse = exports.TokenClaims = exports.UserResponse = exports.User = void 0;
const graphql_1 = require("@nestjs/graphql");
const auth_types_1 = require("../auth/auth.types");
let User = class User {
    id;
    fullName;
    email;
    countryCode;
    phoneNumber;
    role;
    verifiyCode;
    createAt;
    updateAt;
};
exports.User = User;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], User.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], User.prototype, "fullName", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], User.prototype, "email", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], User.prototype, "countryCode", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], User.prototype, "phoneNumber", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], User.prototype, "role", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], User.prototype, "verifiyCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], User.prototype, "createAt", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], User.prototype, "updateAt", void 0);
exports.User = User = __decorate([
    (0, graphql_1.ObjectType)()
], User);
let UserResponse = class UserResponse {
    statusCode;
    message;
    user;
    tokens;
    accessToken;
    verifyCode;
    response;
};
exports.UserResponse = UserResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], UserResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => User, { nullable: true }),
    __metadata("design:type", User)
], UserResponse.prototype, "user", void 0);
__decorate([
    (0, graphql_1.Field)(() => auth_types_1.AuthTokens, { nullable: true }),
    __metadata("design:type", auth_types_1.AuthTokens)
], UserResponse.prototype, "tokens", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UserResponse.prototype, "accessToken", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UserResponse.prototype, "verifyCode", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UserResponse.prototype, "response", void 0);
exports.UserResponse = UserResponse = __decorate([
    (0, graphql_1.ObjectType)()
], UserResponse);
let TokenClaims = class TokenClaims {
    userId;
    email;
    role;
    exp;
};
exports.TokenClaims = TokenClaims;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TokenClaims.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TokenClaims.prototype, "email", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TokenClaims.prototype, "role", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TokenClaims.prototype, "exp", void 0);
exports.TokenClaims = TokenClaims = __decorate([
    (0, graphql_1.ObjectType)()
], TokenClaims);
let TokenValidationResponse = class TokenValidationResponse {
    statusCode;
    message;
    claims;
};
exports.TokenValidationResponse = TokenValidationResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TokenValidationResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], TokenValidationResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => TokenClaims, { nullable: true }),
    __metadata("design:type", TokenClaims)
], TokenValidationResponse.prototype, "claims", void 0);
exports.TokenValidationResponse = TokenValidationResponse = __decorate([
    (0, graphql_1.ObjectType)()
], TokenValidationResponse);
let CreateUserInput = class CreateUserInput {
    fullName;
    email;
    password;
    countryCode;
    phoneNumber;
    role;
    verifiyCode;
};
exports.CreateUserInput = CreateUserInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateUserInput.prototype, "fullName", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateUserInput.prototype, "email", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateUserInput.prototype, "password", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreateUserInput.prototype, "countryCode", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreateUserInput.prototype, "phoneNumber", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreateUserInput.prototype, "role", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreateUserInput.prototype, "verifiyCode", void 0);
exports.CreateUserInput = CreateUserInput = __decorate([
    (0, graphql_1.InputType)()
], CreateUserInput);
let UpdateUserInput = class UpdateUserInput {
    email;
    fullName;
    countryCode;
    phoneNumber;
    role;
};
exports.UpdateUserInput = UpdateUserInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateUserInput.prototype, "email", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateUserInput.prototype, "fullName", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateUserInput.prototype, "countryCode", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateUserInput.prototype, "phoneNumber", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateUserInput.prototype, "role", void 0);
exports.UpdateUserInput = UpdateUserInput = __decorate([
    (0, graphql_1.InputType)()
], UpdateUserInput);
//# sourceMappingURL=user.types.js.map