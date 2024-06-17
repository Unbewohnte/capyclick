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

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type Game struct {
	WorkingDir          string
	Config              conf.Configuration
	Save                save.Save
	AudioPlayers        map[string]*audio.Player
	FontFace            font.Face
	PassiveIncomeTicker int
	Capybara            *Capybara
	Background          *Sprite
	MandarinRain        *MandarinRain
}

func NewGame() Game {
	audioCtx := audio.NewContext(44000)
	fnt := resources.GetFont("PixeloidSans-Bold.otf")

	return Game{
		WorkingDir: ".",
		Config:     conf.Default(),
		Save:       save.Default(),
		AudioPlayers: map[string]*audio.Player{
			"boop":                    resources.GetAudioPlayer(audioCtx, "boop.wav"),
			"woop":                    resources.GetAudioPlayer(audioCtx, "woop.wav"),
			"menu_switch":             resources.GetAudioPlayer(audioCtx, "menu_switch.wav"),
			"levelup":                 resources.GetAudioPlayer(audioCtx, "levelup.wav"),
			"mandarin_box_full":       resources.GetAudioPlayer(audioCtx, "mandarin_box_full.wav"),
			"orange_put":              resources.GetAudioPlayer(audioCtx, "orange_put.wav"),
			"mandarin_rain_completed": resources.GetAudioPlayer(audioCtx, "mandarin_rain_completed.wav"),
		},
		Capybara:   NewCapybara(NewSpriteFromFile("capybara_1.png")),
		Background: NewSpriteFromFile("background_1.png"),
		FontFace: util.NewFace(fnt, &opentype.FaceOptions{
			Size:    32,
			DPI:     72,
			Hinting: font.HintingVertical,
		}),
		PassiveIncomeTicker: 0,
		MandarinRain:        NewMandarinRain(3, 8),
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

func (g *Game) Update() error {
	if ebiten.IsWindowBeingClosed() {
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Exit
		return ebiten.Termination
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyF12) {
		g.ToggleFullscreen()
	}

	g.SaveWindowGeometry()

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
		// Decrease volume
		g.DecreaseVolume(0.2)
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
		// Increase volume
		g.IncreaseVolume(0.2)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) ||
		len(inpututil.AppendJustPressedTouchIDs(nil)) != 0 {
		// Click!
		g.Save.TimesClicked++
		g.Save.Points++
		g.PlaySound("woop")
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

	// Capybara animation update
	g.Capybara.Update()

	if !g.MandarinRain.InProgress && g.Save.TimesClicked > 0 && g.Save.TimesClicked%100 == 0 {
		// Have some oranges!
		g.MandarinRain.Run()
		logger.Info("Started mandarin rain at %d points!", g.Save.Points)
	}

	if g.MandarinRain.InProgress {
		// Calculate mandarin rain logic for this step
		g.MandarinRain.Update(g)
	}

	if g.MandarinRain.Completed {
		// Prepare a new mandarin rain
		g.MandarinRain = NewMandarinRain(3, 8)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Background
	screen.Fill(color.Black)

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(
		float64(screen.Bounds().Dx())/float64(g.Background.Img.Bounds().Dx()),
		float64(screen.Bounds().Dy())/float64(g.Background.Img.Bounds().Dy()),
	)
	screen.DrawImage(g.Background.Img, op)

	// Capybara
	g.Capybara.Draw(screen, g.Save.Level)

	// Mandarin rain
	if g.MandarinRain.InProgress {
		g.MandarinRain.Draw(screen)
	}

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
