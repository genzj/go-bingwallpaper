package bingwallpaper

type VideoSource []string

type Video struct {
	Sources     []VideoSource
	Image       string
	Caption     string
	CaptionLink string
}

type Image struct {
	StartDate          string
	FullStartDate      string
	EndDate            string
	URL                string
	URLBase            string
	Copyright          string
	CopyrightLink      string
	Quiz               string
	WallpaperAvailable bool  `json:"wp"`
	Video              Video `json:"vid"`
}

type ImageMeta struct {
	Images   []Image
	Tooltips map[string]string
}
