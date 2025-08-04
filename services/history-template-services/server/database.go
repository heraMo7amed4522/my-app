package server

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	pb "history-template-services/proto"

	_ "github.com/lib/pq"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Database struct {
	conn *sql.DB
}

func NewDatabase() *Database {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "password"
	}
	if dbName == "" {
		dbName = "myapp"
	}

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database")
	return &Database{conn: db}
}

func (db *Database) GetTemplateByID(id string) (*pb.HistoryTemplate, error) {
	query := `
		SELECT id, title, description, era, dynasty, pharaoh_id, difficulty, 
		       thumbnail_url, language, is_active, version, published_at, created_at, updated_at
		FROM history_templates 
		WHERE id = $1 AND is_active = true
	`

	var template pb.HistoryTemplate
	var publishedAt, createdAt, updatedAt sql.NullTime

	err := db.conn.QueryRow(query, id).Scan(
		&template.Id,
		&template.Title,
		&template.Description,
		&template.Era,
		&template.Dynasty,
		&template.PharaohId,
		&template.Difficulty,
		&template.ThumbnailUrl,
		&template.Language,
		&template.IsActive,
		&template.Version,
		&publishedAt,
		&createdAt,
		&updatedAt,
	)

	if err != nil {
		return nil, err
	}

	// Convert timestamps
	if publishedAt.Valid {
		template.PublishedAt = timestamppb.New(publishedAt.Time)
	}
	if createdAt.Valid {
		template.CreatedAt = timestamppb.New(createdAt.Time)
	}
	if updatedAt.Valid {
		template.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	// Get sections for this template
	sections, err := db.getTemplateSections(id)
	if err != nil {
		log.Printf("Warning: Failed to get sections for template %s: %v", id, err)
	} else {
		template.Sections = sections
	}

	return &template, nil
}

func (db *Database) GetAllTemplates(page, limit int32, sortBy, order string) ([]*pb.HistoryTemplate, int32, error) {
	// Count total templates
	countQuery := "SELECT COUNT(*) FROM history_templates WHERE is_active = true"
	var totalCount int32
	err := db.conn.QueryRow(countQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Build query with pagination
	offset := (page - 1) * limit
	query := `
		SELECT id, title, description, era, dynasty, pharaoh_id, difficulty, 
		       thumbnail_url, language, is_active, version, published_at, created_at, updated_at
		FROM history_templates 
		WHERE is_active = true
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := db.conn.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*pb.HistoryTemplate
	for rows.Next() {
		var template pb.HistoryTemplate
		var publishedAt, createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&template.Id,
			&template.Title,
			&template.Description,
			&template.Era,
			&template.Dynasty,
			&template.PharaohId,
			&template.Difficulty,
			&template.ThumbnailUrl,
			&template.Language,
			&template.IsActive,
			&template.Version,
			&publishedAt,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert timestamps
		if publishedAt.Valid {
			template.PublishedAt = timestamppb.New(publishedAt.Time)
		}
		if createdAt.Valid {
			template.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			template.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		templates = append(templates, &template)
	}

	return templates, totalCount, nil
}

func (db *Database) CreateTemplate(template *pb.HistoryTemplate) (*pb.HistoryTemplate, error) {
	query := `
		INSERT INTO history_templates (id, title, description, era, dynasty, pharaoh_id, 
		                              difficulty, thumbnail_url, language, is_active, version, 
		                              published_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
		RETURNING id
	`

	var id string
	err := db.conn.QueryRow(query,
		template.Id,
		template.Title,
		template.Description,
		template.Era,
		template.Dynasty,
		template.PharaohId,
		template.Difficulty,
		template.ThumbnailUrl,
		template.Language,
		template.IsActive,
		template.Version,
		template.PublishedAt.AsTime(),
		template.CreatedAt.AsTime(),
		template.UpdatedAt.AsTime(),
	).Scan(&id)

	if err != nil {
		return nil, err
	}

	return db.GetTemplateByID(id)
}

func (db *Database) UpdateTemplate(id string, template *pb.HistoryTemplate) (*pb.HistoryTemplate, error) {
	query := `
		UPDATE history_templates 
		SET title = $2, description = $3, era = $4, dynasty = $5, pharaoh_id = $6, 
		    difficulty = $7, thumbnail_url = $8, language = $9, is_active = $10, 
		    version = $11, updated_at = $12
		WHERE id = $1 AND is_active = true
		RETURNING id
	`

	var returnedID string
	err := db.conn.QueryRow(query,
		id,
		template.Title,
		template.Description,
		template.Era,
		template.Dynasty,
		template.PharaohId,
		template.Difficulty,
		template.ThumbnailUrl,
		template.Language,
		template.IsActive,
		template.Version,
		template.UpdatedAt.AsTime(),
	).Scan(&returnedID)

	if err != nil {
		return nil, err
	}

	return db.GetTemplateByID(returnedID)
}

func (db *Database) DeleteTemplate(id string) error {
	query := `UPDATE history_templates SET is_active = false WHERE id = $1 AND is_active = true`
	result, err := db.conn.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (db *Database) GetTemplatesByEra(era string, page, limit int32) ([]*pb.HistoryTemplate, int32, error) {
	// Count total templates
	countQuery := "SELECT COUNT(*) FROM history_templates WHERE era = $1 AND is_active = true"
	var totalCount int32
	err := db.conn.QueryRow(countQuery, era).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get templates with pagination
	offset := (page - 1) * limit
	query := `
		SELECT id, title, description, era, dynasty, pharaoh_id, difficulty, 
		       thumbnail_url, language, is_active, version, published_at, created_at, updated_at
		FROM history_templates 
		WHERE era = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.conn.Query(query, era, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*pb.HistoryTemplate
	for rows.Next() {
		var template pb.HistoryTemplate
		var publishedAt, createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&template.Id,
			&template.Title,
			&template.Description,
			&template.Era,
			&template.Dynasty,
			&template.PharaohId,
			&template.Difficulty,
			&template.ThumbnailUrl,
			&template.Language,
			&template.IsActive,
			&template.Version,
			&publishedAt,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert timestamps
		if publishedAt.Valid {
			template.PublishedAt = timestamppb.New(publishedAt.Time)
		}
		if createdAt.Valid {
			template.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			template.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		templates = append(templates, &template)
	}

	return templates, totalCount, nil
}

func (db *Database) GetTemplatesByDynasty(dynasty int32, page, limit int32) ([]*pb.HistoryTemplate, int32, error) {
	// Count total templates
	countQuery := "SELECT COUNT(*) FROM history_templates WHERE dynasty = $1 AND is_active = true"
	var totalCount int32
	err := db.conn.QueryRow(countQuery, dynasty).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get templates with pagination
	offset := (page - 1) * limit
	query := `
		SELECT id, title, description, era, dynasty, pharaoh_id, difficulty, 
		       thumbnail_url, language, is_active, version, published_at, created_at, updated_at
		FROM history_templates 
		WHERE dynasty = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.conn.Query(query, dynasty, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*pb.HistoryTemplate
	for rows.Next() {
		var template pb.HistoryTemplate
		var publishedAt, createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&template.Id,
			&template.Title,
			&template.Description,
			&template.Era,
			&template.Dynasty,
			&template.PharaohId,
			&template.Difficulty,
			&template.ThumbnailUrl,
			&template.Language,
			&template.IsActive,
			&template.Version,
			&publishedAt,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert timestamps
		if publishedAt.Valid {
			template.PublishedAt = timestamppb.New(publishedAt.Time)
		}
		if createdAt.Valid {
			template.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			template.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		templates = append(templates, &template)
	}

	return templates, totalCount, nil
}

func (db *Database) GetTemplatesByPharaoh(pharaohId string, page, limit int32) ([]*pb.HistoryTemplate, int32, error) {
	// Count total templates
	countQuery := "SELECT COUNT(*) FROM history_templates WHERE pharaoh_id = $1 AND is_active = true"
	var totalCount int32
	err := db.conn.QueryRow(countQuery, pharaohId).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get templates with pagination
	offset := (page - 1) * limit
	query := `
		SELECT id, title, description, era, dynasty, pharaoh_id, difficulty, 
		       thumbnail_url, language, is_active, version, published_at, created_at, updated_at
		FROM history_templates 
		WHERE pharaoh_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.conn.Query(query, pharaohId, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*pb.HistoryTemplate
	for rows.Next() {
		var template pb.HistoryTemplate
		var publishedAt, createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&template.Id,
			&template.Title,
			&template.Description,
			&template.Era,
			&template.Dynasty,
			&template.PharaohId,
			&template.Difficulty,
			&template.ThumbnailUrl,
			&template.Language,
			&template.IsActive,
			&template.Version,
			&publishedAt,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert timestamps
		if publishedAt.Valid {
			template.PublishedAt = timestamppb.New(publishedAt.Time)
		}
		if createdAt.Valid {
			template.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			template.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		templates = append(templates, &template)
	}

	return templates, totalCount, nil
}

func (db *Database) GetTemplatesByDifficulty(difficulty string, page, limit int32) ([]*pb.HistoryTemplate, int32, error) {
	// Count total templates
	countQuery := "SELECT COUNT(*) FROM history_templates WHERE difficulty = $1 AND is_active = true"
	var totalCount int32
	err := db.conn.QueryRow(countQuery, difficulty).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get templates with pagination
	offset := (page - 1) * limit
	query := `
		SELECT id, title, description, era, dynasty, pharaoh_id, difficulty, 
		       thumbnail_url, language, is_active, version, published_at, created_at, updated_at
		FROM history_templates 
		WHERE difficulty = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.conn.Query(query, difficulty, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*pb.HistoryTemplate
	for rows.Next() {
		var template pb.HistoryTemplate
		var publishedAt, createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&template.Id,
			&template.Title,
			&template.Description,
			&template.Era,
			&template.Dynasty,
			&template.PharaohId,
			&template.Difficulty,
			&template.ThumbnailUrl,
			&template.Language,
			&template.IsActive,
			&template.Version,
			&publishedAt,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert timestamps
		if publishedAt.Valid {
			template.PublishedAt = timestamppb.New(publishedAt.Time)
		}
		if createdAt.Valid {
			template.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			template.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		templates = append(templates, &template)
	}

	return templates, totalCount, nil
}

func (db *Database) SearchTemplates(query string, fields []string, page, limit int32) ([]*pb.HistoryTemplate, int32, error) {
	// Build search query
	searchQuery := "%" + query + "%"

	// Count total templates
	countQuery := `
		SELECT COUNT(*) FROM history_templates 
		WHERE (title ILIKE $1 OR description ILIKE $1 OR era ILIKE $1) 
		AND is_active = true
	`
	var totalCount int32
	err := db.conn.QueryRow(countQuery, searchQuery).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get templates with pagination
	offset := (page - 1) * limit
	sqlQuery := `
		SELECT id, title, description, era, dynasty, pharaoh_id, difficulty, 
		       thumbnail_url, language, is_active, version, published_at, created_at, updated_at
		FROM history_templates 
		WHERE (title ILIKE $1 OR description ILIKE $1 OR era ILIKE $1) 
		AND is_active = true
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.conn.Query(sqlQuery, searchQuery, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*pb.HistoryTemplate
	for rows.Next() {
		var template pb.HistoryTemplate
		var publishedAt, createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&template.Id,
			&template.Title,
			&template.Description,
			&template.Era,
			&template.Dynasty,
			&template.PharaohId,
			&template.Difficulty,
			&template.ThumbnailUrl,
			&template.Language,
			&template.IsActive,
			&template.Version,
			&publishedAt,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert timestamps
		if publishedAt.Valid {
			template.PublishedAt = timestamppb.New(publishedAt.Time)
		}
		if createdAt.Valid {
			template.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			template.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		templates = append(templates, &template)
	}

	return templates, totalCount, nil
}

func (db *Database) GetTemplatesByTag(tag string, page, limit int32) ([]*pb.HistoryTemplate, int32, error) {
	// Count total templates with the tag
	countQuery := `
		SELECT COUNT(*) FROM history_templates ht
		JOIN template_tags tt ON ht.id = tt.template_id
		WHERE tt.tag = $1 AND ht.is_active = true
	`
	var totalCount int32
	err := db.conn.QueryRow(countQuery, tag).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get templates with pagination
	offset := (page - 1) * limit
	query := `
		SELECT DISTINCT ht.id, ht.title, ht.description, ht.era, ht.dynasty, ht.pharaoh_id, 
		       ht.difficulty, ht.thumbnail_url, ht.language, ht.is_active, ht.version, 
		       ht.published_at, ht.created_at, ht.updated_at
		FROM history_templates ht
		JOIN template_tags tt ON ht.id = tt.template_id
		WHERE tt.tag = $1 AND ht.is_active = true
		ORDER BY ht.created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := db.conn.Query(query, tag, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*pb.HistoryTemplate
	for rows.Next() {
		var template pb.HistoryTemplate
		var publishedAt, createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&template.Id,
			&template.Title,
			&template.Description,
			&template.Era,
			&template.Dynasty,
			&template.PharaohId,
			&template.Difficulty,
			&template.ThumbnailUrl,
			&template.Language,
			&template.IsActive,
			&template.Version,
			&publishedAt,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert timestamps
		if publishedAt.Valid {
			template.PublishedAt = timestamppb.New(publishedAt.Time)
		}
		if createdAt.Valid {
			template.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			template.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		templates = append(templates, &template)
	}

	return templates, totalCount, nil
}

func (db *Database) GetRelatedTemplates(templateId string, limit int32) ([]*pb.HistoryTemplate, int32, error) {
	// Get related templates based on same era or pharaoh
	query := `
		SELECT id, title, description, era, dynasty, pharaoh_id, difficulty, 
		       thumbnail_url, language, is_active, version, published_at, created_at, updated_at
		FROM history_templates 
		WHERE id != $1 AND is_active = true
		AND (era = (SELECT era FROM history_templates WHERE id = $1) 
		     OR pharaoh_id = (SELECT pharaoh_id FROM history_templates WHERE id = $1))
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := db.conn.Query(query, templateId, limit)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var templates []*pb.HistoryTemplate
	for rows.Next() {
		var template pb.HistoryTemplate
		var publishedAt, createdAt, updatedAt sql.NullTime

		err := rows.Scan(
			&template.Id,
			&template.Title,
			&template.Description,
			&template.Era,
			&template.Dynasty,
			&template.PharaohId,
			&template.Difficulty,
			&template.ThumbnailUrl,
			&template.Language,
			&template.IsActive,
			&template.Version,
			&publishedAt,
			&createdAt,
			&updatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		// Convert timestamps
		if publishedAt.Valid {
			template.PublishedAt = timestamppb.New(publishedAt.Time)
		}
		if createdAt.Valid {
			template.CreatedAt = timestamppb.New(createdAt.Time)
		}
		if updatedAt.Valid {
			template.UpdatedAt = timestamppb.New(updatedAt.Time)
		}

		templates = append(templates, &template)
	}

	return templates, int32(len(templates)), nil
}

func (db *Database) getTemplateSections(templateID string) ([]*pb.TemplateSection, error) {
	query := `
		SELECT id, template_id, title, subtitle, content_type, content, 
		       order_index, optional, created_at
		FROM template_sections 
		WHERE template_id = $1 
		ORDER BY order_index ASC
	`

	rows, err := db.conn.Query(query, templateID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sections []*pb.TemplateSection
	for rows.Next() {
		var section pb.TemplateSection
		var createdAt sql.NullTime

		err := rows.Scan(
			&section.Id,
			&section.TemplateId,
			&section.Title,
			&section.Subtitle,
			&section.ContentType,
			&section.Content,
			&section.OrderIndex,
			&section.Optional,
			&createdAt,
		)
		if err != nil {
			return nil, err
		}

		if createdAt.Valid {
			section.CreatedAt = timestamppb.New(createdAt.Time)
		}

		sections = append(sections, &section)
	}

	return sections, nil
}
