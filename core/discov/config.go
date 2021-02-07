package discov

import "errors"

type EtcdConf struct {
	Hosts []string
	Key   string
	User  string `json:",optional"`
	Pass  string `json:",optional"`
}

func (c EtcdConf) Validate() error {
	if len(c.Hosts) == 0 {
		return errors.New("empty etcd hosts")
	} else if len(c.Key) == 0 {
		return errors.New("empty etcd key")
	} else if len(c.User) > 0 || len(c.Pass) > 0 {
		if len(c.User) == 0 || len(c.Pass) == 0 {
			return errors.New("etcd user and pass must pairs")
		}
		return nil
	} else {
		return nil
	}
}

func (c EtcdConf) EnableAuth() bool {
	return len(c.User) > 0 && len(c.Pass) > 0
}
