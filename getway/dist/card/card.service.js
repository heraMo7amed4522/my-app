"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.CardService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let CardService = class CardService {
    cardServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/card.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const cardProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.CARD_SERVICE_URL || 'localhost:50053';
        this.cardServiceClient = new cardProto.cards.CardService(serviceUrl, grpc.credentials.createInsecure());
    }
    async getCardById(id, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.cardServiceClient.GetCardByID({ id }, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async createCard(input, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.cardServiceClient.CreateNewCard(input, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updateCard(input, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.cardServiceClient.UpdateCardData(input, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async deleteCard(id, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.cardServiceClient.DeleteCardData({ id }, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updateCardStatus(input, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.cardServiceClient.UpdateCardStatus(input, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getCardsByUser(input, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.cardServiceClient.GetCardByUserID(input, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async searchCards(input, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.cardServiceClient.SearchCard(input, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async findCardsByStatus(input, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.cardServiceClient.FindCardByStatus(input, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async processStripePayment(input, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.cardServiceClient.StripePayment(input, metadata, (error, response) => {
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
exports.CardService = CardService;
exports.CardService = CardService = __decorate([
    (0, common_1.Injectable)()
], CardService);
//# sourceMappingURL=card.service.js.map