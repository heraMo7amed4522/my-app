import { CardService } from './card.service';
import { CardResponse, CardsResponse, PaymentResponse, CreateCardInput, UpdateCardInput, UpdateCardStatusInput, SearchCardInput, StripePaymentInput, GetCardsByUserInput, FindCardsByStatusInput } from './card.types';
export declare class CardResolver {
    private readonly cardService;
    constructor(cardService: CardService);
    getCardById(id: string, context: any): Promise<CardResponse>;
    createCard(input: CreateCardInput, context: any): Promise<CardResponse>;
    updateCard(input: UpdateCardInput, context: any): Promise<CardResponse>;
    deleteCard(id: string, context: any): Promise<CardResponse>;
    updateCardStatus(input: UpdateCardStatusInput, context: any): Promise<CardResponse>;
    getCardsByUser(input: GetCardsByUserInput, context: any): Promise<CardsResponse>;
    searchCards(input: SearchCardInput, context: any): Promise<CardsResponse>;
    findCardsByStatus(input: FindCardsByStatusInput, context: any): Promise<CardsResponse>;
    processStripePayment(input: StripePaymentInput, context: any): Promise<PaymentResponse>;
}
