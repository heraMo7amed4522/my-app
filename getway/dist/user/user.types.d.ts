import { AuthTokens } from '../auth/auth.types';
export declare class User {
    id: string;
    fullName: string;
    email: string;
    countryCode?: string;
    phoneNumber?: string;
    role: string;
    verifiyCode?: string;
    createAt: string;
    updateAt: string;
}
export declare class UserResponse {
    statusCode: number;
    message: string;
    user?: User;
    tokens?: AuthTokens;
    accessToken?: string;
    verifyCode?: string;
    response?: string;
}
export declare class TokenClaims {
    userId: string;
    email: string;
    role: string;
    exp: number;
}
export declare class TokenValidationResponse {
    statusCode: number;
    message: string;
    claims?: TokenClaims;
}
export declare class CreateUserInput {
    fullName: string;
    email: string;
    password: string;
    countryCode?: string;
    phoneNumber?: string;
    role?: string;
    verifiyCode?: string;
}
export declare class UpdateUserInput {
    email: string;
    fullName?: string;
    countryCode?: string;
    phoneNumber?: string;
    role?: string;
}
