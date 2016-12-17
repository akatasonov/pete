package main

import (
    "github.com/valyala/fasthttp"
    "html/template"
    "log"
)

type Image struct {
    Name string
    Url string
}

func main() {
    t := template.Must(template.ParseFiles("./views/app.html",
            "./views/gallery.html", "./views/image.html"))

    // Setup FS handler
    fs := &fasthttp.FS {
        Root:               "./",
        IndexNames:         []string{"index.html"},
        GenerateIndexPages: false,
        Compress:           false,
        AcceptByteRange:    false,
    }

    fsHandler := fs.NewRequestHandler()

    requestHandler := func (t *template.Template) func (ctx *fasthttp.RequestCtx) {
        return func(ctx * fasthttp.RequestCtx) {
            switch string(ctx.Path()) {
                case "/gallery":
                    ctx.SetContentType("text/html; charset=utf8")

                    data := map[int]Image {
                       0 : {"pete_1", "https://flic.kr/p/9Yd37s"},
                        1 : {"pete_2", "https://flic.kr/p/xjFhnR"},
                        2 : {"pete_3", "https://flic.kr/p/nVP5fh"},
                        3 : {"pete_4", "https://flic.kr/p/vvC6hq"},
                    }

                    err := t.ExecuteTemplate(ctx, "app", data)
                    if err != nil {
                        panic(err)
                    }
                case "/assets":
                default:
                    fsHandler(ctx)
            }
        }
    }

    if err := fasthttp.ListenAndServe(":8080", requestHandler(t)); err != nil {
        log.Fatalf("Error in ListenAndServe: %s", err)
    }
}
