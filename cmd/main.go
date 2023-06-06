package main

import (
	"fmt"

	"github.com/binarycraft007/perplexity"
)

func main() {
	s, err := perplexity.NewSession()
	if err != nil {
		panic(err)
	}

	if err = s.Check(); err != nil {
		panic(err)
	}

	if err = s.GetSid(); err != nil {
		panic(err)
	}

	if err = s.InitWss(); err != nil {
		panic(err)
	}
	defer s.Close()

	msg := "write helloworld in golang"
	if err = s.Ask(msg); err != nil {
		panic(err)
	}

	answer, err := s.ReadAnswer()
	if err != nil {
		panic(err)
	}
	fmt.Println(answer.Answer)
}
