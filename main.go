package main

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"google.golang.org/appengine"

	"github.com/joho/godotenv"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/line/line-bot-sdk-go/linebot/httphandler"
	"github.com/mjibson/goon"

	"golang.org/x/net/context"

	"github.com/gin-gonic/gin"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/taskqueue"
	"google.golang.org/appengine/urlfetch"
)

type DateSet struct {
	Id   string    `datastore:"-" goon:"id"`
	Date time.Time `datastore:"date"`
}

type Msg struct {
	Type string
	Text string
}
type Req struct {
	To       string
	Messages []Msg
}

var botHandler *httphandler.WebhookHandler

func init() {
	r := gin.New()

	r.StaticFile("/", "./static")
	http.Handle("/", r)

	err := godotenv.Load("line.env")
	if err != nil {
		panic(err)
	}

	botHandler, err = httphandler.New(
		os.Getenv("LINE_BOT_CHANNEL_SECRET"),
		os.Getenv("LINE_BOT_CHANNEL_TOKEN"),
	)
	botHandler.HandleEvents(handleCallback)

	// TODO:
	http.Handle("/callback", botHandler)
	http.HandleFunc("/task", handleTask)

}

func registLatestTime(c *gin.Context, t time.Time) {
	g := goon.NewGoon(c.Request)
	ctx := appengine.NewContext(c.Request)
	post := DateSet{Id: "1", Date: t}
	if _, err := g.Put(&post); err != nil {
		serveError(ctx, c.Writer, err)
		return
	}
}

func serveError(ctx context.Context, w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain")
	io.WriteString(w, "Internal Server Error")
	log.Errorf(ctx, "%v", err)
}

func handleCallback(evs []*linebot.Event, r *http.Request) {
	c := newContext(r)
	ts := make([]*taskqueue.Task, len(evs))
	for i, e := range evs {
		j, err := json.Marshal(e)
		if err != nil {
			errorf(c, "json.Marshal: %v", err)
			return
		}
		data := base64.StdEncoding.EncodeToString(j)
		t := taskqueue.NewPOSTTask("/task", url.Values{"data": {data}})
		ts[i] = t
	}
	taskqueue.AddMulti(c, ts, "")
}

func handleTask(w http.ResponseWriter, r *http.Request) {
	c := newContext(r)
	data := r.FormValue("data")
	if data == "" {
		errorf(c, "No data")
		return
	}

	j, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		errorf(c, "base64 DecodeString: %v", err)
		return
	}

	e := new(linebot.Event)
	err = json.Unmarshal(j, e)
	if err != nil {
		errorf(c, "json.Unmarshal: %v", err)
		return
	}

	bot, err := newLINEBot(c)
	if err != nil {
		errorf(c, "newLINEBot: %v", err)
		return
	}

	bot.PushMessage(e.Source.UserID)
	// latestFeed := getLatestFeed(w, r, "https://queryfeed.net/instagram?q=yui_makino0119")
	// msg := "最新の牧野さんのInstagramの記事はこちらですぉ　" + latestFeed.Link
	latestFeed := getLatestFeed(w, r, "FEED_URL HERE")
	msg := "LINE MSG HERE" + latestFeed.Link

	m := linebot.NewTextMessage(msg)

	logf(c, "EventType: %s\nMessage: %#v", e.Type, e.Message)
	logf(c, "UserIdType: %s\nMessage: %#v", e.Type, e.Source.UserID)

	if _, err = bot.ReplyMessage(e.ReplyToken, m).WithContext(c).Do(); err != nil {
		errorf(c, "ReplayMessage: %v", err)
		return
	}

	w.WriteHeader(200)
}

func newLINEBot(c context.Context) (*linebot.Client, error) {
	return botHandler.NewClient(
		linebot.WithHTTPClient(urlfetch.Client(c)),
	)
}
