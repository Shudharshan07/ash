package prompt

import "os"

type Prompt struct {
	Text string
}

func NewPrompt() *Prompt {
	pwd, err := os.Getwd() // convert this to []byte , mayebe that can be used for perofromace
	if err != nil {
		pwd = "$"
	}
	pwd += " >"
	prompt := &Prompt{
		Text: pwd,
	}

	return prompt
}

func (p *Prompt) Width() int {
	return len([]rune(p.Text))
}
