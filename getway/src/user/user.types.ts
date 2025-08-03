import { ObjectType, Field, InputType, Int } from '@nestjs/graphql';
import { AuthTokens } from '../auth/auth.types';

@ObjectType()
export class User {
  @Field()
  id: string;

  @Field()
  fullName: string;

  @Field()
  email: string;

  @Field({ nullable: true })
  countryCode?: string;

  @Field({ nullable: true })
  phoneNumber?: string;

  @Field()
  role: string;

  @Field({ nullable: true })
  verifiyCode?: string;

  @Field()
  createAt: string;

  @Field()
  updateAt: string;
}



@ObjectType()
export class UserResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => User, { nullable: true })
  user?: User;

  // Updated: Use AuthTokens instead of simple accessToken
  @Field(() => AuthTokens, { nullable: true })
  tokens?: AuthTokens;

  // Add direct accessToken field for GraphQL compatibility
  @Field({ nullable: true })
  accessToken?: string;

  @Field({ nullable: true })
  verifyCode?: string;

  @Field({ nullable: true })
  response?: string;
}

// NEW: Add token validation types
@ObjectType()
export class TokenClaims {
  @Field()
  userId: string;

  @Field()
  email: string;

  @Field()
  role: string;

  @Field(() => Int)
  exp: number;
}

@ObjectType()
export class TokenValidationResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => TokenClaims, { nullable: true })
  claims?: TokenClaims;
}

@InputType()
export class CreateUserInput {
  @Field()
  fullName: string;

  @Field()
  email: string;

  @Field()
  password: string;

  @Field({ nullable: true })
  countryCode?: string;

  @Field({ nullable: true })
  phoneNumber?: string;

  @Field({ nullable: true })
  role?: string;

  @Field({ nullable: true })
  verifiyCode?: string;
}

@InputType()
export class UpdateUserInput {
  @Field()
  email: string;

  @Field({ nullable: true })
  fullName?: string;

  @Field({ nullable: true })
  countryCode?: string;

  @Field({ nullable: true })
  phoneNumber?: string;

  @Field({ nullable: true })
  role?: string;
}