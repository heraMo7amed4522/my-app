export declare class UserProgress {
    id: string;
    userId: string;
    templateId: string;
    sectionId: string;
    progress: number;
    completed: boolean;
    lastViewed: string;
    createdAt: string;
}
export declare class Achievement {
    id: string;
    code: string;
    name: string;
    description: string;
    badgeUrl?: string;
    createdAt: string;
}
export declare class UserAchievement {
    id: string;
    userId: string;
    achievementId: string;
    achievement: Achievement;
    unlockedAt: string;
}
export declare class LearningPath {
    userId: string;
    recommendedTemplateIds: string[];
    overallProgress: number;
    completedTemplates: number;
    totalTemplates: number;
}
export declare class UserProgressList {
    progress: UserProgress[];
    totalCount: number;
}
export declare class UserAchievementList {
    achievements: UserAchievement[];
    totalCount: number;
}
export declare class UserProgressResponse {
    statusCode: number;
    message: string;
    progress?: UserProgress;
    progressList?: UserProgressList;
    achievement?: UserAchievement;
    achievements?: UserAchievementList;
    learningPath?: LearningPath;
    success?: boolean;
}
export declare class UpdateProgressInput {
    userId: string;
    templateId: string;
    sectionId: string;
    progress: number;
    completed: boolean;
}
export declare class UnlockAchievementInput {
    userId: string;
    achievementId: string;
}
