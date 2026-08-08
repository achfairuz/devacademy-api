package quiz

type QuizOptionRequest struct {
	OptionText string `json:"option_text" binding:"required"`
	IsCorrect  bool   `json:"is_correct"`
}

type QuizQuestionRequest struct {
	Question     string              `json:"question" binding:"required"`
	QuestionType string              `json:"question_type" binding:"omitempty,oneof=multiple_choice true_false"`
	Options      []QuizOptionRequest `json:"options"`
}

type QuizRequest struct {
	Title        string                `json:"title" binding:"required"`
	PassingScore int                   `json:"passing_score"`
	Questions    []QuizQuestionRequest `json:"questions"`
}
