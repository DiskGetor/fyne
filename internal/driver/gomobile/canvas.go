package gomobile

import (
	"image"

	"fyne.io/fyne"
	"fyne.io/fyne/internal/painter/gl"
	"fyne.io/fyne/theme"
)

type canvas struct {
	content, overlay fyne.CanvasObject
	painter          gl.Painter
	scale            float32
	size             fyne.Size

	focused fyne.Focusable
	padded  bool

	typedRune func(rune)
	typedKey  func(event *fyne.KeyEvent)
	shortcut  fyne.ShortcutHandler

	inited, dirty bool
	refreshQueue  chan fyne.CanvasObject
}

func (c *canvas) Content() fyne.CanvasObject {
	return c.content
}

func (c *canvas) SetContent(content fyne.CanvasObject) {
	c.content = content

	if c.padded {
		content.Resize(c.Size().Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
		content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	} else {
		content.Resize(c.Size())
		content.Move(fyne.NewPos(0, 0))
	}
}

func (c *canvas) Refresh(obj fyne.CanvasObject) {
	select {
	case c.refreshQueue <- obj:
		// all good
	default:
		// queue is full, ignore
	}
	c.dirty = true
}

func (c *canvas) Resize(size fyne.Size) {
	if size == c.size {
		return
	}

	c.size = size
	if c.padded {
		c.content.Resize(c.Size().Subtract(fyne.NewSize(theme.Padding()*2, theme.Padding()*2)))
		c.content.Move(fyne.NewPos(theme.Padding(), theme.Padding()))
	} else {
		c.content.Resize(c.Size())
		c.content.Move(fyne.NewPos(0, 0))
	}

	if c.overlay != nil {
		c.overlay.Resize(size)
	}
}

func (c *canvas) Focus(obj fyne.Focusable) {
	if c.focused != nil {
		c.focused.FocusLost()
	}

	c.focused = obj
	if obj != nil {
		obj.FocusGained()
	}
}

func (c *canvas) Unfocus() {
	if c.focused != nil {
		c.focused.FocusLost()
	}
	c.focused = nil
}

func (c *canvas) Focused() fyne.Focusable {
	return c.focused
}

func (c *canvas) Size() fyne.Size {
	return c.size
}

func (c *canvas) Scale() float32 {
	return c.scale
}

func (c *canvas) SetScale(scale float32) {
	if scale == fyne.SettingsScaleAuto {
		scale = c.detectScale()
	}
	c.scale = scale
}

func (c *canvas) detectScale() float32 {
	return 2 // TODO real detection
}

func (c *canvas) Overlay() fyne.CanvasObject {
	return c.overlay
}

func (c *canvas) SetOverlay(overlay fyne.CanvasObject) {
	c.overlay = overlay

	if c.overlay != nil {
		c.overlay.Resize(c.size)
	}
}

func (c *canvas) OnTypedRune() func(rune) {
	return c.typedRune
}

func (c *canvas) SetOnTypedRune(typed func(rune)) {
	c.typedRune = typed
}

func (c *canvas) OnTypedKey() func(*fyne.KeyEvent) {
	return c.typedKey
}

func (c *canvas) SetOnTypedKey(typed func(*fyne.KeyEvent)) {
	c.typedKey = typed
}

func (c *canvas) AddShortcut(shortcut fyne.Shortcut, handler func(shortcut fyne.Shortcut)) {
	c.shortcut.AddShortcut(shortcut, handler)
}

func (c *canvas) Capture() image.Image {
	return c.painter.Capture(c)
}

// NewCanvas creates a new gomobile canvas. This is a canvas that will render on a mobile device using OpenGL.
func NewCanvas() fyne.Canvas {
	ret := &canvas{padded: true}
	ret.scale = ret.detectScale()
	ret.refreshQueue = make(chan fyne.CanvasObject, 1024)

	return ret
}
