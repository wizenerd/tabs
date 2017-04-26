package tabs

import (
	"github.com/gopherjs/vecty"
	"github.com/gopherjs/vecty/elem"
	"github.com/gopherjs/vecty/event"
	"github.com/gopherjs/vecty/prop"
)

type Tabs struct {
	vecty.Core
	IsJS     bool
	Rippled  bool
	Children vecty.MarkupOrComponentOrHTML
	Panels   []*Panel
	Bar      *Bar
	active   string
}

func (t *Tabs) OnActive(id string) {
	t.deactivate(t.active)
	t.activate(id)
	vecty.Rerender(t)
}
func (t *Tabs) deactivate(id string) {
	if id != "" {
		for i := 0; i < len(t.Panels); i++ {
			if t.Panels[i].ID == id {
				t.Panels[i].IsActive = false
				if t.active == id {
					t.active = ""
				}
				t.Bar.Links[i].IsActive = false
			}
		}
	}

}

func (t *Tabs) activate(id string) {
	if id != "" {
		for i := 0; i < len(t.Panels); i++ {
			if t.Panels[i].ID == id {
				t.Panels[i].IsActive = true
				t.Bar.Links[i].IsActive = true
				t.active = id
			}
		}
	}

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
	if t.Bar == nil {
		t.Bar = &Bar{}
		t.Bar.ActiveTab = t.OnActive
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
	OnActive func(id string)
}

func (l *Link) Render() *vecty.HTML {
	c := make(vecty.ClassMap)
	c["mdl-tabs__tab"] = true
	if l.IsActive {
		c["is-active"] = true
	}
	return elem.Anchor(
		prop.Href(l.ID),
		c,
		vecty.Text(l.Name),
		event.Click(func(e *vecty.Event) {
			if !l.IsActive {
				l.IsActive = true
				if l.OnActive != nil {
					l.OnActive(l.ID)
				}
				vecty.Rerender(l)
			}
		}),
	)
}

type Bar struct {
	vecty.Core
	Links     []*Link
	ActiveTab func(string)
	Children  vecty.MarkupOrComponentOrHTML
}

func (b *Bar) Render() *vecty.HTML {
	if b.ActiveTab != nil {
		for i := 0; i < len(b.Links); i++ {
			b.Links[i].OnActive = b.ActiveTab
		}
	}
	var l vecty.List
	for i := 0; i < len(b.Links); i++ {
		l = append(l, b.Links[i])
	}
	return elem.Div(
		prop.Class("mdl-tabs__tab-bar"), l,
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
