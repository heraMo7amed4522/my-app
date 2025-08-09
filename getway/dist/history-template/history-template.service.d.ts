import { OnModuleInit } from '@nestjs/common';
export declare class HistoryTemplateService implements OnModuleInit {
    private historyTemplateServiceClient;
    onModuleInit(): Promise<void>;
    getTemplateById(id: string): Promise<any>;
    getAllTemplates(page: number, limit: number, sortBy?: string, order?: string): Promise<any>;
    getTemplatesByEra(era: string, page: number, limit: number): Promise<any>;
    getTemplatesByDynasty(dynasty: number, page: number, limit: number): Promise<any>;
    getTemplatesByPharaoh(pharaohId: string, page: number, limit: number): Promise<any>;
    getTemplatesByDifficulty(difficulty: string, page: number, limit: number): Promise<any>;
    searchTemplates(query: string, fields: string[], page: number, limit: number): Promise<any>;
    createTemplate(templateData: any): Promise<any>;
    updateTemplate(id: string, templateData: any): Promise<any>;
    deleteTemplate(id: string): Promise<any>;
    getTemplatesByTag(tag: string, page: number, limit: number): Promise<any>;
    getRelatedTemplates(templateId: string, limit: number): Promise<any>;
}
