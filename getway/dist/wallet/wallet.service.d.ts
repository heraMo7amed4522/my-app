import { OnModuleInit } from '@nestjs/common';
export declare class WalletService implements OnModuleInit {
    private walletServiceClient;
    onModuleInit(): Promise<void>;
    getWalletByUserID(userID: string, token: string): Promise<any>;
    createWallet(walletData: any, token: string): Promise<any>;
    fundWallet(fundData: any, token: string): Promise<any>;
    deductFromWallet(deductData: any, token: string): Promise<any>;
    getWalletBalance(userID: string, token: string): Promise<any>;
    getWalletTransactions(transactionData: any, token: string): Promise<any>;
    refundToWallet(refundData: any, token: string): Promise<any>;
}
