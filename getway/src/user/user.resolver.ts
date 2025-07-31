import { Resolver, Query, Mutation, Args, Context } from '@nestjs/graphql';
import { UserService } from './user.service';
import {
  User,
  UserResponse,
  CreateUserInput,
  UpdateUserInput,
  TokenValidationResponse,
} from './user.types';

@Resolver(() => User)
export class UserResolver {
  constructor(private readonly userService: UserService) {}

  @Query(() => UserResponse)
  async getUserByEmail(@Args('email') email: string): Promise<UserResponse> {
    return await this.userService.getUserByEmail(email);
  }

  @Mutation(() => UserResponse)
  async createUser(
    @Args('input') input: CreateUserInput,
  ): Promise<UserResponse> {
    return await this.userService.createNewUser(input);
  }

  @Mutation(() => UserResponse)
  async loginUser(
    @Args('email') email: string,
    @Args('password') password: string,
  ): Promise<UserResponse> {
    return await this.userService.loginUser(email, password);
  }

  @Mutation(() => UserResponse)
  async forgetPassword(
    @Args('email') email: string,
    @Args('type') type: string,
  ): Promise<UserResponse> {
    return await this.userService.forgetPassword(email, type);
  }

  @Mutation(() => UserResponse)
  async verifyCode(
    @Args('email') email: string,
    @Args('verifyCode') verifyCode: string,
  ): Promise<UserResponse> {
    return await this.userService.verifyCode(email, verifyCode);
  }

  @Mutation(() => UserResponse)
  async resetPassword(
    @Args('email') email: string,
    @Args('password') password: string,
  ): Promise<UserResponse> {
    return await this.userService.resetPassword(email, password);
  }

  @Mutation(() => UserResponse)
  async updateUser(
    @Args('input') input: UpdateUserInput,
    @Context() context: any,
  ): Promise<UserResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.userService.updateUserData(input, token);
  }

  @Mutation(() => UserResponse)
  async deleteUser(
    @Args('userId') userId: string,
    @Context() context: any,
  ): Promise<UserResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.userService.deleteUserData(userId, token);
  }
  @Query(() => TokenValidationResponse)
  async validateToken(
    @Args('token') token: string,
  ): Promise<TokenValidationResponse> {
    return await this.userService.validateToken(token);
  }
}
