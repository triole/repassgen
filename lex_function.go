package main

import "fmt"

type encoderFunctionCallGenerator struct {
	funcName   string
	argPattern string
	entropy    *float64
}

func baseFunctionCallGenerator(
	s *State,
	funcName string,
	funcObj func(in []rune) ([]rune, error),
	argPattern string,
) (*GenerateOutput, error) {
	argOut, err := Generate(GenerateInput{
		Pattern: argPattern,
	})
	if err != nil {
		lexErr, ok := err.(*LexError)
		if ok {
			lexErr.MovePos(int(s.patternBuffStart + 1))
			return nil, lexErr
		}
		return nil, s.errorUnknown(err.Error())
	}
	result, err := funcObj(argOut.Password)
	if err != nil {
		lexErr, ok := err.(*LexError)
		if ok {
			lexErr.MovePos(int(s.patternBuffStart))
			lexErr.PrependMsg("function " + funcName)
			return nil, lexErr
		}
		return nil, s.errorUnknown("%v returned error: %v", funcName, err)
	}
	err = s.addOutputNonRepeatable(result)
	if err != nil {
		return nil, err
	}
	s.patternEntropy += argOut.PatternEntropy
	return argOut, nil
}

func (g *encoderFunctionCallGenerator) Generate(s *State) error {
	funcName := g.funcName
	funcObj, ok := encoderFunctions[funcName]
	if !ok {
		return s.errorValue("invalid function '%v'", funcName)
	}
	argOut, err := baseFunctionCallGenerator(s, funcName, funcObj, g.argPattern)
	if err != nil {
		return err
	}
	g.entropy = &argOut.PatternEntropy
	return nil
}

func (g *encoderFunctionCallGenerator) Level() int {
	return 0
}

func (g *encoderFunctionCallGenerator) Entropy() (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, fmt.Errorf("entropy is not calculated")
}

func getFuncGenerator(s *State, funcName string, arg string) (generatorIface, error) {
	if _, ok := encoderFunctions[funcName]; ok {
		return &encoderFunctionCallGenerator{
			funcName:   funcName,
			argPattern: arg,
		}, nil
	}
	switch funcName {
	case "bip39word":
		return newBIP99WordGenerator(arg)
	}
	return nil, s.errorValue("invalid function '%v'", funcName)
}

func lexIdentFuncCall(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'(' not closed")
	}
	n := uint(len(s.patternBuff))
	// "$a()"  -->  c.patternBuffStart == 1
	c := s.pattern[s.patternPos]
	s.patternPos++
	switch c {
	case '(':
		s.openParenth++
	case ')':
		s.openParenth--
		if s.openParenth > 0 {
			break
		}
		funcName := string(s.patternBuff[:s.patternBuffStart])
		if funcName == "" {
			return nil, s.errorSyntax("missing function name")
		}
		arg := string(s.patternBuff[s.patternBuffStart:n])
		gen, err := getFuncGenerator(s, funcName, arg)
		if err != nil {
			return nil, err
		}
		err = gen.Generate(s)
		if err != nil {
			return nil, err
		}
		s.patternBuff = nil
		s.lastGen = gen
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexIdentFuncCall, nil
}
