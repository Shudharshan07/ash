package prompt

import "os"

type Prompt struct {
	// Text string
}

func NewPrompt() *Prompt {
	return &Prompt{}
}

func (p *Prompt) Text() string {
	pwd, err := os.Getwd() // convert this to []byte , mayebe that can be used for perofromace
	if err != nil {
		return "? >"
	}

	pwd += " >"

	return pwd
}

func (p *Prompt) Width() int {
	return len([]rune(p.Text()))
}
