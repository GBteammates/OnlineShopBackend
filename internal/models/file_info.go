package models

type FileInfo struct {
	Name       string `json:"Name"`
	Path       string `json:"Path"`
	CreateDate string `json:"CreateDate"`
	ModifyDate string `json:"ModifyDate"`
}