SELECT qr.user_id, qr.quiz_id, qr.score, qr.completed_in
FROM quiz_result qr
WHERE qr.quiz_id = ? and qr.user_id = ?