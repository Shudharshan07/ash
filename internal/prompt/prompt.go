package prompt

type Prompt struct {
	Text string
}

func NewPrompt() *Prompt {
	prompt := &Prompt{
		Text: "$ ",
	}

	return prompt
}

func (p *Prompt) Width() int {
	return len([]rune(p.Text))
}
