package xmlmacro

import (
	"encoding/xml"
	"io"

	"github.com/foxesknow/tango/text"
)

type XmlExpander struct {
}

type expansionState struct {
	expandTextConfig text.ExpandTextConfig
}

func NewXmlExpander() *XmlExpander {
	return &XmlExpander{}
}

func (ex *XmlExpander) Expand(reader io.Reader) ([]xml.Token, error) {
	tokens := make([]xml.Token, 0, 0)
	decoder := xml.NewDecoder(reader)
	state := &expansionState{
		expandTextConfig: text.ExpandTextConfig{},
	}

	return expand(decoder, state, tokens)
}

func expand(decoder *xml.Decoder, state *expansionState, tokens []xml.Token) ([]xml.Token, error) {

	keepGoing := true
	for token, err := decoder.Token(); err != io.EOF && keepGoing; token, err = decoder.Token() {
		var e error
		if element, ok := token.(xml.StartElement); ok {
			tokens, e = processStartElement(element, decoder, state, tokens)
		} else if element, ok := token.(xml.EndElement); ok {
			tokens, e = processEndElement(element, decoder, state, tokens)
			keepGoing = false
		} else if element, ok := token.(xml.CharData); ok {
			tokens, e = processCharData(element, decoder, state, tokens)
		} else if element, ok := token.(xml.Comment); ok {
			tokens = append(tokens, element)
		} else if element, ok := token.(xml.ProcInst); ok {
			tokens = append(tokens, element)
		} else if element, ok := token.(xml.Directive); ok {
			tokens = append(tokens, element)
		}

		if e != nil {
			return nil, e
		}
	}

	return tokens, nil
}

func processStartElement(element xml.StartElement, decoder *xml.Decoder, state *expansionState, tokens []xml.Token) ([]xml.Token, error) {
	copy := xml.CopyToken(element).(xml.StartElement)
	tokens = append(tokens, copy)

	// Expand the attributes
	for _, attr := range copy.Attr {
		if expandedText, err := expandText(attr.Value); err == nil {
			attr.Value = expandedText
		} else {
			return tokens, err
		}
	}

	return expand(decoder, state, tokens)
}

func processEndElement(element xml.EndElement, decoder *xml.Decoder, state *expansionState, tokens []xml.Token) ([]xml.Token, error) {
	copy := xml.CopyToken(element).(xml.EndElement)
	tokens = append(tokens, copy)

	return tokens, nil
}

func processCharData(element xml.CharData, decoder *xml.Decoder, state *expansionState, tokens []xml.Token) ([]xml.Token, error) {
	copy := xml.CopyToken(element).(xml.CharData)

	asString := string(copy)

	if expandedText, err := expandText(asString); err == nil {
		tokens = append(tokens, xml.CharData(expandedText))
	} else {

		return tokens, err
	}

	return tokens, nil
}

func expandText(value string) (string, error) {
	expandTextConfig := text.ExpandTextConfig{
		Begin: "@{",
		End:   "}",
	}

	return text.ExpandText(value, expandTextConfig)
}
