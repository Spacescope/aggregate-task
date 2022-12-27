package utils

import "context"

func BatchWrite(ctx context.Context, obj interface{}, height int64, result interface{}) error {
	session := EngineGroup[DBOBTask].NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	_, err := session.Where("height = ?", height).Delete(obj)
	if err != nil {
		return err
	}

	_, err = session.Insert(result)
	if err != nil {
		return err
	}

	session.Commit()

	return nil
}

func SingleWrite(ctx context.Context, obj interface{}, height int64, result interface{}) error {
	session := EngineGroup[DBOBTask].NewSession()
	defer session.Close()

	if err := session.Begin(); err != nil {
		return err
	}

	_, err := session.Where("height = ?", height).Delete(obj)
	if err != nil {
		return err
	}

	_, err = session.InsertOne(result)
	if err != nil {
		return err
	}

	session.Commit()

	return nil
}
