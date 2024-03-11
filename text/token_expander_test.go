package text

import (
	"strings"
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

func Test_ExpandText_UnknownVariableLookup_With_Namespace(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
		UnknownVariableLookup: func(name string) (any, bool) {
			if name == "lost:age" {
				return 41, true
			}

			return nil, false
		},
	}

	value, err := ExpandText("Jack is ${lost:age}", config)
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

func Test_ExpandText_UnknownVariableLookup_Variable_With_Namespace_Not_Found(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
		UnknownVariableLookup: func(name string) (any, bool) {
			return nil, false
		},
	}

	value, err := ExpandText("Jack is ${foo:age}", config)
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

	if strings.ContainsAny(value, "${}") {
		t.Fatal("text did not expand")
	}
}

func Test_ExpandText_ToUpper(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
		UnknownVariableLookup: func(name string) (any, bool) {
			if name == "name" {
				return "Jack", true
			}

			return nil, false
		},
	}

	value, err := ExpandText("${name||toupper}", config)
	if err != nil {
		t.Fatal("expansion failed")
	}

	if value != "JACK" {
		t.Fatal("text did not expand")
	}
}

func Test_ExpandText_ToLower(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
		UnknownVariableLookup: func(name string) (any, bool) {
			if name == "name" {
				return "Jack", true
			}

			return nil, false
		},
	}

	value, err := ExpandText("${name||tolower}", config)
	if err != nil {
		t.Fatal("expansion failed")
	}

	if value != "jack" {
		t.Fatal("text did not expand")
	}
}

func Test_ExpandText_Trim(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
		UnknownVariableLookup: func(name string) (any, bool) {
			if name == "name" {
				return "  Jack ", true
			}

			return nil, false
		},
	}

	value, err := ExpandText("${name||trim}", config)
	if err != nil {
		t.Fatal("expansion failed")
	}

	if value != "Jack" {
		t.Fatal("text did not expand")
	}
}

func Test_ExpandText_Unknown_Action(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
		UnknownVariableLookup: func(name string) (any, bool) {
			if name == "name" {
				return "  Jack ", true
			}

			return nil, false
		},
	}

	_, err := ExpandText("${name||foo}", config)
	if err == nil {
		t.Fatal("expansion should have failed")
	}
}

func Test_ExpandText_No_End_Token(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
	}

	_, err := ExpandText("hello, ${ world", config)
	if err == nil {
		t.Fatal("expansion should have failed")
	}
}

func Test_ExpandText_Format_Not_Supported(t *testing.T) {
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

	_, err := ExpandText("Jack is ${age|fmt}", config)
	if err == nil {
		t.Fatal("expansion should have failed")
	}
}

func Test_ExpandText_DateTime(t *testing.T) {
	config := ExpandTextConfig{
		Begin: "${",
		End:   "}",
	}

	value, err := ExpandText("${datetime:now|YYYYMMDD}", config)
	if err != nil {
		t.Fatal("expansion failed")
	}

	if len(value) != 8 {
		t.Fatal("expansion failed")
	}
}
