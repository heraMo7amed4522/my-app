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

-- Add trigger to update updated_at timestamp for wallets
CREATE TRIGGER update_wallets_updated_at BEFORE UPDATE ON wallets
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- uploads table for S3 file management
CREATE TABLE uploads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    file_key VARCHAR(500) NOT NULL UNIQUE,
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL, -- 'images' or 'videos'
    content_type VARCHAR(100) NOT NULL, -- MIME type
    file_size BIGINT NOT NULL,
    file_url TEXT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Foreign key constraint
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for uploads table
CREATE INDEX idx_uploads_user_id ON uploads (user_id);
CREATE INDEX idx_uploads_file_type ON uploads (file_type);
CREATE INDEX idx_uploads_uploaded_at ON uploads (uploaded_at);
CREATE INDEX idx_uploads_deleted_at ON uploads (deleted_at);
CREATE INDEX idx_uploads_file_key ON uploads (file_key);

-- Add trigger to update updated_at timestamp for uploads
CREATE TRIGGER update_uploads_updated_at BEFORE UPDATE ON uploads
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- stories table
CREATE TABLE stories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    author VARCHAR(100),
    status VARCHAR(20) DEFAULT 'DRAFT', -- DRAFT, PUBLISHED, ARCHIVED
    category VARCHAR(100),
    tags JSONB, -- array of tags
    image_url TEXT,
    language VARCHAR(10) DEFAULT 'en',
    is_featured BOOLEAN DEFAULT false,
    view_count INTEGER DEFAULT 0,
    like_count INTEGER DEFAULT 0,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    -- Foreign key constraint
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for stories table
CREATE INDEX idx_stories_user_id ON stories (user_id);
CREATE INDEX idx_stories_status ON stories (status);
CREATE INDEX idx_stories_category ON stories (category);
CREATE INDEX idx_stories_published_at ON stories (published_at);
CREATE INDEX idx_stories_is_featured ON stories (is_featured);
CREATE INDEX idx_stories_language ON stories (language);

-- Add trigger to update updated_at timestamp for stories
CREATE TRIGGER update_stories_updated_at BEFORE UPDATE ON stories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Optional: Story likes/reactions table
CREATE TABLE story_reactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    story_id UUID NOT NULL,
    user_id UUID NOT NULL,
    reaction_type VARCHAR(20) NOT NULL, -- LIKE, LOVE, DISLIKE
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (story_id) REFERENCES stories(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE(story_id, user_id) -- One reaction per user per story
);

CREATE INDEX idx_story_reactions_story_id ON story_reactions (story_id);
CREATE INDEX idx_story_reactions_user_id ON story_reactions (user_id);

-- pharaohs table
CREATE TABLE
    pharaohs (
        id VARCHAR(50) PRIMARY KEY, -- ph_001, ph_002, etc.
        name VARCHAR(255) NOT NULL,
        birth_name VARCHAR(255),
        throne_name VARCHAR(255),
        epithet TEXT,
        
        -- Historical & Chronological
        dynasty INTEGER,
        period VARCHAR(100),
        reign_start INTEGER, -- BCE year (negative)
        reign_end INTEGER,
        length_of_reign_years INTEGER,
        predecessor_id VARCHAR(50),
        successor_id VARCHAR(50),
        
        -- Family & Lineage
        father VARCHAR(255),
        mother VARCHAR(255),
        consorts JSONB, -- array of strings
        children_count INTEGER,
        notable_children JSONB, -- array of strings
        
        -- Rule & Governance
        capital VARCHAR(255),
        major_achievements JSONB, -- array of strings
        military_campaigns JSONB, -- array of strings
        building_projects JSONB, -- array of strings
        political_style VARCHAR(100),
        
        -- Religious & Divine Role
        divine_association JSONB, -- array of strings
        temple_affiliations JSONB, -- array of strings
        religious_reforms TEXT,
        pharaoh_as_god BOOLEAN DEFAULT true,
        
        -- Burial & Afterlife
        burial_site VARCHAR(255),
        tomb_discovered BOOLEAN DEFAULT false,
        discovery_year INTEGER,
        tomb_guardian VARCHAR(100),
        funerary_text TEXT,
        
        -- Artifacts & Treasures
        famous_artifacts JSONB, -- array of objects
        treasure_status VARCHAR(50),
        
        -- Media & Presentation
        image_url TEXT,
        statue_count INTEGER,
        mummy_location VARCHAR(255),
        audio_narration_url TEXT,
        video_documentary_url TEXT,
        
        -- User Engagement
        popularity_score DECIMAL(3,1) DEFAULT 0.0,
        user_rating DECIMAL(2,1) DEFAULT 0.0,
        unlock_in_game BOOLEAN DEFAULT false,
        rarity VARCHAR(50),
        traits JSONB, -- object with leadership, military, diplomacy scores
        
        -- Metadata
        source TEXT,
        verified BOOLEAN DEFAULT false,
        language VARCHAR(10) DEFAULT 'en',
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        
        -- Foreign key constraints
        FOREIGN KEY (predecessor_id) REFERENCES pharaohs(id),
        FOREIGN KEY (successor_id) REFERENCES pharaohs(id)
    );

-- Indexes for pharaohs table
CREATE INDEX idx_pharaoh_dynasty ON pharaohs (dynasty);
CREATE INDEX idx_pharaoh_period ON pharaohs (period);
CREATE INDEX idx_pharaoh_reign_start ON pharaohs (reign_start);
CREATE INDEX idx_pharaoh_popularity ON pharaohs (popularity_score);
CREATE INDEX idx_pharaoh_rarity ON pharaohs (rarity);
CREATE INDEX idx_pharaoh_verified ON pharaohs (verified);

-- Add trigger to update updated_at timestamp for pharaohs
CREATE TRIGGER update_pharaohs_updated_at BEFORE UPDATE ON pharaohs
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Tier 1: Core Templates & Sections
CREATE TABLE history_templates (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    title VARCHAR(255) NOT NULL,
    description TEXT,
    era VARCHAR(100),
    dynasty INT,
    pharaoh_id VARCHAR(50),
    difficulty VARCHAR(50) DEFAULT 'Basic',
    thumbnail_url TEXT,
    language VARCHAR(10) DEFAULT 'en',
    is_active BOOLEAN DEFAULT true,
    version INT DEFAULT 1,
    published_at TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (pharaoh_id) REFERENCES pharaohs(id)
);

CREATE TABLE template_sections (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL,
    title VARCHAR(255),
    subtitle VARCHAR(255),
    content_type VARCHAR(50) NOT NULL,  -- TEXT, IMAGE, VIDEO, AUDIO, CHART, TIMELINE
    content TEXT NOT NULL,
    metadata JSONB,           -- e.g., {"caption": "...", "source": "..."}
    order_index INT NOT NULL,
    optional BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES history_templates(id) ON DELETE CASCADE
);

-- Tier 2: Interactive Quizzes, Timeline Events, Media
CREATE TABLE template_quizzes (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id UUID NOT NULL,
    question TEXT NOT NULL,
    options JSONB NOT NULL,
    correct_answer TEXT NOT NULL,
    explanation TEXT,
    difficulty VARCHAR(20) DEFAULT 'Easy',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (section_id) REFERENCES template_sections(id) ON DELETE CASCADE
);

CREATE TABLE template_timelines (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    section_id UUID NOT NULL,
    event_date DATE NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (section_id) REFERENCES template_sections(id) ON DELETE CASCADE
);

-- Tier 3: Source Linking & Knowledge Graph
CREATE TABLE content_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL,
    section_id UUID,
    source_title VARCHAR(255),
    source_url TEXT,
    source_type VARCHAR(50),  -- PRIMARY, SECONDARY, ACADEMIC
    citation_reference TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES history_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (section_id) REFERENCES template_sections(id) ON DELETE CASCADE
);

CREATE TABLE template_references (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL,
    related_template_id UUID NOT NULL,
    relation_type VARCHAR(50),  -- e.g. "See Also", "Prequel", "Consequence"
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES history_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (related_template_id) REFERENCES history_templates(id) ON DELETE CASCADE
);

-- Tier 4: User Learning Flow, Feedback & Gamification
CREATE TABLE user_template_progress (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    template_id UUID NOT NULL,
    section_id UUID,
    progress DECIMAL(5,2) DEFAULT 0.0,
    completed BOOLEAN DEFAULT false,
    last_viewed TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES history_templates(id) ON DELETE CASCADE,
    FOREIGN KEY (section_id) REFERENCES template_sections(id) ON DELETE CASCADE
);

CREATE TABLE user_quiz_responses (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    quiz_id UUID NOT NULL,
    selected_option TEXT NOT NULL,
    is_correct BOOLEAN,
    answered_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (quiz_id) REFERENCES template_quizzes(id) ON DELETE CASCADE
);

CREATE TABLE template_feedback (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    template_id UUID NOT NULL,
    rating DECIMAL(2,1),
    comment TEXT,
    language VARCHAR(10) DEFAULT 'en',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES history_templates(id) ON DELETE CASCADE
);

CREATE TABLE achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    criteria JSONB, -- e.g. {"templates_completed": 5, "quiz_score_avg": 80}
    badge_url TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE user_achievements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    achievement_id UUID NOT NULL,
    unlocked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (achievement_id) REFERENCES achievements(id) ON DELETE CASCADE
);

ALTER TABLE user_achievements ADD CONSTRAINT unique_user_achievement UNIQUE (user_id, achievement_id);

-- Tier 5: Multi-Language & Version Control
CREATE TABLE template_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    base_template_id UUID NOT NULL,
    language VARCHAR(10) NOT NULL,
    title VARCHAR(255),
    description TEXT,
    content_language_links JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (base_template_id) REFERENCES history_templates(id) ON DELETE CASCADE
);

ALTER TABLE template_translations ADD CONSTRAINT unique_template_language UNIQUE (base_template_id, language);

CREATE TABLE template_versions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL,
    version INT NOT NULL,
    change_log TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES history_templates(id) ON DELETE CASCADE
);

-- Tier 6: Tagging, Recommendation, Analytics
CREATE TABLE template_tags (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    template_id UUID NOT NULL,
    tag VARCHAR(50) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (template_id) REFERENCES history_templates(id) ON DELETE CASCADE
);

ALTER TABLE template_tags ADD CONSTRAINT unique_template_tag UNIQUE (template_id, tag);

CREATE TABLE recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    template_id UUID NOT NULL,
    score DECIMAL(5,3),  -- Recommendation score
    context JSONB,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (template_id) REFERENCES history_templates(id) ON DELETE CASCADE
);

-- Indexes for history templates system
CREATE INDEX idx_history_templates_pharaoh_id ON history_templates (pharaoh_id);
CREATE INDEX idx_history_templates_dynasty ON history_templates (dynasty);
CREATE INDEX idx_history_templates_era ON history_templates (era);
CREATE INDEX idx_history_templates_difficulty ON history_templates (difficulty);
CREATE INDEX idx_history_templates_is_active ON history_templates (is_active);
CREATE INDEX idx_template_sections_template_id ON template_sections (template_id);
CREATE INDEX idx_template_sections_order ON template_sections (template_id, order_index);
CREATE INDEX idx_template_quizzes_section_id ON template_quizzes (section_id);
CREATE INDEX idx_template_timelines_section_id ON template_timelines (section_id);
CREATE INDEX idx_user_template_progress_user_id ON user_template_progress (user_id);
CREATE INDEX idx_user_template_progress_template_id ON user_template_progress (template_id);
CREATE INDEX idx_user_quiz_responses_user_id ON user_quiz_responses (user_id);
CREATE INDEX idx_template_feedback_template_id ON template_feedback (template_id);
CREATE INDEX idx_template_tags_template_id ON template_tags (template_id);
CREATE INDEX idx_template_tags_tag ON template_tags (tag);
CREATE INDEX idx_recommendations_user_id ON recommendations (user_id);
CREATE INDEX idx_recommendations_score ON recommendations (score DESC);

-- Add triggers to update updated_at timestamp for new tables
CREATE TRIGGER update_history_templates_updated_at BEFORE UPDATE ON history_templates
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_template_translations_updated_at BEFORE UPDATE ON template_translations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();


INSERT INTO pharaohs (
    id, name, birth_name, throne_name, epithet,
    dynasty, period, reign_start, reign_end, length_of_reign_years,
    father, mother, consorts, children_count, notable_children,
    capital, major_achievements, military_campaigns, building_projects,
    political_style, divine_association, temple_affiliations,
    pharaoh_as_god, burial_site, tomb_discovered, discovery_year,
    famous_artifacts, treasure_status, image_url, statue_count,
    mummy_location, popularity_score, user_rating, unlock_in_game,
    rarity, traits, source, verified
) VALUES

-- 1. Ramses II
(
    'ph_001', 'Ramses II', 'Ra-mes-su', 'Usermaatre Setepenre', 'The Great, Builder of Monuments',
    19, 'New Kingdom', -1279, -1213, 66,
    'Seti I', 'Tuya', '["Nefertari", "Isetnofret"]', 100, '["Merneptah", "Khaemweset"]',
    'Pi-Ramesses', '["Battle of Kadesh", "Signed first peace treaty", "Built Abu Simbel"]',
    '["Kadesh", "Libya", "Nubia"]', '["Ramesseum", "Abu Simbel", "Karnak expansions"]',
    'Authoritarian', '["Amun-Ra", "Horus"]', '["Karnak", "Luxor", "Abu Simbel"]',
    true, 'Valley of the Kings (KV7)', true, 1817,
    '[{"name": "Statue of Ramses II", "museum": "Grand Egyptian Museum"}]',
    'Looted', 'https://example.com/ramses.jpg', 100,
    'Egyptian Museum, Cairo', 9.7, 4.8, true, 'Legendary',
    '{"leadership": 9, "military": 9, "building": 10, "diplomacy": 8}',
    'Temple inscriptions, Manetho, Modern archaeology', true
),

-- 2. Tutankhamun
(
    'ph_002', 'Tutankhamun', 'Nebkheperure', 'Tutankhamun', 'The Boy King',
    18, 'New Kingdom', -1332, -1323, 9,
    'Akhenaten', 'Ankhesenamun (mother?)', '["Ankhesenamun"]', 0, '[]',
    'Thebes', '["Restored worship of Amun", "Repaired temples after Aten era"]',
    '[]', '["Karnak Temple repairs", "Luxor Temple"]',
    'Restorationist', '["Amun-Ra", "Horus"]', '["Karnak", "Luxor"]',
    true, 'Valley of the Kings (KV62)', true, 1922,
    '[{"name": "Golden Mask", "museum": "Grand Egyptian Museum"}, {"name": "Golden Coffin", "museum": "GEM"}]',
    'Intact', 'https://example.com/tut.jpg', 10,
    'Grand Egyptian Museum', 10.0, 4.9, true, 'Legendary',
    '{"leadership": 6, "military": 4, "religion": 9, "symbolism": 10}',
    'Howard Carter discovery, Tomb goods', true
),

-- 3. Hatshepsut
(
    'ph_003', 'Hatshepsut', 'Maatkare', 'Hatshepsut', 'The Female Pharaoh',
    18, 'New Kingdom', -1479, -1458, 21,
    'Thutmose I', 'Ahmes', '["Thutmose II"]', 1, '["Neferure"]',
    'Thebes', '["Built Deir el-Bahari", "Trade expedition to Punt"]',
    '[]', '["Deir el-Bahari", "Karnak"]',
    'Diplomatic', '["Amun-Ra", "Hathor"]', '["Deir el-Bahari", "Karnak"]',
    true, 'Valley of the Kings (KV20)', true, 1902,
    '[{"name": "Obelisks at Karnak", "museum": "On-site"}]',
    'Partially Preserved', 'https://example.com/hatshepsut.jpg', 25,
    'Valley of the Kings (reburied)', 8.9, 4.7, true, 'Epic',
    '{"leadership": 9, "diplomacy": 10, "building": 9, "military": 5}',
    'Inscriptions at Deir el-Bahari, Thutmose III records', true
),

-- 4. Khufu (Cheops)
(
    'ph_004', 'Khufu', 'Khnum-khufwy', 'Khufu', 'Builder of the Great Pyramid',
    4, 'Old Kingdom', -2589, -2566, 23,
    'Sneferu', 'Hetepheres I', '["Meritites I", "Henutsen"]', 8, '["Djedefre", "Khafre"]',
    'Memphis', '["Built the Great Pyramid of Giza"]',
    '[]', '["Great Pyramid", "Pyramid complex"]',
    'Monumental', '["Ra", "Horus"]', '["Giza Plateau"]',
    true, 'Great Pyramid, Giza', true, 1835,
    '[{"name": "Solar Boat", "museum": "Giza Solar Boat Museum"}]',
    'Looted', 'https://example.com/khufu.jpg', 5,
    'Unknown (likely destroyed)', 8.5, 4.5, true, 'Epic',
    '{"leadership": 8, "building": 10, "military": 5, "legacy": 10}',
    'Herodotus, Pyramid inscriptions', true
),

-- 5. Akhenaten
(
    'ph_005', 'Akhenaten', 'Amenhotep IV', 'Neferkheperure Waenre', 'The Heretic King',
    18, 'New Kingdom', -1353, -1336, 17,
    'Amenhotep III', 'Tiye', '["Nefertiti"]', 6, '["Meritaten", "Ankhesenamun"]',
    'Amarna (Akhetaten)', '["Introduced Aten worship", "Founded new capital"]',
    '[]', '["Temple of Aten", "Amarna city"]',
    'Revolutionary', '["Aten"]', '["Amarna"]',
    true, 'Valley of the Kings (KV55)', true, 1907,
    '[{"name": "Amarna Bust of Nefertiti", "museum": "Neues Museum, Berlin"}]',
    'Looted', 'https://example.com/akhenaten.jpg', 3,
    'Possibly KV55 (disputed)', 7.8, 4.3, true, 'Rare',
    '{"religion": 10, "diplomacy": 6, "tradition": 3, "innovation": 10}',
    'Amarna letters, Temple reliefs', true
),

-- 6. Cleopatra VII
(
    'ph_006', 'Cleopatra VII', 'Cleopatra Thea Philopator', 'Cleopatra VII', 'The Last Pharaoh',
    0, 'Ptolemaic Period', -51, -30, 21,
    'Ptolemy XII Auletes', 'Unknown', '["Julius Caesar", "Mark Antony"]', 4, '["Caesarion", "Alexander Helios"]',
    'Alexandria', '["Alliance with Rome", "Preserved Egyptian independence for a time"]',
    '["Naval Battle of Actium"]', '["Palace of Alexandria", "Temple renovations"]',
    'Strategic', '["Isis"]', '["Temple of Isis at Philae"]',
    false, 'Tomb unknown (possibly near Alexandria)', false, NULL,
    '[{"name": "Possible statue in Hermitage", "museum": "Hermitage Museum"}]',
    'Lost', 'https://example.com/cleopatra.jpg', 15,
    'Unknown', 9.2, 4.6, true, 'Legendary',
    '{"diplomacy": 10, "leadership": 8, "survival": 9, "military": 5}',
    'Plutarch, Cassius Dio, Archaeological hints', true
),
(
    'ph_007', 'Thutmose III', 'Menkheperre', 'Thutmose III', 'The Napoleon of Egypt',
    18, 'New Kingdom', -1479, -1425, 54,
    'Thutmose II', 'Iset', '["Merytre-Hatshepsut", "Nebtu"]', 25, '["Amenhotep II"]',
    'Thebes', '["17 military campaigns in Asia", "Expanded Egyptian empire to its peak"]',
    '["Megiddo", "Kadesh", "Syria"]', '["Deir el-Bahari", "Karnak obelisks"]',
    'Military Strategist', '["Amun-Ra", "Montu"]', '["Karnak", "Deir el-Bahari"]',
    true, 'Valley of the Kings (KV34)', true, 1898,
    '[{"name": "Gold Mask of Thutmose III", "museum": "Egyptian Museum"}]',
    'Looted', 'https://example.com/thutmose3.jpg', 30,
    'Deir el-Bahari Cache (DB320)', 8.7, 4.6, true, 'Epic',
    '{"military": 10, "leadership": 9, "strategy": 10, "building": 7}',
    'Karnak war annals, Royal cache records', true
),

-- 8. Amenhotep III (The Magnificent)
(
    'ph_008', 'Amenhotep III', 'Nebmaatre', 'Amenhotep III', 'The Magnificent',
    18, 'New Kingdom', -1386, -1349, 37,
    'Thutmose IV', 'Mutemwiya', '["Tiye", "Gilukhepa"]', 15, '["Akhenaten", "Sitamun"]',
    'Thebes', '["Golden Age of art and diplomacy", "Luxor Temple construction"]',
    '[]', '["Luxor Temple", "Mortuary Temple (Colossi of Memnon)", "Malqata Palace"]',
    'Diplomatic', '["Amun-Ra", "Ra-Horakhty"]', '["Karnak", "Luxor", "West Thebes"]',
    true, 'Valley of the Kings (WV22)', true, 1898,
    '[{"name": "Colossi of Memnon", "museum": "On-site"}]',
    'Partially Preserved', 'https://example.com/amenhotep3.jpg', 200,
    'KV35 Cache (Deir el-Bahari)', 8.4, 4.5, true, 'Epic',
    '{"diplomacy": 10, "art": 10, "building": 9, "military": 5}',
    'Amarna letters, Colossi inscriptions', true
),

-- 9. Seti I (The Great Father)
(
    'ph_009', 'Seti I', 'Menmaatre', 'Seti I', 'The Great Father',
    19, 'New Kingdom', -1290, -1279, 11,
    'Ramesses I', 'Sitre', '["Tuya", "Nefertari?"]', 2, '["Ramses II"]',
    'Thebes / Pi-Ramesses', '["Restored temples", "Military campaigns in Canaan and Libya"]',
    '["Canaan", "Libya", "Hittite frontier"]', '["Abydos Temple", "Karnak Hypostyle Hall"]',
    'Restorationist', '["Amun-Ra", "Osiris"]', '["Abydos", "Karnak", "Thebes"]',
    true, 'Valley of the Kings (KV17)', true, 1817,
    '[{"name": "Sarcophagus of Seti I", "museum": "Sir John Soane Museum, London"}]',
    'Looted', 'https://example.com/seti1.jpg', 50,
    'KV17 (reburied in KV17)', 8.6, 4.6, true, 'Epic',
    '{"leadership": 9, "religion": 9, "building": 8, "military": 8}',
    'Temple reliefs at Abydos and Karnak', true
),

-- 10. Merneptah (Victor of Israel)
(
    'ph_010', 'Merneptah', 'Ba-en-re', 'Merneptah', 'Victor of Israel',
    19, 'New Kingdom', -1213, -1203, 10,
    'Ramses II', 'Isetnofret', '["Isetnofret II", "Wedjemet"]', 12, '["Seti II"]',
    'Pi-Ramesses', '["Defeated the Libyans and Sea Peoples", "First mention of Israel"]',
    '["Libya", "Sea Peoples"]', '["Mortuary Temple (Thebes)", "Karnak"]',
    'Defensive', '["Amun-Ra", "Ptah"]', '["Thebes", "Karnak"]',
    true, 'Valley of the Kings (KV8)', true, 1898,
    '[{"name": "Merneptah Stele", "museum": "Egyptian Museum"}]',
    'Looted', 'https://example.com/merneptah.jpg', 15,
    'KV35 Cache', 7.5, 4.2, true, 'Rare',
    '{"military": 8, "diplomacy": 6, "legacy": 7, "religion": 7}',
    'Merneptah Stele, KV8 inscriptions', true
),

-- 11. Nefertiti (Possible Female Pharaoh - Neferneferuaten)
(
    'ph_011', 'Nefertiti', 'Neferneferuaten Nefertiti', 'Ankhkheperure', 'The Beautiful One Has Come',
    18, 'New Kingdom', -1334, -1330, 4,
    'Unknown (possibly Ay)', 'Tey', '["Akhenaten"]', 6, '["Meritaten", "Ankhesenamun"]',
    'Amarna', '["Co-ruler with Akhenaten", "Possible sole ruler after his death"]',
    '[]', '["Amarna city", "Temple of Aten"]',
    'Revolutionary', '["Aten"]', '["Amarna"]',
    true, 'Tomb unknown (possibly Nefertiti’s bust in Berlin)', true, 1912,
    '[{"name": "Bust of Nefertiti", "museum": "Neues Museum, Berlin"}]',
    'Lost', 'https://example.com/nefertiti.jpg', 1,
    'Unknown', 9.0, 4.7, true, 'Legendary',
    '{"charisma": 10, "influence": 10, "religion": 9, "mystery": 10}',
    'Amarna art, Bust discovery by Borchardt', true
),

-- 12. Ramesses III (Last Great Warrior Pharaoh)
(
    'ph_012', 'Ramesses III', 'Usermaatre Meryamun', 'Ramesses III', 'The Last Great Warrior',
    20, 'New Kingdom', -1186, -1155, 31,
    'Setnakhte', 'Tiy-Merenese', '["Iset Ta-Hemdjert", "Tiy"]', 20, '["Ramesses IV", "Ramesses VI"]',
    'Pi-Ramesses / Thebes', '["Defeated the Sea Peoples", "Built Medinet Habu"]',
    '["Sea Peoples", "Libya"]', '["Medinet Habu", "Temple at Karnak"]',
    'Defensive', '["Amun-Ra", "Horus"]', '["Medinet Habu", "Karnak"]',
    true, 'Valley of the Kings (KV11)', true, 1881,
    '[{"name": "Medinet Habu Reliefs", "museum": "On-site"}]',
    'Looted', 'https://example.com/ramses3.jpg', 40,
    'KV35 Cache', 8.3, 4.4, true, 'Epic',
    '{"military": 9, "survival": 9, "administration": 8, "building": 8}',
    'Medinet Habu reliefs, Judicial Papyrus', true
);
-- === HISTORY TEMPLATES ===
INSERT INTO history_templates (id, title, description, era, dynasty, pharaoh_id, difficulty, thumbnail_url, language, is_active, version, published_at, created_at)
VALUES
  ('00000000-0000-0000-0000-000000000101', 'The Rise of the Old Kingdom', 'Overview of early royal establishment.', 'Old Kingdom', 3, 'ph_004', 'Intermediate', 'https://example.com/thumb1.jpg', 'en', true, 1, NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000102', 'Building the Pyramids', 'Engineering marvels and labor.', 'Old Kingdom', 4, 'ph_004', 'Intermediate', 'https://example.com/thumb2.jpg', 'en', true, 1, NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000103', 'Life in the New Kingdom', 'Daily life and society.', 'New Kingdom', 18, 'ph_001', 'Intermediate', 'https://example.com/thumb3.jpg', 'en', true, 1, NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000104', 'The Decline of the Pharaohs', 'Late period challenges.', 'Late Period', 25, NULL, 'Intermediate', 'https://example.com/thumb4.jpg', 'en', true, 1, NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000105', 'Gods and Rituals of Ancient Egypt', 'Religious beliefs and practices.', 'New Kingdom', 18, 'ph_005', 'Intermediate', 'https://example.com/thumb5.jpg', 'en', true, 1, NOW(), NOW()),
  ('00000000-0000-0000-0000-000000000106', 'Trade and Warfare in the Middle Kingdom', 'Economics and conflict.', 'Middle Kingdom', 12, NULL, 'Intermediate', 'https://example.com/thumb6.jpg', 'en', true, 1, NOW(), NOW());

-- === TEMPLATE SECTIONS, QUIZZES, TIMELINES, CONTENT SOURCES ===
INSERT INTO template_sections (id, template_id, title, subtitle, content_type, content, metadata, order_index, optional, created_at)
VALUES
  ('10000000-0000-0000-0000-000000000111', '00000000-0000-0000-0000-000000000101', 'Origins & Unification', 'What led to the first dynasty?', 'TEXT', 'Content about political unification.', '{"caption":"Unified Egypt","source":"Ancient Chronicle"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-000000000112', '00000000-0000-0000-0000-000000000101', 'Monumental Architecture', 'Early royal tombs and mastabas.', 'TEXT', 'Content about early tombs.', '{"caption":"Mastaba tomb","source":"Archaeological site"}', 2, false, NOW()),

  ('10000000-0000-0000-0000-000000000121', '00000000-0000-0000-0000-000000000102', 'Labor & Logistics', 'How pyramids were built.', 'TEXT', 'Engineering and workforce.', '{"caption":"Construction","source":"Worker inscriptions"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-000000000122', '00000000-0000-0000-0000-000000000102', 'Symbolism & Design', 'Religious meaning behind pyramids.', 'TEXT', 'Symbolic architecture.', '{"caption":"Pyramid alignment","source":"Astronomical study"}', 2, false, NOW()),

  ('10000000-0000-0000-0000-000000000131', '00000000-0000-0000-0000-000000000103', 'Pharaohs & Administration', 'New Kingdom rulers and government.', 'TEXT', 'Administration structure.', '{"caption":"Royal decrees","source":"Papyrus"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-000000000132', '00000000-0000-0000-0000-000000000103', 'Everyday Life', 'Commoners, artisans, and households.', 'TEXT', 'Social customs and daily life.', '{"caption":"Market scene","source":"Tomb painting"}', 2, false, NOW()),

  ('10000000-0000-0000-0000-000000000141', '00000000-0000-0000-0000-000000000104', 'Foreign Invasions', 'Decline due to external threats.', 'TEXT', 'Military challenges.', '{"caption":"Battle relief","source":"Temple carving"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-000000000142', '00000000-0000-0000-0000-000000000104', 'Political Fragmentation', 'Regional powers emerge.', 'TEXT', 'Rise of local rulers.', '{"caption":"Governor","source":"Inscription"}', 2, false, NOW()),

  ('10000000-0000-0000-0000-000000000151', '00000000-0000-0000-0000-000000000105', 'Pantheon & Myths', 'Major gods and myths.', 'TEXT', 'Descriptions of gods.', '{"caption":"Ankh symbol","source":"Relief"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-000000000152', '00000000-0000-0000-0000-000000000105', 'Ritual Practice', 'Temples, offerings, festivals.', 'TEXT', 'Religious ceremonies.', '{"caption":"Temple offerings","source":"Papyrus"}', 2, false, NOW()),

  ('10000000-0000-0000-0000-000000000161', '00000000-0000-0000-0000-000000000106', 'Trade Routes', 'Economic connections abroad.', 'TEXT', 'Import/export networks.', '{"caption":"Trade caravan","source":"Graffiti"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-000000000162', '00000000-0000-0000-0000-000000000106', 'Warfare & Diplomacy', 'Middle Kingdom conflicts and treaties.', 'TEXT', 'Military and alliances.', '{"caption":"Battle map","source":"Temple wall"}', 2, false, NOW());

-- QUIZZES
INSERT INTO template_quizzes (id, section_id, question, options, correct_answer, explanation, difficulty, created_at)
VALUES
  ('20000000-0000-0000-0000-000000000111', '10000000-0000-0000-0000-000000000111', 'Which dynasty unified Upper and Lower Egypt?', '{"a":"1st","b":"3rd","c":"5th"}', 'a', 'Narmer founded the 1st dynasty.', 'Easy', NOW()),
  ('20000000-0000-0000-0000-000000000112', '10000000-0000-0000-0000-000000000112', 'What is a mastaba?', '{"a":"Temple","b":"Flat tomb","c":"Pyramid"}', 'b', 'Mastabas were predecessor tombs.', 'Easy', NOW()),
  ('20000000-0000-0000-0000-000000000121', '10000000-0000-0000-0000-000000000121', 'Which workforce built the pyramids?', '{"a":"Slave labor","b":"Paid laborers","c":"Conscripts"}', 'b', 'Archaeology shows paid workforce.', 'Intermediate', NOW()),
  ('20000000-0000-0000-0000-000000000122', '10000000-0000-0000-0000-000000000122', 'Pyramid alignment aligns with which star?', '{"a":"Sirius","b":"Polaris","c":"Orion’s Belt"}', 'c', 'Aligned with Orion’s Belt.', 'Intermediate', NOW()),
  ('20000000-0000-0000-0000-000000000131', '10000000-0000-0000-0000-000000000131', 'Who ruled during the New Kingdom’s height?', '{"a":"Amenhotep","b":"Ramses II","c":"Tutankhamun"}', 'b', 'Ramses II oversaw major expansion.', 'Easy', NOW()),
  ('20000000-0000-0000-0000-000000000132', '10000000-0000-0000-0000-000000000132', 'Which social class lived in walled towns?', '{"a":"Priests","b":"Artisans","c":"Farmers"}', 'b', 'Artisans clustered in towns.', 'Easy', NOW()),
  ('20000000-0000-0000-0000-000000000141', '10000000-0000-0000-0000-000000000141', 'Which people invaded Egypt in Late Period?', '{"a":"Hyksos","b":"Libyans","c":"Persians"}', 'c', 'Persians conquered Egypt.', 'Intermediate', NOW()),
  ('20000000-0000-0000-0000-000000000142', '10000000-0000-0000-0000-000000000142', 'What led to regional rule in Egypt?', '{"a":"Economic collapse","b":"Central weakens","c":"Religious split"}', 'b', 'The central power weakened.', 'Intermediate', NOW()),
  ('20000000-0000-0000-0000-000000000151', '10000000-0000-0000-0000-000000000151', 'Which god symbolizes life and resurrection?', '{"a":"Osiris","b":"Horus","c":"Anubis"}', 'a', 'Osiris is the god of resurrection.', 'Easy', NOW()),
  ('20000000-0000-0000-0000-000000000152', '10000000-0000-0000-0000-000000000152', 'Where were annual offerings made?', '{"a":"Home","b":"Temple","c":"Marketplace"}', 'b', 'Offerings were at temples.', 'Easy', NOW()),
  ('20000000-0000-0000-0000-000000000161', '10000000-0000-0000-0000-000000000161', 'Trade network included which region?', '{"a":"Nubia","b":"Greece","c":"China"}', 'a', 'Nubia was a key trading partner.', 'Intermediate', NOW()),
  ('20000000-0000-0000-0000-000000000162', '10000000-0000-0000-0000-000000000162', 'Diplomacy was done via what?', '{"a":"Letters","b":"Battles","c":"Marriage alliances"}', 'c', 'Royal marriages sealed treaties.', 'Intermediate', NOW());

-- TIMELINE EVENTS
INSERT INTO template_timelines (id, section_id, event_date, title, description, created_at)
VALUES
  ('30000000-0000-0000-0000-000000000111', '10000000-0000-0000-0000-000000000111', '3100-01-01', 'Narmer unification', 'Legendary unification of Upper and Lower Egypt.', NOW()),
  ('30000000-0000-0000-0000-000000000112', '10000000-0000-0000-0000-000000000112', '2650-01-01', 'First pyramid built', 'Step Pyramid of Djoser.', NOW()),
  ('30000000-0000-0000-0000-000000000121', '10000000-0000-0000-0000-000000000121', '2580-01-01', 'Great Pyramid complete', 'Khufu’s pyramid finished.', NOW()),
  ('30000000-0000-0000-0000-000000000122', '10000000-0000-0000-0000-000000000122', '2500-01-01', 'Pyramids aligned astronomically', 'Celestial alignment insights.', NOW()),
  ('30000000-0000-0000-0000-000000000131', '10000000-0000-0000-0000-000000000131', '1279-01-01', 'Ramses II accession', 'Ramses II begins reign.', NOW()),
  ('30000000-0000-0000-0000-000000000132', '10000000-0000-0000-0000-000000000132', '1250-01-01', 'Artisan settlements flourish', 'Urban trade centers emerge.', NOW()),
  ('30000000-0000-0000-0000-000000000141', '10000000-0000-0000-0000-000000000141', '525-01-01', 'Persian conquest', 'Persians take Egypt.', NOW()),
  ('30000000-0000-0000-0000-000000000142', '10000000-0000-0000-0000-000000000142', '600-01-01', 'Regional governors rise', 'Decentralized control.', NOW()),
  ('30000000-0000-0000-0000-000000000151', '10000000-0000-0000-0000-000000000151', '1500-01-01', 'Horizon of Amun', 'God of Thebes emerges.', NOW()),
  ('30000000-0000-0000-0000-000000000152', '10000000-0000-0000-0000-000000000152', '1250-01-01', 'Temple festivals begin', 'Annual ceremonies start.', NOW()),
  ('30000000-0000-0000-0000-000000000161', '10000000-0000-0000-0000-000000000161', '2000-01-01', 'Nubian gold trade', 'Trade with Nubia peaks.', NOW()),
  ('30000000-0000-0000-0000-000000000162', '10000000-0000-0000-0000-000000000162', '1900-01-01', 'Peace treaty signed', 'Diplomatic marriage alliance.', NOW());

-- CONTENT SOURCES
INSERT INTO content_sources (id, template_id, section_id, source_title, source_url, source_type, citation_reference, created_at)
VALUES
  ('40000000-0000-0000-0000-000000000111', '00000000-0000-0000-0000-000000000101', '10000000-0000-0000-0000-000000000111', 'Papyrus Chronicle', 'https://source.example.com/ancient', 'PRIMARY', 'Chronicle‑I', NOW()),
  ('40000000-0000-0000-0000-000000000112', '00000000-0000-0000-0000-000000000101', '10000000-0000-0000-0000-000000000112', 'Excavation Journal', 'https://source.example.com/archaeology', 'SECONDARY', 'JournalX', NOW()),
  -- similarly add for each section as needed...

-- TEMPLATE TRANSLATIONS
INSERT INTO template_translations (id, base_template_id, language, title, description, content_language_links, created_at, updated_at)
VALUES
  ('50000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000101', 'es', 'Ascenso del Reino Antiguo', 'Descripción en español.', '{}', NOW(), NOW()),
  -- repeat for each template...

-- TEMPLATE VERSIONS
INSERT INTO template_versions (id, template_id, version, change_log, created_at)
VALUES
  ('60000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000101', 1, 'Initial version', NOW()),
  -- repeat for each template...

-- TEMPLATE TAGS
INSERT INTO template_tags (id, template_id, tag, created_at)
VALUES
  ('70000000-0000-0000-0000-000000000101', '00000000-0000-0000-0000-000000000101', 'Architecture', NOW()),
  ('70000000-0000-0000-0000-000000000102', '00000000-0000-0000-0000-000000000102', 'Labor', NOW()),
  ('70000000-0000-0000-0000-000000000103', '00000000-0000-0000-0000-000000000103', 'Daily Life', NOW()),
  ('70000000-0000-0000-0000-000000000104', '00000000-0000-0000-0000-000000000104', 'Conflict', NOW()),
  ('70000000-0000-0000-0000-000000000105', '00000000-0000-0000-0000-000000000105', 'Religion', NOW()),
  ('70000000-0000-0000-0000-000000000106', '00000000-0000-0000-0000-000000000106', 'Trade', NOW());

-- TEMPLATE REFERENCES
INSERT INTO template_references (id, template_id, related_template_id, relation_type, created_at)
VALUES
  ('80000000-0000-0000-0000-000000000102', '00000000-0000-0000-0000-000000000102', '00000000-0000-0000-0000-000000000101', 'Prequel', NOW()),
  ('80000000-0000-0000-0000-000000000103', '00000000-0000-0000-0000-000000000103', '00000000-0000-0000-0000-000000000102', 'See Also', NOW()),


-- === SIX NEW HISTORY TEMPLATES ===
INSERT INTO history_templates (id, title, description, era, dynasty, pharaoh_id, difficulty, thumbnail_url, language, is_active, version, published_at, created_at)
VALUES
  -- Template 7
  ('00000000-0000-0000-0000-000000000107', 'Queen Hatshepsut’s Reign', 'Her life and projects.', 'New Kingdom', 18, 'PHARAOH7', 'Intermediate', 'https://example.com/thumb7.jpg', 'en', true, 1, NOW(), NOW()),
  -- Template 8
  ('00000000-0000-0000-0000-000000000108', 'Tutankhamun’s Tomb Discovery', 'Discovery and legacy.', 'New Kingdom', 18, 'PHARAOH8', 'Intermediate', 'https://example.com/thumb8.jpg', 'en', true, 1, NOW(), NOW()),
  -- Template 9
  ('00000000-0000-0000-0000-000000000109', 'Hyksos Period & Innovations', 'Foreign rulers and technology.', 'Second Intermediate Period', 15, 'PHARAOH9', 'Intermediate', 'https://example.com/thumb9.jpg', 'en', true, 1, NOW(), NOW()),
  -- Template 10
  ('00000000-0000-0000-0000-00000000010A', 'Akhenaten & Aten Worship', 'Religious revolution.', 'Amarna Period', 18, 'PHARAOH10', 'Intermediate', 'https://example.com/thumb10.jpg', 'en', true, 1, NOW(), NOW()),


-- === SECTIONS FOR EACH TEMPLATE (2 each) ===
INSERT INTO template_sections (id, template_id, title, subtitle, content_type, content, metadata, order_index, optional, created_at)
VALUES
  -- Template 7 sections
  ('10000000-0000-0000-0000-000000000171', '00000000-0000-0000-0000-000000000107', 'Rise to Power', 'How Hatshepsut claimed the throne.', 'TEXT', 'Details of her accession.', '{"caption":"Hatshepsut statue","source":"Reliefs"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-000000000172', '00000000-0000-0000-0000-000000000107', 'Building Projects', 'Karnak, Deir el‑Bahri.', 'TEXT', 'Her construction legacy.', '{"caption":"Mortuary temple","source":"Archaeology"}', 2, false, NOW()),

  -- Template 8 sections
  ('10000000-0000-0000-0000-000000000181', '00000000-0000-0000-0000-000000000108', 'Discovery by Carter', 'Howard Carter’s 1922 find.', 'TEXT', 'Expedition story.', '{"caption":"Tut’s tomb","source":"Diary"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-000000000182', '00000000-0000-0000-0000-000000000108', 'Cultural Impact', 'King Tut in popular culture.', 'TEXT', 'Global fascination.', '{"caption":"Exhibition poster","source":"Museum archives"}', 2, false, NOW()),

  -- Template 9 sections
  ('10000000-0000-0000-0000-000000000191', '00000000-0000-0000-0000-000000000109', 'Hyksos Origins', 'Who were the foreign rulers?', 'TEXT', 'Origins and arrival.', '{"caption":"Hyksos pottery","source":"Excavations"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-000000000192', '00000000-0000-0000-0000-000000000109', 'Technological Advances', 'Bronze and chariots.', 'TEXT', 'New military tech.', '{"caption":"Chariot relief","source":"Temple art"}', 2, false, NOW()),

  -- Template 10 sections
  ('10000000-0000-0000-0000-0000000001A1', '00000000-0000-0000-0000-00000000010A', 'Religious Shift', 'Aten replaces Amun.', 'TEXT', 'Monotheistic experiment.', '{"caption":"Aten disk","source":"Tomb art"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-0000000001A2', '00000000-0000-0000-0000-00000000010A', 'Art Style Changes', 'Amarna art revolution.', 'TEXT', 'Style of depictions.', '{"caption":"Amarna portrait","source":"Art history"}', 2, false, NOW()),

  -- Template 11 sections
  ('10000000-0000-0000-0000-0000000001B1', '00000000-0000-0000-0000-00000000010B', 'Political Alliances', 'Rome and Egypt.', 'TEXT', 'Cleopatra and Caesar.', '{"caption":"Bust of Cleopatra","source":"Numismatic"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-0000000001B2', '00000000-0000-0000-0000-00000000010B', 'Cultural Syncretism', 'Greek‑Egyptian fusion.', 'TEXT', 'Ptolemaic blending.', '{"caption":"Temple relief","source":"Hellenistic art"}', 2, false, NOW()),

  -- Template 12 sections
  ('10000000-0000-0000-0000-0000000001C1', '00000000-0000-0000-0000-00000000010C', 'Discovery Context', 'Napoleon to British expeditions.', 'TEXT', 'How Rosetta was found.', '{"caption":"Rosetta Stone","source":"British Museum"}', 1, false, NOW()),
  ('10000000-0000-0000-0000-0000000001C2', '00000000-0000-0000-0000-00000000010C', 'Decipherment', 'Champollion’s breakthrough.', 'TEXT', 'Decoding hieroglyphs.', '{"caption":"Champollion","source":"Linguistic history"}', 2, false, NOW());

-- === QUIZZES (one per section) ===
INSERT INTO template_quizzes (id, section_id, question, options, correct_answer, explanation, difficulty, created_at)
VALUES
  -- For template 7
  ('20000000-0000-0000-0000-000000000171', '10000000-0000-0000-0000-000000000171', 'Was Hatshepsut a pharaoh by lineage or marriage?', '{"a":"Lineage","b":"Marriage","c":"Self‑proclaim"}', 'c', 'She declared herself pharaoh.', 'Intermediate', NOW()),
  ('20000000-0000-0000-0000-000000000172', '10000000-0000-0000-0000-000000000172', 'Which temple was built by Hatshepsut?', '{"a":"Luxor","b":"Karnak","c":"Deir el‑Bahri"}', 'c', 'Her mortuary temple is at Deir el‑Bahri.', 'Easy', NOW()),
  -- template 8
  ('20000000-0000-0000-0000-000000000181', '10000000-0000-0000-0000-000000000181', 'Who discovered Tut’s tomb?', '{"a":"Carter","b":"Howard","c":"Kohl"}', 'a', 'Howard Carter led the excavation.', 'Easy', NOW()),
  ('20000000-0000-0000-0000-000000000182', '10000000-0000-0000-0000-000000000182', 'Tut became a global icon through what?', '{"a":"Movies","b":"Exhibits","c":"Books"}', 'b', 'Major museum exhibits sparked fame.', 'Easy', NOW()),
  -- template 9
  ('20000000-0000-0000-0000-000000000191', '10000000-0000-0000-0000-000000000191', 'Hyksos introduced which weapon?', '{"a":"Bronze swords","b":"Iron weapons","c":"Chariots"}', 'c', 'They brought chariots to Egypt.', 'Intermediate', NOW()),
  ('20000000-0000-0000-0000-000000000192', '10000000-0000-0000-0000-000000000192', 'Hyksos settlement ruled from where?', '{"a":"Avaris","b":"Memphis","c":"Thebes"}', 'a', 'Their capital was Avaris.', 'Intermediate', NOW()),
  -- template 10
  ('20000000-0000-0000-0000-0000000001A1', '10000000-0000-0000-0000-0000000001A1', 'Aten worship was monotheistic: true or false?', '{"a":"True","b":"False"}', 'a', 'He created monotheism around Aten.', 'Intermediate', NOW()),
  ('20000000-0000-0000-0000-0000000001A2', '10000000-0000-0000-0000-0000000001A2', 'Amarna art emphasizes what?', '{"a":"Realism","b":"Formality","c":"Symmetry"}', 'a', 'Art showed more realism and intimacy.', 'Intermediate', NOW()),
  -- template 11
  ('20000000-0000-0000-0000-0000000001B1', '10000000-0000-0000-0000-0000000001B1', 'Cleopatra allied with which Roman leader?', '{"a":"Octavian","b":"Mark Antony","c":"Pompey"}', 'b', 'She famously allied with Mark Antony.', 'Easy', NOW()),
  ('20000000-0000-0000-0000-0000000001B2', '10000000-0000-0000-0000-0000000001B2', 'Ptolemaic art fused which styles?', '{"a":"Egyptian & Greek","b":"Egyptian & Roman","c":"Greek & Roman"}', 'a', 'A blend of Egyptian and Greek motifs.', 'Intermediate', NOW()),
  -- template 12
  ('20000000-0000-0000-0000-0000000001C1', '10000000-0000-0000-0000-0000000001C1', 'Who found the Rosetta Stone?', '{"a":"Napoléon","b":"Britishers","c":"Egyptians"}', 'a', 'French troops discovered it.', 'Easy', NOW()),
  ('20000000-0000-0000-0000-0000000001C2', '10000000-0000-0000-0000-0000000001C2', 'Champollion deciphered hieroglyphs by comparing which languages?', '{"a":"Greek & Demotic","b":"Coptic & Latin","c":"Greek & Latin"}', 'a', 'Greek alongside Demotic led to the breakthrough.', 'Intermediate', NOW());

-- === TIMELINES ===
INSERT INTO template_timelines (id, section_id, event_date, title, description, created_at)
VALUES
  -- template 7
  ('30000000-0000-0000-0000-000000000171', '10000000-0000-0000-0000-000000000171', '1473-01-01', 'Hatshepsut coronation', 'Proclaimed herself pharaoh.', NOW()),
  ('30000000-0000-0000-0000-000000000172', '10000000-0000-0000-0000-000000000172', '1460-01-01', 'Temple built', 'Mortuary temple launched.', NOW()),
  -- template 8
  ('30000000-0000-0000-0000-000000000181', '10000000-0000-0000-0000-000000000181', '1922-11-04', 'Tutankhamun’s tomb opened', 'Discovery by Carter.', NOW()),
  ('30000000-0000-0000-0000-000000000182', '10000000-0000-0000-0000-000000000182', '1972-01-01', 'Tut exhibitions global', 'Museum shows launch worldwide.', NOW()),
  -- template 9
  ('30000000-0000-0000-0000-000000000191', '10000000-0000-0000-0000-000000000191', '1650-01-01', 'Hyksos rule begins', 'Hyksos invade Eastern Delta.', NOW()),
  ('30000000-0000-0000-0000-000000000192', '10000000-0000-0000-0000-000000000192', '1550-01-01', 'Hyksos expelled', 'Egypt regains control.', NOW()),
  -- template 10
  ('30000000-0000-0000-0000-0000000001A1', '10000000-0000-0000-0000-0000000001A1', '1353-01-01', 'Akhenaten reign begins', 'Introduced Aten worship.', NOW()),
  ('30000000-0000-0000-0000-0000000001A2', '10000000-0000-0000-0000-0000000001A2', '1336-01-01', 'Amarna city rises', 'New capital built.', NOW()),
  -- template 11
  ('30000000-0000-0000-0000-0000000001B1', '10000000-0000-0000-0000-0000000001B1', '51 BC-01-01', 'Cleopatra ascends', 'Begins sole rule.', NOW()),
  ('30000000-0000-0000-0000-0000000001B2', '30000000-0000-0000-0000-0000000001B2', '30 BC-01-01', 'Death & Egypt annexed', 'Rome annexes Egypt.', NOW()),
  -- template 12
  ('30000000-0000-0000-0000-0000000001C1', '10000000-0000-0000-0000-0000000001C1', '1799-07-15', 'Rosetta Stone found', 'Napoleon’s Egyptian campaign.', NOW()),
  ('30000000-0000-0000-0000-0000000001C2', '30000000-0000-0000-0000-0000000001C2', '1822-09-14', 'Champollion deciphers', 'Published decipherment.', NOW());
