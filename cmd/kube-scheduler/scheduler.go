/*
Copyright 2014 The Kubernetes Authors.

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

package main

import (
	"fmt"
	"math/rand"
	"os"
	"path"
	"time"

	"github.com/spf13/pflag"

	"io/ioutil"

	"github.com/Microsoft/KubeDevice/device-scheduler/device"
	"github.com/Microsoft/KubeDevice/kube-scheduler/cmd/app"
	"github.com/Microsoft/KubeDevice/logger"

	cliflag "k8s.io/component-base/cli/flag"
	"k8s.io/component-base/logs"
	"k8s.io/klog"
	_ "k8s.io/kubernetes/pkg/util/prometheusclientgo" // load all the prometheus client-go plugins
	_ "k8s.io/kubernetes/pkg/version/prometheus"      // for version metric registration
)

func init() {
	logger.SetLogger()
}

func main() {
	rand.Seed(time.Now().UnixNano())

	command := app.NewSchedulerCommand()

	// TODO: once we switch everything over to Cobra commands, we can go back to calling
	// utilflag.InitFlags() (by removing its pflag.Parse() call). For now, we have to set the
	// normalize func and add the go flag set by hand.
	pflag.CommandLine.SetNormalizeFunc(cliflag.WordSepNormalizeFunc)
	// utilflag.InitFlags()
	logs.InitLogs()
	defer logs.FlushLogs()

	// add the device schedulers
	var deviceSchedulerPlugins []string
	pluginPath := "/schedulerplugins"
	devPlugins, err := ioutil.ReadDir(pluginPath)
	if err != nil {
		klog.Errorf("Cannot read plugins - skipping")
	}
	for _, pluginFile := range devPlugins {
		deviceSchedulerPlugins = append(deviceSchedulerPlugins, path.Join(pluginPath, pluginFile.Name()))
	}
	device.DeviceScheduler.AddDevicesSchedulerFromPlugins(deviceSchedulerPlugins)

	if err := command.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
