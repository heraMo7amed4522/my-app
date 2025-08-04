"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.AuthService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let AuthService = class AuthService {
    authServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/auth.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const authProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.AUTH_SERVICE_URL || 'localhost:50052';
        this.authServiceClient = new authProto.auth.AuthService(serviceUrl, grpc.credentials.createInsecure());
    }
    async authenticateWithFirebase(idToken, provider) {
        return new Promise((resolve, reject) => {
            this.authServiceClient.AuthenticateWithFirebase({ idToken, provider }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
};
exports.AuthService = AuthService;
exports.AuthService = AuthService = __decorate([
    (0, common_1.Injectable)()
], AuthService);
//# sourceMappingURL=auth.service.js.map