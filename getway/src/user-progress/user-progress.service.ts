import { Injectable, OnModuleInit } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';

@Injectable()
export class UserProgressService implements OnModuleInit {
  private userProgressServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/user_progress.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const userProgressProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.USER_PROGRESS_SERVICE_URL || 'localhost:50060';
    
    this.userProgressServiceClient = new userProgressProto.user_progress.UserProgressService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  async getUserProgress(userId: string, templateId?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userProgressServiceClient.GetUserProgress(
        { user_id: userId, template_id: templateId || '' },
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

  async updateProgress(progressData: any): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userProgressServiceClient.UpdateProgress(
        progressData,
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

  async getCompletedTemplates(userId: string, page: number, limit: number): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userProgressServiceClient.GetCompletedTemplates(
        { user_id: userId, page, limit },
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

  async getUserAchievements(userId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userProgressServiceClient.GetUserAchievements(
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

  async unlockAchievement(userId: string, achievementId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userProgressServiceClient.UnlockAchievement(
        { user_id: userId, achievement_id: achievementId },
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

  async getLearningPath(userId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.userProgressServiceClient.GetLearningPath(
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
}