package xmlmacro

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/foxesknow/tango/text"
)

const (
	NS       = "urn:tango.xmlmacro"
	ExpandNS = "urn:tango.xmlmacro.def"
)

type XmlExpander struct {
}

type expansionState struct {
	expandTextConfig text.ExpandTextConfig
	tokens           []xml.Token
	variables        *nameScope
}

func (state *expansionState) addToken(token xml.Token) {
	state.tokens = append(state.tokens, token)
}

func NewXmlExpander() *XmlExpander {
	return &XmlExpander{}
}

func (ex *XmlExpander) Expand(reader io.Reader) ([]xml.Token, error) {
	decoder := xml.NewDecoder(reader)
	state := &expansionState{
		expandTextConfig: text.ExpandTextConfig{},
		tokens:           make([]xml.Token, 0, 0),
		variables:        newNameScope(),
	}

	err := expand(decoder, state)

	if err != nil {
		return nil, err
	}

	return state.tokens, nil
}

func expand(decoder *xml.Decoder, state *expansionState) error {

	keepGoing := true
	for token, err := decoder.Token(); err != io.EOF && keepGoing; token, err = decoder.Token() {
		var e error

		switch element := token.(type) {
		case xml.StartElement:
			e = processStartElement(element, decoder, state)

		case xml.EndElement:
			e = processEndElement(element, decoder, state)
			keepGoing = false

		case xml.CharData:
			e = processCharData(element, decoder, state)

		case xml.Comment:
			state.addToken(xml.CopyToken(element))

		case xml.Directive:
			state.addToken(xml.CopyToken(element))

		default:
			return errors.New("unsupported xml token")
		}

		if e != nil {
			return e
		}
	}

	return nil
}

func processStartElement(element xml.StartElement, decoder *xml.Decoder, state *expansionState) error {

	if element.Name.Space == NS {
		// It's a macro command
		if action, found := macroCommandMap[element.Name.Local]; found {
			action(element, decoder, state)
			return nil
		} else {
			return fmt.Errorf("could not find xml command '%s'", element.Name.Local)
		}
	}

	copy := xml.CopyToken(element).(xml.StartElement)
	state.addToken(copy)

	// Expand the attributes
	for _, attr := range copy.Attr {
		if expandedText, err := expandText(attr.Value); err == nil {
			attr.Value = expandedText
		} else {
			return err
		}
	}

	return expand(decoder, state)
}

func processEndElement(element xml.EndElement, decoder *xml.Decoder, state *expansionState) error {
	copy := xml.CopyToken(element).(xml.EndElement)
	state.addToken(copy)

	return nil
}

func processCharData(element xml.CharData, decoder *xml.Decoder, state *expansionState) error {
	copy := xml.CopyToken(element).(xml.CharData)

	asString := string(copy)

	if expandedText, err := expandText(asString); err == nil {
		state.addToken(xml.CharData(expandedText))
	} else {
		return err
	}

	return nil
}

func expandText(value string) (string, error) {
	expandTextConfig := text.ExpandTextConfig{
		Begin: "@{",
		End:   "}",
	}

	return text.ExpandText(value, expandTextConfig)
}

// Extracts the entire element into a token slice
func extractElement(element xml.StartElement, decoder *xml.Decoder) []xml.Token {
	tokens := []xml.Token{xml.CopyToken(element)}

	// We'll need to walk the xml child tree until we find the end element at our level
	keepGoing := true
	depth := 1

	for token, err := decoder.Token(); err != io.EOF && keepGoing; token, err = decoder.Token() {
		switch next := token.(type) {
		case xml.StartElement:
			tokens = append(tokens, xml.CopyToken(next))
			depth++

		case xml.EndElement:
			tokens = append(tokens, xml.CopyToken(next))
			depth--

			if depth == 0 {
				keepGoing = false
			}

		default:
			tokens = append(tokens, xml.CopyToken(next))
		}
	}

	return tokens
}

func innerText(tokens []xml.Token) string {
	var builder strings.Builder

	for _, token := range tokens {
		if charData, ok := token.(xml.CharData); ok {
			builder.WriteString(string(charData))
		}
	}

	return builder.String()
}

// Ensures the necessary attributes exist on an element
func requireAttributes(construct string, element xml.StartElement, names ...string) ([]string, error) {
	values := make([]string, len(names))

	for i, name := range names {
		index := slices.IndexFunc(element.Attr, func(a xml.Attr) bool { return a.Name.Local == name })
		if index == -1 {
			return nil, fmt.Errorf("could not find '%s' in '%s'", name, construct)
		}

		values[i] = element.Attr[index].Value
	}

	return values, nil
}
