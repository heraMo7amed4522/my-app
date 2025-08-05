import { Resolver, Query, Mutation, Args, Context, Int } from '@nestjs/graphql';
import { UserProgressService } from './user-progress.service';
import {
  UserProgress,
  UserProgressResponse,
  UpdateProgressInput,
  UnlockAchievementInput,
  UserProgressList,
  UserAchievementList,
  LearningPath,
} from './user-progress.types';

@Resolver(() => UserProgress)
export class UserProgressResolver {
  constructor(private readonly userProgressService: UserProgressService) {}

  @Query(() => UserProgressResponse)
  async getUserProgress(
    @Args('userId') userId: string,
    @Args('templateId', { nullable: true }) templateId?: string,
  ): Promise<UserProgressResponse> {
    try {
      const response = await this.userProgressService.getUserProgress(userId, templateId);

      if (response.result?.progress) {
        const progressList: UserProgressList = {
          progress: response.result.progress.progress.map((p: any) => ({
            id: p.id,
            userId: p.user_id,
            templateId: p.template_id,
            sectionId: p.section_id,
            progress: p.progress,
            completed: p.completed,
            lastViewed: p.last_viewed?.seconds 
              ? new Date(p.last_viewed.seconds * 1000).toISOString()
              : new Date().toISOString(),
            createdAt: p.created_at?.seconds 
              ? new Date(p.created_at.seconds * 1000).toISOString()
              : new Date().toISOString(),
          })),
          totalCount: response.result.progress.total_count,
        };

        return {
          statusCode: response.status_code,
          message: response.message,
          progressList,
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
        message: 'Failed to retrieve user progress',
        success: false,
      };
    }
  }

  @Mutation(() => UserProgressResponse)
  async updateProgress(
    @Args('input') input: UpdateProgressInput,
    @Context() context: any,
  ): Promise<UserProgressResponse> {
    try {
      const response = await this.userProgressService.updateProgress({
        user_id: input.userId,
        template_id: input.templateId,
        section_id: input.sectionId,
        progress: input.progress,
        completed: input.completed,
      });

      if (response.result?.progress) {
        return {
          statusCode: response.status_code,
          message: response.message,
          progress: {
            id: response.result.progress.id,
            userId: response.result.progress.user_id,
            templateId: response.result.progress.template_id,
            sectionId: response.result.progress.section_id,
            progress: response.result.progress.progress,
            completed: response.result.progress.completed,
            lastViewed: response.result.progress.last_viewed?.seconds 
              ? new Date(response.result.progress.last_viewed.seconds * 1000).toISOString()
              : new Date().toISOString(),
            createdAt: response.result.progress.created_at?.seconds 
              ? new Date(response.result.progress.created_at.seconds * 1000).toISOString()
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
        message: 'Failed to update progress',
        success: false,
      };
    }
  }

  @Query(() => UserProgressResponse)
  async getCompletedTemplates(
    @Args('userId') userId: string,
    @Args('page', { type: () => Int, defaultValue: 1 }) page: number,
    @Args('limit', { type: () => Int, defaultValue: 10 }) limit: number,
  ): Promise<UserProgressResponse> {
    try {
      const response = await this.userProgressService.getCompletedTemplates(userId, page, limit);

      if (response.result?.completed) {
        const progressList: UserProgressList = {
          progress: response.result.completed.progress.map((p: any) => ({
            id: p.id,
            userId: p.user_id,
            templateId: p.template_id,
            sectionId: p.section_id,
            progress: p.progress,
            completed: p.completed,
            lastViewed: p.last_viewed?.seconds 
              ? new Date(p.last_viewed.seconds * 1000).toISOString()
              : new Date().toISOString(),
            createdAt: p.created_at?.seconds 
              ? new Date(p.created_at.seconds * 1000).toISOString()
              : new Date().toISOString(),
          })),
          totalCount: response.result.completed.total_count,
        };

        return {
          statusCode: response.status_code,
          message: response.message,
          progressList,
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
        message: 'Failed to retrieve completed templates',
        success: false,
      };
    }
  }

  @Query(() => UserProgressResponse)
  async getUserAchievements(
    @Args('userId') userId: string,
  ): Promise<UserProgressResponse> {
    try {
      const response = await this.userProgressService.getUserAchievements(userId);

      if (response.result?.achievements) {
        const achievements: UserAchievementList = {
          achievements: response.result.achievements.achievements.map((a: any) => ({
            id: a.id,
            userId: a.user_id,
            achievementId: a.achievement_id,
            achievement: {
              id: a.achievement.id,
              code: a.achievement.code,
              name: a.achievement.name,
              description: a.achievement.description,
              badgeUrl: a.achievement.badge_url,
              createdAt: a.achievement.created_at?.seconds 
                ? new Date(a.achievement.created_at.seconds * 1000).toISOString()
                : new Date().toISOString(),
            },
            unlockedAt: a.unlocked_at?.seconds 
              ? new Date(a.unlocked_at.seconds * 1000).toISOString()
              : new Date().toISOString(),
          })),
          totalCount: response.result.achievements.total_count,
        };

        return {
          statusCode: response.status_code,
          message: response.message,
          achievements,
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
        message: 'Failed to retrieve user achievements',
        success: false,
      };
    }
  }

  @Mutation(() => UserProgressResponse)
  async unlockAchievement(
    @Args('input') input: UnlockAchievementInput,
    @Context() context: any,
  ): Promise<UserProgressResponse> {
    try {
      const response = await this.userProgressService.unlockAchievement(
        input.userId,
        input.achievementId,
      );

      if (response.result?.achievement) {
        return {
          statusCode: response.status_code,
          message: response.message,
          achievement: {
            id: response.result.achievement.id,
            userId: response.result.achievement.user_id,
            achievementId: response.result.achievement.achievement_id,
            achievement: {
              id: response.result.achievement.achievement.id,
              code: response.result.achievement.achievement.code,
              name: response.result.achievement.achievement.name,
              description: response.result.achievement.achievement.description,
              badgeUrl: response.result.achievement.achievement.badge_url,
              createdAt: response.result.achievement.achievement.created_at?.seconds 
                ? new Date(response.result.achievement.achievement.created_at.seconds * 1000).toISOString()
                : new Date().toISOString(),
            },
            unlockedAt: response.result.achievement.unlocked_at?.seconds 
              ? new Date(response.result.achievement.unlocked_at.seconds * 1000).toISOString()
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
        message: 'Failed to unlock achievement',
        success: false,
      };
    }
  }

  @Query(() => UserProgressResponse)
  async getLearningPath(
    @Args('userId') userId: string,
  ): Promise<UserProgressResponse> {
    try {
      const response = await this.userProgressService.getLearningPath(userId);

      if (response.result?.path) {
        const learningPath: LearningPath = {
          userId: response.result.path.user_id,
          recommendedTemplateIds: response.result.path.recommended_template_ids,
          overallProgress: response.result.path.overall_progress,
          completedTemplates: response.result.path.completed_templates,
          totalTemplates: response.result.path.total_templates,
        };

        return {
          statusCode: response.status_code,
          message: response.message,
          learningPath,
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
        message: 'Failed to generate learning path',
        success: false,
      };
    }
  }
}