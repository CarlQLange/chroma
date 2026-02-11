package formatters

import (
	"strings"
	"testing"

	assert "github.com/alecthomas/assert/v2"
	"github.com/alecthomas/chroma/v2"
)

func TestClosestColour(t *testing.T) {
	actual := findClosest(ttyTables[256], chroma.MustParseColour("#e06c75"))
	assert.Equal(t, chroma.MustParseColour("#d75f87"), actual)
}

func TestNoneColour(t *testing.T) {
	formatter := TTY256
	tokenType := chroma.None

	style, err := chroma.NewStyle("test", chroma.StyleEntries{
		chroma.Background: "#D0ab1e",
	})
	assert.NoError(t, err)

	stringBuilder := strings.Builder{}
	err = formatter.Format(&stringBuilder, style, chroma.Literator(chroma.Token{
		Type:  tokenType,
		Value: "WORD",
	}))
	assert.NoError(t, err)

	// "178" = #d7af00 approximates #d0ab1e
	//
	// 178 color ref: https://jonasjacek.github.io/colors/
	assert.Equal(t, "\033[38;5;178mWORD\033[0m", stringBuilder.String())
}

func TestTermColourTrueColourFormatter(t *testing.T) {
	style, err := chroma.NewStyle("test", chroma.StyleEntries{
		chroma.Keyword: "term-4 bold",
	})
	assert.NoError(t, err)

	var buf strings.Builder
	err = TTY16m.Format(&buf, style, chroma.Literator(chroma.Token{
		Type:  chroma.Keyword,
		Value: "if",
	}))
	assert.NoError(t, err)

	// Should use indexed escape code \033[38;5;4m not truecolour RGB
	assert.Contains(t, buf.String(), "\033[38;5;4m")
	assert.NotContains(t, buf.String(), "\033[38;2;")
}

func TestTermColourTrueColourFormatterBackground(t *testing.T) {
	style, err := chroma.NewStyle("test", chroma.StyleEntries{
		chroma.Keyword: "term-2 bg:term-0",
	})
	assert.NoError(t, err)

	var buf strings.Builder
	err = TTY16m.Format(&buf, style, chroma.Literator(chroma.Token{
		Type:  chroma.Keyword,
		Value: "if",
	}))
	assert.NoError(t, err)

	assert.Contains(t, buf.String(), "\033[38;5;2m")
	assert.Contains(t, buf.String(), "\033[48;5;0m")
}

func TestTermColour256Formatter(t *testing.T) {
	style, err := chroma.NewStyle("test", chroma.StyleEntries{
		chroma.Keyword: "term-4",
	})
	assert.NoError(t, err)

	var buf strings.Builder
	err = TTY256.Format(&buf, style, chroma.Literator(chroma.Token{
		Type:  chroma.Keyword,
		Value: "if",
	}))
	assert.NoError(t, err)

	// Index 4 should emit \033[34m (standard colour)
	assert.Contains(t, buf.String(), "\033[34m")
}

func TestTermColour16Formatter(t *testing.T) {
	style, err := chroma.NewStyle("test", chroma.StyleEntries{
		chroma.Keyword: "term-4",
	})
	assert.NoError(t, err)

	var buf strings.Builder
	err = TTY16.Format(&buf, style, chroma.Literator(chroma.Token{
		Type:  chroma.Keyword,
		Value: "if",
	}))
	assert.NoError(t, err)

	// Index 4 should use \033[34m for foreground
	assert.Contains(t, buf.String(), "\033[34m")
}

func TestTermColour16FormatterBright(t *testing.T) {
	style, err := chroma.NewStyle("test", chroma.StyleEntries{
		chroma.Keyword: "term-12",
	})
	assert.NoError(t, err)

	var buf strings.Builder
	err = TTY16.Format(&buf, style, chroma.Literator(chroma.Token{
		Type:  chroma.Keyword,
		Value: "if",
	}))
	assert.NoError(t, err)

	// Index 12 (bright blue) uses \033[94m
	assert.Contains(t, buf.String(), "\033[94m")
}

func TestTermColourHighIndex(t *testing.T) {
	style, err := chroma.NewStyle("test", chroma.StyleEntries{
		chroma.Keyword: "term-200",
	})
	assert.NoError(t, err)

	var buf strings.Builder
	err = TTY256.Format(&buf, style, chroma.Literator(chroma.Token{
		Type:  chroma.Keyword,
		Value: "if",
	}))
	assert.NoError(t, err)

	// Index >= 16 should use \033[38;5;Nm
	assert.Contains(t, buf.String(), "\033[38;5;200m")
}
