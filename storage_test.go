package vz_test

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/Code-Hex/vz/v2"
)

func TestBlockDeviceIdentifier(t *testing.T) {
	if vz.Available(12.3) {
		t.Skip("VirtioBlockDeviceConfiguration.SetBlockDeviceIdentifier is supported from macOS 12.3")
	}
	dir := t.TempDir()
	path := filepath.Join(dir, "disk.img")
	if err := vz.CreateDiskImage(path, 512); err != nil {
		t.Fatal(err)
	}

	attachment, err := vz.NewDiskImageStorageDeviceAttachment(path, false)
	if err != nil {
		t.Fatal(err)
	}
	config, err := vz.NewVirtioBlockDeviceConfiguration(attachment)
	if err != nil {
		t.Fatal(err)
	}
	got1, err := config.BlockDeviceIdentifier()
	if err != nil {
		t.Fatal(err)
	}
	if got1 != "" {
		t.Fatalf("want empty by default: %q", got1)
	}

	invalidID := strings.Repeat("h", 25)
	if err := config.SetBlockDeviceIdentifier(invalidID); err == nil {
		t.Fatal("want error")
	} else {
		nserr, ok := err.(*vz.NSError)
		if !ok {
			t.Fatalf("unexpected error: %v", err)
		}
		if nserr.Domain != "VZErrorDomain" {
			t.Errorf("unexpected NSError domain: %v", nserr)
		}
		if nserr.Code != int(vz.ErrorInvalidVirtualMachineConfiguration) {
			t.Errorf("unexpected NSError code: %v", nserr)
		}
	}

	want := "hello"
	if err := config.SetBlockDeviceIdentifier(want); err != nil {
		t.Fatal(err)
	}
	got2, err := config.BlockDeviceIdentifier()
	if err != nil {
		t.Fatal(err)
	}
	if got2 != want {
		t.Fatalf("want %q but got %q", want, got2)
	}
}
