import { FeedbackService } from './feedback.service';
import { FeedbackResponse, SubmitFeedbackInput, UpdateFeedbackInput } from './feedback.types';
export declare class FeedbackResolver {
    private readonly feedbackService;
    constructor(feedbackService: FeedbackService);
    submitFeedback(input: SubmitFeedbackInput, context: any): Promise<FeedbackResponse>;
    getTemplateFeedback(templateId: string, page: number, limit: number): Promise<FeedbackResponse>;
    getFeedbackStatistics(templateId: string): Promise<FeedbackResponse>;
    updateFeedback(id: string, input: UpdateFeedbackInput, context: any): Promise<FeedbackResponse>;
    deleteFeedback(id: string, context: any): Promise<FeedbackResponse>;
}
