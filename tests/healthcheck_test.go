/*
Copyright 2017 Crunchy Data Solutions, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package tests

import (
	"bytes"
	"log"
	"net/http"
	"os/exec"
	"testing"
	"time"
)

func TestHealthcheck(t *testing.T) {
	const SLEEP_TIME = 6
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("TestHealthcheck was called")
	var startTime = time.Now()
	conn, err := Connect()
	defer conn.Close()
	if err != nil {
		t.FailNow()
	}

	//read the proxy event stream which equates to a healthcheck
	//has just been performed...after which we will have a retry test
	//window
	_, err = http.Get("http://localhost:10000/api/stream")
	if err != nil {
		t.FailNow()
	}
	log.Println("after the GET")

	//shut down the replica
	var cmdStdout, cmdStderr bytes.Buffer
	var cmd *exec.Cmd
	cmd = exec.Command("docker", "stop", "replica")
	cmd.Stdout = &cmdStdout
	cmd.Stderr = &cmdStderr
	err = cmd.Run()
	if err != nil {
		log.Println("docker stop stdout=" + cmdStdout.String())
		log.Println("docker stop stderr=" + cmdStderr.String())
		t.FailNow()
	}
	log.Println("docker stop stdout=" + cmdStdout.String())

	//sleep a bit to give the replica time to stop
	log.Println("sleeping to let replica shutdown")
	time.Sleep(time.Duration(SLEEP_TIME) * time.Second)

	//read the proxy event stream which should test that healthcheck
	//still works after a backend has been shutdown
	_, err = http.Get("http://localhost:10000/api/stream")
	if err != nil {
		t.FailNow()
	}
	log.Println("after the 2nd GET")

	//restart the replica after the test
	cmd = exec.Command("docker", "start", "replica")
	err = cmd.Run()
	if err != nil {
		log.Println("docker start stdout=" + cmdStdout.String())
		log.Println("docker start stderr=" + cmdStderr.String())
		t.FailNow()
	}
	log.Println("docker start stdout=" + cmdStdout.String())
	//sleep a bit to give the replica time to restart
	log.Println("sleeping to let replica restart")
	time.Sleep(time.Duration(SLEEP_TIME) * time.Second)

	log.Printf("Duration %s\n", time.Since(startTime))

}
