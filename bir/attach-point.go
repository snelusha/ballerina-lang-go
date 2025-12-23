package bir

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

func (p Point) String() string {
	return string(p)
}

type AttachPoint struct {
	Point  Point
	Source bool
}

func NewAttachPoint(point Point, source bool) *AttachPoint {
	return &AttachPoint{
		Point:  point,
		Source: source,
	}
}

func GetAttachmentPoint(value string, source bool) *AttachPoint {
	points := []Point{
		PointType, PointObject, PointFunction, PointObjectMethod,
		PointServiceRemote, PointParameter, PointReturn, PointService,
		PointField, PointObjectField, PointRecordField, PointListener,
		PointAnnotation, PointExternal, PointVar, PointConst,
		PointWorker, PointClass,
	}

	for _, point := range points {
		if point.GetValue() == value {
			return NewAttachPoint(point, source)
		}
	}
	return nil
}
