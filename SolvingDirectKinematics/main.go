package main

/*
The code is a simulation of a robotic arm. It defines the dimensions and characteristics
of the arm and allows the user to input commands through the keyboard. The program then
updates the arm's position and displays it on the screen.

The code first imports necessary packages and defines constants for the screen dimensions.
It then creates a manipulator object with specific characteristics such as link lengths and starting position.

The App struct defines the properties of the application such as the screen dimensions
and the text displayed on the screen. The Update function updates the text displayed
on the screen based on user input. The Draw function draws the robotic arm on the
screen and updates its position based on the manipulator object. The Layout function
sets the dimensions of the screen.

The main function prompts the user to input whether they want to solve the inverse
task and if yes, prompts for the x and y coordinates. It then sets the window size
and runs the game using the App struct.

Overall, the code creates a simulation of a robotic arm that can be controlled through user input.
*/

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	cscreenWidth  = 800
	cscreenHeight = 800
)

var RoboManip = NewManipulator(
	float64(cscreenWidth),    //window dimensions
	float64(cscreenHeight),   //window dimensions
	float64(cscreenWidth/2),  //working point
	float64(cscreenHeight/2), //working point
	3,
	[]float64{50, 50, 50}, //link lengths
	[]float64{
		math.Pi * 5,
		math.Pi / 8,
		math.Pi / 8}, //starting position
	[]float64{
		0.01,
		0.001,
		0.005}, //rotation speed
	[]float64{0, 0, 0}) //acceleration of rotation

type App struct {
	angle        float64
	screenWidth  int
	screenHeight int
	runes        []rune
	text         string
	counter      int
}

func (g *App) Update() error {
	// Add runes that are input by the user by AppendInputChars.
	// Note that AppendInputChars result changes every frame, so you need to call this
	// every frame.
	g.runes = ebiten.AppendInputChars(g.runes[:0])
	g.text += string(g.runes)

	// Adjust the string to be at most 10 lines.
	ss := strings.Split(g.text, "\n")
	if len(ss) > 10 {
		g.text = strings.Join(ss[len(ss)-10:], "\n")
	}

	// If the enter key is pressed, add a line break.
	if repeatingKeyPressed(ebiten.KeyEnter) || repeatingKeyPressed(ebiten.KeyNumpadEnter) {
		g.text += "\n"
	}

	// If the backspace key is pressed, remove one character.
	if repeatingKeyPressed(ebiten.KeyBackspace) {
		if len(g.text) >= 1 {
			g.text = g.text[:len(g.text)-1]
		}
	}

	g.counter++
	return nil
}
func (a *App) Draw(screen *ebiten.Image) {
	screen.Fill(color.White)
	a.draw_grid(screen)

	RoboManip.calculate_movement()
	a.draw_manip(screen, RoboManip)
	a.draw_text(screen)
	a.draw_iterface(RoboManip.direct_kinematics(cscreenWidth, cscreenHeight), RoboManip.working_area, screen)

	//fmt.Println(RoboManip.direct_kinematics(cscreenWidth, cscreenHeight)["xHand"], " - ", RoboManip.robotStruct[2]["endX"])

}
func (a *App) Layout(outsideWidth, outsideHeight int) (int, int) {
	return a.screenWidth, a.screenHeight
}
func main() {
	var inverse_task string
	fmt.Println("Solving the inverse task yes/no?")
	fmt.Scanf("%s\n", &inverse_task)
	if inverse_task == "yes" {
		var task_x float64
		var task_t float64
		fmt.Println("Enter the x coordinate.")
		fmt.Scanf("%f\n", &task_x)
		fmt.Println("Enter the y coordinate.")
		fmt.Scanf("%f\n", &task_t)

		RoboManip.invers_kinematics()
	}

	ebiten.SetWindowSize(cscreenWidth, cscreenHeight)
	ebiten.SetWindowTitle("Robo - arm")
	App := &App{
		angle:        0,
		screenWidth:  cscreenWidth,
		screenHeight: cscreenHeight,
		text:         "Type on the keyboard:\n",
		counter:      0}

	if err := ebiten.RunGame(App); err != nil {
		log.Fatal(err)
	}
}
