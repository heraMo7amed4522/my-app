import { OnModuleInit } from '@nestjs/common';
export declare class UserProgressService implements OnModuleInit {
    private userProgressServiceClient;
    onModuleInit(): Promise<void>;
    getUserProgress(userId: string, templateId?: string): Promise<any>;
    updateProgress(progressData: any): Promise<any>;
    getCompletedTemplates(userId: string, page: number, limit: number): Promise<any>;
    getUserAchievements(userId: string): Promise<any>;
    unlockAchievement(userId: string, achievementId: string): Promise<any>;
    getLearningPath(userId: string): Promise<any>;
}
