import { OnModuleInit } from '@nestjs/common';
export declare class UserService implements OnModuleInit {
    private userServiceClient;
    onModuleInit(): Promise<void>;
    getUserByEmail(email: string): Promise<any>;
    createNewUser(userData: any): Promise<any>;
    loginUser(email: string, password: string): Promise<any>;
    forgetPassword(email: string, type: string): Promise<any>;
    verifyCode(email: string, verifyCode: string): Promise<any>;
    resetPassword(email: string, password: string): Promise<any>;
    updateUserData(userData: any, token?: string): Promise<any>;
    deleteUserData(userId: string, token?: string): Promise<any>;
    validateToken(token: string): Promise<any>;
}
