package diyv

import (
    "fmt"
    "testing"
)

func TestValidator_single_string_field(t *testing.T) {
    vd := NewValidator()

    type Foo struct {
        Bar string `validate_as:"hello"`
    }

    vd.Register("hello", func(i interface{}) error {
        s, ok := i.(string)
        if !ok {
            return fmt.Errorf("Not a string")
        }
        if s != "Hello World!" {
            return fmt.Errorf("'%s' != 'Hello World!'", s)
        }
        return nil
    })

    f := Foo{"Hello World!"}
    err := vd.Validate(f)
    if err != nil {
        t.Fail()
    }

    f = Foo{"malkmf2nff2efn2oefnlf"}
    err = vd.Validate(f)
    if err == nil {
        t.Fail()
    }
}

func TestValidator_not_nil(t *testing.T) {
    vd := NewValidator()
    type Foo struct {
        Bar *string `valid:"not_nil"`
    }

    f := Foo{}
    err := vd.Validate(f)
    if err == nil {
        t.Fail()
    }

    s := "foo bar"
    f = Foo{&s}
    err = vd.Validate(f)
    if err != nil {
        t.Fail()
    }
}

func TestValidator_skip_nil(t *testing.T) {
    vd := NewValidator()
    type Foo struct {
        Bar *string `valid:"skip_nil,alwaysfail"`
    }
    vd.Register("alwaysfail", func(i interface{}) error {
        return fmt.Errorf("fail")
    })

    f := Foo{}
    err := vd.Validate(f)
    if err != nil {
        t.Fail()
    }
}

func TestValidator_undefined_validator(t *testing.T) {
    vd := NewValidator()
    type Foo struct {
        Bar string `validate_as:"name"`
    }
    f := Foo{"foo bar"}
    err := vd.Validate(f)
    if err == nil {
        t.Fail()
    }
}

// Extras

func TestValidator_StringLengthValidator(t *testing.T) {
    vd := NewValidator()

    vd.Register("shortstr", func(i interface{}) error {
        return ValidateStringLength(i, 1, 10)
    })

    type Foo struct {
        Bar string `validate_as:"shortstr"`
    }

    f := Foo{"hello"}
    err := vd.Validate(f)
    if err != nil {
        t.Fail()
    }

    f = Foo{"The quick brown fox"}
    err = vd.Validate(f)
    if err == nil {
        t.Fail()
    }

    type Bar struct {
        Foo *string `validate_as:"shortstr"`
    }

    s := "hello"
    b := Bar{&s}
    err = vd.Validate(b)
    if err != nil {
        t.Fail()
    }

    s = "The quick brown fox"
    b = Bar{&s}
    err = vd.Validate(b)
    if err == nil {
        t.Fail()
    }
}

func TestValidator_ValidateStringExact(t *testing.T) {
    vd := NewValidator()
    vd.Register("arch", func(i interface{}) error {
        return ValidateStringExact(i, "x86", "x86_64")
    })
    type Foo struct {
        Arch string `validate_as:"arch"`
    }
    f := Foo{Arch: "x86"}
    err := vd.Validate(f)
    if err != nil {
        t.Fail()
    }
    f2 := Foo{Arch: "x86_64"}
    err = vd.Validate(f2)
    if err != nil {
        t.Fail()
    }
    f3 := Foo{Arch: "amd64"}
    err = vd.Validate(f3)
    if err == nil {
        t.Fail()
    }
}
