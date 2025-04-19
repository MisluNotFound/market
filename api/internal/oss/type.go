package oss

import "time"

type FileInfo struct {
	Size         int64
	LastModified time.Time
}

type Path struct {
	Path  string
	IsDir bool
}

type OSS interface {
	// Save saves data into path key
	Save(key string, data []byte) error
	// Load loads data from path key
	Load(key string) ([]byte, error)
	// Exists checks if the data exists in the path key
	Exists(key string) (bool, error)
	// State gets the state of the data in the path key
	State(key string) (FileInfo, error)
	// List lists all the data with the given prefix, and all the paths are absolute paths
	List(prefix string) ([]Path, error)
	// Delete deletes the data in the path key
	Delete(key string) error
}
