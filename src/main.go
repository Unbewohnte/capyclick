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
	"strings"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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

type AnimationData struct {
	Squish              float64
	Theta               float64
	BounceDirectionFlag bool
}

type Game struct {
	WorkingDir          string
	Config              conf.Configuration
	Save                save.Save
	AudioContext        *audio.Context
	AudioPlayers        map[string]*audio.Player
	ImageResources      map[string]*ebiten.Image
	Font                font.Face
	AnimationData       AnimationData
	PassiveIncomeTicker int
}

func NewGame() *Game {
	audioCtx := audio.NewContext(44000)
	fnt := ResourceGetFont("PixeloidSans-Bold.otf")

	return &Game{
		WorkingDir:   ".",
		Config:       conf.Default(),
		Save:         save.Default(),
		AudioContext: audioCtx,
		AudioPlayers: map[string]*audio.Player{
			"boop":        GetAudioPlayer(audioCtx, "boop.wav"),
			"woop":        GetAudioPlayer(audioCtx, "woop.wav"),
			"menu_switch": GetAudioPlayer(audioCtx, "menu_switch.wav"),
			"levelup":     GetAudioPlayer(audioCtx, "levelup.wav"),
		},
		ImageResources: map[string]*ebiten.Image{
			"capybara1":   ebiten.NewImageFromImage(ImageFromFile("capybara_1.png")),
			"capybara2":   ebiten.NewImageFromImage(ImageFromFile("capybara_2.png")),
			"capybara3":   ebiten.NewImageFromImage(ImageFromFile("capybara_3.png")),
			"background1": ebiten.NewImageFromImage(ImageFromFile("background_1.png")),
			"background2": ebiten.NewImageFromImage(ImageFromFile("background_2.png")),
		},
		Font: util.NewFont(fnt, &opentype.FaceOptions{
			Size:    24,
			DPI:     72,
			Hinting: font.HintingVertical,
		}),
		AnimationData: AnimationData{
			Theta:               0.0,
			BounceDirectionFlag: true,
			Squish:              0,
		},
		PassiveIncomeTicker: 0,
	}
}

// Plays sound and rewinds the player
func (g *Game) PlaySound(soundKey string) {
	if strings.TrimSpace(soundKey) != "" {
		g.AudioPlayers[soundKey].Rewind()
		g.AudioPlayers[soundKey].Play()
	}
}

// Saves configuration information and game data
func (g *Game) SaveData() error {
	// Save configuration information and game data
	err := save.Create(filepath.Join(g.WorkingDir, SaveFileName), g.Save)
	if err != nil {
		logger.Error("[SaveData] Failed to save game data before closing: %s!", err)
		return err
	}

	err = conf.Create(filepath.Join(g.WorkingDir, ConfigurationFileName), g.Config)
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

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		// Decrease volume
		if g.Config.Volume-0.2 >= 0 {
			g.Config.Volume -= 0.2
			for _, player := range g.AudioPlayers {
				player.SetVolume(g.Config.Volume)
			}
		}
	} else if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		// Increase volume
		if g.Config.Volume+0.2 <= 1.0 {
			g.Config.Volume += 0.2
			for _, player := range g.AudioPlayers {
				player.SetVolume(g.Config.Volume)
			}
		}
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		// Click!
		g.Save.TimesClicked++
		g.Save.Points++
		g.AnimationData.Squish += 0.5
		g.PlaySound("boop")
	}

	if g.Save.Points > 0 && g.Save.Points%100 == 0 {
		// Level progression
		g.Save.Level++
		g.Save.PassiveIncome++
		// Bump points number immediately in order to not mistakengly level up on the next update
		g.Save.Points++
		g.PlaySound("levelup")
	}

	// Capybara Animation
	if g.AnimationData.Theta >= 0.03 {
		g.AnimationData.BounceDirectionFlag = false
	} else if g.AnimationData.Theta <= -0.03 {
		g.AnimationData.BounceDirectionFlag = true
	}

	if g.AnimationData.Squish >= 0 {
		g.AnimationData.Squish -= 0.05
	}

	// Passive points income
	if g.PassiveIncomeTicker == ebiten.TPS() {
		g.PassiveIncomeTicker = 0
		g.Save.Points += g.Save.PassiveIncome
	} else {
		g.PassiveIncomeTicker++
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.Black)

	backBounds := g.ImageResources["background1"].Bounds()
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(
		float64(screen.Bounds().Dx())/float64(backBounds.Dx()),
		float64(screen.Bounds().Dy())/float64(backBounds.Dy()),
	)
	screen.DrawImage(g.ImageResources["background1"], op)

	// Capybara
	var capybaraKey string
	switch g.Save.Level {
	case 1:
		capybaraKey = "capybara1"
	case 2:
		capybaraKey = "capybara2"
	case 3:
		capybaraKey = "capybara3"
	default:
		capybaraKey = "capybara3"
	}

	scale := 10.0
	op = &ebiten.DrawImageOptions{}
	if g.AnimationData.BounceDirectionFlag {
		g.AnimationData.Theta += 0.001
	} else {
		g.AnimationData.Theta -= 0.001
	}

	op.GeoM.Scale(scale+g.AnimationData.Squish, scale-g.AnimationData.Squish)

	op.GeoM.Rotate(g.AnimationData.Theta)

	width := g.ImageResources[capybaraKey].Bounds().Dx() * int(scale)
	height := g.ImageResources[capybaraKey].Bounds().Dy() * int(scale)
	op.GeoM.Translate(
		float64(screen.Bounds().Dx()/2)-float64(width/2),
		float64(screen.Bounds().Dy()/2)-float64(height/2),
	)
	screen.DrawImage(g.ImageResources[capybaraKey], op)

	// Points
	msg := fmt.Sprintf("Points: %d", g.Save.Points)
	text.Draw(
		screen,
		msg,
		g.Font,
		10,
		30,
		color.White,
	)

	// Level
	msg = fmt.Sprintf("Level: %d", g.Save.Level)
	text.Draw(
		screen,
		msg,
		g.Font,
		screen.Bounds().Dx()-len(msg)*19,
		30,
		color.White,
	)

	// Times Clicked
	msg = fmt.Sprintf("Clicks: %d", g.Save.TimesClicked)
	text.Draw(
		screen,
		msg,
		g.Font,
		10,
		// screen.Bounds().Dy()-30,
		screen.Bounds().Dy()-24*2,
		color.White,
	)

	// Volume
	msg = fmt.Sprintf("Volume: %d%%", int(g.Config.Volume*100.0))
	text.Draw(
		screen,
		msg,
		g.Font,
		// screen.Bounds().Dx()-len(msg)*17,
		10,
		screen.Bounds().Dy()-24,
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
	var game *Game = NewGame()

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
	ebiten.SetWindowIcon(util.GenerateIcons(ImageFromFile("capybara_2.png"), [][2]uint{
		{32, 32},
	}))
	ebiten.SetWindowClosingHandled(true) // So we can save data
	ebiten.SetRunnableOnUnfocused(true)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSizeLimits(380, 576, -1, -1)
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
			game.SaveData()
		}
		os.Exit(0)
	} else {
		logger.Error("[Main] Fatal game error: %s", err)
		os.Exit(1)
	}
}
