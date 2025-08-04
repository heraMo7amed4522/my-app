import { AuthService } from './auth.service';
import { AuthResponse, FirebaseAuthInput } from './auth.types';
export declare class AuthResolver {
    private readonly authService;
    constructor(authService: AuthService);
    authenticateWithFirebase(input: FirebaseAuthInput): Promise<AuthResponse>;
}
