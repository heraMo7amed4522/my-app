import { Resolver, Query, Mutation, Args, Context } from '@nestjs/graphql';
import { CardService } from './card.service';
import {
  Card,
  CardResponse,
  CardsResponse,
  PaymentResponse,
  CreateCardInput,
  UpdateCardInput,
  UpdateCardStatusInput,
  SearchCardInput,
  StripePaymentInput,
  GetCardsByUserInput,
  FindCardsByStatusInput,
} from './card.types';

@Resolver(() => Card)
export class CardResolver {
  constructor(private readonly cardService: CardService) {}

  @Query(() => CardResponse)
  async getCardById(
    @Args('id') id: string,
    @Context() context: any,
  ): Promise<CardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.cardService.getCardById(id, token);
  }

  @Mutation(() => CardResponse)
  async createCard(
    @Args('input') input: CreateCardInput,
    @Context() context: any,
  ): Promise<CardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.cardService.createCard(input, token);
  }

  @Mutation(() => CardResponse)
  async updateCard(
    @Args('input') input: UpdateCardInput,
    @Context() context: any,
  ): Promise<CardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.cardService.updateCard(input, token);
  }

  @Mutation(() => CardResponse)
  async deleteCard(
    @Args('id') id: string,
    @Context() context: any,
  ): Promise<CardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.cardService.deleteCard(id, token);
  }

  @Mutation(() => CardResponse)
  async updateCardStatus(
    @Args('input') input: UpdateCardStatusInput,
    @Context() context: any,
  ): Promise<CardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.cardService.updateCardStatus(input, token);
  }

  @Query(() => CardsResponse)
  async getCardsByUser(
    @Args('input') input: GetCardsByUserInput,
    @Context() context: any,
  ): Promise<CardsResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.cardService.getCardsByUser(input, token);
  }

  @Query(() => CardsResponse)
  async searchCards(
    @Args('input') input: SearchCardInput,
    @Context() context: any,
  ): Promise<CardsResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.cardService.searchCards(input, token);
  }

  @Query(() => CardsResponse)
  async findCardsByStatus(
    @Args('input') input: FindCardsByStatusInput,
    @Context() context: any,
  ): Promise<CardsResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.cardService.findCardsByStatus(input, token);
  }

  @Mutation(() => PaymentResponse)
  async processStripePayment(
    @Args('input') input: StripePaymentInput,
    @Context() context: any,
  ): Promise<PaymentResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.cardService.processStripePayment(input, token);
  }
}