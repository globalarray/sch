SELECT qp.user_id, qp.quiz_id, qp.question_id, qp.answer, qp.correct
FROM quiz_progress qp
WHERE qp.quiz_id = ?