package olxclient

import (
	"gopkg.in/mgo.v2"
)

var (
	IsDrop = true
)

func Run() error {
	session, err := mgo.Dial("0.0.0.0")
	if err != nil {
		return err
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	// Drop Database
	if IsDrop {
		err = session.DB("olxparser").DropDatabase()
		if err != nil {
			return err
		}
	}

	c := session.DB("olxparser").C("adverts")
	index := mgo.Index{
		Key:        []string{"id"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err = c.EnsureIndex(index)
	if err != nil {
		return err
	}

	adverts := parse()
	// required for mgo
	s := make([]interface{}, len(adverts))
	for i, v := range adverts {
		s[i] = v
	}
	// ===

	err = c.Insert(s...)
	if err != nil {
		return err
	}

	return nil
}
