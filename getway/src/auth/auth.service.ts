import { Injectable, OnModuleInit } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';

@Injectable()
export class AuthService implements OnModuleInit {
  private authServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/auth.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const authProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.AUTH_SERVICE_URL || 'localhost:50052';
    
    this.authServiceClient = new authProto.auth.AuthService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  async authenticateWithFirebase(idToken: string, provider: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.authServiceClient.AuthenticateWithFirebase(
        { idToken, provider },
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