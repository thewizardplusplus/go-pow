package powValueTypes

type Payload struct {
	rawValue string
}

func NewPayload(rawValue string) Payload {
	return Payload{
		rawValue: rawValue,
	}
}

func (value Payload) ToString() string {
	return value.rawValue
}
