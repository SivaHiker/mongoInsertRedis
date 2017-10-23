package main

import(
	"github.com/go-redis/redis"

	"sync"
	"encoding/json"
	"gopkg.in/mgo.v2"
	"fmt"
)

func main() {

	var counter int
	session, err := mgo.Dial("10.15.0.149")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("userlist").C("httpuserdata")
	iter := GetRedisInstanceGCP().SScan("new_user_info_prefix_three", 0, "*", 10000).Iterator()
	for iter.Next() {
		var userhttpdata UserHTTPData
		err := json.Unmarshal([]byte(iter.Val()), &userhttpdata)
		userinsert := UserInfo{HttpUserData: userhttpdata, HttpFlag: false}
		if err == nil {
			c.Insert(userinsert)
			counter++
		}
		//GetRedisInstanceGCP().SAdd("ums_user_key", iter.Val()).Result()
	}
	if err := iter.Err(); err != nil {
		panic(err)
	}

	fmt.Println("Total Inserted Document Count from Redis",counter)

}

func GetRedisInstanceGCP() *redis.Client {
	var onceGCP sync.Once
	var instanceGCP *redis.Client
	onceGCP.Do(func() {
		client := redis.NewClient(&redis.Options{
			Addr:     "redis-13511.c2.asia-southeast-1-1.gce.cloud.redislabs.com:13511",
			//Addr:     "localhost:6379",
			Password: "",
			DB:       0,
			PoolSize: 100,
			/*
					PoolTimeout:  10 * time.Minute,
					IdleTimeout:  5 * time.Minute,
					ReadTimeout:  2 * time.Second,
					WriteTimeout: 10 * time.Second,
			*/
		})
		instanceGCP = client
	})
	return instanceGCP
}

type UserInfo struct {
	HttpUserData UserHTTPData `json:"UserData"`
	HttpFlag bool `json:"flag"`
}


type UserHTTPData struct {
	EncryptedToken string `json:"encrypted_token"`
	Msisdn         string `json:"msisdn"`
	PlatformToken  string `json:"platform_token"`
	PlatformUID    string `json:"platform_uid"`
	PubKey         string `json:"pub_key"`
	RsaKey         string `json:"rsa_key"`
	Token          string `json:"token"`
	UID            string `json:"uid"`
	UUID           string `json:"uuid"`
}
