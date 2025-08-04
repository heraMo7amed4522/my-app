import { Injectable, OnModuleInit } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';

@Injectable()
export class HistoryTemplateService implements OnModuleInit {
  private historyTemplateServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/history_template.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const historyTemplateProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.HISTORY_TEMPLATE_SERVICE_URL || 'localhost:50055';
    
    this.historyTemplateServiceClient = new historyTemplateProto.history_templates.HistoryTemplateService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  async getTemplateById(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.GetTemplateByID(
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

  async getAllTemplates(page: number, limit: number, sortBy?: string, order?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.GetAllTemplates(
        { page, limit, sortBy, order },
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

  async getTemplatesByEra(era: string, page: number, limit: number): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.GetTemplatesByEra(
        { era, page, limit },
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

  async getTemplatesByDynasty(dynasty: number, page: number, limit: number): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.GetTemplatesByDynasty(
        { dynasty, page, limit },
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

  async getTemplatesByPharaoh(pharaohId: string, page: number, limit: number): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.GetTemplatesByPharaoh(
        { pharaoh_id: pharaohId, page, limit },
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

  async getTemplatesByDifficulty(difficulty: string, page: number, limit: number): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.GetTemplatesByDifficulty(
        { difficulty, page, limit },
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

  async searchTemplates(query: string, fields: string[], page: number, limit: number): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.SearchTemplates(
        { query, fields, page, limit },
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

  async createTemplate(templateData: any): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.CreateTemplate(
        { template: templateData },
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

  async updateTemplate(id: string, templateData: any): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.UpdateTemplate(
        { id, template: templateData },
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

  async deleteTemplate(id: string): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.DeleteTemplate(
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

  async getTemplatesByTag(tag: string, page: number, limit: number): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.GetTemplatesByTag(
        { tag, page, limit },
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

  async getRelatedTemplates(templateId: string, limit: number): Promise<any> {
    return new Promise((resolve, reject) => {
      this.historyTemplateServiceClient.GetRelatedTemplates(
        { template_id: templateId, limit },
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