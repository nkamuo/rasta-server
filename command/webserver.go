package command

import "github.com/nkamuo/rasta-server/web"

func StartWebServer() (err error) {
	r, err := web.BuildWebServer()
	if err != nil {
		return err
	}
	err = r.Run(":8090")
	return
}
