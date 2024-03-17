package xmlmacro

import (
	"bytes"
	"encoding/xml"
	"strings"
	"testing"
)

func Test_Expander(t *testing.T) {
	data := `
        <Root>
            <Name>@{os:pid}</Name>
        </Root>
	`

	expander := NewXmlExpander()
	reader := strings.NewReader(data)
	tokens, err := expander.Expand(reader)

	if err != nil {
		t.Fatal("expansion failed")
	}

	if tokens == nil {
		t.Fatal("expansion failed")
	}

	writer := &bytes.Buffer{}
	encoder := xml.NewEncoder(writer)

	for _, token := range tokens {
		if err = encoder.EncodeToken(token); err != nil {
			t.Fatal("oops")
		}
	}

	encoder.Close()

	value := writer.String()
	if len(value) == 0 {
		t.Fatal("oops")
	}
}

func Test_Expander_Declare(t *testing.T) {
	data := `
        <Root xmlns:m="urn:tango.xmlmacro" xmlns:x="urn:tango.xmlmacro.def">
			<m:Declare name="who">Ben</m:Declare>
            <Name>@{os:pid}</Name>
        </Root>
	`

	expander := NewXmlExpander()
	reader := strings.NewReader(data)
	tokens, err := expander.Expand(reader)

	if err != nil {
		t.Fatal("expansion failed")
	}

	if tokens == nil {
		t.Fatal("expansion failed")
	}

	writer := &bytes.Buffer{}
	encoder := xml.NewEncoder(writer)

	for _, token := range tokens {
		if err = encoder.EncodeToken(token); err != nil {
			t.Fatal("oops")
		}
	}

	encoder.Close()

	value := writer.String()
	if len(value) == 0 {
		t.Fatal("oops")
	}
}
