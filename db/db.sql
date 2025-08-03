-- users table
CREATE TABLE
    users (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
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

-- cards table
CREATE TABLE
    cards (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID NOT NULL,
        encrypted_cardholder_name TEXT NOT NULL,
        encrypted_card_number TEXT NOT NULL,
        encrypted_cvv TEXT NOT NULL,
        masked_card_number VARCHAR(20) NOT NULL,
        expiration_date VARCHAR(7) NOT NULL, -- MM/YYYY format
        card_type VARCHAR(50) NOT NULL,
        status VARCHAR(20) DEFAULT 'ACTIVE',
        balance DECIMAL(15,2) DEFAULT 0.00,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        -- Foreign key constraint
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

-- Indexes for cards table
CREATE INDEX idx_user_id ON cards (user_id);
CREATE INDEX idx_card_type ON cards (card_type);
CREATE INDEX idx_status ON cards (status);
CREATE INDEX idx_masked_card_number ON cards (masked_card_number);

-- Add trigger to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_cards_updated_at BEFORE UPDATE ON cards
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- transactions table
CREATE TABLE
    transactions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
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
        -- Foreign key constraints
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
        FOREIGN KEY (card_id) REFERENCES cards(id) ON DELETE CASCADE
    );

-- Indexes for transactions table
CREATE INDEX idx_user_id_trans ON transactions (user_id);
CREATE INDEX idx_card_id_trans ON transactions (card_id);
CREATE INDEX idx_status_trans ON transactions (status);
CREATE INDEX idx_merchant_id ON transactions (merchant_id);
CREATE INDEX idx_created_at ON transactions (created_at);

-- Add trigger to update updated_at timestamp for transactions
CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- wallets table
CREATE TABLE
    wallets (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        user_id UUID UNIQUE NOT NULL,
        encrypted_balance TEXT NOT NULL,
        currency VARCHAR(10) NOT NULL DEFAULT 'USD',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        -- Foreign key constraint
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
    );

-- Indexes for wallets table
CREATE INDEX idx_wallet_user_id ON wallets (user_id);

-- wallet_transactions table
CREATE TABLE
    wallet_transactions (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
        wallet_id UUID NOT NULL,
        type VARCHAR(20) NOT NULL, -- 'FUND', 'DEDUCT', 'REFUND'
        amount DECIMAL(15,2) NOT NULL,
        currency VARCHAR(10) NOT NULL,
        description TEXT,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        -- Foreign key constraint
        FOREIGN KEY (wallet_id) REFERENCES wallets(id) ON DELETE CASCADE
    );

-- Indexes for wallet_transactions table
CREATE INDEX idx_wallet_trans_wallet_id ON wallet_transactions (wallet_id);
CREATE INDEX idx_wallet_trans_type ON wallet_transactions (type);
CREATE INDEX idx_wallet_trans_created_at ON wallet_transactions (created_at);

-- Add triggers to update updated_at timestamp for wallets
CREATE TRIGGER update_wallets_updated_at BEFORE UPDATE ON wallets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();