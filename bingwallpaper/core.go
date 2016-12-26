package bingwallpaper

import (
	"image/jpeg"
	"os"

	"github.com/genzj/gobingwallpaper/i18n"
	"github.com/genzj/gobingwallpaper/log"
	"github.com/genzj/gobingwallpaper/util"
)

func SimpleDownload() error {
	var imageMeta ImageMeta

	provider := DefaultBaseProvider
	provider.Market = "zh-cn"
	err := util.HTTPGetJSON(provider.FormURL().String(), &imageMeta)

	if err != nil {
		log.Fatal(err)
		return err
	}
	log.WithFields(log.Fields{
		"ImagesCount": len(imageMeta.Images),
	}).Info(i18n.T("image_info_parsed"))

	url := provider.GetBaseURL()
	image := imageMeta.Images[1]
	url.Path = image.URL
	log.WithFields(log.Fields{
		"Start":               image.StartDate,
		"Title":               image.Copyright,
		"Video":               image.Video,
		"Wallpaper Available": image.WallpaperAvailable,
		"URL": url,
	}).Info(i18n.T("image_download_begin"))
	img, err := util.HTTPGetJpeg(url.String())
	if err != nil {
		log.Fatal(err)
		return err
	}
	log.WithFields(log.Fields{
		"width":  img.Bounds().Size().X,
		"height": img.Bounds().Size().Y,
	}).Info(i18n.T("image_downloaded"))

	f, err := os.OpenFile("out.jpg", os.O_WRONLY|os.O_CREATE, 0644)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
		return err
	}

	err = jpeg.Encode(f, img, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}
