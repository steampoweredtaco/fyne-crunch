package pb2

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"sync"
)

const pinSvg = `<svg id="Layer_1" height="42" width="42" data-name="Layer 1" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 122.48 122.88"><defs><style>.cls-1{fill:#2470bd;}.cls-1,.cls-2{fill-rule:evenodd;}.cls-2{fill:#1a1a1a;}</style></defs><title>push-pin-blue</title><path class="cls-1" d="M121.21,36.53,85.92,1.23c-3-3-7.77.1-9.2,2.74-.24.45.19.86-.2,3.92A46.27,46.27,0,0,1,73.8,19.21L58.11,34.91c-6.27,6.26-15.23,3.48-22.87-.32-1.62-.8-3.69-2.57-5.48-.78l-6.64,6.64a2.49,2.49,0,0,0,0,3.53L78.9,99.76a2.5,2.5,0,0,0,3.53,0l6.64-6.64c1.77-1.77-.49-4.06-1.41-6-3.4-7-6.45-16.41-.78-22.08l16.39-16.39a84.14,84.14,0,0,1,11.35-2.57c3.09-.49,3.47-.11,3.91-.4,2.71-1.74,5.7-6.15,2.68-9.17Z"/><polygon class="cls-2" points="53.48 82.11 40.77 69.4 0 120.96 1.92 122.88 53.48 82.11 53.48 82.11"/></svg>`

type pinRender struct {
	p  *Pin
	bg fyne.CanvasObject
}

func (p *pinRender) Destroy() {
}

func (p *pinRender) Layout(size fyne.Size) {
	p.p.b.Resize(size)
}

func (p *pinRender) MinSize() fyne.Size {
	return p.p.b.MinSize()
}

func (p *pinRender) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{
		p.bg,
		p.p.b,
	}
}

func (p *pinRender) Refresh() {
	p.p.RLock()
	if p.p.pinned {
		p.p.b.Icon = p.p.pinnedIcon
	} else {
		p.p.b.Icon = p.p.unpinnedIcon
	}
	p.p.b.Resize(p.p.b.Size())
	p.p.b.Refresh()
}

type Pin struct {
	widget.BaseWidget
	b *widget.Button
	sync.RWMutex
	unpinnedIcon fyne.Resource
	pinnedIcon   fyne.Resource
	pinned       bool
	item         *PinBoardItem
}

func (p *Pin) CreateRenderer() fyne.WidgetRenderer {
	return &pinRender{p: p}
}

func (p *Pin) SetPinnedIcon(f fyne.Resource) {
	p.Lock()
	p.pinnedIcon = f
	p.Unlock()
	p.Refresh()
	p.item.Refresh()
}

func (p *Pin) SetUnpinnedIcon(f fyne.Resource) {
	p.Lock()
	p.unpinnedIcon = f
	p.Unlock()
	p.Refresh()
	p.item.Refresh()
}

func newPin(item *PinBoardItem) *Pin {
	p := &Pin{}
	p.ExtendBaseWidget(p)
	p.pinnedIcon = fyne.NewStaticResource("pin.svg", []byte(pinSvg[:]))
	p.unpinnedIcon = theme.NewColoredResource(p.pinnedIcon, theme.ColorNamePrimary)
	p.pinnedIcon = theme.NewColoredResource(theme.RadioButtonCheckedIcon(), theme.ColorNamePrimary)
	p.b = widget.NewButton("", func() {
		p.RLock()
		p.pinned = !p.pinned
		p.RUnlock()
		p.Refresh()
		p.item.Refresh()
	})
	p.item = item
	return p
}
