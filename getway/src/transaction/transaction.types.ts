import { ObjectType, Field, InputType, Int } from '@nestjs/graphql';

@ObjectType()
export class Transaction {
  @Field()
  id: string;

  @Field()
  userID: string;

  @Field()
  cardID: string;

  @Field()
  merchantID: string;

  @Field()
  merchantName: string;

  @Field()
  cardNumber: string;

  @Field()
  merchantCategory: string;

  @Field()
  amount: string;

  @Field()
  currency: string;

  @Field()
  status: string;

  @Field()
  createAt: string;

  @Field()
  updateAt: string;
}

@ObjectType()
export class TransactionResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => Transaction, { nullable: true })
  transaction?: Transaction;
}

@ObjectType()
export class TransactionsResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => [Transaction], { nullable: true })
  transactions?: Transaction[];
}

@InputType()
export class CreateTransactionInput {
  @Field()
  userID: string;

  @Field()
  cardID: string;

  @Field()
  merchantID: string;

  @Field()
  merchantName: string;

  @Field()
  cardNumber: string;

  @Field()
  merchantCategory: string;

  @Field()
  amount: string;

  @Field()
  currency: string;
}