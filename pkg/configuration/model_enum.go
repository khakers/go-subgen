// Code generated by "go-enum -type=Model -all=false -string=true -new=true -string=true -text=true -json=true -yaml=false"; DO NOT EDIT.

// Install go-enum by `go get install github.com/searKing/golang/tools/go-enum`
package configuration

import (
	"encoding"
	"encoding/json"
	"fmt"
	"strconv"
)

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Tiny_en-0]
	_ = x[Tiny-1]
	_ = x[Base_en-2]
	_ = x[Base-3]
	_ = x[Small_en-4]
	_ = x[Small-5]
	_ = x[Medium_en-6]
	_ = x[Medium-7]
	_ = x[Large_v1-8]
	_ = x[Large-9]
	_ = x[Large_v3-10]
	_ = x[Large_v2-11]
}

const _Model_name = "Tiny_enTinyBase_enBaseSmall_enSmallMedium_enMediumLarge_v1LargeLarge_v3Large_v2"

var _Model_index = [...]uint8{0, 7, 11, 18, 22, 30, 35, 44, 50, 58, 63, 71, 79}

func _() {
	var _nil_Model_value = func() (val Model) { return }()

	// An "cannot convert Model literal (type Model) to type fmt.Stringer" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ fmt.Stringer = _nil_Model_value
}

func (i Model) String() string {
	if i >= Model(len(_Model_index)-1) {
		return "Model(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Model_name[_Model_index[i]:_Model_index[i+1]]
}

// New returns a pointer to a new addr filled with the Model value passed in.
func (i Model) New() *Model {
	clone := i
	return &clone
}

var _Model_values = []Model{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}

var _Model_name_to_values = map[string]Model{
	_Model_name[0:7]:   0,
	_Model_name[7:11]:  1,
	_Model_name[11:18]: 2,
	_Model_name[18:22]: 3,
	_Model_name[22:30]: 4,
	_Model_name[30:35]: 5,
	_Model_name[35:44]: 6,
	_Model_name[44:50]: 7,
	_Model_name[50:58]: 8,
	_Model_name[58:63]: 9,
	_Model_name[63:71]: 10,
	_Model_name[71:79]: 11,
}

// ParseModelString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func ParseModelString(s string) (Model, error) {
	if val, ok := _Model_name_to_values[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to Model values", s)
}

// ModelValues returns all values of the enum
func ModelValues() []Model {
	return _Model_values
}

// IsAModel returns "true" if the value is listed in the enum definition. "false" otherwise
func (i Model) Registered() bool {
	for _, v := range _Model_values {
		if i == v {
			return true
		}
	}
	return false
}

func _() {
	var _nil_Model_value = func() (val Model) { return }()

	// An "cannot convert Model literal (type Model) to type json.Marshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ json.Marshaler = _nil_Model_value

	// An "cannot convert Model literal (type Model) to type encoding.Unmarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ json.Unmarshaler = &_nil_Model_value
}

// MarshalJSON implements the json.Marshaler interface for Model
func (i Model) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for Model
func (i *Model) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Model should be a string, got %s", data)
	}

	var err error
	*i, err = ParseModelString(s)
	return err
}

func _() {
	var _nil_Model_value = func() (val Model) { return }()

	// An "cannot convert Model literal (type Model) to type encoding.TextMarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.TextMarshaler = _nil_Model_value

	// An "cannot convert Model literal (type Model) to type encoding.TextUnmarshaler" compiler error signifies that the base type have changed.
	// Re-run the go-enum command to generate them again.
	var _ encoding.TextUnmarshaler = &_nil_Model_value
}

// MarshalText implements the encoding.TextMarshaler interface for Model
func (i Model) MarshalText() ([]byte, error) {
	return []byte(i.String()), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for Model
func (i *Model) UnmarshalText(text []byte) error {
	var err error
	*i, err = ParseModelString(string(text))
	return err
}
