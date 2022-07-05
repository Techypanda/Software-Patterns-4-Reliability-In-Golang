package model

const DEFAULT_CHANCE = 0.0
const DEFAULT_MODE = "NO_PATTERN" // i.e no pattern!
const BACKOFF = "BACKOFF"
const EXPONENTIAL_BACKOFF = "EXPONENTIAL_BACKOFF"
const EXP_WITH_JITTER = "EXP_WITH_JITTER"
const CIRCUIT_BREAKER = "CIRCUIT_BREAKER"
const ALL = "ALL" // no time to do this

type State struct {
	Mode          string
	FailureChance float64
}

func NewState() State {
	return State{Mode: DEFAULT_MODE, FailureChance: DEFAULT_CHANCE}
}
