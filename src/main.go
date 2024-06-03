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

package main

import (
	"Unbewohnte/capyclick/conf"
	"Unbewohnte/capyclick/game"
	"Unbewohnte/capyclick/logger"
	"Unbewohnte/capyclick/resources"
	"Unbewohnte/capyclick/save"
	"Unbewohnte/capyclick/util"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
)

const Version string = "v0.1"

var (
	silent    *bool = flag.Bool("silent", false, "Set to true in order to discard all logging")
	version   *bool = flag.Bool("version", false, "Prints version information")
	saveFiles *bool = flag.Bool("saveFiles", false, "Run the game with configuration and save files")
)

const (
	ConfigurationFileName string = "capyclickConfig.json"
	SaveFileName          string = "capyclickSave.json"
)

func main() {
	// Set logging output
	logger.SetOutput(os.Stdout)

	// Parse flags
	flag.Parse()

	if *version {
		fmt.Fprintf(os.Stdout, "Capyclick %s\n(c) 2024 Kasianov Nikolai Alexeevich (Unbewohnte)\n", Version)
		os.Exit(0)
	}

	if *silent {
		// Do not output logs
		logger.SetOutput(io.Discard)
	}

	// Create a game instance
	var game *game.Game = game.NewGame()

	if *saveFiles {
		// Work out working directory
		exeDir, err := os.Executable()
		if err != nil {
			logger.Error("[Init] Failed to get executable's path: %s", err)
			os.Exit(1)
		}
		game.WorkingDir = filepath.Dir(exeDir)
	} else {
		game.WorkingDir = ""
	}

	if *saveFiles {
		// Open/Create configuration file
		var config *conf.Configuration
		config, err := conf.FromFile(filepath.Join(game.WorkingDir, ConfigurationFileName))
		if err != nil {
			err = conf.Create(filepath.Join(game.WorkingDir, ConfigurationFileName), game.Config)
			if err != nil {
				logger.Error("[Init] Failed to create a new configuration file: %s", err)
				os.Exit(1)
			}
			logger.Info("[Init] Created a new configuration file")
			// Proceed with a newly created configuration file
		}

		// Replace default config with an opened one (if exists)
		if config != nil {
			game.Config = *config
		}
	}

	// Set up window options
	ebiten.SetWindowIcon(util.GenerateIcons(resources.ImageFromFile("capybara_2.png"), [][2]uint{
		{32, 32},
	}))
	ebiten.SetWindowClosingHandled(true) // So we can save data
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSizeLimits(512, 576, -1, -1)
	ebiten.SetTPS(60)
	ebiten.SetWindowSize(game.Config.WindowSize[0], game.Config.WindowSize[1])
	ebiten.SetWindowPosition(game.Config.LastWindowPosition[0], game.Config.LastWindowPosition[1])
	ebiten.SetWindowTitle(fmt.Sprintf("Capyclick %s", Version))

	if *saveFiles {
		// Open/Create save file
		gameSave, err := save.FromFile(filepath.Join(game.WorkingDir, SaveFileName))
		if err != nil {
			err = save.Create(filepath.Join(game.WorkingDir, SaveFileName), game.Save)
			if err != nil {
				logger.Error("[Init] Failed to create a new save file \"%s\": %s", SaveFileName, err)
				os.Exit(1)
			}
			logger.Info("[Init] Created a new save file \"%s\"", SaveFileName)
			// Proceed with a new save file
		}

		// Replace a blank save with an existing one (if exists)
		if gameSave != nil {
			gameSave.LastOpenedUnix = uint64(time.Now().Unix())
			game.Save = *gameSave
		}
	}

	// Set each player's volume to the saved value
	for _, player := range game.AudioPlayers {
		player.SetVolume(game.Config.Volume)
	}

	// Run the game
	err := ebiten.RunGame(game)
	if err == ebiten.Termination || err == nil {
		logger.Info("[Main] Shutting down!")
		if *saveFiles {
			game.SaveData(SaveFileName, ConfigurationFileName)
		}
		os.Exit(0)
	} else {
		logger.Error("[Main] Fatal game error: %s", err)
		os.Exit(1)
	}
}
