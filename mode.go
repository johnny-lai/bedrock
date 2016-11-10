package bedrock

const (
	DebugMode   = iota // c0 == 0
	ReleaseMode = iota // c1 == 1
)

var bedrockMode = DebugMode

func SetMode(mode int) {
	bedrockMode = mode
}

func Mode() int {
	return bedrockMode
}
