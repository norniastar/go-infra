package apiporttype

type Value string

func (m Value) String() string {
	return string(m)
}

const (
	Bg     Value = "bg"     // Bg  后台
	H5     Value = "h5"     // H5  H5
	Inside Value = "inside" // Inside  内部服务
	Mobile Value = "mobile" // Mobile  移动端
	Open   Value = "open"   // Open  第三方
)
