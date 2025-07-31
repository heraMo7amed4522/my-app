import { Resolver, Mutation, Args } from '@nestjs/graphql';
import { AuthService } from './auth.service';
import { AuthResponse, FirebaseAuthInput } from './auth.types';

@Resolver()
export class AuthResolver {
  constructor(private readonly authService: AuthService) {}

  @Mutation(() => AuthResponse)
  async authenticateWithFirebase(
    @Args('input') input: FirebaseAuthInput,
  ): Promise<AuthResponse> {
    return await this.authService.authenticateWithFirebase(
      input.idToken,
      input.provider,
    );
  }
}