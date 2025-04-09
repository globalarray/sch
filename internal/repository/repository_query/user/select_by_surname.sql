SELECT u.tg_id, u.name, u.patronymic, u.surname, u.role
FROM users u
WHERE u.surname = ?