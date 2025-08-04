import { UserService } from './user.service';
import { UserResponse, CreateUserInput, UpdateUserInput, TokenValidationResponse } from './user.types';
export declare class UserResolver {
    private readonly userService;
    constructor(userService: UserService);
    getUserByEmail(email: string): Promise<UserResponse>;
    createUser(input: CreateUserInput): Promise<UserResponse>;
    loginUser(email: string, password: string): Promise<UserResponse>;
    forgetPassword(email: string, type: string): Promise<UserResponse>;
    verifyCode(email: string, verifyCode: string): Promise<UserResponse>;
    resetPassword(email: string, password: string): Promise<UserResponse>;
    updateUser(input: UpdateUserInput, context: any): Promise<UserResponse>;
    deleteUser(userId: string, context: any): Promise<UserResponse>;
    validateToken(token: string): Promise<TokenValidationResponse>;
}
