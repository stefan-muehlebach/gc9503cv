// auto-generated
// **** THIS FILE IS AUTO-GENERATED, PLEASE DO NOT EDIT IT **** //

package binding

import (
    "fmt"
)

type stringFromBool struct {
    base

    format string

    from Bool
}

// BoolToString creates a binding that connects a Bool data item to a String.
// Changes to the Bool will be pushed to the String and setting the string will parse and set the
// Bool if the parse was successful.
//
func BoolToString(v Bool) String {
    str := &stringFromBool{from: v}
    v.AddListener(str)
    return str
}

// BoolToStringWithFormat creates a binding that connects a Bool data item to a String and is
// presented using the specified format. Changes to the Bool will be pushed to the String and setting
// the string will parse and set the Bool if the string matches the format and its parse was successful.
//
func BoolToStringWithFormat(v Bool, format string) String {
    if format == "%t" { // Same as not using custom formatting.
        return BoolToString(v)
    }

    str := &stringFromBool{from: v, format: format}
    v.AddListener(str)
    return str
}

func (s *stringFromBool) Get() (string) {
    val := s.from.Get()

    if s.format != "" {
        return fmt.Sprintf(s.format, val)
    }
    return formatBool(val)
}

func (s *stringFromBool) Set(str string) {
    var val bool
    if s.format != "" {
        safe := stripFormatPrecision(s.format)
        n, _ := fmt.Sscanf(str, safe+" ", &val) // " " denotes match to end of string
        if n != 1 {
            return
        }
    } else {
        new, _ := parseBool(str)
        val = new
    }

    old := s.from.Get()
    if val == old {
        return
    }
    s.from.Set(val)
    s.DataChanged(s.super)
}

func (s *stringFromBool) DataChanged(data DataItem) {
    s.lock.RLock()
    defer s.lock.RUnlock()
    s.trigger()
}

type stringFromFloat struct {
    base

    format string

    from Float
}

// FloatToString creates a binding that connects a Float data item to a String.
// Changes to the Float will be pushed to the String and setting the string will parse and set the
// Float if the parse was successful.
//
func FloatToString(v Float) String {
    str := &stringFromFloat{from: v}
    v.AddListener(str)
    return str
}

// FloatToStringWithFormat creates a binding that connects a Float data item to a String and is
// presented using the specified format. Changes to the Float will be pushed to the String and setting
// the string will parse and set the Float if the string matches the format and its parse was successful.
//
func FloatToStringWithFormat(v Float, format string) String {
    if format == "%f" { // Same as not using custom formatting.
        return FloatToString(v)
    }

    str := &stringFromFloat{from: v, format: format}
    v.AddListener(str)
    return str
}

func (s *stringFromFloat) Get() (string) {
    val := s.from.Get()

    if s.format != "" {
        return fmt.Sprintf(s.format, val)
    }
    return formatFloat(val)
}

func (s *stringFromFloat) Set(str string) {
    var val float64
    if s.format != "" {
        safe := stripFormatPrecision(s.format)
        n, _ := fmt.Sscanf(str, safe+" ", &val) // " " denotes match to end of string
        if n != 1 {
            return
        }
    } else {
        new, _ := parseFloat(str)
        val = new
    }

    old := s.from.Get()
    if val == old {
        return
    }
    s.from.Set(val)
    s.DataChanged(s.super)
}

func (s *stringFromFloat) DataChanged(data DataItem) {
    s.lock.RLock()
    defer s.lock.RUnlock()
    s.trigger()
}

type stringFromInt struct {
    base

    format string

    from Int
}

// IntToString creates a binding that connects a Int data item to a String.
// Changes to the Int will be pushed to the String and setting the string will parse and set the
// Int if the parse was successful.
//
func IntToString(v Int) String {
    str := &stringFromInt{from: v}
    v.AddListener(str)
    return str
}

// IntToStringWithFormat creates a binding that connects a Int data item to a String and is
// presented using the specified format. Changes to the Int will be pushed to the String and setting
// the string will parse and set the Int if the string matches the format and its parse was successful.
//
func IntToStringWithFormat(v Int, format string) String {
    if format == "%d" { // Same as not using custom formatting.
        return IntToString(v)
    }

    str := &stringFromInt{from: v, format: format}
    v.AddListener(str)
    return str
}

func (s *stringFromInt) Get() (string) {
    val := s.from.Get()

    if s.format != "" {
        return fmt.Sprintf(s.format, val)
    }
    return formatInt(val)
}

func (s *stringFromInt) Set(str string) {
    var val int
    if s.format != "" {
        safe := stripFormatPrecision(s.format)
        n, _ := fmt.Sscanf(str, safe+" ", &val) // " " denotes match to end of string
        if n != 1 {
            return
        }
    } else {
        new, _ := parseInt(str)
        val = new
    }

    old := s.from.Get()
    if val == old {
        return
    }
    s.from.Set(val)
    s.DataChanged(s.super)
}

func (s *stringFromInt) DataChanged(data DataItem) {
    s.lock.RLock()
    defer s.lock.RUnlock()
    s.trigger()
}

type stringToBool struct {
    base

    format string

    from String
}

// StringToBool creates a binding that connects a String data item to a Bool.
// Changes to the String will be parsed and pushed to the Bool if the parse was successful, and setting
// the Bool update the String binding.
//
func StringToBool(str String) Bool {
    v := &stringToBool{from: str}
    str.AddListener(v)
    return v
}

// StringToBoolWithFormat creates a binding that connects a String data item to a Bool and is
// presented using the specified format. Changes to the Bool will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Bool will push a formatted value
// into the String.
//
func StringToBoolWithFormat(str String, format string) Bool {
    if format == "%t" { // Same as not using custom format.
        return StringToBool(str)
    }

    v := &stringToBool{from: str, format: format}
    str.AddListener(v)
    return v
}

func (s *stringToBool) Get() (bool) {
    str := s.from.Get()
    if str == "" {
        return false
    }

    var val bool
    if s.format != "" {
        n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
        if err != nil {
            return false
        }
        if n != 1 {
            return false
        }
    } else {
        new, err := parseBool(str)
        if err != nil {
            return false
        }
        val = new
    }
    return val
}

func (s *stringToBool) Set(val bool) {
    var str string
    if s.format != "" {
        str = fmt.Sprintf(s.format, val)
    } else {
        str = formatBool(val)
    }

    old := s.from.Get()
    if str == old {
        return
    }
    s.from.Set(str)
    s.DataChanged(s.super)
}

func (s *stringToBool) DataChanged(data DataItem) {
    s.lock.RLock()
    defer s.lock.RUnlock()
    s.trigger()
}

type stringToFloat struct {
    base

    format string

    from String
}

// StringToFloat creates a binding that connects a String data item to a Float.
// Changes to the String will be parsed and pushed to the Float if the parse was successful, and setting
// the Float update the String binding.
//
func StringToFloat(str String) Float {
    v := &stringToFloat{from: str}
    str.AddListener(v)
    return v
}

// StringToFloatWithFormat creates a binding that connects a String data item to a Float and is
// presented using the specified format. Changes to the Float will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Float will push a formatted value
// into the String.
//
func StringToFloatWithFormat(str String, format string) Float {
    if format == "%f" { // Same as not using custom format.
        return StringToFloat(str)
    }

    v := &stringToFloat{from: str, format: format}
    str.AddListener(v)
    return v
}

func (s *stringToFloat) Get() (float64) {
    str := s.from.Get()
    if str == "" {
        return 0.0
    }

    var val float64
    if s.format != "" {
        n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
        if err != nil {
            return 0.0
        }
        if n != 1 {
            return 0.0
        }
    } else {
        new, err := parseFloat(str)
        if err != nil {
            return 0.0
        }
        val = new
    }
    return val
}

func (s *stringToFloat) Set(val float64) {
    var str string
    if s.format != "" {
        str = fmt.Sprintf(s.format, val)
    } else {
        str = formatFloat(val)
    }

    old := s.from.Get()
    if str == old {
        return
    }
    s.from.Set(str)
    s.DataChanged(s.super)
}

func (s *stringToFloat) DataChanged(data DataItem) {
    s.lock.RLock()
    defer s.lock.RUnlock()
    s.trigger()
}

type stringToInt struct {
    base

    format string

    from String
}

// StringToInt creates a binding that connects a String data item to a Int.
// Changes to the String will be parsed and pushed to the Int if the parse was successful, and setting
// the Int update the String binding.
//
func StringToInt(str String) Int {
    v := &stringToInt{from: str}
    str.AddListener(v)
    return v
}

// StringToIntWithFormat creates a binding that connects a String data item to a Int and is
// presented using the specified format. Changes to the Int will be parsed and if the format matches and
// the parse is successful it will be pushed to the String. Setting the Int will push a formatted value
// into the String.
//
func StringToIntWithFormat(str String, format string) Int {
    if format == "%d" { // Same as not using custom format.
        return StringToInt(str)
    }

    v := &stringToInt{from: str, format: format}
    str.AddListener(v)
    return v
}

func (s *stringToInt) Get() (int) {
    str := s.from.Get()
    if str == "" {
        return 0
    }

    var val int
    if s.format != "" {
        n, err := fmt.Sscanf(str, s.format+" ", &val) // " " denotes match to end of string
        if err != nil {
            return 0
        }
        if n != 1 {
            return 0
        }
    } else {
        new, err := parseInt(str)
        if err != nil {
            return 0
        }
        val = new
    }
    return val
}

func (s *stringToInt) Set(val int) {
    var str string
    if s.format != "" {
        str = fmt.Sprintf(s.format, val)
    } else {
        str = formatInt(val)
    }

    old := s.from.Get()
    if str == old {
        return
    }
    s.from.Set(str)
    s.DataChanged(s.super)
}

func (s *stringToInt) DataChanged(data DataItem) {
    s.lock.RLock()
    defer s.lock.RUnlock()
    s.trigger()
}
