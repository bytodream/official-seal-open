// 本文件是本章的独立副本：绘制实现原属公共 core 包，现已按章各持一份。
// 目的是让章模块不依赖公共层即可单独交付；公共实现变更时需同步本文件。
package enterprise

import (
	"image"
	"image/color"
	stddraw "image/draw"
	"math"

	"github.com/bytodream/official-seal-open/seal/internal/draw/base"
	"github.com/bytodream/official-seal-open/seal/model"
	"github.com/bytodream/official-seal-open/seal/utils"
	xdraw "golang.org/x/image/draw"
	"golang.org/x/image/font"
	"golang.org/x/image/math/f64"
	"golang.org/x/image/math/fixed"
)

// generateSealFiveStar 绘制方形防伪纹底图（竖线、字掩膜、横线、斜线、空心五角星）；字号约为画布高度 25%，最低 8px。
func generateSealFiveStar(info *model.AntifakeInfo, antiWidth, antiHeight int) (*image.RGBA, error) {
	if info == nil {
		return nil, model.ErrSealAttributeNotConfigured
	}
	fontSize := int(float64(antiHeight) * 0.25)
	if fontSize < 8 {
		fontSize = 8
	}

	img := image.NewRGBA(image.Rect(0, 0, antiWidth, antiHeight))
	fillRect(img, img.Bounds(), color.White)

	spacing := info.Space()
	lineSize := float64(info.LineSize)

	for x := lineSize / 2; x < float64(antiWidth); x += spacing + lineSize {
		drawStrokeLineRGBA(img, int(x), 0, int(x), antiHeight, int(lineSize), info.Color)
	}

	face, err := base.LoadFontFace(info.FontFile, float64(fontSize), info.Dpi)
	if err != nil {
		return nil, err
	}
	defer face.Close()

	metrics := face.Metrics()
	d := &font.Drawer{Face: face}
	textW := int(math.Round(float64(d.MeasureString(info.Content).Ceil()) * antifakeTextScaleX))
	x := float64(antiWidth-textW) / 2
	y := float64(antiHeight)/2 + float64(metrics.Ascent.Ceil())/4.0 - 10

	drawTextOn(img, face, info.Content, x, y, color.White)

	content := image.NewRGBA(image.Rect(0, 0, antiWidth, antiHeight))
	fillRect(content, content.Bounds(), color.White)
	drawTextOn(content, face, info.Content, x-(spacing+lineSize)/2, y, info.Color)
	fillGratingMaskColumns(content, spacing, lineSize, antiWidth, antiHeight)
	content = imageOpaqueNonWhite(content)
	stddraw.Draw(img, img.Bounds(), content, image.Point{}, stddraw.Over)

	for y0 := 0.0; y0 < float64(antiHeight); y0 += spacing + lineSize {
		drawStrokeLineRGBA(img, 0, int(y0), antiWidth, int(y0), 1, info.Color)
	}

	layer := image.NewRGBA(image.Rect(0, 0, antiWidth, antiHeight))
	drawDiagonalRightLines(layer, spacing, 1, 30, info.Color)
	stddraw.Draw(img, img.Bounds(), layer, image.Point{}, stddraw.Over)

	layer = image.NewRGBA(image.Rect(0, 0, antiWidth, antiHeight))
	drawDiagonalLeftLines(layer, spacing, 1, 30, info.Color)
	stddraw.Draw(img, img.Bounds(), layer, image.Point{}, stddraw.Over)

	layer = image.NewRGBA(image.Rect(0, 0, antiWidth, antiHeight))
	drawDiagonalRightLines(layer, spacing, 1, 45, info.Color)
	stddraw.Draw(img, img.Bounds(), layer, image.Point{}, stddraw.Over)

	drawAntifakeEmptyStar(img)
	return img, nil
}

// antifakeTextScaleX 掩膜文字横向拉伸比。
const antifakeTextScaleX = 1.2

// drawTextOn 在 (x,y) 绘制加粗并横向拉伸的防伪字；先在小块 tile 上合成加粗再贴回，避免膨胀污染底色。
func drawTextOn(dst *image.RGBA, face font.Face, s string, x, y float64, col color.Color) {
	pad := simulatedBoldRadius(face)
	metrics := face.Metrics()
	ascent, descent := metrics.Ascent.Ceil(), metrics.Descent.Ceil()
	w := (&font.Drawer{Face: face}).MeasureString(s).Ceil() + 2*pad
	h := ascent + descent + 2*pad
	if w < 1 || h < 1 {
		return
	}
	tile := image.NewRGBA(image.Rect(0, 0, w, h))
	td := &font.Drawer{
		Dst:  tile,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.P(pad, ascent+pad),
	}
	drawBoldString(tile, td, s)

	nw := int(math.Round(float64(w) * antifakeTextScaleX))
	if nw < 1 {
		nw = 1
	}
	stretched := image.NewRGBA(image.Rect(0, 0, nw, h))
	xdraw.NearestNeighbor.Scale(stretched, stretched.Bounds(), tile, tile.Bounds(), xdraw.Over, nil)

	left := int(math.Round(x)) - int(math.Round(float64(pad)*antifakeTextScaleX))
	top := int(math.Round(y)) - ascent - pad
	dstRect := image.Rect(left, top, left+nw, top+h)
	stddraw.Draw(dst, dstRect, stretched, image.Point{}, stddraw.Over)
}

// fillRect 用纯色填充矩形区域。
func fillRect(img *image.RGBA, r image.Rectangle, c color.Color) {
	stddraw.Draw(img, r, &image.Uniform{c}, image.Point{}, stddraw.Src)
}

// fillGratingMaskColumns 按光栅间距把竖条之间的间隙填白，形成字掩膜后与竖线叠加。
func fillGratingMaskColumns(img *image.RGBA, spacing, lineSize float64, w, h int) {
	rects := []image.Rectangle{
		image.Rect(0, 0, int((spacing+lineSize)/2-1), h),
	}
	for start := (spacing + 3*lineSize) / 2; start < float64(w); start += spacing + lineSize {
		rects = append(rects, image.Rect(int(start), 0, int(start+spacing-1), h))
	}
	for _, r := range rects {
		fillRect(img, r, color.White)
	}
}

// imageOpaqueNonWhite 保留非白不透明像素并强制 alpha=255。
func imageOpaqueNonWhite(src *image.RGBA) *image.RGBA {
	b := src.Bounds()
	out := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.RGBAModel.Convert(src.At(x, y)).(color.RGBA)
			if c.A == 0 {
				continue
			}
			if c.R > 225 && c.G > 225 && c.B > 225 {
				continue
			}
			out.SetRGBA(x, y, color.RGBA{c.R, c.G, c.B, 255})
		}
	}
	return out
}

// drawStrokeLineRGBA 绘制线宽为 w 像素的线段。
func drawStrokeLineRGBA(img *image.RGBA, x0, y0, x1, y1, w int, c color.Color) {
	if w < 1 {
		w = 1
	}
	dx := float64(x1 - x0)
	dy := float64(y1 - y0)
	length := math.Hypot(dx, dy)
	steps := int(length)
	if steps < 1 {
		steps = 1
	}
	lo := -(w - 1) / 2
	hi := w - 1 + lo
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x0 + int(math.Round(t*dx))
		y := y0 + int(math.Round(t*dy))
		for sy := lo; sy <= hi; sy++ {
			for sx := lo; sx <= hi; sx++ {
				if InnerBounds(img.Bounds(), x+sx, y+sy) {
					img.Set(x+sx, y+sy, c)
				}
			}
		}
	}
}

// drawDiagonalRightLines 绘制绕左上向右下的平行斜线族，rotateDeg 为线与水平轴夹角。
func drawDiagonalRightLines(img *image.RGBA, spacing, lineSize float64, rotateDeg int, col color.Color) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	rad := float64(rotateDeg) * math.Pi / 180
	sin, cos := math.Sincos(rad)
	if sin == 0 {
		return
	}
	for curLen := spacing; ; curLen += spacing {
		curWidth := curLen / sin
		curHeight := curLen / cos
		if curWidth >= float64(w*2) && curHeight >= float64(h*2) {
			break
		}
		drawStrokeLineRGBA(img, int(curWidth), 0, 0, int(curHeight), int(lineSize), col)
	}
}

// drawDiagonalLeftLines 绘制绕左下向左上的平行斜线族，rotateDeg 为线与水平轴夹角。
func drawDiagonalLeftLines(img *image.RGBA, spacing, lineSize float64, rotateDeg int, col color.Color) {
	w, h := img.Bounds().Dx(), img.Bounds().Dy()
	rad := float64(rotateDeg) * math.Pi / 180
	sin, cos := math.Sincos(rad)
	if sin == 0 {
		return
	}
	for curLen := spacing; ; curLen += spacing {
		curWidth := curLen / sin
		curHeight := curLen / cos
		if curWidth >= float64(w*2) && curHeight >= float64(h*2) {
			break
		}
		drawStrokeLineRGBA(img, int(curWidth), h, 0, h-int(curHeight), int(lineSize), col)
	}
}

// drawAntifakeEmptyStar 在画布中心按边长绘制空心五角星白色遮罩。
func drawAntifakeEmptyStar(buff *image.RGBA) {
	w := buff.Bounds().Dx()
	h := buff.Bounds().Dy()
	drawEmptyStarMask(buff, float64(w)/2, float64(h)/2, w)
}

// drawEmptyStarMask 用五个多边形拼出空心五角星，(d,e) 为星心，length 为外接正方形边长。
func drawEmptyStarMask(g *image.RGBA, d, e float64, length int) {
	r := length / 2
	ch := 1.2566370614359172
	x2 := d - math.Sin(ch)*float64(r)
	x3 := d + math.Sin(ch)*float64(r)
	x4 := d - math.Sin(ch/2)*float64(r)
	x5 := d + math.Sin(ch/2)*float64(r)
	y1 := e - float64(r)
	y2 := e - math.Cos(ch)*float64(r)
	y4 := e + math.Cos(ch/2)*float64(r)
	cx2 := d - math.Tan(ch/2)*math.Cos(ch)*float64(r)
	cx1 := d + math.Tan(ch/2)*math.Cos(ch)*float64(r)
	cx3 := d - math.Sin(ch)*(math.Cos(ch)*float64(r)/math.Cos(ch/2))
	cx4 := d + math.Sin(ch)*(math.Cos(ch)*float64(r)/math.Cos(ch/2))
	cy1 := e - math.Cos(ch)*float64(r)
	cy3 := e + math.Cos(ch)*(math.Cos(ch)*float64(r)/math.Cos(ch/2))
	cy5 := e + math.Cos(ch)*float64(r)/math.Cos(ch/2)

	polys := [][]image.Point{
		{{0, 0}, {int(d), int(y1)}, {int(cx2), int(cy1)}, {0, int(y2)}},
		{{int(d), int(y1)}, {int(cx1), int(cy1)}, {int(2 * d), int(y2)}, {int(2 * d), 0}},
		{{int(2 * d), 0}, {int(x3), int(y2)}, {int(cx4), int(cy3)}, {int(x5), int(y4)}, {int(2 * d), int(2 * e)}},
		{{int(2 * d), int(2 * e)}, {int(x5), int(y4)}, {int(d), int(cy5)}, {int(x4), int(y4)}, {0, int(2 * e)}},
		{{0, int(2 * e)}, {int(x4), int(y4)}, {int(cx3), int(cy3)}, {int(x2), int(y2)}, {0, int(y2)}},
	}
	for _, pts := range polys {
		FillPolygon(g, pts, color.White)
	}
}

// DrawAntifakeLogo 在圆章中心绘制防伪纹或实心五角星。无 model.AntifakeInfo 时退化为 DrawFiveStar。
// 贴图矩形严格为 logoSize×logoSize、左上角 (cx-logoSize/2, cy-logoSize/2)
// 尺寸匹配时 1:1 拷贝，避免二次重采样糊掉细线纹理。
func DrawAntifakeLogo(dst *image.RGBA, m base.SealPaintModel) error {
	info := m.AntifakeInfo
	logoSize := m.LogoSize
	if logoSize < 1 {
		return nil
	}
	if info == nil {
		DrawFiveStar(dst, m.Width/2, m.Height/2-int(m.LogoUpmove), logoSize, model.SealRed)
		return nil
	}
	orig := info.Dpi * logoSize / m.Dpi
	if orig < 16 {
		orig = 16
	}
	patch, err := generateSealFiveStar(info, orig, orig)
	if err != nil {
		return err
	}
	scale := float64(m.Dpi) / float64(info.Dpi)
	scaled := scaleImage(patch, scale)
	patchAlpha := Translucent(scaled, m.Transparency)

	cx := m.Width / 2
	cy := m.Height/2 - int(m.LogoUpmove)
	left := cx - logoSize/2
	top := cy - logoSize/2
	dstRect := image.Rect(left, top, left+logoSize, top+logoSize)
	fit := image.NewRGBA(image.Rect(0, 0, dstRect.Dx(), dstRect.Dy()))
	if patchAlpha.Bounds().Dx() == fit.Bounds().Dx() && patchAlpha.Bounds().Dy() == fit.Bounds().Dy() {
		stddraw.Draw(fit, fit.Bounds(), patchAlpha, patchAlpha.Bounds().Min, stddraw.Src)
	} else {
		xdraw.BiLinear.Scale(fit, fit.Bounds(), patchAlpha, patchAlpha.Bounds(), xdraw.Over, nil)
	}
	xdraw.Draw(dst, dstRect, fit, image.Point{}, xdraw.Over)
	return nil
}

// scaleImage 按 scale 双三次缩放 RGBA。
// 白底红纹的抗锯齿体现在 RGB 变浅而非 alpha 降低；用 Aff3 Transform 对齐采样中心。
func scaleImage(src *image.RGBA, scale float64) *image.RGBA {
	if scale <= 0 || math.Abs(scale-1) < 0.001 {
		return src
	}
	b := src.Bounds()
	nw := int(float64(b.Dx())*scale + 0.5)
	nh := int(float64(b.Dy())*scale + 0.5)
	if nw < 1 {
		nw = 1
	}
	if nh < 1 {
		nh = 1
	}
	out := image.NewRGBA(image.Rect(0, 0, nw, nh))
	sx := float64(nw) / float64(b.Dx())
	sy := float64(nh) / float64(b.Dy())
	aff := f64.Aff3{sx, 0, 0.5*sx - 0.5, 0, sy, 0.5*sy - 0.5}
	xdraw.CatmullRom.Transform(out, aff, src, b, xdraw.Over, nil)
	return out
}

// FootValid 判断脚注属性是否含有效文字。
func FootValid(f *base.FootAttr) bool {
	return f != nil && !utils.IsBlank(f.Text)
}

// DrawFootCenter 在圆章底部居中绘制脚注横排文字。
// 既不做网格拟合，也要把 zoom 作用到字形本身（而不是只改字距）。
func DrawFootCenter(img *image.RGBA, m base.SealPaintModel, col color.Color) error {
	f := m.Foot
	runes := []rune(f.Text)
	if len(runes) == 0 {
		return nil
	}
	face, err := base.LoadFontFaceHinting(m.FontFile, float64(f.FontSizePx)*ArcTextSS, m.Dpi, font.HintingNone)
	if err != nil {
		return err
	}
	defer face.Close()

	cx := float64(m.Width) / 2
	baseline := float64(m.Height)/2 + float64(f.DistancePx) + float64(f.FontSizePx)

	d := &font.Drawer{Face: face}
	fullW := float64(d.MeasureString(f.Text).Ceil()) / ArcTextSS
	charW := fullW / float64(len(runes))
	// WidthPx<=0（未指定排宽，如 DIY 默认值）时按文字自然宽度居中，避免缩放为 0 不可见。
	width := float64(f.WidthPx)
	if width <= 0 {
		width = fullW
	}
	avgW := width / float64(len(runes))
	zoom := 0.75
	if avgW < charW {
		zoom = avgW / charW
		if zoom > 0.75 {
			zoom = 0.75
		}
	}
	charW *= zoom

	for i := len(runes) - 1; i >= 0; i-- {
		ch := string(runes[i])
		x := width/2 - (avgW/2 + float64(len(runes)-1-i)*avgW) - charW/2
		drawRotatedStringScaled(img, face, ch, cx+float64(int(x)), baseline, 0, col, zoom, ArcTextSS)
	}
	return nil
}

// DrawAnnulus 圆环描边。
func DrawAnnulus(img *image.RGBA, cx, cy int, midRadius, edge float64, c color.Color) {
	rOut := midRadius + edge/2
	rIn := midRadius - edge/2
	if rIn < 0 {
		rIn = 0
	}
	bound := int(rOut) + 2
	for dy := -bound; dy <= bound; dy++ {
		for dx := -bound; dx <= bound; dx++ {
			d := math.Hypot(float64(dx), float64(dy))
			if d >= rIn && d <= rOut {
				x, y := cx+dx, cy+dy
				if InnerBounds(img.Bounds(), x, y) {
					img.Set(x, y, c)
				}
			}
		}
	}
}

// FillDisk 填充实心圆盘。
func FillDisk(img *image.RGBA, cx, cy, r int, c color.Color) {
	for dy := -r; dy <= r; dy++ {
		for dx := -r; dx <= r; dx++ {
			if dx*dx+dy*dy <= r*r {
				x, y := cx+dx, cy+dy
				if InnerBounds(img.Bounds(), x, y) {
					img.Set(x, y, c)
				}
			}
		}
	}
}

// InnerBounds 判断坐标是否在图像边界内。
func InnerBounds(b image.Rectangle, x, y int) bool {
	return x >= b.Min.X && x < b.Max.X && y >= b.Min.Y && y < b.Max.Y
}

// DrawAntifakeGapLines 在圆章上绘制防伪断裂短线。
func DrawAntifakeGapLines(img *image.RGBA, width int, edge float32) {
	rng := base.NewRand()
	lineCounts := []int{5, 6, 7, 8, 9, 10}
	lineWidths := []int{2, 3, 4, 5, 6, 7, 8, 9, 10}
	rotates := []int{-5, -4, -3, -2, -1, 1, 2, 3, 4, 5}
	n := lineCounts[rng.Intn(len(lineCounts))]
	half := float64(width) / 2
	maxR := half - float64(edge)
	for i := 0; i < n; i++ {
		lw := lineWidths[rng.Intn(len(lineWidths))]
		theta := rng.Float64() * 2 * math.Pi
		x0 := half * math.Cos(theta)
		y0 := half * math.Sin(theta)
		rot := float64(rotates[rng.Intn(len(rotates))]) * math.Pi / 180
		theta1 := theta - rot
		x1 := maxR * math.Cos(theta1)
		y1 := maxR * math.Sin(theta1)
		DrawThickLine(img, int(half+x0), int(half+y0), int(half+x1), int(half+y1), lw, color.White)
	}
}

// DrawThickLine 用圆盘取样绘制粗线段。
func DrawThickLine(img *image.RGBA, x0, y0, x1, y1, w int, c color.Color) {
	dx := float64(x1 - x0)
	dy := float64(y1 - y0)
	length := math.Hypot(dx, dy)
	if length < 1 {
		return
	}
	steps := int(length)
	for i := 0; i <= steps; i++ {
		t := float64(i) / float64(steps)
		x := x0 + int(t*dx)
		y := y0 + int(t*dy)
		FillDisk(img, x, y, w/2+1, c)
	}
}

// DrawFiveStar 在 (x,y) 绘制实心五角星，length 为外接边长。
func DrawFiveStar(img *image.RGBA, x, y, length int, c color.Color) {
	r := length / 2
	ch := 1.2566370614359172
	x2 := int(float64(x) - math.Sin(ch)*float64(r))
	x3 := int(float64(x) + math.Sin(ch)*float64(r))
	x4 := int(float64(x) - math.Sin(ch/2)*float64(r))
	x5 := int(float64(x) + math.Sin(ch/2)*float64(r))
	y1 := y - r
	y2 := int(float64(y) - math.Cos(ch)*float64(r))
	y4 := int(float64(y) + math.Cos(ch/2)*float64(r))
	bx := int(float64(x) + math.Cos(ch)*math.Tan(ch/2)*float64(r))
	FillTriangle(img, []image.Point{{x2, y2}, {x5, y4}, {bx, y2}}, c)
	FillTriangle(img, []image.Point{{x, y1}, {bx, y2}, {x3, y2}, {x4, y4}}, c)
}

// FillPolygon 用射线法逐像素填充多边形（支持凹多边形）。
func FillPolygon(img *image.RGBA, pts []image.Point, c color.Color) {
	if len(pts) < 3 {
		return
	}
	minX, maxX := pts[0].X, pts[0].X
	minY, maxY := pts[0].Y, pts[0].Y
	for _, p := range pts[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if pointInPolygon(pts, x, y) {
				if InnerBounds(img.Bounds(), x, y) {
					img.Set(x, y, c)
				}
			}
		}
	}
}

// FillTriangle 填充三角形区域。
func FillTriangle(img *image.RGBA, pts []image.Point, c color.Color) {
	minX, maxX := pts[0].X, pts[0].X
	minY, maxY := pts[0].Y, pts[0].Y
	for _, p := range pts[1:] {
		if p.X < minX {
			minX = p.X
		}
		if p.X > maxX {
			maxX = p.X
		}
		if p.Y < minY {
			minY = p.Y
		}
		if p.Y > maxY {
			maxY = p.Y
		}
	}
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			if pointInPolygon(pts, x, y) {
				if InnerBounds(img.Bounds(), x, y) {
					img.Set(x, y, c)
				}
			}
		}
	}
}

// pointInPolygon 射线法判断点是否在多边形内。
func pointInPolygon(pts []image.Point, x, y int) bool {
	inside := false
	j := len(pts) - 1
	for i := 0; i < len(pts); i++ {
		xi, yi := pts[i].X, pts[i].Y
		xj, yj := pts[j].X, pts[j].Y
		if ((yi > y) != (yj > y)) && (float64(x) < float64(xj-xi)*float64(y-yi)/(float64(yj-yi)+1e-9)+float64(xi)) {
			inside = !inside
		}
		j = i
	}
	return inside
}

// DrawNestedInnerCircle 绘制嵌套内圆环与内圈文字。
func DrawNestedInnerCircle(img *image.RGBA, m base.SealPaintModel, col color.Color) error {
	inner := m.InnerCircle
	if inner == nil {
		return nil
	}
	cx, cy := m.Width/2, m.Height/2
	diam := float64(inner.WidthPx) - float64(inner.EdgeWidthPx)
	if diam < 1 {
		return nil
	}
	radius := diam / 2
	DrawAnnulus(img, cx, cy, radius, float64(inner.EdgeWidthPx), col)
	if utils.IsBlank(inner.Content) {
		return nil
	}
	innerModel := base.SealPaintModel{
		Dpi:      m.Dpi,
		FontFile: m.FontFile,
		Width:    m.Width,
		Height:   m.Height,
		Content:  inner.Content,
		Arc:      inner.Arc,
		Margin:   inner.MarginPx,
		FontSize: inner.FontSizePx,
	}
	return DrawArcContent(img, innerModel, radius, col)
}

// ArcTextSS 环排文字超采样倍数；先高分辨率栅格化再缩回，减轻二次插值导致的硬边。
const ArcTextSS = 2

// DrawArcContent 沿上弧绘制主环绕文字。
func DrawArcContent(img *image.RGBA, m base.SealPaintModel, radius float64, col color.Color) error {
	runes := []rune(m.Content)
	if len(runes) == 0 {
		return nil
	}
	face, err := base.LoadFontFaceHinting(m.FontFile, float64(m.FontSize)*ArcTextSS, m.Dpi, font.HintingNone)
	if err != nil {
		return err
	}
	defer face.Close()

	d := &font.Drawer{Face: face}
	fullAdvance := int(math.Round(float64(d.MeasureString(m.Content).Ceil()) / ArcTextSS))
	// 上伸度量取字体逻辑上边界（= ascender + lineGap/2），非字形墨迹上边界。
	ascent := base.StringBoundsAscent(face, ArcTextSS)
	fontRadius := radius - ascent - float64(m.Margin)

	charW := float64(fullAdvance) / float64(len(runes))
	zoom := 0.75
	avgLen := float64(m.Arc*2) * math.Pi * fontRadius / (float64(len(runes)) * 360)
	if avgLen < charW {
		zoom = avgLen / charW
		if zoom > 0.75 {
			zoom = 0.75
		}
	}
	charW *= zoom

	startAngle := 270.0 - float64(m.Arc)/2.0
	stepAngle := float64(m.Arc) / float64(len(runes))
	cx, cy := float64(m.Width)/2, float64(m.Height)/2

	for i, r := range runes {
		ch := string(r)
		angle := startAngle + stepAngle/2 + float64(i)*stepAngle
		ark := angle * math.Pi / 180
		anglex := fontRadius * math.Cos(ark)
		angley := fontRadius * math.Sin(ark)

		var rotateDeg float64
		if angle >= 90 && angle < 270 {
			rotateDeg = 90 + angle
		} else {
			rotateDeg = angle - 270
		}

		x := anglex + charW/2*math.Sin(ark)
		y := angley - charW/2*math.Cos(ark)
		drawRotatedStringScaled(img, face, ch, cx+x, cy+y, rotateDeg*math.Pi/180, col, zoom, ArcTextSS)
	}
	return nil
}

// DrawEnterpriseCode 沿下半弧逆序绘制企业编码文字。
func DrawEnterpriseCode(img *image.RGBA, m base.SealPaintModel, radius float64, col color.Color) error {
	ec := m.EnterpriseCode
	code := ec.Code
	runes := []rune(code)
	if len(runes) == 0 {
		return nil
	}
	face, err := base.LoadFontFaceHinting(m.FontFile, float64(ec.FontSizePx)*ArcTextSS, m.Dpi, font.HintingNone)
	if err != nil {
		return err
	}
	defer face.Close()

	d := &font.Drawer{Face: face}
	fullW := int(math.Round(float64(d.MeasureString(code).Ceil()) / ArcTextSS))
	charW := float64(fullW) / float64(len(runes))
	// 即 margin 同样要参与 /2（原实现是 radius - ascent/2 - margin）；ascent 取字体逻辑上边界。
	ascent := base.StringBoundsAscent(face, ArcTextSS)
	fontRadius := radius - (ascent+float64(m.Margin))/2
	arc := ec.Arc
	if arc <= 0 {
		arc = 90
	}
	avgLen := float64(arc) * math.Pi * fontRadius / (float64(len(runes)) * 180)
	zoom := 1.0
	if avgLen < charW {
		zoom = avgLen / charW
	}
	charW *= zoom

	startAngle := 90.0 - float64(arc)/2
	step := float64(arc) / float64(len(runes))
	cx, cy := float64(m.Width)/2, float64(m.Height)/2

	for i := len(runes) - 1; i >= 0; i-- {
		ch := string(runes[i])
		idx := len(runes) - 1 - i
		angle := startAngle + step/2 + float64(idx)*step
		ark := angle * math.Pi / 180
		x := fontRadius*math.Cos(ark) - charW/2*math.Sin(ark)
		y := fontRadius*math.Sin(ark) + charW/2*math.Cos(ark)
		rot := (angle - 90) * math.Pi / 180
		drawRotatedStringScaled(img, face, ch, cx+x, cy+y, rot, col, zoom, ArcTextSS)
	}
	return nil
}

// drawRotatedStringScaled 的 ss 表示 face 已按 ss 倍字号加载，坐标换算回 1x 空间。
func drawRotatedStringScaled(img *image.RGBA, face font.Face, s string, x, y, rot float64, col color.Color, zoom float64, ss int) {
	if ss < 1 {
		ss = 1
	}
	d := &font.Drawer{Face: face}
	w := d.MeasureString(s).Ceil()
	metrics := face.Metrics()
	baseline := metrics.Ascent.Ceil()
	h := baseline + metrics.Descent.Ceil()
	if w < 1 {
		w = 1
	}
	if h < 1 {
		h = 1
	}
	// 合成加粗会向四周扩张，tile 需要预留 padding，否则右/下笔画被裁掉。
	pad := simulatedBoldRadius(face)
	tile := image.NewRGBA(image.Rect(0, 0, w+2*pad, h+2*pad))
	td := &font.Drawer{
		Dst:  tile,
		Src:  image.NewUniform(col),
		Face: face,
		Dot:  fixed.P(pad, baseline+pad),
	}
	drawBoldString(tile, td, s)
	tile = scaleTileHorizontal(tile, zoom)

	bounds := tile.Bounds()
	w = bounds.Dx()
	pivot := image.Pt(int(math.Round(float64(pad)*zoom)), baseline+pad)
	// 旋转与超采样降采样合并成一次重采样：分成两步会让笔画先被双线性插值再被降采样
	// 既慢又比单次变换更糊。
	rotated, pivotX, pivotY := rotateAroundPivot(tile, rot, pivot, 1/float64(ss))
	dstPt := image.Pt(int(math.Round(x-pivotX)), int(math.Round(y-pivotY)))
	xdraw.Draw(img, rotated.Bounds().Add(dstPt), rotated, image.Point{}, xdraw.Over)
}

// drawBoldString 绘制加粗文字：先光栅化再对 alpha 做圆盘膨胀（可分离最大值滤波，避免重复绘制）。
func drawBoldString(tile *image.RGBA, td *font.Drawer, s string) {
	base := td.Dot
	td.DrawString(s)
	td.Dot = base
	r := simulatedBoldRadius(td.Face)
	if r <= 0 {
		return
	}
	fill := color.Color(color.Black)
	if u, ok := td.Src.(*image.Uniform); ok {
		fill = u.C
	}
	dilateAlpha(tile, r, fill)
}

// simulatedBoldRadius 返回模拟加粗所需的圆盘膨胀半径（像素）。
func simulatedBoldRadius(face font.Face) int {
	m := face.Metrics()
	em := float64(m.Ascent.Ceil() + m.Descent.Ceil())
	if em <= 0 {
		return 0
	}
	r := int(math.Round(em * 0.008))
	if r < 1 {
		r = 1
	}
	return r
}

// scaleTileHorizontal 水平缩放字形位图 tile。
func scaleTileHorizontal(tile *image.RGBA, zoom float64) *image.RGBA {
	if zoom >= 0.999 || tile.Bounds().Dx() < 2 {
		return tile
	}
	nw := int(math.Round(float64(tile.Bounds().Dx()) * zoom))
	if nw < 1 {
		nw = 1
	}
	out := image.NewRGBA(image.Rect(0, 0, nw, tile.Bounds().Dy()))
	xdraw.ApproxBiLinear.Scale(out, out.Bounds(), tile, tile.Bounds(), xdraw.Over, nil)
	return out
}

// rotateAroundPivot 替代已移除的 golang.org/x/image/transform.Rotate（x/image ≥ v0.35 无 transform 子包）。
// 返回旋转后图像以及 pivot 在目标图上的坐标。
// scale 用于把超采样图一次性缩回 1x（<1 时用 CatmullRom，避免二次重采样）。
func rotateAroundPivot(src image.Image, theta float64, pivot image.Point, scale float64) (*image.RGBA, float64, float64) {
	if scale <= 0 {
		scale = 1
	}
	b := src.Bounds()
	sin, cos := math.Sincos(theta)
	px, py := float64(pivot.X), float64(pivot.Y)
	t := f64.Aff3{
		cos * scale, -sin * scale, (px - cos*px + sin*py) * scale,
		sin * scale, cos * scale, (py - sin*px - cos*py) * scale,
	}

	corners := []image.Point{
		b.Min,
		image.Pt(b.Max.X-1, b.Min.Y),
		image.Pt(b.Max.X-1, b.Max.Y-1),
		image.Pt(b.Min.X, b.Max.Y-1),
	}
	minX, minY := math.MaxFloat64, math.MaxFloat64
	maxX, maxY := -math.MaxFloat64, -math.MaxFloat64
	for _, p := range corners {
		dx := t[0]*float64(p.X) + t[1]*float64(p.Y) + t[2]
		dy := t[3]*float64(p.X) + t[4]*float64(p.Y) + t[5]
		minX = math.Min(minX, dx)
		minY = math.Min(minY, dy)
		maxX = math.Max(maxX, dx)
		maxY = math.Max(maxY, dy)
	}

	t[2] -= minX
	t[5] -= minY
	dstW := int(math.Ceil(maxX - minX))
	dstH := int(math.Ceil(maxY - minY))
	if dstW < 1 {
		dstW = 1
	}
	if dstH < 1 {
		dstH = 1
	}
	pivotX := t[0]*px + t[1]*py + t[2]
	pivotY := t[3]*px + t[4]*py + t[5]

	dst := image.NewRGBA(image.Rect(0, 0, dstW, dstH))
	if scale < 0.999 {
		xdraw.CatmullRom.Transform(dst, t, src, b, xdraw.Over, nil)
	} else {
		xdraw.BiLinear.Transform(dst, t, src, b, xdraw.Over, nil)
	}
	return dst, pivotX, pivotY
}

// dilateAlpha 以 (2r+1)² 的窗口对 alpha 做可分离最大值滤波，并按 fill 还原预乘颜色。
func dilateAlpha(g *image.RGBA, r int, fill color.Color) {
	if r < 1 {
		return
	}
	b := g.Bounds()
	w, h := b.Dx(), b.Dy()
	if w < 1 || h < 1 {
		return
	}
	alpha := make([]uint8, w*h)
	// 膨胀只影响墨迹外扩 r 像素的范围，先求墨迹包围盒，避免对整个 tile 做滤波。
	x0, y0, x1, y1 := w, h, -1, -1
	for y := 0; y < h; y++ {
		row := g.Pix[y*g.Stride:]
		base := y * w
		for x := 0; x < w; x++ {
			a := row[x*4+3]
			alpha[base+x] = a
			if a > 0 {
				if x < x0 {
					x0 = x
				}
				if x > x1 {
					x1 = x
				}
				if y < y0 {
					y0 = y
				}
				if y > y1 {
					y1 = y
				}
			}
		}
	}
	if x1 < 0 {
		return
	}
	loX, hiX := x0-r, x1+r
	loY, hiY := y0-r, y1+r
	if loX < 0 {
		loX = 0
	}
	if loY < 0 {
		loY = 0
	}
	if hiX > w-1 {
		hiX = w - 1
	}
	if hiY > h-1 {
		hiY = h - 1
	}
	tmp := make([]uint8, w*h)
	for y := y0; y <= y1; y++ {
		base := y * w
		for x := loX; x <= hiX; x++ {
			lo, hi := x-r, x+r
			if lo < 0 {
				lo = 0
			}
			if hi > w-1 {
				hi = w - 1
			}
			var m uint8
			for k := lo; k <= hi; k++ {
				if v := alpha[base+k]; v > m {
					m = v
				}
			}
			tmp[base+x] = m
		}
	}
	cr, cg, cb, _ := fill.RGBA()
	fr, fg, fb := uint32(cr>>8), uint32(cg>>8), uint32(cb>>8)
	for y := loY; y <= hiY; y++ {
		lo, hi := y-r, y+r
		if lo < 0 {
			lo = 0
		}
		if hi > h-1 {
			hi = h - 1
		}
		row := g.Pix[y*g.Stride:]
		for x := loX; x <= hiX; x++ {
			var m uint8
			for k := lo; k <= hi; k++ {
				if v := tmp[k*w+x]; v > m {
					m = v
				}
			}
			if m == 0 {
				continue
			}
			a := uint32(m)
			row[x*4+0] = uint8(fr * a / 255)
			row[x*4+1] = uint8(fg * a / 255)
			row[x*4+2] = uint8(fb * a / 255)
			row[x*4+3] = uint8(a)
		}
	}
}

// Translucent 将章面色带写入指定 Alpha，白色背景像素丢弃。
// 避免把预乘样本当成非预乘写回导致 PNG 编码反预乘时溢出（红色 255 会退化成 15）。
func Translucent(src image.Image, transparency int) *image.RGBA {
	b := src.Bounds()
	out := image.NewRGBA(b)
	// 钳制到 uint8 范围，避免越界值回绕（如 -127 → 129）。
	if transparency < 0 {
		transparency = 0
	}
	if transparency > 255 {
		transparency = 255
	}
	alpha := uint8(transparency)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := color.NRGBAModel.Convert(src.At(x, y)).(color.NRGBA)
			if c.A == 0 {
				continue
			}
			if c.R > 225 && c.G > 225 && c.B > 225 {
				continue
			}
			out.Set(x, y, color.NRGBA{R: c.R, G: c.G, B: c.B, A: alpha})
		}
	}
	return out
}

// DefaultSampleWatermarkText 预览样章标识默认文字。
const DefaultSampleWatermarkText = "样章"

// sampleWatermarkColor 预览样章标识颜色（半透明灰，叠加在章面之上）。
var sampleWatermarkColor = color.NRGBA{R: 0x80, G: 0x80, B: 0x80, A: 0xB4}

// ApplySampleWatermark 预览时（m.ShowWatermark 为 true）在章面中心叠加「样章」标识。
func ApplySampleWatermark(img *image.RGBA, m base.SealPaintModel) error {
	if img == nil || !m.ShowWatermark {
		return nil
	}
	return DrawSampleWatermark(img, m)
}

// DrawSampleWatermark 在章面中心水平绘制「样章」标识。
// 标识文字取 m.WatermarkText，为空时使用 DefaultSampleWatermarkText。
func DrawSampleWatermark(img *image.RGBA, m base.SealPaintModel) error {
	text := m.WatermarkText
	if text == "" {
		text = DefaultSampleWatermarkText
	}
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}
	// 字号按章面宽度自适应：每字约占章宽的 1/(字数+1)，并以章高的 1/4 为上限，避免标识过大。
	sizePx := float64(m.Width) / float64(len(runes)+1)
	if maxPx := float64(m.Height) / 4; sizePx > maxPx {
		sizePx = maxPx
	}
	if sizePx < 12 {
		sizePx = 12
	}
	face, err := base.LoadFontFaceHinting(m.FontFile, sizePx, m.Dpi, font.HintingFull)
	if err != nil {
		return err
	}
	defer face.Close()

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(sampleWatermarkColor),
		Face: face,
	}
	metrics := face.Metrics()
	ascent, descent := metrics.Ascent.Ceil(), metrics.Descent.Ceil()
	textWidth := d.MeasureString(text).Ceil()
	d.Dot = fixed.P((m.Width-textWidth)/2, (m.Height+ascent-descent)/2)
	d.DrawString(text)
	return nil
}

// Finish 章面收尾：叠加「样章」标识（预览时）后整体透明化。
func Finish(img *image.RGBA, m base.SealPaintModel) (*image.RGBA, error) {
	if err := ApplySampleWatermark(img, m); err != nil {
		return nil, err
	}
	return Translucent(img, m.Transparency), nil
}
