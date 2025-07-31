# User Service

A gRPC-based user management microservice built with Go that provides comprehensive user operations including authentication, profile management, and password reset functionality.

## Features

- **User Registration**: Create new user accounts with validation
- **User Authentication**: Login with email and password
- **Profile Management**: Update user information
- **Password Reset**: Forget password with verification code
- **User Deletion**: Remove user accounts
- **Input Validation**: Comprehensive validation for email, phone, and names
- **Password Security**: Bcrypt hashing for secure password storage
- **Database Integration**: PostgreSQL with proper schema

## API Methods

### gRPC Service Methods

1. **GetUserByEmail** - Retrieve user by email address
2. **CreateNewUser** - Register a new user account
3. **ForgetPassword** - Generate verification code for password reset
4. **VerifyCode** - Verify the provided verification code
5. **ResetPassword** - Reset user password after verification
6. **LoginUser** - Authenticate user and return access token
7. **UpdateUserData** - Update user profile information
8. **DeleteUserData** - Delete user account

## Validation Rules

### Full Name
- 2-100 characters
- Only letters, spaces, hyphens, and apostrophes
- Must contain at least one letter

### Email
- Valid email format
- 5-254 characters
- Standard email regex validation

### Password
- 8-128 characters
- At least one uppercase letter
- At least one lowercase letter
- At least one digit

### Phone Number
- International format with optional country code
- Optional field
- Supports + prefix

### Country Code
- 1-4 digits
- Optional + prefix
- Optional field

## Database Schema

The service uses PostgreSQL with the following user table structure:

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    full_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    country_code VARCHAR(10),
    phone_number VARCHAR(20),
    role VARCHAR(50) DEFAULT 'USER',
    verify_code VARCHAR(50),
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## Environment Variables

- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5432)
- `DB_USER`: Database username (default: postgres)
- `DB_PASSWORD`: Database password (default: 2521)
- `DB_NAME`: Database name (default: userdb)

## Running the Service

### Using Docker Compose (Recommended)

```bash
# From the project root directory
docker-compose up --build
```

### Running Locally

1. Install dependencies:
```bash
go mod tidy
```

2. Set up PostgreSQL database and run the schema from `../../db/db.sql`

3. Set environment variables or use defaults

4. Run the service:
```bash
go run main.go
```

The service will start on port 50051.

## Project Structure

```
user-services/
├── main.go              # Application entry point
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── Dockerfile           # Docker configuration
├── README.md            # This file
├── proto/               # Protocol buffer files
│   ├── user.proto       # Service definition
│   ├── user.pb.go       # Generated protobuf code
│   └── user_grpc.pb.go  # Generated gRPC code
└── server/              # Server implementation
    ├── database.go      # Database connection
    ├── user_server.go   # gRPC service implementation
    └── validation.go    # Input validation functions
```

## Security Features

- Password hashing using bcrypt
- Input validation and sanitization
- SQL injection prevention with parameterized queries
- UUID-based user IDs
- Secure error handling without information leakage

## Error Handling

The service returns structured error responses with:
- HTTP-style status codes
- Descriptive error messages
- Detailed error information
- Timestamps for debugging

## Dependencies

- `google.golang.org/grpc` - gRPC framework
- `google.golang.org/protobuf` - Protocol buffers
- `github.com/lib/pq` - PostgreSQL driver
- `github.com/google/uuid` - UUID generation
- `golang.org/x/crypto/bcrypt` - Password hashing