-- =============================================================
-- Atur Dana — Seed Data for Testing
-- User ID = 1 | ~500 transactions | Last 3 months
-- Date range: 2025-12-19 to 2026-03-19
-- Run: psql -U <user> -d <db> -f docs/seed.sql
-- =============================================================

-- Clear existing seed data (safe re-run)
DELETE FROM budgets      WHERE user_id = 1;
DELETE FROM transactions WHERE user_id = 1;
DELETE FROM categories   WHERE user_id = 1;
DELETE FROM users        WHERE id = 1;

-- Reset sequences to avoid gaps (optional)
-- SELECT setval('users_id_seq', MAX(id)) FROM users;

-- -------------------------------------------------------------
-- User
-- password = "password123" (bcrypt cost 10)
-- -------------------------------------------------------------
INSERT INTO users (id, created_at, updated_at, deleted_at, username, password_hash, email)
VALUES (
    1,
    '2025-12-01 08:00:00+07',
    '2025-12-01 08:00:00+07',
    NULL,
    'testuser',
    '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lh3y',
    'testuser@aturdana.dev'
);

-- -------------------------------------------------------------
-- Categories (user_id = 1)
-- Income: 1–5 | Expense: 6–15
-- -------------------------------------------------------------
INSERT INTO categories (id, created_at, updated_at, deleted_at, user_id, name, is_active) VALUES
-- Income
(1,  '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Salary',         TRUE),
(2,  '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Freelance',      TRUE),
(3,  '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Investment',     TRUE),
(4,  '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Bonus',          TRUE),
(5,  '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Side Business',  TRUE),
-- Expense
(6,  '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Food & Dining',  TRUE),
(7,  '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Transportation', TRUE),
(8,  '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Shopping',       TRUE),
(9,  '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Entertainment',  TRUE),
(10, '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Utilities',      TRUE),
(11, '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Healthcare',     TRUE),
(12, '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Groceries',      TRUE),
(13, '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Personal Care',  TRUE),
(14, '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Travel',         TRUE),
(15, '2025-12-01 08:00:00+07', '2025-12-01 08:00:00+07', NULL, 1, 'Education',      TRUE);

-- Bump sequences so future inserts don't collide
SELECT setval('users_id_seq',     (SELECT MAX(id) FROM users));
SELECT setval('categories_id_seq',(SELECT MAX(id) FROM categories));

-- -------------------------------------------------------------
-- Transactions (~500 rows) — generated via DO block
-- Distribution per day: ~5–6 transactions (mostly expense, some income)
-- Amount in IDR (Rupiah)
-- -------------------------------------------------------------
DO $$
DECLARE
    -- Lookup arrays (1-indexed)
    income_cats   INT[]   := ARRAY[1,2,3,4,5];
    expense_cats  INT[]   := ARRAY[6,7,8,9,10,11,12,13,14,15];

    income_descs  TEXT[]  := ARRAY[
        'Monthly salary deposit',
        'Salary transfer',
        'Mid-month advance',
        'Freelance project payment',
        'Web design contract',
        'Mobile app development',
        'Consulting fee',
        'Dividend income',
        'Stock return',
        'Mutual fund redemption',
        'Year-end bonus',
        'Performance bonus',
        'Holiday allowance',
        'Online store revenue',
        'Marketplace commission'
    ];

    expense_descs TEXT[]  := ARRAY[
        'Lunch at warung',
        'Dinner with family',
        'Coffee and snack',
        'Fast food order',
        'Restaurant dinner',
        'Online food delivery',
        'Grab/Gojek ride',
        'Monthly commuter pass',
        'Parking fee',
        'Fuel top-up',
        'Shirt purchase',
        'Online marketplace order',
        'Shoes and accessories',
        'Monthly utility bill',
        'Electricity bill',
        'Water bill',
        'Internet subscription',
        'Netflix subscription',
        'Spotify subscription',
        'Movie tickets',
        'Doctor visit',
        'Pharmacy / medicine',
        'Lab test fee',
        'Weekly groceries',
        'Supermarket shopping',
        'Fresh produce market',
        'Haircut',
        'Skincare product',
        'Flight ticket',
        'Hotel stay',
        'Online course',
        'Book purchase',
        'Workshop registration'
    ];

    -- Amount bands (IDR)
    small_amounts  FLOAT8[] := ARRAY[15000,20000,25000,30000,35000,40000,45000,50000];
    medium_amounts FLOAT8[] := ARRAY[75000,100000,125000,150000,175000,200000,250000,300000];
    large_amounts  FLOAT8[] := ARRAY[350000,400000,500000,600000,750000,1000000,1250000,1500000];
    income_amounts FLOAT8[] := ARRAY[500000,750000,1000000,1500000,2000000,3000000,5000000,8000000,10000000,15000000,20000000];

    cur_date   DATE;
    end_date   DATE    := '2026-03-19';
    trx_count  INT     := 0;
    target     INT     := 500;
    day_offset INT;
    i          INT;
    rand_val   FLOAT8;
    trx_type   TEXT;
    cat_id     INT;
    amount     FLOAT8;
    descr      TEXT;
    trx_time   TIMESTAMPTZ;
    hour_val   INT;
    minute_val INT;
BEGIN
    cur_date := '2025-12-19';

    WHILE cur_date <= end_date AND trx_count < target LOOP
        day_offset := EXTRACT(DOY FROM cur_date)::INT + EXTRACT(YEAR FROM cur_date)::INT * 366;

        -- Generate 4–7 transactions per day
        FOR i IN 1 .. (4 + (day_offset % 4)) LOOP
            EXIT WHEN trx_count >= target;

            rand_val   := (EXTRACT(EPOCH FROM NOW()) + trx_count * 17 + i * 31 + day_offset * 7)::BIGINT % 100 / 100.0;
            hour_val   := 7 + ((trx_count * 3 + i * 5 + day_offset) % 14);
            minute_val := (trx_count * 7 + i * 11) % 60;
            trx_time   := (cur_date::TIMESTAMP + (hour_val || ' hours')::INTERVAL + (minute_val || ' minutes')::INTERVAL)::TIMESTAMPTZ AT TIME ZONE 'Asia/Jakarta';

            -- ~20% income, ~80% expense
            -- Salary on 1st and 25th of each month (override)
            IF EXTRACT(DAY FROM cur_date) = 1 AND i = 1 THEN
                trx_type := 'income';
                cat_id   := 1; -- Salary
                amount   := income_amounts[1 + (EXTRACT(MONTH FROM cur_date)::INT % 3)]; -- 500k,750k,1M rotated; main salary much higher
                amount   := 8000000 + (EXTRACT(MONTH FROM cur_date)::INT % 3) * 500000;
                descr    := 'Monthly salary deposit';
            ELSIF EXTRACT(DAY FROM cur_date) = 25 AND i = 1 THEN
                trx_type := 'income';
                cat_id   := 1;
                amount   := 7500000;
                descr    := 'Salary transfer (advance)';
            ELSIF rand_val < 0.20 THEN
                trx_type := 'income';
                cat_id   := income_cats[1 + ((trx_count + i) % 5)];
                amount   := income_amounts[1 + ((trx_count * 3 + i) % 9)];
                descr    := income_descs[1 + ((trx_count + i * 2) % array_length(income_descs, 1))];
            ELSE
                trx_type := 'expense';
                -- Weight towards food, transport, groceries
                IF rand_val < 0.45 THEN
                    cat_id := expense_cats[1 + ((trx_count + i) % 3)]; -- food/transport/shopping
                    amount := small_amounts[1 + ((trx_count * 2 + i) % array_length(small_amounts, 1))];
                ELSIF rand_val < 0.70 THEN
                    cat_id := expense_cats[1 + ((trx_count + i) % 5)];
                    amount := medium_amounts[1 + ((trx_count + i * 3) % array_length(medium_amounts, 1))];
                ELSE
                    cat_id := expense_cats[1 + ((trx_count * 2 + i) % 10)];
                    amount := large_amounts[1 + ((trx_count + i) % array_length(large_amounts, 1))];
                END IF;
                descr := expense_descs[1 + ((trx_count * 2 + i * 3 + day_offset) % array_length(expense_descs, 1))];
            END IF;

            INSERT INTO transactions (created_at, updated_at, deleted_at, user_id, type, amount, description, category_id, date)
            VALUES (trx_time, trx_time, NULL, 1, trx_type, amount, descr, cat_id, trx_time);

            trx_count := trx_count + 1;
        END LOOP;

        cur_date := cur_date + INTERVAL '1 day';
    END LOOP;

    RAISE NOTICE 'Inserted % transactions', trx_count;
END;
$$;

SELECT setval('transactions_id_seq', (SELECT MAX(id) FROM transactions));

-- -------------------------------------------------------------
-- Budgets — monthly per expense category (3 months)
-- Dec 2025, Jan 2026, Feb 2026, Mar 2026
-- -------------------------------------------------------------
INSERT INTO budgets (created_at, updated_at, deleted_at, user_id, category_id, amount, start_date, end_date)
VALUES
-- Dec 2025
('2025-12-01 00:00:00+07','2025-12-01 00:00:00+07',NULL,1, 6,  1500000,'2025-12-01','2025-12-31'),
('2025-12-01 00:00:00+07','2025-12-01 00:00:00+07',NULL,1, 7,   800000,'2025-12-01','2025-12-31'),
('2025-12-01 00:00:00+07','2025-12-01 00:00:00+07',NULL,1, 8,  1000000,'2025-12-01','2025-12-31'),
('2025-12-01 00:00:00+07','2025-12-01 00:00:00+07',NULL,1, 9,   500000,'2025-12-01','2025-12-31'),
('2025-12-01 00:00:00+07','2025-12-01 00:00:00+07',NULL,1,10,   600000,'2025-12-01','2025-12-31'),
('2025-12-01 00:00:00+07','2025-12-01 00:00:00+07',NULL,1,11,   400000,'2025-12-01','2025-12-31'),
('2025-12-01 00:00:00+07','2025-12-01 00:00:00+07',NULL,1,12,  1200000,'2025-12-01','2025-12-31'),
('2025-12-01 00:00:00+07','2025-12-01 00:00:00+07',NULL,1,13,   300000,'2025-12-01','2025-12-31'),
('2025-12-01 00:00:00+07','2025-12-01 00:00:00+07',NULL,1,14,  2000000,'2025-12-01','2025-12-31'),
('2025-12-01 00:00:00+07','2025-12-01 00:00:00+07',NULL,1,15,   500000,'2025-12-01','2025-12-31'),
-- Jan 2026
('2026-01-01 00:00:00+07','2026-01-01 00:00:00+07',NULL,1, 6,  1500000,'2026-01-01','2026-01-31'),
('2026-01-01 00:00:00+07','2026-01-01 00:00:00+07',NULL,1, 7,   800000,'2026-01-01','2026-01-31'),
('2026-01-01 00:00:00+07','2026-01-01 00:00:00+07',NULL,1, 8,  1000000,'2026-01-01','2026-01-31'),
('2026-01-01 00:00:00+07','2026-01-01 00:00:00+07',NULL,1, 9,   500000,'2026-01-01','2026-01-31'),
('2026-01-01 00:00:00+07','2026-01-01 00:00:00+07',NULL,1,10,   600000,'2026-01-01','2026-01-31'),
('2026-01-01 00:00:00+07','2026-01-01 00:00:00+07',NULL,1,11,   400000,'2026-01-01','2026-01-31'),
('2026-01-01 00:00:00+07','2026-01-01 00:00:00+07',NULL,1,12,  1200000,'2026-01-01','2026-01-31'),
('2026-01-01 00:00:00+07','2026-01-01 00:00:00+07',NULL,1,13,   300000,'2026-01-01','2026-01-31'),
('2026-01-01 00:00:00+07','2026-01-01 00:00:00+07',NULL,1,14,  2000000,'2026-01-01','2026-01-31'),
('2026-01-01 00:00:00+07','2026-01-01 00:00:00+07',NULL,1,15,   500000,'2026-01-01','2026-01-31'),
-- Feb 2026
('2026-02-01 00:00:00+07','2026-02-01 00:00:00+07',NULL,1, 6,  1500000,'2026-02-01','2026-02-28'),
('2026-02-01 00:00:00+07','2026-02-01 00:00:00+07',NULL,1, 7,   800000,'2026-02-01','2026-02-28'),
('2026-02-01 00:00:00+07','2026-02-01 00:00:00+07',NULL,1, 8,  1000000,'2026-02-01','2026-02-28'),
('2026-02-01 00:00:00+07','2026-02-01 00:00:00+07',NULL,1, 9,   500000,'2026-02-01','2026-02-28'),
('2026-02-01 00:00:00+07','2026-02-01 00:00:00+07',NULL,1,10,   600000,'2026-02-01','2026-02-28'),
('2026-02-01 00:00:00+07','2026-02-01 00:00:00+07',NULL,1,11,   400000,'2026-02-01','2026-02-28'),
('2026-02-01 00:00:00+07','2026-02-01 00:00:00+07',NULL,1,12,  1200000,'2026-02-01','2026-02-28'),
('2026-02-01 00:00:00+07','2026-02-01 00:00:00+07',NULL,1,13,   300000,'2026-02-01','2026-02-28'),
('2026-02-01 00:00:00+07','2026-02-01 00:00:00+07',NULL,1,14,  2000000,'2026-02-01','2026-02-28'),
('2026-02-01 00:00:00+07','2026-02-01 00:00:00+07',NULL,1,15,   500000,'2026-02-01','2026-02-28'),
-- Mar 2026
('2026-03-01 00:00:00+07','2026-03-01 00:00:00+07',NULL,1, 6,  1500000,'2026-03-01','2026-03-31'),
('2026-03-01 00:00:00+07','2026-03-01 00:00:00+07',NULL,1, 7,   800000,'2026-03-01','2026-03-31'),
('2026-03-01 00:00:00+07','2026-03-01 00:00:00+07',NULL,1, 8,  1000000,'2026-03-01','2026-03-31'),
('2026-03-01 00:00:00+07','2026-03-01 00:00:00+07',NULL,1, 9,   500000,'2026-03-01','2026-03-31'),
('2026-03-01 00:00:00+07','2026-03-01 00:00:00+07',NULL,1,10,   600000,'2026-03-01','2026-03-31'),
('2026-03-01 00:00:00+07','2026-03-01 00:00:00+07',NULL,1,11,   400000,'2026-03-01','2026-03-31'),
('2026-03-01 00:00:00+07','2026-03-01 00:00:00+07',NULL,1,12,  1200000,'2026-03-01','2026-03-31'),
('2026-03-01 00:00:00+07','2026-03-01 00:00:00+07',NULL,1,13,   300000,'2026-03-01','2026-03-31'),
('2026-03-01 00:00:00+07','2026-03-01 00:00:00+07',NULL,1,14,  2000000,'2026-03-01','2026-03-31'),
('2026-03-01 00:00:00+07','2026-03-01 00:00:00+07',NULL,1,15,   500000,'2026-03-01','2026-03-31');

SELECT setval('budgets_id_seq', (SELECT MAX(id) FROM budgets));

-- -------------------------------------------------------------
-- Quick verification
-- -------------------------------------------------------------
SELECT type, COUNT(*) AS count, SUM(amount) AS total
FROM   transactions
WHERE  user_id = 1
GROUP  BY type;
