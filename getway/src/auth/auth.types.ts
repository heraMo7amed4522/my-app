import { ObjectType, Field, InputType, Int } from '@nestjs/graphql';

@ObjectType()
export class UserInfo {
  @Field()
  id: string;

  @Field()
  email: string;

  @Field()
  fullName: string;

  @Field()
  role: string;
}

@ObjectType()
export class AuthTokens {
  @Field()
  accessToken: string;

  @Field(() => Int)
  expiresIn: number;

  @Field(() => UserInfo)
  user: UserInfo;
}

@ObjectType()
export class AuthResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => AuthTokens, { nullable: true })
  tokens?: AuthTokens;

  @Field({ nullable: true })
  error?: string;
}

@InputType()
export class FirebaseAuthInput {
  @Field()
  idToken: string;

  @Field()
  provider: string;
}