package main

import (
	"Unbewohnte/capyclick/conf"
	"Unbewohnte/capyclick/logger"
	"Unbewohnte/capyclick/save"
	"Unbewohnte/capyclick/util"
	"flag"
	"fmt"
	"image/color"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const Version string = "v0.1"

var (
	silent  *bool = flag.Bool("silent", false, "Set to true in order to discard all logging")
	version *bool = flag.Bool("version", false, "Prints version information")
	noFiles *bool = flag.Bool("no-files", false, "Run the game without outputting/reading configuration or save files")
)

const (
	ConfigurationFileName string = "capyclickConfig.json"
	SaveFileName          string = "capyclickSave.json"
)

var WorkingDir string = "."

type Game struct {
	Config conf.Configuration
	Save   save.Save
}

// Resource images
var (
	ImgCapybara1 = ebiten.NewImageFromImage(ImageFromFile("capybara_1.png"))
	ImgCapybara2 = ebiten.NewImageFromImage(ImageFromFile("capybara_2.png"))
	ImgCapybara3 = ebiten.NewImageFromImage(ImageFromFile("capybara_3.png"))

	GameFont = util.NewFont(ResourceGetFont("PixeloidSans-Bold.otf"), &opentype.FaceOptions{
		Size:    10,
		DPI:     72,
		Hinting: font.HintingVertical,
	})
)

// Saves configuration information and game data
func SaveData(game *Game) error {
	// Save configuration information and game data
	err := save.Create(filepath.Join(WorkingDir, SaveFileName), game.Save)
	if err != nil {
		logger.Error("[SaveData] Failed to save game data before closing: %s!", err)
		return err
	}

	err = conf.Create(filepath.Join(WorkingDir, ConfigurationFileName), game.Config)
	if err != nil {
		logger.Error("[SaveData] Failed to save game configuration before closing: %s!", err)
		return err
	}

	return nil
}

func (g *Game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		return ebiten.Termination
	}

	// Update configuration and save information
	width, height := ebiten.WindowSize()
	g.Config.WindowSize = [2]int{width, height}

	x, y := ebiten.WindowPosition()
	g.Config.LastWindowPosition = [2]int{x, y}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Exit
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		if ebiten.IsFullscreen() {
			// Turn fullscreen off
			ebiten.SetFullscreen(false)
		} else {
			// Go fullscreen
			ebiten.SetFullscreen(true)
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Click!
		g.Save.TimesClicked++
		g.Save.Points++
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.Black)

	// Capybara
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(15, 15)
	op.GeoM.Translate(float64(screen.Bounds().Dx()/3), float64(screen.Bounds().Dy()/3))
	screen.DrawImage(ImgCapybara1, op)

	// Points
	msg := fmt.Sprintf("Points: %d", g.Save.Points)
	text.Draw(
		screen,
		msg,
		GameFont,
		10,
		20,
		color.White,
	)

	// Level
	msg = fmt.Sprintf("Level: %d", g.Save.Level)
	text.Draw(
		screen,
		msg,
		GameFont,
		screen.Bounds().Dx()-len(msg)*10,
		20,
		color.White,
	)

	// Times Clicked
	msg = fmt.Sprintf("Times Clicked: %d", g.Save.TimesClicked)
	text.Draw(
		screen,
		msg,
		GameFont,
		10,
		screen.Bounds().Dy()-10,
		color.White,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	scaleFactor := ebiten.DeviceScaleFactor()
	return int(float64(outsideWidth) * scaleFactor), int(float64(outsideHeight) * scaleFactor)
}

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
	var game *Game = &Game{
		Config: conf.Default(),
		Save:   save.Default(),
	}

	if !*noFiles {
		// Work out working directory
		exeDir, err := os.Executable()
		if err != nil {
			logger.Error("[Init] Failed to get executable's path: %s", err)
			os.Exit(1)
		}
		WorkingDir = filepath.Dir(exeDir)
	} else {
		WorkingDir = ""
	}

	if !*noFiles {
		// Open/Create configuration file
		var config *conf.Configuration
		config, err := conf.FromFile(filepath.Join(WorkingDir, ConfigurationFileName))
		if err != nil {
			err = conf.Create(filepath.Join(WorkingDir, ConfigurationFileName), game.Config)
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
	ebiten.SetWindowIcon(util.GenerateIcons(ImageFromFile("capybara_2.png"), [][2]uint{
		{32, 32},
	}))
	ebiten.SetWindowClosingHandled(true) // So we can save data
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSizeLimits(480, 320, -1, -1)
	ebiten.SetTPS(60)
	ebiten.SetWindowSize(game.Config.WindowSize[0], game.Config.WindowSize[1])
	ebiten.SetWindowPosition(game.Config.LastWindowPosition[0], game.Config.LastWindowPosition[1])
	ebiten.SetWindowTitle(fmt.Sprintf("Capyclick %s", Version))

	if !*noFiles {
		// Open/Create save file
		gameSave, err := save.FromFile(filepath.Join(WorkingDir, SaveFileName))
		if err != nil {
			err = save.Create(filepath.Join(WorkingDir, SaveFileName), game.Save)
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

	// Run the game
	err := ebiten.RunGame(game)
	if err == ebiten.Termination || err == nil {
		logger.Info("[Main] Shutting down!")
		SaveData(game)
		os.Exit(0)
	} else {
		logger.Error("[Main] Fatal game error: %s", err)
		os.Exit(1)
	}
}