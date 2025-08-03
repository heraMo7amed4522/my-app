import { ObjectType, Field, Int } from '@nestjs/graphql';

@ObjectType()
export class ErrorDetails {
  @Field(() => Int)
  code: number;

  @Field()
  message: string;

  @Field(() => [String])
  details: string[];

  @Field()
  timestamp: string;
}