import { TransactionService } from './transaction.service';
import { TransactionResponse, TransactionsResponse } from './transaction.types';
export declare class TransactionResolver {
    private readonly transactionService;
    constructor(transactionService: TransactionService);
    getTransactionByID(id: string, context: any): Promise<TransactionResponse>;
    getTransactionsByUserID(userID: string, context: any): Promise<TransactionsResponse>;
    getTransactionsByCardID(cardID: string, context: any): Promise<TransactionsResponse>;
    getTransactionsByStatus(status: string, context: any): Promise<TransactionsResponse>;
    getTransactionsByDate(date: string, context: any): Promise<TransactionsResponse>;
}
