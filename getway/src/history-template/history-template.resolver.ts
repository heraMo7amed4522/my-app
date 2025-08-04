import { Resolver, Query, Mutation, Args, Context, Int } from '@nestjs/graphql';
import { HistoryTemplateService } from './history-template.service';
import {
  HistoryTemplate,
  HistoryTemplateResponse,
  CreateTemplateInput,
  UpdateTemplateInput,
} from './history-template.types';

@Resolver(() => HistoryTemplate)
export class HistoryTemplateResolver {
  constructor(private readonly historyTemplateService: HistoryTemplateService) {}

  @Query(() => HistoryTemplateResponse)
  async getTemplateById(@Args('id') id: string): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.getTemplateById(id);
  }

  @Query(() => HistoryTemplateResponse)
  async getAllTemplates(
    @Args('page', { type: () => Int, defaultValue: 1 }) page: number,
    @Args('limit', { type: () => Int, defaultValue: 10 }) limit: number,
    @Args('sortBy', { nullable: true }) sortBy?: string,
    @Args('order', { nullable: true }) order?: string,
  ): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.getAllTemplates(page, limit, sortBy, order);
  }

  @Query(() => HistoryTemplateResponse)
  async getTemplatesByEra(
    @Args('era') era: string,
    @Args('page', { type: () => Int, defaultValue: 1 }) page: number,
    @Args('limit', { type: () => Int, defaultValue: 10 }) limit: number,
  ): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.getTemplatesByEra(era, page, limit);
  }

  @Query(() => HistoryTemplateResponse)
  async getTemplatesByDynasty(
    @Args('dynasty', { type: () => Int }) dynasty: number,
    @Args('page', { type: () => Int, defaultValue: 1 }) page: number,
    @Args('limit', { type: () => Int, defaultValue: 10 }) limit: number,
  ): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.getTemplatesByDynasty(dynasty, page, limit);
  }

  @Query(() => HistoryTemplateResponse)
  async getTemplatesByPharaoh(
    @Args('pharaohId') pharaohId: string,
    @Args('page', { type: () => Int, defaultValue: 1 }) page: number,
    @Args('limit', { type: () => Int, defaultValue: 10 }) limit: number,
  ): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.getTemplatesByPharaoh(pharaohId, page, limit);
  }

  @Query(() => HistoryTemplateResponse)
  async getTemplatesByDifficulty(
    @Args('difficulty') difficulty: string,
    @Args('page', { type: () => Int, defaultValue: 1 }) page: number,
    @Args('limit', { type: () => Int, defaultValue: 10 }) limit: number,
  ): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.getTemplatesByDifficulty(difficulty, page, limit);
  }

  @Query(() => HistoryTemplateResponse)
  async searchTemplates(
    @Args('query') query: string,
    @Args('fields', { type: () => [String], nullable: true }) fields: string[],
    @Args('page', { type: () => Int, defaultValue: 1 }) page: number,
    @Args('limit', { type: () => Int, defaultValue: 10 }) limit: number,
  ): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.searchTemplates(query, fields || [], page, limit);
  }

  @Query(() => HistoryTemplateResponse)
  async getTemplatesByTag(
    @Args('tag') tag: string,
    @Args('page', { type: () => Int, defaultValue: 1 }) page: number,
    @Args('limit', { type: () => Int, defaultValue: 10 }) limit: number,
  ): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.getTemplatesByTag(tag, page, limit);
  }

  @Query(() => HistoryTemplateResponse)
  async getRelatedTemplates(
    @Args('templateId') templateId: string,
    @Args('limit', { type: () => Int, defaultValue: 5 }) limit: number,
  ): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.getRelatedTemplates(templateId, limit);
  }

  @Mutation(() => HistoryTemplateResponse)
  async createTemplate(
    @Args('input') input: CreateTemplateInput,
    @Context() context: any,
  ): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.createTemplate(input);
  }

  @Mutation(() => HistoryTemplateResponse)
  async updateTemplate(
    @Args('id') id: string,
    @Args('input') input: UpdateTemplateInput,
    @Context() context: any,
  ): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.updateTemplate(id, input);
  }

  @Mutation(() => HistoryTemplateResponse)
  async deleteTemplate(
    @Args('id') id: string,
    @Context() context: any,
  ): Promise<HistoryTemplateResponse> {
    return await this.historyTemplateService.deleteTemplate(id);
  }
}