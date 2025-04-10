package powValueTypes

type SerializedPayload struct {
	rawValue string
}

func NewSerializedPayload(rawValue string) SerializedPayload {
	return SerializedPayload{
		rawValue: rawValue,
	}
}

func (value SerializedPayload) ToString() string {
	return value.rawValue
}
