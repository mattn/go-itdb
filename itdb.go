package itdb

/*
#include <gpod/itdb.h>
#include <glib.h>
#include <stdlib.h>

static inline gchar* toGstr(const char* s) { return (gchar*)s; }
static inline char* toCstr(const gchar* s) { return (char*)s; }
static inline void freeCstr(char* s) { free(s); }
*/
// #cgo pkg-config: libgpod-1.0
import "C"
import (
	"errors"
	"github.com/mattn/go-gtk/glib"
	"path/filepath"
	"runtime"
	"unsafe"
)

func gstring(s *C.char) *C.gchar { return C.toGstr(s) }
func cstring(s *C.gchar) *C.char { return C.toCstr(s) }
func gostring(s *C.gchar) string { return C.GoString(cstring(s)) }
func cfree(s *C.char)            { C.freeCstr(s) }

type IPod struct {
	p  string
	db *C.Itdb_iTunesDB
}

func New(mp string) (*IPod, error) {
	mp = filepath.Clean(mp)
	ptr := C.CString(mp)
	defer cfree(ptr)
	var gerror *C.GError
	db := C.itdb_parse(C.toGstr(ptr), &gerror)
	if gerror != nil {
		return nil, glib.ErrorFromNative(unsafe.Pointer(gerror))
	}
	iPod := &IPod{mp, db}
	runtime.SetFinalizer(iPod, func(i *IPod) {
		C.itdb_free(i.db)
	})
	return iPod, nil
}

type Track struct {
	Title  string
	Artist string
}

func (i *IPod) Tracks() (ret []Track, err error) {
	for e := i.db.tracks; e != nil; e = e.next {
		t := (*C.Itdb_Track)(e.data)
		ret = append(ret, Track{
			Title:  gostring(t.title),
			Artist: gostring(t.artist),
		})
	}
	return
}

func (i *IPod) DBPath() (string, error) {
	ptr := C.CString(i.p)
	defer cfree(ptr)
	ret := C.itdb_get_itunesdb_path(C.toGstr(ptr))
	if ret == nil {
		return "", errors.New("iPod Not Found")
	}
	return gostring(ret), nil
}
