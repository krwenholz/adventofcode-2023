package scanner

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/participle/v2"
)

type PuzzleScanner[G any] struct {
	parser  *participle.Parser[G]
	scanner *bufio.Scanner
}

func NewScanner[G any](parser *participle.Parser[G], path string) *PuzzleScanner[G] {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	return &PuzzleScanner[G]{
		parser:  parser,
		scanner: bufio.NewScanner(f),
	}
}

func (p *PuzzleScanner[G]) Scan() bool {
	return p.scanner.Scan()
}

func (p *PuzzleScanner[G]) Struct() *G {
	g, err := p.parser.ParseBytes("", p.scanner.Bytes())
	if err != nil {
		fmt.Println(p.scanner.Text())
		log.Fatal(g, err)
	}

	return g
}
