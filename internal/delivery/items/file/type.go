package file

type FilesInfo struct {
	Name       string `json:"name" example:"20221213125935.jpeg"`
	Path       string `json:"path" example:"storage\\files\\categories\\d0d3df2d-f6c8-4956-9d76-998ee1ec8a39\\20221213125935.jpeg"`
	CreateDate string `json:"createdDate" example:"2022-12-13 12:46:16.0964549 +0300 MSK"`
	ModifyDate string `json:"modifyDate" example:"2022-12-13 12:46:16.0964549 +0300 MSK"`
}

type FileListResponse struct {
	Files []FilesInfo `json:"files"`
}
