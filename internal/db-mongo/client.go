package db_mongo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	MODE_STANDALONE  = "standalone"
	MODE_REPLICATION = "replication"

	AUTH_BASIC = "basic"
)

type Client interface {
	Connect(ctx context.Context) error
	Ping(ctx context.Context, rp *readpref.ReadPref) error
	Database(name string, opts ...*options.DatabaseOptions) *mongo.Database
}

type ClientOption = func(*ClientParam)

type ClientParam struct {
	DbMode string
	DbName string

	AuthMode     string
	AuthSource   string
	AuthUser     string
	AuthPassword string

	StdHost string
	StdPort int

	RsName  string
	RsHosts []string
}

type ClientLocation struct {
	Host string
	Port int
}

type ClientConfig struct {
	DbName   string
	DbMode   string
	AuthMode string
}

func (p *ClientParam) ModeStandalone() bool {
	return p.DbMode == MODE_STANDALONE
}

func (p *ClientParam) ModeReplication() bool {
	return p.DbMode == MODE_REPLICATION
}

func (p *ClientParam) ModeSupported() bool {
	return p.ModeStandalone() || p.ModeReplication()
}

func (p *ClientParam) AuthBasic() bool {
	return p.AuthMode == AUTH_BASIC
}

func (p *ClientParam) AuthSupported() bool {
	return p.AuthBasic()
}

func UsingStandalone(host string, port int) ClientOption {
	return func(cp *ClientParam) {
		cp.DbMode = MODE_STANDALONE
		cp.StdHost = host
		cp.StdPort = port
	}
}

func UsingReplication(rsName string, rsHosts []string) ClientOption {
	return func(cp *ClientParam) {
		cp.DbMode = MODE_REPLICATION
		cp.RsName = rsName
		cp.RsHosts = rsHosts
	}
}

func WithBasicAuth(username, password, source string) ClientOption {
	return func(cp *ClientParam) {
		cp.AuthUser = username
		cp.AuthPassword = password
		cp.AuthSource = source
		cp.AuthMode = AUTH_BASIC
	}
}

func WithConfig(cfg ClientConfig) ClientOption {
	return func(cp *ClientParam) {
		if cfg.DbName != "" {
			cp.DbName = cfg.DbName
		}
		if cfg.AuthMode != "" {
			cp.AuthMode = cfg.AuthMode
		}
		if cfg.DbMode != "" {
			cp.DbMode = cfg.DbMode
		}
	}
}

func NewClient(opts ...ClientOption) (*mongo.Client, error) {
	p := ClientParam{}
	for _, opt := range opts {
		opt(&p)
	}

	if !p.ModeSupported() {
		return nil, fmt.Errorf("mode is not supported")
	}
	if !p.AuthSupported() {
		return nil, fmt.Errorf("auth is not supported")
	}

	mongoOption := options.Client()

	if p.ModeStandalone() {
		mongoOption.SetHosts([]string{fmt.Sprintf("%s:%d", p.StdHost, p.StdPort)})
	} else if p.ModeReplication() {
		mongoOption.
			SetHosts(p.RsHosts).
			SetReplicaSet(p.RsName).
			SetReadPreference(readpref.Secondary())
	}

	if p.AuthBasic() {
		mongoOption.SetAuth(options.Credential{
			Username:   p.AuthUser,
			Password:   p.AuthPassword,
			AuthSource: p.AuthSource,
		})
	}

	return mongo.NewClient(mongoOption)
}
