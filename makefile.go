//+build mage

//Build, test and more ... a developers everyday tool
package main

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/bclicn/color"
	"github.com/gosuri/uilive"
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var (
	backend     = os.Getenv("FTA_BACKEND")
	environment = "local"
	services    []string
)

const (
	all  = "all"
	done = "done"
)

type Helper mg.Namespace
type Build mg.Namespace
type Docker mg.Namespace
type Clean mg.Namespace

// Check for valid make arguments
func (Helper) args() {
	services = []string{}
	service := os.Getenv("service")
	// Single service to be processed
	if service != all {
		services = append(services, service)
	}
	// All services to be processed
	if service == all {
		// List all folders
		cmdOut, err := sh.Output("bash", "-c", "ls -d ${FTA_BACKEND}/*/server")
		if err != nil {
			fmt.Printf("args bash ls -d */server error: %s", err)
		}
		services = strings.Split(cmdOut, "\n")
		// Clean up prefix and postfix of the found path
		servicesClean := []string{}
		for _, service := range services {
			serviceClean := strings.TrimPrefix(service, backend+"/")
			serviceClean = strings.TrimSuffix(serviceClean, "/server")
			servicesClean = append(servicesClean, serviceClean)
		}
		services = servicesClean
	}
}

// Remove the .temp folder
func (Clean) temp() error {
	mg.Deps(Helper.args)
	wg := sync.WaitGroup{}
	// Loop over all services
	for _, service := range services {
		// Sync goroutines by +1 to the waiting group
		wg.Add(1)
		// Prallelize execution
		go func(service string) {
			tempFolder := backend + fmt.Sprintf("/%s/.temp", service)
			// Clean any existing .temp folder
			if _, err := os.Stat(tempFolder); !os.IsNotExist(err) {
				err := sh.Run("rm", "-rf", tempFolder)
				if err != nil {
					fmt.Printf(color.Red("Clean:Temp rm ./.temp error: %s\n"), err)
				}
			}
			// Create new .temp folder for future building and packing
			err := os.Mkdir(tempFolder, 0755)
			if err != nil {
				fmt.Printf(color.Red("Clean:Temp Mkdir error: %s\n"), err)
			}
			// Sync goroutines by -1 to the waiting group
			wg.Done()
		}(service)
	}
	// Wait for all goroutines to finish
	wg.Wait()
	return nil
}

// Transpile protobuffer definitions
func (Build) Protoc() error {
	// find all *.protoc files
	cmdOut, err := sh.Output("bash", "-c", fmt.Sprintf("find %s -type f -name \"*.proto\" -exec ls {} \\;", "./"))
	if err != nil {
		fmt.Printf(color.Red("Build:Protoc find .*.proto error: %s\n"), err)
		return err
	}
	files := strings.Split(cmdOut, "\n")
	// Loop over all found protoc files
	for _, file := range files {
		// Skip vendor transpiling
		if strings.Contains(file, "vendor/") || strings.Contains(file, "go/pkg") {
			continue
		}
		// Transpile .proto file
		err := sh.Run("protoc", "--proto_path=.", "--twirp_out=.", "--go_out=.", file)
		if err != nil {
			fmt.Printf(color.Red("Build:Protoc protoc transpile error: %s\n"), err)
			return err
		}
	}
	return nil
}

// Compile the service
func (Build) Build() error {
	mg.Deps(Helper.args)
	mg.Deps(Clean.temp)
	// Create multi-line updateable stdout
	writer := uilive.New()
	writer.Start()
	// Printing build status to stdout
	buildStatus := map[string]string{}
	// Loop over all service
	for _, service := range services {
		servicePath := backend + "/" + service
		err := sh.Run("env", "CGO_ENABLED=0", "GOOS=linux", "go", "build", "-o", servicePath+"/.temp/"+service, servicePath+"/server/"+service)
		if err != nil {
			fmt.Printf(color.Red("Build:Build go build ... error: %s\n"), err)
			return err
		}
		// Copy Dockerfile to .temp folder
		err = sh.Run("cp", servicePath+"/Dockerfile", servicePath+"/.temp")
		if err != nil {
			fmt.Printf(color.Red("Build:Image cp Dockerfile error: %s\n"), err)
			return err
		}
		// Set build status to 'done'
		buildStatus[service] = done
	}
	writer.Stop() // flush and stop rendering
	return nil
}

// Create a Docker image
func (Docker) Image() error {
	mg.Deps(Build.Build)
	wg := sync.WaitGroup{}
	// Create multi-line updateable stdout
	writer := uilive.New()
	writer.Start()
	// Printing imaging status to stdout
	imageStatus := map[string]string{}
	for _, service := range services {
		fmt.Fprintf(writer, "Creating image for %s service...\n", service)
	}
	// Loop over all services
	for _, service := range services {
		// Sync goroutines by +1 to the waiting group
		wg.Add(1)
		// Parallelize execution
		go func(service string) {
			servicePath := backend + "/" + service
			// Build image
			err := sh.Run("docker", "build", "-t", fmt.Sprintf("%s:%s", service, environment), "-f", fmt.Sprintf("%s/Dockerfile", servicePath), fmt.Sprintf("%s/.temp/.", servicePath))
			if err != nil {
				fmt.Printf(color.Red("Build:Image docker build error: %s\n"), err)
			}
			// Cleanup .temp folder
			err = sh.Run("rm", "-rf", fmt.Sprint("%s/.temp", servicePath))
			if err != nil {
				fmt.Printf(color.Red("Build:Image cp certificate  ... error: %s\n"), err)
			}
			// Set build status to 'done'
			imageStatus[service] = done
			// Update imaging status to stdout
			for _, s := range services {
				fmt.Fprintf(writer, "Created image for %s... %s\n", s, imageStatus[s])
			}
			// sync goroutines by -1 to the waiting group
			wg.Done()
		}(service)
	}
	// wait for all goroutines to finish
	wg.Wait()
	writer.Stop() // flush and stop rendering
	return nil
}

// Integration Test the service
func (Build) Test() error {
	mg.Deps(Helper.args)
	if len(services) != 1 {
		return errors.New("Testing should only invoke one service")
	}
	service := services[0]
	fmt.Printf("%s testing service %s...\n", color.Cyan("FTA"), service)
	cmdOut, err := sh.Output("go", "run", ""+backend+"/"+service+"/test/main.go")
	fmt.Printf("Result: %+v", cmdOut)
	return err
}

// Unit Test the service
func (Build) Unittest() error {
	mg.Deps(Helper.args)
	// Create multi-line updateable stdout
	writer := uilive.New()
	writer.Start()
	// Printing imaging status to stdout
	testStatus := map[string]string{}
	for _, service := range services {
		fmt.Fprintf(writer, "Unit Testing %s service...\n", service)
	}
	// Loop over all services
	for _, service := range services {
		// Parallelize execution
		testPath := fmt.Sprintf("%s/%s/...", backend, service)
		if service != "webserver" {
			testPath = fmt.Sprintf("%s/%s/server/...", backend, service)
		}
		// UnitTest service
		cmdOut, _ := sh.Output("go", "test", testPath)
		fmt.Printf("Result: %+v", cmdOut)
		// Set test status to 'done'
		testStatus[service] = done
	}
	writer.Stop() // flush and stop rendering
	return nil
}
