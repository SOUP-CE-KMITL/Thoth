// TODO: need to change
package main

import (
	"github.com/influxdb/influxdb/client/v2"
)

const (
	MyDB = "thoth"
	username = "thoth_user"
	password = "thoth_secret"
)

func runProfile() {
	
}

func queryDB(clnt client.Client, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: MyDB,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return res, nil
}
