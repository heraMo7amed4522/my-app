import { ErrorDetails } from '../shared/shared.types';
export declare class Card {
    id: string;
    userId: string;
    cardNumber: string;
    cardHolderName: string;
    expirationDate: string;
    cvv: string;
    status: string;
    cardType: string;
    balance: string;
    createAt: string;
    updateAt: string;
}
export declare class CardResponse {
    statusCode: number;
    message: string;
    card?: Card;
    error?: ErrorDetails;
}
export declare class CardsResponse {
    statusCode: number;
    message: string;
    cards?: Card[];
    error?: ErrorDetails;
    totalCount?: number;
    currentPage?: number;
    totalPages?: number;
}
export declare class PaymentResponse {
    statusCode: number;
    message: string;
    paymentIntentId?: string;
    clientSecret?: string;
    error?: ErrorDetails;
}
export declare class CreateCardInput {
    userId: string;
    cardNumber: string;
    cardHolderName: string;
    expirationDate: string;
    cvv: string;
    cardType?: string;
    balance?: string;
}
export declare class UpdateCardInput {
    id: string;
    cardHolderName?: string;
    expirationDate?: string;
    cvv?: string;
    cardType?: string;
    balance?: string;
}
export declare class UpdateCardStatusInput {
    id: string;
    status: string;
}
export declare class SearchCardInput {
    userId?: string;
    cardHolderName?: string;
    cardType?: string;
    status?: string;
    page?: number;
    limit?: number;
}
export declare class StripePaymentInput {
    cardId: string;
    amount: number;
    currency: string;
    description?: string;
}
export declare class GetCardsByUserInput {
    userId: string;
    page?: number;
    limit?: number;
}
export declare class FindCardsByStatusInput {
    status: string;
    page?: number;
    limit?: number;
}
