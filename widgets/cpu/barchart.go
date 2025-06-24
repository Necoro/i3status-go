package cpu

type barchart struct {
	chars   []rune
	binSize float64
}

var barchars = newBarchart([]rune{'▁', '▂', '▃', '▄', '▅', '▆', '▇', '█'})

func (b barchart) Bar(val float64) rune {
	bin := int(val / b.binSize)
	return b.chars[bin]
}

func newBarchart(chars []rune) barchart {
	return barchart{chars: chars, binSize: 100. / float64(len(chars))}
}
