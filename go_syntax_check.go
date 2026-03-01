package main

type goSyntaxChecker struct {
	lastPath   string
	lastSource string
	lastLines  int
	lineErrors map[int]struct{}
	lineMsgs   map[int]string
}

func newGoSyntaxChecker() *goSyntaxChecker {
	return &goSyntaxChecker{}
}

func (c *goSyntaxChecker) lineErrorsFor(_ string, _ []rune) map[int]struct{} {
	if c == nil {
		return nil
	}
	c.lineErrors = nil
	c.lineMsgs = nil
	return nil
}
