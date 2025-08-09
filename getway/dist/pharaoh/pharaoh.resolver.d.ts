import { PharaohService } from './pharaoh.service';
import { PharaohResponse, PharaohsResponse, DeletePharaohResponse, CreatePharaohInput, UpdatePharaohInput, GetAllPharaohsInput, GetPharaohsByDynastyInput, GetPharaohsByPeriodInput, SearchPharaohsInput, GetPharaohsByRarityInput, UpdatePharaohRatingInput } from './pharaoh.types';
export declare class PharaohResolver {
    private readonly pharaohService;
    constructor(pharaohService: PharaohService);
    getPharaohById(id: string): Promise<PharaohResponse>;
    getAllPharaohs(input?: GetAllPharaohsInput): Promise<PharaohsResponse>;
    getPharaohsByDynasty(input: GetPharaohsByDynastyInput): Promise<PharaohsResponse>;
    getPharaohsByPeriod(input: GetPharaohsByPeriodInput): Promise<PharaohsResponse>;
    searchPharaohs(input: SearchPharaohsInput): Promise<PharaohsResponse>;
    getPharaohsByRarity(input: GetPharaohsByRarityInput): Promise<PharaohsResponse>;
    createPharaoh(input: CreatePharaohInput): Promise<PharaohResponse>;
    updatePharaoh(input: UpdatePharaohInput): Promise<PharaohResponse>;
    deletePharaoh(id: string): Promise<DeletePharaohResponse>;
    updatePharaohRating(input: UpdatePharaohRatingInput): Promise<PharaohResponse>;
}
