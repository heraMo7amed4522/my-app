import { ObjectType, Field, InputType, Int } from '@nestjs/graphql';

@ObjectType()
export class Card {
  @Field()
  id: string;

  @Field()
  userId: string;

  @Field()
  cardNumber: string;

  @Field()
  cardHolderName: string;

  @Field()
  expirationDate: string;

  @Field()
  cvv: string;

  @Field()
  status: string;

  @Field()
  cardType: string;

  @Field()
  balance: string;

  @Field()
  createAt: string;

  @Field()
  updateAt: string;
}

@ObjectType()
export class ErrorDetails {
  @Field(() => Int)
  code: number;

  @Field()
  message: string;

  @Field(() => [String])
  details: string[];

  @Field()
  timestamp: string;
}

@ObjectType()
export class CardResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => Card, { nullable: true })
  card?: Card;

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;
}

@ObjectType()
export class CardsResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => [Card], { nullable: true })
  cards?: Card[];

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;

  @Field(() => Int, { nullable: true })
  totalCount?: number;

  @Field(() => Int, { nullable: true })
  currentPage?: number;

  @Field(() => Int, { nullable: true })
  totalPages?: number;
}

@ObjectType()
export class PaymentResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field({ nullable: true })
  paymentIntentId?: string;

  @Field({ nullable: true })
  clientSecret?: string;

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;
}

// Input Types
@InputType()
export class CreateCardInput {
  @Field()
  userId: string;

  @Field()
  cardNumber: string;

  @Field()
  cardHolderName: string;

  @Field()
  expirationDate: string;

  @Field()
  cvv: string;

  @Field({ nullable: true })
  cardType?: string;

  @Field({ nullable: true })
  balance?: string;
}

@InputType()
export class UpdateCardInput {
  @Field()
  id: string;

  @Field({ nullable: true })
  cardHolderName?: string;

  @Field({ nullable: true })
  expirationDate?: string;

  @Field({ nullable: true })
  cvv?: string;

  @Field({ nullable: true })
  cardType?: string;

  @Field({ nullable: true })
  balance?: string;
}

@InputType()
export class UpdateCardStatusInput {
  @Field()
  id: string;

  @Field()
  status: string;
}

@InputType()
export class SearchCardInput {
  @Field({ nullable: true })
  userId?: string;

  @Field({ nullable: true })
  cardHolderName?: string;

  @Field({ nullable: true })
  cardType?: string;

  @Field({ nullable: true })
  status?: string;

  @Field(() => Int, { nullable: true })
  page?: number;

  @Field(() => Int, { nullable: true })
  limit?: number;
}

@InputType()
export class StripePaymentInput {
  @Field()
  cardId: string;

  @Field(() => Int)
  amount: number;

  @Field()
  currency: string;

  @Field({ nullable: true })
  description?: string;
}

@InputType()
export class GetCardsByUserInput {
  @Field()
  userId: string;

  @Field(() => Int, { nullable: true })
  page?: number;

  @Field(() => Int, { nullable: true })
  limit?: number;
}

@InputType()
export class FindCardsByStatusInput {
  @Field()
  status: string;

  @Field(() => Int, { nullable: true })
  page?: number;

  @Field(() => Int, { nullable: true })
  limit?: number;
}