SELECT qr.quiz_id, qr.user_id, qr.score, qr.completed_in
FROM quiz_result qr
WHERE qr.quiz_id = ?
ORDER BY qr.completed_in DESC