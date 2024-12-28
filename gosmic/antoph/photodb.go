package antoph

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/webp"
)

type photoInfo struct {
	Version string `json:"version"`
	Image   struct {
		Name              string `json:"name"`
		BaseName          string `json:"baseName"`
		Permissions       int    `json:"permissions"`
		Format            string `json:"format"`
		FormatDescription string `json:"formatDescription"`
		MimeType          string `json:"mimeType"`
		Class             string `json:"class"`
		Geometry          struct {
			Width  int `json:"width"`
			Height int `json:"height"`
			X      int `json:"x"`
			Y      int `json:"y"`
		} `json:"geometry"`
		Resolution struct {
			X int `json:"x"`
			Y int `json:"y"`
		} `json:"resolution"`
		PrintSize struct {
			X float64 `json:"x"`
			Y float64 `json:"y"`
		} `json:"printSize"`
		Units        string `json:"units"`
		Type         string `json:"type"`
		BaseType     string `json:"baseType"`
		Endianness   string `json:"endianness"`
		Colorspace   string `json:"colorspace"`
		Depth        int    `json:"depth"`
		BaseDepth    int    `json:"baseDepth"`
		ChannelDepth struct {
			Red   int `json:"red"`
			Green int `json:"green"`
			Blue  int `json:"blue"`
		} `json:"channelDepth"`
		Pixels          int `json:"pixels"`
		ImageStatistics struct {
			Overall struct {
				Min               int     `json:"min"`
				Max               int     `json:"max"`
				Mean              float64 `json:"mean"`
				Median            float64 `json:"median"`
				StandardDeviation float64 `json:"standardDeviation"`
				Kurtosis          float64 `json:"kurtosis"`
				Skewness          float64 `json:"skewness"`
				Entropy           float64 `json:"entropy"`
			} `json:"Overall"`
		} `json:"imageStatistics"`
		ChannelStatistics struct {
			Red struct {
				Min               int     `json:"min"`
				Max               int     `json:"max"`
				Mean              float64 `json:"mean"`
				Median            int     `json:"median"`
				StandardDeviation float64 `json:"standardDeviation"`
				Kurtosis          float64 `json:"kurtosis"`
				Skewness          float64 `json:"skewness"`
				Entropy           float64 `json:"entropy"`
			} `json:"red"`
			Green struct {
				Min               int     `json:"min"`
				Max               int     `json:"max"`
				Mean              float64 `json:"mean"`
				Median            int     `json:"median"`
				StandardDeviation float64 `json:"standardDeviation"`
				Kurtosis          float64 `json:"kurtosis"`
				Skewness          float64 `json:"skewness"`
				Entropy           float64 `json:"entropy"`
			} `json:"green"`
			Blue struct {
				Min               int     `json:"min"`
				Max               int     `json:"max"`
				Mean              float64 `json:"mean"`
				Median            int     `json:"median"`
				StandardDeviation float64 `json:"standardDeviation"`
				Kurtosis          float64 `json:"kurtosis"`
				Skewness          float64 `json:"skewness"`
				Entropy           float64 `json:"entropy"`
			} `json:"blue"`
		} `json:"channelStatistics"`
		RenderingIntent string  `json:"renderingIntent"`
		Gamma           float64 `json:"gamma"`
		Chromaticity    struct {
			RedPrimary struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"redPrimary"`
			GreenPrimary struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"greenPrimary"`
			BluePrimary struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"bluePrimary"`
			WhitePrimary struct {
				X float64 `json:"x"`
				Y float64 `json:"y"`
			} `json:"whitePrimary"`
		} `json:"chromaticity"`
		MatteColor       string `json:"matteColor"`
		BackgroundColor  string `json:"backgroundColor"`
		BorderColor      string `json:"borderColor"`
		TransparentColor string `json:"transparentColor"`
		Interlace        string `json:"interlace"`
		Intensity        string `json:"intensity"`
		Compose          string `json:"compose"`
		PageGeometry     struct {
			Width  int `json:"width"`
			Height int `json:"height"`
			X      int `json:"x"`
			Y      int `json:"y"`
		} `json:"pageGeometry"`
		Dispose     string `json:"dispose"`
		Iterations  int    `json:"iterations"`
		Compression string `json:"compression"`
		Quality     int    `json:"quality"`
		Orientation string `json:"orientation"`
		Properties  struct {
			DateCreate                               time.Time `json:"date:create"`
			DateModify                               time.Time `json:"date:modify"`
			DateTimestamp                            time.Time `json:"date:timestamp"`
			ExifApertureValue                        string    `json:"exif:ApertureValue"`
			ExifBrightnessValue                      string    `json:"exif:BrightnessValue"`
			ExifColorSpace                           string    `json:"exif:ColorSpace"`
			ExifContrast                             string    `json:"exif:Contrast"`
			ExifCustomRendered                       string    `json:"exif:CustomRendered"`
			ExifDateTime                             string    `json:"exif:DateTime"`
			ExifDateTimeDigitized                    string    `json:"exif:DateTimeDigitized"`
			ExifDateTimeOriginal                     string    `json:"exif:DateTimeOriginal"`
			ExifDigitalZoomRatio                     string    `json:"exif:DigitalZoomRatio"`
			ExifExifOffset                           string    `json:"exif:ExifOffset"`
			ExifExifVersion                          string    `json:"exif:ExifVersion"`
			ExifExposureBiasValue                    string    `json:"exif:ExposureBiasValue"`
			ExifExposureMode                         string    `json:"exif:ExposureMode"`
			ExifExposureProgram                      string    `json:"exif:ExposureProgram"`
			ExifExposureTime                         string    `json:"exif:ExposureTime"`
			ExifFileSource                           string    `json:"exif:FileSource"`
			ExifFlash                                string    `json:"exif:Flash"`
			ExifFNumber                              string    `json:"exif:FNumber"`
			ExifFocalLength                          string    `json:"exif:FocalLength"`
			ExifFocalLengthIn35MmFilm                string    `json:"exif:FocalLengthIn35mmFilm"`
			ExifFocalPlaneResolutionUnit             string    `json:"exif:FocalPlaneResolutionUnit"`
			ExifFocalPlaneXResolution                string    `json:"exif:FocalPlaneXResolution"`
			ExifFocalPlaneYResolution                string    `json:"exif:FocalPlaneYResolution"`
			ExifLensModel                            string    `json:"exif:LensModel"`
			ExifLensSpecification                    string    `json:"exif:LensSpecification"`
			ExifLightSource                          string    `json:"exif:LightSource"`
			ExifMake                                 string    `json:"exif:Make"`
			ExifMaxApertureValue                     string    `json:"exif:MaxApertureValue"`
			ExifMeteringMode                         string    `json:"exif:MeteringMode"`
			ExifModel                                string    `json:"exif:Model"`
			ExifOffsetTime                           string    `json:"exif:OffsetTime"`
			ExifOffsetTimeDigitized                  string    `json:"exif:OffsetTimeDigitized"`
			ExifOffsetTimeOriginal                   string    `json:"exif:OffsetTimeOriginal"`
			ExifPhotographicSensitivity              string    `json:"exif:PhotographicSensitivity"`
			ExifRecommendedExposureIndex             string    `json:"exif:RecommendedExposureIndex"`
			ExifSaturation                           string    `json:"exif:Saturation"`
			ExifSceneCaptureType                     string    `json:"exif:SceneCaptureType"`
			ExifSceneType                            string    `json:"exif:SceneType"`
			ExifSensitivityType                      string    `json:"exif:SensitivityType"`
			ExifSharpness                            string    `json:"exif:Sharpness"`
			ExifShutterSpeedValue                    string    `json:"exif:ShutterSpeedValue"`
			ExifSoftware                             string    `json:"exif:Software"`
			ExifThumbnailCompression                 string    `json:"exif:thumbnail:Compression"`
			ExifThumbnailJPEGInterchangeFormat       string    `json:"exif:thumbnail:JPEGInterchangeFormat"`
			ExifThumbnailJPEGInterchangeFormatLength string    `json:"exif:thumbnail:JPEGInterchangeFormatLength"`
			ExifThumbnailResolutionUnit              string    `json:"exif:thumbnail:ResolutionUnit"`
			ExifThumbnailXResolution                 string    `json:"exif:thumbnail:XResolution"`
			ExifThumbnailYResolution                 string    `json:"exif:thumbnail:YResolution"`
			ExifWhiteBalance                         string    `json:"exif:WhiteBalance"`
			IccCopyright                             string    `json:"icc:copyright"`
			IccDescription                           string    `json:"icc:description"`
			IccManufacturer                          string    `json:"icc:manufacturer"`
			IccModel                                 string    `json:"icc:model"`
			JpegColorspace                           string    `json:"jpeg:colorspace"`
			JpegSamplingFactor                       string    `json:"jpeg:sampling-factor"`
			Signature                                string    `json:"signature"`
		} `json:"properties"`
		Profiles struct {
			EightBim struct {
				Length int `json:"length"`
			} `json:"8bim"`
			Exif struct {
				Length int `json:"length"`
			} `json:"exif"`
			Icc struct {
				Length int `json:"length"`
			} `json:"icc"`
			Iptc map[string]json.RawMessage `json:"iptc"`
			Xmp  struct {
				Length int `json:"length"`
			} `json:"xmp"`
		} `json:"profiles"`
		Tainted         bool   `json:"tainted"`
		Filesize        string `json:"filesize"`
		NumberPixels    string `json:"numberPixels"`
		PixelsPerSecond string `json:"pixelsPerSecond"`
		UserTime        string `json:"userTime"`
		ElapsedTime     string `json:"elapsedTime"`
		Version         string `json:"version"`
	} `json:"image"`
}

func parseExifDate(s string) time.Time {
	t, err := time.Parse("2006:01:02 15:04:05", s)
	if err != nil {
		panic(err)
	}
	return t
}

func extractMeta(p string) ImgMeta {
	js, err := os.Open(path.Join(p, "info.json"))
	if err != nil {
		panic(err)
	}
	defer js.Close()

	img, err := os.Open(path.Join(p, "l.webp"))
	if err != nil {
		panic(err)
	}
	defer img.Close()
	imgCfg, err := webp.DecodeConfig(img)
	if err != nil {
		panic(err)
	}

	var infos []photoInfo
	err = json.NewDecoder(js).Decode(&infos)
	if err != nil {
		panic(err)
	}

	props := infos[0].Image.Properties
	iptc := infos[0].Image.Profiles.Iptc
	var keywords []string
	if iptcKeywords, found := iptc["Keyword[2,25]"]; found {
		if err := json.Unmarshal(iptcKeywords, &keywords); err != nil {
			panic(err)
		}
	}

	var dt time.Time
	d := props.ExifDateTimeOriginal
	if len(d) == 0 {
		d = props.ExifDateTime
	}
	if len(d) == 0 {
		dt = props.DateCreate
	} else {
		dt = parseExifDate(d)
	}

	camera := fmt.Sprintf("%s %s", props.ExifMake, props.ExifModel)
	lens := props.ExifLensModel
	iso := props.ExifPhotographicSensitivity
	aperture := parseNumber(props.ExifFNumber)
	shutterSpeed := props.ExifExposureTime

	return ImgMeta{
		Width:        imgCfg.Width,
		Height:       imgCfg.Height,
		Date:         dt,
		Camera:       camera,
		Lens:         lens,
		ISO:          iso,
		ShutterSpeed: shutterSpeed,
		Aperture:     aperture,
		Keywords:     keywords,
	}
}

func parseNumber(s string) string {
	if len(s) == 0 {
		return ""
	}
	parts := strings.SplitN(s, "/", 2)
	n, err := strconv.Atoi(parts[0])
	if err != nil {
		panic(err)
	}
	d, err := strconv.Atoi(parts[1])
	if err != nil {
		panic(err)
	}
	return strconv.Itoa(n / d)
}
