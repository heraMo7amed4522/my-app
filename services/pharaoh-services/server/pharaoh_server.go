package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	pb "pharaoh-services/proto"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PharaohServer struct {
	pb.UnimplementedPharaohServiceServer
	db *sql.DB
}

// NewPharaohServer creates a new pharaoh server instance
func NewPharaohServer() (*PharaohServer, error) {
	db, err := NewDatabaseConnection()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	return &PharaohServer{
		db: db,
	}, nil
}

// Close closes the database connection
func (s *PharaohServer) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// GetPharaohByID retrieves a pharaoh by ID
func (s *PharaohServer) GetPharaohByID(ctx context.Context, req *pb.GetPharaohByIDRequest) (*pb.GetPharaohByIDResponse, error) {
	log.Printf("Getting pharaoh by ID: %s", req.Id)

	if req.Id == "" {
		return &pb.GetPharaohByIDResponse{
			StatusCode: 400,
			Message:    "Pharaoh ID is required",
			Result: &pb.GetPharaohByIDResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Pharaoh ID is required",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	query := `
		SELECT id, name, birth_name, throne_name, epithet, dynasty, period, 
		       reign_start, reign_end, length_of_reign_years, predecessor_id, successor_id,
		       father, mother, consorts, children_count, notable_children,
		       capital, major_achievements, military_campaigns, building_projects, political_style,
		       divine_association, temple_affiliations, religious_reforms, pharaoh_as_god,
		       burial_site, tomb_discovered, discovery_year, tomb_guardian, funerary_text,
		       famous_artifacts, treasure_status, image_url, statue_count, mummy_location,
		       audio_narration_url, video_documentary_url, popularity_score, user_rating,
		       unlock_in_game, rarity, traits, source, verified, language, created_at, updated_at
		FROM pharaohs WHERE id = $1
	`

	row := s.db.QueryRowContext(ctx, query, req.Id)

	pharaoh := &pb.Pharaoh{}
	var consorts, notableChildren, majorAchievements, militaryCampaigns, buildingProjects []byte
	var divineAssociation, templeAffiliations, famousArtifacts, traits []byte
	var predecessorID, successorID, birthName, throneName, epithet sql.NullString
	var dynasty, reignStart, reignEnd, lengthOfReignYears, childrenCount sql.NullInt32
	var period, father, mother, capital, politicalStyle, religiousReforms sql.NullString
	var pharaohAsGod, tombDiscovered, unlockInGame, verified sql.NullBool
	var discoveryYear, statueCount sql.NullInt32
	var burialSite, tombGuardian, funeraryText, treasureStatus, imageUrl sql.NullString
	var mummyLocation, audioNarrationUrl, videoDocumentaryUrl, rarity, source, language sql.NullString
	var popularityScore, userRating sql.NullFloat64
	var createdAt, updatedAt time.Time

	err := row.Scan(
		&pharaoh.Id, &pharaoh.Name, &birthName, &throneName, &epithet,
		&dynasty, &period, &reignStart, &reignEnd, &lengthOfReignYears,
		&predecessorID, &successorID, &father, &mother, &consorts,
		&childrenCount, &notableChildren, &capital, &majorAchievements,
		&militaryCampaigns, &buildingProjects, &politicalStyle,
		&divineAssociation, &templeAffiliations, &religiousReforms, &pharaohAsGod,
		&burialSite, &tombDiscovered, &discoveryYear, &tombGuardian, &funeraryText,
		&famousArtifacts, &treasureStatus, &imageUrl, &statueCount, &mummyLocation,
		&audioNarrationUrl, &videoDocumentaryUrl, &popularityScore, &userRating,
		&unlockInGame, &rarity, &traits, &source, &verified, &language,
		&createdAt, &updatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetPharaohByIDResponse{
				StatusCode: 404,
				Message:    "Pharaoh not found",
				Result: &pb.GetPharaohByIDResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "Pharaoh not found",
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, status.Errorf(codes.Internal, "Database error: %v", err)
	}

	// Handle nullable fields
	if birthName.Valid {
		pharaoh.BirthName = birthName.String
	}
	if throneName.Valid {
		pharaoh.ThroneName = throneName.String
	}
	if epithet.Valid {
		pharaoh.Epithet = epithet.String
	}
	if dynasty.Valid {
		pharaoh.Dynasty = dynasty.Int32
	}
	if period.Valid {
		pharaoh.Period = period.String
	}
	if reignStart.Valid {
		pharaoh.ReignStart = reignStart.Int32
	}
	if reignEnd.Valid {
		pharaoh.ReignEnd = reignEnd.Int32
	}
	if lengthOfReignYears.Valid {
		pharaoh.LengthOfReignYears = lengthOfReignYears.Int32
	}
	if predecessorID.Valid {
		pharaoh.PredecessorId = predecessorID.String
	}
	if successorID.Valid {
		pharaoh.SuccessorId = successorID.String
	}
	if father.Valid {
		pharaoh.Father = father.String
	}
	if mother.Valid {
		pharaoh.Mother = mother.String
	}
	if childrenCount.Valid {
		pharaoh.ChildrenCount = childrenCount.Int32
	}
	if capital.Valid {
		pharaoh.Capital = capital.String
	}
	if politicalStyle.Valid {
		pharaoh.PoliticalStyle = politicalStyle.String
	}
	if religiousReforms.Valid {
		pharaoh.ReligiousReforms = religiousReforms.String
	}
	if pharaohAsGod.Valid {
		pharaoh.PharaohAsGod = pharaohAsGod.Bool
	}
	if burialSite.Valid {
		pharaoh.BurialSite = burialSite.String
	}
	if tombDiscovered.Valid {
		pharaoh.TombDiscovered = tombDiscovered.Bool
	}
	if discoveryYear.Valid {
		pharaoh.DiscoveryYear = discoveryYear.Int32
	}
	if tombGuardian.Valid {
		pharaoh.TombGuardian = tombGuardian.String
	}
	if funeraryText.Valid {
		pharaoh.FuneraryText = funeraryText.String
	}
	if treasureStatus.Valid {
		pharaoh.TreasureStatus = treasureStatus.String
	}
	if imageUrl.Valid {
		pharaoh.ImageUrl = imageUrl.String
	}
	if statueCount.Valid {
		pharaoh.StatueCount = statueCount.Int32
	}
	if mummyLocation.Valid {
		pharaoh.MummyLocation = mummyLocation.String
	}
	if audioNarrationUrl.Valid {
		pharaoh.AudioNarrationUrl = audioNarrationUrl.String
	}
	if videoDocumentaryUrl.Valid {
		pharaoh.VideoDocumentaryUrl = videoDocumentaryUrl.String
	}
	if popularityScore.Valid {
		pharaoh.PopularityScore = float32(popularityScore.Float64)
	}
	if userRating.Valid {
		pharaoh.UserRating = float32(userRating.Float64)
	}
	if unlockInGame.Valid {
		pharaoh.UnlockInGame = unlockInGame.Bool
	}
	if rarity.Valid {
		pharaoh.Rarity = rarity.String
	}
	if source.Valid {
		pharaoh.Source = source.String
	}
	if verified.Valid {
		pharaoh.Verified = verified.Bool
	}
	if language.Valid {
		pharaoh.Language = language.String
	}

	// Parse JSON fields
	if len(consorts) > 0 {
		var consortsList []string
		if err := json.Unmarshal(consorts, &consortsList); err == nil {
			pharaoh.Consorts = consortsList
		}
	}

	if len(notableChildren) > 0 {
		var childrenList []string
		if err := json.Unmarshal(notableChildren, &childrenList); err == nil {
			pharaoh.NotableChildren = childrenList
		}
	}

	if len(majorAchievements) > 0 {
		var achievementsList []string
		if err := json.Unmarshal(majorAchievements, &achievementsList); err == nil {
			pharaoh.MajorAchievements = achievementsList
		}
	}

	if len(militaryCampaigns) > 0 {
		var campaignsList []string
		if err := json.Unmarshal(militaryCampaigns, &campaignsList); err == nil {
			pharaoh.MilitaryCampaigns = campaignsList
		}
	}

	if len(buildingProjects) > 0 {
		var projectsList []string
		if err := json.Unmarshal(buildingProjects, &projectsList); err == nil {
			pharaoh.BuildingProjects = projectsList
		}
	}

	if len(divineAssociation) > 0 {
		var divineList []string
		if err := json.Unmarshal(divineAssociation, &divineList); err == nil {
			pharaoh.DivineAssociation = divineList
		}
	}

	if len(templeAffiliations) > 0 {
		var templeList []string
		if err := json.Unmarshal(templeAffiliations, &templeList); err == nil {
			pharaoh.TempleAffiliations = templeList
		}
	}

	if len(famousArtifacts) > 0 {
		var artifactsList []map[string]interface{}
		if err := json.Unmarshal(famousArtifacts, &artifactsList); err == nil {
			for _, artifact := range artifactsList {
				pbArtifact := &pb.Artifact{}
				if name, ok := artifact["name"].(string); ok {
					pbArtifact.Name = name
				}
				if museum, ok := artifact["museum"].(string); ok {
					pbArtifact.Museum = museum
				}
				if description, ok := artifact["description"].(string); ok {
					pbArtifact.Description = description
				}
				pharaoh.FamousArtifacts = append(pharaoh.FamousArtifacts, pbArtifact)
			}
		}
	}

	if len(traits) > 0 {
		var traitsMap map[string]interface{}
		if err := json.Unmarshal(traits, &traitsMap); err == nil {
			pharaoh.Traits = &pb.Traits{}
			if leadership, ok := traitsMap["leadership"].(float64); ok {
				pharaoh.Traits.Leadership = int32(leadership)
			}
			if military, ok := traitsMap["military"].(float64); ok {
				pharaoh.Traits.Military = int32(military)
			}
			if diplomacy, ok := traitsMap["diplomacy"].(float64); ok {
				pharaoh.Traits.Diplomacy = int32(diplomacy)
			}
			if wisdom, ok := traitsMap["wisdom"].(float64); ok {
				pharaoh.Traits.Wisdom = int32(wisdom)
			}
			if charisma, ok := traitsMap["charisma"].(float64); ok {
				pharaoh.Traits.Charisma = int32(charisma)
			}
		}
	}

	pharaoh.CreatedAt = createdAt.Format(time.RFC3339)
	pharaoh.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &pb.GetPharaohByIDResponse{
		StatusCode: 200,
		Message:    "Pharaoh retrieved successfully",
		Result: &pb.GetPharaohByIDResponse_Pharaoh{
			Pharaoh: pharaoh,
		},
	}, nil
}

// GetAllPharaohs retrieves all pharaohs with pagination
func (s *PharaohServer) GetAllPharaohs(ctx context.Context, req *pb.GetAllPharaohsRequest) (*pb.GetAllPharaohsResponse, error) {
	log.Printf("Getting all pharaohs with page: %d, limit: %d", req.Page, req.Limit)

	// Set defaults
	page := req.Page
	if page <= 0 {
		page = 1
	}
	limit := req.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100 // Max limit
	}

	sortBy := req.SortBy
	if sortBy == "" {
		sortBy = "dynasty"
	}

	order := req.Order
	if order == "" {
		order = "asc"
	}

	// Validate sort fields
	allowedSortFields := map[string]bool{
		"dynasty":          true,
		"reign_start":      true,
		"popularity_score": true,
		"name":             true,
		"created_at":       true,
	}

	if !allowedSortFields[sortBy] {
		sortBy = "dynasty"
	}

	if order != "asc" && order != "desc" {
		order = "asc"
	}

	offset := (page - 1) * limit

	// Count total records
	countQuery := "SELECT COUNT(*) FROM pharaohs"
	var totalCount int32
	err := s.db.QueryRowContext(ctx, countQuery).Scan(&totalCount)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to count pharaohs: %v", err)
	}

	// Get pharaohs with pagination
	query := fmt.Sprintf(`
		SELECT id, name, birth_name, throne_name, epithet, dynasty, period,
		       reign_start, reign_end, length_of_reign_years, popularity_score, user_rating,
		       rarity, image_url, verified, created_at, updated_at
		FROM pharaohs 
		ORDER BY %s %s 
		LIMIT $1 OFFSET $2
	`, sortBy, order)

	rows, err := s.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to query pharaohs: %v", err)
	}
	defer rows.Close()

	var pharaohs []*pb.Pharaoh
	for rows.Next() {
		pharaoh := &pb.Pharaoh{}
		var birthName, throneName, epithet, period sql.NullString
		var dynasty, reignStart, reignEnd, lengthOfReignYears sql.NullInt32
		var popularityScore, userRating sql.NullFloat64
		var rarity, imageUrl sql.NullString
		var verified sql.NullBool
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&pharaoh.Id, &pharaoh.Name, &birthName, &throneName, &epithet,
			&dynasty, &period, &reignStart, &reignEnd, &lengthOfReignYears,
			&popularityScore, &userRating, &rarity, &imageUrl, &verified,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to scan pharaoh: %v", err)
		}

		// Handle nullable fields
		if birthName.Valid {
			pharaoh.BirthName = birthName.String
		}
		if throneName.Valid {
			pharaoh.ThroneName = throneName.String
		}
		if epithet.Valid {
			pharaoh.Epithet = epithet.String
		}
		if dynasty.Valid {
			pharaoh.Dynasty = dynasty.Int32
		}
		if period.Valid {
			pharaoh.Period = period.String
		}
		if reignStart.Valid {
			pharaoh.ReignStart = reignStart.Int32
		}
		if reignEnd.Valid {
			pharaoh.ReignEnd = reignEnd.Int32
		}
		if lengthOfReignYears.Valid {
			pharaoh.LengthOfReignYears = lengthOfReignYears.Int32
		}
		if popularityScore.Valid {
			pharaoh.PopularityScore = float32(popularityScore.Float64)
		}
		if userRating.Valid {
			pharaoh.UserRating = float32(userRating.Float64)
		}
		if rarity.Valid {
			pharaoh.Rarity = rarity.String
		}
		if imageUrl.Valid {
			pharaoh.ImageUrl = imageUrl.String
		}
		if verified.Valid {
			pharaoh.Verified = verified.Bool
		}

		pharaoh.CreatedAt = createdAt.Format(time.RFC3339)
		pharaoh.UpdatedAt = updatedAt.Format(time.RFC3339)

		pharaohs = append(pharaohs, pharaoh)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Error iterating pharaohs: %v", err)
	}

	return &pb.GetAllPharaohsResponse{
		StatusCode: 200,
		Message:    "Pharaohs retrieved successfully",
		Result: &pb.GetAllPharaohsResponse_Pharaohs{
			Pharaohs: &pb.PharaohList{
				Pharaohs:   pharaohs,
				TotalCount: totalCount,
				Page:       page,
				Limit:      limit,
			},
		},
	}, nil
}

// GetPharaohsByDynasty retrieves pharaohs by dynasty
func (s *PharaohServer) GetPharaohsByDynasty(ctx context.Context, req *pb.GetPharaohsByDynastyRequest) (*pb.GetPharaohsByDynastyResponse, error) {
	log.Printf("Getting pharaohs by dynasty: %d", req.Dynasty)

	if req.Dynasty <= 0 {
		return &pb.GetPharaohsByDynastyResponse{
			StatusCode: 400,
			Message:    "Dynasty must be a positive number",
			Result: &pb.GetPharaohsByDynastyResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Dynasty must be a positive number",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	query := `
		SELECT id, name, birth_name, throne_name, epithet, dynasty, period,
		       reign_start, reign_end, length_of_reign_years, popularity_score, user_rating,
		       rarity, image_url, verified, created_at, updated_at
		FROM pharaohs 
		WHERE dynasty = $1
		ORDER BY reign_start ASC
	`

	rows, err := s.db.QueryContext(ctx, query, req.Dynasty)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to query pharaohs by dynasty: %v", err)
	}
	defer rows.Close()

	var pharaohs []*pb.Pharaoh
	for rows.Next() {
		pharaoh := &pb.Pharaoh{}
		var birthName, throneName, epithet, period sql.NullString
		var dynasty, reignStart, reignEnd, lengthOfReignYears sql.NullInt32
		var popularityScore, userRating sql.NullFloat64
		var rarity, imageUrl sql.NullString
		var verified sql.NullBool
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&pharaoh.Id, &pharaoh.Name, &birthName, &throneName, &epithet,
			&dynasty, &period, &reignStart, &reignEnd, &lengthOfReignYears,
			&popularityScore, &userRating, &rarity, &imageUrl, &verified,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to scan pharaoh: %v", err)
		}

		// Handle nullable fields
		if birthName.Valid {
			pharaoh.BirthName = birthName.String
		}
		if throneName.Valid {
			pharaoh.ThroneName = throneName.String
		}
		if epithet.Valid {
			pharaoh.Epithet = epithet.String
		}
		if dynasty.Valid {
			pharaoh.Dynasty = dynasty.Int32
		}
		if period.Valid {
			pharaoh.Period = period.String
		}
		if reignStart.Valid {
			pharaoh.ReignStart = reignStart.Int32
		}
		if reignEnd.Valid {
			pharaoh.ReignEnd = reignEnd.Int32
		}
		if lengthOfReignYears.Valid {
			pharaoh.LengthOfReignYears = lengthOfReignYears.Int32
		}
		if popularityScore.Valid {
			pharaoh.PopularityScore = float32(popularityScore.Float64)
		}
		if userRating.Valid {
			pharaoh.UserRating = float32(userRating.Float64)
		}
		if rarity.Valid {
			pharaoh.Rarity = rarity.String
		}
		if imageUrl.Valid {
			pharaoh.ImageUrl = imageUrl.String
		}
		if verified.Valid {
			pharaoh.Verified = verified.Bool
		}

		pharaoh.CreatedAt = createdAt.Format(time.RFC3339)
		pharaoh.UpdatedAt = updatedAt.Format(time.RFC3339)

		pharaohs = append(pharaohs, pharaoh)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Error iterating pharaohs: %v", err)
	}

	return &pb.GetPharaohsByDynastyResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Found %d pharaohs from dynasty %d", len(pharaohs), req.Dynasty),
		Result: &pb.GetPharaohsByDynastyResponse_Pharaohs{
			Pharaohs: &pb.PharaohList{
				Pharaohs:   pharaohs,
				TotalCount: int32(len(pharaohs)),
				Page:       1,
				Limit:      int32(len(pharaohs)),
			},
		},
	}, nil
}

// GetPharaohsByPeriod retrieves pharaohs by historical period
func (s *PharaohServer) GetPharaohsByPeriod(ctx context.Context, req *pb.GetPharaohsByPeriodRequest) (*pb.GetPharaohsByPeriodResponse, error) {
	log.Printf("Getting pharaohs by period: %s", req.Period)

	if req.Period == "" {
		return &pb.GetPharaohsByPeriodResponse{
			StatusCode: 400,
			Message:    "Period is required",
			Result: &pb.GetPharaohsByPeriodResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Period is required",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	query := `
		SELECT id, name, birth_name, throne_name, epithet, dynasty, period,
		       reign_start, reign_end, length_of_reign_years, popularity_score, user_rating,
		       rarity, image_url, verified, created_at, updated_at
		FROM pharaohs 
		WHERE LOWER(period) = LOWER($1)
		ORDER BY dynasty ASC, reign_start ASC
	`

	rows, err := s.db.QueryContext(ctx, query, req.Period)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to query pharaohs by period: %v", err)
	}
	defer rows.Close()

	var pharaohs []*pb.Pharaoh
	for rows.Next() {
		pharaoh := &pb.Pharaoh{}
		var birthName, throneName, epithet, period sql.NullString
		var dynasty, reignStart, reignEnd, lengthOfReignYears sql.NullInt32
		var popularityScore, userRating sql.NullFloat64
		var rarity, imageUrl sql.NullString
		var verified sql.NullBool
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&pharaoh.Id, &pharaoh.Name, &birthName, &throneName, &epithet,
			&dynasty, &period, &reignStart, &reignEnd, &lengthOfReignYears,
			&popularityScore, &userRating, &rarity, &imageUrl, &verified,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to scan pharaoh: %v", err)
		}

		// Handle nullable fields (same as above)
		if birthName.Valid {
			pharaoh.BirthName = birthName.String
		}
		if throneName.Valid {
			pharaoh.ThroneName = throneName.String
		}
		if epithet.Valid {
			pharaoh.Epithet = epithet.String
		}
		if dynasty.Valid {
			pharaoh.Dynasty = dynasty.Int32
		}
		if period.Valid {
			pharaoh.Period = period.String
		}
		if reignStart.Valid {
			pharaoh.ReignStart = reignStart.Int32
		}
		if reignEnd.Valid {
			pharaoh.ReignEnd = reignEnd.Int32
		}
		if lengthOfReignYears.Valid {
			pharaoh.LengthOfReignYears = lengthOfReignYears.Int32
		}
		if popularityScore.Valid {
			pharaoh.PopularityScore = float32(popularityScore.Float64)
		}
		if userRating.Valid {
			pharaoh.UserRating = float32(userRating.Float64)
		}
		if rarity.Valid {
			pharaoh.Rarity = rarity.String
		}
		if imageUrl.Valid {
			pharaoh.ImageUrl = imageUrl.String
		}
		if verified.Valid {
			pharaoh.Verified = verified.Bool
		}

		pharaoh.CreatedAt = createdAt.Format(time.RFC3339)
		pharaoh.UpdatedAt = updatedAt.Format(time.RFC3339)

		pharaohs = append(pharaohs, pharaoh)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Error iterating pharaohs: %v", err)
	}

	return &pb.GetPharaohsByPeriodResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Found %d pharaohs from %s period", len(pharaohs), req.Period),
		Result: &pb.GetPharaohsByPeriodResponse_Pharaohs{
			Pharaohs: &pb.PharaohList{
				Pharaohs:   pharaohs,
				TotalCount: int32(len(pharaohs)),
				Page:       1,
				Limit:      int32(len(pharaohs)),
			},
		},
	}, nil
}

// SearchPharaohs searches pharaohs by query across specified fields
func (s *PharaohServer) SearchPharaohs(ctx context.Context, req *pb.SearchPharaohsRequest) (*pb.SearchPharaohsResponse, error) {
	log.Printf("Searching pharaohs with query: %s", req.Query)

	if req.Query == "" {
		return &pb.SearchPharaohsResponse{
			StatusCode: 400,
			Message:    "Search query is required",
			Result: &pb.SearchPharaohsResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Search query is required",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Default search fields if none specified
	searchFields := req.Fields
	if len(searchFields) == 0 {
		searchFields = []string{"name", "epithet", "major_achievements"}
	}

	// Build dynamic WHERE clause based on fields
	var whereClauses []string
	var args []interface{}
	argIndex := 1

	for _, field := range searchFields {
		switch field {
		case "name":
			whereClauses = append(whereClauses, fmt.Sprintf("LOWER(name) LIKE LOWER($%d)", argIndex))
			args = append(args, "%"+req.Query+"%")
			argIndex++
		case "epithet":
			whereClauses = append(whereClauses, fmt.Sprintf("LOWER(epithet) LIKE LOWER($%d)", argIndex))
			args = append(args, "%"+req.Query+"%")
			argIndex++
		case "achievements":
			whereClauses = append(whereClauses, fmt.Sprintf("LOWER(major_achievements::text) LIKE LOWER($%d)", argIndex))
			args = append(args, "%"+req.Query+"%")
			argIndex++
		case "birth_name":
			whereClauses = append(whereClauses, fmt.Sprintf("LOWER(birth_name) LIKE LOWER($%d)", argIndex))
			args = append(args, "%"+req.Query+"%")
			argIndex++
		case "throne_name":
			whereClauses = append(whereClauses, fmt.Sprintf("LOWER(throne_name) LIKE LOWER($%d)", argIndex))
			args = append(args, "%"+req.Query+"%")
			argIndex++
		}
	}

	if len(whereClauses) == 0 {
		return &pb.SearchPharaohsResponse{
			StatusCode: 400,
			Message:    "No valid search fields specified",
			Result: &pb.SearchPharaohsResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "No valid search fields specified",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	query := fmt.Sprintf(`
		SELECT id, name, birth_name, throne_name, epithet, dynasty, period,
		       reign_start, reign_end, length_of_reign_years, popularity_score, user_rating,
		       rarity, image_url, verified, created_at, updated_at
		FROM pharaohs 
		WHERE %s
		ORDER BY popularity_score DESC, name ASC
	`, strings.Join(whereClauses, " OR "))

	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to search pharaohs: %v", err)
	}
	defer rows.Close()

	var pharaohs []*pb.Pharaoh
	for rows.Next() {
		pharaoh := &pb.Pharaoh{}
		var birthName, throneName, epithet, period sql.NullString
		var dynasty, reignStart, reignEnd, lengthOfReignYears sql.NullInt32
		var popularityScore, userRating sql.NullFloat64
		var rarity, imageUrl sql.NullString
		var verified sql.NullBool
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&pharaoh.Id, &pharaoh.Name, &birthName, &throneName, &epithet,
			&dynasty, &period, &reignStart, &reignEnd, &lengthOfReignYears,
			&popularityScore, &userRating, &rarity, &imageUrl, &verified,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to scan pharaoh: %v", err)
		}

		// Handle nullable fields (same pattern as above)
		if birthName.Valid {
			pharaoh.BirthName = birthName.String
		}
		if throneName.Valid {
			pharaoh.ThroneName = throneName.String
		}
		if epithet.Valid {
			pharaoh.Epithet = epithet.String
		}
		if dynasty.Valid {
			pharaoh.Dynasty = dynasty.Int32
		}
		if period.Valid {
			pharaoh.Period = period.String
		}
		if reignStart.Valid {
			pharaoh.ReignStart = reignStart.Int32
		}
		if reignEnd.Valid {
			pharaoh.ReignEnd = reignEnd.Int32
		}
		if lengthOfReignYears.Valid {
			pharaoh.LengthOfReignYears = lengthOfReignYears.Int32
		}
		if popularityScore.Valid {
			pharaoh.PopularityScore = float32(popularityScore.Float64)
		}
		if userRating.Valid {
			pharaoh.UserRating = float32(userRating.Float64)
		}
		if rarity.Valid {
			pharaoh.Rarity = rarity.String
		}
		if imageUrl.Valid {
			pharaoh.ImageUrl = imageUrl.String
		}
		if verified.Valid {
			pharaoh.Verified = verified.Bool
		}

		pharaoh.CreatedAt = createdAt.Format(time.RFC3339)
		pharaoh.UpdatedAt = updatedAt.Format(time.RFC3339)

		pharaohs = append(pharaohs, pharaoh)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Error iterating search results: %v", err)
	}

	return &pb.SearchPharaohsResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Found %d pharaohs matching '%s'", len(pharaohs), req.Query),
		Result: &pb.SearchPharaohsResponse_Pharaohs{
			Pharaohs: &pb.PharaohList{
				Pharaohs:   pharaohs,
				TotalCount: int32(len(pharaohs)),
				Page:       1,
				Limit:      int32(len(pharaohs)),
			},
		},
	}, nil
}

// CreatePharaoh creates a new pharaoh record
func (s *PharaohServer) CreatePharaoh(ctx context.Context, req *pb.CreatePharaohRequest) (*pb.CreatePharaohResponse, error) {
	log.Printf("Creating new pharaoh: %s", req.Pharaoh.Name)

	if req.Pharaoh == nil {
		return &pb.CreatePharaohResponse{
			StatusCode: 400,
			Message:    "Pharaoh data is required",
			Result: &pb.CreatePharaohResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Pharaoh data is required",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.Pharaoh.Id == "" || req.Pharaoh.Name == "" {
		return &pb.CreatePharaohResponse{
			StatusCode: 400,
			Message:    "Pharaoh ID and name are required",
			Result: &pb.CreatePharaohResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Pharaoh ID and name are required",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Check if pharaoh already exists
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pharaohs WHERE id = $1)", req.Pharaoh.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to check pharaoh existence: %v", err)
	}

	if exists {
		return &pb.CreatePharaohResponse{
			StatusCode: 409,
			Message:    "Pharaoh with this ID already exists",
			Result: &pb.CreatePharaohResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      409,
					Message:   "Pharaoh with this ID already exists",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Prepare JSON fields
	consorts, _ := json.Marshal(req.Pharaoh.Consorts)
	notableChildren, _ := json.Marshal(req.Pharaoh.NotableChildren)
	majorAchievements, _ := json.Marshal(req.Pharaoh.MajorAchievements)
	militaryCampaigns, _ := json.Marshal(req.Pharaoh.MilitaryCampaigns)
	buildingProjects, _ := json.Marshal(req.Pharaoh.BuildingProjects)
	divineAssociation, _ := json.Marshal(req.Pharaoh.DivineAssociation)
	templeAffiliations, _ := json.Marshal(req.Pharaoh.TempleAffiliations)

	// Handle artifacts
	var artifactsJSON []byte
	if len(req.Pharaoh.FamousArtifacts) > 0 {
		artifactsList := make([]map[string]string, len(req.Pharaoh.FamousArtifacts))
		for i, artifact := range req.Pharaoh.FamousArtifacts {
			artifactsList[i] = map[string]string{
				"name":        artifact.Name,
				"museum":      artifact.Museum,
				"description": artifact.Description,
			}
		}
		artifactsJSON, _ = json.Marshal(artifactsList)
	}

	// Handle traits
	var traitsJSON []byte
	if req.Pharaoh.Traits != nil {
		traitsMap := map[string]int32{
			"leadership": req.Pharaoh.Traits.Leadership,
			"military":   req.Pharaoh.Traits.Military,
			"diplomacy":  req.Pharaoh.Traits.Diplomacy,
			"wisdom":     req.Pharaoh.Traits.Wisdom,
			"charisma":   req.Pharaoh.Traits.Charisma,
		}
		traitsJSON, _ = json.Marshal(traitsMap)
	}

	query := `
		INSERT INTO pharaohs (
			id, name, birth_name, throne_name, epithet, dynasty, period,
			reign_start, reign_end, length_of_reign_years, predecessor_id, successor_id,
			father, mother, consorts, children_count, notable_children,
			capital, major_achievements, military_campaigns, building_projects, political_style,
			divine_association, temple_affiliations, religious_reforms, pharaoh_as_god,
			burial_site, tomb_discovered, discovery_year, tomb_guardian, funerary_text,
			famous_artifacts, treasure_status, image_url, statue_count, mummy_location,
			audio_narration_url, video_documentary_url, popularity_score, user_rating,
			unlock_in_game, rarity, traits, source, verified, language
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17,
			$18, $19, $20, $21, $22, $23, $24, $25, $26, $27, $28, $29, $30, $31, $32,
			$33, $34, $35, $36, $37, $38, $39, $40, $41, $42, $43, $44, $45, $46
		) RETURNING created_at, updated_at
	`

	var createdAt, updatedAt time.Time
	err = s.db.QueryRowContext(ctx, query,
		req.Pharaoh.Id, req.Pharaoh.Name, req.Pharaoh.BirthName, req.Pharaoh.ThroneName, req.Pharaoh.Epithet,
		req.Pharaoh.Dynasty, req.Pharaoh.Period, req.Pharaoh.ReignStart, req.Pharaoh.ReignEnd, req.Pharaoh.LengthOfReignYears,
		req.Pharaoh.PredecessorId, req.Pharaoh.SuccessorId, req.Pharaoh.Father, req.Pharaoh.Mother, consorts,
		req.Pharaoh.ChildrenCount, notableChildren, req.Pharaoh.Capital, majorAchievements, militaryCampaigns,
		buildingProjects, req.Pharaoh.PoliticalStyle, divineAssociation, templeAffiliations, req.Pharaoh.ReligiousReforms,
		req.Pharaoh.PharaohAsGod, req.Pharaoh.BurialSite, req.Pharaoh.TombDiscovered, req.Pharaoh.DiscoveryYear,
		req.Pharaoh.TombGuardian, req.Pharaoh.FuneraryText, artifactsJSON, req.Pharaoh.TreasureStatus, req.Pharaoh.ImageUrl,
		req.Pharaoh.StatueCount, req.Pharaoh.MummyLocation, req.Pharaoh.AudioNarrationUrl, req.Pharaoh.VideoDocumentaryUrl,
		req.Pharaoh.PopularityScore, req.Pharaoh.UserRating, req.Pharaoh.UnlockInGame, req.Pharaoh.Rarity,
		traitsJSON, req.Pharaoh.Source, req.Pharaoh.Verified, req.Pharaoh.Language,
	).Scan(&createdAt, &updatedAt)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to create pharaoh: %v", err)
	}

	// Update timestamps in response
	req.Pharaoh.CreatedAt = createdAt.Format(time.RFC3339)
	req.Pharaoh.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &pb.CreatePharaohResponse{
		StatusCode: 201,
		Message:    "Pharaoh created successfully",
		Result: &pb.CreatePharaohResponse_Pharaoh{
			Pharaoh: req.Pharaoh,
		},
	}, nil
}

// UpdatePharaoh updates an existing pharaoh record
func (s *PharaohServer) UpdatePharaoh(ctx context.Context, req *pb.UpdatePharaohRequest) (*pb.UpdatePharaohResponse, error) {
	log.Printf("Updating pharaoh: %s", req.Id)

	if req.Id == "" {
		return &pb.UpdatePharaohResponse{
			StatusCode: 400,
			Message:    "Pharaoh ID is required",
			Result: &pb.UpdatePharaohResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Pharaoh ID is required",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.Pharaoh == nil {
		return &pb.UpdatePharaohResponse{
			StatusCode: 400,
			Message:    "Pharaoh data is required",
			Result: &pb.UpdatePharaohResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Pharaoh data is required",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Check if pharaoh exists
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pharaohs WHERE id = $1)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to check pharaoh existence: %v", err)
	}

	if !exists {
		return &pb.UpdatePharaohResponse{
			StatusCode: 404,
			Message:    "Pharaoh not found",
			Result: &pb.UpdatePharaohResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      404,
					Message:   "Pharaoh not found",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Prepare JSON fields
	consorts, _ := json.Marshal(req.Pharaoh.Consorts)
	notableChildren, _ := json.Marshal(req.Pharaoh.NotableChildren)
	majorAchievements, _ := json.Marshal(req.Pharaoh.MajorAchievements)
	militaryCampaigns, _ := json.Marshal(req.Pharaoh.MilitaryCampaigns)
	buildingProjects, _ := json.Marshal(req.Pharaoh.BuildingProjects)
	divineAssociation, _ := json.Marshal(req.Pharaoh.DivineAssociation)
	templeAffiliations, _ := json.Marshal(req.Pharaoh.TempleAffiliations)

	// Handle artifacts
	var artifactsJSON []byte
	if len(req.Pharaoh.FamousArtifacts) > 0 {
		artifactsList := make([]map[string]string, len(req.Pharaoh.FamousArtifacts))
		for i, artifact := range req.Pharaoh.FamousArtifacts {
			artifactsList[i] = map[string]string{
				"name":        artifact.Name,
				"museum":      artifact.Museum,
				"description": artifact.Description,
			}
		}
		artifactsJSON, _ = json.Marshal(artifactsList)
	}

	// Handle traits
	var traitsJSON []byte
	if req.Pharaoh.Traits != nil {
		traitsMap := map[string]int32{
			"leadership": req.Pharaoh.Traits.Leadership,
			"military":   req.Pharaoh.Traits.Military,
			"diplomacy":  req.Pharaoh.Traits.Diplomacy,
			"wisdom":     req.Pharaoh.Traits.Wisdom,
			"charisma":   req.Pharaoh.Traits.Charisma,
		}
		traitsJSON, _ = json.Marshal(traitsMap)
	}

	query := `
		UPDATE pharaohs SET
			name = $2, birth_name = $3, throne_name = $4, epithet = $5, dynasty = $6, period = $7,
			reign_start = $8, reign_end = $9, length_of_reign_years = $10, predecessor_id = $11, successor_id = $12,
			father = $13, mother = $14, consorts = $15, children_count = $16, notable_children = $17,
			capital = $18, major_achievements = $19, military_campaigns = $20, building_projects = $21, political_style = $22,
			divine_association = $23, temple_affiliations = $24, religious_reforms = $25, pharaoh_as_god = $26,
			burial_site = $27, tomb_discovered = $28, discovery_year = $29, tomb_guardian = $30, funerary_text = $31,
			famous_artifacts = $32, treasure_status = $33, image_url = $34, statue_count = $35, mummy_location = $36,
			audio_narration_url = $37, video_documentary_url = $38, popularity_score = $39, user_rating = $40,
			unlock_in_game = $41, rarity = $42, traits = $43, source = $44, verified = $45, language = $46,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING updated_at
	`

	var updatedAt time.Time
	err = s.db.QueryRowContext(ctx, query,
		req.Id, req.Pharaoh.Name, req.Pharaoh.BirthName, req.Pharaoh.ThroneName, req.Pharaoh.Epithet,
		req.Pharaoh.Dynasty, req.Pharaoh.Period, req.Pharaoh.ReignStart, req.Pharaoh.ReignEnd, req.Pharaoh.LengthOfReignYears,
		req.Pharaoh.PredecessorId, req.Pharaoh.SuccessorId, req.Pharaoh.Father, req.Pharaoh.Mother, consorts,
		req.Pharaoh.ChildrenCount, notableChildren, req.Pharaoh.Capital, majorAchievements, militaryCampaigns,
		buildingProjects, req.Pharaoh.PoliticalStyle, divineAssociation, templeAffiliations, req.Pharaoh.ReligiousReforms,
		req.Pharaoh.PharaohAsGod, req.Pharaoh.BurialSite, req.Pharaoh.TombDiscovered, req.Pharaoh.DiscoveryYear,
		req.Pharaoh.TombGuardian, req.Pharaoh.FuneraryText, artifactsJSON, req.Pharaoh.TreasureStatus, req.Pharaoh.ImageUrl,
		req.Pharaoh.StatueCount, req.Pharaoh.MummyLocation, req.Pharaoh.AudioNarrationUrl, req.Pharaoh.VideoDocumentaryUrl,
		req.Pharaoh.PopularityScore, req.Pharaoh.UserRating, req.Pharaoh.UnlockInGame, req.Pharaoh.Rarity,
		traitsJSON, req.Pharaoh.Source, req.Pharaoh.Verified, req.Pharaoh.Language,
	).Scan(&updatedAt)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update pharaoh: %v", err)
	}

	// Update timestamp in response
	req.Pharaoh.Id = req.Id
	req.Pharaoh.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &pb.UpdatePharaohResponse{
		StatusCode: 200,
		Message:    "Pharaoh updated successfully",
		Result: &pb.UpdatePharaohResponse_Pharaoh{
			Pharaoh: req.Pharaoh,
		},
	}, nil
}

// DeletePharaoh deletes a pharaoh record
func (s *PharaohServer) DeletePharaoh(ctx context.Context, req *pb.DeletePharaohRequest) (*pb.DeletePharaohResponse, error) {
	log.Printf("Deleting pharaoh: %s", req.Id)

	if req.Id == "" {
		return &pb.DeletePharaohResponse{
			StatusCode: 400,
			Message:    "Pharaoh ID is required",
			Result: &pb.DeletePharaohResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Pharaoh ID is required",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Check if pharaoh exists
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pharaohs WHERE id = $1)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to check pharaoh existence: %v", err)
	}

	if !exists {
		return &pb.DeletePharaohResponse{
			StatusCode: 404,
			Message:    "Pharaoh not found",
			Result: &pb.DeletePharaohResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      404,
					Message:   "Pharaoh not found",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Delete the pharaoh
	result, err := s.db.ExecContext(ctx, "DELETE FROM pharaohs WHERE id = $1", req.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to delete pharaoh: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to get affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return &pb.DeletePharaohResponse{
			StatusCode: 404,
			Message:    "Pharaoh not found",
			Result: &pb.DeletePharaohResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      404,
					Message:   "Pharaoh not found",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.DeletePharaohResponse{
		StatusCode: 200,
		Message:    "Pharaoh deleted successfully",
		Result: &pb.DeletePharaohResponse_Success{
			Success: true,
		},
	}, nil
}

// GetPharaohsByRarity retrieves pharaohs by rarity level
func (s *PharaohServer) GetPharaohsByRarity(ctx context.Context, req *pb.GetPharaohsByRarityRequest) (*pb.GetPharaohsByRarityResponse, error) {
	log.Printf("Getting pharaohs by rarity: %s", req.Rarity)

	if req.Rarity == "" {
		return &pb.GetPharaohsByRarityResponse{
			StatusCode: 400,
			Message:    "Rarity is required",
			Result: &pb.GetPharaohsByRarityResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Rarity is required",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Validate rarity values
	validRarities := map[string]bool{
		"Common":    true,
		"Rare":      true,
		"Epic":      true,
		"Legendary": true,
	}

	if !validRarities[req.Rarity] {
		return &pb.GetPharaohsByRarityResponse{
			StatusCode: 400,
			Message:    "Invalid rarity. Must be one of: Common, Rare, Epic, Legendary",
			Result: &pb.GetPharaohsByRarityResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid rarity. Must be one of: Common, Rare, Epic, Legendary",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	query := `
		SELECT id, name, birth_name, throne_name, epithet, dynasty, period,
		       reign_start, reign_end, length_of_reign_years, popularity_score, user_rating,
		       rarity, image_url, verified, created_at, updated_at
		FROM pharaohs 
		WHERE rarity = $1
		ORDER BY popularity_score DESC, name ASC
	`

	rows, err := s.db.QueryContext(ctx, query, req.Rarity)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to query pharaohs by rarity: %v", err)
	}
	defer rows.Close()

	var pharaohs []*pb.Pharaoh
	for rows.Next() {
		pharaoh := &pb.Pharaoh{}
		var birthName, throneName, epithet, period sql.NullString
		var dynasty, reignStart, reignEnd, lengthOfReignYears sql.NullInt32
		var popularityScore, userRating sql.NullFloat64
		var rarity, imageUrl sql.NullString
		var verified sql.NullBool
		var createdAt, updatedAt time.Time

		err := rows.Scan(
			&pharaoh.Id, &pharaoh.Name, &birthName, &throneName, &epithet,
			&dynasty, &period, &reignStart, &reignEnd, &lengthOfReignYears,
			&popularityScore, &userRating, &rarity, &imageUrl, &verified,
			&createdAt, &updatedAt,
		)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "Failed to scan pharaoh: %v", err)
		}

		// Handle nullable fields
		if birthName.Valid {
			pharaoh.BirthName = birthName.String
		}
		if throneName.Valid {
			pharaoh.ThroneName = throneName.String
		}
		if epithet.Valid {
			pharaoh.Epithet = epithet.String
		}
		if dynasty.Valid {
			pharaoh.Dynasty = dynasty.Int32
		}
		if period.Valid {
			pharaoh.Period = period.String
		}
		if reignStart.Valid {
			pharaoh.ReignStart = reignStart.Int32
		}
		if reignEnd.Valid {
			pharaoh.ReignEnd = reignEnd.Int32
		}
		if lengthOfReignYears.Valid {
			pharaoh.LengthOfReignYears = lengthOfReignYears.Int32
		}
		if popularityScore.Valid {
			pharaoh.PopularityScore = float32(popularityScore.Float64)
		}
		if userRating.Valid {
			pharaoh.UserRating = float32(userRating.Float64)
		}
		if rarity.Valid {
			pharaoh.Rarity = rarity.String
		}
		if imageUrl.Valid {
			pharaoh.ImageUrl = imageUrl.String
		}
		if verified.Valid {
			pharaoh.Verified = verified.Bool
		}

		pharaoh.CreatedAt = createdAt.Format(time.RFC3339)
		pharaoh.UpdatedAt = updatedAt.Format(time.RFC3339)

		pharaohs = append(pharaohs, pharaoh)
	}

	if err = rows.Err(); err != nil {
		return nil, status.Errorf(codes.Internal, "Error iterating pharaohs: %v", err)
	}

	return &pb.GetPharaohsByRarityResponse{
		StatusCode: 200,
		Message:    fmt.Sprintf("Found %d %s pharaohs", len(pharaohs), req.Rarity),
		Result: &pb.GetPharaohsByRarityResponse_Pharaohs{
			Pharaohs: &pb.PharaohList{
				Pharaohs:   pharaohs,
				TotalCount: int32(len(pharaohs)),
				Page:       1,
				Limit:      int32(len(pharaohs)),
			},
		},
	}, nil
}

// UpdatePharaohRating updates the user rating for a specific pharaoh
func (s *PharaohServer) UpdatePharaohRating(ctx context.Context, req *pb.UpdatePharaohRatingRequest) (*pb.UpdatePharaohRatingResponse, error) {
	log.Printf("Updating pharaoh rating for ID: %s, new rating: %f", req.Id, req.Rating)

	if req.Id == "" {
		return &pb.UpdatePharaohRatingResponse{
			StatusCode: 400,
			Message:    "Pharaoh ID is required",
			Result: &pb.UpdatePharaohRatingResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Pharaoh ID is required",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.Rating < 0 || req.Rating > 5 {
		return &pb.UpdatePharaohRatingResponse{
			StatusCode: 400,
			Message:    "Invalid rating. Must be between 0 and 5",
			Result: &pb.UpdatePharaohRatingResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid rating. Must be between 0 and 5",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Check if pharaoh exists
	var exists bool
	err := s.db.QueryRowContext(ctx, "SELECT EXISTS(SELECT 1 FROM pharaohs WHERE id = $1)", req.Id).Scan(&exists)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to check pharaoh existence: %v", err)
	}

	if !exists {
		return &pb.UpdatePharaohRatingResponse{
			StatusCode: 404,
			Message:    "Pharaoh not found",
			Result: &pb.UpdatePharaohRatingResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      404,
					Message:   "Pharaoh not found",
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Update the rating
	query := `
		UPDATE pharaohs 
		SET user_rating = $1, updated_at = CURRENT_TIMESTAMP 
		WHERE id = $2
		RETURNING id, name, birth_name, throne_name, epithet, dynasty, period,
		         reign_start, reign_end, length_of_reign_years, popularity_score, user_rating,
		         rarity, image_url, verified, created_at, updated_at
	`

	row := s.db.QueryRowContext(ctx, query, req.Rating, req.Id)

	pharaoh := &pb.Pharaoh{}
	var birthName, throneName, epithet, period sql.NullString
	var dynasty, reignStart, reignEnd, lengthOfReignYears sql.NullInt32
	var popularityScore, userRating sql.NullFloat64
	var rarity, imageUrl sql.NullString
	var verified sql.NullBool
	var createdAt, updatedAt time.Time

	err = row.Scan(
		&pharaoh.Id, &pharaoh.Name, &birthName, &throneName, &epithet,
		&dynasty, &period, &reignStart, &reignEnd, &lengthOfReignYears,
		&popularityScore, &userRating, &rarity, &imageUrl, &verified,
		&createdAt, &updatedAt,
	)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to update pharaoh rating: %v", err)
	}

	// Handle nullable fields
	if birthName.Valid {
		pharaoh.BirthName = birthName.String
	}
	if throneName.Valid {
		pharaoh.ThroneName = throneName.String
	}
	if epithet.Valid {
		pharaoh.Epithet = epithet.String
	}
	if dynasty.Valid {
		pharaoh.Dynasty = dynasty.Int32
	}
	if period.Valid {
		pharaoh.Period = period.String
	}
	if reignStart.Valid {
		pharaoh.ReignStart = reignStart.Int32
	}
	if reignEnd.Valid {
		pharaoh.ReignEnd = reignEnd.Int32
	}
	if lengthOfReignYears.Valid {
		pharaoh.LengthOfReignYears = lengthOfReignYears.Int32
	}
	if popularityScore.Valid {
		pharaoh.PopularityScore = float32(popularityScore.Float64)
	}
	if userRating.Valid {
		pharaoh.UserRating = float32(userRating.Float64)
	}
	if rarity.Valid {
		pharaoh.Rarity = rarity.String
	}
	if imageUrl.Valid {
		pharaoh.ImageUrl = imageUrl.String
	}
	if verified.Valid {
		pharaoh.Verified = verified.Bool
	}

	pharaoh.CreatedAt = createdAt.Format(time.RFC3339)
	pharaoh.UpdatedAt = updatedAt.Format(time.RFC3339)

	return &pb.UpdatePharaohRatingResponse{
		StatusCode: 200,
		Message:    "Pharaoh rating updated successfully",
		Result: &pb.UpdatePharaohRatingResponse_Pharaoh{
			Pharaoh: pharaoh,
		},
	}, nil
}
