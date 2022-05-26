package types

import (
	"bytes"
	h "github.com/meshplus/crypto-standard/hash"
	"sort"
)

// Inode struct
type Inode struct {
	Key   []byte
	Value []byte
	Hash  []byte
}

// Inodes struct
type Inodes []*Inode

// ProofNode struct
type ProofNode struct {
	IsData bool
	Key    []byte
	Hash   []byte
	Inodes Inodes
	Index  int
}

// ProofPath struct
type ProofPath []*ProofNode

// Validate validates poof with given key
func Validate(key []byte, proof ProofPath) bool {
	if len(proof) == 0 || !proof[len(proof)-1].IsData {
		return false
	}
	var nextHash []byte
	for _, elem := range proof {
		current := elem
		if len(nextHash) != 0 && !bytes.Equal(nextHash, current.Hash) {
			return false
		}
		index := sort.Search(len(current.Inodes), func(i int) bool { return bytes.Compare(current.Inodes[i].Key, key) != -1 })
		exact := len(current.Inodes) > 0 && index < len(current.Inodes) && bytes.Equal(current.Inodes[index].Key, key)
		if !exact {
			index--
		}
		if index != elem.Index || (current.IsData && !bytes.Equal(current.Inodes[index].Key, key)) {
			return false
		}
		res := CalProofNodeHash(current)
		if !bytes.Equal(res, current.Hash) {
			return false
		}
		nextHash = current.Inodes[index].Hash
	}
	return true
}

// CalProofNodeHash calculate hash for given ProofNode
func CalProofNodeHash(node *ProofNode) []byte {
	buff := make([]byte, 0)
	if node.IsData {
		for _, in := range node.Inodes {
			buff = append(buff, in.Key...)
			buff = append(buff, in.Value...)
		}
	} else {
		for _, in := range node.Inodes {
			buff = append(buff, in.Hash...)
		}
	}
	hasher := h.NewHasher(h.KECCAK_256)
	res, _ := hasher.Hash(buff)
	return res
}
