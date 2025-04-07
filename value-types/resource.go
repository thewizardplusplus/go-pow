package powValueTypes

import (
	"fmt"
	"net/url"
)

type Resource struct {
	rawValue *url.URL
}

func NewResource(rawValue *url.URL) Resource {
	return Resource{
		rawValue: rawValue,
	}
}

func ParseResource(rawValue string) (Resource, error) {
	parsedRawValue, err := url.Parse(rawValue)
	if err != nil {
		return Resource{}, fmt.Errorf("unable to parse the URL: %w", err)
	}

	return NewResource(parsedRawValue), nil
}

func (value Resource) ToURL() *url.URL {
	return value.rawValue
}

func (value Resource) ToString() string {
	return value.rawValue.String()
}
