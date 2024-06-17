package game

// Returns how many points required to be considered of level
func pointsForLevel(level uint32) uint64 {
	return 25 * uint64(level*level)
}
