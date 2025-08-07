import { Resolver, Query, Mutation, Args, Context } from '@nestjs/graphql';
import { ChatService } from './chat.service';
import {
  // Message Types
  ChatMessage,
  SendMessageResponse,
  GetChatHistoryResponse,
  StandardResponse,
  
  // User Types
  GetUsersByUserEmailResponse,
  UserInfo,
  
  // Group Types
  CreateGroupResponse,
  GetAllGroupsByUserEmailResponse,
  GroupInfo,
  
  // Call Types
  InitiateCallResponse,
  GetCallHistoryResponse,
  CallInfo,
  
  // Notification Types
  AddNotificationResponse,
  GetUnreadNotificationCountResponse,
  NotificationUpdate,
  
  // Scheduled Message Types
  AddScheduleMessageResponse,
  GetScheduledMessagesResponse,
  ScheduledMessage,
  
  // Input Types
  SendMessageInput,
  EditMessageInput,
  DeleteMessageInput,
  GetChatHistoryInput,
  CreateGroupInput,
  UpdateGroupInput,
  InitiateCallInput,
  NotificationInput,
  AddScheduleMessageInput,
  UpdateScheduleMessageInput,
  
  // Enums
  ReactionType,
} from './chat.types';

@Resolver(() => ChatMessage)
export class ChatResolver {
  constructor(private readonly chatService: ChatService) {}

  // ================================================ MESSAGE HANDLING =======================================
  
  @Mutation(() => SendMessageResponse)
  async sendMessage(
    @Args('input') input: SendMessageInput,
    @Context() context: any,
  ): Promise<SendMessageResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.sendMessage(input, token);
  }

  @Mutation(() => StandardResponse)
  async editMessage(
    @Args('input') input: EditMessageInput,
    @Context() context: any,
  ): Promise<StandardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.editMessage(input, token);
  }

  @Mutation(() => StandardResponse)
  async deleteMessage(
    @Args('input') input: DeleteMessageInput,
    @Context() context: any,
  ): Promise<StandardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.deleteMessage(input, token);
  }

  @Query(() => GetChatHistoryResponse)
  async getChatHistory(
    @Args('input') input: GetChatHistoryInput,
    @Context() context: any,
  ): Promise<GetChatHistoryResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.getChatHistory(input, token);
  }

  @Mutation(() => StandardResponse)
  async markAsRead(
    @Args('userId') userId: string,
    @Args('peerId') peerId: string,
    @Args('isGroup') isGroup: boolean,
    @Context() context: any,
  ): Promise<StandardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.markAsRead({ userId, peerId, isGroup }, token);
  }

  @Query(() => GetChatHistoryResponse)
  async searchMessages(
    @Args('userId') userId: string,
    @Args('query') query: string,
    @Args('peerId', { nullable: true }) peerId?: string,
    @Args('groupId', { nullable: true }) groupId?: string,
    @Args('isGroup', { defaultValue: false }) isGroup?: boolean,
    @Args('limit', { nullable: true }) limit?: number,
    @Args('offset', { nullable: true }) offset?: number,
    @Context() context?: any,
  ): Promise<GetChatHistoryResponse> {
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

  @Mutation(() => GetChatHistoryResponse)
  async forwardMessage(
    @Args('messageId') messageId: string,
    @Args('senderId') senderId: string,
    @Args('receiverIds', { type: () => [String] }) receiverIds: string[],
    @Args('groupIds', { type: () => [String] }) groupIds: string[],
    @Context() context: any,
  ): Promise<GetChatHistoryResponse> {
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

  @Mutation(() => StandardResponse)
  async pinMessage(
    @Args('messageId') messageId: string,
    @Args('userId') userId: string,
    @Args('chatId') chatId: string,
    @Args('isGroup') isGroup: boolean,
    @Context() context: any,
  ): Promise<StandardResponse> {
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

  @Mutation(() => StandardResponse)
  async unpinMessage(
    @Args('messageId') messageId: string,
    @Args('userId') userId: string,
    @Args('chatId') chatId: string,
    @Args('isGroup') isGroup: boolean,
    @Context() context: any,
  ): Promise<StandardResponse> {
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

  @Query(() => GetChatHistoryResponse)
  async getPinnedMessages(
    @Args('chatId') chatId: string,
    @Args('isGroup') isGroup: boolean,
    @Context() context: any,
  ): Promise<GetChatHistoryResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.getPinnedMessages({ chatId, isGroup }, token);
  }

  // ================================================ MESSAGE REACTIONS =======================================
  
  @Mutation(() => StandardResponse)
  async addLikeMessage(
    @Args('messageId') messageId: string,
    @Args('userId') userId: string,
    @Context() context: any,
  ): Promise<StandardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.addLikeMessage({ messageId, userId }, token);
  }

  @Mutation(() => StandardResponse)
  async updateLikedMessage(
    @Args('messageId') messageId: string,
    @Args('userId') userId: string,
    @Args('reactionType', { type: () => ReactionType }) reactionType: ReactionType,
    @Context() context: any,
  ): Promise<StandardResponse> {
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

  @Query(() => GetChatHistoryResponse)
  async getLikedMessages(
    @Args('userId') userId: string,
    @Context() context: any,
  ): Promise<GetChatHistoryResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.getLikedMessages({ userId }, token);
  }

  @Query(() => GetChatHistoryResponse)
  async getLastMessages(
    @Args('userId') userId: string,
    @Context() context: any,
  ): Promise<GetChatHistoryResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.getLastMessages({ userId }, token);
  }

  // ================================================ USER MANAGEMENT =======================================
  
  @Query(() => GetUsersByUserEmailResponse)
  async getUsersByUserEmail(
    @Args('userEmails', { type: () => [String] }) userEmails: string[],
    @Context() context: any,
  ): Promise<GetUsersByUserEmailResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.getUsersByUserEmail({ userEmail: userEmails }, token);
  }

  @Query(() => GetUsersByUserEmailResponse)
  async getUsersInGroup(
    @Args('groupId') groupId: string,
    @Context() context: any,
  ): Promise<GetUsersByUserEmailResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.getUsersInGroup({ groupId }, token);
  }

  // ================================================ GROUP MANAGEMENT =======================================
  
  @Mutation(() => CreateGroupResponse)
  async createGroup(
    @Args('input') input: CreateGroupInput,
    @Context() context: any,
  ): Promise<CreateGroupResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.createGroup(input, token);
  }

  @Mutation(() => StandardResponse)
  async joinGroup(
    @Args('userId') userId: string,
    @Args('groupId') groupId: string,
    @Context() context: any,
  ): Promise<StandardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.joinGroup({ userId, groupId }, token);
  }

  @Mutation(() => StandardResponse)
  async leaveGroup(
    @Args('userId') userId: string,
    @Args('groupId') groupId: string,
    @Context() context: any,
  ): Promise<StandardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.leaveGroup({ userId, groupId }, token);
  }

  @Mutation(() => CreateGroupResponse)
  async updateGroup(
    @Args('input') input: UpdateGroupInput,
    @Context() context: any,
  ): Promise<CreateGroupResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.updateGroup(input, token);
  }

  @Query(() => GetAllGroupsByUserEmailResponse)
  async getAllGroupsByUserEmail(
    @Args('userEmail') userEmail: string,
    @Context() context: any,
  ): Promise<GetAllGroupsByUserEmailResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.getAllGroupsByUserEmail({ userEmail }, token);
  }

  // ================================================ CALL MANAGEMENT =======================================
  
  @Mutation(() => InitiateCallResponse)
  async initiateCall(
    @Args('input') input: InitiateCallInput,
    @Context() context: any,
  ): Promise<InitiateCallResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.initiateCall(input, token);
  }

  @Mutation(() => InitiateCallResponse)
  async acceptCall(
    @Args('callId') callId: string,
    @Args('userId') userId: string,
    @Context() context: any,
  ): Promise<InitiateCallResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.acceptCall({ callId, userId }, token);
  }

  @Mutation(() => StandardResponse)
  async rejectCall(
    @Args('callId') callId: string,
    @Args('userId') userId: string,
    @Context() context: any,
  ): Promise<StandardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.rejectCall({ callId, userId }, token);
  }

  @Mutation(() => InitiateCallResponse)
  async endCall(
    @Args('callId') callId: string,
    @Args('userId') userId: string,
    @Context() context: any,
  ): Promise<InitiateCallResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.endCall({ callId, userId }, token);
  }

  @Query(() => GetCallHistoryResponse)
  async getCallHistory(
    @Args('userId') userId: string,
    @Args('limit', { nullable: true }) limit?: number,
    @Args('offset', { nullable: true }) offset?: number,
    @Context() context?: any,
  ): Promise<GetCallHistoryResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.getCallHistory({ userId, limit, offset }, token);
  }

  // ================================================ NOTIFICATIONS =======================================
  
  @Mutation(() => AddNotificationResponse)
  async addNotification(
    @Args('notification') notification: NotificationInput,
    @Context() context: any,
  ): Promise<AddNotificationResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.addNotification({ notification }, token);
  }

  @Mutation(() => AddNotificationResponse)
  async updateNotification(
    @Args('notification') notification: NotificationInput,
    @Context() context: any,
  ): Promise<AddNotificationResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.updateNotification({ notification }, token);
  }

  @Query(() => AddNotificationResponse)
  async getNotification(
    @Args('notificationId') notificationId: string,
    @Args('userId') userId: string,
    @Context() context: any,
  ): Promise<AddNotificationResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.getNotification({ notificationId, userId }, token);
  }

  @Mutation(() => StandardResponse)
  async markNotificationAsRead(
    @Args('notificationId') notificationId: string,
    @Args('userId') userId: string,
    @Context() context: any,
  ): Promise<StandardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.markNotificationAsRead({ notificationId, userId }, token);
  }

  @Query(() => GetUnreadNotificationCountResponse)
  async getUnreadNotificationCount(
    @Args('userId') userId: string,
    @Context() context: any,
  ): Promise<GetUnreadNotificationCountResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.getUnreadNotificationCount({ userId }, token);
  }

  // ================================================ SCHEDULED MESSAGES =======================================
  
  @Mutation(() => AddScheduleMessageResponse)
  async addScheduleMessage(
    @Args('input') input: AddScheduleMessageInput,
    @Context() context: any,
  ): Promise<AddScheduleMessageResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.addScheduleMessage(input, token);
  }

  @Mutation(() => StandardResponse)
  async updateScheduleMessage(
    @Args('input') input: UpdateScheduleMessageInput,
    @Context() context: any,
  ): Promise<StandardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.updateScheduleMessage(input, token);
  }

  @Mutation(() => StandardResponse)
  async cancelScheduledMessage(
    @Args('scheduledMessageId') scheduledMessageId: string,
    @Args('userId') userId: string,
    @Context() context: any,
  ): Promise<StandardResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.cancelScheduledMessage({ scheduledMessageId, userId }, token);
  }

  @Query(() => GetScheduledMessagesResponse)
  async getScheduledMessages(
    @Args('userId') userId: string,
    @Args('chatId', { nullable: true }) chatId?: string,
    @Context() context?: any,
  ): Promise<GetScheduledMessagesResponse> {
    const authHeader = context.req?.headers?.authorization;
    const token = authHeader?.startsWith('Bearer ')
      ? authHeader.substring(7)
      : authHeader;
    return await this.chatService.getScheduledMessages({ userId, chatId }, token);
  }
}