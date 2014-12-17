package sandbox

import (
	"fmt"
	"strconv"
)

type Instruction struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Tutorial struct {
	Instructions []Instruction `json:"instructions"`
}

func NewTutorial(instructions []map[string]string) Tutorial {
	tutorial := Tutorial{}

	for _, inst := range instructions {
		i := Instruction{Title: inst["title"], Body: inst["body"]}
		tutorial.AddInstruction(i)
	}

	return tutorial
}

func (t *Tutorial) AddInstruction(instruction Instruction) {
	inst := append(t.Instructions, instruction)
	t.Instructions = inst
}

func (t *Tutorial) ToCommands() string {
	var commands string

	for index, inst := range t.Instructions {
		stepNum := strconv.Itoa(index + 1)
		//echo "1: What is this?"
		l1 := fmt.Sprintf("echo \"%s: %s\"", stepNum, inst.Title)
		l2 := fmt.Sprintf("echo \"%s\"", inst.Body)

		//commands = append(commands, l1, l2)
		commands += l1 + "\n" + l2 + "\n\n"
	}

	return commands
}
