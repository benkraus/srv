// Package migrationdisk holds the cold-disk migration messages exchanged by
// srv-front, the srv service and the VM runner. It stays free of runtime
// dependencies so srv-front builds on platforms without the VM stack.
package migrationdisk

import (
	"srv/internal/coldstore"
	"srv/internal/storage"
)

// Transfer describes immutable transfer objects, not a space reservation.
// Base is retained at the destination; Disk is discarded after publication.
type Transfer struct {
	Disk coldstore.Object  `json:"disk"`
	Base *coldstore.Object `json:"base,omitempty"`
}

type Request struct {
	Action       string                    `json:"action"`
	Name         string                    `json:"name"`
	InstanceID   string                    `json:"instance_id"`
	Token        string                    `json:"token"`
	ExpectedDisk storage.SubvolumeIdentity `json:"expected_disk"`
	Transfer     Transfer                  `json:"transfer"`
}

type Result struct {
	Transfer          Transfer `json:"transfer"`
	BasePresent       bool     `json:"base_present"`
	SerializedPayload bool     `json:"serialized_payload"`
}
