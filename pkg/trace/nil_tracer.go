package trace

type nilTracer struct {
}

func (t *nilTracer) Trace(a ...interface{}) {

}

func Off() Tracer {
	return &nilTracer{}
}
