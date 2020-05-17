INSERT INTO app.todo_status (status) VALUES
    ('new'),
    ('progress'),
    ('completed'),
    ('archived');

-- seed the todo table with fake todos
INSERT INTO app.todo (title, description, status_id, updated_dt) VALUES
    ('task 1','description 1', 1, NULL),
    ('task 2','description 2', 2, NOW()),
    ('task 3','description 3', 3, NOW()),
    ('task 4','description 4', 4, NULL),
    ('task 5','description 5', 4, NULL),
    ('task 6','description 6', 3, NOW()),
    ('task 7','description 7', 1, NULL),
    ('task 8','description 8', 2, NOW()),
    ('task 9','description 9', 3, NOW()),
    ('task 10','description 10', 2, NOW());

COMMIT;
