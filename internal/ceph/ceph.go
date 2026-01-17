// Package ceph provides utilities for interacting with Ceph RBD.
package ceph

import (
	"strings"
	"time"
)

// RBDImageInfo represents information about an RBD image.
type RBDImageInfo struct {
	ID     string              `json:"id"`
	Parent *RBDImageInfoParent `json:"parent,omitempty"`
}

// RBDImageInfoParent represents the parent image information.
type RBDImageInfoParent struct {
	Pool     string `json:"pool"`
	Image    string `json:"image"`
	Snapshot string `json:"snapshot"`
}

// RBDTimeStamp represents an RBD timestamp.
type RBDTimeStamp struct {
	time.Time
}

// NewRBDTimeStamp creates a new RBDTimeStamp from a time.Time.
func NewRBDTimeStamp(t time.Time) RBDTimeStamp {
	return RBDTimeStamp{t}
}

// UnmarshalJSON implements json.Unmarshaler for RBDTimeStamp.
func (t *RBDTimeStamp) UnmarshalJSON(data []byte) error {
	var err error
	t.Time, err = time.Parse("Mon Jan  2 15:04:05 2006", strings.Trim(string(data), `"`))

	return err
}

// RBDLock represents an RBD lock.
type RBDLock struct {
	LockID  string `json:"id"`
	Locker  string `json:"locker"`
	Address string `json:"address"`
}

// RBDSnapshot represents an RBD snapshot.
type RBDSnapshot struct {
	ID        int          `json:"id,omitempty"`
	Name      string       `json:"name,omitempty"`
	Size      int64        `json:"size,omitempty"`
	Protected bool         `json:"protected,string,omitempty"`
	Timestamp RBDTimeStamp `json:"timestamp,omitzero"`
}

// CephCmd is the interface for Ceph command operations.
type CephCmd interface {
	RBDClone(pool, srcImage, srcSnap, dstImage, features string) error
	RBDInfo(pool, image string) (*RBDImageInfo, error)
	RBDLs(pool string) ([]string, error)
	RBDLockAdd(pool, image, lockID string) error
	RBDLockLs(pool, image string) ([]*RBDLock, error)
	RBDLockRm(pool, image string, lock *RBDLock) error
	RBDRm(pool, image string) error
	RBDTrashMv(pool, image string) error
	CephRBDTaskAddTrashRemove(pool, image string) error
	RBDSnapCreate(pool, image, snap string) error
	RBDSnapLs(pool, image string) ([]RBDSnapshot, error)
	RBDSnapRm(pool, image, snap string) error
}

type cephCmdImpl struct {
	command command
}

// NewCephCmd creates a new CephCmd.
func NewCephCmd() CephCmd {
	return &cephCmdImpl{
		command: newCommand(),
	}
}

// NewCephCmdWithToolsAndCustomKubectl creates a new CephCmd using kubectl.
func NewCephCmdWithToolsAndCustomKubectl(kubectl []string, namespace string) CephCmd {
	return &cephCmdImpl{
		command: newCommandTools(kubectl, namespace),
	}
}

// NewCephCmdWithTools creates a new CephCmd using the KUBECTL environment variable.
func NewCephCmdWithTools(namespace string) CephCmd {
	if len(envKubectlPath) == 0 {
		panic("KUBECTL environment variable is not set")
	}

	return NewCephCmdWithToolsAndCustomKubectl([]string{envKubectlPath}, namespace)
}
