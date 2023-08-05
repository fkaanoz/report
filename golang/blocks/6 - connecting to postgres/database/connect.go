package database

import (
	"database/sql"
	"net/url"
)

type DBConfig struct {
	Host        string
	Port        string
	Username    string
	Password    string
	Database    string
	SSLMode     string
	MaxIdleConn int
	MaxOpenConn int
}

func Open(conf DBConfig) (*sql.DB, error) {
	q := url.Values{}
	q.Set("timezone", "utc")
	q.Set("sslmode", conf.SSLMode)

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(conf.Username, conf.Password),
		Host:     conf.Host,
		Path:     conf.Database,
		RawQuery: q.Encode(),
	}

	conn, err := sql.Open("postgres", u.String())
	if err != nil {
		return nil, err
	}

	conn.SetMaxIdleConns(conf.MaxIdleConn)
	conn.SetMaxOpenConns(conf.MaxOpenConn)

	return conn, nil
}
