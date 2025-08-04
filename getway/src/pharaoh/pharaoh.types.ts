import { ObjectType, Field, InputType, Int, Float } from '@nestjs/graphql';
import { ErrorDetails } from '../shared/shared.types';

@ObjectType()
export class Artifact {
  @Field()
  name: string;

  @Field()
  museum: string;

  @Field()
  description: string;
}

@ObjectType()
export class Traits {
  @Field(() => Int)
  leadership: number;

  @Field(() => Int)
  military: number;

  @Field(() => Int)
  diplomacy: number;

  @Field(() => Int)
  wisdom: number;

  @Field(() => Int)
  charisma: number;
}

@ObjectType()
export class Pharaoh {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field({ nullable: true })
  birthName?: string;

  @Field({ nullable: true })
  throneName?: string;

  @Field({ nullable: true })
  epithet?: string;

  @Field(() => Int, { nullable: true })
  dynasty?: number;

  @Field({ nullable: true })
  period?: string;

  @Field(() => Int, { nullable: true })
  reignStart?: number;

  @Field(() => Int, { nullable: true })
  reignEnd?: number;

  @Field(() => Int, { nullable: true })
  lengthOfReignYears?: number;

  @Field({ nullable: true })
  predecessorId?: string;

  @Field({ nullable: true })
  successorId?: string;

  @Field({ nullable: true })
  father?: string;

  @Field({ nullable: true })
  mother?: string;

  @Field(() => [String], { nullable: true })
  consorts?: string[];

  @Field(() => Int, { nullable: true })
  childrenCount?: number;

  @Field(() => [String], { nullable: true })
  notableChildren?: string[];

  @Field({ nullable: true })
  capital?: string;

  @Field(() => [String], { nullable: true })
  majorAchievements?: string[];

  @Field(() => [String], { nullable: true })
  militaryCampaigns?: string[];

  @Field(() => [String], { nullable: true })
  buildingProjects?: string[];

  @Field({ nullable: true })
  politicalStyle?: string;

  @Field(() => [String], { nullable: true })
  divineAssociation?: string[];

  @Field(() => [String], { nullable: true })
  templeAffiliations?: string[];

  @Field({ nullable: true })
  religiousReforms?: string;

  @Field({ nullable: true })
  pharaohAsGod?: boolean;

  @Field({ nullable: true })
  burialSite?: string;

  @Field({ nullable: true })
  tombDiscovered?: boolean;

  @Field(() => Int, { nullable: true })
  discoveryYear?: number;

  @Field({ nullable: true })
  tombGuardian?: string;

  @Field({ nullable: true })
  funeraryText?: string;

  @Field(() => [Artifact], { nullable: true })
  famousArtifacts?: Artifact[];

  @Field({ nullable: true })
  treasureStatus?: string;

  @Field({ nullable: true })
  imageUrl?: string;

  @Field(() => Int, { nullable: true })
  statueCount?: number;

  @Field({ nullable: true })
  mummyLocation?: string;

  @Field({ nullable: true })
  audioNarrationUrl?: string;

  @Field({ nullable: true })
  videoDocumentaryUrl?: string;

  @Field(() => Float, { nullable: true })
  popularityScore?: number;

  @Field(() => Float, { nullable: true })
  userRating?: number;

  @Field({ nullable: true })
  unlockInGame?: boolean;

  @Field({ nullable: true })
  rarity?: string;

  @Field(() => Traits, { nullable: true })
  traits?: Traits;

  @Field({ nullable: true })
  source?: string;

  @Field({ nullable: true })
  verified?: boolean;

  @Field({ nullable: true })
  language?: string;

  @Field()
  createdAt: string;

  @Field()
  updatedAt: string;
}

@ObjectType()
export class PharaohResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => Pharaoh, { nullable: true })
  pharaoh?: Pharaoh;

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;
}

@ObjectType()
export class PharaohsResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => [Pharaoh], { nullable: true })
  pharaohs?: Pharaoh[];

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;

  @Field(() => Int, { nullable: true })
  totalCount?: number;

  @Field(() => Int, { nullable: true })
  page?: number;

  @Field(() => Int, { nullable: true })
  limit?: number;
}

@ObjectType()
export class DeletePharaohResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field({ nullable: true })
  success?: boolean;

  @Field(() => ErrorDetails, { nullable: true })
  error?: ErrorDetails;
}

// Input Types
@InputType()
export class ArtifactInput {
  @Field()
  name: string;

  @Field()
  museum: string;

  @Field()
  description: string;
}

@InputType()
export class TraitsInput {
  @Field(() => Int)
  leadership: number;

  @Field(() => Int)
  military: number;

  @Field(() => Int)
  diplomacy: number;

  @Field(() => Int)
  wisdom: number;

  @Field(() => Int)
  charisma: number;
}

@InputType()
export class CreatePharaohInput {
  @Field()
  name: string;

  @Field({ nullable: true })
  birthName?: string;

  @Field({ nullable: true })
  throneName?: string;

  @Field({ nullable: true })
  epithet?: string;

  @Field(() => Int, { nullable: true })
  dynasty?: number;

  @Field({ nullable: true })
  period?: string;

  @Field(() => Int, { nullable: true })
  reignStart?: number;

  @Field(() => Int, { nullable: true })
  reignEnd?: number;

  @Field(() => Int, { nullable: true })
  lengthOfReignYears?: number;

  @Field({ nullable: true })
  predecessorId?: string;

  @Field({ nullable: true })
  successorId?: string;

  @Field({ nullable: true })
  father?: string;

  @Field({ nullable: true })
  mother?: string;

  @Field(() => [String], { nullable: true })
  consorts?: string[];

  @Field(() => Int, { nullable: true })
  childrenCount?: number;

  @Field(() => [String], { nullable: true })
  notableChildren?: string[];

  @Field({ nullable: true })
  capital?: string;

  @Field(() => [String], { nullable: true })
  majorAchievements?: string[];

  @Field(() => [String], { nullable: true })
  militaryCampaigns?: string[];

  @Field(() => [String], { nullable: true })
  buildingProjects?: string[];

  @Field({ nullable: true })
  politicalStyle?: string;

  @Field(() => [String], { nullable: true })
  divineAssociation?: string[];

  @Field(() => [String], { nullable: true })
  templeAffiliations?: string[];

  @Field({ nullable: true })
  religiousReforms?: string;

  @Field({ nullable: true })
  pharaohAsGod?: boolean;

  @Field({ nullable: true })
  burialSite?: string;

  @Field({ nullable: true })
  tombDiscovered?: boolean;

  @Field(() => Int, { nullable: true })
  discoveryYear?: number;

  @Field({ nullable: true })
  tombGuardian?: string;

  @Field({ nullable: true })
  funeraryText?: string;

  @Field(() => [ArtifactInput], { nullable: true })
  famousArtifacts?: ArtifactInput[];

  @Field({ nullable: true })
  treasureStatus?: string;

  @Field({ nullable: true })
  imageUrl?: string;

  @Field(() => Int, { nullable: true })
  statueCount?: number;

  @Field({ nullable: true })
  mummyLocation?: string;

  @Field({ nullable: true })
  audioNarrationUrl?: string;

  @Field({ nullable: true })
  videoDocumentaryUrl?: string;

  @Field(() => Float, { nullable: true })
  popularityScore?: number;

  @Field(() => Float, { nullable: true })
  userRating?: number;

  @Field({ nullable: true })
  unlockInGame?: boolean;

  @Field({ nullable: true })
  rarity?: string;

  @Field(() => TraitsInput, { nullable: true })
  traits?: TraitsInput;

  @Field({ nullable: true })
  source?: string;

  @Field({ nullable: true })
  verified?: boolean;

  @Field({ nullable: true })
  language?: string;
}

@InputType()
export class UpdatePharaohInput {
  @Field()
  id: string;

  @Field(() => CreatePharaohInput)
  pharaoh: CreatePharaohInput;
}

@InputType()
export class GetAllPharaohsInput {
  @Field(() => Int, { nullable: true })
  page?: number;

  @Field(() => Int, { nullable: true })
  limit?: number;

  @Field({ nullable: true })
  sortBy?: string;

  @Field({ nullable: true })
  order?: string;
}

@InputType()
export class GetPharaohsByDynastyInput {
  @Field(() => Int)
  dynasty: number;
}

@InputType()
export class GetPharaohsByPeriodInput {
  @Field()
  period: string;
}

@InputType()
export class SearchPharaohsInput {
  @Field()
  query: string;

  @Field(() => [String], { nullable: true })
  fields?: string[];
}

@InputType()
export class GetPharaohsByRarityInput {
  @Field()
  rarity: string;
}

@InputType()
export class UpdatePharaohRatingInput {
  @Field()
  id: string;

  @Field(() => Float)
  rating: number;
}