// Package SVG takes the Contributor struct and creates an SVG from it

package svg

import (
	"bytes"
	"fmt"
	"math"

	"ccollage/internal/client/github"

	svgo "ccollage/third_party/svgo"
)

// BuildCollage takes a contributor array and returns a bytes buffer for use in the HTTP server
func BuildCollage(contributors []github.Contributor, width int, padding int, maxWidth int) bytes.Buffer {
	var avatarWidth, avatarHeight int = width, width
	var avatarHorzPadding = padding
	var avatarVertPadding = padding
	var canvasMaxWidth = maxWidth
	var canvasPadding = padding
	var canvasRadius = avatarWidth / 2

	var canvasCols = canvasMaxWidth / (avatarWidth + avatarHorzPadding)
	var canvasRows = int(math.Ceil(float64(len(contributors)) / float64(canvasCols)))

	var buf bytes.Buffer

	canvas := svgo.New(&buf)

	canvas.Start((avatarWidth+avatarHorzPadding*2)*canvasCols+canvasPadding*2, (avatarHeight+avatarVertPadding*2)*canvasRows+canvasPadding*2)

	for i := range contributors {
		var imgRowPos = i / canvasCols
		var imgColPos = i % canvasCols

		var imgPosX = canvasPadding + imgColPos*(avatarWidth+avatarHorzPadding*2)
		var imgPosY = canvasPadding + imgRowPos*(avatarHeight+avatarVertPadding*2)

		cAvatar := fmt.Sprintf("%s&amp;s=%d", contributors[i].Avatar, avatarWidth)

		canvas.ClipPath(fmt.Sprintf(`id="circle%d"`, i))
		canvas.Circle(canvasPadding+imgPosX+canvasRadius, canvasPadding+imgPosY+canvasRadius, canvasRadius)
		canvas.ClipEnd()

		canvas.Circle(canvasPadding+imgPosX+canvasRadius, canvasPadding+imgPosY+canvasRadius, canvasRadius, "fill:none;stroke:black")
		canvas.Group(fmt.Sprintf(`clip-path="url(#circle%d)"`, i))
		canvas.Link(contributors[i].URL, contributors[i].Username)
		canvas.Image(imgPosX, imgPosY, canvasPadding+avatarWidth, canvasPadding+avatarHeight, cAvatar)
		canvas.LinkEnd()
		canvas.Gend()
	}

	canvas.End()

	return buf
}
