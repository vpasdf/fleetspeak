// +build windows

package config

import (
	"github.com/google/fleetspeak/fleetspeak/src/comtesting"
	"github.com/google/fleetspeak/fleetspeak/src/windows/regutil"
	"golang.org/x/sys/windows/registry"
	"io/ioutil"
	"path/filepath"
	"reflect"
	"sort"
	"testing"
)

func TestReadServices(t *testing.T) {
	defer registry.DeleteKey(registry.LOCAL_MACHINE, `SOFTWARE\TestingOnly`)

	tmpPath, fin := comtesting.GetTempDir("ReadServices")
	defer fin()

	fooConfigPath := filepath.Join(tmpPath, "config.txt")

	if err := ioutil.WriteFile(fooConfigPath, []byte(`name: "foo1"`), 0644); err != nil {
		t.Fatal(err)
	}

	if err := regutil.WriteStringValue(`HKEY_LOCAL_MACHINE\Software\TestingOnly\textservices`, "foo", fooConfigPath); err != nil {
		t.Fatal(err)
	}

	if err := regutil.WriteStringValue(`HKEY_LOCAL_MACHINE\Software\TestingOnly\textservices`, "bar", `name: "bar1"`); err != nil {
		t.Fatal(err)
	}

	handler, err := NewWindowsRegistryPersistenceHandler(`HKEY_LOCAL_MACHINE\Software\TestingOnly`, false)
	if err != nil {
		t.Fatal(err)
	}

	services, err := handler.ReadServices()
	if err != nil {
		t.Fatal(err)
	}
	if len(services) != 2 {
		t.Fatalf("Got too few services: %v", len(services))
	}

	actualServiceNames := []string{services[0].Name, services[1].Name}
	sort.Strings(actualServiceNames)
	expectedServiceNames := []string{"bar1", "foo1"}

	if !reflect.DeepEqual(actualServiceNames, expectedServiceNames) {
		t.Fatalf("Failed, expected: %v, actual: %v", expectedServiceNames, actualServiceNames)
	}
}
