package yaac_frontend_main

import (
	"fmt"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	pages "github.com/DHBW-SE-2023/YAAC/internal/frontend/main/pages"
)

// Sidebar
type YaacSidebar struct {
	widget.BaseWidget
	PagesIndex  *map[string][]string
	Pages       *map[string]pages.Page
	Items       *[]*YaacSidebarItem
	BGColor     color.Color
	SelectColor color.Color
	SetPage     func(page pages.Page)
}

func NewYaacSidebar(pagesIndexMap *map[string][]string, pagesMap *map[string]pages.Page, setPage func(page pages.Page), bgColor color.Color, selectColor color.Color) *YaacSidebar {
	item := &YaacSidebar{
		BGColor:     bgColor,
		SelectColor: selectColor,
		SetPage:     setPage,
	}
	item.ExtendBaseWidget(item)
	item.UpdatePages(pagesIndexMap, pagesMap)
	return item
}

func (item *YaacSidebar) UpdatePages(pagesIndexMap *map[string][]string, pagesMap *map[string]pages.Page) {
	item.PagesIndex = pagesIndexMap
	item.Pages = pagesMap
	titles := make([]*YaacSidebarItem, 0, 10)
	for _, p := range *item.Pages {
		titles = append(titles, NewYaacSidebarItem(p.Title, item.SetPage, p, item.BGColor, item.SelectColor))
	}
	item.Items = &titles
}

func (item *YaacSidebar) CreateRenderer() fyne.WidgetRenderer {
	box := container.NewVBox()
	for _, w := range *item.Items {
		box.Add(w)
	}
	rec := canvas.NewRectangle(item.BGColor)

	return widget.NewSimpleRenderer(container.NewStack(rec, box))
}

// Sidebar Button
type YaacSidebarItem struct {
	widget.BaseWidget
	Title       string
	SetPage     func(page pages.Page)
	Page        pages.Page
	BGColor     color.Color
	SelectColor color.Color
}

func NewYaacSidebarItem(title string, setPage func(page pages.Page), page pages.Page, bgColor color.Color, selectColor color.Color) *YaacSidebarItem {
	item := &YaacSidebarItem{
		BGColor:     bgColor,
		SelectColor: selectColor,
		SetPage:     setPage,
		Page:        page,
	}
	item.ExtendBaseWidget(item)
	item.Updateitem(title)
	return item
}

func (item *YaacSidebarItem) Updateitem(title string) {
	item.Title = title
}

func (item *YaacSidebarItem) CreateRenderer() fyne.WidgetRenderer {
	bg := canvas.NewRectangle(item.BGColor)
	txt := canvas.NewText(item.Title, color.Black)
	// FIXME - Design Item here ontop of bg
	// FIXME - Handle Hover event
	c := container.NewStack(
		NewYaacSidebarButton(item.Title, item.SetPage, item.Page),
		container.NewStack(bg, txt),
	)
	return widget.NewSimpleRenderer(c)
}

// Alt. Sidebar Button
type YaacSidebarButton struct {
	widget.Button
	SetPage func(page pages.Page)
	Page    pages.Page
}

func NewYaacSidebarButton(title string, setPage func(page pages.Page), page pages.Page) *YaacSidebarButton {
	item := &YaacSidebarButton{}
	item.ExtendBaseWidget(item)

	item.SetText(title)
	item.SetPage = setPage
	item.Page = page

	return item
}

func (item *YaacSidebarButton) Tapped(_ *fyne.PointEvent) {
	fmt.Println("Click!")
	item.SetPage(item.Page)
}

func (item *YaacSidebarButton) TappedSecondary(_ *fyne.PointEvent) {
	//fmt.Println("Click 2!")
}
