import { ErrorDetails } from '../shared/shared.types';
export declare class Artifact {
    name: string;
    museum: string;
    description: string;
}
export declare class Traits {
    leadership: number;
    military: number;
    diplomacy: number;
    wisdom: number;
    charisma: number;
}
export declare class Pharaoh {
    id: string;
    name: string;
    birthName?: string;
    throneName?: string;
    epithet?: string;
    dynasty?: number;
    period?: string;
    reignStart?: number;
    reignEnd?: number;
    lengthOfReignYears?: number;
    predecessorId?: string;
    successorId?: string;
    father?: string;
    mother?: string;
    consorts?: string[];
    childrenCount?: number;
    notableChildren?: string[];
    capital?: string;
    majorAchievements?: string[];
    militaryCampaigns?: string[];
    buildingProjects?: string[];
    politicalStyle?: string;
    divineAssociation?: string[];
    templeAffiliations?: string[];
    religiousReforms?: string;
    pharaohAsGod?: boolean;
    burialSite?: string;
    tombDiscovered?: boolean;
    discoveryYear?: number;
    tombGuardian?: string;
    funeraryText?: string;
    famousArtifacts?: Artifact[];
    treasureStatus?: string;
    imageUrl?: string;
    statueCount?: number;
    mummyLocation?: string;
    audioNarrationUrl?: string;
    videoDocumentaryUrl?: string;
    popularityScore?: number;
    userRating?: number;
    unlockInGame?: boolean;
    rarity?: string;
    traits?: Traits;
    source?: string;
    verified?: boolean;
    language?: string;
    createdAt: string;
    updatedAt: string;
}
export declare class PharaohResponse {
    statusCode: number;
    message: string;
    pharaoh?: Pharaoh;
    error?: ErrorDetails;
}
export declare class PharaohsResponse {
    statusCode: number;
    message: string;
    pharaohs?: Pharaoh[];
    error?: ErrorDetails;
    totalCount?: number;
    page?: number;
    limit?: number;
}
export declare class DeletePharaohResponse {
    statusCode: number;
    message: string;
    success?: boolean;
    error?: ErrorDetails;
}
export declare class ArtifactInput {
    name: string;
    museum: string;
    description: string;
}
export declare class TraitsInput {
    leadership: number;
    military: number;
    diplomacy: number;
    wisdom: number;
    charisma: number;
}
export declare class CreatePharaohInput {
    name: string;
    birthName?: string;
    throneName?: string;
    epithet?: string;
    dynasty?: number;
    period?: string;
    reignStart?: number;
    reignEnd?: number;
    lengthOfReignYears?: number;
    predecessorId?: string;
    successorId?: string;
    father?: string;
    mother?: string;
    consorts?: string[];
    childrenCount?: number;
    notableChildren?: string[];
    capital?: string;
    majorAchievements?: string[];
    militaryCampaigns?: string[];
    buildingProjects?: string[];
    politicalStyle?: string;
    divineAssociation?: string[];
    templeAffiliations?: string[];
    religiousReforms?: string;
    pharaohAsGod?: boolean;
    burialSite?: string;
    tombDiscovered?: boolean;
    discoveryYear?: number;
    tombGuardian?: string;
    funeraryText?: string;
    famousArtifacts?: ArtifactInput[];
    treasureStatus?: string;
    imageUrl?: string;
    statueCount?: number;
    mummyLocation?: string;
    audioNarrationUrl?: string;
    videoDocumentaryUrl?: string;
    popularityScore?: number;
    userRating?: number;
    unlockInGame?: boolean;
    rarity?: string;
    traits?: TraitsInput;
    source?: string;
    verified?: boolean;
    language?: string;
}
export declare class UpdatePharaohInput {
    id: string;
    pharaoh: CreatePharaohInput;
}
export declare class GetAllPharaohsInput {
    page?: number;
    limit?: number;
    sortBy?: string;
    order?: string;
}
export declare class GetPharaohsByDynastyInput {
    dynasty: number;
}
export declare class GetPharaohsByPeriodInput {
    period: string;
}
export declare class SearchPharaohsInput {
    query: string;
    fields?: string[];
}
export declare class GetPharaohsByRarityInput {
    rarity: string;
}
export declare class UpdatePharaohRatingInput {
    id: string;
    rating: number;
}
