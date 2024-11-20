package internal

import "fyne.io/fyne/v2"

type BaseRenderer struct {
	objects []fyne.CanvasObject
}

func NewBaseRenderer(objects []fyne.CanvasObject) BaseRenderer {
	return BaseRenderer{objects}
}

func (r *BaseRenderer) Destroy() {
}

func (r *BaseRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *BaseRenderer) SetObjects(objects []fyne.CanvasObject) {
	r.objects = objects
}
