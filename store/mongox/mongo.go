package mongox

import "github.com/qiniu/qmgo"

type Client struct {
	*qmgo.Client
	Config *qmgo.Config
}

type Cli struct {
	*qmgo.QmgoClient
	Config *qmgo.Config
}
