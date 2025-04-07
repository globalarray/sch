SELECT q.id, q.name, q.creation, q.created_by
FROM quiz q
WHERE q.created_by = ?
ORDER BY q.creation