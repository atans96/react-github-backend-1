package routes

import (
	"bytes"
	"context"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/chai2010/webp"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/image/bmp"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"regexp"
	"server_go/src/service"
	"server_go/src/service/validation"
	"strconv"
	"strings"
	"sync"
	"unicode/utf8"
)

var httpRegex = regexp.MustCompile("^(https:|http:|www\\.)\\S*")

func convertToWebP(src string) ([]byte, int, int) {
	buf := &bytes.Buffer{}
	var img image.Image
	var imgHeight int
	var imgWidth int
	makeNew, err := http.NewRequest("GET", src, nil)
	if err != nil {
		panic(err)
	}
	makeNew.Header.Set("User-Agent", "Chrome: 91")
	request := service.MakeRequest{Req: makeNew}
	request.Request()
	if strings.Contains(request.ContentType, "jpeg") {
		img, _ = jpeg.Decode(bytes.NewReader(request.Result))
	} else if strings.Contains(request.ContentType, "png") {
		img, _ = png.Decode(bytes.NewReader(request.Result))
	} else if strings.Contains(request.ContentType, "bmp") {
		img, _ = bmp.Decode(bytes.NewReader(request.Result))
	} else if strings.Contains(request.ContentType, "gif") {
		img, _ = gif.Decode(bytes.NewReader(request.Result))
	} else {
		imgInfo := &service.SVG{}
		wg := sync.WaitGroup{}
		wg.Add(2)
		go func() {
			img, _ = service.Decode(bytes.NewReader(request.Result))
			wg.Done()
		}()
		go func() {
			if utf8.Valid(request.Result) {
				if err != nil {
					fmt.Println(err)
				}
			}
			wg.Done()
		}()
		wg.Wait()
		height, _ := strconv.Atoi(imgInfo.Height)
		width, _ := strconv.Atoi(imgInfo.Width)
		if !(height > 40 && width > 120) {
			return buf.Bytes(), height, width
		}
		if err = webp.Encode(buf, img, &webp.Options{Lossless: true, Quality: 50}); err != nil {
			fmt.Println(err)
		}
		request.Clear()
		return buf.Bytes(), height, width
	}
	request.Clear()
	if img != nil {
		imgWidth = img.Bounds().Dx()
		imgHeight = img.Bounds().Dy()
		if !(imgHeight > 40 && imgWidth > 120) {
			return buf.Bytes(), imgHeight, imgWidth
		}
		if err = webp.Encode(buf, img, &webp.Options{Lossless: true, Quality: 50}); err != nil {
			return buf.Bytes(), imgHeight, imgWidth
		}
	}
	return buf.Bytes(), imgHeight, imgWidth
}

type imgProcessing struct {
	img        *validation.ImagesMarkDown
	outputChan chan []map[string]interface{}
}

func (handler *imgProcessing) Handle(ctx context.Context) {
	var output []map[string]interface{}
	res, err := http.Get("https://github.com/" + handler.img.Data.Value.FullName)
	defer res.Body.Close()
	for {
		select {
		case <-ctx.Done():
			handler.outputChan <- output
			return
		default:
			if err != nil {
				handler.outputChan <- output
				return
			}
			if res.StatusCode != 200 {
				handler.outputChan <- output
				return
			}

			// Load the HTML document
			doc, err := goquery.NewDocumentFromReader(res.Body)
			if err != nil {
				panic(err)
			}
			// Find the review items
			doc.Find(".markdown-body").Each(func(i int, s *goquery.Selection) {
				// For each item found, get the title
				nodes := s.Find("img[src]")
				for i := range nodes.Nodes {
					single := nodes.Eq(i)
					src, exist := single.Attr("src")
					if exist {
						if !httpRegex.MatchString(src) {
							src = "https://github.com" + src
						}
						webP, height, width := convertToWebP(src)
						if height > 40 && width > 120 && webP != nil {
							mp1 := map[string]interface{}{
								"webP":   webP,
								"height": height,
								"width":  width,
							}
							output = append(output, mp1)
						}
					}
				}
			})
			handler.outputChan <- output
			return
		}
	}
}
func ImagesFromMarkdown(c *fiber.Ctx) error {
	img := new(validation.ImagesMarkDown)
	outputChan := make(chan []map[string]interface{})
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func(c *fiber.Ctx, img *validation.ImagesMarkDown) {
		if err := c.BodyParser(img); err != nil {
			fn := c.Locals("cancelFn").(context.CancelFunc)
			fn()
			cancel()
			return
		}
		if len(string(c.Request().Body())) == 0 {
			fn := c.Locals("cancelFn").(context.CancelFunc)
			fn()
			cancel()
			return
		}
		errs := service.Validate(img)
		if errs != nil {
			panic(errs)
		}
		handler := imgProcessing{
			img:        img,
			outputChan: outputChan,
		}
		go handler.Handle(ctx)
	}(c, img)
	select {
	case output := <-outputChan:
		return c.JSON(fiber.Map{"id": img.Data.ID, "value": output})
	}
}
