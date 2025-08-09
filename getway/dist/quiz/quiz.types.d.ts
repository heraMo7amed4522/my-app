export declare class Quiz {
    id: string;
    sectionId: string;
    question: string;
    options: QuizOption[];
    correctAnswer: string;
    explanation?: string;
    difficulty: string;
    createdAt: string;
}
export declare class QuizOption {
    key: string;
    value: string;
}
export declare class QuizResponse {
    id: string;
    userId: string;
    quizId: string;
    selectedOption: string;
    isCorrect: boolean;
    answeredAt: string;
}
export declare class QuizList {
    quizzes: Quiz[];
    totalCount: number;
    page: number;
    limit: number;
}
export declare class QuizStatistics {
    totalAttempts: number;
    correctAnswers: number;
    uniqueUsers: number;
    successRate: string;
}
export declare class QuizServiceResponse {
    statusCode: number;
    message: string;
    quiz?: Quiz;
    quizzes?: QuizList;
    quizResponses?: QuizResponse[];
    statistics?: QuizStatistics;
    success?: boolean;
}
export declare class CreateQuizInput {
    sectionId: string;
    question: string;
    options: QuizOptionInput[];
    correctAnswer: string;
    explanation?: string;
    difficulty: string;
}
export declare class QuizOptionInput {
    key: string;
    value: string;
}
export declare class UpdateQuizInput {
    sectionId?: string;
    question?: string;
    options?: QuizOptionInput[];
    correctAnswer?: string;
    explanation?: string;
    difficulty?: string;
}
export declare class SubmitQuizAnswerInput {
    userId: string;
    quizId: string;
    selectedOption: string;
}
