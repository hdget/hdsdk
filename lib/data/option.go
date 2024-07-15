package data

type DataOption func(*dataManagerImpl)

func WithTempDir(tempDir string) DataOption {
	return func(impl *dataManagerImpl) {
		impl.tempDir = tempDir
	}
}

func WithFsDir(fsDir string) DataOption {
	return func(impl *dataManagerImpl) {
		impl.fsDir = fsDir
	}
}
