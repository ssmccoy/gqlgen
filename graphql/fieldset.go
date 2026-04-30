package graphql

import (
	"context"
	"io"
)

type FieldSet struct {
	fields   []CollectedField
	Values   []Marshaler
	Invalids uint32
	delayed  []delayedResult
}

type delayedResult struct {
	i int
	f func(context.Context) Marshaler
}

func NewFieldSet(fields []CollectedField) *FieldSet {
	return &FieldSet{
		fields: fields,
		Values: make([]Marshaler, len(fields)),
	}
}

func (m *FieldSet) AddField(field CollectedField) {
	m.fields = append(m.fields, field)
	m.Values = append(m.Values, nil)
}

func (m *FieldSet) Concurrently(i int, f func(context.Context) Marshaler) {
	m.delayed = append(m.delayed, delayedResult{i: i, f: f})
}

func (m *FieldSet) Dispatch(ctx context.Context) {
	for _, d := range m.delayed {
		m.Values[d.i] = d.f(ctx)
	}
}

func (m *FieldSet) MarshalGQL(writer io.Writer) {
	writer.Write(openBrace)
	writtenFields := make(map[string]bool, len(m.fields))
	for i, field := range m.fields {
		if writtenFields[field.Alias] {
			continue
		}
		if i != 0 {
			writer.Write(comma)
		}
		writeQuotedString(writer, field.Alias)
		writer.Write(colon)
		m.Values[i].MarshalGQL(writer)
		writtenFields[field.Alias] = true
	}
	writer.Write(closeBrace)
}
