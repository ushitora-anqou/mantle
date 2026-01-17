package ceph

import (
	"errors"
	"slices"
)

// ErrSnapshotNotFound is returned when a snapshot is not found.
var ErrSnapshotNotFound = errors.New("snapshot not found")

// IsRBDSnapshotNamed returns a predicate function that checks if a snapshot has the given name.
func IsRBDSnapshotNamed(name string) func(RBDSnapshot) bool {
	return func(s RBDSnapshot) bool {
		return s.Name == name
	}
}

// FindRBDSnapshot finds an RBD snapshot by name.
func FindRBDSnapshot(cmd CephCmd, pool, image, snap string) (*RBDSnapshot, error) {
	snaps, err := cmd.RBDSnapLs(pool, image)
	if err != nil {
		return nil, err
	}

	i := slices.IndexFunc(snaps, IsRBDSnapshotNamed(snap))
	if i == -1 {
		return nil, ErrSnapshotNotFound
	}

	return &snaps[i], nil
}
