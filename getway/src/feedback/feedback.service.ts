import { Injectable, OnModuleInit } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';

@Injectable()
export class FeedbackService implements OnModuleInit {
  private feedbackServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/feedback.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const feedbackProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.FEEDBACK_SERVICE_URL || 'localhost:50059';
    
    this.feedbackServiceClient = new feedbackProto.feedback.FeedbackService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  async submitFeedback(feedbackData: any): Promise<any> {
    return new Promise((resolve, reject) => {
      this.feedbackServiceClient.SubmitFeedback(
        feedbackData,
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

  async getTemplateFeedback(templateId: string, page: number, limit: number): Promise<any> {
    return new Promise((resolve, reject) => {
      this.feedbackServiceClient.GetTemplateFeedback(
        { template_id: templateId, page, limit },
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

  async getFeedbackStatistics(templateId: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.feedbackServiceClient.GetFeedbackStatistics(
        { template_id: templateId },
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

  async updateFeedback(id: string, feedbackData: any): Promise<any> {
    return new Promise((resolve, reject) => {
      this.feedbackServiceClient.UpdateFeedback(
        { id, ...feedbackData },
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

  async deleteFeedback(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.feedbackServiceClient.DeleteFeedback(
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
}