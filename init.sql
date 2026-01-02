-- Create database if it doesn't exist
SELECT 'CREATE DATABASE "CodeStreaks"'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'CodeStreaks')\gexec

-- Connect to the database
\c CodeStreaks

-- Create extensions if needed
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- The tables will be created automatically by GORM AutoMigrate