package adapter_deconz

type Option func(*Adapter)

func WithBufferCap(cap int) Option {
	return func(a *Adapter) {
		a.bufferCap = cap
	}
}
