package evmmodel

var (
	Tables []interface{}
)

func init() {
	Tables = append(Tables,
		new(GasOutputs),
	)
}
