package symbols

type Annotatable interface {
	AddAnnotation(annotation AnnotationAttachmentSymbol)
	GetAnnotations() []AnnotationAttachmentSymbol
}
