package quizz

import (
	"fmt"
	"io"
	"os"
)

type Cli struct {
	game Game
	out  io.Writer
	in   io.Reader
}

func NewCLI(gameLocation string) (*Cli, error) {
	game, err := NewQuizz(gameLocation)

	if err != nil {
		return nil, fmt.Errorf("Error creating CLI for the game %v", err)
	}

	cli := Cli{game, os.Stdout, os.Stdin}

	return &cli, nil
}

const loadingPrompt = "Loading default game..."

func (c *Cli) Play() {
	fmt.Fprintln(c.out, loadingPrompt)

	c.game.Start(c.out, c.in)

	c.game.Finish(c.out)
}