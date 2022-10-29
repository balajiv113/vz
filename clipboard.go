package vz

/*
#cgo darwin CFLAGS: -x objective-c -fno-objc-arc
#cgo darwin LDFLAGS: -lobjc -framework Foundation -framework Virtualization
# include "virtualization_13.h"
*/
import "C"
import "runtime"

// SpiceAgentPortAttachment is an attachment point that enables
// the Spice clipboard sharing capability.
//
// see: https://developer.apple.com/documentation/virtualization/vzspiceagentportattachment?language=objc
type SpiceAgentPortAttachment struct {
	pointer

	*baseSerialPortAttachment

	enabledSharesClipboard bool
}

var _ SerialPortAttachment = (*SpiceAgentPortAttachment)(nil)

// NewSpiceAgentPortAttachment creates a new Spice agent port attachment.
//
// This is only supported on macOS 13 and newer, ErrUnsupportedOSVersion will
// be returned on older versions.
func NewSpiceAgentPortAttachment() (*SpiceAgentPortAttachment, error) {
	if macosMajorVersionLessThan(13) {
		return nil, ErrUnsupportedOSVersion
	}
	spiceAgent := &SpiceAgentPortAttachment{
		pointer: pointer{
			ptr: C.newVZSpiceAgentPortAttachment(),
		},
		enabledSharesClipboard: true,
	}
	runtime.SetFinalizer(spiceAgent, func(self *SpiceAgentPortAttachment) {
		self.Release()
	})
	return spiceAgent, nil
}

// SetSharesClipboard sets enable the Spice agent clipboard sharing capability.
func (s *SpiceAgentPortAttachment) SetSharesClipboard(enable bool) {
	C.setSharesClipboardVZSpiceAgentPortAttachment(
		s.Ptr(),
		C.bool(enable),
	)
	s.enabledSharesClipboard = enable
}

// SharesClipboard returns enable the Spice agent clipboard sharing capability.
func (s *SpiceAgentPortAttachment) SharesClipboard() bool { return s.enabledSharesClipboard }

// SpiceAgentPortAttachmentName returns the Spice agent port name.
func SpiceAgentPortAttachmentName() (string, error) {
	if macosMajorVersionLessThan(13) {
		return "", ErrUnsupportedOSVersion
	}
	cstring := (*char)(C.getSpiceAgentPortName())
	return cstring.String(), nil
}
