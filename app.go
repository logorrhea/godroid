package main

import (
	"encoding/binary"
	"log"

	"golang.org/x/mobile/app"
	"golang.org/x/mobile/app/debug"
	"golang.org/x/mobile/event"
	"golang.org/x/mobile/f32"
	"golang.org/x/mobile/geom"
	"golang.org/x/mobile/gl"
	"golang.org/x/mobile/gl/glutil"
)

var (
	program  gl.Program
	position gl.Attrib
	offset   gl.Uniform
	color    gl.Uniform
	buf      gl.Buffer

	green    float32
	touchLoc geom.Point
)

func main() {
	app.Run(app.Callbacks{
		Start: start,
		Draw:  draw,
		Stop:  stop,
		Touch: touch,
	})
}

func start() {
	var err error
	_, err = glutil.CreateProgram(vShader, fShader)

	if err != nil {
		log.Printf("Error creating GL program: %v", err)
		return
	}

	buf = gl.GenBuffer()
	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.BufferData(gl.ARRAY_BUFFER, gl.STATIC_DRAW, triangleData)

	position = gl.GetAttribLocation(program, "position")
	color = gl.GetUniformLocation(program, "color")
	offset = gl.GetUniformLocation(program, "offset")
	touchLoc = geom.Point{geom.Width / 2, geom.Height / 2}
}

func stop() {
	gl.DeleteProgram(program)
	gl.DeleteBuffer(buf)
}

func touch(t event.Touch) {
	touchLoc = t.Loc
}

func draw() {
	gl.ClearColor(1, green, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	gl.UseProgram(program)

	green += 0.01
	if green > 1 {
		green = 0
	}
	gl.Uniform4f(color, 0, green, 0, 1)
	gl.Uniform2f(offset, float32(touchLoc.X/geom.Width), float32(touchLoc.Y/geom.Height))

	gl.BindBuffer(gl.ARRAY_BUFFER, buf)
	gl.EnableVertexAttribArray(position)
	gl.VertexAttribPointer(position, coordsPerVertex, gl.FLOAT, false, 0, 0)
	gl.DrawArrays(gl.TRIANGLES, 0, vertexCount)
	gl.DisableVertexAttribArray(position)

	debug.DrawFPS()
}

var triangleData = f32.Bytes(binary.LittleEndian,
	0.0, 0.4, 0.0, // top left
	0.0, 0.0, 0.0, // bottom left
	0.4, 0.0, 0.0, // bottom right
)

const (
	coordsPerVertex = 3
	vertexCount     = 3
)

const vShader = `#version 100
uniform vec2 offset;

attribute vec4 position;
void main() {
    // offset comes in with x/y values between 0 and 1
    // position bounds are -1 to 1
    vec4 offset4 = vec4(2.0*offset.x-1.0, 1.0-2.0*offset.y, 0, 0);
    gl_Position = position + offset4;
}`

const fShader = `#version 100
precision mediump float;
uniform vec4 color;
void main() {
    gl_FragColor = color;
}`
