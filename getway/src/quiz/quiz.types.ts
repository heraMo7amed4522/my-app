import { ObjectType, Field, InputType, Int } from '@nestjs/graphql';

@ObjectType()
export class Quiz {
  @Field()
  id: string;

  @Field()
  sectionId: string;

  @Field()
  question: string;

  @Field(() => [QuizOption])
  options: QuizOption[];

  @Field()
  correctAnswer: string;

  @Field({ nullable: true })
  explanation?: string;

  @Field()
  difficulty: string;

  @Field()
  createdAt: string;
}

@ObjectType()
export class QuizOption {
  @Field()
  key: string;

  @Field()
  value: string;
}

@ObjectType()
export class QuizResponse {
  @Field()
  id: string;

  @Field()
  userId: string;

  @Field()
  quizId: string;

  @Field()
  selectedOption: string;

  @Field()
  isCorrect: boolean;

  @Field()
  answeredAt: string;
}

@ObjectType()
export class QuizList {
  @Field(() => [Quiz])
  quizzes: Quiz[];

  @Field(() => Int)
  totalCount: number;

  @Field(() => Int)
  page: number;

  @Field(() => Int)
  limit: number;
}

@ObjectType()
export class QuizStatistics {
  @Field(() => Int)
  totalAttempts: number;

  @Field(() => Int)
  correctAnswers: number;

  @Field(() => Int)
  uniqueUsers: number;

  @Field()
  successRate: string;
}

@ObjectType()
export class QuizServiceResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => Quiz, { nullable: true })
  quiz?: Quiz;

  @Field(() => QuizList, { nullable: true })
  quizzes?: QuizList;

  @Field(() => [QuizResponse], { nullable: true })
  quizResponses?: QuizResponse[];

  @Field(() => QuizStatistics, { nullable: true })
  statistics?: QuizStatistics;

  @Field({ nullable: true })
  success?: boolean;
}

@InputType()
export class CreateQuizInput {
  @Field()
  sectionId: string;

  @Field()
  question: string;

  @Field(() => [QuizOptionInput])
  options: QuizOptionInput[];

  @Field()
  correctAnswer: string;

  @Field({ nullable: true })
  explanation?: string;

  @Field()
  difficulty: string;
}

@InputType()
export class QuizOptionInput {
  @Field()
  key: string;

  @Field()
  value: string;
}

@InputType()
export class UpdateQuizInput {
  @Field({ nullable: true })
  sectionId?: string;

  @Field({ nullable: true })
  question?: string;

  @Field(() => [QuizOptionInput], { nullable: true })
  options?: QuizOptionInput[];

  @Field({ nullable: true })
  correctAnswer?: string;

  @Field({ nullable: true })
  explanation?: string;

  @Field({ nullable: true })
  difficulty?: string;
}

@InputType()
export class SubmitQuizAnswerInput {
  @Field()
  userId: string;

  @Field()
  quizId: string;

  @Field()
  selectedOption: string;
}