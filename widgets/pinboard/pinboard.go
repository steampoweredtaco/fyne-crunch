package pinboard

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/huandu/skiplist"
	"github.com/steampoweredtaco/fyne-crunch/widgets/internal"
	"image/color"
	"math"
	"sync"
)

// pinboardLayout will be associated with container in the scroller to get scroll
// events and update the layouts of all the objects in the container
type pinboardLayout struct {
	pb        *PinBoard
	container *fyne.Container
	scroller  *container.Scroll

	topBG        *canvas.Rectangle
	topShadow    fyne.CanvasObject
	bottomShadow fyne.CanvasObject
	bottomBG     *canvas.Rectangle
	bg           []*canvas.Rectangle
	fg           []*canvas.Rectangle

	offset float32
}

func (p *pinboardLayout) Layout(objects []fyne.CanvasObject, size fyne.Size) {
	if p.pb == nil {
		return
	}
	p.updateLayout()
}

func (p *pinboardLayout) updateLayout() {
	p.pb.RLock()
	defer p.pb.RUnlock()

	if len(p.fg) < p.pb.pinned.Len() {
		for i := len(p.bg); i < p.pb.items.Len(); i++ {
			p.fg = append(p.fg, canvas.NewRectangle(color.Transparent))
			p.fg[i].StrokeColor = theme.ColorForWidget(theme.ColorNameBackground, p.pb)
			p.fg[i].StrokeWidth = 4
		}
	}
	p.fg = p.fg[:p.pb.pinned.Len()]

	if len(p.bg) < p.pb.pinned.Len() {
		for i := len(p.bg); i < p.pb.items.Len(); i++ {
			p.bg = append(p.bg, canvas.NewRectangle(theme.ColorForWidget(theme.ColorNameBackground, p.pb)))
		}
	}
	p.bg = p.bg[:p.pb.pinned.Len()]

	viewportSize := p.scroller.Size()
	padding := theme.Padding()
	top := fyne.NewPos(0, 0)
	topPinDepth := fyne.NewPos(0, p.offset)
	pinnedTotalHeightRemaining := fyne.NewPos(0, 0)

	if p.pb.pinned.Len() > 0 {
		for pinnedI := p.pb.pinned.Front(); pinnedI != nil; pinnedI = pinnedI.Next() {
			pinnedTotalHeightRemaining.Y += p.pb.minSizes[pinnedI.Key().(int)].Height + padding
		}

		if viewportSize.Subtract(pinnedTotalHeightRemaining).Height <= 0 {
			// simplest case only pinned items can be in view doesn't matter to
			// update all the other items
			for pinnedI := p.pb.pinned.Front(); pinnedI != nil; pinnedI = pinnedI.Next() {
				pbi := pinnedI.Value.(*PinBoardItem)
				pbi.container.Move(topPinDepth)
				topPinDepth.Y += p.pb.minSizes[pinnedI.Key().(int)].Height + padding
			}
			return
		}
	}

	itemPositions := make([]fyne.Position, p.pb.items.Len())

	item := p.pb.items.Front()
	pinnedNotTop := make([]int, 0, p.pb.pinned.Len())
	var topPinnedLen = 0
	for ; item != nil; item = item.Next() {
		k := item.Key().(int)
		pinnedI := p.pb.pinned.Get(k)
		if pinnedI != nil {
			if top.Y <= topPinDepth.Y {
				itemPositions[k] = topPinDepth
				topPinDepth.Y += p.pb.minSizes[k].Height + padding
				top.Y += p.pb.minSizes[k].Height + padding
				pinnedTotalHeightRemaining.Y -= p.pb.minSizes[k].Height - padding
				topPinnedLen++
				continue
			}

			itemPositions[k].Y = top.Y
			top.Y += p.pb.minSizes[k].Height + padding
			pinnedNotTop = append(pinnedNotTop, k)
			continue
		}
		itemPositions[k].Y = top.Y
		top.Y += p.pb.minSizes[k].Height + padding
	}
	if topPinnedLen != 0 {
		p.topShadow.Resize(fyne.NewSize(viewportSize.Width, 0))
		p.topShadow.Move(fyne.NewPos(0, topPinDepth.Y))
		p.topShadow.Show()
	} else {
		p.topShadow.Hide()
	}
	p.topShadow.Refresh()

	p.topBG.Resize(fyne.NewSize(viewportSize.Width, topPinDepth.Y))
	p.topBG.Move(fyne.NewPos(0, 0))
	p.topBG.Show()
	p.topBG.Refresh()

	bottomY := viewportSize.Height
	for i := range pinnedNotTop {
		v := pinnedNotTop[len(pinnedNotTop)-1-i]
		if itemPositions[v].Y-p.offset+p.pb.minSizes[v].Height >= bottomY {
			itemPositions[v].Y = bottomY - p.pb.minSizes[v].Height - padding + p.offset
			bottomY -= p.pb.minSizes[v].Height + padding
		}
	}

	p.bottomBG.Resize(fyne.NewSize(viewportSize.Width, viewportSize.Height-bottomY))
	p.bottomBG.Move(fyne.NewPos(0, bottomY+p.offset))
	p.bottomBG.Show()
	var shadowPos float32 = bottomY
	// Now lower them when the top pins overlapped
	if len(pinnedNotTop) != 0 && itemPositions[pinnedNotTop[len(pinnedNotTop)-1]].Y < topPinDepth.Y {
		offset := topPinDepth.Y - itemPositions[pinnedNotTop[len(pinnedNotTop)-1]].Y + padding
		for _, v := range pinnedNotTop {
			itemPositions[v].Y += offset
			shadowPos = itemPositions[v].Y
		}
	}

	if len(pinnedNotTop) != 0 {
		p.bottomShadow.Resize(fyne.NewSize(viewportSize.Width, 0))
		p.bottomShadow.Move(fyne.NewPos(0, shadowPos+p.offset))
		p.bottomShadow.Show()
	} else {
		p.bottomShadow.Hide()
	}
	p.bottomShadow.Refresh()

	for i, pos := range itemPositions {
		item := p.pb.items.Get(i)
		pbi := item.Value.(*PinBoardItem)
		pbi.container.Move(pos)
		s := p.pb.minSizes[item.Key().(int)]
		pbi.container.Resize(fyne.NewSize(viewportSize.Width, s.Height))
	}
	bgIndex := 0

	for pinnedI := p.pb.pinned.Front(); pinnedI != nil; pinnedI = pinnedI.Next() {
		i := pinnedI.Key().(int)
		s := p.pb.minSizes[i]
		p.bg[bgIndex].Resize(fyne.NewSize(viewportSize.Width+padding/2, s.Height+padding/2))
		p.bg[bgIndex].Move(itemPositions[i])
		p.bg[bgIndex].Refresh()
		p.fg[bgIndex].Resize(fyne.NewSize(viewportSize.Width, s.Height))
		p.fg[bgIndex].Move(itemPositions[i])
		p.fg[bgIndex].Refresh()
		bgIndex++
	}

	objects := []fyne.CanvasObject{}
	for i := p.pb.items.Front(); i != nil; i = i.Next() {
		if _, ok := p.pb.pinned.GetValue(i.Key()); ok {
			continue
		}
		k := i.Key().(int)
		if !((itemPositions[k].Y >= p.offset && itemPositions[k].Y <= p.scroller.Size().Height+p.offset) ||
			(itemPositions[k].Y+p.pb.minSizes[k].Height >= p.offset && itemPositions[k].Y+p.pb.minSizes[k].Height <= p.scroller.Size().Height+p.offset)) {
			// not visable
			continue
		}

		pbi := i.Value.(*PinBoardItem)
		objects = append(objects, pbi.container)
	}
	for _, bg := range p.bg {
		objects = append(objects, bg)
	}
	objects = append(objects, p.topBG, p.bottomBG, p.topShadow, p.bottomShadow)
	for i := p.pb.pinned.Front(); i != nil; i = i.Next() {
		pbi := i.Value.(*PinBoardItem)
		objects = append(objects, pbi.container)
	}
	for _, fg := range p.fg {
		objects = append(objects, fg)
	}
	p.container.Objects = objects
}

func (p *pinboardLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	p.pb.RLock()
	defer p.pb.RUnlock()
	if p.pb.items.Len() == 0 {
		return fyne.NewSize(0, 0)
	}
	var minSize fyne.Size
	for _, s := range p.pb.minSizes {
		minSize.Height += s.Height
		minSize.Width = float32(math.Max(float64(minSize.Width), float64(s.Width)))
	}
	minSize.Height += theme.Padding() * float32(len(p.pb.minSizes))
	return minSize
}

func (p *pinboardLayout) updateOffset(position fyne.Position) {
	if p.offset == position.Y {
		return
	}
	p.offset = position.Y
	p.updateLayout()
}

type pinboardRender struct {
	pb       *PinBoard
	scroller *container.Scroll
	layout   *fyne.Container
	objects  []fyne.CanvasObject
}

func (p *pinboardRender) Destroy() {
	return
}

func (p *pinboardRender) Layout(size fyne.Size) {
	p.scroller.Resize(size)
	// If the resize is smaller than the scroller's content it won't actually
	// cause a Layout trigger to update the bottom pins, so do that here.
	delta := p.scroller.Size().Subtract(p.layout.MinSize())
	if delta.Height < 0 {
		// layout is special and is pinboard aware and doesn't use the parameters.
		p.layout.Layout.Layout(nil, fyne.Size{})
	}
}

func (p *pinboardRender) updateOffset(offset fyne.Position) {
}
func (p *pinboardRender) MinSize() fyne.Size {
	return p.scroller.MinSize()
}

func (p *pinboardRender) Objects() []fyne.CanvasObject {
	return p.objects
}

func (p *pinboardRender) Refresh() {
	if p.scroller == nil {
		canvas.Refresh(p.pb)
		return
	}
	p.Layout(p.pb.Size())
	p.scroller.Refresh()
	l := p.layout.Layout.(*pinboardLayout)
	l.updateLayout()
	canvas.Refresh(p.pb)
}

type PinBoard struct {
	widget.BaseWidget
	sync.RWMutex
	scrollContainer *container.Scroll
	bg              []canvas.Rectangle
	minSizes        []fyne.Size
	items           *skiplist.SkipList
	pinned          *skiplist.SkipList
	expanded        *skiplist.SkipList
}

func (p *PinBoard) CreateRenderer() fyne.WidgetRenderer {
	p.scrollContainer = container.NewVScroll(nil)
	pinboardContainer := &fyne.Container{
		Objects: []fyne.CanvasObject{p.scrollContainer},
	}
	layout := &pinboardLayout{
		pb:           p,
		container:    pinboardContainer,
		scroller:     p.scrollContainer,
		topBG:        canvas.NewRectangle(theme.ColorForWidget(theme.ColorNameBackground, p)),
		bottomBG:     canvas.NewRectangle(theme.ColorForWidget(theme.ColorNameBackground, p)),
		topShadow:    container.NewStack(internal.NewShadow(internal.ShadowBottom, internal.DialogLevel)),
		bottomShadow: container.NewStack(internal.NewShadow(internal.ShadowTop, internal.DialogLevel)),
		offset:       0,
	}
	pinboardContainer.Layout = layout
	pinboardContainer.Resize(layout.MinSize(nil))
	p.scrollContainer.Content = pinboardContainer
	p.scrollContainer.OnScrolled = layout.updateOffset

	r := &pinboardRender{
		pb:       p,
		objects:  []fyne.CanvasObject{p.scrollContainer},
		scroller: p.scrollContainer,
		layout:   pinboardContainer,
	}
	return r
}

func (p *PinBoard) Destroy() {
	defer _l("W: Destroy")()
	for item := p.items.Front(); item != nil; item = item.Next() {
		pbi := item.Value.(*PinBoardItem)
		pbi.Lock()
		pbi.pb = nil
	}
}

func (p *PinBoard) Refresh() {
	func() {
		p.Lock()
		defer p.Unlock()
		pinned := skiplist.New(skiplist.Int)
		for item := p.items.Front(); item != nil; item = item.Next() {
			pbi := item.Value.(*PinBoardItem)
			pbi.RLock()
			if pbi.Pinned {
				pinned.Set(item.Key(), item.Value)
			}
			pbi.RUnlock()
			p.minSizes[item.Key().(int)] = pbi.container.MinSize()
		}
		p.pinned = pinned
	}()
	p.BaseWidget.Refresh()
}

func (p *PinBoard) AddItem(item *PinBoardItem) {
	p.Lock()
	item.pb = p
	key := p.items.Len()
	p.items.Set(key, item)
	p.minSizes = append(p.minSizes, item.container.MinSize())
	p.Unlock()
	p.Refresh()
}

// NewPinBoard returns the PinBoard controller and a canvas object that should be used.
// PinBoard provides its own scroll container so it would not make sense to put it
// directly inside another fyne scroll container.
//
// TODO: It is confusing that PinBoard is also a canvas object which will be broken and
// scrollable and it will not work inside a fyne scroll container as expected. Make this
// return only PinBoard which is the controller and correct canvas object.
func NewPinBoard(pinBoardItems ...*PinBoardItem) *PinBoard {
	ret := &PinBoard{}
	ret.ExtendBaseWidget(ret)
	ret.items = skiplist.New(skiplist.Int)
	ret.pinned = skiplist.New(skiplist.Int)
	ret.expanded = skiplist.New(skiplist.Int)
	ret.minSizes = make([]fyne.Size, len(pinBoardItems))
	for i, item := range pinBoardItems {
		item.Lock()
		ret.items.Set(i, item)
		item.pb = ret
		if item.Pinned {
			ret.pinned.Set(i, item)
		}
		if item.Expanded {
			ret.expanded.Set(i, item)
		}
		item.container.Show()
		item.Unlock()
		ret.minSizes[i] = item.container.MinSize()
	}
	ret.Show()
	return ret
}
