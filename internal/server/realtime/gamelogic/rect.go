package gamelogic

type Rect struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`

	Dynamic bool `json:"dynamic"`
	public  bool
}

func NewRect(x, y, w, h int, dynamic, public bool) *Rect {
	return &Rect{
		X: x,
		Y: y,
		W: w,
		H: h,

		Dynamic: dynamic,
		public:  public,
	}
}
