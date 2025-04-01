SELECT s.key, s.name, s.patronymic, s.surname, s.creation, s.expiration, s.role, s.created_by
FROM secret s
WHERE s.key = ?