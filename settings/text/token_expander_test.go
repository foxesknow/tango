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

func Test_ExpandText_UnknownVariableLookup(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
		UnknownVariableLookup: func(name string) (any, bool) {
			if name == "age" {
				return 41, true
			}

			return nil, false
		},
	}

	value, err := ExpandText("Jack is ${age}", config)
	if err != nil {
		t.Fatal("expansion failed")
	}

	if value != "Jack is 41" {
		t.Fatal("unexpected text expansion")
	}
}

func Test_ExpandText_UnknownVariableLookup_Variable_Not_Found(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
		UnknownVariableLookup: func(name string) (any, bool) {
			return nil, false
		},
	}

	value, err := ExpandText("Jack is ${age}", config)
	if err == nil {
		t.Fatal("expansion should have failed")
	}

	if value != "" {
		t.Fatal("should have returned an empty string")
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
