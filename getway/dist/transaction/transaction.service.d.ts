import { OnModuleInit } from '@nestjs/common';
export declare class TransactionService implements OnModuleInit {
    private transactionServiceClient;
    onModuleInit(): Promise<void>;
    getTransactionByID(id: string, token?: string): Promise<any>;
    getTransactionByUserID(userID: string, token?: string): Promise<any>;
    getTransactionByCardID(cardID: string, token?: string): Promise<any>;
    getTransactionByStatus(status: string, token?: string): Promise<any>;
    getTransactionByDate(date: string, token?: string): Promise<any>;
}
