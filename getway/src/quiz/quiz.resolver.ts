import { Resolver, Query, Mutation, Args, Context } from '@nestjs/graphql';
import { QuizService } from './quiz.service';
import {
  Quiz,
  QuizServiceResponse,
  CreateQuizInput,
  UpdateQuizInput,
  SubmitQuizAnswerInput,
} from './quiz.types';

@Resolver(() => Quiz)
export class QuizResolver {
  constructor(private readonly quizService: QuizService) {}

  @Query(() => QuizServiceResponse)
  async getQuizzesBySection(@Args('sectionId') sectionId: string): Promise<QuizServiceResponse> {
    return await this.quizService.getQuizzesBySection(sectionId);
  }

  @Query(() => QuizServiceResponse)
  async getQuizById(@Args('id') id: string): Promise<QuizServiceResponse> {
    return await this.quizService.getQuizById(id);
  }

  @Query(() => QuizServiceResponse)
  async getUserQuizHistory(@Args('userId') userId: string): Promise<QuizServiceResponse> {
    return await this.quizService.getUserQuizHistory(userId);
  }

  @Query(() => QuizServiceResponse)
  async getQuizStatistics(@Args('quizId') quizId: string): Promise<QuizServiceResponse> {
    return await this.quizService.getQuizStatistics(quizId);
  }

  @Mutation(() => QuizServiceResponse)
  async createQuiz(
    @Args('input') input: CreateQuizInput,
    @Context() context: any,
  ): Promise<QuizServiceResponse> {
    // Convert options array to map format expected by gRPC
    const optionsMap: { [key: string]: string } = {};
    input.options.forEach(option => {
      optionsMap[option.key] = option.value;
    });

    const quizData = {
      section_id: input.sectionId,
      question: input.question,
      options: optionsMap,
      correct_answer: input.correctAnswer,
      explanation: input.explanation,
      difficulty: input.difficulty,
    };

    return await this.quizService.createQuiz(quizData);
  }

  @Mutation(() => QuizServiceResponse)
  async updateQuiz(
    @Args('id') id: string,
    @Args('input') input: UpdateQuizInput,
    @Context() context: any,
  ): Promise<QuizServiceResponse> {
    // Convert options array to map format if provided
    const updateData: any = {
      section_id: input.sectionId,
      question: input.question,
      correct_answer: input.correctAnswer,
      explanation: input.explanation,
      difficulty: input.difficulty,
    };

    if (input.options) {
      const optionsMap: { [key: string]: string } = {};
      input.options.forEach(option => {
        optionsMap[option.key] = option.value;
      });
      updateData.options = optionsMap;
    }

    return await this.quizService.updateQuiz(id, updateData);
  }

  @Mutation(() => QuizServiceResponse)
  async deleteQuiz(
    @Args('id') id: string,
    @Context() context: any,
  ): Promise<QuizServiceResponse> {
    return await this.quizService.deleteQuiz(id);
  }

  @Mutation(() => QuizServiceResponse)
  async submitQuizAnswer(
    @Args('input') input: SubmitQuizAnswerInput,
    @Context() context: any,
  ): Promise<QuizServiceResponse> {
    return await this.quizService.submitQuizAnswer(
      input.userId,
      input.quizId,
      input.selectedOption,
    );
  }
}