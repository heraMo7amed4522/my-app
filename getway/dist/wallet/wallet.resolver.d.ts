import { WalletService } from './wallet.service';
import { WalletResponse, WalletTransactionsResponse, WalletBalanceResponse, CreateWalletInput, FundWalletInput, DeductFromWalletInput, GetWalletTransactionsInput, RefundToWalletInput } from './wallet.types';
export declare class WalletResolver {
    private readonly walletService;
    constructor(walletService: WalletService);
    getWalletByUserID(userID: string, context: any): Promise<WalletResponse>;
    createWallet(input: CreateWalletInput, context: any): Promise<WalletResponse>;
    fundWallet(input: FundWalletInput, context: any): Promise<WalletResponse>;
    deductFromWallet(input: DeductFromWalletInput, context: any): Promise<WalletResponse>;
    getWalletBalance(userID: string, context: any): Promise<WalletBalanceResponse>;
    getWalletTransactions(input: GetWalletTransactionsInput, context: any): Promise<WalletTransactionsResponse>;
    refundToWallet(input: RefundToWalletInput, context: any): Promise<WalletResponse>;
}
