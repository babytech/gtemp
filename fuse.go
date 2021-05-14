package main

import (
	"bazil.org/fuse"
	"bazil.org/fuse/fs"
	"bazil.org/fuse/fuseutil"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

var inode uint64

type Node struct {
	inode uint64
	name  string
}

func NewInode() uint64 {
	inode += 1
	return inode
}

type FS struct {
	root *Dir
}

func (f *FS) Root() (fs.Node, error) {
	return f.root, nil
}

type Dir struct {
	Node
	files       *[]*File
	directories *[]*Dir
}

func (d *Dir) Attr(_ context.Context, a *fuse.Attr) error {
	log.Println("Requested Attr for Directory", d.name)
	a.Inode = d.inode
	a.Mode = os.ModeDir | 0444
	return nil
}

func (d *Dir) Lookup(_ context.Context, name string) (fs.Node, error) {
	log.Println("Requested lookup for ", name)
	if d.files != nil {
		for _, n := range *d.files {
			if n.name == name {
				log.Println("Found match for directory lookup with size", len(n.data))
				return n, nil
			}
		}
	}
	if d.directories != nil {
		for _, n := range *d.directories {
			if n.name == name {
				log.Println("Found match for directory lookup")
				return n, nil
			}
		}
	}
	return nil, syscall.ENOENT
}

func (d *Dir) ReadDirAll(context.Context) ([]fuse.Dirent, error) {
	log.Println("Reading all dirs")
	var children []fuse.Dirent
	if d.files != nil {
		for _, f := range *d.files {
			children = append(children, fuse.Dirent{Inode: f.inode, Type: fuse.DT_File, Name: f.name})
		}
	}
	if d.directories != nil {
		for _, dir := range *d.directories {
			children = append(children, fuse.Dirent{Inode: dir.inode, Type: fuse.DT_Dir, Name: dir.name})
		}
		log.Println(len(children), " children for dir", d.name)
	}
	return children, nil
}

func (d *Dir) Create(_ context.Context, req *fuse.CreateRequest, _ *fuse.CreateResponse) (fs.Node, fs.Handle, error) {
	log.Println("Create request for name", req.Name)
	f := &File{Node: Node{name: req.Name, inode: NewInode()}}
	files := []*File{f}
	if d.files != nil {
		files = append(files, *d.files...)
	}
	d.files = &files
	return f, f, nil
}

func (d *Dir) Remove(_ context.Context, req *fuse.RemoveRequest) error {
	log.Println("Remove request for ", req.Name)
	if req.Dir && d.directories != nil {
		var newDirs []*Dir
		for _, dir := range *d.directories {
			if dir.name != req.Name {
				newDirs = append(newDirs, dir)
			}
		}
		d.directories = &newDirs
		return nil
	} else if !req.Dir && *d.files != nil {
		var newFiles []*File
		for _, f := range *d.files {
			if f.name != req.Name {
				newFiles = append(newFiles, f)
			}
		}
		d.files = &newFiles
		return nil
	}
	return syscall.ENONET
}

func (d *Dir) Mkdir(_ context.Context, req *fuse.MkdirRequest) (fs.Node, error) {
	log.Println("Mkdir request for name", req.Name)
	dir := &Dir{Node: Node{name: req.Name, inode: NewInode()}}
	directories := []*Dir{dir}
	if d.directories != nil {
		directories = append(*d.directories, directories...)
	}
	d.directories = &directories
	return dir, nil
}

type File struct {
	Node
	data []byte
	item string
}

func (f *File) Attr(_ context.Context, a *fuse.Attr) error {
	log.Println("Requested Attr for File", f.name, "has data size", len(f.data))
	a.Inode = f.inode
	a.Mode = 0777
	a.Size = uint64(len(f.data))
	return nil
}

func (f *File) Read(_ context.Context, req *fuse.ReadRequest, resp *fuse.ReadResponse) error {
	log.Println("Requested Read on File", f.name)
	fmt.Printf("[FUSE]===> f.data: %v\n", f.data)
	fuseutil.HandleRead(req, resp, f.data)
	return nil
}

func (f *File) ReadAll(context.Context) ([]byte, error) {
	log.Println("[FUSE]===> Reading all of file", f.name)
	content, _ := ReadFile(f.item)
	fmt.Printf("[FUSE]===> File: %s, Content: %s\n", f.item, content)
	str := strings.Replace(content, "\n", "", -1)
	value, err := strconv.Atoi(str)
	fmt.Printf("[FUSE]===> value: %d\n", value)
	if err != nil {
		return nil, err
	}
	degrees := value / 1000
	mDegrees := value % 1000
	strContent := strconv.Itoa(degrees) + "." + strconv.Itoa(mDegrees)
	fmt.Printf("[FUSE]===> strContent: %v\n", strContent)
	f.data = []byte(strContent)
	fmt.Printf("[FUSE]===> f.data: %v\n", f.data)
	return f.data, nil
}

func (f *File) Write(_ context.Context, req *fuse.WriteRequest, resp *fuse.WriteResponse) error {
	log.Println("Trying to write to ", f.name, "offset", req.Offset, "dataSize:", len(req.Data), "data: ", string(req.Data))
	resp.Size = len(req.Data)
	f.data = req.Data
	log.Println("Wrote to file", f.name)
	return nil
}
func (f *File) Flush(context.Context, *fuse.FlushRequest) error {
	log.Println("Flushing file", f.name)
	return nil
}
func (f *File) Open(context.Context, *fuse.OpenRequest, *fuse.OpenResponse) (fs.Handle, error) {
	log.Println("Open call on file", f.name)
	return f, nil
}

func (f *File) Release(context.Context, *fuse.ReleaseRequest) error {
	log.Println("Release requested on file", f.name)
	return nil
}

func (f *File) Fsync(context.Context, *fuse.FsyncRequest) error {
	log.Println("Fsync call on file", f.name)
	return nil
}
