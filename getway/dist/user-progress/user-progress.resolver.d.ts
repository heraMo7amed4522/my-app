import { UserProgressService } from './user-progress.service';
import { UserProgressResponse, UpdateProgressInput, UnlockAchievementInput } from './user-progress.types';
export declare class UserProgressResolver {
    private readonly userProgressService;
    constructor(userProgressService: UserProgressService);
    getUserProgress(userId: string, templateId?: string): Promise<UserProgressResponse>;
    updateProgress(input: UpdateProgressInput, context: any): Promise<UserProgressResponse>;
    getCompletedTemplates(userId: string, page: number, limit: number): Promise<UserProgressResponse>;
    getUserAchievements(userId: string): Promise<UserProgressResponse>;
    unlockAchievement(input: UnlockAchievementInput, context: any): Promise<UserProgressResponse>;
    getLearningPath(userId: string): Promise<UserProgressResponse>;
}
