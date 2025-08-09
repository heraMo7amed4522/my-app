"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ChatService = void 0;
const common_1 = require("@nestjs/common");
const grpc = require("@grpc/grpc-js");
const protoLoader = require("@grpc/proto-loader");
const path_1 = require("path");
let ChatService = class ChatService {
    chatServiceClient;
    async onModuleInit() {
        const PROTO_PATH = (0, path_1.join)(__dirname, '../../proto/chat.proto');
        const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
            keepCase: true,
            longs: String,
            enums: String,
            defaults: true,
            oneofs: true,
        });
        const chatProto = grpc.loadPackageDefinition(packageDefinition);
        const serviceUrl = process.env.CHAT_SERVICE_URL || 'localhost:50054';
        this.chatServiceClient = new chatProto.chat.ChatService(serviceUrl, grpc.credentials.createInsecure());
    }
    async sendMessage(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.SendMessage(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async editMessage(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.EditMessage(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async deleteMessage(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.DeleteMessage(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getChatHistory(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetChatHistory(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async markAsRead(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.MarkAsRead(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async searchMessages(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.SearchMessages(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async forwardMessage(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.ForwardMessage(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async pinMessage(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.PinMessage(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async unpinMessage(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.UnpinMessage(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getPinnedMessages(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetPinnedMessages(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async addLikeMessage(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.AddLikeMessage(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updateLikedMessage(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.UpdateLikedMessage(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getLikedMessages(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetLikedMessages(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getLastMessages(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetLastMessages(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getUsersByUserEmail(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetUsersByUserEmail(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getUsersInGroup(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetUsersInGroup(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getUserStatus(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetUserStatus(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async createGroup(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.CreateGroup(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async joinGroup(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.JoinGroup(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async leaveGroup(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.LeaveGroup(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updateGroup(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.UpdateGroup(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getAllGroupsByUserEmail(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetAllGroupsByUserEmail(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async initiateCall(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.InitiateCall(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async acceptCall(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.AcceptCall(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async rejectCall(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.RejectCall(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async endCall(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.EndCall(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getCallHistory(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetCallHistory(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async addNotification(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.AddNotification(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updateNotification(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.UpdateNotification(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getNotification(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetNotification(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async markNotificationAsRead(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.MarkNotificationAsRead(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getUnreadNotificationCount(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetUnreadNotificationCount(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async addScheduleMessage(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.AddScheduleMessage(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async updateScheduleMessage(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.UpdateScheduleMessage(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async cancelScheduledMessage(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.CancelScheduledMessage(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
    async getScheduledMessages(request, token) {
        return new Promise((resolve, reject) => {
            const metadata = new grpc.Metadata();
            if (token) {
                metadata.add('authorization', `Bearer ${token}`);
            }
            this.chatServiceClient.GetScheduledMessages(request, metadata, (error, response) => {
                if (error) {
                    reject(error);
                }
                else {
                    resolve(response);
                }
            });
        });
    }
};
exports.ChatService = ChatService;
exports.ChatService = ChatService = __decorate([
    (0, common_1.Injectable)()
], ChatService);
//# sourceMappingURL=chat.service.js.map