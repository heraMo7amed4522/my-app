export declare class TemplateSection {
    id: string;
    templateId: string;
    title: string;
    subtitle?: string;
    contentType: string;
    content: string;
    metadata?: string[];
    orderIndex: number;
    optional: boolean;
    createdAt: string;
}
export declare class HistoryTemplate {
    id: string;
    title: string;
    description: string;
    era: string;
    dynasty: number;
    pharaohId: string;
    difficulty: string;
    thumbnailUrl?: string;
    language: string;
    isActive: boolean;
    version: number;
    publishedAt?: string;
    createdAt: string;
    updatedAt: string;
    sections?: TemplateSection[];
    tags?: string[];
}
export declare class TemplateList {
    templates: HistoryTemplate[];
    totalCount: number;
    page: number;
    limit: number;
}
export declare class HistoryTemplateResponse {
    statusCode: number;
    message: string;
    template?: HistoryTemplate;
    templates?: TemplateList;
    success?: boolean;
}
export declare class CreateTemplateInput {
    title: string;
    description: string;
    era: string;
    dynasty: number;
    pharaohId: string;
    difficulty: string;
    thumbnailUrl?: string;
    language: string;
    isActive?: boolean;
    tags?: string[];
}
export declare class UpdateTemplateInput {
    title?: string;
    description?: string;
    era?: string;
    dynasty?: number;
    pharaohId?: string;
    difficulty?: string;
    thumbnailUrl?: string;
    language?: string;
    isActive?: boolean;
    tags?: string[];
}
