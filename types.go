package caa

// CoverArtInfo is the unmarshaled representation of a JSON file in the Cover Art Archive.
// See https://musicbrainz.org/doc/Cover_Art_Archive/API#Cover_Art_Archive_Metadata for an example.
type CoverArtInfo struct {
	Images  []CoverArtImageInfo
	Release string
}

// CoverArtImageInfo is the unmarshaled representation of a single images metadata in a CAA JSON file.
// See https://musicbrainz.org/doc/Cover_Art_Archive/API#Cover_Art_Archive_Metadata for an example.
type CoverArtImageInfo struct {
	Types      []string
	Front      bool
	Back       bool
	Comment    string
	Image      string
	Thumbnails thumbnailMap
	Approved   bool
	Edit       int
	Id         string
}

// CoverArtImage is a wrapper around an image from the CAA, containing its binary data and mimetype information.
type CoverArtImage struct {
	Mimetype string
	Data     []byte
}

type thumbnailMap map[string]string
