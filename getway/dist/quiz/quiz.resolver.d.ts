import { QuizService } from './quiz.service';
import { QuizServiceResponse, CreateQuizInput, UpdateQuizInput, SubmitQuizAnswerInput } from './quiz.types';
export declare class QuizResolver {
    private readonly quizService;
    constructor(quizService: QuizService);
    getQuizzesBySection(sectionId: string): Promise<QuizServiceResponse>;
    getQuizById(id: string): Promise<QuizServiceResponse>;
    getUserQuizHistory(userId: string): Promise<QuizServiceResponse>;
    getQuizStatistics(quizId: string): Promise<QuizServiceResponse>;
    createQuiz(input: CreateQuizInput, context: any): Promise<QuizServiceResponse>;
    updateQuiz(id: string, input: UpdateQuizInput, context: any): Promise<QuizServiceResponse>;
    deleteQuiz(id: string, context: any): Promise<QuizServiceResponse>;
    submitQuizAnswer(input: SubmitQuizAnswerInput, context: any): Promise<QuizServiceResponse>;
}
