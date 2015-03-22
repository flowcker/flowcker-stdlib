/*  Copyright (C) 2015 Leopoldo Lara Vazquez.

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package stdlib

import (
	"encoding/json"
	"time"

	"github.com/davecgh/go-spew/spew"
	fc "github.com/flowcker/flowcker/common"
	"github.com/garyburd/redigo/redis"
)

type updateRedisKeyConfig struct {
	Addr     string
	DB       string
	Key      string
	Debounce uint32
}

// UpdateRedisKey for each incoming IP in port `input`, saves the IP content in
// the Redis DB.
// Configuration:
//  - addr string // address of Redis server, IP:Port
//  - db string // DB number to select
//  - key string // redis key to set to the IP content
func UpdateRedisKey(atom *fc.Atom, in chan fc.IPInbound) (out chan fc.IPOutbound, err error) {
	var log = setupLogging(atom.ID)

	out = make(chan fc.IPOutbound)
	var config updateRedisKeyConfig

	config.Key = "output"
	config.DB = "0"
	config.Addr = "127.0.0.1:49153"
	config.Debounce = 500
	json.Unmarshal(*atom.Config, &config)

	go func() {
		defer log.Debug("UpdateRedisKey element: exiting")
		defer close(out)

		log.Debug("Starting UpdateRedisKey element")

		log.Debug("UpdateRedisKey: connecting to " + config.Addr + " DB " + config.DB)
		conn, err := redis.DialTimeout("tcp", config.Addr, time.Second, time.Second, time.Second)
		if err != nil {
			log.Fatal("UpdateRedisKey: error connecting to Redis")
			panic(err)
		}
		log.Debug("UpdateRedisKey connected to Redis")

		log.Debug("UpdateRedisKey selecting DB")
		_, err = conn.Do("SELECT", config.DB)
		if err != nil {
			log.Fatal("UpdateRedisKey: error selecting DB")
			panic(err)
		}
		log.Debug("UpdateRedisKey selected DB")

		// TODO move this debouncedSet stuff to a better place, or extract it
		debouncedSet := make(chan []byte)
		if config.Debounce != 0 {
			go func() {
				var buffer []byte
				var ok bool

				buffer, ok = <-debouncedSet
				if !ok {
					return
				}
				for {
					select {
					case buffer, ok = <-debouncedSet:
						if !ok {
							return
						}
					case <-time.After(time.Duration(config.Debounce) * time.Millisecond):
						_, err := conn.Do("SET", config.Key, buffer)
						if err != nil {
							log.Fatal("UpdateRedisKey: error seting key value")
							panic(err)
						}
						log.Debug("UpdateRedisKey element: data saved \n%s", spew.Sdump(buffer))
						buffer, ok = <-debouncedSet
						if !ok {
							return
						}
					}
				}
			}()
		} else {
			go func() {
				var buffer []byte
				var ok bool

				for {
					buffer, ok = <-debouncedSet
					if !ok {
						return
					}
					_, err := conn.Do("SET", config.Key, buffer)
					if err != nil {
						log.Fatal("UpdateRedisKey: error seting key value")
						panic(err)
					}
				}
			}()
		}

		for incoming := range in {
			switch incoming.GetTo().Name {
			case "input":
				if len(incoming.GetData()) == 0 {
					log.Debug("UpdateRedisKey element: received empty IP")
					continue
				}
				log.Debug("UpdateRedisKey element: received data")
				debouncedSet <- incoming.GetData()
			}
		}
	}()

	return out, nil
}
