export declare class UserInfo {
    id: string;
    email: string;
    fullName: string;
    role: string;
}
export declare class AuthTokens {
    accessToken: string;
    expiresIn: number;
    user: UserInfo;
}
export declare class AuthResponse {
    statusCode: number;
    message: string;
    tokens?: AuthTokens;
    error?: string;
}
export declare class FirebaseAuthInput {
    idToken: string;
    provider: string;
}
