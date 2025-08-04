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
exports.FirebaseAuthInput = exports.AuthResponse = exports.AuthTokens = exports.UserInfo = void 0;
const graphql_1 = require("@nestjs/graphql");
let UserInfo = class UserInfo {
    id;
    email;
    fullName;
    role;
};
exports.UserInfo = UserInfo;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserInfo.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserInfo.prototype, "email", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserInfo.prototype, "fullName", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserInfo.prototype, "role", void 0);
exports.UserInfo = UserInfo = __decorate([
    (0, graphql_1.ObjectType)()
], UserInfo);
let AuthTokens = class AuthTokens {
    accessToken;
    expiresIn;
    user;
};
exports.AuthTokens = AuthTokens;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], AuthTokens.prototype, "accessToken", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], AuthTokens.prototype, "expiresIn", void 0);
__decorate([
    (0, graphql_1.Field)(() => UserInfo),
    __metadata("design:type", UserInfo)
], AuthTokens.prototype, "user", void 0);
exports.AuthTokens = AuthTokens = __decorate([
    (0, graphql_1.ObjectType)()
], AuthTokens);
let AuthResponse = class AuthResponse {
    statusCode;
    message;
    tokens;
    error;
};
exports.AuthResponse = AuthResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], AuthResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], AuthResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => AuthTokens, { nullable: true }),
    __metadata("design:type", AuthTokens)
], AuthResponse.prototype, "tokens", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], AuthResponse.prototype, "error", void 0);
exports.AuthResponse = AuthResponse = __decorate([
    (0, graphql_1.ObjectType)()
], AuthResponse);
let FirebaseAuthInput = class FirebaseAuthInput {
    idToken;
    provider;
};
exports.FirebaseAuthInput = FirebaseAuthInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FirebaseAuthInput.prototype, "idToken", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FirebaseAuthInput.prototype, "provider", void 0);
exports.FirebaseAuthInput = FirebaseAuthInput = __decorate([
    (0, graphql_1.InputType)()
], FirebaseAuthInput);
//# sourceMappingURL=auth.types.js.map