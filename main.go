package main

import (
	"encoding/json"
	"fmt"
	"github.com/affanshahid/oversight/prober"
	"github.com/affanshahid/oversight/prober/httpprobe"
	"github.com/affanshahid/oversight/prober/probe"

	"github.com/jinzhu/gorm"
	"github.com/jinzhu/gorm/dialects/postgres"
	"github.com/spf13/viper"
)

func main() {
	db, err := gorm.Open(
		"postgres",
		fmt.Sprintf(
			"postgres://%s:%s@%s:%s/%s?sslmode=disable",
			viper.GetString("db.username"),
			viper.GetString("db.password"),
			viper.GetString("db.host"),
			viper.GetString("db.port"),
			viper.GetString("db.database"),
		),
	)

	if err != nil {
		panic(err)
	}

	defer db.Close()
	db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)
	db.AutoMigrate(new(probe.Config))

	options := httpprobe.HTTPOptions{
		Method: "GET",
		URL:    "http://www.google.com",
	}

	rawOpts, err := json.Marshal(&options)

	if err != nil {
		panic(err)
	}

	pc := probe.Config{
		Descriminator: "http",
		Interval:      2000,
		Options:       postgres.Jsonb{RawMessage: rawOpts},
	}

	db.Where("id is not null").Unscoped().Delete(new(probe.Config))
	db.Create(&pc)

	options2 := httpprobe.HTTPOptions{
		Method: "GET",
		URL:    "http://www.example.com",
	}

	rawOpts2, err := json.Marshal(&options2)

	if err != nil {
		panic(err)
	}

	pc2 := probe.Config{
		Descriminator: "http",
		Interval:      2000,
		Options:       postgres.Jsonb{RawMessage: rawOpts2},
	}

	db.Create(&pc2)

	prober := prober.New(db)
	prober.Start()

	select {}
}
