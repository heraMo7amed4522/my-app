"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.PharaohService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let PharaohService = class PharaohService {
    pharaohServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/pharaoh.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const pharaohProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.PHARAOH_SERVICE_URL || 'localhost:50055';
        this.pharaohServiceClient = new pharaohProto.pharaohs.PharaohService(serviceUrl, grpc.credentials.createInsecure());
    }
    async getPharaohById(id) {
        return new Promise((resolve, reject) => {
            this.pharaohServiceClient.GetPharaohByID({ id }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getAllPharaohs(input) {
        return new Promise((resolve, reject) => {
            this.pharaohServiceClient.GetAllPharaohs({
                page: input.page || 1,
                limit: input.limit || 10,
                sort_by: input.sortBy || 'dynasty',
                order: input.order || 'asc',
            }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getPharaohsByDynasty(input) {
        return new Promise((resolve, reject) => {
            this.pharaohServiceClient.GetPharaohsByDynasty({ dynasty: input.dynasty }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getPharaohsByPeriod(input) {
        return new Promise((resolve, reject) => {
            this.pharaohServiceClient.GetPharaohsByPeriod({ period: input.period }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async searchPharaohs(input) {
        return new Promise((resolve, reject) => {
            this.pharaohServiceClient.SearchPharaohs({
                query: input.query,
                fields: input.fields || ['name', 'epithet', 'major_achievements'],
            }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async createPharaoh(input) {
        return new Promise((resolve, reject) => {
            this.pharaohServiceClient.CreatePharaoh({ pharaoh: input }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updatePharaoh(input) {
        return new Promise((resolve, reject) => {
            this.pharaohServiceClient.UpdatePharaoh({
                id: input.id,
                pharaoh: input.pharaoh,
            }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async deletePharaoh(id) {
        return new Promise((resolve, reject) => {
            this.pharaohServiceClient.DeletePharaoh({ id }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getPharaohsByRarity(input) {
        return new Promise((resolve, reject) => {
            this.pharaohServiceClient.GetPharaohsByRarity({ rarity: input.rarity }, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updatePharaohRating(input) {
        return new Promise((resolve, reject) => {
            this.pharaohServiceClient.UpdatePharaohRating({
                id: input.id,
                rating: input.rating,
            }, (error, response) => {
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
exports.PharaohService = PharaohService;
exports.PharaohService = PharaohService = __decorate([
    (0, common_1.Injectable)()
], PharaohService);
//# sourceMappingURL=pharaoh.service.js.map