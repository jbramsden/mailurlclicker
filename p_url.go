package backends

import (
	ma "net/mail"
	"strings"
	"github.com/flashmob/go-guerrilla/mail"
	"github.com/mvdan/xurls"
	"io/ioutil"
	//"net/http"
	"context"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/client"
	"time"
)

// ----------------------------------------------------------------------------------
// Processor Name: URLparser
// ----------------------------------------------------------------------------------
// Description   : finds URL in the mail and goes to that url
// ----------------------------------------------------------------------------------
// Config Options: none
// --------------:-------------------------------------------------------------------
// Input         : envelope
// ----------------------------------------------------------------------------------
// Output        : Headers will be populated in e.Header
// ----------------------------------------------------------------------------------
func init() {
	processors["urlparser"] = func() Decorator {
		return URLParser()
	}
}

func URLParser() Decorator {
	return func(p Processor) Processor {
		return ProcessWith(func(e *mail.Envelope, task SelectTask) (Result, error) {
			if task == TaskSaveMail {
				r := strings.NewReader(e.String())
				m, _ := ma.ReadMessage(r)
				body, _ := ioutil.ReadAll(m.Body)
				url := xurls.Relaxed().FindString(string(body))
				Log().Infof("URL: %s",url)

// create context
	ctxt, cancel := context.WithCancel(context.Background())
	defer cancel()

	// create chrome
	c, err := chromedp.New(ctxt, chromedp.WithTargets(client.New().WatchPageTargets(ctxt)), chromedp.WithLog(Log().Infof))
	if err != nil {
		Log().Errorf(err.Error())
	}

	// run task list
	var title string
	err = c.Run(ctxt, text(url, &title))
	if err != nil {
		Log().Errorf(err.Error())
	}
	Log().Infof("Title %s:",title)


			//	res, err := http.Get(url)
			//	if err != nil {
			//		Log().Errorf("could not get url: %s", err.Error())
			//	}else{
			//		output, errn	 := ioutil.ReadAll(res.Body)
			//		if errn != nil {
			//			Log().Errorf("Could not get ioutil: %s", err.Error())
			//		}else{
			//			Log().Infof("HTTP: %s", output)
			//		}
			//	}
			//	res.Body.Close()

				// next processor
				return p.Process(e, task)
			} else {
				// next processor
				return p.Process(e, task)
			}
		})
	}
}

func text(res string, title *string) chromedp.Tasks {
	return chromedp.Tasks{
		chromedp.Navigate(res),
		chromedp.Sleep(2 * time.Second),
		chromedp.Title(title),
	}
}
