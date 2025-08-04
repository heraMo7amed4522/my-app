import { ErrorDetails } from '../shared/shared.types';
export declare class Wallet {
    id: string;
    userID: string;
    balance: string;
    currency: string;
    createdAt: string;
    updatedAt: string;
}
export declare class WalletTransaction {
    id: string;
    walletId: string;
    type: string;
    amount: string;
    currency: string;
    description?: string;
    createdAt: string;
}
export declare class WalletResponse {
    statusCode: number;
    message: string;
    wallet?: Wallet;
    error?: ErrorDetails;
}
export declare class WalletTransactionsResponse {
    statusCode: number;
    message: string;
    transactions?: WalletTransaction[];
    error?: ErrorDetails;
}
export declare class WalletBalanceResponse {
    statusCode: number;
    message: string;
    balance?: string;
    currency?: string;
    error?: ErrorDetails;
}
export declare class CreateWalletInput {
    userID: string;
    currency?: string;
}
export declare class FundWalletInput {
    userID: string;
    cardID: string;
    amount: string;
    currency: string;
}
export declare class DeductFromWalletInput {
    userID: string;
    amount: string;
    currency: string;
    transactionID: string;
}
export declare class GetWalletTransactionsInput {
    userID: string;
    limit?: number;
    offset?: number;
}
export declare class RefundToWalletInput {
    userID: string;
    amount: string;
    currency: string;
    transactionID: string;
    reason?: string;
}
