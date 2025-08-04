export declare class Transaction {
    id: string;
    userID: string;
    cardID: string;
    merchantID: string;
    merchantName: string;
    cardNumber: string;
    merchantCategory: string;
    amount: string;
    currency: string;
    status: string;
    createAt: string;
    updateAt: string;
}
export declare class TransactionResponse {
    statusCode: number;
    message: string;
    transaction?: Transaction;
}
export declare class TransactionsResponse {
    statusCode: number;
    message: string;
    transactions?: Transaction[];
}
export declare class CreateTransactionInput {
    userID: string;
    cardID: string;
    merchantID: string;
    merchantName: string;
    cardNumber: string;
    merchantCategory: string;
    amount: string;
    currency: string;
}
