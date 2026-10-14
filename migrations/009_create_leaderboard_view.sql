CREATE VIEW leaderboard AS
SELECT
    user_id,
    level_id,
    MAX(score) AS best_score,
    MAX(accuracy) AS best_accuracy
FROM game_attempts
WHERE completed = TRUE
GROUP BY user_id, level_id;
