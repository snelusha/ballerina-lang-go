package elements

type Point string

const (
	PointType          Point = "type"
	PointObject        Point = "objecttype"
	PointFunction      Point = "function"
	PointObjectMethod  Point = "objectfunction"
	PointServiceRemote Point = "serviceremotefunction"
	PointParameter     Point = "parameter"
	PointReturn        Point = "return"
	PointService       Point = "service"
	PointField         Point = "field"
	PointObjectField   Point = "objectfield"
	PointRecordField   Point = "recordfield"
	PointListener      Point = "listener"
	PointAnnotation    Point = "annotation"
	PointExternal      Point = "external"
	PointVar           Point = "var"
	PointConst         Point = "const"
	PointWorker        Point = "worker"
	PointClass         Point = "class"
)

func (p Point) GetValue() string {
	return string(p)
}

type AttachPoint interface {
	GetPoint() Point
	IsSource() bool
}

type attachPointImpl struct {
	point  Point
	source bool
}

func NewAttachPoint(point Point, source bool) AttachPoint {
	return &attachPointImpl{
		point:  point,
		source: source,
	}
}

func GetAttachmentPoint(value string, source bool) AttachPoint {
	point := Point(value)
	switch point {
	case PointType, PointObject, PointFunction, PointObjectMethod, PointServiceRemote,
		PointParameter, PointReturn, PointService, PointField, PointObjectField,
		PointRecordField, PointListener, PointAnnotation, PointExternal, PointVar,
		PointConst, PointWorker, PointClass:
		return NewAttachPoint(point, source)
	default:
		return nil
	}
}

func (a *attachPointImpl) GetPoint() Point {
	return a.point
}

func (a *attachPointImpl) IsSource() bool {
	return a.source
}
