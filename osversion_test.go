package vz

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"syscall"
	"testing"
)

type nopDoer struct{}

func (*nopDoer) Do(func()) {}

func TestAvailableVersion(t *testing.T) {
	majorMinorVersionOnce = &nopDoer{}
	defer func() {
		majorMinorVersion = 0
		majorMinorVersionOnce = &sync.Once{}
	}()

	t.Run("macOS 11", func(t *testing.T) {
		majorMinorVersion = 10
		cases := map[string]func() error{
			"NewLinuxBootLoader": func() error {
				_, err := NewLinuxBootLoader("")
				return err
			},
			"NewVirtualMachineConfiguration": func() error {
				_, err := NewVirtualMachineConfiguration(nil, 0, 0)
				return err
			},
			"NewFileHandleSerialPortAttachment": func() error {
				_, err := NewFileHandleSerialPortAttachment(nil, nil)
				return err
			},
			"NewFileSerialPortAttachment": func() error {
				_, err := NewFileSerialPortAttachment("", false)
				return err
			},
			"NewVirtioConsoleDeviceSerialPortConfiguration": func() error {
				_, err := NewVirtioConsoleDeviceSerialPortConfiguration(nil)
				return err
			},
			"NewVirtioEntropyDeviceConfiguration": func() error {
				_, err := NewVirtioEntropyDeviceConfiguration()
				return err
			},
			"NewVirtioTraditionalMemoryBalloonDeviceConfiguration": func() error {
				_, err := NewVirtioTraditionalMemoryBalloonDeviceConfiguration()
				return err
			},
			"NewNATNetworkDeviceAttachment": func() error {
				_, err := NewNATNetworkDeviceAttachment()
				return err
			},
			"NewBridgedNetworkDeviceAttachment": func() error {
				_, err := NewBridgedNetworkDeviceAttachment(nil)
				return err
			},
			"NewFileHandleNetworkDeviceAttachment": func() error {
				_, err := NewFileHandleNetworkDeviceAttachment(nil)
				return err
			},
			"NewVirtioNetworkDeviceConfiguration": func() error {
				_, err := NewVirtioNetworkDeviceConfiguration(nil)
				return err
			},
			"NewMACAddress": func() error {
				_, err := NewMACAddress(nil)
				return err
			},
			"NewRandomLocallyAdministeredMACAddress": func() error {
				_, err := NewRandomLocallyAdministeredMACAddress()
				return err
			},
			"NewVirtioSocketDeviceConfiguration": func() error {
				_, err := NewVirtioSocketDeviceConfiguration()
				return err
			},
			"(*VirtioSocketDevice).Listen": func() error {
				_, err := (*VirtioSocketDevice)(nil).Listen(1)
				return err
			},
			"NewDiskImageStorageDeviceAttachment": func() error {
				_, err := NewDiskImageStorageDeviceAttachment("", false)
				return err
			},
			"NewVirtioBlockDeviceConfiguration": func() error {
				_, err := NewVirtioBlockDeviceConfiguration(nil)
				return err
			},
			"NewVirtualMachine": func() error {
				_, err := NewVirtualMachine(nil)
				return err
			},
		}
		for name, fn := range cases {
			err := fn()
			if !errors.Is(err, ErrUnsupportedOSVersion) {
				t.Fatalf("unexpected error %v in %s", err, name)
			}
		}
	})

	t.Run("macOS 12", func(t *testing.T) {
		majorMinorVersion = 11
		cases := map[string]func() error{
			"NewVirtioSoundDeviceConfiguration": func() error {
				_, err := NewVirtioSoundDeviceConfiguration()
				return err
			},
			"NewVirtioSoundDeviceHostInputStreamConfiguration": func() error {
				_, err := NewVirtioSoundDeviceHostInputStreamConfiguration()
				return err
			},
			"NewVirtioSoundDeviceHostOutputStreamConfiguration": func() error {
				_, err := NewVirtioSoundDeviceHostOutputStreamConfiguration()
				return err
			},
			"NewUSBKeyboardConfiguration": func() error {
				_, err := NewUSBKeyboardConfiguration()
				return err
			},
			"NewGenericPlatformConfiguration": func() error {
				_, err := NewGenericPlatformConfiguration()
				return err
			},
			"NewUSBScreenCoordinatePointingDeviceConfiguration": func() error {
				_, err := NewUSBScreenCoordinatePointingDeviceConfiguration()
				return err
			},
			"NewVirtioFileSystemDeviceConfiguration": func() error {
				_, err := NewVirtioFileSystemDeviceConfiguration("")
				return err
			},
			"NewSharedDirectory": func() error {
				_, err := NewSharedDirectory("", false)
				return err
			},
			"NewSingleDirectoryShare": func() error {
				_, err := NewSingleDirectoryShare(nil)
				return err
			},
			"NewMultipleDirectoryShare": func() error {
				_, err := NewMultipleDirectoryShare(nil)
				return err
			},
			"(*VirtualMachine).Stop": func() error {
				return (*VirtualMachine)(nil).Stop()
			},
			"(*VirtualMachine).StartGraphicApplication": func() error {
				return (*VirtualMachine)(nil).StartGraphicApplication(0, 0)
			},
		}
		for name, fn := range cases {
			err := fn()
			if !errors.Is(err, ErrUnsupportedOSVersion) {
				t.Fatalf("unexpected error %v in %s", err, name)
			}
		}
	})

	t.Run("macOS 12.3", func(t *testing.T) {
		majorMinorVersion = 12
		cases := map[string]func() error{
			"BlockDeviceIdentifier": func() error {
				_, err := (*VirtioBlockDeviceConfiguration)(nil).BlockDeviceIdentifier()
				return err
			},
			"SetBlockDeviceIdentifier": func() error {
				return (*VirtioBlockDeviceConfiguration)(nil).SetBlockDeviceIdentifier("")
			},
		}
		for name, fn := range cases {
			err := fn()
			if !errors.Is(err, ErrUnsupportedOSVersion) {
				t.Fatalf("unexpected error %v in %s", err, name)
			}
		}
	})
}

func Test_fetchMajorMinorVersion(t *testing.T) {
	tests := []struct {
		name    string
		sysctl  func(string) (string, error)
		want    float64
		wantErr bool
	}{
		{
			name: "valid 12.0",
			sysctl: func(s string) (string, error) {
				return "12.0", nil
			},
			want:    12,
			wantErr: false,
		},
		{
			name: "valid 12.3",
			sysctl: func(s string) (string, error) {
				return "12.3", nil
			},
			want:    12.3,
			wantErr: false,
		},
		{
			name: "valid 12.3.1",
			sysctl: func(s string) (string, error) {
				return "12.3.1", nil
			},
			want:    12.3,
			wantErr: false,
		},
		{
			name: "invalid unknown",
			sysctl: func(s string) (string, error) {
				return "unknown", nil
			},
			want:    0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sysctl = tt.sysctl
			defer func() {
				sysctl = syscall.Sysctl
			}()

			version, err := fetchMajorMinorVersion()
			if (err != nil) != tt.wantErr {
				t.Errorf("fetchMajorMinorVersion() error = %v, wantErr %v", err, tt.wantErr)
			}
			if version != tt.want {
				t.Errorf("want version %.3f but got %.3f", tt.want, version)
			}
		})
	}
}

func Test_macOSBuildTargetAvailable(t *testing.T) {
	maxAllowedVersionOnce = &nopDoer{}
	defer func() {
		maxAllowedVersionOnce = &sync.Once{}
	}()

	cases := []struct {
		// version is specified only 11, 12, 12.3, 13
		version           float64
		maxAllowedVersion int
		wantErr           bool
		wantErrMsg        string
	}{
		{
			version:           11,
			maxAllowedVersion: 0, // undefined case
			wantErr:           true,
			wantErrMsg:        "undefined __MAC_OS_X_VERSION_MAX_ALLOWED",
		},
		{
			version:           11,
			maxAllowedVersion: 100000,
			wantErr:           true,
			wantErrMsg:        "for 11.0",
		},
		{
			version:           11,
			maxAllowedVersion: 110000,
		},
		{
			version:           12,
			maxAllowedVersion: 110000,
			wantErr:           true,
			wantErrMsg:        "for 12.0",
		},
		{
			version:           12,
			maxAllowedVersion: 120000,
		},
		{
			version:           12,
			maxAllowedVersion: 120100, // __MAC_12_1
		},
		{
			version:           12,
			maxAllowedVersion: 120200, // __MAC_12_2
		},
		{
			version:           12,
			maxAllowedVersion: 120300, // __MAC_12_3
		},
		{
			version:           12,
			maxAllowedVersion: 130000, // __MAC_13_0
		},
		{
			version:           12.3,
			maxAllowedVersion: 120000,
			wantErr:           true,
			wantErrMsg:        "for 12.3",
		},
		{
			version:           12.3,
			maxAllowedVersion: 120300, // __MAC_12_3
		},
		{
			version:           12.3,
			maxAllowedVersion: 130000, // __MAC_13_0
		},
		{
			version:           13,
			maxAllowedVersion: 120300,
			wantErr:           true,
			wantErrMsg:        "for 13.0",
		},
		{
			version:           13,
			maxAllowedVersion: 130000, // __MAC_13_0
		},
	}
	for _, tc := range cases {
		prefix := "valid"
		if tc.wantErr {
			prefix = "invalid"
		}
		name := fmt.Sprintf(
			"%s maxAllowedVersion is %d and API target %.1f",
			prefix,
			tc.maxAllowedVersion,
			tc.version,
		)
		t.Run(name, func(t *testing.T) {
			tmp := maxAllowedVersion
			defer func() { maxAllowedVersion = tmp }()
			maxAllowedVersion = tc.maxAllowedVersion

			err := macOSBuildTargetAvailable(tc.version)
			if (err != nil) != tc.wantErr {
				t.Fatalf("macOSBuildTargetAvailable(%.1f) error = %v, wantErr %v", tc.version, err, tc.wantErr)
			}
			if tc.wantErr {
				got := err.Error()
				if !strings.Contains(got, tc.wantErrMsg) {
					t.Errorf("want msg %q but got %q", tc.wantErrMsg, got)
				}
				if !errors.Is(err, ErrBuildTargetOSVersion) {
					t.Errorf("unexpected unwrap error: %v", err)
				}
			}
		})
	}
}
