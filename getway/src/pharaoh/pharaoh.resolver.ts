import { Resolver, Query, Mutation, Args } from '@nestjs/graphql';
import { PharaohService } from './pharaoh.service';
import {
  Pharaoh,
  PharaohResponse,
  PharaohsResponse,
  DeletePharaohResponse,
  CreatePharaohInput,
  UpdatePharaohInput,
  GetAllPharaohsInput,
  GetPharaohsByDynastyInput,
  GetPharaohsByPeriodInput,
  SearchPharaohsInput,
  GetPharaohsByRarityInput,
  UpdatePharaohRatingInput,
} from './pharaoh.types';

@Resolver(() => Pharaoh)
export class PharaohResolver {
  constructor(private readonly pharaohService: PharaohService) {}

  @Query(() => PharaohResponse)
  async getPharaohById(
    @Args('id') id: string,
  ): Promise<PharaohResponse> {
    return await this.pharaohService.getPharaohById(id);
  }

  @Query(() => PharaohsResponse)
  async getAllPharaohs(
    @Args('input', { nullable: true }) input: GetAllPharaohsInput = {},
  ): Promise<PharaohsResponse> {
    return await this.pharaohService.getAllPharaohs(input);
  }

  @Query(() => PharaohsResponse)
  async getPharaohsByDynasty(
    @Args('input') input: GetPharaohsByDynastyInput,
  ): Promise<PharaohsResponse> {
    return await this.pharaohService.getPharaohsByDynasty(input);
  }

  @Query(() => PharaohsResponse)
  async getPharaohsByPeriod(
    @Args('input') input: GetPharaohsByPeriodInput,
  ): Promise<PharaohsResponse> {
    return await this.pharaohService.getPharaohsByPeriod(input);
  }

  @Query(() => PharaohsResponse)
  async searchPharaohs(
    @Args('input') input: SearchPharaohsInput,
  ): Promise<PharaohsResponse> {
    return await this.pharaohService.searchPharaohs(input);
  }

  @Query(() => PharaohsResponse)
  async getPharaohsByRarity(
    @Args('input') input: GetPharaohsByRarityInput,
  ): Promise<PharaohsResponse> {
    return await this.pharaohService.getPharaohsByRarity(input);
  }

  @Mutation(() => PharaohResponse)
  async createPharaoh(
    @Args('input') input: CreatePharaohInput,
  ): Promise<PharaohResponse> {
    return await this.pharaohService.createPharaoh(input);
  }

  @Mutation(() => PharaohResponse)
  async updatePharaoh(
    @Args('input') input: UpdatePharaohInput,
  ): Promise<PharaohResponse> {
    return await this.pharaohService.updatePharaoh(input);
  }

  @Mutation(() => DeletePharaohResponse)
  async deletePharaoh(
    @Args('id') id: string,
  ): Promise<DeletePharaohResponse> {
    return await this.pharaohService.deletePharaoh(id);
  }

  @Mutation(() => PharaohResponse)
  async updatePharaohRating(
    @Args('input') input: UpdatePharaohRatingInput,
  ): Promise<PharaohResponse> {
    return await this.pharaohService.updatePharaohRating(input);
  }
}