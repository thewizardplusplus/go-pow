package powValueTypes

import (
	"bytes"
	"fmt"
	"text/template"
)

type HashDataLayout struct {
	rawValue *template.Template
}

func NewHashDataLayout(rawValue *template.Template) HashDataLayout {
	return HashDataLayout{
		rawValue: rawValue,
	}
}

func ParseHashDataLayout(rawValue string) (HashDataLayout, error) {
	parsedRawValue, err := template.New("").Parse(rawValue)
	if err != nil {
		return HashDataLayout{}, fmt.Errorf(
			"unable to parse the text template: %w",
			err,
		)
	}

	return NewHashDataLayout(parsedRawValue), nil
}

func MustParseHashDataLayout(rawValue string) HashDataLayout {
	value, err := ParseHashDataLayout(rawValue)
	if err != nil {
		panic(fmt.Sprintf("powValueTypes.MustParseHashDataLayout(): %s", err))
	}

	return value
}

func (value HashDataLayout) Execute(data any) (string, error) {
	var buffer bytes.Buffer
	if err := value.rawValue.Execute(&buffer, data); err != nil {
		return "", fmt.Errorf("unable to execute the text template: %w", err)
	}

	return buffer.String(), nil
}

func (value HashDataLayout) ToTemplate() *template.Template {
	return value.rawValue
}

func (value HashDataLayout) ToString() string {
	return value.rawValue.Root.String()
}
