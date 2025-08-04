import { OnModuleInit } from '@nestjs/common';
import { CreateCardInput, UpdateCardInput, UpdateCardStatusInput, SearchCardInput, StripePaymentInput, GetCardsByUserInput, FindCardsByStatusInput } from './card.types';
export declare class CardService implements OnModuleInit {
    private cardServiceClient;
    onModuleInit(): Promise<void>;
    getCardById(id: string, token?: string): Promise<any>;
    createCard(input: CreateCardInput, token?: string): Promise<any>;
    updateCard(input: UpdateCardInput, token?: string): Promise<any>;
    deleteCard(id: string, token?: string): Promise<any>;
    updateCardStatus(input: UpdateCardStatusInput, token?: string): Promise<any>;
    getCardsByUser(input: GetCardsByUserInput, token?: string): Promise<any>;
    searchCards(input: SearchCardInput, token?: string): Promise<any>;
    findCardsByStatus(input: FindCardsByStatusInput, token?: string): Promise<any>;
    processStripePayment(input: StripePaymentInput, token?: string): Promise<any>;
}
