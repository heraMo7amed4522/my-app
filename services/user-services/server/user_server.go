package server

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	pb "user-services/proto"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/metadata"
	"gopkg.in/gomail.v2"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
	db *Database
}

func NewUserServer() *UserServer {
	return &UserServer{
		db: NewDatabase(),
	}
}

// MARK: GetUserByEmail
func (s *UserServer) GetUserByEmail(ctx context.Context, req *pb.GetUserByEmailRequest) (*pb.GetUserByEmialResponse, error) {
	log.Printf("GetUserByEmail called with email: %s", req.Email)

	if !ValidateEmail(req.Email) {
		return &pb.GetUserByEmialResponse{
			StatusCode: 400,
			Message:    "Invalid email format",
			Result: &pb.GetUserByEmialResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid email format need to check you email",
					Details:   []string{"Email must be in valid format must be contain @ and ."},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	query := `SELECT id, full_name, email, country_code, phone_number, role, verify_code, created_at, updated_at FROM users WHERE email = $1`
	row := s.db.DB.QueryRow(query, req.Email)

	var user pb.User
	var countryCode, phoneNumber, verifyCode sql.NullString
	var createdAt, updatedAt time.Time

	err := row.Scan(&user.Id, &user.FullName, &user.Email, &countryCode, &phoneNumber, &user.Role, &verifyCode, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.GetUserByEmialResponse{
				StatusCode: 404,
				Message:    "User not found",
				Result: &pb.GetUserByEmialResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "User not found with this email",
						Details:   []string{"No user found with this email"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	if countryCode.Valid {
		user.CountryCode = countryCode.String
	}
	if phoneNumber.Valid {
		user.PhoneNumber = phoneNumber.String
	}
	if verifyCode.Valid {
		user.VerifiyCode = verifyCode.String
	}
	user.CreateAt = createdAt.Format(time.RFC3339)
	user.UpdateAt = updatedAt.Format(time.RFC3339)

	return &pb.GetUserByEmialResponse{
		StatusCode: 200,
		Message:    "User retrieved successfully",
		Result: &pb.GetUserByEmialResponse_User{
			User: &user,
		},
	}, nil
}

// MARK: CreateNewUser
func (s *UserServer) CreateNewUser(ctx context.Context, req *pb.CreateNewUserRequest) (*pb.CreateNewUserResponse, error) {
	log.Printf("CreateNewUser called with email: %s", req.Email)
	if !ValidateFullName(req.FullName) {
		return &pb.CreateNewUserResponse{
			StatusCode: 400,
			Message:    "Invalid full name",
			Result: &pb.CreateNewUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid full name",
					Details:   []string{"Full name must be 2-100 characters and contain only letters, spaces, hyphens, and apostrophes"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if !ValidateEmail(req.Email) {
		return &pb.CreateNewUserResponse{
			StatusCode: 400,
			Message:    "Invalid email format",
			Result: &pb.CreateNewUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid email format",
					Details:   []string{"Email must be in valid format"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if !ValidatePassword(req.Password) {
		return &pb.CreateNewUserResponse{
			StatusCode: 400,
			Message:    "Invalid password",
			Result: &pb.CreateNewUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid password",
					Details:   []string{"Password must be 8-128 characters with at least one uppercase, lowercase, and digit"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	if req.PhoneNumber != "" && !ValidatePhoneNumberWithCountry(req.PhoneNumber, req.CountryCode) {
		return &pb.CreateNewUserResponse{
			StatusCode: 400,
			Message:    "Invalid phone number",
			Result: &pb.CreateNewUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid phone number",
					Details:   []string{"Phone number must be in valid international format for the specified country"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.CountryCode != "" && !ValidateCountryCode(req.CountryCode) {
		return &pb.CreateNewUserResponse{
			StatusCode: 400,
			Message:    "Invalid country code",
			Result: &pb.CreateNewUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid country code",
					Details:   []string{"Country code must be 1-4 digits"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	var existingID string
	err := s.db.DB.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&existingID)
	if err == nil {
		return &pb.CreateNewUserResponse{
			StatusCode: 409,
			Message:    "User already exists",
			Result: &pb.CreateNewUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      409,
					Message:   "User already exists",
					Details:   []string{"A user with this email already exists"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	userID := uuid.New().String()
	now := time.Now()
	role := req.Role
	if role == "" {
		role = "USER"
	}
	query := `INSERT INTO users (id, full_name, email, country_code, phone_number, role, verify_code, password_hash, created_at, updated_at) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	_, err = s.db.DB.Exec(query, userID, req.FullName, req.Email,
		nullString(req.CountryCode), nullString(req.PhoneNumber),
		role, nullString(req.VerifiyCode), string(hashedPassword), now, now)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %v", err)
	}
	user := &pb.User{
		Id:          userID,
		FullName:    req.FullName,
		Email:       req.Email,
		CountryCode: req.CountryCode,
		PhoneNumber: req.PhoneNumber,
		Role:        role,
		VerifiyCode: req.VerifiyCode,
		CreateAt:    now.Format(time.RFC3339),
		UpdateAt:    now.Format(time.RFC3339),
	}

	return &pb.CreateNewUserResponse{
		StatusCode: 201,
		Message:    "User created successfully",
		Result: &pb.CreateNewUserResponse_User{
			User: user,
		},
	}, nil
}

func nullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: s != ""}
}

// MARK: ForgetPassword
func (s *UserServer) ForgetPassword(ctx context.Context, req *pb.ForgetPasswordRequest) (*pb.ForgetPasswordResponse, error) {
	log.Printf("ForgetPassword called with email: %s, type: %s", req.Email, req.Type)

	if !ValidateEmail(req.Email) {
		return &pb.ForgetPasswordResponse{
			StatusCode: 400,
			Message:    "Invalid email format",
			Result: &pb.ForgetPasswordResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid email format",
					Details:   []string{"Email must be in valid format"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	var userID, phoneNumber, countryCode string
	var phoneNumberNull, countryCodeNull sql.NullString
	query := `SELECT id, phone_number, country_code FROM users WHERE email = $1`
	err := s.db.DB.QueryRow(query, req.Email).Scan(&userID, &phoneNumberNull, &countryCodeNull)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.ForgetPasswordResponse{
				StatusCode: 404,
				Message:    "User not found",
				Result: &pb.ForgetPasswordResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "User not found",
						Details:   []string{"No user found with this email"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	if phoneNumberNull.Valid {
		phoneNumber = phoneNumberNull.String
	}
	if countryCodeNull.Valid {
		countryCode = countryCodeNull.String
	}

	verifyCode := generateVerificationCode()
	_, err = s.db.DB.Exec("UPDATE users SET verify_code = $1, updated_at = $2 WHERE email = $3",
		verifyCode, time.Now(), req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to update verification code: %v", err)
	}
	if req.Type == "1" {
		// Send SMS via Twilio
		if phoneNumber == "" {
			return &pb.ForgetPasswordResponse{
				StatusCode: 400,
				Message:    "Phone number not found",
				Result: &pb.ForgetPasswordResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      400,
						Message:   "Phone number not found",
						Details:   []string{"User does not have a phone number for SMS verification"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		err = sendSMSVerificationCode(phoneNumber, countryCode, verifyCode)
		if err != nil {
			log.Printf("Failed to send SMS: %v", err)
			return &pb.ForgetPasswordResponse{
				StatusCode: 500,
				Message:    "Failed to send SMS",
				Result: &pb.ForgetPasswordResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      500,
						Message:   "Failed to send SMS verification code",
						Details:   []string{"Unable to send SMS at this time"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		return &pb.ForgetPasswordResponse{
			StatusCode: 200,
			Message:    "Verification code sent via SMS successfully",
			Result: &pb.ForgetPasswordResponse_VerifyCode{
				VerifyCode: verifyCode,
			},
		}, nil
	} else {
		// Send Email
		err = sendEmailVerificationCode(req.Email, verifyCode)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			return &pb.ForgetPasswordResponse{
				StatusCode: 500,
				Message:    "Failed to send email",
				Result: &pb.ForgetPasswordResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      500,
						Message:   "Failed to send email verification code",
						Details:   []string{"Unable to send email at this time"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}

		return &pb.ForgetPasswordResponse{
			StatusCode: 200,
			Message:    "Verification code sent via email successfully",
			Result: &pb.ForgetPasswordResponse_VerifyCode{
				VerifyCode: verifyCode,
			},
		}, nil
	}
}

// Helper function to send SMS via Twilio
func sendSMSVerificationCode(phoneNumber, countryCode, verifyCode string) error {
	accountSid := os.Getenv("TWILIO_ACCOUNT_SID")
	authToken := os.Getenv("TWILIO_AUTH_TOKEN")
	twilioPhoneNumber := os.Getenv("TWILIO_PHONE_NUMBER")
	if accountSid == "" || authToken == "" || twilioPhoneNumber == "" {
		return fmt.Errorf("Twilio credentials not configured")
	}
	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSid,
		Password: authToken,
	})
	fullPhoneNumber := phoneNumber
	if countryCode != "" && !strings.HasPrefix(phoneNumber, "+") {
		// Remove leading zero if present
		if strings.HasPrefix(phoneNumber, "0") {
			phoneNumber = phoneNumber[1:]
		}
		fullPhoneNumber = countryCode + phoneNumber
	}
	messageBody := fmt.Sprintf("Your password reset verification code is: %s", verifyCode)
	params := &openapi.CreateMessageParams{}
	params.SetTo(fullPhoneNumber)
	params.SetFrom(twilioPhoneNumber)
	params.SetBody(messageBody)

	_, err := client.Api.CreateMessage(params)
	if err != nil {
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	log.Printf("SMS sent successfully to %s", fullPhoneNumber)
	return nil
}

// Helper function to send email verification code
func sendEmailVerificationCode(toEmail, verifyCode string) error {
	senderEmail := os.Getenv("SEND_EMAIL")
	senderPassword := os.Getenv("EMAIL_PASSWORD")
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")

	if senderEmail == "" || senderPassword == "" || smtpHost == "" || smtpPortStr == "" {
		return fmt.Errorf("email configuration not complete")
	}

	smtpPort, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", senderEmail)
	m.SetHeader("To", toEmail)
	m.SetHeader("Subject", "Password Reset Verification Code")

	body := fmt.Sprintf(`
		<html>
		<body>
			<div style="text-align: center;">
				<h1>Hera Company</h1>
			</div>
			<h2>Password Reset Request</h2>
			<p>You have requested to reset your password. Please use the following verification code:</p>
			<h3 style="color: #007bff; font-size: 24px; letter-spacing: 2px;">%s</h3>
			<p>This code will expire in 15 minutes.</p>
			<p>If you did not request this password reset, please ignore this email.</p>
		</body>
		</html>
	`, verifyCode)

	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, senderEmail, senderPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %v", err)
	}

	log.Printf("Email sent successfully to %s", toEmail)
	return nil
}

// MARK: VerifyCode
func (s *UserServer) VerifyCode(ctx context.Context, req *pb.VerifyCodeRequest) (*pb.VerifyCodeResponse, error) {
	log.Printf("VerifyCode called with email: %s", req.Email)

	if !ValidateEmail(req.Email) {
		return &pb.VerifyCodeResponse{
			StatusCode: 400,
			Message:    "Invalid email format",
			Result: &pb.VerifyCodeResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid email format",
					Details:   []string{"Email must be in valid format"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	var storedCode sql.NullString
	err := s.db.DB.QueryRow("SELECT verify_code FROM users WHERE email = $1", req.Email).Scan(&storedCode)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.VerifyCodeResponse{
				StatusCode: 404,
				Message:    "User not found",
				Result: &pb.VerifyCodeResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "User not found",
						Details:   []string{"No user found with this email"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	if !storedCode.Valid || storedCode.String != req.VerifyCode {
		return &pb.VerifyCodeResponse{
			StatusCode: 400,
			Message:    "Invalid verification code",
			Result: &pb.VerifyCodeResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid verification code",
					Details:   []string{"The provided verification code is incorrect"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.VerifyCodeResponse{
		StatusCode: 200,
		Message:    "Verification code is valid",
		Result: &pb.VerifyCodeResponse_Response{
			Response: "Verification successful",
		},
	}, nil
}

// MARK: ResetPassword
func (s *UserServer) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.ResetPasswordResponse, error) {
	log.Printf("ResetPassword called with email: %s", req.Email)
	if !ValidateEmail(req.Email) {
		return &pb.ResetPasswordResponse{
			StatusCode: 400,
			Message:    "Invalid email format",
			Result: &pb.ResetPasswordResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid email format",
					Details:   []string{"Email must be in valid format"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if !ValidatePassword(req.Password) {
		return &pb.ResetPasswordResponse{
			StatusCode: 400,
			Message:    "Invalid password",
			Result: &pb.ResetPasswordResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid password",
					Details:   []string{"Password must be 8-128 characters with at least one uppercase, lowercase, and digit"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	var userID string
	err := s.db.DB.QueryRow("SELECT id FROM users WHERE email = $1", req.Email).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.ResetPasswordResponse{
				StatusCode: 404,
				Message:    "User not found",
				Result: &pb.ResetPasswordResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "User not found",
						Details:   []string{"No user found with this email"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %v", err)
	}
	_, err = s.db.DB.Exec("UPDATE users SET password_hash = $1, verify_code = NULL, updated_at = $2 WHERE email = $3",
		string(hashedPassword), time.Now(), req.Email)
	if err != nil {
		return nil, fmt.Errorf("failed to reset password: %v", err)
	}

	return &pb.ResetPasswordResponse{
		StatusCode: 200,
		Message:    "Password reset successfully",
		Result: &pb.ResetPasswordResponse_Response{
			Response: "Password has been reset successfully",
		},
	}, nil
}

// MARK: LoginUser
func (s *UserServer) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	log.Printf("LoginUser called with email: %s", req.Email)
	if !ValidateEmail(req.Email) {
		return &pb.LoginUserResponse{
			StatusCode: 400,
			Message:    "Invalid email format",
			Result: &pb.LoginUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid email format",
					Details:   []string{"Email must be in valid format"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	var user pb.User
	var passwordHash string
	var countryCode, phoneNumber, verifyCode sql.NullString
	var createdAt, updatedAt time.Time

	query := `SELECT id, full_name, email, country_code, phone_number, role, verify_code, password_hash, created_at, updated_at FROM users WHERE email = $1`
	row := s.db.DB.QueryRow(query, req.Email)

	err := row.Scan(&user.Id, &user.FullName, &user.Email, &countryCode, &phoneNumber, &user.Role, &verifyCode, &passwordHash, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.LoginUserResponse{
				StatusCode: 401,
				Message:    "Invalid credentials",
				Result: &pb.LoginUserResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      401,
						Message:   "Invalid credentials",
						Details:   []string{"Email or password is incorrect"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(req.Password))
	if err != nil {
		return &pb.LoginUserResponse{
			StatusCode: 401,
			Message:    "Invalid credentials",
			Result: &pb.LoginUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Invalid credentials",
					Details:   []string{"Email or password is incorrect"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	if countryCode.Valid {
		user.CountryCode = countryCode.String
	}
	if phoneNumber.Valid {
		user.PhoneNumber = phoneNumber.String
	}
	if verifyCode.Valid {
		user.VerifiyCode = verifyCode.String
	}
	user.CreateAt = createdAt.Format(time.RFC3339)
	user.UpdateAt = updatedAt.Format(time.RFC3339)
	// Generate both access and refresh tokens
	accessToken, err := generateTokens(user.Id, user.Email, user.Role)
	if err != nil {
		return &pb.LoginUserResponse{
			StatusCode: 500,
			Message:    "Failed to generate tokens",
			Result: &pb.LoginUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      500,
					Message:   "Failed to generate tokens",
					Details:   []string{"Internal server error"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	return &pb.LoginUserResponse{
		StatusCode: 200,
		Message:    "Login successful",
		Tokens: &pb.AuthTokens{
			AccessToken: accessToken,
		},
		Result: &pb.LoginUserResponse_User{
			User: &user,
		},
	}, nil
}

// MARK: UpdateUserData
func (s *UserServer) UpdateUserData(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	log.Printf("UpdateUserData called with email: %s", req.Email)
	metadata, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.UpdateUserResponse{
			StatusCode: 401,
			Message:    "Missing authentication",
			Result: &pb.UpdateUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Missing authentication",
					Details:   []string{"Authorization header is required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	authorization := metadata.Get("authorization")
	if len(authorization) == 0 {
		return &pb.UpdateUserResponse{
			StatusCode: 401,
			Message:    "Missing authorization header",
			Result: &pb.UpdateUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Missing authorization header",
					Details:   []string{"Bearer token is required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	authHeader := authorization[0]
	if !strings.HasPrefix(authHeader, "Bearer") {
		return &pb.UpdateUserResponse{
			StatusCode: 401,
			Message:    "Invalid authorization format",
			Result: &pb.UpdateUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Invalid authorization format",
					Details:   []string{"Authorization must be in format: Bearer <token>"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	token := strings.TrimPrefix(authHeader, "Bearer ")
	_, err := ValidateAccessToken(token)
	if err != nil {
		return &pb.UpdateUserResponse{
			StatusCode: 401,
			Message:    "Invalid or expired token",
			Result: &pb.UpdateUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Invalid or expired token",
					Details:   []string{"Token validation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	if req.FullName != "" && !ValidateFullName(req.FullName) {
		return &pb.UpdateUserResponse{
			StatusCode: 400,
			Message:    "Invalid full name",
			Result: &pb.UpdateUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid full name",
					Details:   []string{"Full name must be 2-100 characters and contain only letters, spaces, hyphens, and apostrophes"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.Email != "" && !ValidateEmail(req.Email) {
		return &pb.UpdateUserResponse{
			StatusCode: 400,
			Message:    "Invalid email format",
			Result: &pb.UpdateUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid email format",
					Details:   []string{"Email must be in valid format"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	if req.PhoneNumber != "" && !ValidatePhoneNumberWithCountry(req.PhoneNumber, req.CountryCode) {
		return &pb.UpdateUserResponse{
			StatusCode: 400,
			Message:    "Invalid phone number",
			Result: &pb.UpdateUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid phone number",
					Details:   []string{"Phone number must be in valid international format for the specified country"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	if req.CountryCode != "" && !ValidateCountryCode(req.CountryCode) {
		return &pb.UpdateUserResponse{
			StatusCode: 400,
			Message:    "Invalid country code",
			Result: &pb.UpdateUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid country code",
					Details:   []string{"Country code must be 1-4 digits"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	var existingUser pb.User
	var countryCode, phoneNumber, verifyCode sql.NullString
	var createdAt, updatedAt time.Time
	query := `SELECT id, full_name, email, country_code, phone_number, role, verify_code, created_at, updated_at FROM users WHERE email = $1`
	row := s.db.DB.QueryRow(query, req.Email)
	err = row.Scan(&existingUser.Id, &existingUser.FullName, &existingUser.Email, &countryCode, &phoneNumber, &existingUser.Role, &verifyCode, &createdAt, &updatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.UpdateUserResponse{
				StatusCode: 404,
				Message:    "User not found",
				Result: &pb.UpdateUserResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "User not found",
						Details:   []string{"No user found with this email"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}
	updateFields := []string{}
	updateValues := []interface{}{}
	paramIndex := 1

	if req.FullName != "" {
		updateFields = append(updateFields, fmt.Sprintf("full_name = $%d", paramIndex))
		updateValues = append(updateValues, req.FullName)
		paramIndex++
	}

	if req.CountryCode != "" {
		updateFields = append(updateFields, fmt.Sprintf("country_code = $%d", paramIndex))
		updateValues = append(updateValues, req.CountryCode)
		paramIndex++
	}

	if req.PhoneNumber != "" {
		updateFields = append(updateFields, fmt.Sprintf("phone_number = $%d", paramIndex))
		updateValues = append(updateValues, req.PhoneNumber)
		paramIndex++
	}

	if req.Role != "" {
		updateFields = append(updateFields, fmt.Sprintf("role = $%d", paramIndex))
		updateValues = append(updateValues, req.Role)
		paramIndex++
	}

	if len(updateFields) == 0 {
		return &pb.UpdateUserResponse{
			StatusCode: 400,
			Message:    "No fields to update",
			Result: &pb.UpdateUserResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "No fields to update",
					Details:   []string{"At least one field must be provided for update"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}
	updateFields = append(updateFields, fmt.Sprintf("updated_at = $%d", paramIndex))
	updateValues = append(updateValues, time.Now())
	paramIndex++
	updateValues = append(updateValues, req.Email)
	updateQuery := fmt.Sprintf("UPDATE users SET %s WHERE email = $%d",
		strings.Join(updateFields, ", "), paramIndex)

	_, err = s.db.DB.Exec(updateQuery, updateValues...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %v", err)
	}
	updatedUser, err := s.GetUserByEmail(ctx, &pb.GetUserByEmailRequest{Email: req.Email})
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %v", err)
	}

	return &pb.UpdateUserResponse{
		StatusCode: 200,
		Message:    "User updated successfully",
		Result: &pb.UpdateUserResponse_User{
			User: updatedUser.GetUser(),
		},
	}, nil
}

// MARK: DeleteUserData
func (s *UserServer) DeleteUserData(ctx context.Context, req *pb.DeleteUserDatatRequest) (*pb.DeleteUserDataResponse, error) {
	log.Printf("DeleteUserData called with ID: %s", req.Id)

	// Validate access token from metadata
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return &pb.DeleteUserDataResponse{
			StatusCode: 401,
			Message:    "Missing authentication",
			Result: &pb.DeleteUserDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Missing authentication",
					Details:   []string{"Authorization header is required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	authorization := md.Get("authorization")
	if len(authorization) == 0 {
		return &pb.DeleteUserDataResponse{
			StatusCode: 401,
			Message:    "Missing authorization header",
			Result: &pb.DeleteUserDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Missing authorization header",
					Details:   []string{"Bearer token is required"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Extract token from "Bearer <token>" format
	authHeader := authorization[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return &pb.DeleteUserDataResponse{
			StatusCode: 401,
			Message:    "Invalid authorization format",
			Result: &pb.DeleteUserDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Invalid authorization format",
					Details:   []string{"Authorization must be in format: Bearer <token>"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	// Validate the token
	claims, err := ValidateAccessToken(token)
	if err != nil {
		return &pb.DeleteUserDataResponse{
			StatusCode: 401,
			Message:    "Invalid or expired token",
			Result: &pb.DeleteUserDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Invalid or expired token",
					Details:   []string{"Token validation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Extract user ID from token
	tokenUserID, ok := (*claims)["user_id"].(string)
	if !ok {
		return &pb.DeleteUserDataResponse{
			StatusCode: 401,
			Message:    "Invalid token claims",
			Result: &pb.DeleteUserDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Invalid token claims",
					Details:   []string{"Token does not contain valid user ID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Check if user is trying to delete their own account or has admin role
	userRole, _ := (*claims)["role"].(string)
	if tokenUserID != req.Id && userRole != "ADMIN" {
		return &pb.DeleteUserDataResponse{
			StatusCode: 403,
			Message:    "Forbidden",
			Result: &pb.DeleteUserDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      403,
					Message:   "Forbidden",
					Details:   []string{"You can only delete your own account or need admin privileges"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Validate user ID format
	if _, err := uuid.Parse(req.Id); err != nil {
		return &pb.DeleteUserDataResponse{
			StatusCode: 400,
			Message:    "Invalid user ID format",
			Result: &pb.DeleteUserDataResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      400,
					Message:   "Invalid user ID format",
					Details:   []string{"User ID must be a valid UUID"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	// Check if user exists
	var userID string
	err = s.db.DB.QueryRow("SELECT id FROM users WHERE id = $1", req.Id).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return &pb.DeleteUserDataResponse{
				StatusCode: 404,
				Message:    "User not found",
				Result: &pb.DeleteUserDataResponse_Error{
					Error: &pb.ErrorDetails{
						Code:      404,
						Message:   "User not found",
						Details:   []string{"No user found with this ID"},
						Timestamp: time.Now().Format(time.RFC3339),
					},
				},
			}, nil
		}
		return nil, fmt.Errorf("database error: %v", err)
	}

	// Delete the user
	_, err = s.db.DB.Exec("DELETE FROM users WHERE id = $1", req.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %v", err)
	}

	return &pb.DeleteUserDataResponse{
		StatusCode: 200,
		Message:    "User deleted successfully",
		Result: &pb.DeleteUserDataResponse_Response{
			Response: "User has been deleted successfully",
		},
	}, nil
}

// MARK : Generate Verification Code
func generateVerificationCode() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(900000) + 100000)
}

// Standardize all JWT secret handling
func getJWTSecret() string {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}
	return jwtSecret
}

func generateTokens(userID, email, role string) (string, error) {
	jwtSecret := getJWTSecret()
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable in production.")
	}
	// Generate access token (1 hour expiry)
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}
	return accessTokenString, nil
}

// Standardized token validation for system-wide use
func ValidateAccessToken(tokenString string) (*jwt.MapClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-super-secret-jwt-key-change-this-in-production"
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

func (s *UserServer) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	claims, err := ValidateAccessToken(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{
			StatusCode: 401,
			Message:    "Invalid token",
			Result: &pb.ValidateTokenResponse_Error{
				Error: &pb.ErrorDetails{
					Code:      401,
					Message:   "Invalid token",
					Details:   []string{"Token validation failed"},
					Timestamp: time.Now().Format(time.RFC3339),
				},
			},
		}, nil
	}

	userID, _ := (*claims)["user_id"].(string)
	email, _ := (*claims)["email"].(string)
	role, _ := (*claims)["role"].(string)
	exp, _ := (*claims)["exp"].(float64)

	return &pb.ValidateTokenResponse{
		StatusCode: 200,
		Message:    "Token is valid",
		Result: &pb.ValidateTokenResponse_Claims{
			Claims: &pb.TokenClaims{
				UserId: userID,
				Email:  email,
				Role:   role,
				Exp:    int64(exp),
			},
		},
	}, nil
}
