package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/values"
	"helm.sh/helm/v3/pkg/downloader"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/release"
)

func findChartDir(root string) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {

		if err != nil {

			fmt.Println(err)
			return nil
		}
		fmt.Println(info.Name())
		if !info.IsDir() && strings.ToLower(info.Name()) == "values.yaml" {
			files = append(files, path)
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	// 別のフォルダに同じ名前のファイルが見つかったらどうするか
	for _, file := range files {
		fmt.Printf("chart dir %s\n", filepath.Dir(file))
	}
}

func main() {
	settings := cli.New()
	settings.Debug = false

	actionConfig := new(action.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	client := action.NewInstall(actionConfig)

	var args []string

	var vals []string
	vals = append(vals, "favorite.drink=tea")
	vals = append(vals, "hoge.myServiceName=my-service-1")
	vals = append(vals, "hoge.myServicePort=7654")
	var valueFiles []string
	valueFiles = append(valueFiles, "F:/dev/helm-v3.9.0/mychart/values.yaml")

	// 存在しないと後続でnil pointerになる
	// 大文字小文字は区別される
	valueOpts := &values.Options{Values: vals, ValueFiles: valueFiles}

	var out bytes.Buffer

	client.DryRun = true
	//client.ReleaseName = "my-release-name-1"
	client.Replace = true // Skip the name check
	client.ClientOnly = true
	//client.APIVersions = chartutil.VersionSet(extraAPIs)
	client.IncludeCRDs = false

	args = append(args, "my-release-1")
	args = append(args, "F:/dev/helm-v3.9.0/mychart")

	findChartDir("F:/dev/helm-v3.9.0")

	rel, err := runInstall(args, client, valueOpts, &out, settings)

	if err != nil && !settings.Debug {
		if rel != nil {
			fmt.Printf("%v\n\nUse --debug flag to render out invalid YAML", err)
			return
		}
		fmt.Printf("[ERROR] %v", err)
		return
	}

	// テンプレートファイルが複数ある場合、---で区切られてくる
	fmt.Printf("%s", rel.Manifest)
}

func runInstall(args []string, client *action.Install, valueOpts *values.Options, out io.Writer, settings *cli.EnvSettings) (*release.Release, error) {
	fmt.Printf("[DEBUG] Original chart version: %q\n", client.Version)
	if client.Version == "" && client.Devel {
		fmt.Printf("setting version to >0.0.0-0")
		client.Version = ">0.0.0-0"
	}

	name, chart, err := client.NameAndChart(args)
	if err != nil {
		return nil, err
	}
	client.ReleaseName = name

	cp, err := client.ChartPathOptions.LocateChart(chart, settings)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[DEBUG] CHART PATH: %s\n", cp)

	p := getter.All(settings)
	vals, err := valueOpts.MergeValues(p)
	if err != nil {
		return nil, err
	}

	// Check chart dependencies to make sure all are present in /charts
	chartRequested, err := loader.Load(cp)
	if err != nil {
		return nil, err
	}

	if err := checkIfInstallable(chartRequested); err != nil {
		return nil, err
	}

	if chartRequested.Metadata.Deprecated {
		fmt.Printf("This chart is deprecated")
	}

	if req := chartRequested.Metadata.Dependencies; req != nil {
		// If CheckDependencies returns an error, we have unfulfilled dependencies.
		// As of Helm 2.4.0, this is treated as a stopping condition:
		// https://github.com/helm/helm/issues/2209
		if err := action.CheckDependencies(chartRequested, req); err != nil {
			err = errors.Wrap(err, "An error occurred while checking for chart dependencies. You may need to run `helm dependency build` to fetch missing dependencies")
			if client.DependencyUpdate {
				man := &downloader.Manager{
					Out:              out,
					ChartPath:        cp,
					Keyring:          client.ChartPathOptions.Keyring,
					SkipUpdate:       false,
					Getters:          p,
					RepositoryConfig: settings.RepositoryConfig,
					RepositoryCache:  settings.RepositoryCache,
					Debug:            settings.Debug,
				}
				if err := man.Update(); err != nil {
					return nil, err
				}
				// Reload the chart with the updated Chart.lock file.
				if chartRequested, err = loader.Load(cp); err != nil {
					return nil, errors.Wrap(err, "failed reloading chart after repo update")
				}
			} else {
				return nil, err
			}
		}
	}

	client.Namespace = settings.Namespace()

	// Create context and prepare the handle of SIGTERM
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	cSignal := make(chan os.Signal, 2)
	signal.Notify(cSignal, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-cSignal
		fmt.Fprintf(out, "Release %s has been cancelled.\n", args[0])
		cancel()
	}()

	return client.RunWithContext(ctx, chartRequested, vals)
}

// checkIfInstallable validates if a chart can be installed
//
// Application chart type is only installable
func checkIfInstallable(ch *chart.Chart) error {
	switch ch.Metadata.Type {
	case "", "application":
		return nil
	}
	return errors.Errorf("%s charts are not installable", ch.Metadata.Type)
}
