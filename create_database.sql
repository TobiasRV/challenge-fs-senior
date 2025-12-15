SELECT 'CREATE DATABASE tasks-manager' 
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'tasks-manager')