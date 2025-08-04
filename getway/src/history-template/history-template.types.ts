import { ObjectType, Field, InputType, Int } from '@nestjs/graphql';

@ObjectType()
export class TemplateSection {
  @Field()
  id: string;

  @Field()
  templateId: string;

  @Field()
  title: string;

  @Field({ nullable: true })
  subtitle?: string;

  @Field()
  contentType: string;

  @Field()
  content: string;

  @Field(() => [String], { nullable: true })
  metadata?: string[];

  @Field(() => Int)
  orderIndex: number;

  @Field()
  optional: boolean;

  @Field()
  createdAt: string;
}

@ObjectType()
export class HistoryTemplate {
  @Field()
  id: string;

  @Field()
  title: string;

  @Field()
  description: string;

  @Field()
  era: string;

  @Field(() => Int)
  dynasty: number;

  @Field()
  pharaohId: string;

  @Field()
  difficulty: string;

  @Field({ nullable: true })
  thumbnailUrl?: string;

  @Field()
  language: string;

  @Field()
  isActive: boolean;

  @Field(() => Int)
  version: number;

  @Field({ nullable: true })
  publishedAt?: string;

  @Field()
  createdAt: string;

  @Field()
  updatedAt: string;

  @Field(() => [TemplateSection], { nullable: true })
  sections?: TemplateSection[];

  @Field(() => [String], { nullable: true })
  tags?: string[];
}

@ObjectType()
export class TemplateList {
  @Field(() => [HistoryTemplate])
  templates: HistoryTemplate[];

  @Field(() => Int)
  totalCount: number;

  @Field(() => Int)
  page: number;

  @Field(() => Int)
  limit: number;
}

@ObjectType()
export class HistoryTemplateResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => HistoryTemplate, { nullable: true })
  template?: HistoryTemplate;

  @Field(() => TemplateList, { nullable: true })
  templates?: TemplateList;

  @Field({ nullable: true })
  success?: boolean;
}

@InputType()
export class CreateTemplateInput {
  @Field()
  title: string;

  @Field()
  description: string;

  @Field()
  era: string;

  @Field(() => Int)
  dynasty: number;

  @Field()
  pharaohId: string;

  @Field()
  difficulty: string;

  @Field({ nullable: true })
  thumbnailUrl?: string;

  @Field()
  language: string;

  @Field({ nullable: true })
  isActive?: boolean;

  @Field(() => [String], { nullable: true })
  tags?: string[];
}

@InputType()
export class UpdateTemplateInput {
  @Field({ nullable: true })
  title?: string;

  @Field({ nullable: true })
  description?: string;

  @Field({ nullable: true })
  era?: string;

  @Field(() => Int, { nullable: true })
  dynasty?: number;

  @Field({ nullable: true })
  pharaohId?: string;

  @Field({ nullable: true })
  difficulty?: string;

  @Field({ nullable: true })
  thumbnailUrl?: string;

  @Field({ nullable: true })
  language?: string;

  @Field({ nullable: true })
  isActive?: boolean;

  @Field(() => [String], { nullable: true })
  tags?: string[];
}