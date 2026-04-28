package jobs

import (
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"time"
)

type StorageCleaner struct {
	Dirs   []string
	MaxAge time.Duration
}

func StorageCleanup(dirs []string, maxAge time.Duration) *StorageCleaner {
	return &StorageCleaner{
		Dirs:   dirs,
		MaxAge: maxAge,
	}
}

func (s *StorageCleaner) clean() {
	for _, dir := range s.Dirs {
		if dir == "" {
			continue
		}
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				log.Println("⚠️ Error accessing:", path, err)
				return nil
			}
			if d.IsDir() {
				return nil
			}
			info, err := d.Info()
			if err != nil {
				return nil
			}
			if time.Since(info.ModTime()) > s.MaxAge {
				if err := os.Remove(path); err == nil {
					log.Println("Deleted old file:", path)
				} else {
					log.Println("Failed to delete:", path, err)
				}
			}
			return nil
		})
	}
}

func (s *StorageCleaner) Start() {
	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			s.clean()
		}
	}()
}
