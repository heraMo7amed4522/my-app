# Transaction Service

A gRPC-based transaction service built with Go that manages financial transactions for the payment system.

## Features

- **Transaction Management**: Create, retrieve, and manage financial transactions
- **User Authentication**: Token-based authentication via user service
- **Card Validation**: Ensures transactions are only created for user-owned cards
- **Transaction Status Tracking**: Support for PENDING, COMPLETED, and REFUNDED statuses
- **Refund Processing**: Handle transaction refunds with proper validation
- **Database Integration**: PostgreSQL database with proper indexing and foreign keys

## API Endpoints

### gRPC Methods

1. **GetTransactionByID**: Retrieve a specific transaction by ID
2. **GetTransactionByUserID**: Get all transactions for a user
3. **GetTransactionByCardID**: Get all transactions for a specific card
4. **GetTransactionByStatus**: Filter transactions by status
5. **GetTransactionByDate**: Get transactions for a specific date
6. **AddTransaction**: Create a new transaction
7. **RefundTransaction**: Process a refund for an existing transaction

## Database Schema

The service uses a `transactions` table with the following structure:

```sql
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    card_id UUID NOT NULL,
    merchant_id VARCHAR(255) NOT NULL,
    merchant_name VARCHAR(255) NOT NULL,
    card_number VARCHAR(20) NOT NULL,
    merchant_category VARCHAR(100) NOT NULL,
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(3) NOT NULL DEFAULT 'USD',
    status VARCHAR(20) DEFAULT 'PENDING',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE
);
```

## Security

- All endpoints require valid JWT token authentication
- Users can only access their own transactions
- Card ownership validation before transaction creation
- Proper authorization checks for all operations

## Environment Variables

- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 5420)
- `DB_USER`: Database username (default: postgres)
- `DB_PASSWORD`: Database password (default: 2521)
- `DB_NAME`: Database name (default: userdb)
- `USER_SERVICE_ADDR`: User service address for token validation (default: localhost:50051)

## Running the Service

### Local Development
```bash
go run main.go
```

### Docker
```bash
docker build -t transaction-service .
docker run -p 50054:50054 transaction-service
```

### Docker Compose
The service is included in the main docker-compose.yml and will start automatically with:
```bash
docker-compose up
```

## Port

The service runs on port **50054**

## Dependencies

- PostgreSQL database
- User service (for token validation)
- Card service (for card validation)

## Transaction Statuses

- **PENDING**: Transaction created but not yet processed
- **COMPLETED**: Transaction successfully processed
- **REFUNDED**: Transaction has been refunded
- **FAILED**: Transaction processing failed