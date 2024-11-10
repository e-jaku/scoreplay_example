package storage

type Storage interface {
	UploadMedia([]byte) error
}
