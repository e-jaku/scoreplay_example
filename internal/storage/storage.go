package storage

type Storage interface {
	UploadMedia(media []byte) error
}
