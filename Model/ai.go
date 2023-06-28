package Model

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

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

func RunInference(p Prompt) <-chan string {
	out := make(chan string)
	go func() {
		payload, _ := json.Marshal(p)
		resp, err := http.Post("http://llama.her.st/completion", "application/json", bytes.NewBuffer(payload))
		if err != nil {
			fmt.Println("Something went wrong with the completion", err)
			close(out)
			return
		}
		defer resp.Body.Close()
		reader := bufio.NewReader(resp.Body)
		for {
			buffer := make([]byte, 2)
			for i := 0; i < 2; i++ {
				b, err := reader.ReadByte()
				if err != nil {
					close(out)
					return
				}
				buffer[i] = b
			}
			out <- string(buffer)
		}
	}()
	return out
}

func GeneratePrompt(mode string, input string, maxToken int, model string, inclIngest bool, inclStats bool) Prompt {
	builder := strings.Builder{}

	switch mode {
	case "completion":
		builder.WriteString(input)
	case "instruction":
		builder.WriteString("### Instruction: ")
		builder.WriteRune('\n')
		builder.WriteString(input)
		builder.WriteRune('\n')
		builder.WriteString("### Response:")
		builder.WriteRune('\n')
	case "chat":
		builder.WriteString("User: ")
		builder.WriteString(input)
		builder.WriteRune('\n')
		builder.WriteString("AI:")
	}

	p := Prompt{
		MaxTokens:            maxToken,
		Input:                builder.String(),
		Model:                model,
		IgnoreEos:            false,
		TopK:                 20,
		TopP:                 0.9,
		Temperature:          0.1,
		Mirostat:             2,
		Entropy:              3,
		LearningRate:         0.003,
		TailFreeSamplingRate: 1,
		TypicalP:             1,
		PenalizeNewLines:     false,
		PenalizeSpaces:       false,
		RepetitionPenalty:    1.3,
		IncludeIngest:        inclIngest,
		IncludeStatistics:    inclStats,
		OneShot:              true,
	}
	return p
}

func GetModels() []string {
	models := []string{}
	resp, err := http.Get("http://llama.her.st/models")
	if err != nil {
		fmt.Println("Failed to get list of models.", err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan(); i++ {
		m := scanner.Text()
		models = append(models, m)
	}
	return models
}
