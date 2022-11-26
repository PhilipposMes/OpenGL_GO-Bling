package main

import (
	"bufio"
	"math/rand"

	"fmt"

	//"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	gl "github.com/chsc/gogl/gl33"
	"github.com/veandco/go-sdl2/sdl"
)

func createprogram() gl.Uint {
	// VERTEX SHADER
	vs := gl.CreateShader(gl.VERTEX_SHADER)
	vs_source := gl.GLString(vertexShaderSource)
	gl.ShaderSource(vs, 1, &vs_source, nil)
	gl.CompileShader(vs)
	var vs_status gl.Int
	gl.GetShaderiv(vs, gl.COMPILE_STATUS, &vs_status)
	fmt.Printf("Compiled Vertex Shader: %v\n", vs_status)

	// FRAGMENT SHADER
	fs := gl.CreateShader(gl.FRAGMENT_SHADER)
	fs_source := gl.GLString(fragmentShaderSource)
	gl.ShaderSource(fs, 1, &fs_source, nil)
	gl.CompileShader(fs)
	var fstatus gl.Int
	gl.GetShaderiv(fs, gl.COMPILE_STATUS, &fstatus)
	fmt.Printf("Compiled Fragment Shader: %v\n", fstatus)

	// CREATE PROGRAM
	program := gl.CreateProgram()
	gl.AttachShader(program, vs)
	gl.AttachShader(program, fs)
	fragoutstring := gl.GLString("outColor")
	defer gl.GLStringFree(fragoutstring)
	gl.BindFragDataLocation(program, gl.Uint(0), fragoutstring)

	gl.LinkProgram(program)
	var linkstatus gl.Int
	gl.GetProgramiv(program, gl.LINK_STATUS, &linkstatus)
	fmt.Printf("Program Link: %v\n", linkstatus)

	return program
}

var uniRoll float32 = 0.0
var uniYaw float32 = 1.0
var uniPitch float32 = 0.0
var uniscale float32 = 0.3
var yrot float32 = 20.0
var zrot float32 = 0.0
var xrot float32 = 0.0
var UniScale gl.Int

// light
var UniLight gl.Int
var l1 gl.Float = 2
var l2 gl.Float = 1
var l3 gl.Float = 4

// Colour
var UniColour gl.Int
var c1 gl.Float = 1
var c2 gl.Float = 1
var c3 gl.Float = 1

// path of file
var path string

//var path = "shape.obj"
//var path = "untitled.obj"

// arrays for v , f,a the lines of the file
// var v []float32
var v []float64
var f []int
var allLines []string

// project
var colorchange float64 = 0.1
var fvn []int
var vn []float64

func scanner(path string) ([]string, error) {
	readFile, err := os.Open(path) //read file & path
	if err != nil {
		fmt.Println(err)
	}
	defer readFile.Close() //close file

	fileScanner := bufio.NewScanner(readFile) //reading file
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() { //scaning file
		allLines = append(allLines, fileScanner.Text())
	}
	return allLines, fileScanner.Err()
}

func main() {
	var window *sdl.Window
	var context sdl.GLContext
	var event sdl.Event
	var running bool
	var err error

	// // scanner()
	path := os.Args
	allLines, err := scanner(path[1])
	fmt.Println("OBJ file:", path[1])

	// scanner()
	//allLines, err := scanner(path)

	if err != nil {
		// ... handle error

		panic(err)
	}
	for _, lines := range allLines {
		//selecting (scaning) lines that start with v or f
		forv := strings.HasPrefix(lines, "v ")
		//fmt.Println(forv)
		forf := strings.HasPrefix(lines, "f")
		//read := strings.Fields(lines)

		//project
		forvn := strings.HasPrefix(lines, "vn ")

		if forv == true { // if v in line
			//addlinev := strings.Split(lines, "/n")
			addlinev := strings.Fields(lines) // spliting the lines with v
			//for _, add := range addlinev { //for range of lines with v
			//v = append(v, add[2:]) //add to array the line starting with the 2nd character
			//}

			v1 := strings.Fields(addlinev[1])
			v2 := strings.Fields(addlinev[2])
			v3 := strings.Fields(addlinev[3])
			if floatvline1, err := strconv.ParseFloat(v1[0], 64); err == nil {
				v = append(v, floatvline1)
			}
			if floatvline2, err := strconv.ParseFloat(v2[0], 64); err == nil {
				v = append(v, floatvline2)
			}
			if floatvline3, err := strconv.ParseFloat(v3[0], 64); err == nil {
				v = append(v, floatvline3)
			}

		}

		if forf == true { // if f in line
			//addlinef := strings.Split(lines, "/n") // spliting the lines with f
			addlinef := strings.Fields(lines) // spliting the lines with f

			//for _, add := range addlinef { //for range of lines with f

			//f = append(f, add[2:]) //add to array the line starting with the 2nd character

			f1 := strings.Split(addlinef[1], "/")
			f2 := strings.Split(addlinef[2], "/")
			f3 := strings.Split(addlinef[3], "/")

			//fmt.Println("check", f1, f2, f3)

			int1, err := strconv.Atoi(f1[0])
			int2, err := strconv.Atoi(f2[0])
			int3, err := strconv.Atoi(f3[0])

			//fmt.Println("check2", int1, int2, int3)
			f = append(f, int1, int2, int3)

			//project

			int11, err := strconv.Atoi(f1[2])
			int22, err := strconv.Atoi(f2[2])
			int33, err := strconv.Atoi(f3[2])
			//fmt.Println("check2vn", int11, int22, int33)

			if err != nil {
				panic(err)
			}

			fvn = append(fvn, int11, int22, int33)

			//fmt.Println("chewck ffffffffff", f)
			//}
			//f = append(f, addlinef...)
		}
		//project
		if forvn == true { // if vn in line
			addlinevn := strings.Fields(lines) // spliting the lines with vn

			vn1 := strings.Fields(addlinevn[1])
			vn2 := strings.Fields(addlinevn[2])
			vn3 := strings.Fields(addlinevn[3])
			if floatvnline1, err := strconv.ParseFloat(vn1[0], 64); err == nil {
				vn = append(vn, floatvnline1)
			}
			if floatvnline2, err := strconv.ParseFloat(vn2[0], 64); err == nil {
				vn = append(vn, floatvnline2)
			}
			if floatvnline3, err := strconv.ParseFloat(vn3[0], 64); err == nil {
				vn = append(vn, floatvnline3)
			}

			//fmt.Println("vn", vn1, vn2, vn3) //test vn
		}
		for _, fvnnumber := range fvn {

			vnnumber1 := fvnnumber*3 - 3
			vnnumber2 := fvnnumber*3 - 2
			vnnumber3 := fvnnumber*3 - 1

			tvn1 := gl.Float(vn[vnnumber1])
			tvn2 := gl.Float(vn[vnnumber2])
			tvn3 := gl.Float(vn[vnnumber3])

			triangle_normals = append(triangle_normals, tvn1, tvn2, tvn3)

			//fmt.Println("vnn", fvnnumber, tvn1, tvn2, tvn3) //test vn
		}
		//appending everything
		for _, fnumber := range f {

			vnumber1 := fnumber*3 - 3 //*3 for  line -3 if first -2 for second ...
			//fmt.Println("no1", f[1], num1)

			vnumber2 := fnumber*3 - 2
			//fmt.Println(num2)
			vnumber3 := fnumber*3 - 1

			tv1 := gl.Float(v[vnumber1])
			//fmt.Println("no1", num1, t1)
			tv2 := gl.Float(v[vnumber2])
			//fmt.Println("no2", num2, t2)
			tv3 := gl.Float(v[vnumber3])

			triangle_vertices = append(triangle_vertices, tv1, tv2, tv3)
			triangle_colours = append(triangle_colours, 1, 1, 1)
			//triangle_normals = append(triangle_normals,)

			//fmt.Println("f", fnumber)  //test f
		}
	}

	runtime.LockOSThread()
	if err = sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()
	window, err = sdl.CreateWindow(winTitle, sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		winWidth, winHeight, sdl.WINDOW_OPENGL)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()
	context, err = window.GLCreateContext()
	if err != nil {
		panic(err)
	}
	defer sdl.GLDeleteContext(context)

	gl.Init()
	gl.Viewport(0, 0, gl.Sizei(winWidth), gl.Sizei(winHeight))
	// OPENGL FLAGS
	gl.ClearColor(0.0, 0.1, 0.0, 1.0)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.SRC_ALPHA, gl.ONE_MINUS_SRC_ALPHA)

	// VERTEX BUFFER
	var vertexbuffer gl.Uint
	gl.GenBuffers(1, &vertexbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(len(triangle_vertices)*4), gl.Pointer(&triangle_vertices[0]), gl.STATIC_DRAW)

	// COLOUR BUFFER
	var colourbuffer gl.Uint
	gl.GenBuffers(1, &colourbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, colourbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(len(triangle_colours)*4), gl.Pointer(&triangle_colours[0]), gl.STATIC_DRAW)

	// NORMAL BUFFER
	var normalbuffer gl.Uint
	gl.GenBuffers(1, &normalbuffer)
	gl.BindBuffer(gl.ARRAY_BUFFER, normalbuffer)
	gl.BufferData(gl.ARRAY_BUFFER, gl.Sizeiptr(len(triangle_normals)*4),
		gl.Pointer(&triangle_normals[0]), gl.STATIC_DRAW)

	// GUESS WHAT
	program := createprogram()

	// VERTEX ARRAY
	var VertexArrayID gl.Uint
	gl.GenVertexArrays(1, &VertexArrayID)
	gl.BindVertexArray(VertexArrayID)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, vertexbuffer)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, gl.FALSE, 0, nil)

	// VERTEX ARRAY HOOK COLOURS
	gl.EnableVertexAttribArray(1)
	gl.BindBuffer(gl.ARRAY_BUFFER, colourbuffer)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, gl.FALSE, 0, nil)

	gl.EnableVertexAttribArray(2)
	gl.BindBuffer(gl.ARRAY_BUFFER, normalbuffer)
	gl.VertexAttribPointer(2, 3, gl.FLOAT, gl.FALSE, 0, nil)

	//UNIFORM HOOK
	unistring := gl.GLString("scaleMove")
	UniScale = gl.GetUniformLocation(program, unistring)
	fmt.Printf("Uniform Link: %v\n", UniScale+1)

	//UNIFORM HOOK
	unistring1 := gl.GLString("LightSource")
	UniLight = gl.GetUniformLocation(program, unistring1)
	fmt.Printf("Uniform light Link: %v\n", UniLight+1)

	//UNIFORM HOOK
	unistring2 := gl.GLString("ChangeColour")
	UniColour = gl.GetUniformLocation(program, unistring2)
	fmt.Printf("Uniform light Link: %v\n", UniColour+1)

	gl.UseProgram(program)

	running = true
	for running {
		for event = sdl.PollEvent(); event != nil; event =
			sdl.PollEvent() {
			switch t := event.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.KeyboardEvent:
				keyCode := t.Keysym.Sym
				keys := ""
				switch t.Keysym.Mod {
				case sdl.KMOD_LCTRL:
					keys += "CTRL"
					fmt.Println("keys")
					if keyCode < 10000000 {
						if keys != "" {
							keys += " + "
						}
					} else {
						if t.State == sdl.RELEASED {
							keys += string(keyCode) + " released"
							//drawgl()
							//window.GLSwap()
						} else if t.State == sdl.PRESSED {
							keys += string(keyCode) + " pressed"

							c1 = gl.Float(rand.Float64())
							c2 = gl.Float(rand.Float64())
							c3 = gl.Float(rand.Float64())
							fmt.Println("Colour 1:", c1, "Colour 2:", c2, "Colour 3:", c3)

							// for _, i := range f {
							// 	triangle_colours[i] = triangle_colours[i] + 1.5
							// 	fmt.Println(triangle_colours[1])
							// }
						}
					}

				case sdl.KMOD_LSHIFT:
					keys += "SHIFT"
					fmt.Println("keys")
					if keyCode < 10000000 {
						if keys != "" {
							keys += " + "
						}
					} else {
						if t.State == sdl.RELEASED {
							keys += string(keyCode) + " released"
						} else if t.State == sdl.PRESSED {
							keys += string(keyCode) + " pressed"

							c1 = 1
							c2 = 1
							c3 = 1
							fmt.Println("Turned white")
						}
					}
				case sdl.KMOD_CAPS:
					keys += "CAPS"
					fmt.Println("keys")
					if keyCode < 10000000 {
						if keys != "" {
							keys += " + "
						}
					} else {
						if t.State == sdl.RELEASED {
							keys += string(keyCode) + " released"
						} else if t.State == sdl.PRESSED {
							keys += string(keyCode) + " pressed"

							l1 = gl.Float(rand.Float64())
							l2 = gl.Float(rand.Float64())
							l3 = gl.Float(rand.Float64())
							fmt.Println("Lightposition 1:", c1, "Lightposition 2:", c2, "Lightposition 3:", c3)
						}
					}
				}
				if keys != "" {
					fmt.Println(keys)
				}

				// etc

			case *sdl.MouseMotionEvent:

				xrot = float32(t.Y) / 2
				yrot = float32(t.X) / 2
				fmt.Printf("[%dms]MouseMotion\tid:%d\tx:%d\ty:%d\txrel:%d\tyrel:%d\n", t.Timestamp, t.Which, t.X, t.Y, t.XRel, t.YRel)
			}

		}
		drawgl()
		window.GLSwap()

	}

}

func drawgl() {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)

	uniYaw = yrot * (math.Pi / 180.0)
	yrot = yrot - 1.0
	uniPitch = zrot * (math.Pi / 180.0)
	zrot = zrot - 0.5
	uniRoll = xrot * (math.Pi / 180.0)
	xrot = xrot - 0.2

	gl.Uniform4f(UniScale, gl.Float(uniRoll), gl.Float(uniYaw), gl.Float(uniPitch), gl.Float(uniscale))
	//gl.Uniform3f(UniLight, gl.Float(2), gl.Float(1), gl.Float(4)) //light position
	gl.Uniform3f(UniLight, gl.Float(l1), gl.Float(l2), gl.Float(l3)) //light position
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.DrawArrays(gl.TRIANGLES, gl.Int(0), gl.Sizei(len(triangle_vertices)*4))
	gl.Uniform3f(UniColour, c1, c2, c3) //colour

	time.Sleep(50 * time.Millisecond)

}

const (
	winTitle           = "OpenGL Shader"
	winWidth           = 950
	winHeight          = 540
	vertexShaderSource = `
#version 330
layout (location = 0) in vec3 Position;
layout(location = 1) in vec3 vertexColor;
layout(location = 2) in vec3 normal;



uniform vec4 scaleMove;
uniform vec3 LightSource;
uniform vec3 ChangeColour;

out vec3 fragmentColor;
out vec3 outnormal;

struct Lights
{
	vec3 position;
	vec3 diffuse;
};

// float lambert(vec3 N , vec3 L)
// {
// 	vec3 nN = normalize(N);
// 	vec3 nL = normalize(L);
// 	float result = dot(nN, nL);
// 	return max(result, 0.0);
// }

void main()
{ 
// YOU CAN OPTIMISE OUT cos(scaleMove.x) AND sin(scaleMove.y) AND UNIFORM THE VALUES IN
    vec3 scale = Position.xyz * scaleMove.w;
	

	

// rotate on z pole
   vec3 rotatez = vec3((scale.x * cos(scaleMove.x) - scale.y * sin(scaleMove.x)), (scale.x * sin(scaleMove.x) + scale.y * cos(scaleMove.x)), scale.z);

  
// rotate on y pole
    vec3 rotatey = vec3((rotatez.x * cos(scaleMove.y) - rotatez.z * sin(scaleMove.y)), rotatez.y, (rotatez.x * sin(scaleMove.y) + rotatez.z * cos(scaleMove.y)));

	
// rotate on x pole
    vec3 rotatex = vec3(rotatey.x, (rotatey.y * cos(scaleMove.z) - rotatey.z * sin(scaleMove.z)), (rotatey.y * sin(scaleMove.z) + rotatey.z * cos(scaleMove.z)));

	
// move
vec3 move = vec3(rotatex.xy, rotatex.z - 0.2);


// terrible perspective transform
vec3 persp = vec3( move.x  / ( (move.z + 2) / 3 ) ,
		   move.y  / ( (move.z + 2) / 3 ) ,
		     move.z);

    gl_Position = vec4(persp, 1.0);

	Lights lights;
	lights.diffuse = vec3(1.0,1.0,1.0);
	//lights.position = vec3(10.0,0,10.0);

    //fragmentColor = vertexColor;
	//fragmentColor = ChangeColour*lights.diffuse* lambert(normal,LightSource);;

	outnormal = normal;
	fragmentColor = ChangeColour*lights.diffuse;
}
`
	fragmentShaderSource = `
#version 330
out vec4 outColor;
in vec3 fragmentColor;

//for lambert in fragmentShaderSource
in vec3 outnormal;
uniform vec3 LightSource;

float lambert(vec3 N , vec3 L)
{
	vec3 nN = normalize(N);
	vec3 nL = normalize(L);
	float result = dot(nN, nL);
	return max(result, 0.0);
}

void main()
{
	//outColor = vec4(fragmentColor, 1.0);
	outColor = vec4(fragmentColor*lambert(outnormal,LightSource), 1.0);
}
`
)

var triangle_vertices = []gl.Float{}

var triangle_colours = []gl.Float{}

var triangle_normals = []gl.Float{}
