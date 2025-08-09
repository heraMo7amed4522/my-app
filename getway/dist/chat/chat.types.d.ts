export declare enum MessageType {
    MSG_TEXT = "MSG_TEXT",
    MSG_IMAGE = "MSG_IMAGE",
    MSG_VIDEO = "MSG_VIDEO",
    MSG_FILE = "MSG_FILE",
    MSG_AUDIO = "MSG_AUDIO",
    MSG_LOCATION = "MSG_LOCATION",
    MSG_STICKER = "MSG_STICKER",
    MSG_GIF = "MSG_GIF",
    MSG_VOICE_NOTE = "MSG_VOICE_NOTE",
    MSG_SYSTEM = "MSG_SYSTEM",
    MSG_POLL = "MSG_POLL",
    MSG_CONTACT = "MSG_CONTACT",
    MSG_DOCUMENT = "MSG_DOCUMENT",
    MSG_LINK_PREVIEW = "MSG_LINK_PREVIEW"
}
export declare enum CallType {
    CALL_VOICE = "CALL_VOICE",
    CALL_VIDEO = "CALL_VIDEO"
}
export declare enum CallStatus {
    CALL_INITIATED = "CALL_INITIATED",
    CALL_RINGING = "CALL_RINGING",
    CALL_ACCEPTED = "CALL_ACCEPTED",
    CALL_REJECTED = "CALL_REJECTED",
    CALL_ENDED = "CALL_ENDED",
    CALL_MISSED = "CALL_MISSED",
    CALL_BUSY = "CALL_BUSY"
}
export declare enum ReactionType {
    LIKE = "LIKE",
    LOVE = "LOVE",
    LAUGH = "LAUGH",
    WOW = "WOW",
    SAD = "SAD",
    ANGRY = "ANGRY"
}
export declare enum NotificationType {
    NOTIFICATION_MESSAGE = "NOTIFICATION_MESSAGE",
    NOTIFICATION_MENTION = "NOTIFICATION_MENTION",
    NOTIFICATION_REACTION = "NOTIFICATION_REACTION",
    NOTIFICATION_CALL = "NOTIFICATION_CALL",
    NOTIFICATION_GROUP_INVITE = "NOTIFICATION_GROUP_INVITE",
    NOTIFICATION_SYSTEM = "NOTIFICATION_SYSTEM"
}
export declare enum GroupRole {
    MEMBER = "MEMBER",
    ADMIN = "ADMIN",
    OWNER = "OWNER"
}
export declare enum MessageStatus {
    SENT = "SENT",
    DELIVERED = "DELIVERED",
    READ = "READ",
    FAILED = "FAILED"
}
export declare class FileMetadata {
    fileName: string;
    fileType: string;
    fileSize: number;
    fileUrl: string;
    thumbnailUrl?: string;
    duration?: number;
    width?: number;
    height?: number;
}
export declare class LocationData {
    latitude: number;
    longitude: number;
    address?: string;
    placeName?: string;
}
export declare class PollOption {
    id: string;
    text: string;
    voterIds: string[];
    voteCount: number;
}
export declare class PollData {
    question: string;
    options: PollOption[];
    allowMultipleAnswers: boolean;
    expiresAt: string;
    isAnonymous: boolean;
    totalVotes: number;
}
export declare class MessageReaction {
    userId: string;
    reactionType: ReactionType;
    timestamp: string;
}
export declare class ChatMessage {
    messageId: string;
    senderId: string;
    receiverId?: string;
    groupId?: string;
    content: string;
    type: MessageType;
    timestamp: string;
    isGroup: boolean;
    isRead: boolean;
    isEdited: boolean;
    likedBy: string[];
    replyToMessageId?: string;
    attachments: string[];
    status: MessageStatus;
    deliveredAt?: string;
    readAt?: string;
    reactions: MessageReaction[];
    isPinned: boolean;
    pinnedAt?: string;
    pinnedBy?: string;
    forwardCount: number;
    originalMessageId?: string;
    fileMetadata?: FileMetadata;
    locationData?: LocationData;
    pollData?: PollData;
    threadId?: string;
    parentMessageId?: string;
    threadReplyCount: number;
    editedAt?: string;
    editHistory: string[];
    isScheduled: boolean;
    scheduledAt?: string;
    mentionedUserIds: string[];
    isSystemMessage: boolean;
    deviceInfo?: string;
    clientVersion?: string;
    isEncrypted: boolean;
    encryptionKeyId?: string;
}
export declare class ChatMessageList {
    messages: ChatMessage[];
}
export declare class UserInfo {
    id: string;
    fullName: string;
    email: string;
    avatarUrl?: string;
    isOnline: boolean;
}
export declare class UserInfoList {
    users: UserInfo[];
}
export declare class GroupMember {
    userId: string;
    role: GroupRole;
    joinedAt: string;
}
export declare class GroupInfo {
    id: string;
    name: string;
    description?: string;
    avatarUrl?: string;
    creatorId: string;
    members: GroupMember[];
    createdAt: string;
    updatedAt: string;
    maxMembers: number;
    isPrivate: boolean;
}
export declare class GroupInfoList {
    groups: GroupInfo[];
}
export declare class CallInfo {
    callId: string;
    callerId: string;
    receiverId?: string;
    groupId?: string;
    callType: CallType;
    status: CallStatus;
    startTime: string;
    endTime?: string;
    duration: number;
    isGroup: boolean;
    participants: string[];
}
export declare class CallInfoList {
    calls: CallInfo[];
}
export declare class NotificationUpdate {
    id: string;
    type: NotificationType;
    title: string;
    content: string;
    senderId?: string;
    chatId?: string;
    messageId?: string;
    timestamp: string;
    isRead: boolean;
}
export declare class ScheduledMessage {
    id: string;
    chatId: string;
    senderId: string;
    content: string;
    scheduledAt: string;
    type: MessageType;
    attachments: string[];
    createdAt: string;
}
export declare class ScheduledMessageList {
    messages: ScheduledMessage[];
}
export declare class StateMessage {
    statusCode: number;
    message: string;
    timestamp: string;
}
export declare class ErrorMessage {
    code: number;
    message: string;
    details: string[];
    timestamp: string;
}
export declare class SendMessageResponse {
    statusCode: number;
    message: string;
    savedMessage?: ChatMessage;
    error?: ErrorMessage;
}
export declare class GetChatHistoryResponse {
    statusCode: number;
    message: string;
    messages?: ChatMessageList;
    error?: ErrorMessage;
}
export declare class StandardResponse {
    statusCode: number;
    message: string;
    status?: StateMessage;
    error?: ErrorMessage;
}
export declare class GetUsersByUserEmailResponse {
    statusCode: number;
    message: string;
    users?: UserInfoList;
    error?: ErrorMessage;
}
export declare class CreateGroupResponse {
    statusCode: number;
    message: string;
    group?: GroupInfo;
    error?: ErrorMessage;
}
export declare class GetAllGroupsByUserEmailResponse {
    statusCode: number;
    message: string;
    groups?: GroupInfoList;
    error?: ErrorMessage;
}
export declare class InitiateCallResponse {
    statusCode: number;
    message: string;
    call?: CallInfo;
    error?: ErrorMessage;
}
export declare class GetCallHistoryResponse {
    statusCode: number;
    message: string;
    calls?: CallInfoList;
    error?: ErrorMessage;
}
export declare class AddNotificationResponse {
    statusCode: number;
    message: string;
    notification?: NotificationUpdate;
    error?: ErrorMessage;
}
export declare class GetUnreadNotificationCountResponse {
    statusCode: number;
    message: string;
    count?: number;
    error?: ErrorMessage;
}
export declare class AddScheduleMessageResponse {
    statusCode: number;
    message: string;
    scheduledMessageId?: string;
    error?: ErrorMessage;
}
export declare class GetScheduledMessagesResponse {
    statusCode: number;
    message: string;
    messages?: ScheduledMessageList;
    error?: ErrorMessage;
}
export declare class SendMessageInput {
    senderId: string;
    receiverId?: string;
    groupId?: string;
    content: string;
    type: MessageType;
    isGroup: boolean;
}
export declare class EditMessageInput {
    messageId: string;
    userId: string;
    newContent: string;
}
export declare class DeleteMessageInput {
    messageId: string;
    userId: string;
}
export declare class GetChatHistoryInput {
    userId: string;
    peerId: string;
    isGroup: boolean;
    limit?: number;
    offset?: number;
}
export declare class CreateGroupInput {
    creatorId: string;
    groupName: string;
    description?: string;
    memberIds: string[];
    avatarUrl?: string;
}
export declare class UpdateGroupInput {
    groupId: string;
    userId: string;
    groupName?: string;
    description?: string;
    avatarUrl?: string;
    addMemberIds?: string[];
    removeMemberIds?: string[];
}
export declare class InitiateCallInput {
    callerId: string;
    receiverId?: string;
    groupId?: string;
    callType: CallType;
    isGroup: boolean;
}
export declare class NotificationInput {
    id: string;
    type: NotificationType;
    title: string;
    content: string;
    senderId?: string;
    chatId?: string;
    messageId?: string;
    timestamp: string;
    isRead: boolean;
}
export declare class AddScheduleMessageInput {
    chatId: string;
    senderId: string;
    content: string;
    scheduledAt: string;
    type: MessageType;
    attachments?: string[];
}
export declare class UpdateScheduleMessageInput {
    scheduledMessageId: string;
    userId: string;
    content?: string;
    scheduledAt?: string;
}
