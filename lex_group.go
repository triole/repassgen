package main

import "fmt"

type groupGenerator struct {
	pattern string
	entropy *float64
}

func (g *groupGenerator) Generate(s *State) error {
	out, err := Generate(GenerateInput{
		Pattern: g.pattern,
	})
	if err != nil {
		lexErr, ok := err.(*LexError)
		if ok {
			lexErr.MovePos(int(s.patternBuffStart))
			return lexErr
		}
		return s.errorUnknown(err.Error())
	}
	err = s.addOutputNonRepeatable(out.Password)
	if err != nil {
		return err
	}
	s.patternEntropy += out.PatternEntropy
	g.entropy = &out.PatternEntropy
	return nil
}

func (g *groupGenerator) Level() int {
	return 0
}

func (g *groupGenerator) Entropy() (float64, error) {
	if g.entropy != nil {
		return *g.entropy, nil
	}
	return 0, fmt.Errorf("entropy is not calculated")
}

func lexGroup(s *State) (LexType, error) {
	if s.end() {
		return nil, s.errorSyntax("'(' not closed")
	}
	c := s.pattern[s.patternPos]
	s.patternPos++
	n := uint(len(s.patternBuff))
	switch c {
	case '\\':
		return lexGroupBackslash, nil
	case ')':
		childPattern := string(s.patternBuff[s.patternBuffStart:n])
		gen := &groupGenerator{pattern: childPattern}
		err := gen.Generate(s)
		if err != nil {
			return nil, err
		}
		s.lastGen = gen
		s.patternBuff = nil
		return LexRoot, nil
	}
	s.patternBuff = append(s.patternBuff, c)
	return lexGroup, nil
}

func lexGroupBackslash(s *State) (LexType, error) {
	c := s.pattern[s.patternPos]
	s.patternPos++
	if c == ')' {
		s.patternBuff = append(s.patternBuff, c)
	} else {
		s.patternBuff = append(s.patternBuff, '\\', c)
	}
	return lexGroup, nil
}
