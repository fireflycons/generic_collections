package local

// https://medium.com/@johnsiilver/writing-an-interface-that-only-sub-packages-can-implement-fe36e7511449

type InternalInter interface {
	internalOnly()
}

type InternalImpl struct{}

func (i *InternalImpl) internalOnly() {}
