import { Injectable, OnModuleInit } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';

@Injectable()
export class TransactionService implements OnModuleInit {
  private transactionServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/transaction.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const transactionProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.TRANSACTION_SERVICE_URL || 'localhost:50054';
    
    this.transactionServiceClient = new transactionProto.transaction.TransactionService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  async getTransactionByID(id: string, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', token);
      }
      
      this.transactionServiceClient.GetTransactionByID(
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

  async getTransactionByUserID(userID: string, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', token);
      }
      
      this.transactionServiceClient.GetTransactionByUserID(
        { userID },
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

  async getTransactionByCardID(cardID: string, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', token);
      }
      
      this.transactionServiceClient.GetTransactionByCardID(
        { cardID },
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

  async getTransactionByStatus(status: string, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', token);
      }
      
      this.transactionServiceClient.GetTransactionByStatus(
        { status },
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

  async getTransactionByDate(date: string, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', token);
      }
      
      this.transactionServiceClient.GetTransactionByDate(
        { date },
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