package schema

import (
	"bytes"
	"fmt"
)

type cursor struct {
	line   int
	column int
}

type Field struct {
	cursor
	Name         string
	FieldType    string
	IsArray      bool
	IsDepricated bool
	Value        *uint32
}

func (f Field) String() string {
	buf := bytes.NewBuffer(nil)
	if f.FieldType != "" {
		buf.WriteString(f.FieldType)
		if f.IsArray {
			buf.WriteString("[]")
		}
		buf.WriteString(" ")
	}
	buf.WriteString(f.Name)
	if f.Value != nil {
		fmt.Fprintf(buf, " = %d", *f.Value)
		if f.IsDepricated {
			buf.WriteString(" [deprecated]")
		}
	}
	buf.WriteString(";")
	return buf.String()
}

type DefinitionKind int

const (
	KIND_ENUM DefinitionKind = iota
	KIND_STRUCT
	KIND_MESSAGE
)

var definitionKind_v2s = []string{"enum", "struct", "message"}

func (dk DefinitionKind) String() string {
	return definitionKind_v2s[int(dk)]
}

type Definition struct {
	cursor
	Name   string
	Kind   DefinitionKind
	Fields []*Field
}

func (d Definition) String() string {
	buf := bytes.NewBuffer(nil)
	fmt.Fprintf(buf, "%s %s {\n", d.Kind, d.Name)
	for _, f := range d.Fields {
		buf.WriteString("\t")
		buf.WriteString(f.String())
		buf.WriteString("\n")
	}
	/*switch d.Kind {
	case KIND_MESSAGE:
		for _, f := range d.Fields {
			arraySuffix := ""
			if f.IsArray {
				arraySuffix = "[]"
			}
			deprecatedSuffix := ""
			if f.IsDepricated {
				deprecatedSuffix = " [deprecated]"
			}
			str += fmt.Sprintf("\t%s%s %s = %d%s;\n", f.FieldType, arraySuffix, f.Name, f.Value, deprecatedSuffix)
		}
	case KIND_STRUCT:
		for _, f := range d.Fields {
			arraySuffix := ""
			if f.IsArray {
				arraySuffix = "[]"
			}
			str += fmt.Sprintf("\t%s%s %s;\n", f.FieldType, arraySuffix, f.Name)
		}
	case KIND_ENUM:
		for _, f := range d.Fields {
			str += fmt.Sprintf("\t%s = %d;\n", f.Name, f.Value)
		}
	}*/
	buf.WriteString("}")
	return buf.String()
}

type Schema struct {
	PackageName string
	Definitions []*Definition
}

func (s Schema) String() string {
	buf := bytes.NewBuffer(nil)
	if s.PackageName != "" {
		buf.WriteString("package ")
		buf.WriteString(s.PackageName)
		buf.WriteString(";\n\n")
	}
	for _, d := range s.Definitions {
		buf.WriteString(d.String())
		buf.WriteString("\n\n")
	}
	return buf.String()
}
