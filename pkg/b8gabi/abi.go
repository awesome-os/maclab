// Package b8gabi provides the Awesome-OS B8G_ABI_MAGIC_V3 and 64-byte aligned
// B8G_ARENA_V1 bridge types plus t2linux/t2xnu kernel target manifest constants
// for maclab Apple Silicon hardware-in-the-loop kernel testing.
package b8gabi

import (
	"encoding/binary"
	"errors"
	"fmt"
)

const (
	// B8G_ABI_MAGIC_V3 is the canonical Awesome-OS B8G v3 ABI magic ("BGV3" LE).
	B8G_ABI_MAGIC_V3 uint32 = 0x33564742

	// ArenaAlignment is the mandatory 64-byte alignment for B8G_ARENA_V1 slabs.
	ArenaAlignment = 64

	// ArenaVersionV1 is the B8G_ARENA_V1 header version.
	ArenaVersionV1 uint16 = 1

	// KernelTargetT2Linux identifies awesome-os/t2linux Linux kernel artifacts.
	KernelTargetT2Linux = "awesome-os/t2linux"

	// KernelTargetT2XNU identifies awesome-os/t2xnu XNU/Darwin kernel artifacts.
	KernelTargetT2XNU = "awesome-os/t2xnu"

	// BootloaderLimineUKI targets Limine unified kernel images on Apple Silicon ESP.
	BootloaderLimineUKI = "limine-uki"

	// BootloaderGRUBESP targets GRUB one-shot maclab.env on FAT ESP.
	BootloaderGRUBESP = "grub-esp"

	// BootloaderM1N1XNU targets m1n1 / OOB DFU staging for t2xnu bring-up.
	BootloaderM1N1XNU = "m1n1-xnu"
)

// ArenaHeaderV1 is the 64-byte aligned B8G_ARENA_V1 header used to exchange
// zero-copy telemetry, kernel artifacts, and checkpoint buffers across B8G peers.
type ArenaHeaderV1 struct {
	Magic      uint32
	Version    uint16
	Flags      uint16
	ArenaBytes uint64
	UsedBytes  uint64
	Generation uint64
	TargetKind uint32
	ArchID     uint32
	Reserved   [3]uint64
}

// KernelTargetManifest describes a maclab test target for t2linux or t2xnu.
type KernelTargetManifest struct {
	Target     string
	Bootloader string
	Arch       string
	ABIMagic   uint32
}

// DefaultKernelTargets returns the canonical Awesome-OS kernel test targets.
func DefaultKernelTargets() []KernelTargetManifest {
	return []KernelTargetManifest{
		{
			Target:     KernelTargetT2Linux,
			Bootloader: BootloaderLimineUKI,
			Arch:       "arm64",
			ABIMagic:   B8G_ABI_MAGIC_V3,
		},
		{
			Target:     KernelTargetT2XNU,
			Bootloader: BootloaderM1N1XNU,
			Arch:       "arm64",
			ABIMagic:   B8G_ABI_MAGIC_V3,
		},
	}
}

// NewArenaHeaderV1 constructs a validated 64-byte aligned B8G_ARENA_V1 header.
func NewArenaHeaderV1(arenaBytes uint64, generation uint64) (ArenaHeaderV1, error) {
	h := ArenaHeaderV1{
		Magic:      B8G_ABI_MAGIC_V3,
		Version:    ArenaVersionV1,
		ArenaBytes: arenaBytes,
		Generation: generation,
	}
	if err := h.Validate(); err != nil {
		return ArenaHeaderV1{}, err
	}
	return h, nil
}

// Validate checks B8G_ABI_MAGIC_V3, version, and 64-byte alignment invariants.
func (h ArenaHeaderV1) Validate() error {
	if h.Magic != B8G_ABI_MAGIC_V3 {
		return fmt.Errorf("invalid B8G ABI magic: got 0x%08x, want 0x%08x", h.Magic, B8G_ABI_MAGIC_V3)
	}
	if h.Version != ArenaVersionV1 {
		return fmt.Errorf("unsupported B8G_ARENA_V1 version: %d", h.Version)
	}
	if h.ArenaBytes%ArenaAlignment != 0 {
		return fmt.Errorf("arena size %d is not %d-byte aligned", h.ArenaBytes, ArenaAlignment)
	}
	if h.UsedBytes > h.ArenaBytes {
		return errors.New("used_bytes exceeds arena_bytes")
	}
	return nil
}

// MarshalBinary encodes the header into a 64-byte little-endian buffer.
func (h ArenaHeaderV1) MarshalBinary() ([]byte, error) {
	if err := h.Validate(); err != nil {
		return nil, err
	}
	buf := make([]byte, ArenaAlignment)
	binary.LittleEndian.PutUint32(buf[0:4], h.Magic)
	binary.LittleEndian.PutUint16(buf[4:6], h.Version)
	binary.LittleEndian.PutUint16(buf[6:8], h.Flags)
	binary.LittleEndian.PutUint64(buf[8:16], h.ArenaBytes)
	binary.LittleEndian.PutUint64(buf[16:24], h.UsedBytes)
	binary.LittleEndian.PutUint64(buf[24:32], h.Generation)
	binary.LittleEndian.PutUint32(buf[32:36], h.TargetKind)
	binary.LittleEndian.PutUint32(buf[36:40], h.ArchID)
	binary.LittleEndian.PutUint64(buf[40:48], h.Reserved[0])
	binary.LittleEndian.PutUint64(buf[48:56], h.Reserved[1])
	binary.LittleEndian.PutUint64(buf[56:64], h.Reserved[2])
	return buf, nil
}

// UnmarshalBinary decodes and validates a 64-byte little-endian header.
func (h *ArenaHeaderV1) UnmarshalBinary(data []byte) error {
	if len(data) != ArenaAlignment {
		return fmt.Errorf("expected %d bytes, got %d", ArenaAlignment, len(data))
	}
	decoded := ArenaHeaderV1{
		Magic:      binary.LittleEndian.Uint32(data[0:4]),
		Version:    binary.LittleEndian.Uint16(data[4:6]),
		Flags:      binary.LittleEndian.Uint16(data[6:8]),
		ArenaBytes: binary.LittleEndian.Uint64(data[8:16]),
		UsedBytes:  binary.LittleEndian.Uint64(data[16:24]),
		Generation: binary.LittleEndian.Uint64(data[24:32]),
		TargetKind: binary.LittleEndian.Uint32(data[32:36]),
		ArchID:     binary.LittleEndian.Uint32(data[36:40]),
		Reserved: [3]uint64{
			binary.LittleEndian.Uint64(data[40:48]),
			binary.LittleEndian.Uint64(data[48:56]),
			binary.LittleEndian.Uint64(data[56:64]),
		},
	}
	if err := decoded.Validate(); err != nil {
		return err
	}
	*h = decoded
	return nil
}
