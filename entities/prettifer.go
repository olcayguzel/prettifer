package entities

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strings"
)

type Prettifer struct {
	Settings        Settings
	ExcludedColumns []string
	columns         []*Column
}

func CreateNew() Prettifer {
	return Prettifer{
		Settings: Settings{
			AutoGenerateColumns: true,
			ShowVerticalLines:   false,
			Columns:             []*Column{},
		},
	}
}

func (p *Prettifer) AddExcludedField(f string) {
	if p.ExcludedColumns == nil {
		p.ExcludedColumns = make([]string, 0)
	}
	p.ExcludedColumns = append(p.ExcludedColumns, f)
}

func (p *Prettifer) getColumns(v interface{}) []*Column {
	columns := []*Column{}
	if v != nil {
		t := reflect.TypeOf(v)
		if t.Kind() == reflect.Ptr {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return columns
		}
		if p.ExcludedColumns == nil {
			p.ExcludedColumns = make([]string, 0)
		}
		fieldCount := t.NumField()
		for i := 0; i < fieldCount; i++ {
			field := t.Field(i)
			if stringSliceIndexOf(p.ExcludedColumns, field.Name) >= 0 {
				continue
			}
			column := &Column{}
			column.IsAutoGenerated = true
			column.Order = i
			column.Name = field.Name
			column.Header = getAutoGenerateColumnName(field.Name)
			column.Type = getType(v, field.Name)
			column.Align = getAlignment(column.Type)
			column.Format = getDefaultDataFormat(column.Type)
			column.MaxWidth = 0
			column.MinWidth = 0
			columns = append(columns, column)
		}
	}
	return columns
}

func (p *Prettifer) generateColumns(value interface{}) {
	if p.Settings.AutoGenerateColumns {
		p.columns = append(p.columns, p.getColumns(value)...)
	}
	if len(p.Settings.Columns) > 0 {
		for _, c := range p.Settings.Columns {
			c.IsAutoGenerated = false
			if len(c.Type) == 0 {
				c.Type = getType(value, c.Name)
			}
			if c.Align != 1 && c.Align != 2 && c.Align != 3 {
				c.Align = getAlignment(c.Type)
			}
			if len(c.Format) == 0 {
				c.Format = getDefaultDataFormat(c.Type)
			}
			p.columns = append(p.columns, c)
		}
	}
}

func (p *Prettifer) setColumnWidths() {
	if len(p.columns) > 0 {
		maxLength := 0
		for _, c := range p.columns {
			len := len(c.Header)
			if len > c.MaxWidth && c.MaxWidth > 0 {
				c.Width = c.MaxWidth
			} else if len < c.MinWidth && c.MinWidth > 0 {
				c.Width = c.MinWidth
			} else {
				c.Width = len + 2
			}
			if c.Width > maxLength {
				maxLength = c.Width
			}
		}
		for _, c := range p.columns {
			if maxLength > c.Width && (c.MinWidth == 0 || c.MinWidth < maxLength) && (c.MaxWidth == 0 || c.MaxWidth > maxLength) {
				c.Width = maxLength
			}
		}
	}
}

func (p *Prettifer) printColumns(output *os.File) {
	p.setColumnWidths()
	sort.SliceStable(p.columns, func(i int, j int) bool {
		return p.columns[i].Order < p.columns[j].Order || (p.columns[i].IsAutoGenerated && !p.columns[i].IsAutoGenerated)
	})
	fmt.Println()
	fmt.Print("-")
	for i, c := range p.columns {
		if i == len(p.columns)-1 {
			fmt.Print(strings.Repeat("-", c.Width), "-")
		} else {
			fmt.Print(strings.Repeat("-", c.Width), "+")
		}
	}
	fmt.Println()
	fmt.Print("|")
	for _, c := range p.columns {
		fmt.Print(padding(c.Header, c.Width, AL_CENTER), "|")
	}
	fmt.Println()
	fmt.Print("-")
	for i, c := range p.columns {
		if i == len(p.columns)-1 {
			fmt.Print(strings.Repeat("-", c.Width), "-")
		} else {
			fmt.Print(strings.Repeat("-", c.Width), "+")
		}
	}
	fmt.Println()
}

func (p *Prettifer) ToStdOutput(values ...interface{}) {
	if len(values) > 0 {
		p.generateColumns(values[0])
		p.printColumns(os.Stdout)
		for _, v := range values {
			for _, c := range p.columns {
				value := getValue(v, c.Name)
				fmt.Print(" ", padding(obtainStringValue(value, *c), c.Width, c.Align))
			}
			fmt.Println()
		}
	}
}