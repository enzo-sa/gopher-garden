package ui

import (
	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"strconv"
)

type btn struct {
	pressed bool
}

// menu drawing occurs here

// this was rushed, so beware of sloppy code with magic values that just work

var menuImg = paint.NewImageOp(getImg(path + "/resources/menu-btn.png"))
var backImg = paint.NewImageOp(getImg(path + "/resources/back-btn.png"))
var newgameImg = paint.NewImageOp(getImg(path + "/resources/newgame-btn.png"))
var highscoresImg = paint.NewImageOp(getImg(path + "/resources/highscores-btn.png"))
var continueImg = paint.NewImageOp(getImg(path + "/resources/continue-btn.png"))
var exitImg = paint.NewImageOp(getImg(path + "/resources/exit-btn.png"))
var gameoverImg = paint.NewImageOp(getImg(path + "/resources/gameover.png"))
var highscoresTitleImg = paint.NewImageOp(getImg(path + "/resources/highscores-title.png"))

// custom button widget layout
func (b *btn) Layout(gtx layout.Context, img paint.ImageOp, hs bool) layout.Dimensions {
	// avoid pollution
	defer op.Push(gtx.Ops).Pop()
	for _, e := range gtx.Events(b) {
		if e, ok := e.(pointer.Event); ok && (img == backImg || !hs) {
			switch e.Type {
			case pointer.Press:
				b.pressed = true // will set to false later
			}
		}
	}
	len := topMenuPx
	// set area
	// set offset for buttons other than menu to be in center of lawn
	if img != menuImg {
		op.Offset(f32.Pt((WidthPx/2)-(len/2), 0)).Add(gtx.Ops)
		pointer.Rect(image.Rect(0, 0, int(len), int(len))).Add(gtx.Ops)
	} else {
		pointer.Rect(image.Rect(0, 0, int(len), int(len))).Add(gtx.Ops)
	}
	pointer.InputOp{
		Tag:   b,
		Types: pointer.Press,
	}.Add(gtx.Ops)
	// draw
	paint.ImageOp(img).Add(gtx.Ops)
	paint.PaintOp{Rect: f32.Rect(0, 0, len, len)}.Add(gtx.Ops)
	return layout.Dimensions{Size: image.Pt(int(len), int(len*1.4))}
}

// text field for getting user's name
func (u *Ui) getName() {
	for _, e := range u.nameEditor.Events() {
		if e, ok := e.(widget.SubmitEvent); ok {
			u.name = e.Text
			u.nameEditor.SetText("")
		}
	}
	if u.name == "" {
		defer op.Push(u.gtx.Ops).Pop()
		op.Offset(f32.Pt(WidthPx/2.325, HeightPx/2.1)).Add(u.gtx.Ops)
		u.th.Color.Text = color.RGBA{0xFF, 0xFF, 0xFF, 0xFF}
		e := material.Editor(u.th, u.nameEditor, "Enter Name")
		e.Font.Style = text.Italic
		border := widget.Border{Color: color.RGBA{A: 0xFF}, CornerRadius: unit.Dp(8), Width: unit.Dp(2)}
		border.Layout(u.gtx, func(gtx layout.Context) layout.Dimensions {
			return layout.UniformInset(unit.Dp(8)).Layout(u.gtx, e.Layout)
		})
	}
}

// draws main menu
func (u *Ui) mainMenu() {
	defer op.Push(u.gtx.Ops).Pop()
	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(u.gtx,
		// vertical offset dimension and background
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			paint.ColorOp{Color: color.RGBA{0, 100, 0, 0xF}}.Add(gtx.Ops)
			paint.PaintOp{Rect: f32.Rect((WidthPx/2)-(topMenuPx/2)-topMenuPx/4, HeightPx/4.2,
				(WidthPx/2)+(topMenuPx/2)+topMenuPx/4, HeightPx/1.28)}.Add(gtx.Ops)
			return layout.Dimensions{Size: image.Pt(int(topMenuPx), int(HeightPx/4))}
		}),
		// new game button
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return u.newgameBtn.Layout(gtx, newgameImg, u.highscoresBtn.pressed)
		}),
		// highscores button
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return u.highscoresBtn.Layout(gtx, highscoresImg, u.highscoresBtn.pressed)
		}),
		// back button
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return u.backmenuBtn.Layout(gtx, backImg, u.highscoresBtn.pressed)
		}),
		// exit button
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return u.exitBtn.Layout(gtx, exitImg, u.highscoresBtn.pressed)
		}),
	)
}

// layout for highscore+name value pair
func (v *hs) Layout(gtx layout.Context) layout.Dimensions {
	var scr string
	if v.score > -1 {
		scr = strconv.Itoa(v.score)
		if len(scr) > 4 {
			scr = scr[0:4] + "..."
		}
	} else {
		scr = "N\\A"
	}
	var nm string
	if v.name != "" {
		nm = v.name
		if len(nm) > 5 {
			nm = nm[0:5] + "..."
		}
	} else {
		nm = "N\\A"
	}
	th := material.NewTheme(gofont.Collection())
	th.TextSize = unit.Dp(topMenuPx / 5)
	score := material.H2(th, scr)
	score.Color = color.RGBA{50, 205, 50, 0xFF}
	name := material.H2(th, nm)
	name.Color = color.RGBA{0, 66, 37, 0xFF}
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		// uniform spacer (not using inset)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Pt(int(WidthPx/5.5), 0)}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return name.Layout(gtx)
		}),
	)
	layout.Flex{Axis: layout.Horizontal}.Layout(gtx,
		// uniform spacer (not using inset)
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Pt(int(WidthPx/1.5), 0)}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return score.Layout(gtx)
		}),
	)
	return layout.Dimensions{Size: image.Pt(0, int(topMenuPx))}
}

// gets data from highscores.txt file via readHs() and draws data to screen
func (u *Ui) drawHs() {
	// draw high scores
	top := u.readHs()
	defer op.Push(u.gtx.Ops).Pop()
	op.Offset(f32.Pt(-(topMenuPx * 0.7), (HeightPx/4)-topMenuPx)).Add(u.gtx.Ops)
	paint.ColorOp{Color: color.RGBA{69, 150, 76, 0xFF}}.Add(u.gtx.Ops)
	paint.PaintOp{Rect: f32.Rect(WidthPx/5.6, 0, WidthPx/1.02, HeightPx/1.7)}.Add(u.gtx.Ops)
	var children []layout.FlexChild
	children = append(children,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			paint.ColorOp{Color: color.RGBA{0, 100, 0, 0xFF}}.Add(u.gtx.Ops)
			paint.PaintOp{Rect: f32.Rect(WidthPx/5.6, 0, WidthPx/1.02, topMenuPx)}.Add(u.gtx.Ops)
			paint.ImageOp(highscoresTitleImg).Add(gtx.Ops)
			paint.PaintOp{Rect: f32.Rect(WidthPx/6, 0, WidthPx/1.3, topMenuPx)}.Add(u.gtx.Ops)
			tempst := op.Push(u.gtx.Ops)
			op.Offset(f32.Pt(topMenuPx*3, 0)).Add(u.gtx.Ops)
			u.backhighscoresBtn.Layout(u.gtx, backImg, u.highscoresBtn.pressed)
			tempst.Pop()
			return layout.Dimensions{Size: image.Pt(int(WidthPx/1.02), int(topMenuPx))}
		}),
	)
	for i := range top {
		v := &top[i]
		children = append(children,
			layout.Rigid(func(gtx layout.Context) layout.Dimensions {
				return v.Layout(gtx)
			}),
		)
	}
	layout.Flex{Axis: layout.Vertical}.Layout(u.gtx, children...)
}

func (u *Ui) gameOver() {
	defer op.Push(u.gtx.Ops).Pop()
	// greyed out background
	paint.ColorOp{Color: color.RGBA{0, 0, 0, 0xDF}}.Add(u.gtx.Ops)
	paint.PaintOp{Rect: f32.Rect(0, 0, WidthPx, HeightPx)}.Add(u.gtx.Ops)
	// game over text
	paint.ImageOp(gameoverImg).Add(u.gtx.Ops)
	paint.PaintOp{Rect: f32.Rect(0, 0, WidthPx, HeightPx)}.Add(u.gtx.Ops)
	// continue button
	layout.Flex{
		Axis: layout.Vertical,
	}.Layout(u.gtx,
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return layout.Dimensions{Size: image.Pt(int(topMenuPx), int(HeightPx/4))}
		}),
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return u.continueBtn.Layout(gtx, continueImg, u.highscoresBtn.pressed)
		}),
	)
}

// draws top menu
func (u *Ui) topMenu() {
	defer op.Push(u.gtx.Ops).Pop()
	layout.Flex{
		Axis: layout.Horizontal,
	}.Layout(u.gtx,
		// topMenu background
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			paint.ColorOp{Color: color.RGBA{69, 150, 76, 0xFF}}.Add(gtx.Ops)
			paint.PaintOp{Rect: f32.Rect(0, 0, WidthPx, topMenuPx)}.Add(gtx.Ops)
			return layout.Dimensions{Size: image.Pt(0, 0)}
		}),
		// carrot image
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			paint.ImageOp(carrot2Img).Add(gtx.Ops)
			paint.PaintOp{Rect: f32.Rect(0, 0, topMenuPx, topMenuPx)}.Add(gtx.Ops)
			return layout.Dimensions{Size: image.Pt(int(topMenuPx), int(topMenuPx))}
		}),
		// score counter
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			score := material.H2(u.th, strconv.Itoa(u.ga.Score))
			score.Color = color.RGBA{0, 0, 0, 0xFF}
			score.Layout(gtx)
			return layout.Dimensions{Size: image.Pt(int(WidthPx/1.4), int(topMenuPx))}
		}),
		// menu button
		layout.Rigid(func(gtx layout.Context) layout.Dimensions {
			return u.menuBtn.Layout(gtx, menuImg, u.highscoresBtn.pressed)
		}),
	)
}
