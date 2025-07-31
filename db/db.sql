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
        cardholder_name VARCHAR(255) NOT NULL,
        card_brand VARCHAR(50) NOT NULL,
        card_type VARCHAR(50) NOT NULL,
        masked_card_number VARCHAR(20) NOT NULL,
        card_expiry_month INTEGER NOT NULL,
        card_expiry_year INTEGER NOT NULL,
        payment_gateway VARCHAR(100) NOT NULL,
        is_default BOOLEAN DEFAULT FALSE,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        -- Indexes for performance
        INDEX idx_user_id (user_id),
        INDEX idx_card_brand (card_brand),
        INDEX idx_card_type (card_type),
        UNIQUE KEY unique_user_card (
            user_id,
            masked_card_number,
            card_expiry_month,
            card_expiry_year
        )
    );