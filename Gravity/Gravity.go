package main

import (
    "math"
    "fmt"
    "image/color"

    "github.com/faiface/pixel"
    "github.com/faiface/pixel/imdraw"
    "github.com/faiface/pixel/pixelgl"
    "github.com/faiface/pixel/text"
    "golang.org/x/image/font/basicfont"
    "golang.org/x/image/colornames"
)

var (
    win *pixelgl.Window
    imd *imdraw.IMDraw
)

func Closest(c pixel.Circle, p pixel.Vec) pixel.Vec {
    dif := p.Sub(c.Center)
    normalized := dif.Scaled(1/dif.Len())
    point := normalized.Scaled(c.Radius).Add(c.Center)
    return point
}

func AngleToPoint(c pixel.Circle, angle float64) pixel.Vec {
    radians := (-angle+90) * (math.Pi/180)
    point := pixel.V(math.Sin(radians), math.Cos(radians)).Scaled(c.Radius).Add(c.Center)
    return point
}

func run() {
    monitor := pixelgl.PrimaryMonitor()
    PositionX, PositionY  := monitor.Position()
    SizeX, SizeY := monitor.Size()
    screen := pixel.R(PositionX, PositionY, SizeX, SizeY)

    cfg := pixelgl.WindowConfig{
        Title:   "Projection",
        Monitor: pixelgl.PrimaryMonitor(),
        Bounds:  screen,
    }
    var err error
    win, err = pixelgl.NewWindow(cfg)
    if err != nil {
        panic(err)
    }

    imd = imdraw.New(nil)

    Circle := pixel.C(win.Bounds().Center().Sub(pixel.V(150, 150)), 300)
    Point := AngleToPoint(Circle, 0)
    Force := pixel.V(0, -1)

    dif := Point.Sub(Circle.Center)
    Tangent := dif.Normal().Scaled(1/dif.Normal().Len())
    Torque := Tangent.Scaled(Tangent.Dot(Force))

    Atlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
    msg := text.New(pixel.V(500, 100), Atlas)
    msg.Orig = msg.BoundsOf(" ").Size().Scaled(5)

    for !win.Closed() {
        
        win.Clear(colornames.Black)
        imd.Clear()
        msg.Clear()

        if win.JustPressed(pixelgl.KeyEscape) {
            win.SetClosed(true)
        }

        if win.Pressed(pixelgl.MouseButtonLeft) && !win.MousePosition().Eq(Circle.Center) {
            Mouse := win.MousePosition()
            if Mouse.X >= Circle.Center.X && Mouse.Y <= Circle.Center.Y {
                Point = Circle.Center.Add(pixel.V(300, 0))
            } else if Mouse.X <= Circle.Center.X && Mouse.Y >= Circle.Center.Y {
                Point = Circle.Center.Add(pixel.V(0, 300))
            } else if Mouse.X < Circle.Center.X && Mouse.Y < Circle.Center.Y {
                Point = Closest(Circle, win.MousePosition())
                Dif := Point.Sub(Circle.Center)
                Point = pixel.V(-Dif.Y, -Dif.X).Add(Circle.Center)
            } else {
                Point = Closest(Circle, win.MousePosition())
            }

            dif := Point.Sub(Circle.Center)
            Tangent = dif.Normal().Scaled(1/dif.Normal().Len())
            Torque = Tangent.Scaled(Tangent.Dot(Force))
        }

        msg.WriteString(fmt.Sprint(math.Round(Torque.Len()*math.Pow(10, 5))/math.Pow(10, 5)))

        imd.Color = color.RGBA{50, 50, 50, 255}
        imd.Push(Point)
        imd.Circle(10, 0)

        imd.Color = colornames.Dimgray
        imd.Push(Point)
        imd.Push(Point.Add(Force.Scaled(300)))
        imd.Line(3)

        imd.Color = colornames.White
        imd.Push(Point)
        imd.Circle(6, 0)
        

        imd.Push(Circle.Center)
        imd.CircleArc(Circle.Radius, 0, math.Pi/2, 3)

        imd.Push(Point)
        imd.Push(Point.Add(Torque.Scaled(300)))
        imd.Line(3)
        
        msg.Draw(win, pixel.IM.Scaled(msg.Orig, 10))
        imd.Draw(win)
        win.Update()
    }
}

func main() {
    pixelgl.Run(run)
}