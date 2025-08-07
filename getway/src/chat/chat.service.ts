import { Injectable, OnModuleInit } from '@nestjs/common';
import * as grpc from '@grpc/grpc-js';
import * as protoLoader from '@grpc/proto-loader';
import { join } from 'path';

@Injectable()
export class ChatService implements OnModuleInit {
  private chatServiceClient: any;

  async onModuleInit() {
    const PROTO_PATH = join(__dirname, '../../proto/chat.proto');
    
    const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
      keepCase: true,
      longs: String,
      enums: String,
      defaults: true,
      oneofs: true,
    });

    const chatProto = grpc.loadPackageDefinition(packageDefinition) as any;
    
    const serviceUrl = process.env.CHAT_SERVICE_URL || 'localhost:50054';
    
    this.chatServiceClient = new chatProto.chat.ChatService(
      serviceUrl,
      grpc.credentials.createInsecure(),
    );
  }

  // ================================================ MESSAGE HANDLING =======================================
  
  async sendMessage(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.SendMessage(
        request,
        metadata,
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

  async editMessage(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.EditMessage(
        request,
        metadata,
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

  async deleteMessage(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.DeleteMessage(
        request,
        metadata,
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

  async getChatHistory(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetChatHistory(
        request,
        metadata,
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

  async markAsRead(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.MarkAsRead(
        request,
        metadata,
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

  async searchMessages(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.SearchMessages(
        request,
        metadata,
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

  async forwardMessage(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.ForwardMessage(
        request,
        metadata,
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

  async pinMessage(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.PinMessage(
        request,
        metadata,
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

  async unpinMessage(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.UnpinMessage(
        request,
        metadata,
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

  async getPinnedMessages(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetPinnedMessages(
        request,
        metadata,
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

  // ================================================ MESSAGE REACTIONS =======================================
  
  async addLikeMessage(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.AddLikeMessage(
        request,
        metadata,
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

  async updateLikedMessage(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.UpdateLikedMessage(
        request,
        metadata,
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

  async getLikedMessages(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetLikedMessages(
        request,
        metadata,
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

  async getLastMessages(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetLastMessages(
        request,
        metadata,
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

  // ================================================ USER MANAGEMENT =======================================
  
  async getUsersByUserEmail(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetUsersByUserEmail(
        request,
        metadata,
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

  async getUsersInGroup(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetUsersInGroup(
        request,
        metadata,
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

  async getUserStatus(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetUserStatus(
        request,
        metadata,
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

  // ================================================ GROUP MANAGEMENT =======================================
  
  async createGroup(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.CreateGroup(
        request,
        metadata,
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

  async joinGroup(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.JoinGroup(
        request,
        metadata,
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

  async leaveGroup(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.LeaveGroup(
        request,
        metadata,
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

  async updateGroup(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.UpdateGroup(
        request,
        metadata,
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

  async getAllGroupsByUserEmail(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetAllGroupsByUserEmail(
        request,
        metadata,
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

  // ================================================ CALL MANAGEMENT =======================================
  
  async initiateCall(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.InitiateCall(
        request,
        metadata,
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

  async acceptCall(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.AcceptCall(
        request,
        metadata,
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

  async rejectCall(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.RejectCall(
        request,
        metadata,
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

  async endCall(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.EndCall(
        request,
        metadata,
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

  async getCallHistory(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetCallHistory(
        request,
        metadata,
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

  // ================================================ NOTIFICATIONS =======================================
  
  async addNotification(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.AddNotification(
        request,
        metadata,
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

  async updateNotification(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.UpdateNotification(
        request,
        metadata,
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

  async getNotification(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetNotification(
        request,
        metadata,
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

  async markNotificationAsRead(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.MarkNotificationAsRead(
        request,
        metadata,
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

  async getUnreadNotificationCount(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetUnreadNotificationCount(
        request,
        metadata,
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

  // ================================================ SCHEDULED MESSAGES =======================================
  
  async addScheduleMessage(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.AddScheduleMessage(
        request,
        metadata,
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

  async updateScheduleMessage(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.UpdateScheduleMessage(
        request,
        metadata,
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

  async cancelScheduledMessage(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.CancelScheduledMessage(
        request,
        metadata,
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

  async getScheduledMessages(request: any, token?: string): Promise<any> {
    return new Promise((resolve, reject) => {
      const metadata = new grpc.Metadata();
      if (token) {
        metadata.add('authorization', `Bearer ${token}`);
      }
      this.chatServiceClient.GetScheduledMessages(
        request,
        metadata,
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