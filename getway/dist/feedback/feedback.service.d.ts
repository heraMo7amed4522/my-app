import { OnModuleInit } from '@nestjs/common';
export declare class FeedbackService implements OnModuleInit {
    private feedbackServiceClient;
    onModuleInit(): Promise<void>;
    submitFeedback(feedbackData: any): Promise<any>;
    getTemplateFeedback(templateId: string, page: number, limit: number): Promise<any>;
    getFeedbackStatistics(templateId: string): Promise<any>;
    updateFeedback(id: string, feedbackData: any): Promise<any>;
    deleteFeedback(id: string): Promise<any>;
}
