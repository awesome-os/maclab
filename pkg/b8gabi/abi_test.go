package b8gabi

import (
	"testing"
	"unsafe"
)

func TestArenaHeaderV1LayoutAndRoundTrip(t *testing.T) {
	if got := unsafe.Sizeof(ArenaHeaderV1{}); got != ArenaAlignment {
		t.Fatalf("sizeof(ArenaHeaderV1) = %d, want %d", got, ArenaAlignment)
	}
	if B8G_ABI_MAGIC_V3 != 0x33564742 {
		t.Fatalf("B8G_ABI_MAGIC_V3 = 0x%08x, want 0x33564742", B8G_ABI_MAGIC_V3)
	}

	hdr, err := NewArenaHeaderV1(4096, 7)
	if err != nil {
		t.Fatalf("NewArenaHeaderV1 failed: %v", err)
	}
	hdr.UsedBytes = 512

	raw, err := hdr.MarshalBinary()
	if err != nil {
		t.Fatalf("MarshalBinary failed: %v", err)
	}
	if len(raw) != ArenaAlignment {
		t.Fatalf("len(raw) = %d, want %d", len(raw), ArenaAlignment)
	}

	var decoded ArenaHeaderV1
	if err := decoded.UnmarshalBinary(raw); err != nil {
		t.Fatalf("UnmarshalBinary failed: %v", err)
	}
	if decoded != hdr {
		t.Fatalf("round-trip mismatch: got %+v, want %+v", decoded, hdr)
	}
}

func TestArenaHeaderV1RejectsUnalignedSize(t *testing.T) {
	if _, err := NewArenaHeaderV1(100, 1); err == nil {
		t.Fatal("expected error for non-64-byte-aligned arena size")
	}
}

func TestDefaultKernelTargets(t *testing.T) {
	targets := DefaultKernelTargets()
	if len(targets) != 2 {
		t.Fatalf("len(targets) = %d, want 2", len(targets))
	}
	if targets[0].Target != KernelTargetT2Linux || targets[0].ABIMagic != B8G_ABI_MAGIC_V3 {
		t.Fatalf("unexpected t2linux target: %+v", targets[0])
	}
	if targets[1].Target != KernelTargetT2XNU || targets[1].ABIMagic != B8G_ABI_MAGIC_V3 {
		t.Fatalf("unexpected t2xnu target: %+v", targets[1])
	}
}
