package game

import (
	"Unbewohnte/capyclick/conf"
	"Unbewohnte/capyclick/logger"
	"Unbewohnte/capyclick/resources"
	"Unbewohnte/capyclick/save"
	"Unbewohnte/capyclick/util"
	"fmt"
	"image/color"
	"path/filepath"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
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
	FontFace            font.Face
	AnimationData       AnimationData
	PassiveIncomeTicker int
}

func NewGame() *Game {
	audioCtx := audio.NewContext(44000)
	fnt := resources.GetFont("PixeloidSans-Bold.otf")

	return &Game{
		WorkingDir:   ".",
		Config:       conf.Default(),
		Save:         save.Default(),
		AudioContext: audioCtx,
		AudioPlayers: map[string]*audio.Player{
			"boop":        resources.GetAudioPlayer(audioCtx, "boop.wav"),
			"woop":        resources.GetAudioPlayer(audioCtx, "woop.wav"),
			"menu_switch": resources.GetAudioPlayer(audioCtx, "menu_switch.wav"),
			"levelup":     resources.GetAudioPlayer(audioCtx, "levelup.wav"),
		},
		ImageResources: map[string]*ebiten.Image{
			"capybara1":   ebiten.NewImageFromImage(resources.ImageFromFile("capybara_1.png")),
			"capybara2":   ebiten.NewImageFromImage(resources.ImageFromFile("capybara_2.png")),
			"capybara3":   ebiten.NewImageFromImage(resources.ImageFromFile("capybara_3.png")),
			"background1": ebiten.NewImageFromImage(resources.ImageFromFile("background_1.png")),
			"background2": ebiten.NewImageFromImage(resources.ImageFromFile("background_2.png")),
		},
		FontFace: util.NewFace(fnt, &opentype.FaceOptions{
			Size:    32,
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
func (g *Game) SaveData(saveFileName string, configurationFileName string) error {
	// Save configuration information and game data
	err := save.Create(filepath.Join(g.WorkingDir, saveFileName), g.Save)
	if err != nil {
		logger.Error("[SaveData] Failed to save game data before closing: %s!", err)
		return err
	}

	err = conf.Create(filepath.Join(g.WorkingDir, configurationFileName), g.Config)
	if err != nil {
		logger.Error("[SaveData] Failed to save game configuration before closing: %s!", err)
		return err
	}

	return nil
}

// Returns how many points required to be considered of level
func pointsForLevel(level uint32) uint64 {
	return 25 * uint64(level*level)
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

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) != 0 {
		// Click!
		g.Save.TimesClicked++
		g.Save.Points++
		g.AnimationData.Squish += 0.5
		g.PlaySound("woop")
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

	if g.Save.Points > 0 && g.Save.Points >= pointsForLevel(g.Save.Level+1) {
		// Level progression
		g.Save.Level++
		g.Save.PassiveIncome++
		g.PlaySound("levelup")
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

	op = &ebiten.DrawImageOptions{}
	if g.AnimationData.BounceDirectionFlag {
		g.AnimationData.Theta += 0.001
	} else {
		g.AnimationData.Theta -= 0.001
	}

	capybaraBounds := g.ImageResources[capybaraKey].Bounds()
	scale := float64(screen.Bounds().Dx()) / float64(capybaraBounds.Dx()) / 2.5
	op.GeoM.Scale(
		scale+g.AnimationData.Squish,
		scale-g.AnimationData.Squish,
	)
	op.GeoM.Rotate(g.AnimationData.Theta)

	capyWidth := float64(g.ImageResources[capybaraKey].Bounds().Dx()) * scale
	capyHeight := float64(g.ImageResources[capybaraKey].Bounds().Dy()) * scale
	op.GeoM.Translate(
		float64(screen.Bounds().Dx()/2)-capyWidth/2,
		float64(screen.Bounds().Dy()/2)-capyHeight/2,
	)

	screen.DrawImage(g.ImageResources[capybaraKey], op)

	// Points
	msg := fmt.Sprintf("Points: %d", g.Save.Points)
	text.Draw(
		screen,
		msg,
		g.FontFace,
		10,
		g.FontFace.Metrics().Height.Ceil(),
		color.White,
	)

	// Level
	msg = fmt.Sprintf(
		"Level: %d (+%d)",
		g.Save.Level,
		pointsForLevel(g.Save.Level+1)-g.Save.Points,
	)
	text.Draw(
		screen,
		msg,
		g.FontFace,
		10,
		g.FontFace.Metrics().Height.Ceil()*2,
		color.White,
	)

	// Times Clicked
	msg = fmt.Sprintf("Clicks: %d", g.Save.TimesClicked)
	text.Draw(
		screen,
		msg,
		g.FontFace,
		10,
		screen.Bounds().Dy()-g.FontFace.Metrics().Height.Ceil()*2,
		color.White,
	)

	// Volume
	msg = fmt.Sprintf("Volume: %d%% (← or →)", int(g.Config.Volume*100.0))
	text.Draw(
		screen,
		msg,
		g.FontFace,
		10,
		screen.Bounds().Dy()-g.FontFace.Metrics().Height.Ceil(),
		color.White,
	)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	scaleFactor := ebiten.DeviceScaleFactor()
	return int(float64(outsideWidth) * scaleFactor), int(float64(outsideHeight) * scaleFactor)
}
