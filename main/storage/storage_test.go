package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStorage(t *testing.T) {
	mysql := &MySQL{Name: "mysql"}
	mongodb := &MongoDB{Name: "mongodb"}

	mysqlAdapter := &MySQLAdapter{mysql: mysql}
	mongoAdapter := &MongoDBAdapter{mongodb: mongodb}

	assert.Equal(t, "MySQL", mysqlAdapter.AdapterName())
	assert.Equal(t, "MongoDB", mongoAdapter.AdapterName())

}
