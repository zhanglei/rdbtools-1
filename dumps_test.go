package rdbtools

import (
	"os"
	"testing"
)

// This file contains tests for all the dumps in the dumps/ directory
// Those dumps are taken from https://github.com/sripathikrishnan/redis-rdb-tools

func mustOpen(t *testing.T, path string) *os.File {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("Error while opening file '%s'; err=%s", path, err)
	}

	return f
}

func doParse(t *testing.T, p *Parser, path string) {
	err := p.Parse(mustOpen(t, path))
	if err != nil {
		t.Fatalf("Error while parsing '%s'; err=%s", path, err)
	}
}

func TestDumpDictionary(t *testing.T) {
	p := NewParser(ParserContext{
		HashMetadataCh: make(chan HashMetadata),
		HashDataCh:     make(chan StringObject),
	})

	go doParse(t, p, "dumps/dictionary.rdb")

	var i, j int
	stop := false
	for !stop {
		select {
		case _, ok := <-p.ctx.HashMetadataCh:
			if !ok {
				p.ctx.HashMetadataCh = nil
				break
			}
			i++
		case _, ok := <-p.ctx.HashDataCh:
			if !ok {
				p.ctx.HashDataCh = nil
				break
			}
			j++
		}

		if p.ctx.Invalid() {
			break
		}
	}

	equals(t, 1, i)
	equals(t, 1000, j)
}