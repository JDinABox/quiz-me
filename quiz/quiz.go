package quiz

type QuizID string

type Quiz struct {
	ID                  string     `json:"id"`
	Name                string     `json:"name"`
	Description         string     `json:"description"`
	MaxQuestionsPerQuiz int        `json:"maxQuestionsPerQuiz"`
	Questions           []Question `json:"questions"`
}

type Question struct {
	Question    string         `json:"question"`
	Answers     map[int]string `json:"answers"`
	Correct     []int          `json:"correct"`
	Explanation string         `json:"explanation"`
}

type Quizzes map[QuizID]Quiz
