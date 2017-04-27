package tabs

import (
	"context"

	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
	"honnef.co/go/js/dom"
)

var tabUpdates chan string
var doc dom.Document

func init() {
	tabUpdates = make(chan string, 100)
	doc = dom.GetWindow().Document()
}

type Tabs struct {
	vecty.Core
	IsJS     bool
	Rippled  bool
	Children vecty.MarkupOrComponentOrHTML
	Panels   []*Panel
	Bar      *Bar
	active   string
}

func New(ctx context.Context) *Tabs {
	t := &Tabs{}
	return t.watch(ctx)
}

func switchState(id string, active bool) {
	toActive(id+"-bar", active)
	toActive(id, active)
}

func toActive(id string, state bool) {
	e := doc.GetElementByID(id)
	a := "is-active"
	if state {
		e.Class().Add(a)
	} else {
		e.Class().Remove(a)
	}
}

func (t *Tabs) watch(ctx context.Context) *Tabs {
	go func() {
	done:
		for {
			select {
			case id := <-tabUpdates:
				if id != t.active {
					switchState(t.active, false)
					switchState(id, true)
					t.active = id
				}
			case <-ctx.Done():
				break done
			}
		}
	}()
	return t
}

func (t *Tabs) Render() *vecty.HTML {
	c := make(vecty.ClassMap)
	c["mdl-tabs"] = true
	if t.IsJS {
		c["mdl-js-tabs"] = true
		if t.Rippled {
			c[" mdl-js-ripple-effect"] = true
		}
	}
	if t.Bar == nil || len(t.Bar.Links) == 0 {
		t.Bar = &Bar{}
		for i := 0; i < len(t.Panels); i++ {
			l := &Link{
				ID:   t.Panels[i].ID,
				Name: t.Panels[i].Name,
			}
			if t.Panels[i].IsActive {
				l.IsActive = true
				t.active = t.Panels[i].ID
			}
			t.Bar.Links = append(t.Bar.Links, l)
		}
	}
	var p vecty.List
	for i := 0; i < len(t.Panels); i++ {
		p = append(p, t.Panels[i])
	}
	return elem.Div(
		c,
		vecty.List{t.Bar, p},
	)
}

type Link struct {
	vecty.Core
	ID       string
	Name     string
	IsActive bool
}

func (l *Link) Render() *vecty.HTML {
	c := make(vecty.ClassMap)
	c["mdl-tabs__tab"] = true
	if l.IsActive {
		c["is-active"] = true
	}
	return elem.Anchor(
		prop.Href(l.ID),
		prop.ID(l.ID+"-bar"),
		c,
		vecty.Text(l.Name),
		event.Click(func(e *vecty.Event) {
			go func() {
				tabUpdates <- l.ID
			}()
		}),
	)
}

type Bar struct {
	vecty.Core
	Links    []*Link
	Children vecty.MarkupOrComponentOrHTML
}

func (b *Bar) Render() *vecty.HTML {

	var l vecty.List
	for i := 0; i < len(b.Links); i++ {
		l = append(l, b.Links[i])
	}
	return elem.Div(
		prop.Class("mdl-tabs__tab-bar"), l,
		b.Children,
	)
}

type Panel struct {
	vecty.Core
	IsActive bool
	Name     string
	ID       string
	Children vecty.MarkupOrComponentOrHTML
}

func (p *Panel) Render() *vecty.HTML {
	c := make(vecty.ClassMap)
	c["mdl-tabs__panel "] = true
	if p.IsActive {
		c["is-active"] = true
	}
	return elem.Div(
		c,
		p.Children,
		prop.ID(p.ID),
	)
}

func (p *Panel) Activate() {
	if !p.IsActive {
		p.IsActive = true
	}
}
