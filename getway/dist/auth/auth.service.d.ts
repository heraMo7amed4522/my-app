import { OnModuleInit } from '@nestjs/common';
export declare class AuthService implements OnModuleInit {
    private authServiceClient;
    onModuleInit(): Promise<void>;
    authenticateWithFirebase(idToken: string, provider: string): Promise<any>;
}
