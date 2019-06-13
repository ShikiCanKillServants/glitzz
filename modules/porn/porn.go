package porn

import (
	"errors"
	"github.com/lovelaced/glitzz/config"
	"github.com/lovelaced/glitzz/core"
	"github.com/lovelaced/glitzz/modules/porn/pornhub"
	"github.com/lovelaced/glitzz/modules/porn/pornmd"
	"github.com/lovelaced/glitzz/util"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

// New registers the porn module.
func New(sender core.Sender, conf config.Config) (core.Module, error) {
	rv := &porn{
		Base: core.NewBase("porn", sender, conf),
	}

	rv.AddCommand("pornmd", rv.pornMD)
	rv.AddCommand("porn", rv.pornHubComment)
	rv.AddCommand("porntitle", rv.pornHubTitle)
	rv.AddCommand("pornlast", rv.pornHubLast)
	return rv, nil
}

type porn struct {
	core.Base
	lastURL   string
	lastTitle string
}

func (pm *porn) pornMD(arguments core.CommandArguments) ([]string, error) {
	return pornmd.ReturnRandSearch()
}

func (pm *porn) pornHubTitle(arguments core.CommandArguments) ([]string, error) {
	var p *pornhub.Pornhub
	var err error

	if len(arguments.Arguments) == 0 {
		p, err = pornhub.InitRandom()
	} else {
		query := strings.Join(arguments.Arguments, " ")
		results, err := pornhub.Search(query)

		pm.Log.Info("(Porn) User " + arguments.Nick + " searched " + query)

		if err != nil {
			return nil, err
		}
		if len(results) == 0 {
			return []string{"No results"}, nil
		}
		p, err = pornhub.Init(results[rand.Intn(len(results))].URL)
	}
	if err != nil {
		return nil, err
	}
	pm.lastURL = p.URL
	pm.lastTitle = p.Title

	res := []string{}

	if p.Uploader != "" {
		res = append(res, p.Uploader+" - ")
	}
	res = append(res, p.Title)
	if len(p.Categories) != 0 {
		res = append(res, " | 🔖 "+strings.Join(p.Categories, ", "))
	}
	if p.Views > 0 {
		res = append(res, " | 📸 "+strconv.Itoa(p.Views))
	}
	return []string{strings.Join(res, "")}, nil
}

func (pm *porn) pornHubComment(arguments core.CommandArguments) ([]string, error) {
	var p *pornhub.Pornhub
	var withArgs = len(arguments.Arguments) != 0
	var results []pornhub.Result
	var err error

	if withArgs {
		query := strings.Join(arguments.Arguments, " ")
		results, err = pornhub.Search(query)

		pm.Log.Info("(Porn) User " + arguments.Nick + " searched " + query)

		if err != nil {
			return nil, err
		}
		if len(results) == 0 {
			return []string{"No results"}, nil
		}
	}

	// Try to find a random video with comments in N attempts
	var tries = 5
	var i = 0

	for ; i < tries; i++ {
		if withArgs {
			p, err = pornhub.Init(results[rand.Intn(len(results))].URL)
		} else {
			p, err = pornhub.InitRandom()
		}

		if err != nil {
			return nil, err
		}
		if len(p.Comments) != 0 {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}
	if i == tries {
		return nil, errors.New("Couldn't find a video with comments after " + strconv.Itoa(i) + " attempts")
	}

	pm.lastURL = p.URL
	pm.lastTitle = p.Title

	com := p.Comments[rand.Intn(len(p.Comments))]
	res := []string{util.Returntonormal(util.Boldtext(com.Author))}
	if com.Verified {
		res = append(res, "✔")
	}
	res = append(res, " - "+com.Message)
	if com.Score != 0 {
		res = append(res, "		(Score: "+strconv.Itoa(com.Score)+")")
	}

	return []string{strings.Join(res, "")}, nil
}

func (pm *porn) pornHubLast(arguments core.CommandArguments) (result []string, err error) {
	return []string{pm.lastTitle + " - " + pm.lastURL}, err
}
