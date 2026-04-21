package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"golang.org/x/term"
)

// ============================================
// STRUKTUR DATA
// ============================================

type ParamDef struct {
	Name  string
	Value float64
}

type Layer struct {
	Name       string
	Enabled    bool
	Alpha      float64 // Current visibility (0.0 - 1.0)
	TargetAlpha float64 // Target visibility for transition
	Params     []ParamDef
}

// ============================================
// TRANSITION SYSTEM
// ============================================

type TransitionState struct {
	AutoMode         bool    // Auto cycle through sequences
	CurrentLayer     int     // Currently showing layer (auto mode)
	NextLayer        int     // Layer transitioning to (auto mode)
	HoldDuration     float64 // How long to hold each sequence
	HoldTimer        float64 // Current hold timer
	TransitionSpeed  float64 // Speed of fade (applies to both modes)
	TransitionType   int     // 0=fade, 1=crossfade, 2=wipe
}

var transition = TransitionState{
	AutoMode:        false,
	CurrentLayer:    1,
	NextLayer:       2,
	HoldDuration:    150.0, // Frames to hold each sequence (~6 seconds at 25fps)
	HoldTimer:       0.0,
	TransitionSpeed: 0.04,  // Fade speed
	TransitionType:  0,
}

// ============================================
// TOPENG ASCII ART
// ============================================

var topengASCII = []string{
	"    ....                   ....:=+****###########**##########***+==-:...                       .:::.",
	"   ..  ...           ...:-+###*****###***********************+++++****###+-:..                ..:...",
	"....    ...        .:*##**********************#+=--=+#*******+++***********###+:.           ...    .",
	".         ...  ..-*#**************************=-=-==-=***+++++******************#*-..    ...        ",
	"            ..:##*********++++*************#+---====---++++************************##:....          ",
	"           .-##*********+*+++***********##%#=--+*--*+--=####*************************##-.           ",
	"         .:##*********+++******######%%%%%%*---%*--*%---+%%%%%%######******************##..         ",
	"        .+#*********+*++****##%%%%%%%%%%%%%*=--=*==*=--=*%%%%%%%%%%%%%##*********+##*****#+.        ",
	"      .:%#******+=+++####%%%%%%%%%%%%%%%%%%%=---====---=%%%%%%%%%%%%%%%%%%%######*+#******#%:.      ",
	"     .=#*******++++#%%%%%%%%%%%%%%%%%%%%%%%%%=--------=%%%%%%%%%%%%%%%%%%%%%%%%%%#**********#=.     ",
	"    .-%******++*+%%%%%%%%%%%%%%%%%%%%%%%%%%%%%+------+%%%%%%%%%%%%%%%%%%%%%%%%%%%%%##********#-.    ",
	"   .=%#****+*##%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%*----*%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%#****#%=.   ",
	".. :#******#*%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%#==*%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%******#:   ",
	"...*#****#%%%%%%%%%%%%%%%%%%%%%%%%%%%%#%%%%%%%%%%**%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%#****#*. .",
	" .=%****#%%%%%%%%%%%%%%%%%%%%%%%%%%%*++%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%#***##=..",
	"..##***%%%%%%%%%%%%%%%%%%%%%%%%%%%++#%%%%%%%%%%%%%%%%%%%%%%%%%%%%%##%%%%%%%%%%%%%%%%%%%%%%%%%#*##*:.",
	".-#***#%%%%%%%%%%%%%%%%%%%%%%%%%#*#%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%#*##%=.",
	".+#**#%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%*-::-*%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%%#####***#+.",
	".***#%%%%%%%%%%%%%%%%#*+------:.....:-===--:............:--===-:.....:------+*#%%%%%%%%#***####++**.",
	":#*#%%%%%%%%%%%%%+==-==----=======-..............................-=======----==-==+%%%%*#**###**+**:",
	":#*#%%%%%%%%%%%*==-----------------==:...........*+...........:==-----------------=-**#+*####***++#:",
	"-#*%%%%%%%%%%#==---====---------------=-........+%%=:.......-=-----------------=------######****++#-",
	"+##%%%%%%%%#+=----=====-----===---------=-:..===*%#*+==..:-=---------===---------------=*##******=*=",
	"*#%%%%%%%#==---=-:.::::..........-==-------==-..+%%+:.-==-------==-...............:-=-----******++*+",
	"#%%%%%%%%+---=:.......................::--==---==##==--===--::.......................:----=****+++**",
	"#%%%%%%%#--=-.............................*.....====....:+.............................:---+**++++*#",
	"%%%%%%%%===-*=......::-==+*##+:.......:-:##..+-.:=-:.-=..##:-:.......:+*+====-:.......-+:---++++++*%",
	"%%%%%%%=-=-:=*#%%%##****#%%%%%%%+.......=%%-.....==.....-%%+.......+%%%%##%#*****####*+-::---++++=*%",
	"%%%%%%-==-+:.............=:..:-+#%#+:...+#%+............=%#+...:+#%#+-:..:=.............:=:-::++==+%",
	"%%%%%*--=:.:#:.....................=#%*-::%+............=#::-*%#=.....................:+:..:::=+=+*%",
	"%%%%%#==-..................:..:........:--:..............:--:........:.....................:::=+==*%",
	"%%%%%===..........::.......=:.-......................................-::=.......::..........::-===+#",
	"#%%%#-=-...-*#%%%%%%%%%%%###+-*......................................*-+###%%%%%%%%%%%#*:...::.-=-+#",
	"*%%%+-=:.:+++=-:.....:-+*%%+%%%%%+-:............................:-+%%%%%+#%*+-:.....:-=+++:..:::--**",
	"+%%%===:....:-*##*+=:.....:=*##%%#%%%+:......................:+%%%#%%##*=:.....:-+*##*-:.....:::--#+",
	"=%%%-=-..........-+*#%%%#%%#**++===%%%%*:..................:*%%%%===++**#%%%%%%#*+:..........::.--#=",
	":%%#:=-...............::-*%%%%%%%%%%%%%%*:........::::....:*%%%%%%%%%%%%%%*=:...::::..........:.:-%-",
	":#%#:=-.......................:-=+**###%%-.......::::.....:%%###**+=-:..........::::..........:.:+#:",
	".*%*:=-.........................................................................................:#*.",
	".+%#:=-.........................................................................................-%+.",
	" -%#:--.........................................................................................+%-.",
	" .##::=:........................................................................................*%. ",
	"..=#:.=:.......................................................................................:#=..",
	"  .#-:==.......................................................................................-#...",
	"  .+*:-+......=#....................-%*......................*%-....................#-.........+=.. ",
	"   :#-:=-..........................-:-%......................%-:-.............................:#:...",
	"....+#.:=:.........................:..++:..................:*+..:.............................#*.   ",
	"..  .%=.-=.............................=%%%#**+=----=+**#%%%-................................=%.    ",
	"     -#====.....................................:::.........................................-#-     ",
	"     .+*===-................................................................................*+.     ",
	"     .:#=:=-=..............................................................................=#:      ",
	"      .:%-................................................................................-%:.      ",
	"       .-#:...............................................................................#-.       ",
	"        .=%:.............................................................................%+.        ",
	"         .+*:............:******=-------::::::::::::::::::::-------=****++:.............*+.         ",
	"          .=*.............:******+=+++=----------------------=+++=+****++:.............#=..         ",
	"..       ...=#:.............=*******+**+===++++======++++===+**+*******=.............:#+. ...       ",
	" .........  .-#:..............=******+==+*****+++**+++*****+==+******=..............:#-.    ...   ..",
	"   .:::.     .-%=...............-+*******++==+==++++==+==++*******+-...............=%-.       ......",
	"   .:::.       .#+................:+=*********+==+===+**********=.................+#.          .... ",
	"  .... ...      .+%-.................-************************-.................-#+.        ....  ..",
	"...      ...     ..#+..................:*******************+:..................+#.        ...       ",
	"           ..      .=%-..................-****************-..................-#=..       ...        ",
	"             ...     .*#:..................-************-..................:#*.        ..           ",
	"               ...   ..:*+...................-********-...................*#:..     ....            ",
	"                  .......-#*:..................:=**+:..................:*#-.   ......               ",
	"                ........  .-#*:..................::::................:*#:.     .:::.                ",
	"              ..........    .:*#-.................:::..............-#*:.      ......                ",
	"             .......   ...   ...=%*..............................*%=...   .....     ...             ",
	"          .......        ....    .:*#=:.......................=#*:..   ......          ...          ",
	"          ......           ....   ...-#%-..................-%#-..      .......           ...        ",
	"...     .....                 ....   ...-*#=:..........:=#*-...     ........               ....     ",
	"   ..  ..                        ..   ..  .:#%*:....:*%*:.     .. .......                     ..  ..",
	"    ....                          .:::.     ...-*##*-..        ........                        ...  ",
}

// ============================================
// GLOBAL VARIABLES
// ============================================

var (
	frame         int
	selectedLayer = 1
	selectedParam = 0
	layers        [14]Layer
	colorMode     = true
	globalBright  = 1.0
	running       = true

	balls         []Ball
	typistColumns []TypistColumn

	gradients = []string{" ", ".", ":", ";", "=", "+", "*", "#", "%", "@", "█"}

	layerColors = []string{
		"",
		"\033[38;5;51m",  // 1: Cyan (Portal)
		"\033[38;5;226m", // 2: Yellow (Ball Matrix)
		"\033[38;5;196m", // 3: Red (Topeng)
		"\033[38;5;46m",  // 4: Green (Matrix Typist)
		"\033[38;5;208m", // 5: Orange (Tubuh)
		"\033[38;5;75m",  // 6: Electric Blue (Mountain)
		"\033[38;5;255m", // 7: White (Barcode)
		"\033[38;5;21m",  // 8: Blue
		"\033[38;5;201m", // 9: Magenta
		"\033[38;5;220m", // 10: Gold
		"\033[38;5;129m", // 11: Purple
		"\033[38;5;240m", // 12: Dark Gray
		"\033[38;5;231m", // 13: Bright White (Particle Wave)
	}
)

// ============================================
// STRUCTURES
// ============================================

type Ball struct {
	X, Y   float64
	VX, VY float64
	Size   float64
	Char   int
}

type TypistColumn struct {
	X      int
	Chars  []rune
	Speed  float64
	Offset float64
	Length int
}

// ============================================
// TERMINAL SIZE
// ============================================

func getTerminalSize() (int, int) {
	width, height, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 120, 30
	}
	return width, height
}

// ============================================
// EASING FUNCTIONS
// ============================================

func easeInOutCubic(t float64) float64 {
	if t < 0.5 {
		return 4 * t * t * t
	}
	return 1 - math.Pow(-2*t+2, 3)/2
}

func easeInOutQuad(t float64) float64 {
	if t < 0.5 {
		return 2 * t * t
	}
	return 1 - math.Pow(-2*t+2, 2)/2
}

// ============================================
// LAYER ALPHA TRANSITION (BOTH MODES)
// ============================================

func updateLayerAlphas() {
	for i := 1; i <= 13; i++ {
		// Smoothly interpolate Alpha towards TargetAlpha
		diff := layers[i].TargetAlpha - layers[i].Alpha
		
		if math.Abs(diff) < 0.001 {
			layers[i].Alpha = layers[i].TargetAlpha
		} else {
			// Ease the transition
			layers[i].Alpha += diff * transition.TransitionSpeed * 2
		}
		
		// Clamp
		layers[i].Alpha = math.Max(0, math.Min(1, layers[i].Alpha))
	}
}

func updateAutoMode() {
	if !transition.AutoMode {
		return
	}

	// Check if current layer is fully visible and next is fully hidden
	currentAlpha := layers[transition.CurrentLayer].Alpha
	
	if currentAlpha >= 0.99 {
		// Holding current sequence
		transition.HoldTimer += 1.0
		
		if transition.HoldTimer >= transition.HoldDuration {
			// Start transition to next
			transition.HoldTimer = 0
			
			// Fade out current
			layers[transition.CurrentLayer].TargetAlpha = 0
			
			// Fade in next
			layers[transition.NextLayer].TargetAlpha = 1
			
			// Prepare next cycle
			transition.CurrentLayer = transition.NextLayer
			transition.NextLayer = transition.CurrentLayer + 1
			if transition.NextLayer > 13 {
				transition.NextLayer = 1
			}
		}
	}
}

// ============================================
// INITIALIZATION
// ============================================

func initLayers() {
	names := []string{
		"",
		"PORTAL LOOPING",       // 1
		"BALL MATRIX",          // 2
		"TOPENG SAKRAL",        // 3
		"MATRIX DATA TYPIST",   // 4
		"TUBUH YANG MENGINGAT", // 5
		"3D MOUNTAIN 360",      // 6
		"BARCODE GLITCH",       // 7
		"WAVE INTERFERENCE",    // 8
		"HYPER SPIRAL",         // 9
		"VOID ABYSS",           // 10
		"GLITCH STORM",         // 11
		"CHAOS ENTROPY",        // 12
		"PARTICLE WAVE FIELD",  // 13
	}

	for i := 1; i <= 13; i++ {
		layers[i] = Layer{
			Name:        names[i],
			Enabled:     false,
			Alpha:       0.0,
			TargetAlpha: 0.0,
			Params: []ParamDef{
				{"Speed", 0.5},
				{"Density", 0.5},
				{"Brightness", 0.8},
				{"Scale", 0.5},
				{"Chaos", 0.3},
			},
		}
	}
	
	// Start with layer 1 enabled
	layers[1].Enabled = true
	layers[1].Alpha = 1.0
	layers[1].TargetAlpha = 1.0
	
	initBalls(60)
	initTypist(100)
}

func initBalls(count int) {
	balls = make([]Ball, count)
	for i := range balls {
		balls[i] = Ball{
			X:    rand.Float64()*2 - 1,
			Y:    rand.Float64()*2 - 1,
			VX:   (rand.Float64() - 0.5) * 0.03,
			VY:   (rand.Float64() - 0.5) * 0.03,
			Size: rand.Float64()*0.08 + 0.02,
			Char: rand.Intn(len(gradients)-1) + 1,
		}
	}
}

func initTypist(columns int) {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%&*アイウエオカキクケコサシスセソタチツテト")
	typistColumns = make([]TypistColumn, columns)
	for i := range typistColumns {
		length := rand.Intn(20) + 8
		colChars := make([]rune, length)
		for j := range colChars {
			colChars[j] = chars[rand.Intn(len(chars))]
		}
		typistColumns[i] = TypistColumn{
			X:      i,
			Chars:  colChars,
			Speed:  rand.Float64()*0.8 + 0.3,
			Offset: rand.Float64() * 20,
			Length: length,
		}
	}
}

// ============================================
// HELPER FUNCTIONS
// ============================================

func noise2D(x, y float64) float64 {
	n := math.Sin(x*12.9898+y*78.233) * 43758.5453
	return n - math.Floor(n)
}

func getTopengChar(nx, ny, t, scale, chaos float64, width, height int) (rune, float64) {
	topengH := len(topengASCII)
	topengW := 0
	for _, line := range topengASCII {
		if len(line) > topengW {
			topengW = len(line)
		}
	}

	breathScale := 1.0 + math.Sin(t*0.8)*0.05*chaos
	centerX := float64(width) / 2
	centerY := float64(height) / 2

	screenX := int((nx + 1) * 0.5 * float64(width))
	screenY := int((ny + 1) * 0.5 * float64(height))

	actualScale := scale * breathScale * 1.2
	topengX := int((float64(screenX)-centerX)/actualScale + float64(topengW)/2)
	topengY := int((float64(screenY)-centerY)/actualScale + float64(topengH)/2)

	if chaos > 0.5 && rand.Float64() > 0.98 {
		topengX += rand.Intn(5) - 2
		topengY += rand.Intn(3) - 1
	}

	if topengY >= 0 && topengY < topengH && topengX >= 0 && topengX < len(topengASCII[topengY]) {
		char := rune(topengASCII[topengY][topengX])

		intensity := 0.0
		switch char {
		case ' ', '.':
			intensity = 0.1
		case ':', '-', '=':
			intensity = 0.3
		case '+', '*':
			intensity = 0.5
		case '#', '%':
			intensity = 0.7
		case '@', '█':
			intensity = 0.9
		default:
			intensity = 0.4
		}

		pulse := 1.0 + math.Sin(t*2)*0.2
		return char, intensity * pulse
	}

	return ' ', 0
}

// ============================================
// ALL RENDER FUNCTIONS
// ============================================

func renderPortalLooping(nx, ny, t, density, scale, chaos float64) float64 {
	r := math.Sqrt(nx*nx + ny*ny)
	theta := math.Atan2(ny, nx)

	zoom := math.Mod(t*scale*0.8, 1.0)
	rings := 0.0
	numRings := 12

	for i := 0; i < numRings; i++ {
		ringPhase := float64(i) / float64(numRings)
		ringPos := math.Mod(zoom+ringPhase, 1.0)
		ringR := (1.0 - ringPos) * (1.2 + density*0.5)
		thickness := 0.02 + chaos*0.03 + ringPos*0.02
		dist := math.Abs(r - ringR)

		if dist < thickness {
			intensity := 1.0 - (dist / thickness)
			intensity = math.Pow(intensity, 0.6)
			fade := math.Sin(ringPos * math.Pi)
			intensity *= fade
			rotAngle := theta + t*0.3 + float64(i)*0.5
			warp := math.Sin(rotAngle*(4+chaos*4)) * 0.2 * chaos
			intensity *= (1.0 + warp)
			rings = math.Max(rings, intensity)
		}
	}

	centerGlow := math.Exp(-r*6*(1-density*0.3)) * (0.5 + math.Sin(t*2)*0.3)
	arms := 3.0 + chaos*3
	spiralR := r + 0.001
	spiral := math.Sin(theta*arms+math.Log(spiralR)*5-t*3) * 0.25
	spiral *= math.Exp(-r * 1.5)
	edgeEnergy := math.Sin(theta*8+t*2-r*10) * 0.15 * math.Exp(-math.Abs(r-0.9)*8)
	result := rings*0.7 + centerGlow*0.4 + spiral + edgeEnergy
	innerDepth := math.Exp(-r*8) * math.Sin(t*4) * 0.2

	return (result + innerDepth) * scale
}

func renderBallMatrix(nx, ny, t, density, scale, chaos float64) (float64, string) {
	maxVal := 0.0
	resultChar := " "

	ballCount := int(density*50) + 10
	if ballCount > len(balls) {
		ballCount = len(balls)
	}

	for i := 0; i < ballCount; i++ {
		ball := balls[i]
		dx := nx - ball.X
		dy := (ny - ball.Y) * 2
		dist := math.Sqrt(dx*dx + dy*dy)
		size := ball.Size * (0.5 + scale)

		if dist < size {
			intensity := 1.0 - (dist / size)
			intensity = math.Pow(intensity, 0.3)
			pulse := 1.0 + math.Sin(t*4+float64(i))*0.3*chaos
			val := intensity * pulse
			if val > maxVal {
				maxVal = val
				ballChars := []string{".", "o", "O", "●", "◉", "@"}
				idx := int(intensity * float64(len(ballChars)-1))
				if idx >= len(ballChars) {
					idx = len(ballChars) - 1
				}
				resultChar = ballChars[idx]
			}
		}

		if chaos > 0.4 {
			trailX := ball.X - ball.VX*8
			trailY := ball.Y - ball.VY*8
			trailDist := math.Sqrt((nx-trailX)*(nx-trailX) + (ny-trailY)*(ny-trailY)*4)
			if trailDist < size*0.3 && maxVal < 0.2 {
				maxVal = 0.2
				resultChar = "·"
			}
		}
	}

	return maxVal, resultChar
}

func renderTopeng(nx, ny, t, density, scale, chaos float64, width, height int) (float64, rune) {
	char, intensity := getTopengChar(nx, ny, t, scale, chaos, width, height)

	if char == ' ' || intensity < 0.05 {
		r := math.Sqrt(nx*nx + ny*ny)
		if r < 0.8 {
			aura := (0.8 - r) * 0.3 * density
			aura *= 1 + math.Sin(t*2+r*10)*0.3
			return aura, '·'
		}
		return 0, ' '
	}

	eyeL := math.Sqrt(math.Pow(nx+0.15, 2) + math.Pow(ny+0.1, 2))
	eyeR := math.Sqrt(math.Pow(nx-0.15, 2) + math.Pow(ny+0.1, 2))
	if (eyeL < 0.08 || eyeR < 0.08) && chaos > 0.3 {
		glow := 1.0 + math.Sin(t*5)*0.5
		return intensity * glow * 1.5, '◉'
	}

	return intensity * density, char
}

func renderTypistMatrix(nx, ny, t, density, scale, chaos float64, width, height int) (float64, rune) {
	screenX := int((nx + 1) * 0.5 * float64(width))
	colCount := int(density*float64(len(typistColumns))*0.8) + 20
	if colCount > len(typistColumns) {
		colCount = len(typistColumns)
	}

	colSpacing := float64(width) / float64(colCount)
	colIdx := int(float64(screenX) / colSpacing)

	if colIdx < 0 || colIdx >= colCount {
		return 0, ' '
	}

	col := typistColumns[colIdx%len(typistColumns)]
	fallSpeed := col.Speed * scale * 2
	yPos := (ny + 1) * 0.5
	fallOffset := math.Mod(t*fallSpeed+col.Offset, 2.0)
	relPos := yPos - fallOffset + 1
	if relPos < 0 {
		relPos += 2
	}

	charIdx := int(relPos * float64(col.Length) * 0.5)

	if charIdx >= 0 && charIdx < len(col.Chars) {
		headDist := relPos
		if headDist < 0 {
			headDist = -headDist
		}

		intensity := 0.0
		if headDist < 0.3 {
			intensity = 1.0 - headDist*2
		} else if headDist < 0.8 {
			intensity = 0.5 - (headDist - 0.3)
		}

		intensity *= density

		if rand.Float64() > 0.97 && chaos > 0.3 {
			chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%")
			return intensity, chars[rand.Intn(len(chars))]
		}

		if intensity > 0.1 {
			return intensity, col.Chars[charIdx]
		}
	}

	if rand.Float64() > 0.995 {
		return 0.1, rune(".:;="[rand.Intn(4)])
	}

	return 0, ' '
}

func renderTubuhYangMengingat(nx, ny, t, density, scale, chaos float64) float64 {
	r := math.Sqrt(nx*nx + ny*ny)
	theta := math.Atan2(ny, nx)

	pulse1 := math.Sin(r*10*(1+density)-t*2) * math.Cos(theta*3+t*0.5)
	pulse2 := math.Sin(r*8-t*1.5+theta*2) * 0.5
	pulse3 := math.Cos(r*15*(1+chaos)-t*3) * math.Sin(theta*5) * 0.3
	breath := 1.0 + math.Sin(t*0.5)*0.2
	memory := noise2D(nx*5+t*0.1, ny*5) * chaos * 0.4

	return (pulse1 + pulse2 + pulse3) * breath * (1 - r*0.2) + memory
}

// smoothNoise interpolates noise2D with cubic smoothstep — removes blocky artifacts.
func smoothNoise(x, y float64) float64 {
	ix, iy := math.Floor(x), math.Floor(y)
	fx, fy := x-ix, y-iy
	ux := fx * fx * (3 - 2*fx)
	uy := fy * fy * (3 - 2*fy)
	return noise2D(ix, iy)*(1-ux)*(1-uy) +
		noise2D(ix+1, iy)*ux*(1-uy) +
		noise2D(ix, iy+1)*(1-ux)*uy +
		noise2D(ix+1, iy+1)*ux*uy
}

func mountainHeight(wx, wy, t, density, chaos float64) float64 {
	// Light domain warp for organic, non-repetitive terrain
	wx0 := wx + math.Sin(wx*0.6+wy*0.4+t*0.008)*0.35*chaos
	wy0 := wy + math.Cos(wx*0.4-wy*0.6+t*0.011)*0.35*chaos

	// Smooth fBm base — rolling terrain
	h := smoothNoise(wx0*0.9+t*0.006, wy0*0.9)*0.34 +
		smoothNoise(wx0*2.1+t*0.010, wy0*2.1)*0.22 +
		smoothNoise(wx0*4.3+t*0.016, wy0*4.3)*0.14 +
		smoothNoise(wx0*8.7+t*0.022, wy0*8.7)*0.08*(0.5+density*0.5) +
		smoothNoise(wx0*17.5+t*0.028, wy0*17.5)*0.04*density

	// Prominent Gaussian peak — exp(-r²) is the natural shape of a mountain.
	// Period = 12 world units: only ONE peak visible at a time as camera pans.
	rx := wx - math.Round(wx/12.0)*12.0
	ry := wy - 2.0
	mainPeak := 0.76 * math.Exp(-rx*rx*0.45-ry*ry*0.45)

	// Smaller companion ridgeline slightly behind the main peak
	rx2 := wx - math.Round((wx-1.2)/12.0)*12.0 - 1.2
	ry2 := wy - 2.6
	ridge := 0.38 * math.Exp(-rx2*rx2*0.9-ry2*ry2*0.7)

	return math.Max(0, math.Min(1, h/0.82+mainPeak+ridge))
}

func renderMountain3D(nx, ny, t, density, scale, chaos float64) float64 {
	// Perspective: camera at (0,0,eyeH), looking forward.
	// Surface eq: (eyeH-h)/wy = (ny-horizon)/focal
	const horizon = -0.15
	const eyeH = 1.00
	const focal = 0.80

	if ny <= horizon {
		return 0
	}

	pan := t * 0.25 * scale

	// Damped iteration to find terrain surface depth
	wy := eyeH * focal / (ny - horizon)
	h := 0.0
	for i := 0; i < 9; i++ {
		wx := wy*nx/focal + pan
		h = mountainHeight(wx, wy, t, density, chaos)
		nw := (eyeH - h) * focal / (ny - horizon)
		if nw < 0.03 {
			return 0
		}
		wy = wy*0.3 + nw*0.7
	}

	wx := wy*nx/focal + pan

	// Exponential atmospheric haze
	fade := math.Exp(-wy * 0.17)

	// Rectangular mesh grid — like a cloth draped over terrain.
	// Both axes in world space so grid stays square, perspective does the rest.
	gf := 3.0 + density*2.5
	cellX := math.Mod(math.Abs(wx)*gf, 1.0)
	cellY := math.Mod(wy*gf, 1.0)
	dX := math.Min(cellX, 1.0-cellX)
	dY := math.Min(cellY, 1.0-cellY)

	lw := 0.044 + chaos*0.012
	lvX := math.Max(0, 1.0-dX/lw)
	lvY := math.Max(0, 1.0-dY/lw)
	lv := math.Max(lvX, lvY)
	if lvX > 0.08 && lvY > 0.08 {
		lv = math.Min(1.0, lv*1.65) // bright node at each grid vertex
	}

	if lv < 0.01 {
		return 0
	}

	// Peak glow: exponential — valleys near-black, peaks brilliant white
	peakGlow := 0.10 + math.Pow(h, 2.0)*0.90
	if h > 0.68 {
		peakGlow = math.Min(1.0, peakGlow+(h-0.68)*3.2)
	}

	return math.Min(1.0, lv*peakGlow*fade*scale)
}

func renderBarcodeGlitch(nx, ny, t, density, scale, chaos float64) (float64, string) {
	intensity := 0.0
	resultChar := ""

	segments := int(40 + density*60)
	posX := (nx + 1.0) * 0.5
	barWidth := 1.0 / float64(segments)

	segIdx := int(posX / barWidth)
	if segIdx >= segments {
		segIdx = segments - 1
	}
	if segIdx < 0 {
		segIdx = 0
	}

	hash1 := math.Sin(float64(segIdx)*12.9898+t*0.05) * 43758.5453
	hash1 = hash1 - math.Floor(hash1)
	hash2 := math.Sin(float64(segIdx)*78.233+t*0.02) * 43758.5453
	hash2 = hash2 - math.Floor(hash2)

	isBar := hash1 > 0.35
	widthVar := hash2*0.5 + 0.5

	glitchActive := false
	glitchOffset := 0.0

	glitchZone := math.Sin(t*3+ny*8) * chaos
	if glitchZone > 0.6 {
		glitchActive = true
		glitchOffset = math.Sin(t*15+float64(segIdx)*0.5) * 0.2 * chaos
	}

	vGlitch := math.Sin(t*5+nx*15) * chaos
	if vGlitch > 0.7 {
		glitchActive = true
	}

	if glitchActive {
		displaced := posX + glitchOffset
		displaced = math.Mod(displaced+1, 1)
		newSegIdx := int(displaced / barWidth)
		if newSegIdx >= 0 && newSegIdx < segments {
			newHash := math.Sin(float64(newSegIdx)*12.9898+t*0.05) * 43758.5453
			newHash = newHash - math.Floor(newHash)
			isBar = newHash > 0.35
		}
	}

	scanY := math.Mod(t*1.5, 2.0) - 1.0
	scanDist := math.Abs(ny - scanY)
	scanIntensity := 0.0
	if scanDist < 0.04 {
		scanIntensity = (0.04 - scanDist) * 25 * scale
	}

	barHeight := 0.85 + hash2*0.1
	inBarHeight := math.Abs(ny) < barHeight

	if isBar && inBarHeight {
		intensity = 0.9 * scale * widthVar
		if glitchActive && rand.Float64() > 0.6 {
			glitchChars := []string{"█", "▓", "▒", "░", "│", "┃", "║", "▌", "▐", "▀", "▄"}
			resultChar = glitchChars[rand.Intn(len(glitchChars))]
		} else {
			if widthVar > 0.8 {
				resultChar = "█"
			} else if widthVar > 0.6 {
				resultChar = "▓"
			} else if widthVar > 0.4 {
				resultChar = "│"
			} else {
				resultChar = "┃"
			}
		}
	} else {
		intensity = 0.05
		if glitchActive && rand.Float64() > 0.92 {
			noiseChars := []string{".", ":", ";", ",", "'", "`"}
			resultChar = noiseChars[rand.Intn(len(noiseChars))]
			intensity = 0.25
		}
	}

	intensity += scanIntensity

	if ny > 0.88 {
		digitWidth := 1.0 / 13.0
		digitIdx := int(posX / digitWidth)
		if digitIdx >= 0 && digitIdx < 13 {
			digit := int(math.Abs(math.Sin(float64(digitIdx)*7.5+t*0.005)*10)) % 10
			if glitchActive && rand.Float64() > 0.4 {
				digit = rand.Intn(10)
			}
			resultChar = string(rune('0' + digit))
			intensity = 0.7
		}
	}

	if chaos > 0.5 && rand.Float64() > 0.97 {
		corruptChars := []string{"X", "#", "@", "!", "?", "*", "&", "%", "E", "R"}
		resultChar = corruptChars[rand.Intn(len(corruptChars))]
		intensity = 1.0
	}

	return intensity, resultChar
}

func renderWaveInterference(nx, ny, t, density, scale, chaos float64) float64 {
	wave1 := math.Sin(nx*density*15 + t*2)
	wave2 := math.Sin(ny*density*12 - t*1.5)
	wave3 := math.Sin((nx+ny)*density*10 + t)
	wave4 := math.Cos((nx-ny)*density*8 - t*0.5)
	interference := math.Sin(nx*ny*50*chaos + t)

	return (wave1 + wave2 + wave3 + wave4 + interference*0.3) / 4.3 * scale
}

func renderHyperSpiral(nx, ny, t, density, scale, chaos float64) float64 {
	r := math.Sqrt(nx*nx + ny*ny)
	theta := math.Atan2(ny, nx)

	arms := 3 + chaos*5
	spiral := math.Sin(r*density*25 - t*5 + theta*arms)
	spiral2 := math.Sin(r*density*15 + t*3 - theta*(arms+2))

	return math.Tanh((spiral + spiral2*0.5) * 2 * scale)
}

func renderVoid(nx, ny, t, density, scale, chaos float64) float64 {
	r := math.Sqrt(nx*nx + ny*ny)
	void := 1.0 - math.Tanh(r*3*scale)

	theta := math.Atan2(ny, nx)
	disk := math.Sin(theta*8-t*3+r*20) * math.Exp(-math.Abs(r-0.4)*10) * density
	ripple := math.Sin(r*40-t*2) * 0.1 * (1 - r) * chaos

	if rand.Float64() > 0.998 && r < 0.3 {
		return 1.0
	}

	return void*0.8 + disk*0.5 + ripple
}

func renderGlitchStorm(nx, ny, t, density, scale, chaos float64) float64 {
	blockSize := 5 + chaos*15
	blockX := math.Floor(nx * blockSize)
	blockY := math.Floor(ny * blockSize)

	if noise2D(blockX+math.Floor(t*0.8), blockY+math.Floor(t*0.3)) > 0.75-density*0.25 {
		return rand.Float64()
	}

	scanline := 0.0
	if math.Mod(ny*50+t*10, 2) < 0.3 {
		scanline = 0.3
	}

	static := noise2D(nx*100+t, ny*100) * 0.3 * chaos

	if rand.Float64() > 0.998 {
		return 1.0
	}

	return (scanline + static) * scale
}

func renderChaosEntropy(nx, ny, t, density, scale, chaos float64) float64 {
	n1 := noise2D(nx*5*scale+t*0.5, ny*5*scale)
	n2 := noise2D(nx*10*scale-t*0.3, ny*10*scale+t*0.2) * 0.5
	n3 := noise2D(nx*20*scale+t*0.1, ny*20*scale-t*0.4) * 0.25
	n4 := noise2D(nx*40*scale, ny*40*scale+t*0.2) * 0.125

	result := (n1 + n2 + n3 + n4) / 1.875

	if rand.Float64() > 1-chaos*0.1 {
		result = rand.Float64()
	}

	entropy := math.Sin(nx*30+t*2) * math.Sin(ny*30-t*1.5) * chaos * 0.3

	return (result + entropy) * density * 1.5
}

func renderParticleWaveField(nx, ny, t, density, scale, chaos float64) (float64, string) {
	dotDensity := 30.0 + density*40

	warpX := nx + math.Sin(ny*5+t)*0.08*chaos + math.Sin(ny*12+t*1.5)*0.03*chaos
	warpY := ny + math.Sin(nx*4+t*0.8)*0.06*chaos + math.Cos(nx*10+t*1.2)*0.02*chaos

	gridX := math.Floor(warpX * dotDensity)
	gridY := math.Floor(warpY * dotDensity)

	cellX := (warpX * dotDensity) - gridX
	cellY := (warpY * dotDensity) - gridY

	gx := gridX / dotDensity
	gy := gridY / dotDensity

	wave := 0.0
	wave += math.Sin(gx*8*scale+t*1.2) * math.Cos(gy*6*scale-t*0.9) * 0.4
	wave += math.Sin((gx+gy)*5*scale+t*0.7) * 0.3
	wave += math.Cos((gx-gy)*7*scale-t*1.1) * 0.25
	wave += math.Sin(gx*15*scale+gy*12*scale+t*1.8) * 0.2
	wave += math.Cos(gx*20*scale-gy*18*scale-t*1.5) * 0.15

	terrain := math.Sin(gx*4*scale) * math.Sin(gy*3*scale+t*0.3)
	terrain += math.Sin(gx*2*scale+gy*2.5*scale-t*0.2) * 0.5
	wave += terrain * 0.3

	ripple := math.Sin(gx*25*scale+t*3) * math.Sin(gy*22*scale-t*2.5) * 0.1 * chaos
	wave += ripple

	waveNorm := (wave + 1.3) / 2.6
	waveNorm = math.Max(0, math.Min(1, waveNorm))

	offsetX := 0.5 + math.Sin(gx*10+gy*8+t*2)*0.12*chaos
	offsetY := 0.5 + math.Cos(gx*8+gy*10+t*1.5)*0.12*chaos

	dx := cellX - offsetX
	dy := cellY - offsetY
	dist := math.Sqrt(dx*dx + dy*dy)

	baseRadius := 0.32 + waveNorm*0.18
	radius := baseRadius

	intensity := 0.0
	resultChar := " "

	if dist < radius {
		falloff := 1.0 - (dist / radius)
		falloff = math.Pow(falloff, 0.4)
		intensity = falloff * (0.3 + waveNorm*0.7)

		if intensity > 0.75 {
			resultChar = "●"
		} else if intensity > 0.55 {
			resultChar = "◉"
		} else if intensity > 0.4 {
			resultChar = "○"
		} else if intensity > 0.25 {
			resultChar = "◦"
		} else if intensity > 0.12 {
			resultChar = "·"
		} else {
			resultChar = "."
		}

		if chaos > 0.5 && rand.Float64() > 0.97 {
			altChars := []string{"•", "∘", "°", "⋅"}
			resultChar = altChars[rand.Intn(len(altChars))]
		}
	} else if dist < radius*1.5 && waveNorm > 0.6 {
		intensity = 0.1 * (1 - (dist-radius)/(radius*0.5))
		resultChar = "·"
	}

	return intensity, resultChar
}

func updateBalls(speed float64) {
	for i := range balls {
		balls[i].X += balls[i].VX * speed * 2
		balls[i].Y += balls[i].VY * speed * 2

		if balls[i].X > 1 || balls[i].X < -1 {
			balls[i].VX *= -1
			balls[i].X = math.Max(-1, math.Min(1, balls[i].X))
		}
		if balls[i].Y > 1 || balls[i].Y < -1 {
			balls[i].VY *= -1
			balls[i].Y = math.Max(-1, math.Min(1, balls[i].Y))
		}

		if rand.Float64() > 0.95 {
			balls[i].VX += (rand.Float64() - 0.5) * 0.02
			balls[i].VY += (rand.Float64() - 0.5) * 0.02
			balls[i].VX = math.Max(-0.05, math.Min(0.05, balls[i].VX))
			balls[i].VY = math.Max(-0.05, math.Min(0.05, balls[i].VY))
		}
	}
}

// ============================================
// MAIN RENDER
// ============================================

func render() {
	width, height := getTerminalSize()
	var output strings.Builder

	output.WriteString("\033[H")

	// Update transitions
	updateLayerAlphas()
	updateAutoMode()

	// Update physics
	if layers[2].Alpha > 0.01 {
		updateBalls(layers[2].Params[0].Value)
	}

	// UI dimensions
	boxW := 30
	boxH := 11
	boxX := width - boxW - 2
	boxY := 1

	renderH := height - 3
	t := float64(frame) * 0.1

	for y := 0; y < renderH; y++ {
		for x := 0; x < width; x++ {
			inBox := (x >= boxX && x < boxX+boxW && y >= boxY && y < boxY+boxH)

			if inBox {
				relX := x - boxX
				relY := y - boxY

				if relY == 0 {
					if relX == 0 {
						output.WriteString("\033[97m┌")
					} else if relX == boxW-1 {
						output.WriteString("┐\033[0m")
					} else {
						output.WriteString("─")
					}
				} else if relY == boxH-1 {
					if relX == 0 {
						output.WriteString("\033[97m└")
					} else if relX == boxW-1 {
						output.WriteString("┘\033[0m")
					} else {
						output.WriteString("─")
					}
				} else if relX == 0 {
					output.WriteString("\033[97m│\033[0m")
				} else if relX == boxW-1 {
					output.WriteString("\033[97m│\033[0m")
				} else {
					if relY == 1 {
						title := layers[selectedLayer].Name
						if len(title) > boxW-4 {
							title = title[:boxW-4]
						}
						if relX >= 2 && relX < 2+len(title) {
							output.WriteString("\033[93m")
							output.WriteByte(title[relX-2])
							output.WriteString("\033[0m")
						} else {
							output.WriteString(" ")
						}
					} else if relY >= 2 && relY <= 6 {
						paramIdx := relY - 2
						if paramIdx < len(layers[selectedLayer].Params) {
							p := layers[selectedLayer].Params[paramIdx]
							line := fmt.Sprintf("%s:", p.Name)
							sliderLen := 8
							filled := int(p.Value * float64(sliderLen))

							if relX == 1 {
								if paramIdx == selectedParam {
									output.WriteString("\033[93m>\033[0m")
								} else {
									output.WriteString(" ")
								}
							} else if relX >= 2 && relX < 2+len(line) {
								output.WriteString("\033[97m")
								output.WriteByte(line[relX-2])
								output.WriteString("\033[0m")
							} else if relX == 2+len(line) {
								output.WriteString("\033[90m(\033[0m")
							} else if relX > 2+len(line) && relX <= 2+len(line)+sliderLen {
								pos := relX - 3 - len(line)
								if pos < filled {
									output.WriteString("\033[96m█\033[0m")
								} else {
									output.WriteString("\033[90m░\033[0m")
								}
							} else if relX == 3+len(line)+sliderLen {
								output.WriteString("\033[90m)\033[0m")
							} else {
								output.WriteString(" ")
							}
						} else {
							output.WriteString(" ")
						}
					} else if relY == 7 {
						// Alpha display
						alphaLine := fmt.Sprintf("Alpha: %.0f%%", layers[selectedLayer].Alpha*100)
						if relX >= 2 && relX < 2+len(alphaLine) {
							output.WriteString("\033[96m")
							output.WriteByte(alphaLine[relX-2])
							output.WriteString("\033[0m")
						} else {
							output.WriteString(" ")
						}
					} else if relY == 8 {
						// Mode display
						modeStr := "Mode: MANUAL"
						if transition.AutoMode {
							modeStr = fmt.Sprintf("Mode: AUTO [%d→%d]", transition.CurrentLayer, transition.NextLayer)
						}
						if relX >= 2 && relX < 2+len(modeStr) {
							if transition.AutoMode {
								output.WriteString("\033[92m")
							} else {
								output.WriteString("\033[93m")
							}
							output.WriteByte(modeStr[relX-2])
							output.WriteString("\033[0m")
						} else {
							output.WriteString(" ")
						}
					} else if relY == 9 {
						// Controls hint
						hint := "[Enter]=Auto [Space]=Toggle"
						if relX >= 2 && relX < 2+len(hint) {
							output.WriteString("\033[90m")
							output.WriteByte(hint[relX-2])
							output.WriteString("\033[0m")
						} else {
							output.WriteString(" ")
						}
					} else {
						output.WriteString(" ")
					}
				}
			} else {
				nx := (float64(x)/float64(width))*2 - 1
				ny := (float64(y)/float64(renderH))*2 - 1
				nx *= float64(width) / float64(renderH) * 0.5

				totalVal := 0.0
				dominantLayer := 0
				maxContrib := 0.0
				specialChar := ""

				for lIdx := 1; lIdx <= 13; lIdx++ {
					lyr := layers[lIdx]
					
					// Skip if alpha is essentially zero
					if lyr.Alpha < 0.01 {
						continue
					}

					speed := lyr.Params[0].Value
					density := lyr.Params[1].Value
					bright := lyr.Params[2].Value
					scale := lyr.Params[3].Value + 0.5
					chaos := lyr.Params[4].Value
					lt := t * speed

					var val float64
					var tempChar string

					switch lIdx {
					case 1:
						val = renderPortalLooping(nx, ny, lt, density, scale, chaos)
					case 2:
						v, c := renderBallMatrix(nx, ny, lt, density, scale, chaos)
						val = v
						tempChar = c
					case 3:
						v, c := renderTopeng(nx, ny, lt, density, scale, chaos, width, renderH)
						val = v
						if c != ' ' && c != 0 {
							tempChar = string(c)
						}
					case 4:
						v, c := renderTypistMatrix(nx, ny, lt, density, scale, chaos, width, renderH)
						val = v
						if c != ' ' && c != 0 {
							tempChar = string(c)
						}
					case 5:
						val = renderTubuhYangMengingat(nx, ny, lt, density, scale, chaos)
					case 6:
						val = renderMountain3D(nx, ny, lt, density, scale, chaos)
					case 7:
						v, c := renderBarcodeGlitch(nx, ny, lt, density, scale, chaos)
						val = v
						tempChar = c
					case 8:
						val = renderWaveInterference(nx, ny, lt, density, scale, chaos)
					case 9:
						val = renderHyperSpiral(nx, ny, lt, density, scale, chaos)
					case 10:
						val = renderVoid(nx, ny, lt, density, scale, chaos)
					case 11:
						val = renderGlitchStorm(nx, ny, lt, density, scale, chaos)
					case 12:
						val = renderChaosEntropy(nx, ny, lt, density, scale, chaos)
					case 13:
						v, c := renderParticleWaveField(nx, ny, lt, density, scale, chaos)
						val = v
						tempChar = c
					}

					val *= bright
					
					// Apply layer alpha (fade effect)
					val *= lyr.Alpha

					if tempChar != "" && tempChar != " " && val > 0.05 {
						if lyr.Alpha > 0.5 || specialChar == "" {
							specialChar = tempChar
						}
					}

					totalVal += val

					if math.Abs(val) > maxContrib {
						maxContrib = math.Abs(val)
						dominantLayer = lIdx
					}
				}

				totalVal *= globalBright

				// Dim special char if alpha is low
				if specialChar != "" && totalVal > 0.05 {
					if colorMode && dominantLayer > 0 {
						intensity := math.Min(1.0, math.Abs(totalVal))
						if intensity > 0.7 {
							output.WriteString("\033[97m")
						} else if intensity > 0.4 {
							output.WriteString(layerColors[dominantLayer])
						} else {
							output.WriteString("\033[90m")
						}
					}
					output.WriteString(specialChar)
					output.WriteString("\033[0m")
				} else {
					idx := int((math.Abs(totalVal)+0.15)*float64(len(gradients)-1)) - 1
					if idx < 0 {
						idx = 0
					}
					if idx >= len(gradients) {
						idx = len(gradients) - 1
					}

					char := gradients[idx]
					if colorMode && dominantLayer > 0 && idx > 1 {
						output.WriteString(layerColors[dominantLayer] + char + "\033[0m")
					} else {
						output.WriteString(char)
					}
				}
			}
		}
		if y < renderH-1 {
			output.WriteString("\n")
		}
	}

	// Layer selector
	output.WriteString("\n\n")

	spacing := width / 14
	if spacing < 4 {
		spacing = 4
	}

	for i := 1; i <= 13; i++ {
		if i == 1 {
			output.WriteString(strings.Repeat(" ", spacing/2))
		} else {
			gap := spacing - 3
			if gap < 1 {
				gap = 1
			}
			output.WriteString(strings.Repeat(" ", gap))
		}

		displayNum := ""
		if i <= 9 {
			displayNum = fmt.Sprintf("%d", i)
		} else if i == 10 {
			displayNum = "0"
		} else if i == 11 {
			displayNum = "-"
		} else if i == 12 {
			displayNum = "="
		} else if i == 13 {
			displayNum = "+"
		}

		alpha := layers[i].Alpha
		
		if i == selectedLayer {
			output.WriteString(fmt.Sprintf("\033[93m(%s)\033[0m", displayNum))
		} else if alpha > 0.5 {
			output.WriteString(fmt.Sprintf("\033[92m[%s]\033[0m", displayNum))
		} else if alpha > 0.01 {
			// Transitioning - show partial
			output.WriteString(fmt.Sprintf("\033[96m<%s>\033[0m", displayNum))
		} else {
			output.WriteString(fmt.Sprintf("\033[90m{%s}\033[0m", displayNum))
		}
	}

	fmt.Print(output.String())
}

// ============================================
// INPUT HANDLER
// ============================================

func handleInput() {
	if err := keyboard.Open(); err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer keyboard.Close()

	for running {
		char, key, err := keyboard.GetKey()
		if err != nil {
			continue
		}

		switch key {
		case keyboard.KeyEsc:
			running = false
			return
		case keyboard.KeyArrowUp:
			selectedParam = (selectedParam + 4) % 5
		case keyboard.KeyArrowDown:
			selectedParam = (selectedParam + 1) % 5
		case keyboard.KeyArrowRight:
			layers[selectedLayer].Params[selectedParam].Value = math.Min(1.0, layers[selectedLayer].Params[selectedParam].Value+0.05)
		case keyboard.KeyArrowLeft:
			layers[selectedLayer].Params[selectedParam].Value = math.Max(0.0, layers[selectedLayer].Params[selectedParam].Value-0.05)
		case keyboard.KeySpace:
			// Toggle layer with fade transition
			if layers[selectedLayer].TargetAlpha > 0.5 {
				layers[selectedLayer].TargetAlpha = 0
				layers[selectedLayer].Enabled = false
			} else {
				layers[selectedLayer].TargetAlpha = 1
				layers[selectedLayer].Enabled = true
			}
		case keyboard.KeyEnter:
			// Toggle auto mode
			transition.AutoMode = !transition.AutoMode
			if transition.AutoMode {
				// Reset all layers
				for i := 1; i <= 13; i++ {
					layers[i].Alpha = 0
					layers[i].TargetAlpha = 0
					layers[i].Enabled = false
				}
				// Start from selected layer
				transition.CurrentLayer = selectedLayer
				transition.NextLayer = selectedLayer + 1
				if transition.NextLayer > 13 {
					transition.NextLayer = 1
				}
				layers[transition.CurrentLayer].Alpha = 1
				layers[transition.CurrentLayer].TargetAlpha = 1
				layers[transition.CurrentLayer].Enabled = true
				transition.HoldTimer = 0
			}
		}

		switch char {
		case 'q', 'Q':
			running = false
			return
		case '1', '2', '3', '4', '5', '6', '7', '8', '9':
			newLayer := int(char - '0')
			selectedLayer = newLayer
			selectedParam = 0
			if transition.AutoMode {
				// Jump to this layer in auto mode
				for i := 1; i <= 13; i++ {
					layers[i].TargetAlpha = 0
				}
				layers[newLayer].TargetAlpha = 1
				transition.CurrentLayer = newLayer
				transition.NextLayer = newLayer + 1
				if transition.NextLayer > 13 {
					transition.NextLayer = 1
				}
				transition.HoldTimer = 0
			}
		case '0':
			selectedLayer = 10
			selectedParam = 0
		case '-':
			selectedLayer = 11
			selectedParam = 0
		case '=':
			selectedLayer = 12
			selectedParam = 0
		case '+':
			selectedLayer = 13
			selectedParam = 0
		case 'a', 'A':
			// Same as space - toggle with fade
			if layers[selectedLayer].TargetAlpha > 0.5 {
				layers[selectedLayer].TargetAlpha = 0
				layers[selectedLayer].Enabled = false
			} else {
				layers[selectedLayer].TargetAlpha = 1
				layers[selectedLayer].Enabled = true
			}
		case 'c', 'C':
			colorMode = !colorMode
		case '[':
			globalBright = math.Max(0.1, globalBright-0.1)
		case ']':
			globalBright = math.Min(2.0, globalBright+0.1)
		case 'r', 'R':
			initBalls(60)
			initTypist(100)
		case 'f', 'F':
			// Faster transitions
			transition.TransitionSpeed = math.Min(0.2, transition.TransitionSpeed+0.01)
		case 's', 'S':
			// Slower transitions
			transition.TransitionSpeed = math.Max(0.01, transition.TransitionSpeed-0.01)
		case 'h', 'H':
			// Longer hold duration
			transition.HoldDuration = math.Min(500, transition.HoldDuration+25)
		case 'g', 'G':
			// Shorter hold duration
			transition.HoldDuration = math.Max(25, transition.HoldDuration-25)
		case 'n', 'N':
			// Skip to next sequence immediately (auto mode)
			if transition.AutoMode {
				transition.HoldTimer = transition.HoldDuration
			}
		case 'x', 'X':
			// Solo mode - enable only selected, disable others with fade
			for i := 1; i <= 13; i++ {
				if i == selectedLayer {
					layers[i].TargetAlpha = 1
					layers[i].Enabled = true
				} else {
					layers[i].TargetAlpha = 0
					layers[i].Enabled = false
				}
			}
		}
	}
}

// ============================================
// MAIN
// ============================================

func main() {
	rand.Seed(time.Now().UnixNano())
	initLayers()

	fmt.Print("\033[?25l")
	fmt.Print("\033[2J")
	defer func() {
		fmt.Print("\033[?25h")
		fmt.Print("\033[0m")
		fmt.Print("\033[2J")
		fmt.Print("\033[H")
	}()

	go handleInput()

	ticker := time.NewTicker(40 * time.Millisecond)
	defer ticker.Stop()

	for running {
		<-ticker.C
		render()
		frame++
	}
}
