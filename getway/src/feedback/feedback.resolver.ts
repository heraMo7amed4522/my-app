import { Resolver, Query, Mutation, Args, Context, Int, Float } from '@nestjs/graphql';
import { FeedbackService } from './feedback.service';
import {
  Feedback,
  FeedbackResponse,
  SubmitFeedbackInput,
  UpdateFeedbackInput,
} from './feedback.types';

@Resolver(() => Feedback)
export class FeedbackResolver {
  constructor(private readonly feedbackService: FeedbackService) {}

  @Mutation(() => FeedbackResponse)
  async submitFeedback(
    @Args('input') input: SubmitFeedbackInput,
    @Context() context: any,
  ): Promise<FeedbackResponse> {
    try {
      const response = await this.feedbackService.submitFeedback({
        user_id: input.userId,
        template_id: input.templateId,
        rating: input.rating,
        comment: input.comment,
        language: input.language || 'en',
      });

      if (response.result?.feedback) {
        return {
          statusCode: response.status_code,
          message: response.message,
          feedback: {
            id: response.result.feedback.id,
            userId: response.result.feedback.user_id,
            templateId: response.result.feedback.template_id,
            rating: response.result.feedback.rating,
            comment: response.result.feedback.comment,
            language: response.result.feedback.language,
            createdAt: response.result.feedback.created_at?.seconds 
              ? new Date(response.result.feedback.created_at.seconds * 1000).toISOString()
              : new Date().toISOString(),
          },
        };
      } else {
        return {
          statusCode: response.status_code,
          message: response.message,
        };
      }
    } catch (error) {
      return {
        statusCode: 500,
        message: 'Failed to submit feedback',
      };
    }
  }

  @Query(() => FeedbackResponse)
  async getTemplateFeedback(
    @Args('templateId') templateId: string,
    @Args('page', { type: () => Int, defaultValue: 1 }) page: number,
    @Args('limit', { type: () => Int, defaultValue: 10 }) limit: number,
  ): Promise<FeedbackResponse> {
    try {
      const response = await this.feedbackService.getTemplateFeedback(templateId, page, limit);

      if (response.result?.feedback) {
        const feedbackList = {
          feedback: response.result.feedback.feedback.map((f: any) => ({
            id: f.id,
            userId: f.user_id,
            templateId: f.template_id,
            rating: f.rating,
            comment: f.comment,
            language: f.language,
            createdAt: f.created_at?.seconds 
              ? new Date(f.created_at.seconds * 1000).toISOString()
              : new Date().toISOString(),
          })),
          totalCount: response.result.feedback.total_count,
          page: response.result.feedback.page,
          limit: response.result.feedback.limit,
        };

        return {
          statusCode: response.status_code,
          message: response.message,
          feedbackList,
        };
      } else {
        return {
          statusCode: response.status_code,
          message: response.message,
        };
      }
    } catch (error) {
      return {
        statusCode: 500,
        message: 'Failed to get template feedback',
      };
    }
  }

  @Query(() => FeedbackResponse)
  async getFeedbackStatistics(
    @Args('templateId') templateId: string,
  ): Promise<FeedbackResponse> {
    try {
      const response = await this.feedbackService.getFeedbackStatistics(templateId);

      if (response.result?.statistics) {
        const statistics = {
          templateId: response.result.statistics.template_id,
          averageRating: response.result.statistics.average_rating,
          totalRatings: response.result.statistics.total_ratings,
          ratingDistribution: Object.entries(response.result.statistics.rating_distribution || {}).map(
            ([rating, count]) => ({
              rating,
              count: count as number,
            })
          ),
        };

        return {
          statusCode: response.status_code,
          message: response.message,
          statistics,
        };
      } else {
        return {
          statusCode: response.status_code,
          message: response.message,
        };
      }
    } catch (error) {
      return {
        statusCode: 500,
        message: 'Failed to get feedback statistics',
      };
    }
  }

  @Mutation(() => FeedbackResponse)
  async updateFeedback(
    @Args('id') id: string,
    @Args('input') input: UpdateFeedbackInput,
    @Context() context: any,
  ): Promise<FeedbackResponse> {
    try {
      const response = await this.feedbackService.updateFeedback(id, {
        rating: input.rating,
        comment: input.comment,
      });

      if (response.result?.feedback) {
        return {
          statusCode: response.status_code,
          message: response.message,
          feedback: {
            id: response.result.feedback.id,
            userId: response.result.feedback.user_id,
            templateId: response.result.feedback.template_id,
            rating: response.result.feedback.rating,
            comment: response.result.feedback.comment,
            language: response.result.feedback.language,
            createdAt: response.result.feedback.created_at?.seconds 
              ? new Date(response.result.feedback.created_at.seconds * 1000).toISOString()
              : new Date().toISOString(),
          },
        };
      } else {
        return {
          statusCode: response.status_code,
          message: response.message,
        };
      }
    } catch (error) {
      return {
        statusCode: 500,
        message: 'Failed to update feedback',
      };
    }
  }

  @Mutation(() => FeedbackResponse)
  async deleteFeedback(
    @Args('id') id: string,
    @Context() context: any,
  ): Promise<FeedbackResponse> {
    try {
      const response = await this.feedbackService.deleteFeedback(id);

      return {
        statusCode: response.status_code,
        message: response.message,
        success: response.result?.success || false,
      };
    } catch (error) {
      return {
        statusCode: 500,
        message: 'Failed to delete feedback',
        success: false,
      };
    }
  }
}