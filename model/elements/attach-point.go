package elements

type AttachPointKind string

const (
	AttachPointType          AttachPointKind = "type"
	AttachPointObject        AttachPointKind = "objecttype"
	AttachPointFunction      AttachPointKind = "function"
	AttachPointObjectMethod  AttachPointKind = "objectfunction"
	AttachPointServiceRemote AttachPointKind = "serviceremotefunction"
	AttachPointParameter     AttachPointKind = "parameter"
	AttachPointReturn        AttachPointKind = "return"
	AttachPointService       AttachPointKind = "service"
	AttachPointField         AttachPointKind = "field"
	AttachPointObjectField   AttachPointKind = "objectfield"
	AttachPointRecordField   AttachPointKind = "recordfield"
	AttachPointListener      AttachPointKind = "listener"
	AttachPointAnnotation    AttachPointKind = "annotation"
	AttachPointExternal      AttachPointKind = "external"
	AttachPointVar           AttachPointKind = "var"
	AttachPointConst         AttachPointKind = "const"
	AttachPointWorker        AttachPointKind = "worker"
	AttachPointClass         AttachPointKind = "class"
)

type AttachPoint interface {
	GetPoint() AttachPointKind
	GetSource() bool
}

type attachPointImpl struct {
	point  AttachPointKind
	source bool
}

func NewAttachPoint(point AttachPointKind, source bool) AttachPoint {
	return &attachPointImpl{
		point:  point,
		source: source,
	}
}

func GetAttachmentPoint(value string, source bool) AttachPoint {
	return &attachPointImpl{
		point:  AttachPointKind(value),
		source: source,
	}
}

func (a *attachPointImpl) GetPoint() AttachPointKind {
	return a.point
}

func (a *attachPointImpl) GetSource() bool {
	return a.source
}
