package diyv

import (
    "fmt"
    "reflect"
    "strings"
    "unicode/utf8"
)

type Validator struct {
    validator_funcs map[string]func(interface{}) error
}

func (vd *Validator) Validate(o interface{}) error {
    v := reflect.ValueOf(o)
    if v.Kind() == reflect.Ptr {
        v = v.Elem()
    }

    t := v.Type()
    if t == nil || t.Kind() != reflect.Struct {
        return nil
    }

    for i := 0; i < t.NumField(); i++ {
        fieldname := t.Field(i).Name
        if !v.Field(i).CanInterface() {
            continue
        }
        fieldval := v.Field(i).Interface()
        validate_as_tag := strings.TrimSpace(t.Field(i).Tag.Get("validate_as"))
        if validate_as_tag == "" {
            continue
        }
        validator_names := strings.Split(validate_as_tag, ",")
        for _, validator_name := range validator_names {
            validator_name = strings.TrimSpace(validator_name)
            if validator_name == "_struct" {
                err := vd.Validate(fieldval)
                if err != nil {
                    return err
                }
            } else if validator_name == "not_nil" {
                if reflect.ValueOf(fieldval).IsNil() {
                    return fmt.Errorf("%s[%s]: Value is nil.",
                        fieldname, "not_nil")
                }
                continue
            } else if validator_name == "skip_nil" {
                if reflect.ValueOf(fieldval).IsNil() {
                    break
                }
                continue
            } else {
                if vd.validator_funcs[validator_name] == nil {
                    return fmt.Errorf("%s[%s]: Undefined validator %s",
                        fieldname, validator_name, validator_name)
                }
                err := vd.validator_funcs[validator_name](fieldval)
                if err != nil {
                    return fmt.Errorf("%s[%s]: %s",
                        fieldname, validator_name, err)
                }
            }
        }
    }
    return nil
}

func (vd *Validator) Register(validator_name string,
    fn func(interface{}) error) {
    vd.validator_funcs[validator_name] = fn
}

func NewValidator() *Validator {
    vd := Validator{}
    vd.validator_funcs = make(map[string]func(interface{}) error)
    return &vd
}

// Extras

func ValidateStringLength(i interface{}, minchars int, maxchars int) error {
    s, err := strdual(i)
    if err != nil {
        return err
    }
    nchars := utf8.RuneCountInString(s)
    if nchars < minchars || nchars > maxchars {
        return fmt.Errorf("Invalid length")
    }
    return nil
}

func ValidateStringLengthTrimmed(i interface{},
    minchars int,
    maxchars int) error {
    s, err := strdual(i)
    if err != nil {
        return err
    }
    s = strings.TrimSpace(s)
    nchars := utf8.RuneCountInString(s)
    if nchars < minchars || nchars > maxchars {
        return fmt.Errorf("Invalid length")
    }
    return nil
}

func ValidateStringExact(i interface{}, possible_values ...string) error {
    s, err := strdual(i)
    if err != nil {
        return err
    }
    for i := 0; i < len(possible_values); i++ {
        if s == possible_values[i] {
            return nil
        }
    }
    return fmt.Errorf("String did not match any of the possible values")
}

func strdual(i interface{}) (string, error) {
    v := reflect.ValueOf(i)
    if v.Kind() == reflect.Ptr {
        s, ok := i.(*string)
        if !ok {
            return "", fmt.Errorf("Invalid type")
        }
        return *s, nil
    }
    s, ok := i.(string)
    if !ok {
        return "", fmt.Errorf("Invalid type")
    }
    return s, nil
}
