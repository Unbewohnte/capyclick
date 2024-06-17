package game

import "strings"

// Plays sound and rewinds the player
func (g *Game) PlaySound(soundKey string) {
	if strings.TrimSpace(soundKey) != "" {
		g.AudioPlayers[soundKey].Rewind()
		g.AudioPlayers[soundKey].Play()
	}
}

func (g *Game) SetVolume(volume float64) {
	if volume > 1.0 || volume < 0.0 {
		return
	}

	g.Config.Volume = volume
	for _, player := range g.AudioPlayers {
		player.SetVolume(volume)
	}
}

func (g *Game) IncreaseVolume(volumeDelta float64) {
	for _, player := range g.AudioPlayers {
		volume := player.Volume() + volumeDelta
		if volume > 1.0 || volume < 0.0 {
			continue
		}

		player.SetVolume(volume)
		g.Config.Volume = volume
	}
}

func (g *Game) DecreaseVolume(volumeDelta float64) {
	for _, player := range g.AudioPlayers {
		volume := player.Volume() - volumeDelta
		if volume > 1.0 || volume < 0.0 {
			continue
		}

		player.SetVolume(volume)
		g.Config.Volume = volume
	}
}
