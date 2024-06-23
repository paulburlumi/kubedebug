// Package kubedebug provides the means to interactively run "kubectl debug"
package kubedebug

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/huh"
)

// Commander executes commands
type Commander interface {
	Output(name string, arg ...string) (string, error)
	Run(redirect bool, name string, arg ...string) error
}

// KubeDebug interactively runs "kubectl debug"
type KubeDebug struct {
	args          []string
	stderr        io.Writer
	command       Commander
	actions       []func() error
	contextName   string
	namespaceName string
	podName       string
	containerName string
	imageName     string
}

// NewKubeDebug constructs [KubeDebug]
func NewKubeDebug(args []string, stderr io.Writer, command Commander) *KubeDebug {
	kd := &KubeDebug{
		args:    args,
		stderr:  stderr,
		command: command,
	}
	kd.actions = []func() error{
		kd.parseArgs,
		kd.chooseContext,
		kd.setContext,
		kd.chooseNamespace,
		kd.choosePod,
		kd.chooseContainer,
		kd.chooseImage,
		kd.debugContainer,
	}
	return kd
}

// Run the interactive actions
func (kd *KubeDebug) Run() error {
	for _, action := range kd.actions {
		if err := action(); err != nil {
			return fmt.Errorf("could not run action: %w", err)
		}
	}
	return nil
}

func (kd *KubeDebug) parseArgs() error {
	name := filepath.Base(kd.args[0])
	flagSet := flag.NewFlagSet(name, flag.ContinueOnError)
	flagSet.SetOutput(kd.stderr)

	const defaultImage = "busybox"
	const usage = "image name for the debug container"
	flagSet.StringVar(&kd.imageName, "image_name", defaultImage, usage)
	flagSet.StringVar(&kd.imageName, "i", defaultImage, usage+" (shorthand)")

	if err := flagSet.Parse(kd.args[1:]); err != nil {
		return fmt.Errorf("could not parse arguments: %w", err)
	}
	return nil
}

func (kd *KubeDebug) chooseContext() error {
	out, err := kd.command.Output(
		"kubectl",
		"config",
		"get-contexts",
		"--output=name",
	)
	if err != nil {
		return fmt.Errorf("could not output: %w", err)
	}

	contextNames := strings.Split(strings.TrimSpace(out), "\n")

	if len(contextNames) == 0 {
		return errors.New("no kubernetes contexts found")
	}

	kd.contextName = contextNames[0]

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Kubernetes Context").
				Description("In Kubernetes, contexts provide a mechanism to specify the namespace, user and cluster.").
				Options(huh.NewOptions(contextNames...)...).
				Value(&kd.contextName),
		),
	)
	if err := form.Run(); err != nil {
		return fmt.Errorf("could not run form: %w", err)
	}
	return nil
}

func (kd *KubeDebug) setContext() error {
	err := kd.command.Run(
		false,
		"kubectl",
		"config",
		"use-context",
		kd.contextName,
	)
	if err != nil {
		return fmt.Errorf("could not run: %w", err)
	}
	return nil
}

func (kd *KubeDebug) chooseNamespace() error {
	out, err := kd.command.Output(
		"kubectl",
		"get",
		"ns",
		"--no-headers",
		"-o", "custom-columns=:metadata.name")
	if err != nil {
		return fmt.Errorf("could not output: %w", err)
	}

	namespaceNames := strings.Split(strings.TrimSpace(out), "\n")

	if len(namespaceNames) == 0 {
		return errors.New("no kubernetes namespaces found")
	}

	kd.namespaceName = namespaceNames[0]

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Kubernetes Namespace").
				Description("In Kubernetes, namespaces provide a mechanism for isolating groups of resources within a single cluster.").
				Options(huh.NewOptions(namespaceNames...)...).
				Value(&kd.namespaceName),
		),
	)
	if err := form.Run(); err != nil {
		return fmt.Errorf("could not run form: %w", err)
	}
	return nil
}

func (kd *KubeDebug) choosePod() error {
	out, err := kd.command.Output(
		"kubectl",
		"-n", kd.namespaceName,
		"get",
		"pods",
		"-o", "custom-columns=:metadata.name",
	)
	if err != nil {
		return fmt.Errorf("could not output: %w", err)
	}

	podNames := strings.Split(strings.TrimSpace(out), "\n")

	if len(podNames) == 0 {
		return errors.New("no kubernetes pods found")
	}

	kd.podName = podNames[0]

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Kubernetes Pod").
				Description("In Kubernetes, pods are similar to a set of containers with shared namespaces and shared filesystem volumes.").
				Options(huh.NewOptions(podNames...)...).
				Value(&kd.podName),
		),
	)
	if err := form.Run(); err != nil {
		return fmt.Errorf("could not run form: %w", err)
	}
	return nil
}

func (kd *KubeDebug) chooseContainer() error {
	out, err := kd.command.Output(
		"kubectl",
		"-n", kd.namespaceName,
		"get",
		"pods",
		kd.podName,
		"-o", "jsonpath={.spec.containers[*].name}",
	)
	if err != nil {
		return fmt.Errorf("could not output: %w", err)
	}

	containerNames := strings.Split(strings.TrimSpace(out), "\n")

	if len(containerNames) == 0 {
		return errors.New("no kubernetes containers found")
	}

	kd.containerName = containerNames[0]

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Choose a Kubernetes Container").
				Description("In Kubernetes, pods can contain a single container or multiple containers that form a single cohesive unit.").
				Options(huh.NewOptions(containerNames...)...).
				Value(&kd.containerName),
		),
	)
	if err := form.Run(); err != nil {
		return fmt.Errorf("could not run form: %w", err)
	}
	return nil
}

func (kd *KubeDebug) chooseImage() error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Choose an Image").
				Description("The image to use for the debug container.").
				Validate(func(s string) error {
					if len(s) == 0 {
						return errors.New("the container image cannot be empty")
					}
					return nil
				}).
				Value(&kd.imageName),
		),
	)
	if err := form.Run(); err != nil {
		return fmt.Errorf("could not run form: %w", err)
	}
	return nil
}

func (kd *KubeDebug) debugContainer() error {
	err := kd.command.Run(
		true,
		"kubectl",
		"-n", kd.namespaceName,
		"debug",
		kd.podName,
		"-it",
		"--target="+kd.containerName,
		"--image="+kd.imageName,
	)
	if err != nil {
		return fmt.Errorf("could not run: %w", err)
	}
	return nil
}
