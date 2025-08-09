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
exports.UpdatePharaohRatingInput = exports.GetPharaohsByRarityInput = exports.SearchPharaohsInput = exports.GetPharaohsByPeriodInput = exports.GetPharaohsByDynastyInput = exports.GetAllPharaohsInput = exports.UpdatePharaohInput = exports.CreatePharaohInput = exports.TraitsInput = exports.ArtifactInput = exports.DeletePharaohResponse = exports.PharaohsResponse = exports.PharaohResponse = exports.Pharaoh = exports.Traits = exports.Artifact = void 0;
const graphql_1 = require("@nestjs/graphql");
const shared_types_1 = require("../shared/shared.types");
let Artifact = class Artifact {
    name;
    museum;
    description;
};
exports.Artifact = Artifact;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Artifact.prototype, "name", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Artifact.prototype, "museum", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Artifact.prototype, "description", void 0);
exports.Artifact = Artifact = __decorate([
    (0, graphql_1.ObjectType)()
], Artifact);
let Traits = class Traits {
    leadership;
    military;
    diplomacy;
    wisdom;
    charisma;
};
exports.Traits = Traits;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], Traits.prototype, "leadership", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], Traits.prototype, "military", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], Traits.prototype, "diplomacy", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], Traits.prototype, "wisdom", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], Traits.prototype, "charisma", void 0);
exports.Traits = Traits = __decorate([
    (0, graphql_1.ObjectType)()
], Traits);
let Pharaoh = class Pharaoh {
    id;
    name;
    birthName;
    throneName;
    epithet;
    dynasty;
    period;
    reignStart;
    reignEnd;
    lengthOfReignYears;
    predecessorId;
    successorId;
    father;
    mother;
    consorts;
    childrenCount;
    notableChildren;
    capital;
    majorAchievements;
    militaryCampaigns;
    buildingProjects;
    politicalStyle;
    divineAssociation;
    templeAffiliations;
    religiousReforms;
    pharaohAsGod;
    burialSite;
    tombDiscovered;
    discoveryYear;
    tombGuardian;
    funeraryText;
    famousArtifacts;
    treasureStatus;
    imageUrl;
    statueCount;
    mummyLocation;
    audioNarrationUrl;
    videoDocumentaryUrl;
    popularityScore;
    userRating;
    unlockInGame;
    rarity;
    traits;
    source;
    verified;
    language;
    createdAt;
    updatedAt;
};
exports.Pharaoh = Pharaoh;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Pharaoh.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Pharaoh.prototype, "name", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "birthName", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "throneName", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "epithet", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], Pharaoh.prototype, "dynasty", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "period", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], Pharaoh.prototype, "reignStart", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], Pharaoh.prototype, "reignEnd", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], Pharaoh.prototype, "lengthOfReignYears", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "predecessorId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "successorId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "father", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "mother", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], Pharaoh.prototype, "consorts", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], Pharaoh.prototype, "childrenCount", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], Pharaoh.prototype, "notableChildren", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "capital", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], Pharaoh.prototype, "majorAchievements", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], Pharaoh.prototype, "militaryCampaigns", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], Pharaoh.prototype, "buildingProjects", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "politicalStyle", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], Pharaoh.prototype, "divineAssociation", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], Pharaoh.prototype, "templeAffiliations", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "religiousReforms", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], Pharaoh.prototype, "pharaohAsGod", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "burialSite", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], Pharaoh.prototype, "tombDiscovered", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], Pharaoh.prototype, "discoveryYear", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "tombGuardian", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "funeraryText", void 0);
__decorate([
    (0, graphql_1.Field)(() => [Artifact], { nullable: true }),
    __metadata("design:type", Array)
], Pharaoh.prototype, "famousArtifacts", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "treasureStatus", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "imageUrl", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], Pharaoh.prototype, "statueCount", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "mummyLocation", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "audioNarrationUrl", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "videoDocumentaryUrl", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float, { nullable: true }),
    __metadata("design:type", Number)
], Pharaoh.prototype, "popularityScore", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float, { nullable: true }),
    __metadata("design:type", Number)
], Pharaoh.prototype, "userRating", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], Pharaoh.prototype, "unlockInGame", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "rarity", void 0);
__decorate([
    (0, graphql_1.Field)(() => Traits, { nullable: true }),
    __metadata("design:type", Traits)
], Pharaoh.prototype, "traits", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "source", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], Pharaoh.prototype, "verified", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], Pharaoh.prototype, "language", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Pharaoh.prototype, "createdAt", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], Pharaoh.prototype, "updatedAt", void 0);
exports.Pharaoh = Pharaoh = __decorate([
    (0, graphql_1.ObjectType)()
], Pharaoh);
let PharaohResponse = class PharaohResponse {
    statusCode;
    message;
    pharaoh;
    error;
};
exports.PharaohResponse = PharaohResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], PharaohResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], PharaohResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => Pharaoh, { nullable: true }),
    __metadata("design:type", Pharaoh)
], PharaohResponse.prototype, "pharaoh", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], PharaohResponse.prototype, "error", void 0);
exports.PharaohResponse = PharaohResponse = __decorate([
    (0, graphql_1.ObjectType)()
], PharaohResponse);
let PharaohsResponse = class PharaohsResponse {
    statusCode;
    message;
    pharaohs;
    error;
    totalCount;
    page;
    limit;
};
exports.PharaohsResponse = PharaohsResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], PharaohsResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], PharaohsResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)(() => [Pharaoh], { nullable: true }),
    __metadata("design:type", Array)
], PharaohsResponse.prototype, "pharaohs", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], PharaohsResponse.prototype, "error", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], PharaohsResponse.prototype, "totalCount", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], PharaohsResponse.prototype, "page", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], PharaohsResponse.prototype, "limit", void 0);
exports.PharaohsResponse = PharaohsResponse = __decorate([
    (0, graphql_1.ObjectType)()
], PharaohsResponse);
let DeletePharaohResponse = class DeletePharaohResponse {
    statusCode;
    message;
    success;
    error;
};
exports.DeletePharaohResponse = DeletePharaohResponse;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], DeletePharaohResponse.prototype, "statusCode", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], DeletePharaohResponse.prototype, "message", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], DeletePharaohResponse.prototype, "success", void 0);
__decorate([
    (0, graphql_1.Field)(() => shared_types_1.ErrorDetails, { nullable: true }),
    __metadata("design:type", shared_types_1.ErrorDetails)
], DeletePharaohResponse.prototype, "error", void 0);
exports.DeletePharaohResponse = DeletePharaohResponse = __decorate([
    (0, graphql_1.ObjectType)()
], DeletePharaohResponse);
let ArtifactInput = class ArtifactInput {
    name;
    museum;
    description;
};
exports.ArtifactInput = ArtifactInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ArtifactInput.prototype, "name", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ArtifactInput.prototype, "museum", void 0);
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], ArtifactInput.prototype, "description", void 0);
exports.ArtifactInput = ArtifactInput = __decorate([
    (0, graphql_1.InputType)()
], ArtifactInput);
let TraitsInput = class TraitsInput {
    leadership;
    military;
    diplomacy;
    wisdom;
    charisma;
};
exports.TraitsInput = TraitsInput;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TraitsInput.prototype, "leadership", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TraitsInput.prototype, "military", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TraitsInput.prototype, "diplomacy", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TraitsInput.prototype, "wisdom", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], TraitsInput.prototype, "charisma", void 0);
exports.TraitsInput = TraitsInput = __decorate([
    (0, graphql_1.InputType)()
], TraitsInput);
let CreatePharaohInput = class CreatePharaohInput {
    name;
    birthName;
    throneName;
    epithet;
    dynasty;
    period;
    reignStart;
    reignEnd;
    lengthOfReignYears;
    predecessorId;
    successorId;
    father;
    mother;
    consorts;
    childrenCount;
    notableChildren;
    capital;
    majorAchievements;
    militaryCampaigns;
    buildingProjects;
    politicalStyle;
    divineAssociation;
    templeAffiliations;
    religiousReforms;
    pharaohAsGod;
    burialSite;
    tombDiscovered;
    discoveryYear;
    tombGuardian;
    funeraryText;
    famousArtifacts;
    treasureStatus;
    imageUrl;
    statueCount;
    mummyLocation;
    audioNarrationUrl;
    videoDocumentaryUrl;
    popularityScore;
    userRating;
    unlockInGame;
    rarity;
    traits;
    source;
    verified;
    language;
};
exports.CreatePharaohInput = CreatePharaohInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "name", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "birthName", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "throneName", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "epithet", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], CreatePharaohInput.prototype, "dynasty", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "period", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], CreatePharaohInput.prototype, "reignStart", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], CreatePharaohInput.prototype, "reignEnd", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], CreatePharaohInput.prototype, "lengthOfReignYears", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "predecessorId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "successorId", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "father", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "mother", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], CreatePharaohInput.prototype, "consorts", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], CreatePharaohInput.prototype, "childrenCount", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], CreatePharaohInput.prototype, "notableChildren", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "capital", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], CreatePharaohInput.prototype, "majorAchievements", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], CreatePharaohInput.prototype, "militaryCampaigns", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], CreatePharaohInput.prototype, "buildingProjects", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "politicalStyle", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], CreatePharaohInput.prototype, "divineAssociation", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], CreatePharaohInput.prototype, "templeAffiliations", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "religiousReforms", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], CreatePharaohInput.prototype, "pharaohAsGod", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "burialSite", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], CreatePharaohInput.prototype, "tombDiscovered", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], CreatePharaohInput.prototype, "discoveryYear", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "tombGuardian", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "funeraryText", void 0);
__decorate([
    (0, graphql_1.Field)(() => [ArtifactInput], { nullable: true }),
    __metadata("design:type", Array)
], CreatePharaohInput.prototype, "famousArtifacts", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "treasureStatus", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "imageUrl", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], CreatePharaohInput.prototype, "statueCount", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "mummyLocation", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "audioNarrationUrl", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "videoDocumentaryUrl", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float, { nullable: true }),
    __metadata("design:type", Number)
], CreatePharaohInput.prototype, "popularityScore", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float, { nullable: true }),
    __metadata("design:type", Number)
], CreatePharaohInput.prototype, "userRating", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], CreatePharaohInput.prototype, "unlockInGame", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "rarity", void 0);
__decorate([
    (0, graphql_1.Field)(() => TraitsInput, { nullable: true }),
    __metadata("design:type", TraitsInput)
], CreatePharaohInput.prototype, "traits", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "source", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", Boolean)
], CreatePharaohInput.prototype, "verified", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], CreatePharaohInput.prototype, "language", void 0);
exports.CreatePharaohInput = CreatePharaohInput = __decorate([
    (0, graphql_1.InputType)()
], CreatePharaohInput);
let UpdatePharaohInput = class UpdatePharaohInput {
    id;
    pharaoh;
};
exports.UpdatePharaohInput = UpdatePharaohInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdatePharaohInput.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(() => CreatePharaohInput),
    __metadata("design:type", CreatePharaohInput)
], UpdatePharaohInput.prototype, "pharaoh", void 0);
exports.UpdatePharaohInput = UpdatePharaohInput = __decorate([
    (0, graphql_1.InputType)()
], UpdatePharaohInput);
let GetAllPharaohsInput = class GetAllPharaohsInput {
    page;
    limit;
    sortBy;
    order;
};
exports.GetAllPharaohsInput = GetAllPharaohsInput;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], GetAllPharaohsInput.prototype, "page", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int, { nullable: true }),
    __metadata("design:type", Number)
], GetAllPharaohsInput.prototype, "limit", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], GetAllPharaohsInput.prototype, "sortBy", void 0);
__decorate([
    (0, graphql_1.Field)({ nullable: true }),
    __metadata("design:type", String)
], GetAllPharaohsInput.prototype, "order", void 0);
exports.GetAllPharaohsInput = GetAllPharaohsInput = __decorate([
    (0, graphql_1.InputType)()
], GetAllPharaohsInput);
let GetPharaohsByDynastyInput = class GetPharaohsByDynastyInput {
    dynasty;
};
exports.GetPharaohsByDynastyInput = GetPharaohsByDynastyInput;
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Int),
    __metadata("design:type", Number)
], GetPharaohsByDynastyInput.prototype, "dynasty", void 0);
exports.GetPharaohsByDynastyInput = GetPharaohsByDynastyInput = __decorate([
    (0, graphql_1.InputType)()
], GetPharaohsByDynastyInput);
let GetPharaohsByPeriodInput = class GetPharaohsByPeriodInput {
    period;
};
exports.GetPharaohsByPeriodInput = GetPharaohsByPeriodInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetPharaohsByPeriodInput.prototype, "period", void 0);
exports.GetPharaohsByPeriodInput = GetPharaohsByPeriodInput = __decorate([
    (0, graphql_1.InputType)()
], GetPharaohsByPeriodInput);
let SearchPharaohsInput = class SearchPharaohsInput {
    query;
    fields;
};
exports.SearchPharaohsInput = SearchPharaohsInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], SearchPharaohsInput.prototype, "query", void 0);
__decorate([
    (0, graphql_1.Field)(() => [String], { nullable: true }),
    __metadata("design:type", Array)
], SearchPharaohsInput.prototype, "fields", void 0);
exports.SearchPharaohsInput = SearchPharaohsInput = __decorate([
    (0, graphql_1.InputType)()
], SearchPharaohsInput);
let GetPharaohsByRarityInput = class GetPharaohsByRarityInput {
    rarity;
};
exports.GetPharaohsByRarityInput = GetPharaohsByRarityInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], GetPharaohsByRarityInput.prototype, "rarity", void 0);
exports.GetPharaohsByRarityInput = GetPharaohsByRarityInput = __decorate([
    (0, graphql_1.InputType)()
], GetPharaohsByRarityInput);
let UpdatePharaohRatingInput = class UpdatePharaohRatingInput {
    id;
    rating;
};
exports.UpdatePharaohRatingInput = UpdatePharaohRatingInput;
__decorate([
    (0, graphql_1.Field)(),
    __metadata("design:type", String)
], UpdatePharaohRatingInput.prototype, "id", void 0);
__decorate([
    (0, graphql_1.Field)(() => graphql_1.Float),
    __metadata("design:type", Number)
], UpdatePharaohRatingInput.prototype, "rating", void 0);
exports.UpdatePharaohRatingInput = UpdatePharaohRatingInput = __decorate([
    (0, graphql_1.InputType)()
], UpdatePharaohRatingInput);
//# sourceMappingURL=pharaoh.types.js.map