"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.WalletService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let WalletService = class WalletService {
    walletServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/wallet.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const walletProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.WALLET_SERVICE_URL || 'localhost:50055';
        this.walletServiceClient = new walletProto.wallet.WalletService(serviceUrl, grpc.credentials.createInsecure());
    }
    async getWalletByUserID(userID, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            metadata.add('authorization', `Bearer ${token}`);
            this.walletServiceClient.GetWalletByUserID({ userID }, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async createWallet(walletData, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            metadata.add('authorization', `Bearer ${token}`);
            this.walletServiceClient.CreateWallet(walletData, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async fundWallet(fundData, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            metadata.add('authorization', `Bearer ${token}`);
            this.walletServiceClient.FundWallet(fundData, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async deductFromWallet(deductData, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            metadata.add('authorization', `Bearer ${token}`);
            this.walletServiceClient.DeductFromWallet(deductData, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getWalletBalance(userID, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            metadata.add('authorization', `Bearer ${token}`);
            this.walletServiceClient.GetWalletBalance({ userID }, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getWalletTransactions(transactionData, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            metadata.add('authorization', `Bearer ${token}`);
            this.walletServiceClient.GetWalletTransactions(transactionData, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async refundToWallet(refundData, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            metadata.add('authorization', `Bearer ${token}`);
            this.walletServiceClient.RefundToWallet(refundData, metadata, (error, response) => {
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
exports.WalletService = WalletService;
exports.WalletService = WalletService = __decorate([
    (0, common_1.Injectable)()
], WalletService);
//# sourceMappingURL=wallet.service.js.map