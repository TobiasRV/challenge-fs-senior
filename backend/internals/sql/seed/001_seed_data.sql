-- +goose Up

-- Seed Users (without team_id first to avoid FK violation)
INSERT INTO users (id, created_at, updated_at, username, password, email, role, team_id) VALUES
    ('127b93d2-b5e4-4a97-beed-e69442e34183', '2025-12-14 00:03:08.110657', '2025-12-14 00:03:08.110657', 'admin', '$2a$10$9mP65HDoxxiXR0lAmGHsMOvIvePIiEJWvDfvvRsXMAxL6UZwCRq2W', 'admin@admin.com', 'Admin', NULL),
    ('0d5e7086-fcab-4e45-9669-d5553901e837', '2025-12-14 00:04:46.764337', '2025-12-14 00:04:46.764337', 'Manager 1', '$2a$10$qVtwmIoATkjoLde2X2qku.qG7c50Pn5BC3tsWslSQP0Je0kq8x6LO', 'manager@manager.com', 'Manager', NULL),
    ('8bae7c97-c6bf-4245-b150-77a271804533', '2025-12-14 00:04:58.894826', '2025-12-14 00:04:58.894826', 'Manager 2', '$2a$10$h1Op3bn1y6V9vKOB4D8OFen6PvowQznRi8c1O3K5X19K9EGyvfY3a', 'manager2@manager.com', 'Manager', NULL),
    ('efdb7586-7b0f-4ba5-b8f2-96fe387b57c3', '2025-12-14 00:05:23.263517', '2025-12-14 00:05:23.263517', 'Miembro 1', '$2a$10$9uPAxzShNB88SJAdksVvhud3XKMysb4ss3qNrAiT2skmpmfriG0o2', 'member@member.com', 'Member', NULL),
    ('631b0765-dce6-46da-b0a6-d945f3ed6fed', '2025-12-14 00:05:41.464407', '2025-12-14 00:05:41.464407', 'Miembro 2', '$2a$10$trcIioRKDy/6OV0TiHhj6.pwRAru0.4Qzt7qCQZintID/QbFIwlN.', 'member2@member.com', 'Member', NULL);

-- Seed Teams
INSERT INTO teams (id, created_at, updated_at, name, owner_id) VALUES
    ('1df07229-93c4-4c5a-83f1-719878772f9a', '2025-12-14 00:03:15.187194', '2025-12-14 00:03:15.187194', 'Team 1', '127b93d2-b5e4-4a97-beed-e69442e34183');

-- Update Users with team_id
UPDATE users SET team_id = '1df07229-93c4-4c5a-83f1-719878772f9a' WHERE id IN (
    '0d5e7086-fcab-4e45-9669-d5553901e837',
    '8bae7c97-c6bf-4245-b150-77a271804533',
    'efdb7586-7b0f-4ba5-b8f2-96fe387b57c3',
    '631b0765-dce6-46da-b0a6-d945f3ed6fed'
);

-- Seed Projects
INSERT INTO projects (id, created_at, updated_at, name, team_id, manager_id, status) VALUES
    ('34de3c56-faf4-45a3-8e8f-608f6e4cf3e0', '2025-12-14 00:06:45.019592', '2025-12-14 00:06:45.019592', 'Projecto 1 Manager 1', '1df07229-93c4-4c5a-83f1-719878772f9a', '0d5e7086-fcab-4e45-9669-d5553901e837', 'OnHold'),
    ('25718bf8-10c0-4c9c-96ea-085fdf7b671c', '2025-12-14 00:06:52.38892', '2025-12-14 00:06:52.38892', 'Projecto 2 Manager 1', '1df07229-93c4-4c5a-83f1-719878772f9a', '0d5e7086-fcab-4e45-9669-d5553901e837', 'OnHold'),
    ('4a9a691e-92d3-4a81-8c33-166b7b9d4caf', '2025-12-14 00:09:21.646134', '2025-12-14 00:09:21.646134', 'Proyecto 2 Manager 2', '1df07229-93c4-4c5a-83f1-719878772f9a', '8bae7c97-c6bf-4245-b150-77a271804533', 'OnHold'),
    ('0fce547a-5a4e-4590-8a86-4ab61a554300', '2025-12-14 00:08:39.24456', '2025-12-14 00:11:43.952112', 'Proyecto 1 Manager 2', '1df07229-93c4-4c5a-83f1-719878772f9a', '8bae7c97-c6bf-4245-b150-77a271804533', 'OnHold');

-- Seed Tasks
INSERT INTO tasks (id, created_at, updated_at, project_id, user_id, status, title, description) VALUES
    ('064206e2-7a95-4dc2-bb65-3745534249da', '2025-12-14 00:07:27.774294', '2025-12-14 00:07:27.774294', '34de3c56-faf4-45a3-8e8f-608f6e4cf3e0', 'efdb7586-7b0f-4ba5-b8f2-96fe387b57c3', 'ToDo', 'Tarea 1 Projecto 1 Manager 1', 'Esto es una descripcion'),
    ('0d840c2c-ece3-4def-b9ce-09202caf98d8', '2025-12-14 00:07:44.413343', '2025-12-14 00:07:44.413343', '34de3c56-faf4-45a3-8e8f-608f6e4cf3e0', '631b0765-dce6-46da-b0a6-d945f3ed6fed', 'ToDo', 'Tarea 2 Projecto 1 Manager 1', NULL),
    ('9c1458d4-6476-47aa-bbea-c2b441cbf60c', '2025-12-14 00:08:04.307129', '2025-12-14 00:08:04.307129', '25718bf8-10c0-4c9c-96ea-085fdf7b671c', 'efdb7586-7b0f-4ba5-b8f2-96fe387b57c3', 'ToDo', 'Tarea 1 Projecto 2 Manager 1', NULL),
    ('360b31fc-b171-4094-9263-5a946d967140', '2025-12-14 00:08:16.514514', '2025-12-14 00:08:16.514514', '25718bf8-10c0-4c9c-96ea-085fdf7b671c', '631b0765-dce6-46da-b0a6-d945f3ed6fed', 'ToDo', 'Tarea 2 Projecto 2 Manager 1', NULL),
    ('84373ac3-5b54-42d5-8019-67a6d6d60b07', '2025-12-14 00:13:01.692269', '2025-12-14 00:13:01.69227', '0fce547a-5a4e-4590-8a86-4ab61a554300', 'efdb7586-7b0f-4ba5-b8f2-96fe387b57c3', 'ToDo', 'Tarea 1 Proyecto 1 Manager 2', NULL),
    ('b393750b-3d11-4d23-a573-4ffef5eeb0f3', '2025-12-14 00:13:17.125859', '2025-12-14 00:13:17.125859', '0fce547a-5a4e-4590-8a86-4ab61a554300', '631b0765-dce6-46da-b0a6-d945f3ed6fed', 'ToDo', 'Tarea 2 Proyecto 1 Manager 2', NULL),
    ('a35e8570-dbc4-45c1-a771-a5bb588a60ea', '2025-12-14 00:13:33.324957', '2025-12-14 00:13:33.324957', '4a9a691e-92d3-4a81-8c33-166b7b9d4caf', 'efdb7586-7b0f-4ba5-b8f2-96fe387b57c3', 'ToDo', 'Tarea 1 Proyecto 2 Manager 2', NULL),
    ('678436d4-ec82-42ab-a4c0-25ce49d13a39', '2025-12-14 00:13:46.916076', '2025-12-14 00:13:46.916076', '4a9a691e-92d3-4a81-8c33-166b7b9d4caf', '631b0765-dce6-46da-b0a6-d945f3ed6fed', 'ToDo', 'Tarea 2 Proyecto 2 Manager 2', NULL);

-- +goose Down
DELETE FROM tasks;

DELETE FROM projects;

UPDATE users SET team_id = NULL;

DELETE FROM teams;

DELETE FROM users;
