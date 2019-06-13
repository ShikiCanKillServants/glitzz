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
)

const (
	errNoComments = "Didn't return a video with comments"
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
		p, err = pornhub.InitRandom()
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
	var err error

	if len(arguments.Arguments) == 0 {
		p, err = pornhub.InitRandom()
	} else {
		p, err = pornhub.InitRandom()
	}
	if err != nil {
		return nil, err
	}
	pm.lastURL = p.URL
	pm.lastTitle = p.Title

	if len(p.Comments) == 0 {
		return nil, errors.New(errNoComments)
	}

	com := p.Comments[rand.Intn(len(p.Comments))]
	res := []string{util.Returntonormal(util.Boldtext(com.Author))}
	if com.Verified {
		res = append(res, "✔")
	}
	res = append(res, " - "+com.Message)
	if com.Score != 0 {
		res = append(res, "👍 "+util.Returntonormal(util.Greentext(strconv.Itoa(com.Score))))
	}

	return []string{strings.Join(res, "")}, nil
}

func (pm *porn) pornHubLast(arguments core.CommandArguments) (result []string, err error) {
	return []string{pm.lastTitle + " - " + pm.lastURL}, err
}
