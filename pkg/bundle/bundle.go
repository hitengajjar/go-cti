package bundle

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"

	"github.com/acronis/go-cti/pkg/collector"
	"github.com/acronis/go-cti/pkg/filesys"
)

type Bundle struct {
	Index     *Index
	IndexLock *IndexLock

	Registry *collector.CtiRegistry

	BaseDir string
}

// New creates a new bundle from the specified path.
// If the path is empty, the current working directory is used.
func New(path string, options ...InitializeOption) *Bundle {
	b := &Bundle{
		BaseDir: path,
		Index:   &Index{},
		IndexLock: &IndexLock{
			Version:          IndexLockVersion,
			DependentBundles: make(map[string]string),
			SourceInfo:       make(map[string]Info),
		},
	}

	for _, opt := range options {
		opt(b)
	}

	return b
}

type InitializeOption func(*Bundle)

func WithAppCode(appCode string) InitializeOption {
	return func(b *Bundle) {
		// TODO: Validate appCode
		b.Index.AppCode = appCode
	}
}

func WithRamlxVersion(version string) InitializeOption {
	return func(b *Bundle) {
		b.Index.RamlxVersion = version
	}
}
func WithEntities(entities []string) InitializeOption {
	return func(b *Bundle) {
		if entities != nil {
			b.Index.Entities = entities
		}
	}
}

func (b *Bundle) Read() error {
	idx, err := ReadIndex(b.BaseDir)
	if err != nil {
		return fmt.Errorf("read index file: %w", err)
	}
	idxLock, err := ReadIndexLock(b.BaseDir)
	if err != nil {
		return fmt.Errorf("read index lock: %w", err)
	}

	b.Index = idx
	b.IndexLock = idxLock
	return nil
}

func (b *Bundle) SaveIndexLock() error {
	if err := b.IndexLock.Save(b.BaseDir); err != nil {
		return fmt.Errorf("save index lock: %w", err)
	}
	return nil
}

func (b *Bundle) SaveIndex() error {
	if err := b.Index.Save(b.BaseDir); err != nil {
		return fmt.Errorf("save index: %w", err)
	}
	return nil
}

func (b *Bundle) GetDictionaries() (Dictionaries, error) {
	dictionaries := Dictionaries{
		Dictionaries: make(map[LangCode]Entry),
	}

	for _, dict := range b.Index.Dictionaries {
		file, err := os.Open(path.Join(b.BaseDir, dict))
		if err != nil {
			return Dictionaries{}, fmt.Errorf("open dictionary file: %w", err)
		}
		defer file.Close()

		entry, err := validateDictionary(file)
		if err != nil {
			return Dictionaries{}, fmt.Errorf("validate dictionary: %w", err)
		}
		lang := filesys.GetBaseName(file.Name())
		dictionaries.Dictionaries[LangCode(lang)] = entry
	}

	return dictionaries, nil
}

func validateDictionary(input io.Reader) (Entry, error) {
	decoder := json.NewDecoder(input)

	var config Entry
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("decode dictionary: %w", err)
	}

	return config, nil
}
