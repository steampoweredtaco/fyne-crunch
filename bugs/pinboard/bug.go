package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/huandu/skiplist"
	"math"
	"sync"
)

type Item struct {
	sync.RWMutex
	Expanded  bool
	Title     string
	Detail    fyne.CanvasObject
	pb        *AWidget
	container *fyne.Container
}

func (i *Item) Refresh() {
	i.RLock()
	if i.Expanded {
		i.Detail.Show()
	} else {
		i.Detail.Hide()
	}
	defer i.RUnlock()
	i.container.Refresh()
	if i.pb != nil {
		i.pb.Refresh()
	}
}

func NewItem(title string, detail fyne.CanvasObject) *Item {
	detail.Hide()
	p := &Item{
		Title:  title,
		Detail: detail,
	}
	b := widget.NewButton(p.Title, func() {
		p.Lock()
		p.Expanded = !p.Expanded
		p.Unlock()
		p.Refresh()
	})
	b.Text = p.Title
	b.Icon = theme.IconForWidget(theme.IconNameRadioButton, p.pb)
	b.Importance = widget.MediumImportance
	b.Alignment = widget.ButtonAlignLeading
	b.IconPlacement = widget.ButtonIconLeadingText
	b.Refresh()

	p.container = container.NewVBox(b, detail)

	return p
}

type aRender struct {
	aWidget *AWidget
	objects []fyne.CanvasObject
}

func (p *aRender) Destroy() {
	return
}

func (p *aRender) Layout(size fyne.Size) {
	pos := fyne.NewPos(0, 0)
	i := p.aWidget.items.Front()
	for ; i != nil; i = i.Next() {
		pbi := i.Value.(*Item)
		pbi.container.Move(pos)
		pbi.container.Resize(fyne.NewSize(size.Width, pbi.container.MinSize().Height))
		pos.Y += pbi.Detail.MinSize().Height
		pos.Y += p.aWidget.Theme().Size(theme.SizeNamePadding)
	}
	// for _, obj := range p.objects {
	// 	obj.Move(pos)
	// 	obj.Resize(fyne.NewSize(size.Width, obj.MinSize().Height))
	// 	pos.Y += obj.MinSize().Height
	// 	pos.Y += p.aWidget.Theme().Size(theme.SizeNamePadding)
	// }
}

func (p *aRender) MinSize() fyne.Size {
	if len(p.aWidget.objects) == 0 {
		return fyne.NewSize(0, 0)
	}
	// Padding between elements need to be calculated
	var height float32
	var width float32

	padding := int(theme.Padding()) * (len(p.aWidget.objects) - 1)
	height += float32(padding)

	for _, item := range p.aWidget.objects {
		s := item.MinSize()
		height += s.Height
		width = float32(math.Max(float64(s.Width), float64(width)))
	}
	return fyne.NewSize(width, height)
}

func (p *aRender) Objects() []fyne.CanvasObject {
	return p.objects
}

func (p *aRender) Refresh() {
	for _, obj := range p.aWidget.aItems {
		obj := obj
		b := widget.NewButton(obj.Title, func() {
			obj.Lock()
			obj.Expanded = !obj.Expanded
			obj.Unlock()
			obj.Refresh()
		})

		b.Icon = theme.IconForWidget(theme.IconNameRadioButton, b)
		b.Importance = widget.MediumImportance
		b.Alignment = widget.ButtonAlignLeading
		b.IconPlacement = widget.ButtonIconLeadingText
		b.Refresh()
		b.Show()
		obj.container.Objects[0] = b
		if obj.Expanded {
			obj.container.Objects[1].Show()
		} else {
			obj.container.Objects[1].Hide()
		}
	}
	p.Layout(p.MinSize())
	canvas.Refresh(p.aWidget)
}

type AWidget struct {
	widget.BaseWidget
	sync.RWMutex
	objects []fyne.CanvasObject
	aItems  []*Item
	items   *skiplist.SkipList
}

func (p *AWidget) CreateRenderer() fyne.WidgetRenderer {
	r := &aRender{
		aWidget: p,
		objects: p.objects,
	}
	return r
}

func (p *AWidget) Refresh() {
	for _, item := range p.aItems {
		item.Refresh()
	}
	p.BaseWidget.Refresh()
}

func NewAWidget(items ...*Item) *AWidget {
	ret := &AWidget{}
	ret.ExtendBaseWidget(ret)
	ret.items = skiplist.New(skiplist.Int)

	for _, item := range items {
		item.Detail.Hide()
		item.Lock()
		if item.Expanded {
			item.Detail.Show()
		}
		item.container.Show()
		item.Unlock()
	}

	for i, item := range items {
		item.container.Show()
		ret.objects = append(ret.objects, item.container)
		ret.items.Set(i, item)
	}
	ret.aItems = items
	ret.Show()
	ret.Refresh()
	return ret
}

func main() {
	a := app.New()

	w := a.NewWindow("demo")
	i := NewItem("test2", widget.NewLabel("just a test 2"))
	i.Expanded = true
	i.Refresh()
	aWidget := NewAWidget(
		NewItem("test1", widget.NewLabel("just a test 1")),
		i,
	)
	w.Show()

	w.Resize(fyne.NewSize(800, 600))
	b := widget.NewButton("test just a button", func() {
	})

	b.Icon = theme.IconForWidget(theme.IconNameRadioButton, b)
	b.Importance = widget.MediumImportance
	b.Alignment = widget.ButtonAlignLeading
	b.IconPlacement = widget.ButtonIconLeadingText
	b.Refresh()
	w.SetContent(container.NewVSplit(aWidget, container.NewVBox(b)))

	w.ShowAndRun()
}
