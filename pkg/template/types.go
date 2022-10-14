package template

type logger interface {
	Debugf(string, ...interface{})
}

type metadata struct {
	Target string            `json:"target"`
	Tags   map[string]string `json:"tags"`
}

type metadataHit struct {
	Metadata metadata
	Hitrate  float32
}
