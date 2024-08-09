package main

type FileHisRecord struct {
	ID      uint
	FileId  string `json:"file_id"`
	Version uint   `json:"version"`
}

type TestTwoReg struct {
}
