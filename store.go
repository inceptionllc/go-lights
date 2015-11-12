package lights

import (
	"errors"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

// Data Store implementation

// Store is implemented by data storage providers for persistent
// configuration information.
type Store interface {
	// Read a value from the provided collection with a given ID.
	Read(collection, id string) (string, error)
	// Write a value to the provided collection with a given ID.
	Write(collection, id, value string) error
	// Remove a value from the provided collection with a given ID.
	Remove(collection, id string) error
	// Remove all values from the provided collection.
	RemoveAll(collection string) error
	// Load reads all the values out of a collection.
	Load(colection string) ([]string, error)
}

// FileStore implements the Store interface by storing each value in
// a file named after the item ID and folders for each collection.
// Note that IDs and collections must be file name friendly.
// On disk, the files will have a `.txt` file extension added.
type FileStore struct {
	Base string // The path to the file store base
}

// NewFileStore creates a new file-based Store implementation. Pass in
// an optional base path to use for data storage (otherwise the user's home
// directory is used).
func NewFileStore(base ...string) (*FileStore, error) {
	var root string
	if len(base) > 0 {
		root = base[0]
	} else {
		root = filepath.Join(os.Getenv("HOME"), "data")
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return nil, err
	}
	err = os.MkdirAll(root, 0755)
	if err != nil {
		return nil, err
	}
	return &FileStore{root}, nil
}

// Read a value from the provided collection and ID.
func (f *FileStore) Read(collection, id string) (string, error) {
	name := filepath.Join(f.Base, collection, id+".txt")
	text, err := ioutil.ReadFile(name)
	if err != nil {
		return "", err
	}
	return string(text), nil
}

// Write a value to the provided collection and ID.
func (f *FileStore) Write(collection, id, value string) error {
	base := filepath.Join(f.Base, collection)
	err := os.MkdirAll(base, 0755)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filepath.Join(base, id+".txt"), []byte(value), 0660)
}

// Remove a value from the provided collection and ID.
func (f *FileStore) Remove(collection, id string) error {
	return os.Remove(filepath.Join(f.Base, collection, id+".txt"))
}

// RemoveAll removes all items from a collection.
func (f *FileStore) RemoveAll(collection string) error {
	return os.RemoveAll(filepath.Join(f.Base, collection))
}

// Load all the values for a collection.
func (f *FileStore) Load(collection string) ([]string, error) {
	items := []string{}
	base := filepath.Join(f.Base, collection)
	base, err := filepath.Abs(base)
	if err != nil {
		log.Println("Error getting abs path to", collection, err)
		return nil, err
	}
	log.Println("Loading collection", collection, "from", base)
	filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// Ignore
			log.Println("store ignoring walk error", err)
			return filepath.SkipDir
		}
		/*
			if info.IsDir() {
				log.Println("skipping dir", path)
				return filepath.SkipDir
			}
		*/
		if filepath.Ext(path) == ".txt" {
			log.Println("loading item", path)
			text, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			items = append(items, string(text))
		} else {
			log.Println("skipping non item", path, filepath.Ext(path))
		}
		return nil
	})
	return items, nil

}

// MockStore is used to test services that rely on Store implementations.
type MockStore struct {
	Data map[string]map[string]string
}

// Read a value from the provided collection with a given ID.
func (s *MockStore) Read(collection, id string) (string, error) {
	c, ok := s.Data[collection]
	if !ok {
		return "", errors.New("No collection found " + collection)
	}
	item, ok := c[id]
	if !ok {
		return "", errors.New("No item with ID found " + id)
	}
	return item, nil
}

// Write a value to the provided collection with a given ID.
func (s *MockStore) Write(collection, id, value string) error {
	c, ok := s.Data[collection]
	if ok {
		c[id] = value
	} else {
		c = map[string]string{id: value}
		if s.Data == nil {
			s.Data = map[string]map[string]string{collection: c}
		} else {
			s.Data[collection] = c
		}
	}
	return nil
}

// Remove a value from the provided collection with a given ID.
func (s *MockStore) Remove(collection, id string) error {
	c, ok := s.Data[collection]
	if ok {
		delete(c, id)
	}
	return nil
}

// RemoveAll clears all items from a collection.
func (s *MockStore) RemoveAll(collection string) error {
	delete(s.Data, collection)
	return nil
}

// Load all the values from a collection.
func (s *MockStore) Load(collection string) ([]string, error) {
	items := []string{}
	for _, item := range s.Data[collection] {
		items = append(items, item)
	}
	return items, nil
}

// Reset removes all data from the store.
func (s *MockStore) Reset() {
	s.Data = map[string]map[string]string{}
}
