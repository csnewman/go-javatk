package dextk

import (
	"errors"
	"fmt"
	"strings"
)

var (
	ErrEmptyTypeDesc = errors.New("empty type desc")
	ErrBadTypeDesc   = errors.New("malformed type desc")
)

func (r *Reader) ReadTypeAndParse(id uint32) (TypeDescriptor, error) {
	var res TypeDescriptor

	typeDef, err := r.ReadType(id)
	if err != nil {
		return res, err
	}

	typeDesc, err := r.ReadString(typeDef.DescriptorStringID)
	if err != nil {
		return res, err
	}

	return ParseTypeDescriptor(typeDesc)
}

type TypeDescriptor struct {
	Type        uint8
	ArrayLength int
	ClassName   string
}

func ParseTypeDescriptor(value string) (TypeDescriptor, error) {
	var res TypeDescriptor

	l := len(value)

	if l == 0 {
		return res, ErrEmptyTypeDesc
	}

	for value[res.ArrayLength] == '[' {
		res.ArrayLength++

		// Ensure there is a next character
		if res.ArrayLength == l {
			return res, ErrBadTypeDesc
		}
	}

	res.Type = value[res.ArrayLength]

	// Check if a string
	if res.Type != 'L' {
		// Only a single character should appear if not a string
		if res.ArrayLength+1 != l {
			return res, ErrBadTypeDesc
		}

		return res, nil
	}

	if value[l-1] != ';' {
		return res, ErrBadTypeDesc
	}

	if l-2-res.ArrayLength <= 0 {
		return res, ErrBadTypeDesc
	}

	res.ClassName = value[1+res.ArrayLength : l-1]

	return res, nil
}

func (d TypeDescriptor) String() string {
	ap := strings.Repeat("[", d.ArrayLength)

	if d.Type == 'L' {
		return fmt.Sprintf("%vL%s;", ap, d.ClassName)
	}

	return fmt.Sprintf("%v%c", ap, d.Type)
}

func (d TypeDescriptor) Base() TypeDescriptor {
	return TypeDescriptor{
		Type:        d.Type,
		ArrayLength: 0,
		ClassName:   d.ClassName,
	}
}

func (d TypeDescriptor) IsArray() bool {
	return d.ArrayLength != 0
}

func (d TypeDescriptor) IsClass() bool {
	return d.Type == 'L'
}