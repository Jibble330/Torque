package main

import (
    "time"
    "math"

    "github.com/faiface/pixel"
    "github.com/faiface/pixel/imdraw"
    "github.com/faiface/pixel/pixelgl"
    "golang.org/x/image/colornames"
)

const (
    MOTOR = 1.0/3.0
    LOAD = 1.0
    MA = LOAD/MOTOR
    WEIGHT = 2
)

var (
    win *pixelgl.Window
    imd *imdraw.IMDraw
)

type Point struct {
    Pos pixel.Vec
    Radius float64
    Moving bool
}

func NewPoint(Pos pixel.Vec, Radius float64) Point {
    return Point {Pos, Radius, false}
}

func (p *Point) Draw() {
    imd.Push(p.Pos)
    imd.Circle(p.Radius, 0)
}

func (p *Point) Inside(Point pixel.Vec) bool {
    Point = Point.Sub(p.Pos)
    Distance := math.Hypot(Point.X, Point.Y)
    return Distance <= p.Radius
}

type Line struct {
    Start, End pixel.Vec
}

func NewLine(Start, End pixel.Vec) Line {
    return Line{Start, End}
}

func (l *Line) Draw(thickness float64) {
    imd.Push(l.Start, l.End)
    imd.Line(thickness)
}

func (l *Line) Vec() pixel.Vec {
    return l.End.Sub(l.Start)
}

func Lever(Input float64) float64 {
    return Input*MA
}

func Project(v pixel.Vec, u pixel.Vec) pixel.Vec { //Projects v onto u
    ULen := math.Hypot(u.X, u.Y)
    if ULen != 1 {
        u.X /= ULen
        u.Y /= ULen
    }
    Dot := v.Dot(u)
    Projection := u.Scaled(Dot)
    return Projection
}

func run() {
    monitor := pixelgl.PrimaryMonitor()
    PositionX, PositionY  := monitor.Position()
    SizeX, SizeY := monitor.Size()
    screen := pixel.R(PositionX, PositionY, SizeX, SizeY)

    cfg := pixelgl.WindowConfig{
        Title:   "Arm",
        Monitor: pixelgl.PrimaryMonitor(),
        Bounds:  screen,
    }
    var err error
    win, err = pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }

    imd = imdraw.New(nil)
    fps := time.NewTicker(time.Second/60)
    defer fps.Stop()

    

    l1 := NewLine(win.Bounds().Center(), win.Bounds().Center().Add(pixel.V(0, -250)))
    l2 := NewLine(l1.End, l1.End.Add(pixel.V(0, -250)))
    
    p := NewPoint(l1.End, 10)
    

    for !win.Closed() {
        win.Clear(colornames.Black)
        imd.Clear()

        if win.JustPressed(pixelgl.KeyEscape) {
            win.SetClosed(true)
        }

        if win.JustPressed(pixelgl.MouseButtonLeft) && p.Inside(win.MousePosition()) {
            p.Moving = true
        }
        if win.JustReleased(pixelgl.MouseButtonLeft) {
            p.Moving = false
        }

        imd.Color = colornames.Dimgray
        g := pixel.V(0, -WEIGHT)
        
        
        imd.Push(l1.Start, l1.Start.Add(g.Scaled(100)))
        imd.Line(3)

        if p.Moving {
            p.Pos.Y = win.MousePosition().Y
        }

        y := -win.Bounds().Center().Sub(p.Pos).Y
        var x float64
        if y <= 250 && y >= -250 {
            x = math.Sqrt(250*250 - y*y)
        } else {
            y = math.Mod(y, 250)
            x = math.Sqrt(250*250 - y*y)
        }

        l1.End = l1.Start.Add(pixel.V(x, y))
        l2.Start = l1.End
        l2.End = l2.Start.Add(pixel.V(-x, y))
        p.Pos.X = x + win.Bounds().Center().X

        f1 := Project(g, l1.Vec().Unit())
        f2 := Project(f1, l2.Vec().Normal().Unit())

        imd.Color = colornames.White
        l1.Draw(3)
        l2.Draw(3)

        imd.Color = colornames.Dimgray
        imd.Push(l1.End, l1.End.Add(f1.Scaled(100)))
        imd.Line(3)
        imd.Color = colornames.White
        imd.Push(l2.Start, l2.Start.Add(f2.Scaled(100)))
        imd.Line(3)

        imd.Color = colornames.White
        p.Draw()

        imd.Draw(win)
        win.Update()
        <-fps.C
    }
}

func main() {
    pixelgl.Run(run)
}