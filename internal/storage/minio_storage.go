package storage

type MinioStorage struct {
	// not sure yet what will be here
}

func NewMinioStorage() Storage {
	return &MinioStorage{}
}

func (s *MinioStorage) UploadMedia(media []byte) error {
	return nil
}
