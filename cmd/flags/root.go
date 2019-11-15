//Copyright 2019 Expedia, Inc.
//
//Licensed under the Apache License, Version 2.0 (the "License");
//you may not use this file except in compliance with the License.
//You may obtain a copy of the License at
//
//http://www.apache.org/licenses/LICENSE-2.0
//
//Unless required by applicable law or agreed to in writing, software
//distributed under the License is distributed on an "AS IS" BASIS,
//WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//See the License for the specific language governing permissions and
//limitations under the License.

package flags

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"mittens/pkg/grpc"
	"mittens/pkg/http"
	"mittens/pkg/warmup"
	"time"
)

type Root struct {
	TimeoutSeconds           int
	Concurrency              int
	RequestDelayMilliseconds int
	ExitAfterWarmup          bool
	FileProbe                *FileProbe
	ServerProbe              *ServerProbe
	Target                   *Target
	Http                     *Http
	Grpc                     *Grpc
	Profile                  *Profile
}

func NewRoot() *Root {

	return &Root{
		FileProbe:   new(FileProbe),
		ServerProbe: new(ServerProbe),
		Target:      new(Target),
		Http:        new(Http),
		Grpc:        new(Grpc),
		Profile:     new(Profile),
	}
}

func (r *Root) String() string {
	return fmt.Sprintf("%+v", *r)
}

func (r *Root) InitFlags(cmd *cobra.Command) {

	cmd.Flags().IntVar(
		&r.TimeoutSeconds,
		"timeout-seconds",
		60,
		"Time after which warm up will stop making requests",
	)
	cmd.Flags().IntVar(
		&r.Concurrency,
		"concurrency",
		2,
		"Number of concurrent requests for warm up",
	)
	cmd.Flags().IntVar(
		&r.RequestDelayMilliseconds,
		"request-delay-milliseconds",
		50,
		"Delay in milliseconds between requests",
	)
	cmd.Flags().BoolVar(
		&r.ExitAfterWarmup,
		"exit-after-warmup",
		false,
		"If warm up process should finish after completion. This is useful to prevent container restarts.",
	)

	r.FileProbe.InitFlags(cmd)
	r.ServerProbe.InitFlags(cmd)
	r.Target.InitFlags(cmd)
	r.Http.InitFlags(cmd)
	r.Grpc.InitFlags(cmd)
	r.Profile.InitFlags(cmd)
}

func (r *Root) GetWarmupOptions() warmup.Options {

	return warmup.Options{
		TimeoutSeconds: r.TimeoutSeconds,
		Concurrency:    r.Concurrency,
	}
}

func (r *Root) GetReadinessHttpClient() http.Client {
	return r.Target.GetReadinessHttpClient()
}

func (r *Root) GetHttpClient() http.Client {
	return r.Target.GetHttpClient()
}

func (r *Root) GetGrpcClient() grpc.Client {
	return r.Target.GetGrpcClient(r.TimeoutSeconds)
}

func (r *Root) GetWarmupTargetOptions() warmup.TargetOptions {

	options := r.Target.GetWarmupTargetOptions()
	if options.ReadinessTimeoutInSeconds <= 0 {
		log.Printf("readiness timeout in seconds not set, defaulting to timeout in seconds: %ds", r.TimeoutSeconds)
		options.ReadinessTimeoutInSeconds = r.TimeoutSeconds
	}
	return options
}

func (r *Root) GetWarmupHttpHeaders() map[string]string {
	return r.Http.GetWarmupHttpHeaders()
}

func (r *Root) GetWarmupHttpRequests(done <-chan struct{}) (chan warmup.HttpRequest, error) {

	requests, err := r.Http.GetWarmupHttpRequests()
	if err != nil {
		return nil, err
	}

	requestsChan := make(chan warmup.HttpRequest, r.Concurrency)
	go func() {
		if len(requests) == 0 {
			log.Print("no http warm up requests specified")
			close(requestsChan)
			return
		}
		timeout := time.After(time.Duration(r.TimeoutSeconds) * time.Second)
		for {
			for _, v := range requests {
				time.Sleep(time.Duration(r.RequestDelayMilliseconds) * time.Millisecond)
				select {
				case <-timeout:
					log.Printf("timeout %d seconds exceeded", r.TimeoutSeconds)
					close(requestsChan)
					return
				case <-done:
					log.Print("get http requests: received done signal, closing http chan")
					close(requestsChan)
					return
				default:
					requestsChan <- v
				}
			}
		}
	}()
	return requestsChan, nil
}

func (r *Root) GetWarmupGrpcHeaders() []string {
	return r.Grpc.GetWarmupGrpcHeaders()
}

func (r *Root) GetWarmupGrpcRequests(done <-chan struct{}) (chan warmup.GrpcRequest, error) {

	requests, err := r.Grpc.GetWarmupGrpcRequests()
	if err != nil {
		return nil, err
	}

	requestsChan := make(chan warmup.GrpcRequest, r.Concurrency)
	go func() {
		if len(requests) == 0 {
			log.Print("no grpc warm up requests specified")
			close(requestsChan)
			return
		}
		timeout := time.After(time.Duration(r.TimeoutSeconds) * time.Second)
		for {
			for _, v := range requests {
				time.Sleep(time.Duration(r.RequestDelayMilliseconds) * time.Millisecond)
				select {
				case <-timeout:
					log.Printf("timeout %d seconds exceeded", r.TimeoutSeconds)
					close(requestsChan)
					return
				case <-done:
					log.Print("get grpc requests: received done signal, closing grpc chan")
					close(requestsChan)
					return
				default:
					requestsChan <- v
				}
			}
		}
	}()
	return requestsChan, nil
}
