package topology

import (
	"testing"
	"time"

	"utils/mongodb/mongo-go-driver/core/connstring"
	"utils/stretchr/testify/assert"
)

func TestOptionsSetting(t *testing.T) {
	conf := &config{}
	ssts := time.Minute
	assert.Zero(t, conf.cs)

	opt := WithConnString(func(connstring.ConnString) connstring.ConnString {
		return connstring.ConnString{
			ServerSelectionTimeout:    ssts,
			ServerSelectionTimeoutSet: true,
		}

	})

	assert.NoError(t, opt(conf))

	assert.Equal(t, ssts, conf.serverSelectionTimeout)
}