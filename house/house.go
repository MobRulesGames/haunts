package house

import (
  "glop/gui"
)

type Room struct {
  Defname string
  *roomDef
  RoomInst
}

type WallFacing int
const (
  NearLeft WallFacing = iota
  NearRight
  FarLeft
  FarRight
)

type Door struct {
  // Which wall the door is on
  Facing WallFacing

  // How far along this wall the door is located
  Pos int
}

type RoomInst struct {
  // The placement of doors in this room
  Doors []Door

  // The offset of this room on this floor
  X,Y int
}

func (ri *RoomInst) Pos() (x,y int) {
  return ri.X, ri.Y
}

type Floor struct {
  Rooms []*Room
}

type houseDef struct {
  Floors []*Floor

  // The floor that the explorers start on
  Starting_floor int
}

func MakeHouseDef() *houseDef {
  var h houseDef
  return &h
}

type HouseEditor struct {
  *gui.HorizontalTable
  tab *gui.TabFrame
  widgets []tabWidget

  house  *houseDef
  viewer *HouseViewer
}

func (he *HouseEditor) GetViewer() Viewer {
  return he.viewer
}

func (w *HouseEditor) SelectTab(n int) {
  if n < 0 || n >= len(w.widgets) { return }
  if n != w.tab.SelectedTab() {
    w.widgets[w.tab.SelectedTab()].Collapse()
    w.tab.SelectTab(n)
    // w.viewer.SetEditMode(editNothing)
    w.widgets[n].Expand()
  }
}

type houseDataTab struct {
  *gui.VerticalTable

  num_floors *gui.ComboBox
  theme      *gui.ComboBox

  house  *houseDef
  viewer *HouseViewer
}
func makeHouseDataTab(house *houseDef, viewer *HouseViewer) *houseDataTab {
  var hdt houseDataTab
  hdt.VerticalTable = gui.MakeVerticalTable()
  hdt.house = house
  hdt.viewer = viewer

  num_floors_options := []string{ "1 Floor", "2 Floors", "3 Floors", "4 Floors" }
  hdt.num_floors = gui.MakeComboTextBox(num_floors_options, 300)
  hdt.theme = gui.MakeComboTextBox(tags.Themes, 300)

  hdt.VerticalTable.AddChild(hdt.num_floors)
  hdt.VerticalTable.AddChild(hdt.theme)
  return &hdt
}
func (hdt *houseDataTab) Think(ui *gui.Gui, t int64) {
  hdt.VerticalTable.Think(ui, t)
  num_floors := hdt.num_floors.GetComboedIndex() + 1
  if len(hdt.house.Floors) != num_floors {
    for len(hdt.house.Floors) < num_floors {
      hdt.house.Floors = append(hdt.house.Floors, &Floor{})
    }
    if len(hdt.house.Floors) > num_floors {
      hdt.house.Floors = hdt.house.Floors[0 : num_floors]
    }
  }
}
func (hdt *houseDataTab) Collapse() {}
func (hdt *houseDataTab) Expand() {}

func MakeHouseEditorPanel(house *houseDef, datadir string) Editor {
  var he HouseEditor
  he.HorizontalTable = gui.MakeHorizontalTable()
  he.viewer = MakeHouseViewer(house, 62)
  he.HorizontalTable.AddChild(he.viewer)

  r1 := MakeRoom("name")
  r2 := MakeRoom("name")
  r3 := MakeRoom("name")
  r4 := MakeRoom("name")
  r1.X,r1.Y = 0,0
  r2.X,r2.Y = 20,0
  r3.X,r3.Y = 0,15
  r4.X,r4.Y = 20,15
  house.Floors = append(house.Floors, &Floor{ Rooms: []*Room{ r4, r2, r1, r3 }})
  he.widgets = append(he.widgets, makeHouseDataTab(house, he.viewer))
  var tabs []gui.Widget
  for _,w := range he.widgets {
    tabs = append(tabs, w.(gui.Widget))
  }
  he.tab = gui.MakeTabFrame(tabs)
  he.HorizontalTable.AddChild(he.tab)

  return &he
}