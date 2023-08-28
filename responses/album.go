package responses

type AlbumDTO struct {
    Status int `json:"status"`
    Message string `json:"message"`
    Data map[string]interface{} `json:"data"`
}

const (
	ALBUM_DELETED_SUCCESSFULLY = "Album deleted successfully"
)