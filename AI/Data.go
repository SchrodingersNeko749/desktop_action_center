package AI

type Prompt struct {
	MaxTokens            int     `json:"maxTokens"`
	Input                string  `json:"input"`
	Model                string  `json:"model"`
	IgnoreEos            bool    `json:"ignore_eos"`
	TopK                 int     `json:"top_k"`
	TopP                 float64 `json:"top_p"`
	Temperature          float64 `json:"temperature"`
	Mirostat             int     `json:"mirostat"`
	Entropy              int     `json:"entropy"`
	LearningRate         float64 `json:"learningRate"`
	TailFreeSamplingRate int     `json:"tailFreeSamplingRate"`
	TypicalP             int     `json:"typical_p"`
	PenalizeNewLines     bool    `json:"penalizeNewLines"`
	PenalizeSpaces       bool    `json:"penalizeSpaces"`
	RepetitionPenalty    float64 `json:"repetition_penalty"`
	IncludeIngest        bool    `json:"includeIngest"`
	IncludeStatistics    bool    `json:"includeStatistics"`
	OneShot              bool    `json:"oneShot"`
}
