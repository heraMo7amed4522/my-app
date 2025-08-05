import { ObjectType, Field, InputType, Int, Float } from '@nestjs/graphql';

@ObjectType()
export class Feedback {
  @Field()
  id: string;

  @Field()
  userId: string;

  @Field()
  templateId: string;

  @Field(() => Float)
  rating: number;

  @Field()
  comment: string;

  @Field()
  language: string;

  @Field()
  createdAt: string;
}

@ObjectType()
export class FeedbackList {
  @Field(() => [Feedback])
  feedback: Feedback[];

  @Field(() => Int)
  totalCount: number;

  @Field(() => Int)
  page: number;

  @Field(() => Int)
  limit: number;
}

@ObjectType()
export class FeedbackStatistics {
  @Field()
  templateId: string;

  @Field(() => Float)
  averageRating: number;

  @Field(() => Int)
  totalRatings: number;

  @Field(() => [RatingDistribution])
  ratingDistribution: RatingDistribution[];
}

@ObjectType()
export class RatingDistribution {
  @Field()
  rating: string;

  @Field(() => Int)
  count: number;
}

@ObjectType()
export class FeedbackResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => Feedback, { nullable: true })
  feedback?: Feedback;

  @Field(() => FeedbackList, { nullable: true })
  feedbackList?: FeedbackList;

  @Field(() => FeedbackStatistics, { nullable: true })
  statistics?: FeedbackStatistics;

  @Field({ nullable: true })
  success?: boolean;
}

@InputType()
export class SubmitFeedbackInput {
  @Field()
  userId: string;

  @Field()
  templateId: string;

  @Field(() => Float)
  rating: number;

  @Field()
  comment: string;

  @Field({ nullable: true })
  language?: string;
}

@InputType()
export class UpdateFeedbackInput {
  @Field(() => Float)
  rating: number;

  @Field()
  comment: string;
}