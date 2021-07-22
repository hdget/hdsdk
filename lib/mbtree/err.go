package mbtree

import "github.com/pkg/errors"

var (
	ErrNodeNotFound        = errors.New("node not found")
	ErrInvalidNode         = errors.New("invalid node")
	ErrInvalidNodeId       = errors.New("invalid node id")
	ErrNodeAlreadyExist    = errors.New("node already exist")
	ErrInvalidSourceDest   = errors.New("source node cannot be ancestor of destination parent node")
	ErrDeleteRootForbidden = errors.New("root node cannot be deleted")
)
