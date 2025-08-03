# Wallet Service

A gRPC-based wallet service that manages user wallets, balances, and transactions.

## Features

- Create and manage user wallets
- Fund wallets from cards
- Deduct from wallets for transactions
- Get wallet balance and transaction history
- Refund to wallets
- Secure balance encryption
- Token-based authentication

## Architecture

### Service Communication

This service communicates with other services via gRPC:

- **User Service**: For token validation and user authentication
- **Transaction Service**: For transaction processing (via gRPC calls, not proto imports)
- **Card Service**: For card validation during funding

### Why gRPC Communication Instead of Proto Imports?

We use gRPC client connections rather than importing proto files because:

1. **Service Decoupling**: Each service maintains its own proto definitions, reducing tight coupling
2. **Independent Deployment**: Services can be updated independently without affecting others
3. **Network Boundaries**: Services may run on different machines/containers
4. **Version Management**: Each service can evolve its API independently
5. **Security**: gRPC provides built-in authentication and encryption

## Environment Variables

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=5420
DB_USER=postgres
DB_PASSWORD=2521
DB_NAME=userdb

# User Service Configuration
USER_SERVICE_ADDR=localhost:50051

# Encryption Configuration
ENCRYPTION_KEY=wallet-service-encryption-key-2024
```

## Database Schema

### Wallets Table
```sql
CREATE TABLE wallets (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255) UNIQUE NOT NULL,
    encrypted_balance TEXT NOT NULL,
    currency VARCHAR(10) NOT NULL DEFAULT 'USD',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Wallet Transactions Table
```sql
CREATE TABLE wallet_transactions (
    id SERIAL PRIMARY KEY,
    wallet_id INTEGER REFERENCES wallets(id),
    type VARCHAR(20) NOT NULL, -- 'FUND', 'DEDUCT', 'REFUND'
    amount DECIMAL(15,2) NOT NULL,
    currency VARCHAR(10) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## API Endpoints

- `GetWalletByUserID`: Retrieve wallet information for a user
- `CreateWallet`: Create a new wallet for a user
- `FundWallet`: Add funds to a wallet from a card
- `DeductFromWallet`: Deduct funds from a wallet
- `GetWalletBalance`: Get current wallet balance
- `GetWalletTransactions`: Get wallet transaction history
- `RefundToWallet`: Refund money to a wallet

## Running the Service

```bash
# Install dependencies
go mod tidy

# Build the service
go build -o wallet-service .

# Run the service
./wallet-service
```

The service will start on port 50055.

## Security Features

- Balance encryption using AES-GCM
- Token-based authentication
- Input validation
- SQL injection prevention
- Database transactions for consistency