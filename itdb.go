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
	"github.com/mikkyang/id3-go"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
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
	Id         uint
	Title      string
	Album      string
	Artist     string
	Genre      string
	Comment    string
	TrackNr    int
	Year       int
	TrackLen   int
	SampleRate int
}

func (i *IPod) Tracks() (ret []Track, err error) {
	for e := i.db.tracks; e != nil; e = e.next {
		t := (*C.Itdb_Track)(e.data)
		ret = append(ret, Track{
			Id:         uint(t.id),
			Title:      gostring(t.title),
			Album:      gostring(t.album),
			Artist:     gostring(t.artist),
			Genre:      gostring(t.genre),
			Comment:    gostring(t.comment),
			TrackNr:    int(t.track_nr),
			Year:       int(t.year),
			TrackLen:   int(t.tracklen),
			SampleRate: int(t.samplerate),
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

func (i *IPod) RemoveTrack(id uint) error {
	var found *C.Itdb_Track
	for e := i.db.tracks; e != nil; e = e.next {
		t := (*C.Itdb_Track)(e.data)
		if t.id == C.guint32(id) {
			found = t
			break
		}
	}
	if found == nil {
		return errors.New("Unknown track ID")
	}

	C.itdb_playlist_remove_track(nil, found)
	C.itdb_track_remove(found)
	return nil
}

func (i *IPod) CopyTrack(fn string) error {
	f, err := id3.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	ptr := C.CString(fn)
	t := C.itdb_track_new()
	t.userdata = C.gpointer(C.toGstr(ptr))
	t.userdata_duplicate = C.ItdbUserDataDuplicateFunc(C.g_strdup)
	t.userdata_destroy = C.ItdbUserDataDestroyFunc(C.g_free)

	t.transferred = 0
	C.itdb_track_add(i.db, t, -1)

	title := C.CString(f.Title())
	t.title = C.toGstr(title)

	album := C.CString(f.Album())
	t.album = C.toGstr(album)

	artist := C.CString(f.Artist())
	t.artist = C.toGstr(artist)

	genre := C.CString(f.Genre())
	t.genre = C.toGstr(genre)

	comment := C.CString(strings.Join(f.Comments(), "\n"))
	t.comment = C.toGstr(comment)

	t.track_nr = C.gint32(f.Padding())

	year, err := strconv.Atoi(f.Year())
	if err != nil {
		return err
	}
	t.year = C.gint32(year)

	t.tracklen = C.gint32(f.Size())
	t.samplerate = 0

	t.itdb = i.db

	var gerror *C.GError
	C.itdb_cp_track_to_ipod(t, (*C.gchar)(t.userdata), &gerror)
	if gerror != nil {
		return glib.ErrorFromNative(unsafe.Pointer(gerror))
	}

	return nil
}

func (i *IPod) Write() error {
	var gerror *C.GError
	C.itdb_write(i.db, &gerror)
	if gerror != nil {
		return glib.ErrorFromNative(unsafe.Pointer(gerror))
	}
	return nil
}
