// The mapstructure package exposes functionality to convert an
// abitrary map[string]interface{} into a native SafeGo structure.
//
// The SafeGo structure can be arbitrarily complex, containing slices,
// other structs, etc. and the decoder will properly decode nested
// maps and so on into the proper structures in the native SafeGo struct.
// See the examples to see what the decoder is capable of.
package mapkit

import (
	"errors"
	"fmt"
	"github.com/textthree/cvgokit/strkit"
	"github.com/textthree/cvgokit/timekit"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

// Error implements the error interface and can represents multiple
// errors that occur in the course of a single decode.
type Error struct {
	Errors []string
}

func (e *Error) Error() string {
	points := make([]string, len(e.Errors))
	for i, err := range e.Errors {
		points[i] = fmt.Sprintf("* %s", err)
	}

	return fmt.Sprintf(
		"%d error(s) decoding:\n\n%s",
		len(e.Errors), strings.Join(points, "\n"))
}

func appendErrors(errors []string, err error) []string {
	switch e := err.(type) {
	case *Error:
		return append(errors, e.Errors...)
	default:
		return append(errors, e.Error())
	}
}

type DecodeHookFunc func(reflect.Kind, reflect.Kind, interface{}) (interface{}, error)

// DecoderConfig is the configuration that is used to add a new decoder
// and allows customization of various aspects of decoding.
type DecoderConfig struct {
	// DecodeHook, if set, will be called before any decoding and any
	// dto conversion (if WeaklyTypedInput is on). This lets you modify
	// the values before they're set down onto the resulting struct.
	//
	// If an error is returned, the entire decode will fail with that
	// error.
	DecodeHook DecodeHookFunc

	// If ErrorUnused is true, then it is an error for there to exist
	// keys in the original map that were unused in the decoding process
	// (extra keys).
	ErrorUnused bool

	// If WeaklyTypedInput is true, the decoder will make the following
	// "weak" conversions:
	//
	//   - bools to string (true = "1", false = "0")
	//   - numbers to string (base 10)
	//   - bools to int/uint (true = 1, false = 0)
	//   - strings to int/uint (base implied by prefix)
	//   - int to bool (true if value != 0)
	//   - string to bool (accepts: 1, t, T, TRUE, true, True, 0, f, F,
	//     FALSE, false, False. Anything else is an error)
	//   - empty array = empty map and vice versa
	//
	WeaklyTypedInput bool

	// Metadata is the struct that will contain extra metadata about
	// the decoding. If this is nil, then no metadata will be tracked.
	Metadata *Metadata

	// Result is a pointer to the struct that will contain the decoded
	// value.
	Result interface{}

	// The tag name that mapstructure reads for field names. This
	// defaults to "mapstructure"
	TagName string
}

// A Decoder takes a raw interface value and turns it into structured
// data, keeping track of rich error information along the way in case
// anything goes wrong. Unlike the basic top-level Decode method, you can
// more finely control how the Decoder behaves using the DecoderConfig
// structure. The top-level Decode method is just a convenience that sets
// up the most basic Decoder.
type Decoder struct {
	config *DecoderConfig
}

// Metadata contains information about decoding a structure that
// is tedious or difficult to get otherwise.
type Metadata struct {
	// Keys are the keys of the structure which were successfully decoded
	Keys []string

	// Unused is a slice of keys that were found in the raw value but
	// weren't decoded since there was no matching field in the result interface
	Unused []string
}

// Decode takes a map and uses reflection to convert it into the
// given SafeGo native structure. val must be a pointer to a struct.
func Decode(m interface{}, rawVal interface{}) error {
	config := &DecoderConfig{
		Metadata: nil,
		Result:   rawVal,
	}

	decoder, err := NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(m)
}

// DecodePath takes a map and uses reflection to convert it into the
// given SafeGo native structure. Tags are used to specify the mapping
// between fields in the map and structure
func DecodePath(m map[string]interface{}, rawVal interface{}) error {
	config := &DecoderConfig{
		Metadata: nil,
		Result:   nil,
	}

	decoder, err := NewPathDecoder(config)
	if err != nil {
		return err
	}

	_, err = decoder.DecodePath(m, rawVal)
	return err
}

// DecodeSlicePath decodes a slice of maps against a slice of structures that
// contain specified tags
func DecodeSlicePath(ms []map[string]interface{}, rawSlice interface{}) error {
	reflectRawSlice := reflect.TypeOf(rawSlice)
	rawKind := reflectRawSlice.Kind()
	rawElement := reflectRawSlice.Elem()

	if (rawKind == reflect.Ptr && rawElement.Kind() != reflect.Slice) ||
		(rawKind != reflect.Ptr && rawKind != reflect.Slice) {
		return fmt.Errorf("Incompatible Value, Looking For Slice : %v : %v", rawKind, rawElement.Kind())
	}

	config := &DecoderConfig{
		Metadata: nil,
		Result:   nil,
	}

	decoder, err := NewPathDecoder(config)
	if err != nil {
		return err
	}

	// Create a slice large enough to decode all the values
	valSlice := reflect.MakeSlice(rawElement, len(ms), len(ms))

	// Iterate over the maps and decode each one
	for index, m := range ms {
		sliceElementType := rawElement.Elem()
		if sliceElementType.Kind() != reflect.Ptr {
			// A slice of objects
			obj := reflect.New(rawElement.Elem())
			decoder.DecodePath(m, reflect.Indirect(obj))
			indexVal := valSlice.Index(index)
			indexVal.Set(reflect.Indirect(obj))
		} else {
			// A slice of pointers
			obj := reflect.New(rawElement.Elem().Elem())
			decoder.DecodePath(m, reflect.Indirect(obj))
			indexVal := valSlice.Index(index)
			indexVal.Set(obj)
		}
	}

	// AddRoute the new slice
	reflect.ValueOf(rawSlice).Elem().Set(valSlice)
	return nil
}

// NewDecoder returns a new decoder for the given configuration. Once
// a decoder has been returned, the same configuration must not be used
// again.
func NewDecoder(config *DecoderConfig) (*Decoder, error) {
	val := reflect.ValueOf(config.Result)
	if val.Kind() != reflect.Ptr {
		return nil, errors.New("result must be a pointer")
	}

	val = val.Elem()
	if !val.CanAddr() {
		return nil, errors.New("result must be addressable (a pointer)")
	}

	if config.Metadata != nil {
		if config.Metadata.Keys == nil {
			config.Metadata.Keys = make([]string, 0)
		}

		if config.Metadata.Unused == nil {
			config.Metadata.Unused = make([]string, 0)
		}
	}

	if config.TagName == "" {
		config.TagName = "mapstructure"
	}

	result := &Decoder{
		config: config,
	}

	return result, nil
}

// NewPathDecoder returns a new decoder for the given configuration.
// This is used to decode path specific structures
func NewPathDecoder(config *DecoderConfig) (*Decoder, error) {
	if config.Metadata != nil {
		if config.Metadata.Keys == nil {
			config.Metadata.Keys = make([]string, 0)
		}

		if config.Metadata.Unused == nil {
			config.Metadata.Unused = make([]string, 0)
		}
	}

	if config.TagName == "" {
		config.TagName = "mapstructure"
	}

	result := &Decoder{
		config: config,
	}

	return result, nil
}

// Decode decodes the given raw interface to the target pointer specified
// by the configuration.
func (d *Decoder) Decode(raw interface{}) error {
	return d.decode("", raw, reflect.ValueOf(d.config.Result).Elem())
}

// DecodePath decodes the raw interface against the map based on the
// specified tags
func (d *Decoder) DecodePath(m map[string]interface{}, rawVal interface{}) (bool, error) {
	decoded := false

	var val reflect.Value
	reflectRawValue := reflect.ValueOf(rawVal)
	kind := reflectRawValue.Kind()

	// Looking for structs and pointers to structs
	switch kind {
	case reflect.Ptr:
		val = reflectRawValue.Elem()
		if val.Kind() != reflect.Struct {
			return decoded, fmt.Errorf("Incompatible MsgType : %v : Looking For Struct", kind)
		}
	case reflect.Struct:
		var ok bool
		val, ok = rawVal.(reflect.Value)
		if ok == false {
			return decoded, fmt.Errorf("Incompatible MsgType : %v : Looking For reflect.Value", kind)
		}
	default:
		return decoded, fmt.Errorf("Incompatible MsgType : %v", kind)
	}

	// Iterate over the fields in the struct
	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		tag := typeField.Tag
		tagValue := tag.Get("jpath")

		// Is this a field without a tag
		if tagValue == "" {
			if valueField.Kind() == reflect.Struct {
				// We have a struct that may have indivdual tags. Process separately
				d.DecodePath(m, valueField)
				continue
			} else if valueField.Kind() == reflect.Ptr && reflect.TypeOf(valueField).Kind() == reflect.Struct {
				// We have a pointer to a struct
				if valueField.IsNil() {
					// Create the object since it doesn't exist
					valueField.Set(reflect.New(valueField.Type().Elem()))
					decoded, _ = d.DecodePath(m, valueField.Elem())
					if decoded == false {
						// If nothing was decoded for this object return the pointer to nil
						valueField.Set(reflect.NewAt(valueField.Type().Elem(), nil))
					}
					continue
				}

				d.DecodePath(m, valueField.Elem())
				continue
			}
		}

		// Use mapstructure to populate the fields
		keys := strings.Split(tagValue, ".")
		data := d.findData(m, keys)
		if data != nil {
			if valueField.Kind() == reflect.Slice {
				// Ignore a slice of maps - This sucks but not sure how to check
				if strings.Contains(valueField.Type().String(), "map[") {
					goto normal_decode
				}

				// We have a slice
				mapSlice := data.([]interface{})
				if len(mapSlice) > 0 {
					// Test if this is a slice of more maps
					_, ok := mapSlice[0].(map[string]interface{})
					if ok == false {
						goto normal_decode
					}

					// Extract the maps out and run it through DecodeSlicePath
					ms := make([]map[string]interface{}, len(mapSlice))
					for index, m2 := range mapSlice {
						ms[index] = m2.(map[string]interface{})
					}

					DecodeSlicePath(ms, valueField.Addr().Interface())
					continue
				}
			}
		normal_decode:
			decoded = true
			err := d.decode("", data, valueField)
			if err != nil {
				return false, err
			}
		}
	}

	return decoded, nil
}

// Decodes an unknown data dto into a specific reflection value.
func (d *Decoder) decode(name string, data interface{}, val reflect.Value) error {
	if data == nil {
		// If the data is nil, then we don't set anything.
		return nil
	}

	dataVal := reflect.ValueOf(data)
	if !dataVal.IsValid() {
		// If the data value is invalid, then we just set the value
		// to be the zero value.
		val.Set(reflect.Zero(val.Type()))
		return nil
	}

	if d.config.DecodeHook != nil {
		// We have a DecodeHook, so let's pre-process the data.
		var err error
		data, err = d.config.DecodeHook(d.getKind(dataVal), d.getKind(val), data)
		if err != nil {
			return err
		}
	}

	var err error
	dataKind := d.getKind(val)
	switch dataKind {
	case reflect.Bool:
		err = d.decodeBool(name, data, val)
	case reflect.Interface:
		err = d.decodeBasic(name, data, val)
	case reflect.String:
		err = d.decodeString(name, data, val)
	case reflect.Int:
		err = d.decodeInt(name, data, val)
	case reflect.Uint:
		err = d.decodeUint(name, data, val)
	case reflect.Float32:
		err = d.decodeFloat(name, data, val)
	case reflect.Struct:
		err = d.decodeStruct(name, data, val)
	case reflect.Map:
		err = d.decodeMap(name, data, val)
	case reflect.Slice:
		err = d.decodeSlice(name, data, val)
	default:
		// If we reached this point then we weren't able to decode it
		return fmt.Errorf("%s: unsupported dto: %s", name, dataKind)
	}

	// If we reached here, then we successfully decoded SOMETHING, so
	// mark the key as used if we're tracking metadata.
	if d.config.Metadata != nil && name != "" {
		d.config.Metadata.Keys = append(d.config.Metadata.Keys, name)
	}

	return err
}

// findData locates the data by walking the keys down the map
func (d *Decoder) findData(m map[string]interface{}, keys []string) interface{} {
	if len(keys) == 1 {
		if value, ok := m[keys[0]]; ok == true {
			return value
		}
		return nil
	}

	if value, ok := m[keys[0]]; ok == true {
		if m, ok := value.(map[string]interface{}); ok == true {
			return d.findData(m, keys[1:])
		}
	}

	return nil
}

func (d *Decoder) getKind(val reflect.Value) reflect.Kind {
	kind := val.Kind()

	switch {
	case kind >= reflect.Int && kind <= reflect.Int64:
		return reflect.Int
	case kind >= reflect.Uint && kind <= reflect.Uint64:
		return reflect.Uint
	case kind >= reflect.Float32 && kind <= reflect.Float64:
		return reflect.Float32
	default:
		return kind
	}
}

// This decodes a basic dto (bool, int, string, etc.) and sets the
// value to "data" of that dto.
func (d *Decoder) decodeBasic(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	dataValType := dataVal.Type()
	if !dataValType.AssignableTo(val.Type()) {
		return fmt.Errorf(
			"'%s' expected dto '%s', got '%s'",
			name, val.Type(), dataValType)
	}

	val.Set(dataVal)
	return nil
}

func (d *Decoder) decodeString(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	dataKind := d.getKind(dataVal)
	switch {
	case dataKind == reflect.String:
		val.SetString(dataVal.String())
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetString("1")
		} else {
			val.SetString("0")
		}
	case dataKind == reflect.Int && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatInt(dataVal.Int(), 10))
	case dataKind == reflect.Uint && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatUint(dataVal.Uint(), 10))
	case dataKind == reflect.Float32 && d.config.WeaklyTypedInput:
		val.SetString(strconv.FormatFloat(dataVal.Float(), 'f', -1, 64))
	// 2021-06-16 20:31:28 +0800 CST
	case reflect.TypeOf(data).String() == "time.Time":
		val.SetString(timekit.DateTimeFormat(data))
	default:
		return fmt.Errorf(
			"'%s' decodeString expected dto '%s', got unconvertible dto '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeInt(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	dataKind := d.getKind(dataVal)
	switch {
	case dataKind == reflect.Int:
		val.SetInt(dataVal.Int())
	case dataKind == reflect.Uint:
		val.SetInt(int64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetInt(int64(dataVal.Float()))
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetInt(1)
		} else {
			val.SetInt(0)
		}
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		i, err := strconv.ParseInt(dataVal.String(), 0, val.Type().Bits())
		if err == nil {
			val.SetInt(i)
		} else {
			return fmt.Errorf("cannot parse '%s' as int: %s", name, err)
		}
	// 字符串强转int
	case dataKind == reflect.String:
		switch val.Type().String() {
		case "int":
			integerValue, err := strconv.Atoi(dataVal.String())
			val.SetInt(int64(integerValue))
			if err != nil {
				fmt.Println(err)
			}
		case "int8":
			int8, err := strconv.ParseInt(dataVal.String(), 10, 8)
			val.SetInt(int64(int8))
			if err != nil {
				fmt.Println(err)
			}
		case "int16":
			int16, err := strconv.ParseInt(dataVal.String(), 10, 16)
			val.SetInt(int64(int16))
			if err != nil {
				fmt.Println(err)
			}
		case "int32":
			int32, err := strconv.ParseInt(dataVal.String(), 10, 32)
			val.SetInt(int64(int32))
			if err != nil {
				fmt.Println(err)
			}
		case "int64":
			int64, err := strconv.ParseInt(dataVal.String(), 10, 64)
			val.SetInt(int64)
			if err != nil {
				fmt.Println(err)
			}
		}
	default:
		return fmt.Errorf(
			"'%s' decodeInt expected dto '%s', got unconvertible dto '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeUint(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	dataKind := d.getKind(dataVal)

	switch {
	case dataKind == reflect.Int:
		val.SetUint(uint64(dataVal.Int()))
	case dataKind == reflect.Uint:
		val.SetUint(dataVal.Uint())
	case dataKind == reflect.Float32:
		val.SetUint(uint64(dataVal.Float()))
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetUint(1)
		} else {
			val.SetUint(0)
		}
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		i, err := strconv.ParseUint(dataVal.String(), 0, val.Type().Bits())
		if err == nil {
			val.SetUint(i)
		} else {
			return fmt.Errorf("cannot parse '%s' as uint: %s", name, err)
		}
	default:
		return fmt.Errorf(
			"'%s' decodeUint expected dto '%s', got unconvertible dto '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeBool(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	dataKind := d.getKind(dataVal)

	switch {
	case dataKind == reflect.Bool:
		val.SetBool(dataVal.Bool())
	case dataKind == reflect.Int && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Int() != 0)
	case dataKind == reflect.Uint && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Uint() != 0)
	case dataKind == reflect.Float32 && d.config.WeaklyTypedInput:
		val.SetBool(dataVal.Float() != 0)
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		b, err := strconv.ParseBool(dataVal.String())
		if err == nil {
			val.SetBool(b)
		} else if dataVal.String() == "" {
			val.SetBool(false)
		} else {
			return fmt.Errorf("cannot parse '%s' as bool: %s", name, err)
		}
	default:
		return fmt.Errorf(
			"'%s' decodeBool expected dto '%s', got unconvertible dto '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeFloat(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.ValueOf(data)
	dataKind := d.getKind(dataVal)

	switch {
	case dataKind == reflect.Int:
		val.SetFloat(float64(dataVal.Int()))
	case dataKind == reflect.Uint:
		val.SetFloat(float64(dataVal.Uint()))
	case dataKind == reflect.Float32:
		val.SetFloat(float64(dataVal.Float()))
	case dataKind == reflect.Bool && d.config.WeaklyTypedInput:
		if dataVal.Bool() {
			val.SetFloat(1)
		} else {
			val.SetFloat(0)
		}
	case dataKind == reflect.String && d.config.WeaklyTypedInput:
		f, err := strconv.ParseFloat(dataVal.String(), val.Type().Bits())
		if err == nil {
			val.SetFloat(f)
		} else {
			return fmt.Errorf("cannot parse '%s' as float: %s", name, err)
		}
	case dataKind == reflect.String:
		val.SetFloat(strkit.StringToFloat64(dataVal.String()))
	default:
		return fmt.Errorf(
			"'%s' decodeFloat expected dto '%s', got unconvertible dto '%s'",
			name, val.Type(), dataVal.Type())
	}

	return nil
}

func (d *Decoder) decodeMap(name string, data interface{}, val reflect.Value) error {
	valType := val.Type()
	valKeyType := valType.Key()
	valElemType := valType.Elem()

	// Make a new map to hold our result
	mapType := reflect.MapOf(valKeyType, valElemType)
	valMap := reflect.MakeMap(mapType)

	// Check input dto
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	if dataVal.Kind() != reflect.Map {
		// Accept empty array/slice instead of an empty map in weakly typed mode
		if d.config.WeaklyTypedInput &&
			(dataVal.Kind() == reflect.Slice || dataVal.Kind() == reflect.Array) &&
			dataVal.Len() == 0 {
			val.Set(valMap)
			return nil
		} else {
			return fmt.Errorf("'%s' decodeMap expected a map, got '%s'", name, dataVal.Kind())
		}
	}

	// Accumulate errors
	errors := make([]string, 0)

	for _, k := range dataVal.MapKeys() {
		fieldName := fmt.Sprintf("%s[%s]", name, k)

		// First decode the key into the proper dto
		currentKey := reflect.Indirect(reflect.New(valKeyType))
		if err := d.decode(fieldName, k.Interface(), currentKey); err != nil {
			errors = appendErrors(errors, err)
			continue
		}

		// Next decode the data into the proper dto
		v := dataVal.MapIndex(k).Interface()
		currentVal := reflect.Indirect(reflect.New(valElemType))
		if err := d.decode(fieldName, v, currentVal); err != nil {
			errors = appendErrors(errors, err)
			continue
		}

		valMap.SetMapIndex(currentKey, currentVal)
	}

	// AddRoute the built up map to the value
	val.Set(valMap)

	// If we had errors, return those
	if len(errors) > 0 {
		return &Error{errors}
	}

	return nil
}

func (d *Decoder) decodeSlice(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	valType := val.Type()
	valElemType := valType.Elem()

	// Make a new slice to hold our result, same size as the original data.
	sliceType := reflect.SliceOf(valElemType)
	valSlice := reflect.MakeSlice(sliceType, dataVal.Len(), dataVal.Len())

	// Check input dto
	if dataValKind != reflect.Array && dataValKind != reflect.Slice {
		// Accept empty map instead of array/slice in weakly typed mode
		if d.config.WeaklyTypedInput && dataVal.Kind() == reflect.Map && dataVal.Len() == 0 {
			val.Set(valSlice)
			return nil
		} else {
			return fmt.Errorf(
				"'%s': source data must be an array or slice, got %s", name, dataValKind)
		}
	}

	// Accumulate any errors
	errors := make([]string, 0)

	for i := 0; i < dataVal.Len(); i++ {
		currentData := dataVal.Index(i).Interface()
		currentField := valSlice.Index(i)

		fieldName := fmt.Sprintf("%s[%d]", name, i)
		if err := d.decode(fieldName, currentData, currentField); err != nil {
			errors = appendErrors(errors, err)
		}
	}

	// Finally, set the value to the slice we built up
	val.Set(valSlice)

	// If there were errors, we return those
	if len(errors) > 0 {
		return &Error{errors}
	}

	return nil
}

func (d *Decoder) decodeStruct(name string, data interface{}, val reflect.Value) error {
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	dataValKind := dataVal.Kind()
	if dataValKind != reflect.Map {
		return fmt.Errorf("'%s' expected a map, got '%s'", name, dataValKind)
	}

	dataValType := dataVal.Type()
	if kind := dataValType.Key().Kind(); kind != reflect.String && kind != reflect.Interface {
		return fmt.Errorf(
			"'%s' needs a map with string keys, has '%s' keys",
			name, dataValType.Key().Kind())
	}

	dataValKeys := make(map[reflect.Value]struct{})
	dataValKeysUnused := make(map[interface{}]struct{})
	for _, dataValKey := range dataVal.MapKeys() {
		dataValKeys[dataValKey] = struct{}{}
		dataValKeysUnused[dataValKey.Interface()] = struct{}{}
	}

	errors := make([]string, 0)

	// This slice will keep track of all the structs we'll be decoding.
	// There can be more than one struct if there are embedded structs
	// that are squashed.
	structs := make([]reflect.Value, 1, 5)
	structs[0] = val

	// Compile the list of all the fields that we're going to be decoding
	// from all the structs.
	fields := make(map[*reflect.StructField]reflect.Value)
	for len(structs) > 0 {
		structVal := structs[0]
		structs = structs[1:]

		structType := structVal.Type()
		for i := 0; i < structType.NumField(); i++ {
			fieldType := structType.Field(i)

			if fieldType.Anonymous {
				fieldKind := fieldType.Type.Kind()
				if fieldKind != reflect.Struct {
					errors = appendErrors(errors,
						fmt.Errorf("%s: unsupported dto: %s", fieldType.Name, fieldKind))
					continue
				}

				// We have an embedded field. We "squash" the fields down
				// if specified in the tag.
				squash := false
				tagParts := strings.Split(fieldType.Tag.Get(d.config.TagName), ",")
				for _, tag := range tagParts[1:] {
					if tag == "squash" {
						squash = true
						break
					}
				}

				if squash {
					structs = append(structs, val.FieldByName(fieldType.Name))
					continue
				}
			}

			// Normal struct field, store it away
			fields[&fieldType] = structVal.Field(i)
		}
	}

	for fieldType, field := range fields {
		fieldName := fieldType.Name

		tagValue := fieldType.Tag.Get(d.config.TagName)
		tagValue = strings.SplitN(tagValue, ",", 2)[0]
		if tagValue != "" {
			fieldName = tagValue
		}

		rawMapKey := reflect.ValueOf(fieldName)
		rawMapVal := dataVal.MapIndex(rawMapKey)
		if !rawMapVal.IsValid() {
			// Do a slower search by iterating over each key and
			// doing case-insensitive search.
			for dataValKey, _ := range dataValKeys {
				mK, ok := dataValKey.Interface().(string)
				if !ok {
					// Not a string key
					continue
				}

				if strings.EqualFold(mK, fieldName) {
					rawMapKey = dataValKey
					rawMapVal = dataVal.MapIndex(dataValKey)
					break
				}
			}

			if !rawMapVal.IsValid() {
				// There was no matching key in the map for the value in
				// the struct. Just ignore.
				continue
			}
		}

		// Delete the key we're using from the unused map so we stop tracking
		delete(dataValKeysUnused, rawMapKey.Interface())

		if !field.IsValid() {
			// This should never happen
			panic("field is not valid")
		}

		// If we can't set the field, then it is unexported or something,
		// and we just continue onwards.
		if !field.CanSet() {
			continue
		}

		// If the name is empty string, then we're at the root, and we
		// don't dot-join the fields.
		if name != "" {
			fieldName = fmt.Sprintf("%s.%s", name, fieldName)
		}

		if err := d.decode(fieldName, rawMapVal.Interface(), field); err != nil {
			errors = appendErrors(errors, err)
		}
	}

	if d.config.ErrorUnused && len(dataValKeysUnused) > 0 {
		keys := make([]string, 0, len(dataValKeysUnused))
		for rawKey, _ := range dataValKeysUnused {
			keys = append(keys, rawKey.(string))
		}
		sort.Strings(keys)

		err := fmt.Errorf("'%s' has invalid keys: %s", name, strings.Join(keys, ", "))
		errors = appendErrors(errors, err)
	}

	if len(errors) > 0 {
		return &Error{errors}
	}

	// Add the unused keys to the list of unused keys if we're tracking metadata
	if d.config.Metadata != nil {
		for rawKey, _ := range dataValKeysUnused {
			key := rawKey.(string)
			if name != "" {
				key = fmt.Sprintf("%s.%s", name, key)
			}

			d.config.Metadata.Unused = append(d.config.Metadata.Unused, key)
		}
	}

	return nil
}
