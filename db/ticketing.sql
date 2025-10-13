-- ===============================================
-- TWC OTA Ticketing System Database Schema
-- ===============================================
-- Generated from Go entities analysis
-- Date: October 13, 2025
-- ===============================================

-- Enable UUID extension for PostgreSQL
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ===============================================
-- USER MANAGEMENT TABLES
-- ===============================================

-- Users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    type VARCHAR(50) NOT NULL,
    typeid INTEGER,
    users_extras TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Password reset tokens
CREATE TABLE password_resets (
    id SERIAL PRIMARY KEY,
    token VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);

-- ===============================================
-- AGENT MANAGEMENT TABLES
-- ===============================================

-- Master agents table
CREATE TABLE master_agents (
    agent_id SERIAL PRIMARY KEY,
    agent_name VARCHAR(255) NOT NULL,
    agent_address TEXT,
    agent_group_id INTEGER DEFAULT NULL,
    agent_extras TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Agent groups (referenced by master_agents)
CREATE TABLE agent_groups (
    group_id SERIAL PRIMARY KEY,
    group_name VARCHAR(255) NOT NULL,
    group_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ===============================================
-- MASTER DATA TABLES
-- ===============================================

-- Master groups (tourism sites/destinations)
CREATE TABLE master_group (
    group_id SERIAL PRIMARY KEY,
    group_mid VARCHAR(50) UNIQUE NOT NULL,
    group_name VARCHAR(255) NOT NULL,
    group_label VARCHAR(255),
    group_logo VARCHAR(500),
    group_estimate VARCHAR(100),
    description TEXT,
    lat VARCHAR(50),
    long VARCHAR(50),
    adult_age VARCHAR(20),
    child_age VARCHAR(20),
    how_to_use_ticket TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Master tickets
CREATE TABLE master_ticket (
    mtick_id SERIAL PRIMARY KEY,
    mtick_name VARCHAR(255) NOT NULL,
    mtick_code VARCHAR(50),
    mtick_type VARCHAR(50),
    mtick_cat VARCHAR(50),
    mtick_group_id INTEGER REFERENCES master_group(group_id),
    mtick_loc_start INTEGER,
    mtick_loc_finish INTEGER,
    mtick_merchant_code VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Tariff types
CREATE TABLE master_tariff_type (
    trfftype_id SERIAL PRIMARY KEY,
    trfftype_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Master tariffs
CREATE TABLE master_tariff (
    trf_id SERIAL PRIMARY KEY,
    trf_name VARCHAR(255) NOT NULL,
    trf_code VARCHAR(50),
    trf_agent_id INTEGER REFERENCES master_agents(agent_id),
    trf_trftype INTEGER REFERENCES master_tariff_type(trfftype_id),
    trf_currency_code VARCHAR(3) DEFAULT 'IDR',
    trf_value DECIMAL(15,2) NOT NULL,
    trf_label VARCHAR(255),
    trf_start_date DATE,
    trf_end_date DATE,
    trf_priority INTEGER DEFAULT 0,
    trf_release VARCHAR(50),
    trf_qty INTEGER DEFAULT 0,
    trf_tax DECIMAL(15,2) DEFAULT 0,
    trf_assurance DECIMAL(15,2) DEFAULT 0,
    trf_fix_price DECIMAL(15,2) DEFAULT 0,
    trf_admin DECIMAL(15,2) DEFAULT 0,
    trf_others DECIMAL(15,2) DEFAULT 0,
    day_condition VARCHAR(50),
    begin_time TIME,
    end_time TIME,
    card_type VARCHAR(50),
    expired_qr INTEGER DEFAULT 24,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Master tariff details (linking tariffs to tickets)
CREATE TABLE master_tariffdet (
    trfdet_id SERIAL PRIMARY KEY,
    trfdet_trf_id INTEGER REFERENCES master_tariff(trf_id),
    trfdet_mtick_id INTEGER REFERENCES master_ticket(mtick_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Currency rates
CREATE TABLE currency (
    curr_id SERIAL PRIMARY KEY,
    curr_code VARCHAR(3) NOT NULL,
    curr_rate DECIMAL(15,6) NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ===============================================
-- DISCOUNT MANAGEMENT TABLES
-- ===============================================

-- Multi-destination discounts
CREATE TABLE master_discount_multi (
    discm_id SERIAL PRIMARY KEY,
    discm_name VARCHAR(255) NOT NULL,
    discm_start_date DATE NOT NULL,
    discm_end_date DATE NOT NULL,
    discm_destination INTEGER NOT NULL,
    discm_value DECIMAL(15,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ===============================================
-- BOOKING SYSTEM TABLES
-- ===============================================

-- Main booking table
CREATE TABLE booking (
    booking_id SERIAL,
    agent_id INTEGER REFERENCES master_agents(agent_id),
    booking_number VARCHAR(100) NOT NULL,
    booking_date DATE NOT NULL,
    booking_mid VARCHAR(50),
    booking_payment_method VARCHAR(50),
    booking_amount DECIMAL(15,2) NOT NULL,
    booking_emoney INTEGER DEFAULT 0,
    booking_total_payment DECIMAL(15,2) NOT NULL,
    booking_uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    booking_redeem_date DATE DEFAULT NULL,
    booking_invoice VARCHAR(100),
    customer_note TEXT,
    customer_email VARCHAR(255),
    customer_username VARCHAR(255),
    customer_phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Booking details
CREATE TABLE bookingdet (
    bookingdet_id SERIAL,
    bookingdet_booking_id INTEGER,
    bookingdet_trf_id INTEGER REFERENCES master_tariff(trf_id),
    bookingdet_trftype VARCHAR(50),
    bookingdet_amount DECIMAL(15,2) NOT NULL,
    bookingdet_qty INTEGER NOT NULL,
    bookingdet_total DECIMAL(15,2) NOT NULL,
    bookingdet_uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bookingdet_booking_uuid UUID REFERENCES booking(booking_uuid),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Booking list (individual tickets)
CREATE TABLE bookinglist (
    bookinglist_id SERIAL,
    bookinglist_bookingdet_id INTEGER,
    bookinglist_mtick_id INTEGER REFERENCES master_ticket(mtick_id),
    bookinglist_mid VARCHAR(50),
    bookinglist_uuid UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    bookinglist_bookingdet_uuid UUID REFERENCES bookingdet(bookingdet_uuid),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ===============================================
-- TRIP PLANNER TABLES
-- ===============================================

-- Main trip planner table
CREATE TABLE trip_planner (
    tp_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tp_stan INTEGER,
    tp_number VARCHAR(100) NOT NULL,
    tp_src_type INTEGER DEFAULT 1,
    tp_start_date DATE NOT NULL,
    tp_end_date DATE NOT NULL,
    tp_duration INTEGER NOT NULL,
    tp_status INTEGER DEFAULT 0,
    tp_user_id INTEGER REFERENCES users(id),
    tp_adult INTEGER DEFAULT 0,
    tp_child INTEGER DEFAULT 0,
    tp_contact TEXT,
    tp_total_amount DECIMAL(15,2) NOT NULL,
    tp_extras TEXT,
    tp_invoice INTEGER,
    tp_agent_id INTEGER REFERENCES master_agents(agent_id),
    tp_payment_method VARCHAR(50) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT NULL
);

-- Trip planner persons
CREATE TABLE trip_planner_person (
    tpp_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tpp_tp_id UUID REFERENCES trip_planner(tp_id),
    tpp_name VARCHAR(255) NOT NULL,
    tpp_type INTEGER NOT NULL, -- 1=adult, 2=child
    tpp_qr VARCHAR(255),
    tpp_extras TEXT,
    id_number VARCHAR(50),
    title VARCHAR(10),
    type_id VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Trip planner destinations
CREATE TABLE trip_planner_destination (
    tpd_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tpd_tpp_id UUID REFERENCES trip_planner_person(tpp_id),
    tpd_group_mid VARCHAR(50) REFERENCES master_group(group_mid),
    tpd_trf_id INTEGER REFERENCES master_tariff(trf_id),
    tpd_amount DECIMAL(15,2) NOT NULL,
    tpd_date DATE NOT NULL,
    tpd_exp_date DATE NOT NULL,
    tpd_day INTEGER NOT NULL,
    tpd_duration INTEGER NOT NULL,
    tpd_extras TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ===============================================
-- TICKETING SYSTEM TABLES
-- ===============================================

-- Main ticket table
CREATE TABLE ticket (
    tick_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tick_stan INTEGER,
    tick_number VARCHAR(100) NOT NULL,
    tick_mid VARCHAR(50),
    tick_src_type INTEGER DEFAULT 1,
    tick_src_id VARCHAR(100),
    tick_src_inv_num VARCHAR(100),
    tick_amount DECIMAL(15,2) NOT NULL,
    tick_emoney INTEGER DEFAULT 0,
    tick_purc VARCHAR(50),
    tick_issuing VARCHAR(50),
    tick_date DATE NOT NULL,
    tick_total_payment DECIMAL(15,2) NOT NULL,
    tick_payment_method VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Ticket details
CREATE TABLE ticketdet (
    tickdet_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tickdet_tick_id UUID REFERENCES ticket(tick_id),
    tickdet_trf_id INTEGER REFERENCES master_tariff(trf_id),
    tickdet_trftype VARCHAR(50),
    tickdet_amount DECIMAL(15,2) NOT NULL,
    tickdet_qty INTEGER NOT NULL,
    tickdet_total DECIMAL(15,2) NOT NULL,
    tickdet_qr VARCHAR(255),
    ext TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Ticket list (individual ticket items)
CREATE TABLE ticketlist (
    ticklist_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ticklist_tickdet_id UUID REFERENCES ticketdet(tickdet_id),
    ticklist_mtick_id INTEGER REFERENCES master_ticket(mtick_id),
    ticklist_expire TIMESTAMP,
    ticklist_mid VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ===============================================
-- INVENTORY MANAGEMENT TABLES
-- ===============================================

-- OTA Inventory
CREATE TABLE ota_inventory (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    inventory_number VARCHAR(100) NOT NULL,
    agent_id INTEGER REFERENCES master_agents(agent_id),
    agent_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- OTA Inventory Details
CREATE TABLE ota_inventory_detail (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    ota_inventory_id UUID REFERENCES ota_inventory(id),
    group_id INTEGER REFERENCES master_group(group_id),
    group_name VARCHAR(255),
    trf_id INTEGER REFERENCES master_tariff(trf_id),
    trf_name VARCHAR(255),
    expiry_date TIMESTAMP,
    qr VARCHAR(255),
    trf_amount DECIMAL(15,2),
    qr_prefix VARCHAR(20),
    redeem_date TIMESTAMP DEFAULT NULL,
    void_date TIMESTAMP DEFAULT NULL,
    trf_type VARCHAR(50),
    group_mid VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ===============================================
-- USER PREFERENCES TABLES
-- ===============================================

-- User favorites
CREATE TABLE favorite (
    fav_id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    fav_user_id INTEGER REFERENCES users(id),
    fav_data TEXT NOT NULL,
    fav_extras TEXT DEFAULT NULL,
    fav_created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    fav_deleted TIMESTAMP DEFAULT NULL
);

-- ===============================================
-- NOTIFICATION SYSTEM TABLES
-- ===============================================

-- Inbox notifications
CREATE TABLE inbox_notification (
    inbox_id SERIAL PRIMARY KEY,
    agent_id INTEGER REFERENCES master_agents(agent_id),
    inbox_title VARCHAR(255),
    inbox_short_desc TEXT,
    inbox_full_desc TEXT,
    inbox_image_url VARCHAR(500),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ===============================================
-- CONFIGURATION TABLES
-- ===============================================

-- Mobile app configuration
CREATE TABLE mobile_config (
    mconfig_id SERIAL PRIMARY KEY,
    mconfig_src_type INTEGER DEFAULT 1,
    mconfig_value TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- ===============================================
-- INDEXES FOR PERFORMANCE
-- ===============================================

-- User indexes
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_type ON users(type);

-- Agent indexes
CREATE INDEX idx_agents_group ON master_agents(agent_group_id);

-- Master data indexes
CREATE INDEX idx_group_mid ON master_group(group_mid);
CREATE INDEX idx_ticket_group ON master_ticket(mtick_group_id);
CREATE INDEX idx_tariff_agent ON master_tariff(trf_agent_id);
CREATE INDEX idx_tariff_type ON master_tariff(trf_trftype);

-- Booking indexes
CREATE INDEX idx_booking_agent ON booking(agent_id);
CREATE INDEX idx_booking_date ON booking(booking_date);
CREATE INDEX idx_booking_number ON booking(booking_number);
CREATE INDEX idx_bookingdet_booking ON bookingdet(bookingdet_booking_uuid);
CREATE INDEX idx_bookinglist_det ON bookinglist(bookinglist_bookingdet_uuid);

-- Trip planner indexes
CREATE INDEX idx_trip_agent ON trip_planner(tp_agent_id);
CREATE INDEX idx_trip_user ON trip_planner(tp_user_id);
CREATE INDEX idx_trip_date ON trip_planner(tp_start_date, tp_end_date);
CREATE INDEX idx_trip_person_trip ON trip_planner_person(tpp_tp_id);
CREATE INDEX idx_trip_dest_person ON trip_planner_destination(tpd_tpp_id);

-- Ticket indexes
CREATE INDEX idx_ticket_date ON ticket(tick_date);
CREATE INDEX idx_ticket_number ON ticket(tick_number);
CREATE INDEX idx_ticketdet_ticket ON ticketdet(tickdet_tick_id);
CREATE INDEX idx_ticketlist_det ON ticketlist(ticklist_tickdet_id);

-- Inventory indexes
CREATE INDEX idx_inventory_agent ON ota_inventory(agent_id);
CREATE INDEX idx_inventory_detail_inv ON ota_inventory_detail(ota_inventory_id);
CREATE INDEX idx_inventory_qr ON ota_inventory_detail(qr);

-- Favorite indexes
CREATE INDEX idx_favorite_user ON favorite(fav_user_id);

-- ===============================================
-- FOREIGN KEY CONSTRAINTS
-- ===============================================

-- Agent group foreign key
ALTER TABLE master_agents 
ADD CONSTRAINT fk_agent_group 
FOREIGN KEY (agent_group_id) REFERENCES agent_groups(group_id);

-- ===============================================
-- SAMPLE DATA INSERTS (OPTIONAL)
-- ===============================================

-- Insert default currency
INSERT INTO currency (curr_code, curr_rate) VALUES ('IDR', 1.000000);
INSERT INTO currency (curr_code, curr_rate) VALUES ('USD', 15000.000000);

-- Insert default tariff types
INSERT INTO master_tariff_type (trfftype_name) VALUES ('Adult');
INSERT INTO master_tariff_type (trfftype_name) VALUES ('Child');
INSERT INTO master_tariff_type (trfftype_name) VALUES ('Senior');

-- Insert default agent group
INSERT INTO agent_groups (group_name, group_description) 
VALUES ('Default Group', 'Default agent group');

-- ===============================================
-- VIEWS FOR COMMON QUERIES
-- ===============================================

-- View for complete booking information
CREATE VIEW v_booking_complete AS
SELECT 
    b.booking_uuid,
    b.booking_number,
    b.booking_date,
    b.booking_total_payment,
    b.customer_email,
    b.customer_username,
    a.agent_name,
    bd.bookingdet_amount,
    bd.bookingdet_qty,
    mt.trf_name,
    mg.group_name
FROM booking b
LEFT JOIN master_agents a ON b.agent_id = a.agent_id
LEFT JOIN bookingdet bd ON b.booking_uuid = bd.bookingdet_booking_uuid
LEFT JOIN master_tariff mt ON bd.bookingdet_trf_id = mt.trf_id
LEFT JOIN master_tariffdet mtd ON mt.trf_id = mtd.trfdet_trf_id
LEFT JOIN master_ticket mti ON mtd.trfdet_mtick_id = mti.mtick_id
LEFT JOIN master_group mg ON mti.mtick_group_id = mg.group_id;

-- View for trip planner summary
CREATE VIEW v_trip_summary AS
SELECT 
    tp.tp_id,
    tp.tp_number,
    tp.tp_start_date,
    tp.tp_end_date,
    tp.tp_duration,
    tp.tp_total_amount,
    tp.tp_status,
    a.agent_name,
    COUNT(tpp.tpp_id) as total_persons,
    COUNT(DISTINCT tpd.tpd_group_mid) as total_destinations
FROM trip_planner tp
LEFT JOIN master_agents a ON tp.tp_agent_id = a.agent_id
LEFT JOIN trip_planner_person tpp ON tp.tp_id = tpp.tpp_tp_id
LEFT JOIN trip_planner_destination tpd ON tpp.tpp_id = tpd.tpd_tpp_id
GROUP BY tp.tp_id, tp.tp_number, tp.tp_start_date, tp.tp_end_date, 
         tp.tp_duration, tp.tp_total_amount, tp.tp_status, a.agent_name;

-- ===============================================
-- STORED PROCEDURES (POSTGRESQL FUNCTIONS)
-- ===============================================

-- Function to get ticket expiry based on QR settings
CREATE OR REPLACE FUNCTION get_ticket_expiry(trf_id_param INTEGER, issue_date TIMESTAMP)
RETURNS TIMESTAMP AS $$
DECLARE
    expired_qr_hours INTEGER;
BEGIN
    SELECT COALESCE(expired_qr, 24) INTO expired_qr_hours
    FROM master_tariff 
    WHERE trf_id = trf_id_param;
    
    RETURN issue_date + (expired_qr_hours || ' hours')::INTERVAL;
END;
$$ LANGUAGE plpgsql;

-- Function to calculate discount for multiple destinations
CREATE OR REPLACE FUNCTION calculate_multi_destination_discount(
    dest_count INTEGER,
    visit_date DATE,
    base_amount DECIMAL(15,2)
) RETURNS DECIMAL(15,2) AS $$
DECLARE
    discount_value DECIMAL(15,2) := 0;
BEGIN
    SELECT discm_value INTO discount_value
    FROM master_discount_multi
    WHERE discm_destination <= dest_count
      AND visit_date BETWEEN discm_start_date AND discm_end_date
    ORDER BY discm_destination DESC
    LIMIT 1;
    
    RETURN COALESCE(discount_value, 0);
END;
$$ LANGUAGE plpgsql;

-- ===============================================
-- TRIGGERS FOR AUDIT TRAILS
-- ===============================================

-- Trigger to update booking timestamp
CREATE OR REPLACE FUNCTION update_booking_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.created_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER booking_timestamp_trigger
    BEFORE INSERT ON booking
    FOR EACH ROW
    EXECUTE FUNCTION update_booking_timestamp();

-- ===============================================
-- COMMENTS FOR DOCUMENTATION
-- ===============================================

COMMENT ON TABLE users IS 'System users including agents and customers';
COMMENT ON TABLE master_agents IS 'Travel agents who can create bookings';
COMMENT ON TABLE master_group IS 'Tourism destinations/sites';
COMMENT ON TABLE master_ticket IS 'Available ticket types';
COMMENT ON TABLE master_tariff IS 'Pricing information for tickets';
COMMENT ON TABLE booking IS 'Main booking transactions';
COMMENT ON TABLE trip_planner IS 'Multi-day trip planning system';
COMMENT ON TABLE ticket IS 'Individual ticket purchases';
COMMENT ON TABLE ota_inventory IS 'Pre-purchased ticket inventory';
COMMENT ON TABLE favorite IS 'User saved favorite destinations';

-- ===============================================
-- END OF SCHEMA
-- ===============================================