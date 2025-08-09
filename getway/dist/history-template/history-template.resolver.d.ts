import { HistoryTemplateService } from './history-template.service';
import { HistoryTemplateResponse, CreateTemplateInput, UpdateTemplateInput } from './history-template.types';
export declare class HistoryTemplateResolver {
    private readonly historyTemplateService;
    constructor(historyTemplateService: HistoryTemplateService);
    getTemplateById(id: string): Promise<HistoryTemplateResponse>;
    getAllTemplates(page: number, limit: number, sortBy?: string, order?: string): Promise<HistoryTemplateResponse>;
    getTemplatesByEra(era: string, page: number, limit: number): Promise<HistoryTemplateResponse>;
    getTemplatesByDynasty(dynasty: number, page: number, limit: number): Promise<HistoryTemplateResponse>;
    getTemplatesByPharaoh(pharaohId: string, page: number, limit: number): Promise<HistoryTemplateResponse>;
    getTemplatesByDifficulty(difficulty: string, page: number, limit: number): Promise<HistoryTemplateResponse>;
    searchTemplates(query: string, fields: string[], page: number, limit: number): Promise<HistoryTemplateResponse>;
    getTemplatesByTag(tag: string, page: number, limit: number): Promise<HistoryTemplateResponse>;
    getRelatedTemplates(templateId: string, limit: number): Promise<HistoryTemplateResponse>;
    createTemplate(input: CreateTemplateInput, context: any): Promise<HistoryTemplateResponse>;
    updateTemplate(id: string, input: UpdateTemplateInput, context: any): Promise<HistoryTemplateResponse>;
    deleteTemplate(id: string, context: any): Promise<HistoryTemplateResponse>;
}
