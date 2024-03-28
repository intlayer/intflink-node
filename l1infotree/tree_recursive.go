package l1infotree

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

// L1InfoTreeRecursive is a recursive implementation of the L1InfoTree of Feijoa
type L1InfoTreeRecursive struct {
	historicL1InfoTree *L1InfoTree
	snapShot           L1InfoTreeRecursiveSnapshot
}
type L1InfoTreeRecursiveSnapshot struct {
	HistoricL1InfoTreeRoot common.Hash
	L1Data                 common.Hash
	L1InfoTreeRoot         common.Hash
}

// NewL1InfoTreeRecursive creates a new empty L1InfoTreeRecursive
func NewL1InfoTreeRecursive(height uint8) (*L1InfoTreeRecursive, error) {
	historic, err := NewL1InfoTree(height, nil)
	if err != nil {
		return nil, err
	}

	mtr := &L1InfoTreeRecursive{
		historicL1InfoTree: historic,
		snapShot: L1InfoTreeRecursiveSnapshot{
			HistoricL1InfoTreeRoot: common.Hash{},
			L1Data:                 common.Hash{},
			L1InfoTreeRoot:         common.Hash{},
		},
	}

	return mtr, nil
}

// NewL1InfoTreeRecursiveFromLeaves creates a new L1InfoTreeRecursive from leaves as they are
func NewL1InfoTreeRecursiveFromLeaves(height uint8, leaves [][32]byte) (*L1InfoTreeRecursive, error) {
	mtr, err := NewL1InfoTreeRecursive(height)
	if err != nil {
		return nil, err
	}

	for i, leaf := range leaves {
		snapShot, err := mtr.AddLeaf(uint32(i), leaf)
		if err != nil {
			return nil, err
		}
		mtr.snapShot = snapShot
	}
	return mtr, nil
}

// AddLeaf hashes the current historicL1InfoRoot + leaf data into the new leaf value and then adds it to the historicL1InfoTree
func (mt *L1InfoTreeRecursive) AddLeaf(index uint32, leaf [32]byte) (L1InfoTreeRecursiveSnapshot, error) {
	//adds the current l1InfoTreeRoot into the tree to generate the next historicL2InfoTree
	_, err := mt.historicL1InfoTree.AddLeaf(index, mt.snapShot.L1InfoTreeRoot)
	if err != nil {
		return L1InfoTreeRecursiveSnapshot{}, err
	}

	//creates the new snapshot
	snapShot := L1InfoTreeRecursiveSnapshot{}
	snapShot.HistoricL1InfoTreeRoot = mt.historicL1InfoTree.GetRoot()
	snapShot.L1Data = common.BytesToHash(leaf[:])
	snapShot.L1InfoTreeRoot = crypto.Keccak256Hash(snapShot.HistoricL1InfoTreeRoot.Bytes(), snapShot.L1Data.Bytes())
	mt.snapShot = snapShot

	return snapShot, nil
}

// GetRoot returns the root of the L1InfoTreeRecursive
func (mt *L1InfoTreeRecursive) GetRoot() common.Hash {
	return mt.snapShot.L1InfoTreeRoot
}

// GetHistoricRoot returns the root of the HistoricL1InfoTree
func (mt *L1InfoTreeRecursive) GetHistoricRoot() common.Hash {
	return mt.historicL1InfoTree.GetRoot()
}

// ComputeMerkleProof computes the Merkle proof from the leaves
func (mt *L1InfoTreeRecursive) ComputeMerkleProof(gerIndex uint32, leaves [][32]byte) ([][32]byte, common.Hash, error) {
	return mt.historicL1InfoTree.ComputeMerkleProof(gerIndex, leaves)
}
