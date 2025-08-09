import { ObjectType, Field, InputType, Int, registerEnumType } from '@nestjs/graphql';

export enum MessageType {
  MSG_TEXT = 'MSG_TEXT',
  MSG_IMAGE = 'MSG_IMAGE',
  MSG_VIDEO = 'MSG_VIDEO',
  MSG_FILE = 'MSG_FILE',
  MSG_AUDIO = 'MSG_AUDIO',
  MSG_LOCATION = 'MSG_LOCATION',
  MSG_STICKER = 'MSG_STICKER',
  MSG_GIF = 'MSG_GIF',
  MSG_VOICE_NOTE = 'MSG_VOICE_NOTE',
  MSG_SYSTEM = 'MSG_SYSTEM',
  MSG_POLL = 'MSG_POLL',
  MSG_CONTACT = 'MSG_CONTACT',
  MSG_DOCUMENT = 'MSG_DOCUMENT',
  MSG_LINK_PREVIEW = 'MSG_LINK_PREVIEW',
}

export enum CallType {
  CALL_VOICE = 'CALL_VOICE',
  CALL_VIDEO = 'CALL_VIDEO',
}

export enum CallStatus {
  CALL_INITIATED = 'CALL_INITIATED',
  CALL_RINGING = 'CALL_RINGING',
  CALL_ACCEPTED = 'CALL_ACCEPTED',
  CALL_REJECTED = 'CALL_REJECTED',
  CALL_ENDED = 'CALL_ENDED',
  CALL_MISSED = 'CALL_MISSED',
  CALL_BUSY = 'CALL_BUSY',
}

export enum ReactionType {
  LIKE = 'LIKE',
  LOVE = 'LOVE',
  LAUGH = 'LAUGH',
  WOW = 'WOW',
  SAD = 'SAD',
  ANGRY = 'ANGRY',
}

export enum NotificationType {
  NOTIFICATION_MESSAGE = 'NOTIFICATION_MESSAGE',
  NOTIFICATION_MENTION = 'NOTIFICATION_MENTION',
  NOTIFICATION_REACTION = 'NOTIFICATION_REACTION',
  NOTIFICATION_CALL = 'NOTIFICATION_CALL',
  NOTIFICATION_GROUP_INVITE = 'NOTIFICATION_GROUP_INVITE',
  NOTIFICATION_SYSTEM = 'NOTIFICATION_SYSTEM',
}

export enum GroupRole {
  MEMBER = 'MEMBER',
  ADMIN = 'ADMIN',
  OWNER = 'OWNER',
}

export enum MessageStatus {
  SENT = 'SENT',
  DELIVERED = 'DELIVERED',
  READ = 'READ',
  FAILED = 'FAILED',
}

registerEnumType(MessageType, { name: 'MessageType' });
registerEnumType(CallType, { name: 'CallType' });
registerEnumType(CallStatus, { name: 'CallStatus' });
registerEnumType(ReactionType, { name: 'ReactionType' });
registerEnumType(NotificationType, { name: 'NotificationType' });
registerEnumType(GroupRole, { name: 'GroupRole' });
registerEnumType(MessageStatus, { name: 'MessageStatus' });


@ObjectType()
export class FileMetadata {
  @Field()
  fileName: string;

  @Field()
  fileType: string;

  @Field(() => Int)
  fileSize: number;

  @Field()
  fileUrl: string;

  @Field({ nullable: true })
  thumbnailUrl?: string;

  @Field(() => Int, { nullable: true })
  duration?: number;

  @Field(() => Int, { nullable: true })
  width?: number;

  @Field(() => Int, { nullable: true })
  height?: number;
}

@ObjectType()
export class LocationData {
  @Field()
  latitude: number;

  @Field()
  longitude: number;

  @Field({ nullable: true })
  address?: string;

  @Field({ nullable: true })
  placeName?: string;
}

@ObjectType()
export class PollOption {
  @Field()
  id: string;

  @Field()
  text: string;

  @Field(() => [String])
  voterIds: string[];

  @Field(() => Int)
  voteCount: number;
}

@ObjectType()
export class PollData {
  @Field()
  question: string;

  @Field(() => [PollOption])
  options: PollOption[];

  @Field()
  allowMultipleAnswers: boolean;

  @Field()
  expiresAt: string;

  @Field()
  isAnonymous: boolean;

  @Field(() => Int)
  totalVotes: number;
}

@ObjectType()
export class MessageReaction {
  @Field()
  userId: string;

  @Field(() => ReactionType)
  reactionType: ReactionType;

  @Field()
  timestamp: string;
}

@ObjectType()
export class ChatMessage {
  @Field()
  messageId: string;

  @Field()
  senderId: string;

  @Field({ nullable: true })
  receiverId?: string;

  @Field({ nullable: true })
  groupId?: string;

  @Field()
  content: string;

  @Field(() => MessageType)
  type: MessageType;

  @Field()
  timestamp: string;

  @Field()
  isGroup: boolean;

  @Field()
  isRead: boolean;

  @Field()
  isEdited: boolean;

  @Field(() => [String])
  likedBy: string[];

  @Field({ nullable: true })
  replyToMessageId?: string;

  @Field(() => [String])
  attachments: string[];

  @Field(() => MessageStatus)
  status: MessageStatus;

  @Field({ nullable: true })
  deliveredAt?: string;

  @Field({ nullable: true })
  readAt?: string;

  @Field(() => [MessageReaction])
  reactions: MessageReaction[];

  @Field()
  isPinned: boolean;

  @Field({ nullable: true })
  pinnedAt?: string;

  @Field({ nullable: true })
  pinnedBy?: string;

  @Field(() => Int)
  forwardCount: number;

  @Field({ nullable: true })
  originalMessageId?: string;

  @Field(() => FileMetadata, { nullable: true })
  fileMetadata?: FileMetadata;

  @Field(() => LocationData, { nullable: true })
  locationData?: LocationData;

  @Field(() => PollData, { nullable: true })
  pollData?: PollData;

  @Field({ nullable: true })
  threadId?: string;

  @Field({ nullable: true })
  parentMessageId?: string;

  @Field(() => Int)
  threadReplyCount: number;

  @Field({ nullable: true })
  editedAt?: string;

  @Field(() => [String])
  editHistory: string[];

  @Field()
  isScheduled: boolean;

  @Field({ nullable: true })
  scheduledAt?: string;

  @Field(() => [String])
  mentionedUserIds: string[];

  @Field()
  isSystemMessage: boolean;

  @Field({ nullable: true })
  deviceInfo?: string;

  @Field({ nullable: true })
  clientVersion?: string;

  @Field()
  isEncrypted: boolean;

  @Field({ nullable: true })
  encryptionKeyId?: string;
}

@ObjectType()
export class ChatMessageList {
  @Field(() => [ChatMessage])
  messages: ChatMessage[];
}

@ObjectType()
export class ChatUserInfo {
  @Field()
  id: string;

  @Field()
  fullName: string;

  @Field()
  email: string;

  @Field({ nullable: true })
  avatarUrl?: string;

  @Field()
  isOnline: boolean;
}

@ObjectType()
export class UserInfoList {
  @Field(() => [ChatUserInfo])  
  users: ChatUserInfo[];
}

@ObjectType()
export class GroupMember {
  @Field()
  userId: string;

  @Field(() => GroupRole)
  role: GroupRole;

  @Field()
  joinedAt: string;
}

@ObjectType()
export class GroupInfo {
  @Field()
  id: string;

  @Field()
  name: string;

  @Field({ nullable: true })
  description?: string;

  @Field({ nullable: true })
  avatarUrl?: string;

  @Field()
  creatorId: string;

  @Field(() => [GroupMember])
  members: GroupMember[];

  @Field()
  createdAt: string;

  @Field()
  updatedAt: string;

  @Field(() => Int)
  maxMembers: number;

  @Field()
  isPrivate: boolean;
}

@ObjectType()
export class GroupInfoList {
  @Field(() => [GroupInfo])
  groups: GroupInfo[];
}

@ObjectType()
export class CallInfo {
  @Field()
  callId: string;

  @Field()
  callerId: string;

  @Field({ nullable: true })
  receiverId?: string;

  @Field({ nullable: true })
  groupId?: string;

  @Field(() => CallType)
  callType: CallType;

  @Field(() => CallStatus)
  status: CallStatus;

  @Field()
  startTime: string;

  @Field({ nullable: true })
  endTime?: string;

  @Field(() => Int)
  duration: number;

  @Field()
  isGroup: boolean;

  @Field(() => [String])
  participants: string[];
}

@ObjectType()
export class CallInfoList {
  @Field(() => [CallInfo])
  calls: CallInfo[];
}

@ObjectType()
export class NotificationUpdate {
  @Field()
  id: string;

  @Field(() => NotificationType)
  type: NotificationType;

  @Field()
  title: string;

  @Field()
  content: string;

  @Field({ nullable: true })
  senderId?: string;

  @Field({ nullable: true })
  chatId?: string;

  @Field({ nullable: true })
  messageId?: string;

  @Field()
  timestamp: string;

  @Field()
  isRead: boolean;
}

@ObjectType()
export class ScheduledMessage {
  @Field()
  id: string;

  @Field()
  chatId: string;

  @Field()
  senderId: string;

  @Field()
  content: string;

  @Field()
  scheduledAt: string;

  @Field(() => MessageType)
  type: MessageType;

  @Field(() => [String])
  attachments: string[];

  @Field()
  createdAt: string;
}

@ObjectType()
export class ScheduledMessageList {
  @Field(() => [ScheduledMessage])
  messages: ScheduledMessage[];
}

@ObjectType()
export class StateMessage {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field()
  timestamp: string;
}

@ObjectType()
export class ErrorMessage {
  @Field(() => Int)
  code: number;

  @Field()
  message: string;

  @Field(() => [String])
  details: string[];

  @Field()
  timestamp: string;
}


@ObjectType()
export class SendMessageResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => ChatMessage, { nullable: true })
  savedMessage?: ChatMessage;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}

@ObjectType()
export class GetChatHistoryResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => ChatMessageList, { nullable: true })
  messages?: ChatMessageList;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}

@ObjectType()
export class StandardResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => StateMessage, { nullable: true })
  status?: StateMessage;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}

@ObjectType()
export class GetUsersByUserEmailResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => UserInfoList, { nullable: true })
  users?: UserInfoList;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}

@ObjectType()
export class CreateGroupResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => GroupInfo, { nullable: true })
  group?: GroupInfo;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}

@ObjectType()
export class GetAllGroupsByUserEmailResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => GroupInfoList, { nullable: true })
  groups?: GroupInfoList;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}

@ObjectType()
export class InitiateCallResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => CallInfo, { nullable: true })
  call?: CallInfo;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}

@ObjectType()
export class GetCallHistoryResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => CallInfoList, { nullable: true })
  calls?: CallInfoList;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}

@ObjectType()
export class AddNotificationResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => NotificationUpdate, { nullable: true })
  notification?: NotificationUpdate;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}

@ObjectType()
export class GetUnreadNotificationCountResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => Int, { nullable: true })
  count?: number;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}

@ObjectType()
export class AddScheduleMessageResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field({ nullable: true })
  scheduledMessageId?: string;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}

@ObjectType()
export class GetScheduledMessagesResponse {
  @Field(() => Int)
  statusCode: number;

  @Field()
  message: string;

  @Field(() => ScheduledMessageList, { nullable: true })
  messages?: ScheduledMessageList;

  @Field(() => ErrorMessage, { nullable: true })
  error?: ErrorMessage;
}


@InputType()
export class SendMessageInput {
  @Field()
  senderId: string;

  @Field({ nullable: true })
  receiverId?: string;

  @Field({ nullable: true })
  groupId?: string;

  @Field()
  content: string;

  @Field(() => MessageType)
  type: MessageType;

  @Field()
  isGroup: boolean;
}

@InputType()
export class EditMessageInput {
  @Field()
  messageId: string;

  @Field()
  userId: string;

  @Field()
  newContent: string;
}

@InputType()
export class DeleteMessageInput {
  @Field()
  messageId: string;

  @Field()
  userId: string;
}

@InputType()
export class GetChatHistoryInput {
  @Field()
  userId: string;

  @Field()
  peerId: string;

  @Field()
  isGroup: boolean;

  @Field(() => Int, { nullable: true })
  limit?: number;

  @Field(() => Int, { nullable: true })
  offset?: number;
}

@InputType()
export class CreateGroupInput {
  @Field()
  creatorId: string;

  @Field()
  groupName: string;

  @Field({ nullable: true })
  description?: string;

  @Field(() => [String])
  memberIds: string[];

  @Field({ nullable: true })
  avatarUrl?: string;
}

@InputType()
export class UpdateGroupInput {
  @Field()
  groupId: string;

  @Field()
  userId: string;

  @Field({ nullable: true })
  groupName?: string;

  @Field({ nullable: true })
  description?: string;

  @Field({ nullable: true })
  avatarUrl?: string;

  @Field(() => [String], { nullable: true })
  addMemberIds?: string[];

  @Field(() => [String], { nullable: true })
  removeMemberIds?: string[];
}

@InputType()
export class InitiateCallInput {
  @Field()
  callerId: string;

  @Field({ nullable: true })
  receiverId?: string;

  @Field({ nullable: true })
  groupId?: string;

  @Field(() => CallType)
  callType: CallType;

  @Field()
  isGroup: boolean;
}

@InputType()
export class NotificationInput {
  @Field()
  id: string;

  @Field(() => NotificationType)
  type: NotificationType;

  @Field()
  title: string;

  @Field()
  content: string;

  @Field({ nullable: true })
  senderId?: string;

  @Field({ nullable: true })
  chatId?: string;

  @Field({ nullable: true })
  messageId?: string;

  @Field()
  timestamp: string;

  @Field()
  isRead: boolean;
}

@InputType()
export class AddScheduleMessageInput {
  @Field()
  chatId: string;

  @Field()
  senderId: string;

  @Field()
  content: string;

  @Field()
  scheduledAt: string;

  @Field(() => MessageType)
  type: MessageType;

  @Field(() => [String], { nullable: true })
  attachments?: string[];
}

@InputType()
export class UpdateScheduleMessageInput {
  @Field()
  scheduledMessageId: string;

  @Field()
  userId: string;

  @Field({ nullable: true })
  content?: string;

  @Field({ nullable: true })
  scheduledAt?: string;
}