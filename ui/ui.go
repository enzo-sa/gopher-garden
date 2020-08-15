package ui

import (
	"gioui.org/app"
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/key"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"github.com/enzo-sa/gopher-garden/engine"
	"image"
	"image/color"
	"path/filepath"
	"runtime"
	"time"
)

// main ui is handled here

var topMenuPx float32 = engine.LawnLength * 10
// HeightPx represents game display height in pixels (varies as window size changes)
var HeightPx float32 = engine.LawnLength*100 + topMenuPx
// WidthPx represents game display width in pixels (varies as window size changes)
var WidthPx float32 = engine.LawnLength * 100

func getPath() string {
	_, currentFilePath, _, _ := runtime.Caller(0)
	return filepath.Dir(currentFilePath)
}

type ind int

// decoded png images are loaded here for drawing
var path = getPath()
var grassTxtrImg = paint.NewImageOp(getImg(path + "/resources/grass-texture.png"))
var holeImg = paint.NewImageOp(getImg(path + "/resources/hole.png"))
var snakeUImg = paint.NewImageOp(getImg(path + "/resources/snake-u.png"))
var snakeDImg = paint.NewImageOp(getImg(path + "/resources/snake-d.png"))
var snakeRImg = paint.NewImageOp(getImg(path + "/resources/snake-r.png"))
var snakeLImg = paint.NewImageOp(getImg(path + "/resources/snake-l.png"))
var gopherUImg = paint.NewImageOp(getImg(path + "/resources/gopher-u.png"))
var gopherDImg = paint.NewImageOp(getImg(path + "/resources/gopher-d.png"))
var gopherRImg = paint.NewImageOp(getImg(path + "/resources/gopher-r.png"))
var gopherLImg = paint.NewImageOp(getImg(path + "/resources/gopher-l.png"))
var gopherDeadImg = paint.NewImageOp(getImg(path + "/resources/gopher-dead.png"))
var carrotImg = paint.NewImageOp(getImg(path + "/resources/carrot.png"))
var carrot2Img = paint.NewImageOp(getImg(path + "/resources/carrot2.png"))
var grassBkgImg = paint.NewImageOp(getImg(path + "/resources/grass-background.png"))
var titleImg = paint.NewImageOp(getImg(path + "/resources/title-screen.png"))

// Ui struct contains main elements of the ui
type Ui struct {
	w                  *app.Window
	gtx                layout.Context
	th                 *material.Theme
	ga                 *engine.Garden
	name               string
	menuBtn           btn
	continueBtn       btn
	highscoresBtn     btn
	backhighscoresBtn btn
	newgameBtn        btn
	backmenuBtn       btn
	exitBtn           btn
	nameEditor        *widget.Editor
	titleScreen       bool
}

// NewUi returns an initialized Ui struct
func NewUi(w *app.Window) *Ui {
	u := Ui{
		w:  w,
		th: material.NewTheme(gofont.Collection()),
		ga: engine.NewGame(),
	}
	u.th.TextSize = unit.Dp(topMenuPx / 5)
	u.ga.ScaleOffset(WidthPx)
	u.nameEditor = &widget.Editor{
		SingleLine: true,
		Submit:     true,
	}
	u.menuBtn.pressed = true
	u.titleScreen = true
	return &u
}

// draws grass block given lawn ind
func (i ind) grass(gtx layout.Context, lawn *[engine.LawnArea]engine.Grass, dead bool) {
	g := lawn[i]
	draw := func() {
		// scale square size based on width and height values
		paint.PaintOp{Rect: f32.Rect(0+g.Off.X, 0+g.Off.Y+topMenuPx,
			((WidthPx)/engine.LawnLength)+g.Off.X, ((HeightPx-topMenuPx)/engine.LawnLength)+g.Off.Y+topMenuPx)}.Add(gtx.Ops)
	}
	defer op.Push(gtx.Ops).Pop()
	// draws base grass texture
	paint.ImageOp(grassTxtrImg).Add(gtx.Ops)
	draw()
	// conditionally draws other features
	if g.Hole {
		paint.ImageOp(holeImg).Add(gtx.Ops)
		draw()
	}
	if g.Carrot {
		paint.ImageOp(carrotImg).Add(gtx.Ops)
		draw()
	}
	if g.Player.Has {
		if dead {
			paint.ImageOp(gopherDeadImg).Add(gtx.Ops)
		} else {
			playerDirec := map[int]paint.ImageOp{0: gopherUImg, 1: gopherDImg, 2: gopherRImg, 3: gopherLImg}
			paint.ImageOp(playerDirec[g.Player.Direc]).Add(gtx.Ops)
		}
		draw()
	}
	if g.Snake.Has {
		snakeDirec := map[int]paint.ImageOp{0: snakeUImg, 1: snakeDImg, 2: snakeRImg, 3: snakeLImg}
		paint.ImageOp(snakeDirec[g.Snake.Direc]).Add(gtx.Ops)
		draw()
	}
}

// draws all grass
func (u *Ui) full() {
	cs := u.gtx.Constraints
	u.gtx.Constraints.Min = image.Point{}
	// center lawn as biggest possible square in window and define transformation for
	// later ops to be in center
	defer op.Push(u.gtx.Ops).Pop()
	paint.ImageOp(grassBkgImg).Add(u.gtx.Ops)
	paint.PaintOp{Rect: f32.Rect(0, 0, float32(cs.Min.X), float32(cs.Min.Y))}.Add(u.gtx.Ops)
	// set both values to the smallest of both (because the garden is always a square) (plus the menu offset for height)
	if float32(cs.Min.X) != WidthPx || float32(cs.Min.Y) != HeightPx {
		if cs.Min.X < cs.Min.Y-int(topMenuPx) {
			topMenuPx = float32(cs.Min.X) / 10
			WidthPx = float32(cs.Min.X)
			HeightPx = float32(cs.Min.X)
			HeightPx += topMenuPx
			op.Offset(f32.Pt(0, float32(((cs.Min.Y-int(topMenuPx))-cs.Min.X)/2))).Add(u.gtx.Ops)

		} else {
			topMenuPx = float32(cs.Min.Y) / 10
			WidthPx = float32(cs.Min.Y) - topMenuPx
			HeightPx = float32(cs.Min.Y)
			op.Offset(f32.Pt(float32((cs.Min.X-(cs.Min.Y-int(topMenuPx)))/2), 0)).Add(u.gtx.Ops)
		}
		u.th.TextSize = unit.Dp(topMenuPx / 5)
		u.ga.ScaleOffset(WidthPx)
	}
	u.topMenu()
	for i := 0; i < engine.LawnArea; i++ {
		ind(i).grass(u.gtx, u.ga.Lawn, u.ga.Dead)
	}
	if u.ga.Dead && !u.menuBtn.pressed {
		u.gameOver()
	}
	if u.menuBtn.pressed {
		paint.ColorOp{Color: color.RGBA{0, 0, 0, 0xdf}}.Add(u.gtx.Ops)
		paint.PaintOp{Rect: f32.Rect(0, 0, WidthPx, HeightPx)}.Add(u.gtx.Ops)
		if u.titleScreen {
			paint.ImageOp(titleImg).Add(u.gtx.Ops)
			paint.PaintOp{Rect: f32.Rect(0, topMenuPx, WidthPx, HeightPx)}.Add(u.gtx.Ops)
		}
		if !u.newgameBtn.pressed {
			u.mainMenu()
		}
	}
	if u.highscoresBtn.pressed {
		u.drawHs()
	}
	// get name when newgame button is pressed
	if u.newgameBtn.pressed && !u.highscoresBtn.pressed {
		u.titleScreen = false
		u.getName()
	}
}

// Loop is the main ui loop which constantly draws frames and handles input
func (u *Ui) Loop() error {
	// snake ticker
	// channel dec is sent every 5 carrots collected and it decreases interval time
	startInterval := float64(1500)
	dec := make(chan bool)
	reset := make(chan bool)
	ticker := time.NewTicker(time.Duration(startInterval) * time.Millisecond)
	go func() {
		decCounter := 1.0
		for {
			select {
			// accelerate interval frequency / decrease interval time
			case <-dec:
				ticker.Stop()
				decCounter += 0.25
				ticker = time.NewTicker(time.Duration(startInterval/decCounter) * time.Millisecond)
			// send reset every new game. reset sets decCounter back to start
			case <-reset:
				ticker.Stop()
				decCounter = 1.0
				ticker = time.NewTicker(time.Duration(startInterval/decCounter) * time.Millisecond)
			}
		}
	}()
	// main loop
	var ops op.Ops
	var tickerSwitch bool
	for {
		select {
		case e := <-u.w.Events():
			switch e := e.(type) {
			case system.DestroyEvent:
				return e.Err
			case system.FrameEvent:
				u.gtx = layout.NewContext(&ops, e)
				// had the option to only update the squares that were changed,
				// as returned by Update(), but have to redraw everything anyways
				// between frame events and caching an image of the lawn and updating
				// upon that image would end up being less efficient than simply
				// redrawing all of the squares each time
				if u.ga.Dead {
					if u.continueBtn.pressed {
						u.menuBtn.pressed = true
						u.continueBtn.pressed = false
						// locally save name and score data when continue button pressed
						u.updateHs()
						// reset name data
						u.name = ""
					}
				}
				if u.backhighscoresBtn.pressed {
					u.highscoresBtn.pressed = false
					u.backhighscoresBtn.pressed = false
				}
				if u.menuBtn.pressed && !u.highscoresBtn.pressed {
					if u.backmenuBtn.pressed {
						if !u.titleScreen {
							u.menuBtn.pressed = false
						}
						u.backmenuBtn.pressed = false
					}
					if u.newgameBtn.pressed {
						// dont start newgame until name is entered
						if u.name != "" {
							u.ga = engine.NewGame()
							// reset snake ticker on newgame
							reset <- true
							u.menuBtn.pressed = false
							u.newgameBtn.pressed = false
						}
					}
					if u.exitBtn.pressed {
						print("Exiting.\n")
						return nil
					}
				} else {
					u.ga.Update()
				}
				u.full()
				e.Frame(u.gtx.Ops)
				if u.ga.Score%5 != 0 {
					tickerSwitch = true
				}
				// update ticker if 5th carrot eaten
				if u.ga.Score%5 == 0 && u.ga.Score != 0 && tickerSwitch {
					dec <- true
					tickerSwitch = false
				}
			case key.Event:
				if !u.ga.Dead && !u.menuBtn.pressed {
					if ok := u.ga.HandleKey(e.Name); ok {
						u.w.Invalidate()
					}
				}
			case system.ClipboardEvent:
				u.nameEditor.SetText(e.Text)
			}
		case <-ticker.C:
			if !u.menuBtn.pressed {
				u.ga.MoveSnakes()
				u.w.Invalidate()
			}
		}
	}
}
