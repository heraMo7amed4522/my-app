import { Resolver, Query, Mutation, Args, Context } from '@nestjs/graphql';
import { WalletService } from './wallet.service';
import {
  Wallet,
  WalletResponse,
  WalletTransactionsResponse,
  WalletBalanceResponse,
  CreateWalletInput,
  FundWalletInput,
  DeductFromWalletInput,
  GetWalletTransactionsInput,
  RefundToWalletInput,
} from './wallet.types';

@Resolver(() => Wallet)
export class WalletResolver {
  constructor(private readonly walletService: WalletService) {}

  @Query(() => WalletResponse)
  async getWalletByUserID(
    @Args('userID') userID: string,
    @Context() context: any,
  ): Promise<WalletResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.walletService.getWalletByUserID(userID, token);
  }

  @Mutation(() => WalletResponse)
  async createWallet(
    @Args('input') input: CreateWalletInput,
    @Context() context: any,
  ): Promise<WalletResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.walletService.createWallet(input, token);
  }

  @Mutation(() => WalletResponse)
  async fundWallet(
    @Args('input') input: FundWalletInput,
    @Context() context: any,
  ): Promise<WalletResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.walletService.fundWallet(input, token);
  }

  @Mutation(() => WalletResponse)
  async deductFromWallet(
    @Args('input') input: DeductFromWalletInput,
    @Context() context: any,
  ): Promise<WalletResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.walletService.deductFromWallet(input, token);
  }

  @Query(() => WalletBalanceResponse)
  async getWalletBalance(
    @Args('userID') userID: string,
    @Context() context: any,
  ): Promise<WalletBalanceResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.walletService.getWalletBalance(userID, token);
  }

  @Query(() => WalletTransactionsResponse)
  async getWalletTransactions(
    @Args('input') input: GetWalletTransactionsInput,
    @Context() context: any,
  ): Promise<WalletTransactionsResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.walletService.getWalletTransactions(input, token);
  }

  @Mutation(() => WalletResponse)
  async refundToWallet(
    @Args('input') input: RefundToWalletInput,
    @Context() context: any,
  ): Promise<WalletResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.walletService.refundToWallet(input, token);
  }
}