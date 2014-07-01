package gocaa

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"

	"code.google.com/p/go-uuid/uuid"
)

const baseurl = "http://coverartarchive.org"

type CAAClient struct {
	useragent string
	client    http.Client
}

func NewCAAClient(useragent string) (c *CAAClient) {
	c = &CAAClient{useragent: useragent, client: http.Client{}}
	return
}

func buildURL(path string) (url *url.URL) {
	url, err := url.Parse(baseurl)

	if err != nil {
		return
	}

	url.Path = path
	return
}

func (c *CAAClient) get(url *url.URL) (resp *http.Response, err error) {
	req, _ := http.NewRequest("GET", url.String(), nil)
	req.Header.Set("User-Agent", c.useragent)

	resp, err = c.client.Do(req)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return
}

func (c *CAAClient) getAndJson(url *url.URL) (info *CoverArtInfo, err error) {
	resp, err := c.get(url)

	defer resp.Body.Close()

	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	err = json.Unmarshal(body, &info)

	return

}

func (c *CAAClient) getImage(entitytype string, mbid uuid.UUID, imageid string, size int) (image CoverArtImage, err error) {
	var extra string

	if size == Small || size == 250 {
		extra = "-250"
	} else if size == Large || size == 500 {
		extra = "-500"
	} else {
		extra = ""
	}

	url := buildURL(entitytype + "/" + mbid.String() + "/" + imageid + extra)
	resp, err := c.get(url)

	defer resp.Body.Close()

	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = HTTPError{StatusCode: resp.StatusCode, URL: url}
		return
	}

	data, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return
	}

	image.Data = data

	ext := filepath.Ext(resp.Request.URL.String())
	mimetype := mime.TypeByExtension(ext)

	image.Mimetype = mimetype

	return

}

func (c *CAAClient) GetReleaseInfo(mbid uuid.UUID) (info *CoverArtInfo, err error) {
	url := buildURL("release/" + mbid.String())
	info, err = c.getAndJson(url)
	return
}

func (c *CAAClient) GetReleaseFront(mbid uuid.UUID, size int) (image CoverArtImage, err error) {
	image, err = c.getImage("release", mbid, "front", size)
	return
}

func (c *CAAClient) GetReleaseBack(mbid uuid.UUID, size int) (image CoverArtImage, err error) {
	image, err = c.getImage("release", mbid, "back", size)
	return
}

func (c *CAAClient) GetReleaseImage(mbid uuid.UUID, imageid int, size int) (image CoverArtImage, err error) {
	id := strconv.Itoa(imageid)
	image, err = c.getImage("release", mbid, id, size)
	return
}

func (c *CAAClient) GetReleaseGroupInfo(mbid uuid.UUID) (info *CoverArtInfo, err error) {
	url := buildURL("release-group/" + mbid.String())
	info, err = c.getAndJson(url)
	return
}

func (c *CAAClient) GetReleaseGroupFront(mbid uuid.UUID, size int) (image CoverArtImage, err error) {
	if size != Original {
		err = InvalidImageSizeError{Entitytype: "release-group", Size: size}
		return
	}
	image, err = c.getImage("release-group", mbid, "front", size)
	return
}
