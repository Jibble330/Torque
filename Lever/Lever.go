package main

import (
    "math"
    "errors"
    "fmt"
    //"image/color"

    "github.com/faiface/pixel"
    "github.com/faiface/pixel/imdraw"
    "github.com/faiface/pixel/pixelgl"
    "golang.org/x/image/colornames"
)

const (
    LOAD_WEIGHT = 1
)

type Lever struct {
    Pos pixel.Vec
    Len, Angle float64
    Fulcrum, Effort, Load float64 //Between 0 and 1
    Class uint8 //Class 1, 2, or 3
}

func NewLever(Pos pixel.Vec, Len float64, Fulcrum, Effort, Load float64) (Lever, error) {
    var Class uint8
    if (Fulcrum < Effort && Fulcrum > Load) || (Fulcrum > Effort && Fulcrum < Load) {
        Class = 1
    } else if (Effort < Fulcrum && Effort > Load) || (Effort > Fulcrum && Effort < Load) {
        Class = 2
    } else if (Load < Effort && Load > Fulcrum) || (Load > Effort && Load < Fulcrum) {
        Class = 3
    } else {
        return Lever{}, errors.New("Points must be different")
    }

    return Lever{Pos, Len, -90*(math.Pi/180), Fulcrum, Effort, Load, Class}, nil
}

func (l *Lever) Output(Effort float64) pixel.Vec {
    var Force float64
    if l.Class == 3 {
        Force = Effort*(l.Effort/l.Load)
    } else if l.Class == 2 {
        Force = Effort*(l.Load/l.Effort)
    } else if l.Class == 1 {
        EffortDif := math.Abs(l.Effort-l.Fulcrum)
        LoadDif := math.Abs(l.Load-l.Fulcrum)
        Force = Effort*(EffortDif/LoadDif)
    }
    return pixel.V(0, Force).Rotated(-l.Angle)
}

func (l *Lever) Draw() {
    End := l.Pos.Add(pixel.V(l.Len, 0).Rotated(-l.Angle))

    imd.Color = colornames.Gray
    imd.Push(l.Pos)
    imd.Push(End)
    imd.Line(15)

    imd.Color = colornames.Dimgray
    imd.Push(l.Pos)
    imd.Push(End)
    imd.Line(5)

    Fulcrum := l.Pos.Add(pixel.V(l.Len*l.Fulcrum, 0).Rotated(-l.Angle))

    imd.Color = colornames.Lightgray
    imd.Push(Fulcrum)
    imd.Circle(10, 0)

    imd.Color = colornames.Dimgray
    imd.Push(Fulcrum)
    imd.Circle(5, 0)
}

func AngleToPoint(c pixel.Circle, angle float64) pixel.Vec {
    Rad := (-angle+90) * (math.Pi/180)
    Point := pixel.V(math.Sin(Rad), math.Cos(Rad)).Scaled(c.Radius).Add(c.Center)
    return Point
}

var (
    win *pixelgl.Window
    imd *imdraw.IMDraw
)

const MAGIC64 = 0x5FE6EB50C7B537A9

func FastInvSqrt64(n float64) float64 {
    if n < 0 {
        return math.NaN()
    }
    n2, th := n*0.5, float64(1.5)
    b := math.Float64bits(n)
    b = MAGIC64 - (b >> 1)
    f := math.Float64frombits(b)
    f *= th - (n2 * f * f)
    return f
}

func Unit(v pixel.Vec) pixel.Vec {
    squared := v.X*v.X + v.Y*v.Y
    scale := FastInvSqrt64(squared)
    return v.Scaled(scale)
}

func run() {
    monitor := pixelgl.PrimaryMonitor()
    PositionX, PositionY  := monitor.Position()
    SizeX, SizeY := monitor.Size()
    screen := pixel.R(PositionX, PositionY, SizeX, SizeY)

    cfg := pixelgl.WindowConfig{
        Title:   "Lever",
        Monitor: pixelgl.PrimaryMonitor(),
        Bounds:  screen,
    }
    var err error
    win, err = pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }

    imd = imdraw.New(nil)

    Arm, _ := NewLever(win.Bounds().Center().Sub(pixel.V(250, 0)), 500, 0, 0.5, 1)

    fmt.Println(Arm.Output(1))
    fmt.Println(Arm.Output(1).Len())

    for !win.Closed() {
        win.Clear(colornames.Black)
        imd.Clear()

        if win.JustPressed(pixelgl.KeyEscape) {
            win.SetClosed(true)
        }

        if win.Pressed(pixelgl.MouseButtonLeft) {
            dif := win.MousePosition().Sub(Arm.Pos)
            Arm.Angle = -dif.Angle()
        }

        Arm.Draw()

        imd.Draw(win)
        win.Update()
    }
}

func main() {
    pixelgl.Run(run)
}