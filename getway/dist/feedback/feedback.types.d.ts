export declare class Feedback {
    id: string;
    userId: string;
    templateId: string;
    rating: number;
    comment: string;
    language: string;
    createdAt: string;
}
export declare class FeedbackList {
    feedback: Feedback[];
    totalCount: number;
    page: number;
    limit: number;
}
export declare class FeedbackStatistics {
    templateId: string;
    averageRating: number;
    totalRatings: number;
    ratingDistribution: RatingDistribution[];
}
export declare class RatingDistribution {
    rating: string;
    count: number;
}
export declare class FeedbackResponse {
    statusCode: number;
    message: string;
    feedback?: Feedback;
    feedbackList?: FeedbackList;
    statistics?: FeedbackStatistics;
    success?: boolean;
}
export declare class SubmitFeedbackInput {
    userId: string;
    templateId: string;
    rating: number;
    comment: string;
    language?: string;
}
export declare class UpdateFeedbackInput {
    rating: number;
    comment: string;
}
