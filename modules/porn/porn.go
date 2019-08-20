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
	rv.AddCommand("porn", rv.pornHubTitleShort)
	rv.AddCommand("pornfull", rv.pornHubTitleFull)
	rv.AddCommand("pornsay", rv.pornHubComment)
	rv.AddCommand("pornlast", rv.pornHubLast)

	rv.AddCommand("gayporn", rv.gayPornHubTitleShort)
	rv.AddCommand("gaypornfull", rv.gayPornHubTitleFull)
	rv.AddCommand("gaypornsay", rv.gayPornHubComment)
	return rv, nil
}

type porn struct {
	core.Base
	lastURL   string
	lastTitle string
	lastCall  time.Time
}

func (pm *porn) pornMD(arguments core.CommandArguments) ([]string, error) {
	return pornmd.ReturnRandSearch()
}

func (pm *porn) isSpam() (spam bool) {
	spam = !pm.lastCall.IsZero() && time.Since(pm.lastCall) < (5*time.Second)
	if !spam {
		pm.lastCall = time.Now()
	}
	return
}

func (pm *porn) pornHubTitle(arguments core.CommandArguments) (*pornhub.Pornhub, error) {
	var gay = len(arguments.Arguments) > 0 && strings.ToUpper(arguments.Arguments[0]) == "GAY"
	var p *pornhub.Pornhub
	var err error

	if len(arguments.Arguments) == 0 || (gay && len(arguments.Arguments) == 1) {
		p, err = pornhub.InitRandom(gay)
	} else {
		var query string
		if !gay {
			query = strings.Join(arguments.Arguments, " ")
		} else {
			query = strings.Join(arguments.Arguments[1:], " ")
		}
		results, err := pornhub.Search(query, gay)

		pm.Log.Info("User " + arguments.Nick + " searched " + query)

		if err != nil {
			return nil, err
		}
		if len(results) == 0 {
			return nil, errors.New("No results")
		}
		p, err = pornhub.Init(results[rand.Intn(len(results))].URL)
	}
	if err != nil {
		return nil, err
	}
	pm.lastURL = p.URL
	pm.lastTitle = p.Title

	return p, nil
}

func (pm *porn) pornHubTitleShort(arguments core.CommandArguments) ([]string, error) {
	if pm.isSpam() {
		return []string{"Stop flooding"}, nil
	}

	p, err := pm.pornHubTitle(arguments)
	if err != nil {
		return nil, err
	}

	res := []string{}
	if p.Uploader != "" {
		res = append(res, p.Uploader+" - ")
	}
	res = append(res, p.Title)

	return []string{strings.Join(res, "")}, nil
}

func (pm *porn) pornHubTitleFull(arguments core.CommandArguments) ([]string, error) {
	if pm.isSpam() {
		return []string{"Stop flooding"}, nil
	}

	p, err := pm.pornHubTitle(arguments)
	if err != nil {
		return nil, err
	}

	res := []string{}
	if p.Uploader != "" {
		res = append(res, p.Uploader+" - ")
	}
	res = append(res, p.Title)
	if len(p.Categories) != 0 {
		res = append(res, " | 🔖 "+strings.Join(p.Categories, ", "))
	}
	if p.Views > 0 {
		res = append(res, " | 👀 "+strconv.Itoa(p.Views))
	}
	return []string{strings.Join(res, "")}, nil
}

func (pm *porn) pornHubComment(arguments core.CommandArguments) ([]string, error) {
	if pm.isSpam() {
		return []string{"Stop flooding"}, nil
	}

	var gay = len(arguments.Arguments) > 0 && strings.ToUpper(arguments.Arguments[0]) == "GAY"
	var p *pornhub.Pornhub
	var withArgs = len(arguments.Arguments) != 0
	var results []pornhub.Result
	var err error

	if withArgs {
		var query string
		if !gay {
			query = strings.Join(arguments.Arguments, " ")
		} else {
			query = strings.Join(arguments.Arguments[1:], " ")
		}
		results, err = pornhub.Search(query, gay)

		pm.Log.Info("User " + arguments.Nick + " searched " + query)

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
			p, err = pornhub.InitRandom(gay)
		}

		if err != nil {
			return nil, err
		}
		if len(p.Comments) != 0 {
			break
		}

		time.Sleep(200 * time.Millisecond)
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
	res = append(res, ": "+com.Message)
	if com.Score != 0 {
		res = append(res, "	(Score:"+strconv.Itoa(com.Score)+")")
	}

	return []string{strings.Join(res, "")}, nil
}

func (pm *porn) pornHubLast(arguments core.CommandArguments) (result []string, err error) {
	return []string{pm.lastURL + " - " + pm.lastTitle}, err
}

// Gay

func (pm *porn) gayPornHubTitleShort(arguments core.CommandArguments) ([]string, error) {
	args := arguments
	args.Arguments = append([]string{"GAY"}, args.Arguments...)
	return pm.pornHubTitleShort(args)
}

func (pm *porn) gayPornHubTitleFull(arguments core.CommandArguments) ([]string, error) {
	args := arguments
	args.Arguments = append([]string{"GAY"}, args.Arguments...)
	return pm.pornHubTitleFull(args)
}

func (pm *porn) gayPornHubComment(arguments core.CommandArguments) ([]string, error) {
	args := arguments
	args.Arguments = append([]string{"GAY"}, args.Arguments...)
	return pm.pornHubComment(args)
}
