"use strict";
var __decorate = (this && this.__decorate) || function (decorators, target, key, desc) {
    var c = arguments.length, r = c < 3 ? target : desc === null ? desc = Object.getOwnPropertyDescriptor(target, key) : desc, d;
    if (typeof Reflect === "object" && typeof Reflect.decorate === "function") r = Reflect.decorate(decorators, target, key, desc);
    else for (var i = decorators.length - 1; i >= 0; i--) if (d = decorators[i]) r = (c < 3 ? d(r) : c > 3 ? d(target, key, r) : d(target, key)) || r;
    return c > 3 && r && Object.defineProperty(target, key, r), r;
};
var __metadata = (this && this.__metadata) || function (k, v) {
    if (typeof Reflect === "object" && typeof Reflect.metadata === "function") return Reflect.metadata(k, v);
};
var __param = (this && this.__param) || function (paramIndex, decorator) {
    return function (target, key) { decorator(target, key, paramIndex); }
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.ChatResolver = void 0;
const graphql_1 = require("@nestjs/graphql");
const chat_service_1 = require("./chat.service");
const chat_types_1 = require("./chat.types");
let ChatResolver = class ChatResolver {
    chatService;
    constructor(chatService) {
        this.chatService = chatService;
    }
    async sendMessage(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.sendMessage(input, token);
    }
    async editMessage(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.editMessage(input, token);
    }
    async deleteMessage(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.deleteMessage(input, token);
    }
    async getChatHistory(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.getChatHistory(input, token);
    }
    async markAsRead(userId, peerId, isGroup, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.markAsRead({ userId, peerId, isGroup }, token);
    }
    async searchMessages(userId, query, peerId, groupId, isGroup, limit, offset, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.searchMessages({
            userId,
            query,
            peerId,
            groupId,
            isGroup,
            limit,
            offset,
        }, token);
    }
    async forwardMessage(messageId, senderId, receiverIds, groupIds, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.forwardMessage({
            messageId,
            senderId,
            receiverIds,
            groupIds,
        }, token);
    }
    async pinMessage(messageId, userId, chatId, isGroup, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.pinMessage({
            messageId,
            userId,
            chatId,
            isGroup,
        }, token);
    }
    async unpinMessage(messageId, userId, chatId, isGroup, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.unpinMessage({
            messageId,
            userId,
            chatId,
            isGroup,
        }, token);
    }
    async getPinnedMessages(chatId, isGroup, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.getPinnedMessages({ chatId, isGroup }, token);
    }
    async addLikeMessage(messageId, userId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.addLikeMessage({ messageId, userId }, token);
    }
    async updateLikedMessage(messageId, userId, reactionType, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.updateLikedMessage({
            messageId,
            userId,
            reactionType,
        }, token);
    }
    async getLikedMessages(userId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.getLikedMessages({ userId }, token);
    }
    async getLastMessages(userId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.getLastMessages({ userId }, token);
    }
    async getUsersByUserEmail(userEmails, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.getUsersByUserEmail({ userEmail: userEmails }, token);
    }
    async getUsersInGroup(groupId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.getUsersInGroup({ groupId }, token);
    }
    async createGroup(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.createGroup(input, token);
    }
    async joinGroup(userId, groupId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.joinGroup({ userId, groupId }, token);
    }
    async leaveGroup(userId, groupId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.leaveGroup({ userId, groupId }, token);
    }
    async updateGroup(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.updateGroup(input, token);
    }
    async getAllGroupsByUserEmail(userEmail, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.getAllGroupsByUserEmail({ userEmail }, token);
    }
    async initiateCall(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.initiateCall(input, token);
    }
    async acceptCall(callId, userId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.acceptCall({ callId, userId }, token);
    }
    async rejectCall(callId, userId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.rejectCall({ callId, userId }, token);
    }
    async endCall(callId, userId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.endCall({ callId, userId }, token);
    }
    async getCallHistory(userId, limit, offset, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.getCallHistory({ userId, limit, offset }, token);
    }
    async addNotification(notification, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.addNotification({ notification }, token);
    }
    async updateNotification(notification, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.updateNotification({ notification }, token);
    }
    async getNotification(notificationId, userId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.getNotification({ notificationId, userId }, token);
    }
    async markNotificationAsRead(notificationId, userId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.markNotificationAsRead({ notificationId, userId }, token);
    }
    async getUnreadNotificationCount(userId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.getUnreadNotificationCount({ userId }, token);
    }
    async addScheduleMessage(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.addScheduleMessage(input, token);
    }
    async updateScheduleMessage(input, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.updateScheduleMessage(input, token);
    }
    async cancelScheduledMessage(scheduledMessageId, userId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.cancelScheduledMessage({ scheduledMessageId, userId }, token);
    }
    async getScheduledMessages(userId, chatId, context) {
        const authHeader = context.req?.headers?.authorization;
        const token = authHeader?.startsWith('Bearer ')
            ? authHeader.substring(7)
            : authHeader;
        return await this.chatService.getScheduledMessages({ userId, chatId }, token);
    }
};
exports.ChatResolver = ChatResolver;
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.SendMessageResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [chat_types_1.SendMessageInput, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "sendMessage", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [chat_types_1.EditMessageInput, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "editMessage", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [chat_types_1.DeleteMessageInput, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "deleteMessage", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.GetChatHistoryResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [chat_types_1.GetChatHistoryInput, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "getChatHistory", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __param(1, (0, graphql_1.Args)('peerId')),
    __param(2, (0, graphql_1.Args)('isGroup')),
    __param(3, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Boolean, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "markAsRead", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.GetChatHistoryResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __param(1, (0, graphql_1.Args)('query')),
    __param(2, (0, graphql_1.Args)('peerId', { nullable: true })),
    __param(3, (0, graphql_1.Args)('groupId', { nullable: true })),
    __param(4, (0, graphql_1.Args)('isGroup', { defaultValue: false })),
    __param(5, (0, graphql_1.Args)('limit', { nullable: true })),
    __param(6, (0, graphql_1.Args)('offset', { nullable: true })),
    __param(7, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, String, String, Boolean, Number, Number, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "searchMessages", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.GetChatHistoryResponse),
    __param(0, (0, graphql_1.Args)('messageId')),
    __param(1, (0, graphql_1.Args)('senderId')),
    __param(2, (0, graphql_1.Args)('receiverIds', { type: () => [String] })),
    __param(3, (0, graphql_1.Args)('groupIds', { type: () => [String] })),
    __param(4, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Array, Array, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "forwardMessage", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('messageId')),
    __param(1, (0, graphql_1.Args)('userId')),
    __param(2, (0, graphql_1.Args)('chatId')),
    __param(3, (0, graphql_1.Args)('isGroup')),
    __param(4, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, String, Boolean, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "pinMessage", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('messageId')),
    __param(1, (0, graphql_1.Args)('userId')),
    __param(2, (0, graphql_1.Args)('chatId')),
    __param(3, (0, graphql_1.Args)('isGroup')),
    __param(4, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, String, Boolean, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "unpinMessage", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.GetChatHistoryResponse),
    __param(0, (0, graphql_1.Args)('chatId')),
    __param(1, (0, graphql_1.Args)('isGroup')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Boolean, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "getPinnedMessages", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('messageId')),
    __param(1, (0, graphql_1.Args)('userId')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "addLikeMessage", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('messageId')),
    __param(1, (0, graphql_1.Args)('userId')),
    __param(2, (0, graphql_1.Args)('reactionType', { type: () => chat_types_1.ReactionType })),
    __param(3, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "updateLikedMessage", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.GetChatHistoryResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "getLikedMessages", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.GetChatHistoryResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "getLastMessages", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.GetUsersByUserEmailResponse),
    __param(0, (0, graphql_1.Args)('userEmails', { type: () => [String] })),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [Array, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "getUsersByUserEmail", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.GetUsersByUserEmailResponse),
    __param(0, (0, graphql_1.Args)('groupId')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "getUsersInGroup", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.CreateGroupResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [chat_types_1.CreateGroupInput, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "createGroup", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __param(1, (0, graphql_1.Args)('groupId')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "joinGroup", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __param(1, (0, graphql_1.Args)('groupId')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "leaveGroup", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.CreateGroupResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [chat_types_1.UpdateGroupInput, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "updateGroup", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.GetAllGroupsByUserEmailResponse),
    __param(0, (0, graphql_1.Args)('userEmail')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "getAllGroupsByUserEmail", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.InitiateCallResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [chat_types_1.InitiateCallInput, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "initiateCall", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.InitiateCallResponse),
    __param(0, (0, graphql_1.Args)('callId')),
    __param(1, (0, graphql_1.Args)('userId')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "acceptCall", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('callId')),
    __param(1, (0, graphql_1.Args)('userId')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "rejectCall", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.InitiateCallResponse),
    __param(0, (0, graphql_1.Args)('callId')),
    __param(1, (0, graphql_1.Args)('userId')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "endCall", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.GetCallHistoryResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __param(1, (0, graphql_1.Args)('limit', { nullable: true })),
    __param(2, (0, graphql_1.Args)('offset', { nullable: true })),
    __param(3, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Number, Number, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "getCallHistory", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.AddNotificationResponse),
    __param(0, (0, graphql_1.Args)('notification')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [chat_types_1.NotificationInput, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "addNotification", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.AddNotificationResponse),
    __param(0, (0, graphql_1.Args)('notification')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [chat_types_1.NotificationInput, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "updateNotification", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.AddNotificationResponse),
    __param(0, (0, graphql_1.Args)('notificationId')),
    __param(1, (0, graphql_1.Args)('userId')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "getNotification", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('notificationId')),
    __param(1, (0, graphql_1.Args)('userId')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "markNotificationAsRead", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.GetUnreadNotificationCountResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "getUnreadNotificationCount", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.AddScheduleMessageResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [chat_types_1.AddScheduleMessageInput, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "addScheduleMessage", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('input')),
    __param(1, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [chat_types_1.UpdateScheduleMessageInput, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "updateScheduleMessage", null);
__decorate([
    (0, graphql_1.Mutation)(() => chat_types_1.StandardResponse),
    __param(0, (0, graphql_1.Args)('scheduledMessageId')),
    __param(1, (0, graphql_1.Args)('userId')),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "cancelScheduledMessage", null);
__decorate([
    (0, graphql_1.Query)(() => chat_types_1.GetScheduledMessagesResponse),
    __param(0, (0, graphql_1.Args)('userId')),
    __param(1, (0, graphql_1.Args)('chatId', { nullable: true })),
    __param(2, (0, graphql_1.Context)()),
    __metadata("design:type", Function),
    __metadata("design:paramtypes", [String, String, Object]),
    __metadata("design:returntype", Promise)
], ChatResolver.prototype, "getScheduledMessages", null);
exports.ChatResolver = ChatResolver = __decorate([
    (0, graphql_1.Resolver)(() => chat_types_1.ChatMessage),
    __metadata("design:paramtypes", [chat_service_1.ChatService])
], ChatResolver);
//# sourceMappingURL=chat.resolver.js.map