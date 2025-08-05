package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"strings"
	"time"

	pb "upload-services/proto"

	"github.com/google/uuid"
)

type UploadServer struct {
	pb.UnimplementedUploadServiceServer
	db       *Database
	s3Client *S3Client
}

func NewUploadServer() *UploadServer {
	s3Client, err := NewS3Client()
	if err != nil {
		log.Fatalf("Failed to initialize S3 client: %v", err)
	}

	return &UploadServer{
		db:       NewDatabase(),
		s3Client: s3Client,
	}
}

func (s *UploadServer) UploadFile(ctx context.Context, req *pb.UploadFileRequest) (*pb.UploadFileResponse, error) {
	log.Printf("UploadFile called for user: %s, file: %s", req.UserId, req.FileName)

	// Validate input
	if req.UserId == "" {
		return &pb.UploadFileResponse{
			StatusCode: 400,
			Message:    "User ID is required",
			Result: &pb.UploadFileResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "User ID is required",
					Details:   []string{"User ID cannot be empty"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if len(req.FileData) == 0 {
		return &pb.UploadFileResponse{
			StatusCode: 400,
			Message:    "File data is required",
			Result: &pb.UploadFileResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "File data is required",
					Details:   []string{"File data cannot be empty"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Validate file type
	if !isValidFileType(req.ContentType, req.FileType) {
		return &pb.UploadFileResponse{
			StatusCode: 400,
			Message:    "Invalid file type",
			Result: &pb.UploadFileResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid file type",
					Details:   []string{"Only images and videos are allowed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Generate unique file key
	fileID := uuid.New().String()
	fileExt := filepath.Ext(req.FileName)
	fileKey := fmt.Sprintf("%s/%s/%s%s", req.UserId, getFileTypeString(req.FileType), fileID, fileExt)

	// Upload to S3
	fileURL, err := s.s3Client.UploadFile(req.FileData, fileKey, req.ContentType)
	if err != nil {
		log.Printf("Failed to upload file to S3: %v", err)
		return &pb.UploadFileResponse{
			StatusCode: 500,
			Message:    "Failed to upload file",
			Result: &pb.UploadFileResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{"Failed to upload file to storage"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Save to database
	query := `
		INSERT INTO uploads (id, user_id, file_key, file_name, file_type, content_type, file_size, file_url, uploaded_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	uploadedAt := time.Now()
	_, err = s.db.db.Exec(query, fileID, req.UserId, fileKey, req.FileName, 
		getFileTypeString(req.FileType), req.ContentType, len(req.FileData), fileURL, uploadedAt)

	if err != nil {
		log.Printf("Failed to save upload record to database: %v", err)
		// Try to delete from S3 since DB save failed
		s.s3Client.DeleteFile(fileKey)
		return &pb.UploadFileResponse{
			StatusCode: 500,
			Message:    "Failed to save upload record",
			Result: &pb.UploadFileResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{"Failed to save upload record"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.UploadFileResponse{
		StatusCode: 200,
		Message:    "File uploaded successfully",
		Result: &pb.UploadFileResponse_Success{
			Success: &pb.UploadSuccess{
				FileUrl:     fileURL,
				FileKey:     fileKey,
				FileId:      fileID,
				FileSize:    int64(len(req.FileData)),
				UploadedAt:  uploadedAt.Format(time.RFC3339),
			},
		},
	}, nil
}

func (s *UploadServer) DeleteFile(ctx context.Context, req *pb.DeleteFileRequest) (*pb.DeleteFileResponse, error) {
	log.Printf("DeleteFile called for user: %s, file_key: %s", req.UserId, req.FileKey)

	// Validate input
	if req.UserId == "" || req.FileKey == "" {
		return &pb.DeleteFileResponse{
			StatusCode: 400,
			Message:    "User ID and file key are required",
			Result: &pb.DeleteFileResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "User ID and file key are required",
					Details:   []string{"Both user_id and file_key must be provided"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Check if file exists and belongs to user
	var fileExists bool
	var deletedAt sql.NullTime
	query := `SELECT EXISTS(SELECT 1 FROM uploads WHERE file_key = $1 AND user_id = $2), deleted_at FROM uploads WHERE file_key = $1 AND user_id = $2`
	err := s.db.db.QueryRow(query, req.FileKey, req.UserId).Scan(&fileExists, &deletedAt)

	if err != nil && err != sql.ErrNoRows {
		log.Printf("Database error: %v", err)
		return &pb.DeleteFileResponse{
			StatusCode: 500,
			Message:    "Internal server error",
			Result: &pb.DeleteFileResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Database error",
					Details:   []string{"Failed to check file existence"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if !fileExists {
		return &pb.DeleteFileResponse{
			StatusCode: 404,
			Message:    "File not found",
			Result: &pb.DeleteFileResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      404,
					Message:   "File not found",
					Details:   []string{"File does not exist or does not belong to user"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if deletedAt.Valid {
		return &pb.DeleteFileResponse{
			StatusCode: 400,
			Message:    "File already deleted",
			Result: &pb.DeleteFileResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "File already deleted",
					Details:   []string{"This file has already been deleted"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Delete from S3
	err = s.s3Client.DeleteFile(req.FileKey)
	if err != nil {
		log.Printf("Failed to delete file from S3: %v", err)
		return &pb.DeleteFileResponse{
			StatusCode: 500,
			Message:    "Failed to delete file from storage",
			Result: &pb.DeleteFileResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{"Failed to delete file from storage"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Mark as deleted in database
	deletedTime := time.Now()
	updateQuery := `UPDATE uploads SET deleted_at = $1 WHERE file_key = $2 AND user_id = $3`
	_, err = s.db.db.Exec(updateQuery, deletedTime, req.FileKey, req.UserId)

	if err != nil {
		log.Printf("Failed to mark file as deleted in database: %v", err)
		return &pb.DeleteFileResponse{
			StatusCode: 500,
			Message:    "Failed to update delete record",
			Result: &pb.DeleteFileResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{"Failed to update delete record"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.DeleteFileResponse{
		StatusCode: 200,
		Message:    "File deleted successfully",
		Result: &pb.DeleteFileResponse_Success{
			Success: &pb.DeleteSuccess{
				FileKey:   req.FileKey,
				DeletedAt: deletedTime.Format(time.RFC3339),
			},
		},
	}, nil
}

func (s *UploadServer) GetFileURL(ctx context.Context, req *pb.GetFileURLRequest) (*pb.GetFileURLResponse, error) {
	log.Printf("GetFileURL called for file_key: %s", req.FileKey)

	if req.FileKey == "" {
		return &pb.GetFileURLResponse{
			StatusCode: 400,
			Message:    "File key is required",
			Result: &pb.GetFileURLResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "File key is required",
					Details:   []string{"File key cannot be empty"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	expiresIn := time.Duration(req.ExpiresIn) * time.Second
	if expiresIn <= 0 {
		expiresIn = 1 * time.Hour // Default to 1 hour
	}

	// Generate presigned URL
	signedURL, err := s.s3Client.GeneratePresignedURL(req.FileKey, expiresIn)
	if err != nil {
		log.Printf("Failed to generate presigned URL: %v", err)
		return &pb.GetFileURLResponse{
			StatusCode: 500,
			Message:    "Failed to generate file URL",
			Result: &pb.GetFileURLResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Internal server error",
					Details:   []string{"Failed to generate signed URL"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	expiresAt := time.Now().Add(expiresIn)

	return &pb.GetFileURLResponse{
		StatusCode: 200,
		Message:    "File URL generated successfully",
		Result: &pb.GetFileURLResponse_Success{
			Success: &pb.URLSuccess{
				SignedUrl: signedURL,
				ExpiresAt: expiresAt.Format(time.RFC3339),
			},
		},
	}, nil
}

// Helper functions
func isValidFileType(contentType string, fileType pb.FileType) bool {
	switch fileType {
	case pb.FileType_IMAGE:
		return strings.HasPrefix(contentType, "image/")
	case pb.FileType_VIDEO:
		return strings.HasPrefix(contentType, "video/")
	default:
		return false
	}
}

func getFileTypeString(fileType pb.FileType) string {
	switch fileType {
	case pb.FileType_IMAGE:
		return "images"
	case pb.FileType_VIDEO:
		return "videos"
	default:
		return "unknown"
	}
}