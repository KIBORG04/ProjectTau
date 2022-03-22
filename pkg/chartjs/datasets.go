package chartjs

type shape string

const (
	Circle      shape = "circle"
	Cross       shape = "cross"
	Crossrot    shape = "crossrot"
	Dash        shape = "dash"
	Line        shape = "line"
	Rect        shape = "rect"
	Rectrounded shape = "rectrounded"
	Rectrot     shape = "rectrot"
	Star        shape = "star"
	Triangle    shape = "triangle"
)

type Coords map[string]interface{}

type Dataset struct {
	Label string   `json:"label,omitempty"`
	Data  []Coords `json:"data"`

	BackgroundColor        []*RGBA `json:"backgroundColor,omitempty"`
	BorderColor            []*RGBA `json:"borderColor,omitempty"`
	BorderWidth            float64 `json:"borderWidth,omitempty"`
	Fill                   string  `json:"fill,omitempty"`
	SteppedLine            string  `json:"steppedLine,omitempty"`
	LineTension            float64 `json:"lineTension,omitempty"`
	CubicInterpolationMode string  `json:"cubicInterpolationMode,omitempty"`
	PointBackgroundColor   RGBA    `json:"pointBackgroundColor,omitempty"`
	PointBorderColor       RGBA    `json:"pointBorderColor,omitempty"`
	PointBorderWidth       float64 `json:"pointBorderWidth,omitempty"`
	PointRadius            float64 `json:"pointRadius,omitempty"`
	PointHitRadius         float64 `json:"pointHitRadius,omitempty"`
	PointHoverRadius       float64 `json:"pointHoverRadius,omitempty"`
	PointHoverBorderColor  RGBA    `json:"pointHoverBorderColor,omitempty"`
	PointHoverBorderWidth  float64 `json:"pointHoverBorderWidth,omitempty"`
	PointStyle             shape   `json:"pointStyle,omitempty"`
	ShowLine               bool    `json:"showLine,omitempty"`
	SpanGaps               bool    `json:"spanGaps,omitempty"`
}

func BarDataset(label string, coords []Coords) *Dataset {
	backColors := RandomColor(len(coords))
	return &Dataset{
		Label:           label,
		Data:            coords,
		BackgroundColor: backColors,
		BorderColor:     RandomBorder(backColors),
		BorderWidth:     3,
	}
}
