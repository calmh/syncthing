package fs

import (
	"context"
	"sync"
	"time"
)

type InstrumentedFilesystem struct {
	fs        Filesystem
	mut       sync.Mutex
	durations map[string]time.Duration
	counts    map[string]int
}

func NewInstrumentedFilesystem(fs Filesystem) *InstrumentedFilesystem {
	return &InstrumentedFilesystem{
		fs:        fs,
		durations: make(map[string]time.Duration),
		counts:    make(map[string]int),
	}
}

func (f *InstrumentedFilesystem) Durations() map[string]time.Duration {
	m := make(map[string]time.Duration)
	f.mut.Lock()
	for k, v := range f.durations {
		m[k] = v
	}
	f.mut.Unlock()
	return m
}

func (f *InstrumentedFilesystem) Counts() map[string]int {
	m := make(map[string]int)
	f.mut.Lock()
	for k, v := range f.counts {
		m[k] = v
	}
	f.mut.Unlock()
	return m
}

func (f *InstrumentedFilesystem) instrument(op string) func() {
	t0 := time.Now()
	return func() {
		d := time.Since(t0)
		f.mut.Lock()
		f.durations[op] += d
		f.counts[op]++
		f.mut.Unlock()
	}
}

func (f *InstrumentedFilesystem) Chmod(name string, mode FileMode) error {
	defer f.instrument("Chmod")()
	return f.fs.Chmod(name, mode)
}

func (f *InstrumentedFilesystem) Lchown(name string, uid, gid int) error {
	defer f.instrument("Lchown")()
	return f.fs.Lchown(name, uid, gid)
}

func (f *InstrumentedFilesystem) Chtimes(name string, atime time.Time, mtime time.Time) error {
	defer f.instrument("Chtimes")()
	return f.fs.Chtimes(name, atime, mtime)
}

func (f *InstrumentedFilesystem) Create(name string) (File, error) {
	defer f.instrument("Create")()
	return f.fs.Create(name)
}

func (f *InstrumentedFilesystem) CreateSymlink(target, name string) error {
	defer f.instrument("CreateSymlink")()
	return f.fs.CreateSymlink(target, name)
}

func (f *InstrumentedFilesystem) DirNames(name string) ([]string, error) {
	defer f.instrument("DirNames")()
	return f.fs.DirNames(name)
}

func (f *InstrumentedFilesystem) Lstat(name string) (FileInfo, error) {
	defer f.instrument("Lstat")()
	return f.fs.Lstat(name)
}

func (f *InstrumentedFilesystem) Mkdir(name string, perm FileMode) error {
	defer f.instrument("Mkdir")()
	return f.fs.Mkdir(name, perm)
}

func (f *InstrumentedFilesystem) MkdirAll(name string, perm FileMode) error {
	defer f.instrument("MkdirAll")()
	return f.fs.MkdirAll(name, perm)
}

func (f *InstrumentedFilesystem) Open(name string) (File, error) {
	defer f.instrument("Open")()
	return f.fs.Open(name)
}

func (f *InstrumentedFilesystem) OpenFile(name string, flags int, mode FileMode) (File, error) {
	defer f.instrument("OpenFile")()
	return f.fs.OpenFile(name, flags, mode)
}

func (f *InstrumentedFilesystem) ReadSymlink(name string) (string, error) {
	defer f.instrument("ReadSymlink")()
	return f.fs.ReadSymlink(name)
}

func (f *InstrumentedFilesystem) Remove(name string) error {
	defer f.instrument("Remove")()
	return f.fs.Remove(name)
}

func (f *InstrumentedFilesystem) RemoveAll(name string) error {
	defer f.instrument("RemoveAll")()
	return f.fs.RemoveAll(name)
}

func (f *InstrumentedFilesystem) Rename(oldname, newname string) error {
	defer f.instrument("Rename")()
	return f.fs.Rename(oldname, newname)
}

func (f *InstrumentedFilesystem) Stat(name string) (FileInfo, error) {
	defer f.instrument("Stat")()
	return f.fs.Stat(name)
}

func (f *InstrumentedFilesystem) SymlinksSupported() bool {
	defer f.instrument("SymlinksSupported")()
	return f.fs.SymlinksSupported()
}

func (f *InstrumentedFilesystem) Walk(name string, walkFn WalkFunc) error {
	defer f.instrument("Walk")()
	return f.fs.Walk(name, walkFn)
}

func (f *InstrumentedFilesystem) Watch(path string, ignore Matcher, ctx context.Context, ignorePerms bool) (<-chan Event, error) {
	defer f.instrument("Watch")()
	return f.fs.Watch(path, ignore, ctx, ignorePerms)
}

func (f *InstrumentedFilesystem) Hide(name string) error {
	defer f.instrument("Hide")()
	return f.fs.Hide(name)
}

func (f *InstrumentedFilesystem) Unhide(name string) error {
	defer f.instrument("Unhide")()
	return f.fs.Unhide(name)
}

func (f *InstrumentedFilesystem) Glob(pattern string) ([]string, error) {
	defer f.instrument("Glob")()
	return f.fs.Glob(pattern)
}

func (f *InstrumentedFilesystem) Roots() ([]string, error) {
	defer f.instrument("Roots")()
	return f.fs.Roots()
}

func (f *InstrumentedFilesystem) Usage(name string) (Usage, error) {
	defer f.instrument("Usage")()
	return f.fs.Usage(name)
}

func (f *InstrumentedFilesystem) Type() FilesystemType {
	defer f.instrument("Type")()
	return f.fs.Type()
}

func (f *InstrumentedFilesystem) URI() string {
	defer f.instrument("URI")()
	return f.fs.URI()
}

func (f *InstrumentedFilesystem) SameFile(fi1, fi2 FileInfo) bool {
	defer f.instrument("SameFile")()
	return f.fs.SameFile(fi1, fi2)
}
