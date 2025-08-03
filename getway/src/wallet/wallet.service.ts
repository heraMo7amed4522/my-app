import { Injectable, OnModuleInit } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';

@Injectable()
export class WalletService implements OnModuleInit {
  private walletServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/wallet.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const walletProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.WALLET_SERVICE_URL || 'localhost:50055';
    
    this.walletServiceClient = new walletProto.wallet.WalletService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  async getWalletByUserID(userID: string, token: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      metadata.add('authorization', `Bearer ${token}`);
      
      this.walletServiceClient.GetWalletByUserID(
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

  async createWallet(walletData: any, token: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      metadata.add('authorization', `Bearer ${token}`);
      
      this.walletServiceClient.CreateWallet(
        walletData,
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

  async fundWallet(fundData: any, token: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      metadata.add('authorization', `Bearer ${token}`);
      
      this.walletServiceClient.FundWallet(
        fundData,
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

  async deductFromWallet(deductData: any, token: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      metadata.add('authorization', `Bearer ${token}`);
      
      this.walletServiceClient.DeductFromWallet(
        deductData,
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

  async getWalletBalance(userID: string, token: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      metadata.add('authorization', `Bearer ${token}`);
      
      this.walletServiceClient.GetWalletBalance(
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

  async getWalletTransactions(transactionData: any, token: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      metadata.add('authorization', `Bearer ${token}`);
      
      this.walletServiceClient.GetWalletTransactions(
        transactionData,
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

  async refundToWallet(refundData: any, token: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      metadata.add('authorization', `Bearer ${token}`);
      
      this.walletServiceClient.RefundToWallet(
        refundData,
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