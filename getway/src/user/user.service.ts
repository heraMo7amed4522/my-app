import { Injectable, OnModuleInit } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';

@Injectable()
export class UserService implements OnModuleInit {
  private userServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/user.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const userProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.USER_SERVICE_URL || 'localhost:50051';
    
    this.userServiceClient = new userProto.users.UserService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  async getUserByEmail(email: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userServiceClient.GetUserByEmail(
        { email },
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

  async createNewUser(userData: any): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userServiceClient.CreateNewUser(
        userData,
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

  async loginUser(email: string, password: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userServiceClient.LoginUser(
        { email, password },
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

  async forgetPassword(email: string, type: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userServiceClient.ForgetPassword(
        { email, type },
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

  async verifyCode(email: string, verifyCode: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userServiceClient.VerifyCode(
        { email, verifyCode },
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

  async resetPassword(email: string, password: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userServiceClient.ResetPassword(
        { email, password },
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

  async updateUserData(userData: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
       const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.userServiceClient.UpdateUserData(
        userData,
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

  async deleteUserData(userId: string, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
       const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.userServiceClient.DeleteUserData(
        { id: userId },
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
  async validateToken(token: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userServiceClient.ValidateToken(
        { token },
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