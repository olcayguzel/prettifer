package entities

type Settings struct {
	AutoGenerateColumns bool
	ShowVerticalLines   bool
	ShowRowNumbers      bool
	Columns             []*Column
}
