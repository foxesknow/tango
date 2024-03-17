package xmlmacro

import (
	"encoding/xml"
	"fmt"
)

type macroCommand func(element xml.StartElement, decoder *xml.Decoder, state *expansionState) error

var (
	macroCommandMap = map[string]macroCommand{
		"Comment": macroCommandComment,
		"Declare": macroCommandDeclare,
	}
)

func macroCommandComment(_ xml.StartElement, _ *xml.Decoder, _ *expansionState) error {
	return nil
}

func macroCommandDeclare(element xml.StartElement, decoder *xml.Decoder, state *expansionState) error {
	values, err := requireAttributes("Define", element, "name")
	if err != nil {
		return err
	}

	name := values[0]
	tokens := extractElement(element, decoder)
	value := innerText(tokens[1 : len(tokens)-1])

	if !state.variables.declare(name, value) {
		return fmt.Errorf("variable '%s' is already declared", name)
	}

	return nil
}
