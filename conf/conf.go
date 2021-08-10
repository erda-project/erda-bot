package conf

import (
	"github.com/erda-project/erda/pkg/envconf"
)

type Conf struct {
	GitHubActor string `env:"GITHUB_ACTOR" default:"erda-bot"`
	GitHubEmail string `env:"GITHUB_EMAIL" default:"erda@terminus.io"`
	GitHubToken string `env:"GITHUB_TOKEN" required:"true"`

	DingTalkAccessToken string `env:"DINGTALK_ACCESS_TOKEN"`
	DingTalkSecret      string `env:"DINGTALK_SECRET"`
}

var cfg Conf
var bot GitHubBot
var dingtalk DingTalkConf

func Load() {
	envconf.MustLoad(&cfg)
	bot = GitHubBot{
		GitHubActor: cfg.GitHubActor,
		GitHubEmail: cfg.GitHubEmail,
		GitHubToken: cfg.GitHubToken,
	}
	dingtalk = DingTalkConf{
		AccessToken: cfg.DingTalkAccessToken,
		Secret:      cfg.DingTalkSecret,
	}
}

func Bot() GitHubBot {
	return bot
}

func DingTalk() DingTalkConf {
	return dingtalk
}
