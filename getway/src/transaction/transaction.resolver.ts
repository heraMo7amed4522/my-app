import { Resolver, Query, Mutation, Args, Context } from '@nestjs/graphql';
import { TransactionService } from './transaction.service';
import {
  Transaction,
  TransactionResponse,
  TransactionsResponse,
} from './transaction.types';

@Resolver(() => Transaction)
export class TransactionResolver {
  constructor(private readonly transactionService: TransactionService) {}

  @Query(() => TransactionResponse)
  async getTransactionByID(
    @Args('id') id: string,
    @Context() context: any,
  ): Promise<TransactionResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.transactionService.getTransactionByID(id, token);
  }

  @Query(() => TransactionsResponse)
  async getTransactionsByUserID(
    @Args('userID') userID: string,
    @Context() context: any,
  ): Promise<TransactionsResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.transactionService.getTransactionByUserID(userID, token);
  }

  @Query(() => TransactionsResponse)
  async getTransactionsByCardID(
    @Args('cardID') cardID: string,
    @Context() context: any,
  ): Promise<TransactionsResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.transactionService.getTransactionByCardID(cardID, token);
  }

  @Query(() => TransactionsResponse)
  async getTransactionsByStatus(
    @Args('status') status: string,
    @Context() context: any,
  ): Promise<TransactionsResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.transactionService.getTransactionByStatus(status, token);
  }

  @Query(() => TransactionsResponse)
  async getTransactionsByDate(
    @Args('date') date: string,
    @Context() context: any,
  ): Promise<TransactionsResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.transactionService.getTransactionByDate(date, token);
  }
}