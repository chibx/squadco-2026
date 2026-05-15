-- 01.vendors.sql
-- Tables for vendor documents, KYC, verification, payments, and admin

CREATE TABLE business (
    business_id BIGINT PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT NOT NULL,
    tin TEXT,
    nin TEXT NOT NULL UNIQUE,
    business_type TEXT NOT NULL,
    rc_number TEXT NOT NULL,
    address TEXT NOT NULL,
    status TEXT NOT NULL,
    vendor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_business_vendor_id ON business(vendor_id);

CREATE TABLE vendor_documents (
    document_id TEXT PRIMARY KEY,
    vendor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    document_type TEXT NOT NULL,
    document_name TEXT NOT NULL,
    public_id TEXT NOT NULL,
    verification_status TEXT NOT NULL DEFAULT 'pending',
    verified_by TEXT,
    verified_at TIMESTAMP,
    expiry_date TIMESTAMP,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);


CREATE INDEX idx_vendor_documents_vendor_id ON vendor_documents(vendor_id);

CREATE TABLE vendor_kyc (
    kyc_id TEXT PRIMARY KEY,
    vendor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    nin_number TEXT NOT NULL,
    bank_name TEXT NOT NULL,
    account_number TEXT NOT NULL,
    account_name TEXT NOT NULL,
    latitude DOUBLE PRECISION,
    longitude DOUBLE PRECISION,
    location_accuracy DOUBLE PRECISION,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_vendor_kyc_vendor_id ON vendor_kyc(vendor_id);

CREATE TABLE vendor_verifications (
    verification_id TEXT PRIMARY KEY,
    vendor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    trust_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    initial_score DOUBLE PRECISION NOT NULL DEFAULT 0,
    verdict TEXT NOT NULL,
    breakdown_json JSONB,
    flags JSONB,
    verdict_reason TEXT,
    validation_status TEXT NOT NULL DEFAULT 'pending',
    job_id TEXT,
    validated_by TEXT,
    validated_at TIMESTAMP,
    admin_note TEXT,
    admin_decision TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_vendor_verifications_vendor_id ON vendor_verifications(vendor_id);

CREATE TABLE verification_jobs (
    job_id TEXT PRIMARY KEY,
    vendor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'pending',
    current_step TEXT,
    steps_json JSONB,
    error_message TEXT,
    result TEXT,
    estimated_seconds INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_verification_jobs_vendor_id ON verification_jobs(vendor_id);

CREATE TABLE payment_transactions (
    transaction_id TEXT PRIMARY KEY,
    vendor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    transaction_ref TEXT NOT NULL,
    amount DECIMAL(13,2) NOT NULL,
    currency TEXT NOT NULL DEFAULT 'NGN',
    status TEXT NOT NULL DEFAULT 'pending',
    transaction_type TEXT NOT NULL,
    squad_checkout_url TEXT,
    squad_transaction_ref TEXT,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_payment_transactions_vendor_id ON payment_transactions(vendor_id);
CREATE INDEX idx_payment_transactions_transaction_ref ON payment_transactions(transaction_ref);

CREATE TABLE vendor_virtual_accounts (
    virtual_account_id TEXT PRIMARY KEY,
    vendor_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    virtual_account_number TEXT NOT NULL,
    bank_name TEXT NOT NULL,
    customer_identifier TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_vendor_virtual_accounts_vendor_id ON vendor_virtual_accounts(vendor_id);

CREATE TABLE admin_users (
    admin_id BIGINT PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    name TEXT NOT NULL,
    role TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL
);

CREATE INDEX idx_admin_users_email ON admin_users(email);