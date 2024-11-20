package pinboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"sync"
)

type PinBoardItem struct {
	sync.RWMutex
	Title      string
	Detail     fyne.CanvasObject
	Pinned     bool
	Expanded   bool
	OnPinned   func(pinned bool)
	OnExpanded func(pinned bool)
	pb         *PinBoard
	container  *fyne.Container
	bLabel     *widget.Button
	pin        *Pin
}

func (p *PinBoardItem) refresh_() {
	if p.Expanded && !p.Detail.Visible() {
		p.Detail.Show()
		p.Detail.Refresh()
	} else if !p.Expanded && p.Detail.Visible() {
		p.Detail.Hide()
		p.Detail.Refresh()
	}
	if p.bLabel.Text != p.Title {
		p.bLabel.SetText(p.Title)
	}
	if p.OnPinned != nil && !p.Pinned && p.pin.pinned {
		p.OnPinned(p.pin.pinned)
	}
	p.Pinned = p.pin.pinned
	p.container.Refresh()
	if p.pb == nil {
		return
	}
	p.pb.Refresh()
}

func (p *PinBoardItem) Refresh() {
	p.RLock()
	p.refresh_()
	p.RUnlock()
}

func NewPinBoardItem(title string, detail fyne.CanvasObject) *PinBoardItem {
	detail.Hide()
	p := &PinBoardItem{
		Title:    title,
		Detail:   detail,
		Expanded: false,
	}
	bLabel := widget.NewButton(title, func() {
		p.RLock()
		p.Expanded = !p.Expanded
		p.RUnlock()
		p.Refresh()
	})
	p.bLabel = bLabel
	pB := newPin(p)
	p.pin = pB
	p.container = container.NewVBox(container.NewBorder(nil, nil, pB, nil, bLabel), detail)
	p.Refresh()
	return p
}
