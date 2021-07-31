package main

import (
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/joho/godotenv"
	tb "gopkg.in/tucnak/telebot.v2"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		panic(err)
    }
	// connecting to telegram bot
	b, err := tb.NewBot(tb.Settings{
		URL: "https://api.telegram.org",

		Token:  os.Getenv("TOKEN"),
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		panic(err)
	}

	// inquire
	b.Handle("/start", func(m *tb.Message) {
		b.Send(m.Chat, "Hi, I'm a bot which will help you to search for articles on wikipedia more quickly.\nTo find tedious information write me what you want to search for.\n\nAvailable languages:\n• English\n• Russian\n• German\n• French\n• Chinese\n• Arabic\n\nHave fun using the bot)")
	})

	b.Handle(tb.OnText, func(m *tb.Message) {
		if !stringInSlice("|", m.Text) && !stringInSlice("🌐", m.Text) && !stringInSlice("🔍", m.Text) {
			langList := WikiPageSearch(m.Text)
			// add key board
			if len(langList) != 0 && len(langList) != 1 {
				selector := &tb.ReplyMarkup{ResizeReplyKeyboard: true}

				switch len(langList) {
					case 2: selector.Reply(selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[0])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[1])))
					case 3: selector.Reply(selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[0])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[1])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[2])))
					case 4: selector.Reply(selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[0])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[1])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[2])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[3])))
					case 5: selector.Reply(selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[0])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[1])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[2])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[3])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[4])))
					case 6: selector.Reply(selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[0])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[1])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[2])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[3])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[4])), selector.Row(selector.Text("🌐 " + m.Text + " | " + langList[5])))
				}

				b.Send(m.Chat, "🌐 Your query found pages in these languages.", tb.ModeHTML, selector)	
			} else if len(langList) == 1 {
				selector := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
				PageListURL := strings.Split(WikiRequestst(m.Text, langList[0]), "https://" + langList[0] + ".wikipedia.org/wiki/")[1:]
				PageList := []string{}
				// convert page name 
				for _, td := range PageListURL {
					text, _ := url.QueryUnescape(strings.Split(td, `","`)[0])
					exitText := ""
					for index, textSplit := range strings.Split(text, "_") {
						if index == len(strings.Split(text, "_")) - 1 {
							exitText += textSplit
						} else {
							exitText += textSplit + " "
						}
					}
					if stringInSlice(`"`, exitText) {
						PageList = append(PageList, strings.Split(exitText, `"`)[0])
					} else {
						PageList = append(PageList, exitText)
					}
					
				}
				// add key board
				if len(PageList) != 1 {
					switch len(PageList) {
						case 2: selector.Reply(selector.Row(selector.Text("🔍 " + PageList[0] + " | " + langList[0])), selector.Row(selector.Text("🔍 " + PageList[1] + " | " + langList[0])))
						case 3: selector.Reply(selector.Row(selector.Text("🔍 " + PageList[0] + " | " + langList[0])), selector.Row(selector.Text("🔍 " + PageList[1] + " | " + langList[0])), selector.Row(selector.Text("🔍 " + PageList[2] + " | " + langList[0])))
						case 4: selector.Reply(selector.Row(selector.Text("🔍 " + PageList[0] + " | " + langList[0])), selector.Row(selector.Text("🔍 " + PageList[1] + " | " + langList[0])), selector.Row(selector.Text("🔍 " + PageList[2] + " | " + langList[0])), selector.Row(selector.Text("🔍 " + PageList[3] + " | " + langList[0])))
						case 5: selector.Reply(selector.Row(selector.Text("🔍 " + PageList[0] + " | " + langList[0])), selector.Row(selector.Text("🔍 " + PageList[1] + " | " + langList[0])), selector.Row(selector.Text("🔍 " + PageList[2] + " | " + langList[0])), selector.Row(selector.Text("🔍 " + PageList[3] + " | " + langList[0])), selector.Row(selector.Text("🔍 " + PageList[4] + " | " + langList[0])))
					}
					
					b.Send(m.Chat, "🔍 Your query found pages in these themes.", tb.ModeHTML, selector)
				} else {
					selector := &tb.ReplyMarkup{ReplyKeyboardRemove: true}
					b.Send(m.Chat, WikiParser("https://" + langList[0] + ".wikipedia.org/wiki/" + PageList[0]), tb.ModeHTML, selector)
				}
			} else {
				b.Send(m.Chat, "Page not found.", tb.ModeHTML)	
			}
		} else if stringInSlice("|", m.Text) && !stringInSlice("🔍", m.Text) && stringInSlice("🌐", m.Text) {
			message := strings.Split(m.Text, " | ")
			PageListURL := strings.Split(WikiRequestst(strings.Split(message[0], "🌐 ")[1], message[1]), "https://" + message[1] + ".wikipedia.org/wiki/")[1:]
			PageList := []string{}
			// convert page name
			for _, td := range PageListURL {
				text, _ := url.QueryUnescape(strings.Split(td, `","`)[0])
				exitText := ""
				for index, textSplit := range strings.Split(text, "_") {
					if index == len(strings.Split(text, "_")) - 1 {
						exitText += textSplit
					} else {
						exitText += textSplit + " "
					}
				}
				if stringInSlice(`"`, exitText) {
					PageList = append(PageList, strings.Split(exitText, `"`)[0])
				} else {
					PageList = append(PageList, exitText)
				}
			}
			// add key board
			if len(PageList) != 1 {
				selector := &tb.ReplyMarkup{ResizeReplyKeyboard: true}
				
				switch len(PageList) {
					case 2: selector.Reply(selector.Row(selector.Text("🔍 " + PageList[0] + " | " + message[1])), selector.Row(selector.Text("🔍 " + PageList[1] + " | " + message[1])))
					case 3: selector.Reply(selector.Row(selector.Text("🔍 " + PageList[0] + " | " + message[1])), selector.Row(selector.Text("🔍 " + PageList[1] + " | " + message[1])), selector.Row(selector.Text("🔍 " + PageList[2] + " | " + message[1])))
					case 4: selector.Reply(selector.Row(selector.Text("🔍 " + PageList[0] + " | " + message[1])), selector.Row(selector.Text("🔍 " + PageList[1] + " | " + message[1])), selector.Row(selector.Text("🔍 " + PageList[2] + " | " + message[1])), selector.Row(selector.Text("🔍 " + PageList[3] + " | " + message[1])))
					case 5: selector.Reply(selector.Row(selector.Text("🔍 " + PageList[0] + " | " + message[1])), selector.Row(selector.Text("🔍 " + PageList[1] + " | " + message[1])), selector.Row(selector.Text("🔍 " + PageList[2] + " | " + message[1])), selector.Row(selector.Text("🔍 " + PageList[3] + " | " + message[1])), selector.Row(selector.Text("🔍 " + PageList[4] + " | " + message[1])))
				}

				b.Send(m.Chat, "🔍 Your query found pages in these themes.", tb.ModeHTML, selector)
			} else {
				selector := &tb.ReplyMarkup{ReplyKeyboardRemove: true}
				b.Send(m.Chat, WikiParser("https://" + message[1] + ".wikipedia.org/wiki/" + PageList[0]), tb.ModeHTML, selector)
			}
			
		} else if stringInSlice("|", m.Text) && stringInSlice("🔍", m.Text) {
			selector := &tb.ReplyMarkup{ReplyKeyboardRemove: true}
			message := strings.Split(m.Text, " | ")
			b.Send(m.Chat, WikiParser("https://" + message[1] + ".wikipedia.org/wiki/" + strings.Split(message[0], "🔍 ")[1]), tb.ModeHTML, selector)
		}
	})

	b.Start()
}

func WikiPageSearch(text string) []string {
	// Get the code of the page
	langList := []string{} 

	for _, lang := range [6]string{"en", "ru", "de", "fr", "ar", "zh"} {
		doc := WikiRequestst(text, lang)

		if doc != "Bad Request" && strings.Split(doc, `",[`)[1] != "],[],[]]" {
			langList = append(langList, lang)
		} 
	}
	
	return langList
}

func WikiParser(url string) string {
	// page existence check
	response, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	document, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		panic(err)
	}
	// Page text parsing
	ExitText := strings.Split((document.Find(".mw-parser-output").Find("p").Text()), "\n")
	ClearExitText := ""
	for _, PageText := range ExitText {
		if PageText != "" {
			// Deleting [] in text
			for _, ExitText := range strings.Split(PageText, "[") {
				if strings.Contains(ExitText, "]") {
					ClearExitText += strings.Split(ExitText, "]")[1]
				} else {
					ClearExitText += ExitText
				}
			}
			break
		}
	}
	switch {
		case len(strings.Split(ClearExitText, " ")) > 10: return ClearExitText + "\n\nRead more: " + `<a href="` + url + `">click</a>`
		default: return "Read: " + `<a href="` + url + `">click</a>`
	}
}

func stringInSlice(search string, list string) bool {
    for _, words := range strings.Split(list, "") {
        if words == search {
            return true
        }
    }
    return false
}

func WikiRequestst(text, language string) string {
	// query editing 
	exitText := ""

	for index, textSplit := range strings.Split(text, " ") {
		if index == len(strings.Split(text, " ")) - 1 {
			exitText += textSplit
		} else {
			exitText += textSplit + "_"
		}
	}
	
	// Get the code of the page
	res, err := http.Get("https://" + language + ".wikipedia.org/w/api.php?action=opensearch&search=" + url.QueryEscape(exitText) + "&format=json&limit=5") 
	if err != nil {
		panic(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		panic(err)
	}

	return doc.Text()
}
