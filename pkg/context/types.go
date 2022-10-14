package context

type TargetType string

const (
	TargetTypeFile      TargetType = "file"
	TargetTypeDirectory TargetType = "directory"
)

type Context struct {
	TargetType TargetType
	TargetPath string
	Tags       map[string]string
}

type logger interface {
	Debugf(string, ...interface{})
}
