package sdk

import "github.com/it-chain/leveldb-wrapper"

type Cell struct {
	DBHandler *leveldbwrapper.DBHandle
}

func NewCell(name string) *Cell{
	path := "./wsdb"
	dbProvider := leveldbwrapper.CreateNewDBProvider(path)
	return &Cell{
		DBHandler: dbProvider.GetDBHandle(name),
	}
}

func (c Cell) PutData(key string, value []byte) error {
	return c.DBHandler.Put([]byte(key), value, true)
}

func (c Cell) GetData(key string) ([]byte, error) {
	value, err := c.DBHandler.Get([]byte(key))
	return value, err
}