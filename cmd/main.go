package main

import (
	"fmt"
	"perplexity"
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

	msg := "write helloworld in golang"
	if err = s.Ask(msg); err != nil {
		panic(err)
	}
	defer s.Wss.Close()

	answer, err := s.ReadAnswer()
	if err != nil {
		panic(err)
	}
	fmt.Println(answer.Answer)
}
