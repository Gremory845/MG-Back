package leaderboard

type Entry struct {
	UserID       string  `json:"user_id"`
	Username     string  `json:"username"`
	BestScore    int     `json:"best_score"`
	BestAccuracy float64 `json:"best_accuracy"`
}
