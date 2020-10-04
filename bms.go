package main

type BMS struct {
	connection *Connect
}

type cell struct {
	Number int
	AverageVoltage float32
	MinimumVoltage float32
	MaximumVoltage float32
	DeltaVoltage  float32
	AverageDeviation float32
	MinimumDeviation float32
	MaximumDeviation float32
	DeltaDeviation float32
}

// ShowStats contains all output of the 'show stats'
type ShowStats struct {
	MeanCellVoltage float32
	StandardDeviation float32
	Cells             []cell
}



