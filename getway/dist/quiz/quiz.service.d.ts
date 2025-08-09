import { OnModuleInit } from '@nestjs/common';
export declare class QuizService implements OnModuleInit {
    private quizServiceClient;
    onModuleInit(): Promise<void>;
    getQuizzesBySection(sectionId: string): Promise<any>;
    getQuizById(id: string): Promise<any>;
    createQuiz(quizData: any): Promise<any>;
    updateQuiz(id: string, quizData: any): Promise<any>;
    deleteQuiz(id: string): Promise<any>;
    submitQuizAnswer(userId: string, quizId: string, selectedOption: string): Promise<any>;
    getUserQuizHistory(userId: string): Promise<any>;
    getQuizStatistics(quizId: string): Promise<any>;
}
