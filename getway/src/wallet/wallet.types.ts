import { ObjectType, Field, InputType, Int } from '@nestjs/graphql';
import { ErrorDetails } from '../shared/shared.types';

@ObjectType()
export class Wallet {
  @Field()
  id: string;

  @Field()
  userID: string;

  @Field()
  balance: string;

  @Field()
  currency: string;

  @Field()
  createdAt: string;

  @Field()
  updatedAt: string;
}

@ObjectType()
export class WalletTransaction {
  @Field()
  id: string;

  @Field()
  walletId: string;

  @Field()
  type: string;

  @Field()
  amount: string;

  @Field()
  currency: string;

  @Field({ nullable: true })
  description?: string;

  @Field()
  createdAt: string;
}



@ObjectType()
export class WalletResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => Wallet, { nullable: true })
  wallet?: Wallet;

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;
}

@ObjectType()
export class WalletTransactionsResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => [WalletTransaction], { nullable: true })
  transactions?: WalletTransaction[];

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;
}

@ObjectType()
export class WalletBalanceResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field({ nullable: true })
  balance?: string;

  @Field({ nullable: true })
  currency?: string;

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;
}

// Input Types
@InputType()
export class CreateWalletInput {
  @Field()
  userID: string;

  @Field({ nullable: true })
  currency?: string;
}

@InputType()
export class FundWalletInput {
  @Field()
  userID: string;

  @Field()
  cardID: string;

  @Field()
  amount: string;

  @Field()
  currency: string;
}

@InputType()
export class DeductFromWalletInput {
  @Field()
  userID: string;

  @Field()
  amount: string;

  @Field()
  currency: string;

  @Field()
  transactionID: string;
}

@InputType()
export class GetWalletTransactionsInput {
  @Field()
  userID: string;

  @Field(() => Int, { nullable: true })
  limit?: number;

  @Field(() => Int, { nullable: true })
  offset?: number;
}

@InputType()
export class RefundToWalletInput {
  @Field()
  userID: string;

  @Field()
  amount: string;

  @Field()
  currency: string;

  @Field()
  transactionID: string;

  @Field({ nullable: true })
  reason?: string;
}