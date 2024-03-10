package text

import (
	"testing"
)

func Test_ExpandText_NoTokens(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
	}

	value, err := ExpandText("hello, world", config)
	if err != nil {
		t.Fatal("expansion failed")
	}

	if value != "hello, world" {
		t.Fatal("should be hello, world")
	}

}

func Test_ExpandText_OneTokens(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
	}

	value, err := ExpandText("hello, ${os:pid} world", config)
	if err != nil {
		t.Fatal("expansion failed")
	}

	if value != "hello, world" {
		t.Fatal("should be hello, world")
	}

}
