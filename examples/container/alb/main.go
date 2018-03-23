package main

import (
	"log"

	"github.com/IBM-Bluemix/bluemix-go/session"

	v1 "github.com/IBM-Bluemix/bluemix-go/api/container/containerv1"
	"github.com/IBM-Bluemix/bluemix-go/trace"
)

var albConfig = v1.ALBConfig{
	AlbType:   "public",
	Zone:      "ams03",
	ClusterID: "test4",
}

func main() {
	trace.Logger = trace.NewLogger("true")

	sess, err := session.New()
	if err != nil {
		log.Fatal(err)
	}

	albClient, err := v1.New(sess)
	if err != nil {
		log.Fatal(err)
	}
	albAPI := albClient.Albs()

	err = albAPI.DeployALB("test4", albConfig)
	log.Fatal(err)
}
