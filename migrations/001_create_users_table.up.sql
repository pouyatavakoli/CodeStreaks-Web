-- Create users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    handle VARCHAR(50) UNIQUE NOT NULL,
    currentStreak INTEGER DEFAULT 0,
    maxStreak INTEGER DEFAULT 0,
    lastSubmissionDate TIMESTAMP,
    lastUpdatedAt TIMESTAMP NOT NULL,
    createdAt TIMESTAMP NOT NULL DEFAULT NOW(),
    -- Check constraints
    CONSTRAINT chk_streaks_non_negative CHECK (
        currentStreak >= 0
        AND maxStreak >= 0
    ),
    CONSTRAINT chk_max_streak_gte_current CHECK (maxStreak >= currentStreak)
);