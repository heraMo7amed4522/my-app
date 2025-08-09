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
Object.defineProperty(exports, "__esModule", { value: true });
exports.UpdateScheduleMessageInput = exports.AddScheduleMessageInput = exports.NotificationInput = exports.InitiateCallInput = exports.UpdateGroupInput = exports.CreateGroupInput = exports.GetChatHistoryInput = exports.DeleteMessageInput = exports.EditMessageInput = exports.SendMessageInput = exports.GetScheduledMessagesResponse = exports.AddScheduleMessageResponse = exports.GetUnreadNotificationCountResponse = exports.AddNotificationResponse = exports.GetCallHistoryResponse = exports.InitiateCallResponse = exports.GetAllGroupsByUserEmailResponse = exports.CreateGroupResponse = exports.GetUsersByUserEmailResponse = exports.StandardResponse = exports.GetChatHistoryResponse = exports.SendMessageResponse = exports.ErrorMessage = exports.StateMessage = exports.ScheduledMessageList = exports.ScheduledMessage = exports.NotificationUpdate = exports.CallInfoList = exports.CallInfo = exports.GroupInfoList = exports.GroupInfo = exports.GroupMember = exports.UserInfoList = exports.UserInfo = exports.ChatMessageList = exports.ChatMessage = exports.MessageReaction = exports.PollData = exports.PollOption = exports.LocationData = exports.FileMetadata = exports.MessageStatus = exports.GroupRole = exports.NotificationType = exports.ReactionType = exports.CallStatus = exports.CallType = exports.MessageType = void 0;
const graphql_1 = require("@nestjs/graphql");
var MessageType;
(function (MessageType) {
    MessageType["MSG_TEXT"] = "MSG_TEXT";
    MessageType["MSG_IMAGE"] = "MSG_IMAGE";
    MessageType["MSG_VIDEO"] = "MSG_VIDEO";
    MessageType["MSG_FILE"] = "MSG_FILE";
    MessageType["MSG_AUDIO"] = "MSG_AUDIO";
    MessageType["MSG_LOCATION"] = "MSG_LOCATION";
    MessageType["MSG_STICKER"] = "MSG_STICKER";
    MessageType["MSG_GIF"] = "MSG_GIF";
    MessageType["MSG_VOICE_NOTE"] = "MSG_VOICE_NOTE";
    MessageType["MSG_SYSTEM"] = "MSG_SYSTEM";
    MessageType["MSG_POLL"] = "MSG_POLL";
    MessageType["MSG_CONTACT"] = "MSG_CONTACT";
    MessageType["MSG_DOCUMENT"] = "MSG_DOCUMENT";
    MessageType["MSG_LINK_PREVIEW"] = "MSG_LINK_PREVIEW";
})(MessageType || (exports.MessageType = MessageType = {}));
var CallType;
(function (CallType) {
    CallType["CALL_VOICE"] = "CALL_VOICE";
    CallType["CALL_VIDEO"] = "CALL_VIDEO";
})(CallType || (exports.CallType = CallType = {}));
var CallStatus;
(function (CallStatus) {
    CallStatus["CALL_INITIATED"] = "CALL_INITIATED";
    CallStatus["CALL_RINGING"] = "CALL_RINGING";
    CallStatus["CALL_ACCEPTED"] = "CALL_ACCEPTED";
    CallStatus["CALL_REJECTED"] = "CALL_REJECTED";
    CallStatus["CALL_ENDED"] = "CALL_ENDED";
    CallStatus["CALL_MISSED"] = "CALL_MISSED";
    CallStatus["CALL_BUSY"] = "CALL_BUSY";
})(CallStatus || (exports.CallStatus = CallStatus = {}));
var ReactionType;
(function (ReactionType) {
    ReactionType["LIKE"] = "LIKE";
    ReactionType["LOVE"] = "LOVE";
    ReactionType["LAUGH"] = "LAUGH";
    ReactionType["WOW"] = "WOW";
    ReactionType["SAD"] = "SAD";
    ReactionType["ANGRY"] = "ANGRY";
})(ReactionType || (exports.ReactionType = ReactionType = {}));
var NotificationType;
(function (NotificationType) {
    NotificationType["NOTIFICATION_MESSAGE"] = "NOTIFICATION_MESSAGE";
    NotificationType["NOTIFICATION_MENTION"] = "NOTIFICATION_MENTION";
    NotificationType["NOTIFICATION_REACTION"] = "NOTIFICATION_REACTION";
    NotificationType["NOTIFICATION_CALL"] = "NOTIFICATION_CALL";
    NotificationType["NOTIFICATION_GROUP_INVITE"] = "NOTIFICATION_GROUP_INVITE";
    NotificationType["NOTIFICATION_SYSTEM"] = "NOTIFICATION_SYSTEM";
})(NotificationType || (exports.NotificationType = NotificationType = {}));
var GroupRole;
(function (GroupRole) {
    GroupRole["MEMBER"] = "MEMBER";
    GroupRole["ADMIN"] = "ADMIN";
    GroupRole["OWNER"] = "OWNER";
})(GroupRole || (exports.GroupRole = GroupRole = {}));
var MessageStatus;
(function (MessageStatus) {
    MessageStatus["SENT"] = "SENT";
    MessageStatus["DELIVERED"] = "DELIVERED";
    MessageStatus["READ"] = "READ";
    MessageStatus["FAILED"] = "FAILED";
})(MessageStatus || (exports.MessageStatus = MessageStatus = {}));
(0, graphql_1.registerEnumType)(MessageType, { name: 'MessageType' });
(0, graphql_1.registerEnumType)(CallType, { name: 'CallType' });
(0, graphql_1.registerEnumType)(CallStatus, { name: 'CallStatus' });
(0, graphql_1.registerEnumType)(ReactionType, { name: 'ReactionType' });
(0, graphql_1.registerEnumType)(NotificationType, { name: 'NotificationType' });
(0, graphql_1.registerEnumType)(GroupRole, { name: 'GroupRole' });
(0, graphql_1.registerEnumType)(MessageStatus, { name: 'MessageStatus' });
let FileMetadata = class FileMetadata {
    fileName;
    fileType;
    fileSize;
    fileUrl;
    thumbnailUrl;
    duration;
    width;
    height;
};
exports.FileMetadata = FileMetadata;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FileMetadata.prototype, "fileName", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FileMetadata.prototype, "fileType", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], FileMetadata.prototype, "fileSize", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], FileMetadata.prototype, "fileUrl", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], FileMetadata.prototype, "thumbnailUrl", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], FileMetadata.prototype, "duration", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], FileMetadata.prototype, "width", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], FileMetadata.prototype, "height", void 0);
exports.FileMetadata = FileMetadata = __decorate([
    (0, graphql_1.ObjectType)()
], FileMetadata);
let LocationData = class LocationData {
    latitude;
    longitude;
    address;
    placeName;
};
exports.LocationData = LocationData;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Number)
], LocationData.prototype, "latitude", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Number)
], LocationData.prototype, "longitude", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], LocationData.prototype, "address", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], LocationData.prototype, "placeName", void 0);
exports.LocationData = LocationData = __decorate([
    (0, graphql_1.ObjectType)()
], LocationData);
let PollOption = class PollOption {
    id;
    text;
    voterIds;
    voteCount;
};
exports.PollOption = PollOption;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], PollOption.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], PollOption.prototype, "text", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String]),
    __metadata("design:type", Array)
], PollOption.prototype, "voterIds", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], PollOption.prototype, "voteCount", void 0);
exports.PollOption = PollOption = __decorate([
    (0, graphql_1.ObjectType)()
], PollOption);
let PollData = class PollData {
    question;
    options;
    allowMultipleAnswers;
    expiresAt;
    isAnonymous;
    totalVotes;
};
exports.PollData = PollData;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], PollData.prototype, "question", void 0);
__decorate([
    (0, graphql_1.Field)(() => [PollOption]),
    __metadata("design:type", Array)
], PollData.prototype, "options", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], PollData.prototype, "allowMultipleAnswers", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], PollData.prototype, "expiresAt", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], PollData.prototype, "isAnonymous", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], PollData.prototype, "totalVotes", void 0);
exports.PollData = PollData = __decorate([
    (0, graphql_1.ObjectType)()
], PollData);
let MessageReaction = class MessageReaction {
    userId;
    reactionType;
    timestamp;
};
exports.MessageReaction = MessageReaction;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], MessageReaction.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(() => ReactionType),
    __metadata("design:type", String)
], MessageReaction.prototype, "reactionType", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], MessageReaction.prototype, "timestamp", void 0);
exports.MessageReaction = MessageReaction = __decorate([
    (0, graphql_1.ObjectType)()
], MessageReaction);
let ChatMessage = class ChatMessage {
    messageId;
    senderId;
    receiverId;
    groupId;
    content;
    type;
    timestamp;
    isGroup;
    isRead;
    isEdited;
    likedBy;
    replyToMessageId;
    attachments;
    status;
    deliveredAt;
    readAt;
    reactions;
    isPinned;
    pinnedAt;
    pinnedBy;
    forwardCount;
    originalMessageId;
    fileMetadata;
    locationData;
    pollData;
    threadId;
    parentMessageId;
    threadReplyCount;
    editedAt;
    editHistory;
    isScheduled;
    scheduledAt;
    mentionedUserIds;
    isSystemMessage;
    deviceInfo;
    clientVersion;
    isEncrypted;
    encryptionKeyId;
};
exports.ChatMessage = ChatMessage;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ChatMessage.prototype, "messageId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ChatMessage.prototype, "senderId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "receiverId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "groupId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ChatMessage.prototype, "content", void 0);
__decorate([
    (0, graphql_1.Field)(() => MessageType),
    __metadata("design:type", String)
], ChatMessage.prototype, "type", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ChatMessage.prototype, "timestamp", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], ChatMessage.prototype, "isGroup", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], ChatMessage.prototype, "isRead", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], ChatMessage.prototype, "isEdited", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String]),
    __metadata("design:type", Array)
], ChatMessage.prototype, "likedBy", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "replyToMessageId", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String]),
    __metadata("design:type", Array)
], ChatMessage.prototype, "attachments", void 0);
__decorate([
    (0, graphql_1.Field)(() => MessageStatus),
    __metadata("design:type", String)
], ChatMessage.prototype, "status", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "deliveredAt", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "readAt", void 0);
__decorate([
    (0, graphql_1.Field)(() => [MessageReaction]),
    __metadata("design:type", Array)
], ChatMessage.prototype, "reactions", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], ChatMessage.prototype, "isPinned", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "pinnedAt", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "pinnedBy", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], ChatMessage.prototype, "forwardCount", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "originalMessageId", void 0);
__decorate([
    (0, graphql_1.Field)(() => FileMetadata, { nullable: true }),
    __metadata("design:type", FileMetadata)
], ChatMessage.prototype, "fileMetadata", void 0);
__decorate([
    (0, graphql_1.Field)(() => LocationData, { nullable: true }),
    __metadata("design:type", LocationData)
], ChatMessage.prototype, "locationData", void 0);
__decorate([
    (0, graphql_1.Field)(() => PollData, { nullable: true }),
    __metadata("design:type", PollData)
], ChatMessage.prototype, "pollData", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "threadId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "parentMessageId", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], ChatMessage.prototype, "threadReplyCount", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "editedAt", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String]),
    __metadata("design:type", Array)
], ChatMessage.prototype, "editHistory", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], ChatMessage.prototype, "isScheduled", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "scheduledAt", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String]),
    __metadata("design:type", Array)
], ChatMessage.prototype, "mentionedUserIds", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], ChatMessage.prototype, "isSystemMessage", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "deviceInfo", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "clientVersion", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], ChatMessage.prototype, "isEncrypted", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], ChatMessage.prototype, "encryptionKeyId", void 0);
exports.ChatMessage = ChatMessage = __decorate([
    (0, graphql_1.ObjectType)()
], ChatMessage);
let ChatMessageList = class ChatMessageList {
    messages;
};
exports.ChatMessageList = ChatMessageList;
__decorate([
    (0, graphql_1.Field)(() => [ChatMessage]),
    __metadata("design:type", Array)
], ChatMessageList.prototype, "messages", void 0);
exports.ChatMessageList = ChatMessageList = __decorate([
    (0, graphql_1.ObjectType)()
], ChatMessageList);
let UserInfo = class UserInfo {
    id;
    fullName;
    email;
    avatarUrl;
    isOnline;
};
exports.UserInfo = UserInfo;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserInfo.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserInfo.prototype, "fullName", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UserInfo.prototype, "email", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UserInfo.prototype, "avatarUrl", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], UserInfo.prototype, "isOnline", void 0);
exports.UserInfo = UserInfo = __decorate([
    (0, graphql_1.ObjectType)()
], UserInfo);
let UserInfoList = class UserInfoList {
    users;
};
exports.UserInfoList = UserInfoList;
__decorate([
    (0, graphql_1.Field)(() => [UserInfo]),
    __metadata("design:type", Array)
], UserInfoList.prototype, "users", void 0);
exports.UserInfoList = UserInfoList = __decorate([
    (0, graphql_1.ObjectType)()
], UserInfoList);
let GroupMember = class GroupMember {
    userId;
    role;
    joinedAt;
};
exports.GroupMember = GroupMember;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GroupMember.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(() => GroupRole),
    __metadata("design:type", String)
], GroupMember.prototype, "role", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GroupMember.prototype, "joinedAt", void 0);
exports.GroupMember = GroupMember = __decorate([
    (0, graphql_1.ObjectType)()
], GroupMember);
let GroupInfo = class GroupInfo {
    id;
    name;
    description;
    avatarUrl;
    creatorId;
    members;
    createdAt;
    updatedAt;
    maxMembers;
    isPrivate;
};
exports.GroupInfo = GroupInfo;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GroupInfo.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GroupInfo.prototype, "name", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], GroupInfo.prototype, "description", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], GroupInfo.prototype, "avatarUrl", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GroupInfo.prototype, "creatorId", void 0);
__decorate([
    (0, graphql_1.Field)(() => [GroupMember]),
    __metadata("design:type", Array)
], GroupInfo.prototype, "members", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GroupInfo.prototype, "createdAt", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GroupInfo.prototype, "updatedAt", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], GroupInfo.prototype, "maxMembers", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], GroupInfo.prototype, "isPrivate", void 0);
exports.GroupInfo = GroupInfo = __decorate([
    (0, graphql_1.ObjectType)()
], GroupInfo);
let GroupInfoList = class GroupInfoList {
    groups;
};
exports.GroupInfoList = GroupInfoList;
__decorate([
    (0, graphql_1.Field)(() => [GroupInfo]),
    __metadata("design:type", Array)
], GroupInfoList.prototype, "groups", void 0);
exports.GroupInfoList = GroupInfoList = __decorate([
    (0, graphql_1.ObjectType)()
], GroupInfoList);
let CallInfo = class CallInfo {
    callId;
    callerId;
    receiverId;
    groupId;
    callType;
    status;
    startTime;
    endTime;
    duration;
    isGroup;
    participants;
};
exports.CallInfo = CallInfo;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CallInfo.prototype, "callId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CallInfo.prototype, "callerId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CallInfo.prototype, "receiverId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CallInfo.prototype, "groupId", void 0);
__decorate([
    (0, graphql_1.Field)(() => CallType),
    __metadata("design:type", String)
], CallInfo.prototype, "callType", void 0);
__decorate([
    (0, graphql_1.Field)(() => CallStatus),
    __metadata("design:type", String)
], CallInfo.prototype, "status", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CallInfo.prototype, "startTime", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CallInfo.prototype, "endTime", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], CallInfo.prototype, "duration", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], CallInfo.prototype, "isGroup", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String]),
    __metadata("design:type", Array)
], CallInfo.prototype, "participants", void 0);
exports.CallInfo = CallInfo = __decorate([
    (0, graphql_1.ObjectType)()
], CallInfo);
let CallInfoList = class CallInfoList {
    calls;
};
exports.CallInfoList = CallInfoList;
__decorate([
    (0, graphql_1.Field)(() => [CallInfo]),
    __metadata("design:type", Array)
], CallInfoList.prototype, "calls", void 0);
exports.CallInfoList = CallInfoList = __decorate([
    (0, graphql_1.ObjectType)()
], CallInfoList);
let NotificationUpdate = class NotificationUpdate {
    id;
    type;
    title;
    content;
    senderId;
    chatId;
    messageId;
    timestamp;
    isRead;
};
exports.NotificationUpdate = NotificationUpdate;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], NotificationUpdate.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(() => NotificationType),
    __metadata("design:type", String)
], NotificationUpdate.prototype, "type", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], NotificationUpdate.prototype, "title", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], NotificationUpdate.prototype, "content", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], NotificationUpdate.prototype, "senderId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], NotificationUpdate.prototype, "chatId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], NotificationUpdate.prototype, "messageId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], NotificationUpdate.prototype, "timestamp", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], NotificationUpdate.prototype, "isRead", void 0);
exports.NotificationUpdate = NotificationUpdate = __decorate([
    (0, graphql_1.ObjectType)()
], NotificationUpdate);
let ScheduledMessage = class ScheduledMessage {
    id;
    chatId;
    senderId;
    content;
    scheduledAt;
    type;
    attachments;
    createdAt;
};
exports.ScheduledMessage = ScheduledMessage;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ScheduledMessage.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ScheduledMessage.prototype, "chatId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ScheduledMessage.prototype, "senderId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ScheduledMessage.prototype, "content", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ScheduledMessage.prototype, "scheduledAt", void 0);
__decorate([
    (0, graphql_1.Field)(() => MessageType),
    __metadata("design:type", String)
], ScheduledMessage.prototype, "type", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String]),
    __metadata("design:type", Array)
], ScheduledMessage.prototype, "attachments", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ScheduledMessage.prototype, "createdAt", void 0);
exports.ScheduledMessage = ScheduledMessage = __decorate([
    (0, graphql_1.ObjectType)()
], ScheduledMessage);
let ScheduledMessageList = class ScheduledMessageList {
    messages;
};
exports.ScheduledMessageList = ScheduledMessageList;
__decorate([
    (0, graphql_1.Field)(() => [ScheduledMessage]),
    __metadata("design:type", Array)
], ScheduledMessageList.prototype, "messages", void 0);
exports.ScheduledMessageList = ScheduledMessageList = __decorate([
    (0, graphql_1.ObjectType)()
], ScheduledMessageList);
let StateMessage = class StateMessage {
    statusCode;
    message;
    timestamp;
};
exports.StateMessage = StateMessage;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], StateMessage.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], StateMessage.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], StateMessage.prototype, "timestamp", void 0);
exports.StateMessage = StateMessage = __decorate([
    (0, graphql_1.ObjectType)()
], StateMessage);
let ErrorMessage = class ErrorMessage {
    code;
    message;
    details;
    timestamp;
};
exports.ErrorMessage = ErrorMessage;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], ErrorMessage.prototype, "code", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ErrorMessage.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String]),
    __metadata("design:type", Array)
], ErrorMessage.prototype, "details", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ErrorMessage.prototype, "timestamp", void 0);
exports.ErrorMessage = ErrorMessage = __decorate([
    (0, graphql_1.ObjectType)()
], ErrorMessage);
let SendMessageResponse = class SendMessageResponse {
    statusCode;
    message;
    savedMessage;
    error;
};
exports.SendMessageResponse = SendMessageResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], SendMessageResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], SendMessageResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => ChatMessage, { nullable: true }),
    __metadata("design:type", ChatMessage)
], SendMessageResponse.prototype, "savedMessage", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], SendMessageResponse.prototype, "error", void 0);
exports.SendMessageResponse = SendMessageResponse = __decorate([
    (0, graphql_1.ObjectType)()
], SendMessageResponse);
let GetChatHistoryResponse = class GetChatHistoryResponse {
    statusCode;
    message;
    messages;
    error;
};
exports.GetChatHistoryResponse = GetChatHistoryResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], GetChatHistoryResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetChatHistoryResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => ChatMessageList, { nullable: true }),
    __metadata("design:type", ChatMessageList)
], GetChatHistoryResponse.prototype, "messages", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], GetChatHistoryResponse.prototype, "error", void 0);
exports.GetChatHistoryResponse = GetChatHistoryResponse = __decorate([
    (0, graphql_1.ObjectType)()
], GetChatHistoryResponse);
let StandardResponse = class StandardResponse {
    statusCode;
    message;
    status;
    error;
};
exports.StandardResponse = StandardResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], StandardResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], StandardResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => StateMessage, { nullable: true }),
    __metadata("design:type", StateMessage)
], StandardResponse.prototype, "status", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], StandardResponse.prototype, "error", void 0);
exports.StandardResponse = StandardResponse = __decorate([
    (0, graphql_1.ObjectType)()
], StandardResponse);
let GetUsersByUserEmailResponse = class GetUsersByUserEmailResponse {
    statusCode;
    message;
    users;
    error;
};
exports.GetUsersByUserEmailResponse = GetUsersByUserEmailResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], GetUsersByUserEmailResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetUsersByUserEmailResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => UserInfoList, { nullable: true }),
    __metadata("design:type", UserInfoList)
], GetUsersByUserEmailResponse.prototype, "users", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], GetUsersByUserEmailResponse.prototype, "error", void 0);
exports.GetUsersByUserEmailResponse = GetUsersByUserEmailResponse = __decorate([
    (0, graphql_1.ObjectType)()
], GetUsersByUserEmailResponse);
let CreateGroupResponse = class CreateGroupResponse {
    statusCode;
    message;
    group;
    error;
};
exports.CreateGroupResponse = CreateGroupResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], CreateGroupResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateGroupResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => GroupInfo, { nullable: true }),
    __metadata("design:type", GroupInfo)
], CreateGroupResponse.prototype, "group", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], CreateGroupResponse.prototype, "error", void 0);
exports.CreateGroupResponse = CreateGroupResponse = __decorate([
    (0, graphql_1.ObjectType)()
], CreateGroupResponse);
let GetAllGroupsByUserEmailResponse = class GetAllGroupsByUserEmailResponse {
    statusCode;
    message;
    groups;
    error;
};
exports.GetAllGroupsByUserEmailResponse = GetAllGroupsByUserEmailResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], GetAllGroupsByUserEmailResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetAllGroupsByUserEmailResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => GroupInfoList, { nullable: true }),
    __metadata("design:type", GroupInfoList)
], GetAllGroupsByUserEmailResponse.prototype, "groups", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], GetAllGroupsByUserEmailResponse.prototype, "error", void 0);
exports.GetAllGroupsByUserEmailResponse = GetAllGroupsByUserEmailResponse = __decorate([
    (0, graphql_1.ObjectType)()
], GetAllGroupsByUserEmailResponse);
let InitiateCallResponse = class InitiateCallResponse {
    statusCode;
    message;
    call;
    error;
};
exports.InitiateCallResponse = InitiateCallResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], InitiateCallResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], InitiateCallResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => CallInfo, { nullable: true }),
    __metadata("design:type", CallInfo)
], InitiateCallResponse.prototype, "call", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], InitiateCallResponse.prototype, "error", void 0);
exports.InitiateCallResponse = InitiateCallResponse = __decorate([
    (0, graphql_1.ObjectType)()
], InitiateCallResponse);
let GetCallHistoryResponse = class GetCallHistoryResponse {
    statusCode;
    message;
    calls;
    error;
};
exports.GetCallHistoryResponse = GetCallHistoryResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], GetCallHistoryResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetCallHistoryResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => CallInfoList, { nullable: true }),
    __metadata("design:type", CallInfoList)
], GetCallHistoryResponse.prototype, "calls", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], GetCallHistoryResponse.prototype, "error", void 0);
exports.GetCallHistoryResponse = GetCallHistoryResponse = __decorate([
    (0, graphql_1.ObjectType)()
], GetCallHistoryResponse);
let AddNotificationResponse = class AddNotificationResponse {
    statusCode;
    message;
    notification;
    error;
};
exports.AddNotificationResponse = AddNotificationResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], AddNotificationResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], AddNotificationResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => NotificationUpdate, { nullable: true }),
    __metadata("design:type", NotificationUpdate)
], AddNotificationResponse.prototype, "notification", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], AddNotificationResponse.prototype, "error", void 0);
exports.AddNotificationResponse = AddNotificationResponse = __decorate([
    (0, graphql_1.ObjectType)()
], AddNotificationResponse);
let GetUnreadNotificationCountResponse = class GetUnreadNotificationCountResponse {
    statusCode;
    message;
    count;
    error;
};
exports.GetUnreadNotificationCountResponse = GetUnreadNotificationCountResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], GetUnreadNotificationCountResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetUnreadNotificationCountResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], GetUnreadNotificationCountResponse.prototype, "count", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], GetUnreadNotificationCountResponse.prototype, "error", void 0);
exports.GetUnreadNotificationCountResponse = GetUnreadNotificationCountResponse = __decorate([
    (0, graphql_1.ObjectType)()
], GetUnreadNotificationCountResponse);
let AddScheduleMessageResponse = class AddScheduleMessageResponse {
    statusCode;
    message;
    scheduledMessageId;
    error;
};
exports.AddScheduleMessageResponse = AddScheduleMessageResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], AddScheduleMessageResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], AddScheduleMessageResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], AddScheduleMessageResponse.prototype, "scheduledMessageId", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], AddScheduleMessageResponse.prototype, "error", void 0);
exports.AddScheduleMessageResponse = AddScheduleMessageResponse = __decorate([
    (0, graphql_1.ObjectType)()
], AddScheduleMessageResponse);
let GetScheduledMessagesResponse = class GetScheduledMessagesResponse {
    statusCode;
    message;
    messages;
    error;
};
exports.GetScheduledMessagesResponse = GetScheduledMessagesResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], GetScheduledMessagesResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetScheduledMessagesResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => ScheduledMessageList, { nullable: true }),
    __metadata("design:type", ScheduledMessageList)
], GetScheduledMessagesResponse.prototype, "messages", void 0);
__decorate([
    (0, graphql_1.Field)(() => ErrorMessage, { nullable: true }),
    __metadata("design:type", ErrorMessage)
], GetScheduledMessagesResponse.prototype, "error", void 0);
exports.GetScheduledMessagesResponse = GetScheduledMessagesResponse = __decorate([
    (0, graphql_1.ObjectType)()
], GetScheduledMessagesResponse);
let SendMessageInput = class SendMessageInput {
    senderId;
    receiverId;
    groupId;
    content;
    type;
    isGroup;
};
exports.SendMessageInput = SendMessageInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], SendMessageInput.prototype, "senderId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], SendMessageInput.prototype, "receiverId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], SendMessageInput.prototype, "groupId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], SendMessageInput.prototype, "content", void 0);
__decorate([
    (0, graphql_1.Field)(() => MessageType),
    __metadata("design:type", String)
], SendMessageInput.prototype, "type", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], SendMessageInput.prototype, "isGroup", void 0);
exports.SendMessageInput = SendMessageInput = __decorate([
    (0, graphql_1.InputType)()
], SendMessageInput);
let EditMessageInput = class EditMessageInput {
    messageId;
    userId;
    newContent;
};
exports.EditMessageInput = EditMessageInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], EditMessageInput.prototype, "messageId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], EditMessageInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], EditMessageInput.prototype, "newContent", void 0);
exports.EditMessageInput = EditMessageInput = __decorate([
    (0, graphql_1.InputType)()
], EditMessageInput);
let DeleteMessageInput = class DeleteMessageInput {
    messageId;
    userId;
};
exports.DeleteMessageInput = DeleteMessageInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeleteMessageInput.prototype, "messageId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeleteMessageInput.prototype, "userId", void 0);
exports.DeleteMessageInput = DeleteMessageInput = __decorate([
    (0, graphql_1.InputType)()
], DeleteMessageInput);
let GetChatHistoryInput = class GetChatHistoryInput {
    userId;
    peerId;
    isGroup;
    limit;
    offset;
};
exports.GetChatHistoryInput = GetChatHistoryInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetChatHistoryInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetChatHistoryInput.prototype, "peerId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], GetChatHistoryInput.prototype, "isGroup", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], GetChatHistoryInput.prototype, "limit", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], GetChatHistoryInput.prototype, "offset", void 0);
exports.GetChatHistoryInput = GetChatHistoryInput = __decorate([
    (0, graphql_1.InputType)()
], GetChatHistoryInput);
let CreateGroupInput = class CreateGroupInput {
    creatorId;
    groupName;
    description;
    memberIds;
    avatarUrl;
};
exports.CreateGroupInput = CreateGroupInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateGroupInput.prototype, "creatorId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreateGroupInput.prototype, "groupName", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreateGroupInput.prototype, "description", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String]),
    __metadata("design:type", Array)
], CreateGroupInput.prototype, "memberIds", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreateGroupInput.prototype, "avatarUrl", void 0);
exports.CreateGroupInput = CreateGroupInput = __decorate([
    (0, graphql_1.InputType)()
], CreateGroupInput);
let UpdateGroupInput = class UpdateGroupInput {
    groupId;
    userId;
    groupName;
    description;
    avatarUrl;
    addMemberIds;
    removeMemberIds;
};
exports.UpdateGroupInput = UpdateGroupInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateGroupInput.prototype, "groupId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateGroupInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateGroupInput.prototype, "groupName", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateGroupInput.prototype, "description", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateGroupInput.prototype, "avatarUrl", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], UpdateGroupInput.prototype, "addMemberIds", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], UpdateGroupInput.prototype, "removeMemberIds", void 0);
exports.UpdateGroupInput = UpdateGroupInput = __decorate([
    (0, graphql_1.InputType)()
], UpdateGroupInput);
let InitiateCallInput = class InitiateCallInput {
    callerId;
    receiverId;
    groupId;
    callType;
    isGroup;
};
exports.InitiateCallInput = InitiateCallInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], InitiateCallInput.prototype, "callerId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], InitiateCallInput.prototype, "receiverId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], InitiateCallInput.prototype, "groupId", void 0);
__decorate([
    (0, graphql_1.Field)(() => CallType),
    __metadata("design:type", String)
], InitiateCallInput.prototype, "callType", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], InitiateCallInput.prototype, "isGroup", void 0);
exports.InitiateCallInput = InitiateCallInput = __decorate([
    (0, graphql_1.InputType)()
], InitiateCallInput);
let NotificationInput = class NotificationInput {
    id;
    type;
    title;
    content;
    senderId;
    chatId;
    messageId;
    timestamp;
    isRead;
};
exports.NotificationInput = NotificationInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], NotificationInput.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(() => NotificationType),
    __metadata("design:type", String)
], NotificationInput.prototype, "type", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], NotificationInput.prototype, "title", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], NotificationInput.prototype, "content", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], NotificationInput.prototype, "senderId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], NotificationInput.prototype, "chatId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], NotificationInput.prototype, "messageId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], NotificationInput.prototype, "timestamp", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", Boolean)
], NotificationInput.prototype, "isRead", void 0);
exports.NotificationInput = NotificationInput = __decorate([
    (0, graphql_1.InputType)()
], NotificationInput);
let AddScheduleMessageInput = class AddScheduleMessageInput {
    chatId;
    senderId;
    content;
    scheduledAt;
    type;
    attachments;
};
exports.AddScheduleMessageInput = AddScheduleMessageInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], AddScheduleMessageInput.prototype, "chatId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], AddScheduleMessageInput.prototype, "senderId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], AddScheduleMessageInput.prototype, "content", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], AddScheduleMessageInput.prototype, "scheduledAt", void 0);
__decorate([
    (0, graphql_1.Field)(() => MessageType),
    __metadata("design:type", String)
], AddScheduleMessageInput.prototype, "type", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], AddScheduleMessageInput.prototype, "attachments", void 0);
exports.AddScheduleMessageInput = AddScheduleMessageInput = __decorate([
    (0, graphql_1.InputType)()
], AddScheduleMessageInput);
let UpdateScheduleMessageInput = class UpdateScheduleMessageInput {
    scheduledMessageId;
    userId;
    content;
    scheduledAt;
};
exports.UpdateScheduleMessageInput = UpdateScheduleMessageInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateScheduleMessageInput.prototype, "scheduledMessageId", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdateScheduleMessageInput.prototype, "userId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateScheduleMessageInput.prototype, "content", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], UpdateScheduleMessageInput.prototype, "scheduledAt", void 0);
exports.UpdateScheduleMessageInput = UpdateScheduleMessageInput = __decorate([
    (0, graphql_1.InputType)()
], UpdateScheduleMessageInput);
//# sourceMappingURL=chat.types.js.map