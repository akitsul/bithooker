package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrepare_DoesNotModifiesSingleLineData(t *testing.T) {
	test := assert.New(t)

	origin := map[string]string{
		"key": "value",
	}

	actual := prepare(origin)
	test.Len(actual, 1)
	test.Equal("value", actual["key"])
}

func TestPrepare_RemovesNewLinePrefix(t *testing.T) {
	test := assert.New(t)

	origin := map[string]string{
		"key": "\nvalue",
	}

	actual := prepare(origin)
	test.Len(actual, 1)
	test.Equal("value", actual["key"])
}

func TestPrepare_ReplacesDoubleNewLineSuffixToNewLine(t *testing.T) {
	test := assert.New(t)

	origin := map[string]string{
		"key": "value\n\n",
	}

	actual := prepare(origin)
	test.Len(actual, 1)
	test.Equal("value\n", actual["key"])
}

func TestPrepare_RemovesNewLinePrefixAndReplacesDoubleNewLineSuffix(
	t *testing.T,
) {
	test := assert.New(t)

	origin := map[string]string{
		"key": "\nvalue\n\n",
	}

	actual := prepare(origin)
	test.Len(actual, 1)
	test.Equal("value\n", actual["key"])
}
