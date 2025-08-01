import { Injectable, OnModuleInit } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';
import {
  CreateCardInput,
  UpdateCardInput,
  UpdateCardStatusInput,
  SearchCardInput,
  StripePaymentInput,
  GetCardsByUserInput,
  FindCardsByStatusInput,
} from './card.types';

@Injectable()
export class CardService implements OnModuleInit {
  private cardServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/card.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const cardProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.CARD_SERVICE_URL || 'localhost:50053';
    
    this.cardServiceClient = new cardProto.cards.CardService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  async getCardById(id: string, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }

      this.cardServiceClient.GetCardByID(
        { id },
        metadata,
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

  async createCard(input: CreateCardInput, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }

      this.cardServiceClient.CreateNewCard(
        input,
        metadata,
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

  async updateCard(input: UpdateCardInput, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }

      this.cardServiceClient.UpdateCardData(
        input,
        metadata,
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

  async deleteCard(id: string, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }

      this.cardServiceClient.DeleteCardData(
        { id },
        metadata,
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

  async updateCardStatus(input: UpdateCardStatusInput, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }

      this.cardServiceClient.UpdateCardStatus(
        input,
        metadata,
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

  async getCardsByUser(input: GetCardsByUserInput, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }

      this.cardServiceClient.GetCardByUserID(
        input,
        metadata,
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

  async searchCards(input: SearchCardInput, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }

      this.cardServiceClient.SearchCard(
        input,
        metadata,
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

  async findCardsByStatus(input: FindCardsByStatusInput, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }

      this.cardServiceClient.FindCardByStatus(
        input,
        metadata,
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

  async processStripePayment(input: StripePaymentInput, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }

      this.cardServiceClient.StripePayment(
        input,
        metadata,
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