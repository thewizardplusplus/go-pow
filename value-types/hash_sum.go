package powValueTypes

type HashSum struct {
	rawValue []byte
}

func NewHashSum(rawValue []byte) HashSum {
	if rawValue == nil {
		rawValue = []byte{}
	}

	return HashSum{
		rawValue: rawValue,
	}
}

func (value HashSum) Len() int {
	return len(value.rawValue)
}

func (value HashSum) ToBytes() []byte {
	return value.rawValue
}
