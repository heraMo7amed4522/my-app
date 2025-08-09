import { OnModuleInit } from '@nestjs/common';
import { CreatePharaohInput, UpdatePharaohInput, GetAllPharaohsInput, GetPharaohsByDynastyInput, GetPharaohsByPeriodInput, SearchPharaohsInput, GetPharaohsByRarityInput, UpdatePharaohRatingInput } from './pharaoh.types';
export declare class PharaohService implements OnModuleInit {
    private pharaohServiceClient;
    onModuleInit(): Promise<void>;
    getPharaohById(id: string): Promise<any>;
    getAllPharaohs(input: GetAllPharaohsInput): Promise<any>;
    getPharaohsByDynasty(input: GetPharaohsByDynastyInput): Promise<any>;
    getPharaohsByPeriod(input: GetPharaohsByPeriodInput): Promise<any>;
    searchPharaohs(input: SearchPharaohsInput): Promise<any>;
    createPharaoh(input: CreatePharaohInput): Promise<any>;
    updatePharaoh(input: UpdatePharaohInput): Promise<any>;
    deletePharaoh(id: string): Promise<any>;
    getPharaohsByRarity(input: GetPharaohsByRarityInput): Promise<any>;
    updatePharaohRating(input: UpdatePharaohRatingInput): Promise<any>;
}
