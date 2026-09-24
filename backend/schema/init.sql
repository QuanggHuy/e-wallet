CREATE TABLE IF NOT EXISTS accounts (
    id BIGSERIAL PRIMARY KEY,
    account_no VARCHAR(30) UNIQUE NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    balance NUMERIC(18,2) NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Bảng lịch sử giao dịch, partition theo tháng
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL,
    tx_code VARCHAR(30) NOT NULL,
    from_account_id BIGINT NOT NULL,
    to_account_id BIGINT NOT NULL,
    amount NUMERIC(18,2) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    note TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    PRIMARY KEY (id, created_at)
) PARTITION BY RANGE (created_at);

-- Giải thích nhanh từng cột:
-- - from_account_id — ID tài khoản gửi tiền
-- - to_account_id — ID tài khoản nhận tiền
-- - amount — số tiền chuyển
-- - status — trạng thái: pending / completed / failed
-- - note — ghi chú chuyển tiền (có thể để trống nên không có NOT NULL)
-- - PRIMARY KEY (id, created_at) — khi dùng partition, primary key bắt buộc phải bao gồm cột partition (created_at)
-- - PARTITION BY RANGE (created_at) — khai báo chia theo khoảng thời gian

CREATE TABLE IF NOT EXISTS transactions_2026_09
    PARTITION OF transactions
    FOR VALUES FROM ('2026-09-01') TO ('2026-10-01');

CREATE TABLE IF NOT EXISTS transactions_2026_10
    PARTITION OF transactions
    FOR VALUES FROM ('2026-10-01') TO ('2026-11-01');

CREATE TABLE IF NOT EXISTS transactions_2026_11
    PARTITION OF transactions
    FOR VALUES FROM ('2026-11-01') TO ('2026-12-01');

CREATE TABLE IF NOT EXISTS transactions_2026_12
    PARTITION OF transactions
    FOR VALUES FROM ('2026-12-01') TO ('2027-01-01');

-- Index cho accounts
CREATE INDEX IF NOT EXISTS idx_accounts_account_no
    ON accounts(account_no);

-- Index cho transactions
CREATE INDEX IF NOT EXISTS idx_transactions_from
    ON transactions(from_account_id, created_at DESC);

CREATE INDEX IF NOT EXISTS idx_transactions_to
    ON transactions(to_account_id, created_at DESC);

-- Partial index: chỉ index các giao dịch đang pending
CREATE INDEX IF NOT EXISTS idx_transactions_pending
    ON transactions(status) WHERE status = 'pending';

CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password_hash VARCHAR(100) NOT NULL,
    account_id BIGINT UNIQUE NOT NULL REFERENCES accounts(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);