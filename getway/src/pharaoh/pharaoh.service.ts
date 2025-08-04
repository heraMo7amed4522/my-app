import { Injectable, OnModuleInit } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';
import {
  CreatePharaohInput,
  UpdatePharaohInput,
  GetAllPharaohsInput,
  GetPharaohsByDynastyInput,
  GetPharaohsByPeriodInput,
  SearchPharaohsInput,
  GetPharaohsByRarityInput,
  UpdatePharaohRatingInput,
} from './pharaoh.types';

@Injectable()
export class PharaohService implements OnModuleInit {
  private pharaohServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/pharaoh.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const pharaohProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.PHARAOH_SERVICE_URL || 'localhost:50055';
    
    this.pharaohServiceClient = new pharaohProto.pharaohs.PharaohService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  async getPharaohById(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.pharaohServiceClient.GetPharaohByID(
        { id },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async getAllPharaohs(input: GetAllPharaohsInput): Promise<any> {
    return new Promise((resolve, reject) => {
      this.pharaohServiceClient.GetAllPharaohs(
        {
          page: input.page || 1,
          limit: input.limit || 10,
          sort_by: input.sortBy || 'dynasty',
          order: input.order || 'asc',
        },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async getPharaohsByDynasty(input: GetPharaohsByDynastyInput): Promise<any> {
    return new Promise((resolve, reject) => {
      this.pharaohServiceClient.GetPharaohsByDynasty(
        { dynasty: input.dynasty },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async getPharaohsByPeriod(input: GetPharaohsByPeriodInput): Promise<any> {
    return new Promise((resolve, reject) => {
      this.pharaohServiceClient.GetPharaohsByPeriod(
        { period: input.period },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async searchPharaohs(input: SearchPharaohsInput): Promise<any> {
    return new Promise((resolve, reject) => {
      this.pharaohServiceClient.SearchPharaohs(
        {
          query: input.query,
          fields: input.fields || ['name', 'epithet', 'major_achievements'],
        },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async createPharaoh(input: CreatePharaohInput): Promise<any> {
    return new Promise((resolve, reject) => {
      this.pharaohServiceClient.CreatePharaoh(
        { pharaoh: input },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async updatePharaoh(input: UpdatePharaohInput): Promise<any> {
    return new Promise((resolve, reject) => {
      this.pharaohServiceClient.UpdatePharaoh(
        {
          id: input.id,
          pharaoh: input.pharaoh,
        },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async deletePharaoh(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.pharaohServiceClient.DeletePharaoh(
        { id },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async getPharaohsByRarity(input: GetPharaohsByRarityInput): Promise<any> {
    return new Promise((resolve, reject) => {
      this.pharaohServiceClient.GetPharaohsByRarity(
        { rarity: input.rarity },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async updatePharaohRating(input: UpdatePharaohRatingInput): Promise<any> {
    return new Promise((resolve, reject) => {
      this.pharaohServiceClient.UpdatePharaohRating(
        {
          id: input.id,
          rating: input.rating,
        },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }
}