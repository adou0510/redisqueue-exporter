package main

import (
	"fmt"
	"strconv"
	"strings"
)

type RedisDSN struct {
	Addr     string
	Password string
	DB       int
	Key      string
}

// ParseRedisDSN 解析 Redis DSN
func ParseRedisDSN(dsn string) (*RedisDSN, error) {
	// PASSWORD@ADDR/DB/KEY
	password := ""
	var addressParams []string

	params := strings.Split(dsn, "@")

	if len(params) == 2 {
		// 有密码
		password = params[0]
		addressParams = strings.Split(params[1], "/")
	} else {
		addressParams = strings.Split(params[0], "/")
	}

	if len(addressParams) != 3 {
		return nil, fmt.Errorf("invalid dsn")
	}

	addr := addressParams[0]
	db, err := strconv.Atoi(addressParams[1])
	if err != nil {
		return nil, fmt.Errorf("invalid db")
	}

	key := addressParams[2]

	return &RedisDSN{
		Addr:     addr,
		Password: password,
		DB:       db,
		Key:      key,
	}, nil

}
