import { ObjectType, Field, InputType, Int, Float } from '@nestjs/graphql';

@ObjectType()
export class UserProgress {
  @Field()
  id: string;

  @Field()
  userId: string;

  @Field()
  templateId: string;

  @Field()
  sectionId: string;

  @Field(() => Float)
  progress: number;

  @Field()
  completed: boolean;

  @Field()
  lastViewed: string;

  @Field()
  createdAt: string;
}

@ObjectType()
export class Achievement {
  @Field()
  id: string;

  @Field()
  code: string;

  @Field()
  name: string;

  @Field()
  description: string;

  @Field({ nullable: true })
  badgeUrl?: string;

  @Field()
  createdAt: string;
}

@ObjectType()
export class UserAchievement {
  @Field()
  id: string;

  @Field()
  userId: string;

  @Field()
  achievementId: string;

  @Field(() => Achievement)
  achievement: Achievement;

  @Field()
  unlockedAt: string;
}

@ObjectType()
export class LearningPath {
  @Field()
  userId: string;

  @Field(() => [String])
  recommendedTemplateIds: string[];

  @Field(() => Float)
  overallProgress: number;

  @Field(() => Int)
  completedTemplates: number;

  @Field(() => Int)
  totalTemplates: number;
}

@ObjectType()
export class UserProgressList {
  @Field(() => [UserProgress])
  progress: UserProgress[];

  @Field(() => Int)
  totalCount: number;
}

@ObjectType()
export class UserAchievementList {
  @Field(() => [UserAchievement])
  achievements: UserAchievement[];

  @Field(() => Int)
  totalCount: number;
}

@ObjectType()
export class UserProgressResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => UserProgress, { nullable: true })
  progress?: UserProgress;

  @Field(() => UserProgressList, { nullable: true })
  progressList?: UserProgressList;

  @Field(() => UserAchievement, { nullable: true })
  achievement?: UserAchievement;

  @Field(() => UserAchievementList, { nullable: true })
  achievements?: UserAchievementList;

  @Field(() => LearningPath, { nullable: true })
  learningPath?: LearningPath;

  @Field({ nullable: true })
  success?: boolean;
}

@InputType()
export class UpdateProgressInput {
  @Field()
  userId: string;

  @Field()
  templateId: string;

  @Field()
  sectionId: string;

  @Field(() => Float)
  progress: number;

  @Field()
  completed: boolean;
}

@InputType()
export class UnlockAchievementInput {
  @Field()
  userId: string;

  @Field()
  achievementId: string;
}