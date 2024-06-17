/*
  	capyclick - Capybara clicker game
    Copyright (C) 2024  Kasianov Nikolai Alekseevich (Unbewohnte)

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

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
