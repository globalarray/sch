SELECT q.id, q.quiz_id, q.question, q.answers
FROM question q
WHERE q.id = ?