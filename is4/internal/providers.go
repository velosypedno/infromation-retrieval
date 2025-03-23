package internal

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

type DocsByExtsProvider struct {
	Exts []string
}

func (p *DocsByExtsProvider) Provide(dir string) ([]string, error) {
	extsMap := make(map[string]struct{})
	for _, ext := range p.Exts {
		extsMap[ext] = struct{}{}
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		ext := strings.ToLower(filepath.Ext(entry.Name()))
		if _, ok := extsMap[ext]; ok {
			log.Printf("File %v is suitable\n", filepath.Join(dir, entry.Name()))
			files = append(files, filepath.Join(dir, entry.Name()))
		}
	}
	return files, nil
}

type DocsProvider struct {
	Provider Provider[string, []string]
}

func (p *DocsProvider) Provide(dirs []string) ([]string, error) {
	files := []string{}
	for _, dir := range dirs {
		filesFromDir, err := p.Provider.Provide(dir)
		if err != nil {
			return nil, err
		}
		files = append(files, filesFromDir...)
	}
	return files, nil

}

type ReaderProvider struct {
	ReadersMapper Mapper[string, FileReader]
	ExtProvider   Provider[string, string]
	DefaultReader FileReader
}

func (p *ReaderProvider) Provide(path string) (FileReader, error) {
	ext, err := p.ExtProvider.Provide(path)
	if err != nil {
		return nil, err
	}
	reader := p.ReadersMapper.Map(ext)
	if reader == nil {
		return p.DefaultReader, nil
	}
	return reader, nil
}

type ExtByPathProvider struct{}

func (p *ExtByPathProvider) Provide(path string) (string, error) {
	return filepath.Ext(path), nil
}
