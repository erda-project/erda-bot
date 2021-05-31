package conf

import (
	"github.com/erda-project/erda/pkg/envconf"
)

type Conf struct {
	GitHubActor string `env:"GITHUB_ACTOR" default:"erda-bot"`
	GitHubEmail string `env:"GITHUB_EMAIL" default:"erda@terminus.io"`
	GitHubToken string `env:"GITHUB_TOKEN" required:"true"`
}

var cfg Conf
var bot GitHubBot

func Load() {
	envconf.MustLoad(&cfg)
	bot = GitHubBot{
		GitHubActor: cfg.GitHubActor,
		GitHubEmail: cfg.GitHubEmail,
		GitHubToken: cfg.GitHubToken,
	}
}

func Bot() GitHubBot {
	return bot
}
