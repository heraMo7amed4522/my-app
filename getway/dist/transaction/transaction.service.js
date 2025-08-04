"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.TransactionService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let TransactionService = class TransactionService {
    transactionServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/transaction.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const transactionProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.TRANSACTION_SERVICE_URL || 'localhost:50054';
        this.transactionServiceClient = new transactionProto.transaction.TransactionService(serviceUrl, grpc.credentials.createInsecure());
    }
    async getTransactionByID(id, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', token);
            }
            this.transactionServiceClient.GetTransactionByID({ id }, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getTransactionByUserID(userID, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', token);
            }
            this.transactionServiceClient.GetTransactionByUserID({ userID }, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getTransactionByCardID(cardID, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', token);
            }
            this.transactionServiceClient.GetTransactionByCardID({ cardID }, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getTransactionByStatus(status, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', token);
            }
            this.transactionServiceClient.GetTransactionByStatus({ status }, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getTransactionByDate(date, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', token);
            }
            this.transactionServiceClient.GetTransactionByDate({ date }, metadata, (error, response) => {
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
exports.TransactionService = TransactionService;
exports.TransactionService = TransactionService = __decorate([
    (0, common_1.Injectable)()
], TransactionService);
//# sourceMappingURL=transaction.service.js.map