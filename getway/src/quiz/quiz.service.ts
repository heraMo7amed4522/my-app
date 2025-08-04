import { Injectable, OnModuleInit } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';

@Injectable()
export class QuizService implements OnModuleInit {
  private quizServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/quiz.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const quizProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.QUIZ_SERVICE_URL || 'localhost:50056';
    
    this.quizServiceClient = new quizProto.quiz.QuizService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  async getQuizzesBySection(sectionId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.quizServiceClient.GetQuizzesBySection(
        { section_id: sectionId },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async getQuizById(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.quizServiceClient.GetQuizByID(
        { id },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async createQuiz(quizData: any): Promise<any> {
    return new Promise((resolve, reject) => {
      this.quizServiceClient.CreateQuiz(
        { quiz: quizData },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async updateQuiz(id: string, quizData: any): Promise<any> {
    return new Promise((resolve, reject) => {
      this.quizServiceClient.UpdateQuiz(
        { quiz: { id, ...quizData } },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async deleteQuiz(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.quizServiceClient.DeleteQuiz(
        { id },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async submitQuizAnswer(userId: string, quizId: string, selectedOption: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.quizServiceClient.SubmitQuizAnswer(
        { user_id: userId, quiz_id: quizId, selected_option: selectedOption },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async getUserQuizHistory(userId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.quizServiceClient.GetUserQuizHistory(
        { user_id: userId },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }

  async getQuizStatistics(quizId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.quizServiceClient.GetQuizStatistics(
        { quiz_id: quizId },
        (error: any, response: any) => {
          if (error) {
            reject(error);
          } else {
            resolve(response);
          }
        },
      );
    });
  }
}