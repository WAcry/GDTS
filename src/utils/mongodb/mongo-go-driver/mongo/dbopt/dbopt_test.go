package dbopt

import (
	"testing"

	"utils/mongodb/mongo-go-driver/core/readconcern"
	"utils/mongodb/mongo-go-driver/core/readpref"
	"utils/mongodb/mongo-go-driver/core/writeconcern"
)

var rcLocal = readconcern.Local()
var rcMajority = readconcern.Majority()

var wc1 = writeconcern.New(writeconcern.W(10))
var wc2 = writeconcern.New(writeconcern.W(20))

var rpPrimary = readpref.Primary()
var rpSeconadary = readpref.Secondary()

func requireDbEqual(t *testing.T, expected *Database, actual *Database) {
	switch {
	case expected.ReadConcern != actual.ReadConcern:
		t.Errorf("read concerns don't match")
	case expected.WriteConcern != actual.WriteConcern:
		t.Errorf("write concerns don't match")
	case expected.ReadPreference != actual.ReadPreference:
		t.Errorf("read preferences don't match")
	}
}

func createNestedBundle1(t *testing.T) *DatabaseBundle {
	nested := BundleDatabase(ReadConcern(rcMajority))
	testhelpers.RequireNotNil(t, nested, "nested bundle was nil")

	outer := BundleDatabase(ReadConcern(rcMajority), WriteConcern(wc1), nested)
	testhelpers.RequireNotNil(t, nested, "nested bundle was nil")

	return outer
}

func createdNestedBundle2(t *testing.T) *DatabaseBundle {
	b1 := BundleDatabase(WriteConcern(wc2))
	testhelpers.RequireNotNil(t, b1, "b1 was nil")

	b2 := BundleDatabase(ReadPreference(rpPrimary), b1)
	testhelpers.RequireNotNil(t, b2, "b2 was nil")

	outer := BundleDatabase(WriteConcern(wc1), ReadPreference(rpSeconadary), b2)
	testhelpers.RequireNotNil(t, outer, "outer was nil")

	return outer
}

func createNestedBundle3(t *testing.T) *DatabaseBundle {
	b1 := BundleDatabase(ReadConcern(rcMajority))
	testhelpers.RequireNotNil(t, b1, "b1 was nil")

	b2 := BundleDatabase(WriteConcern(wc1), b1)
	testhelpers.RequireNotNil(t, b2, "b1 was nil")

	b3 := BundleDatabase(ReadPreference(rpPrimary))
	testhelpers.RequireNotNil(t, b3, "b1 was nil")

	b4 := BundleDatabase(WriteConcern(wc2), b3)
	testhelpers.RequireNotNil(t, b4, "b1 was nil")

	outer := BundleDatabase(b4, WriteConcern(wc1), b2)
	testhelpers.RequireNotNil(t, outer, "b1 was nil")

	return outer
}

func TestDbOpt(t *testing.T) {
	nilBundle := BundleDatabase()
	var nilDb = &Database{}

	var bundle1 *DatabaseBundle
	bundle1 = bundle1.ReadConcern(rcLocal).ReadConcern(rcLocal)
	testhelpers.RequireNotNil(t, bundle1, "created bundle was nil")
	bundle1Db := &Database{
		ReadConcern: rcLocal,
	}

	bundle2 := BundleDatabase(ReadConcern(rcLocal))
	bundle2Db := &Database{
		ReadConcern: rcLocal,
	}

	bundle3 := BundleDatabase(ReadConcern(rcLocal), ReadPreference(rpPrimary), ReadConcern(rcMajority), ReadPreference(rpSeconadary),
		WriteConcern(wc1), WriteConcern(wc2))
	bundle3Db := &Database{
		ReadConcern:    rcMajority,
		ReadPreference: rpSeconadary,
		WriteConcern:   wc2,
	}

	nested1 := createNestedBundle1(t)
	nested1Db := &Database{
		ReadConcern:  rcMajority,
		WriteConcern: wc1,
	}

	nested2 := createdNestedBundle2(t)
	nested2Db := &Database{
		ReadPreference: rpPrimary,
		WriteConcern:   wc2,
	}

	nested3 := createNestedBundle3(t)
	nested3Db := &Database{
		ReadConcern:    rcMajority,
		ReadPreference: rpPrimary,
		WriteConcern:   wc1,
	}

	t.Run("TestAll", func(t *testing.T) {
		opts := []Option{
			ReadConcern(rcLocal),
			WriteConcern(wc1),
			ReadPreference(rpPrimary),
		}

		db, err := BundleDatabase(opts...).Unbundle()
		testhelpers.RequireNil(t, err, "got non-nil error from unbundle: %s", err)
		requireDbEqual(t, db, &Database{
			ReadConcern:    rcLocal,
			WriteConcern:   wc1,
			ReadPreference: rpPrimary,
		})
	})

	t.Run("Unbundle", func(t *testing.T) {
		var cases = []struct {
			name   string
			bundle *DatabaseBundle
			db     *Database
		}{
			{"NilBundle", nilBundle, nilDb},
			{"Bundle1", bundle1, bundle1Db},
			{"Bundle2", bundle2, bundle2Db},
			{"Bundle3", bundle3, bundle3Db},
			{"Nested1", nested1, nested1Db},
			{"Nested2", nested2, nested2Db},
			{"Nested3", nested3, nested3Db},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				db, err := tc.bundle.Unbundle()
				testhelpers.RequireNil(t, err, "err unbundling db: %s", err)
				requireDbEqual(t, db, tc.db)
			})
		}
	})
}
