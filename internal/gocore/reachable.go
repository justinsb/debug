package gocore

import (
	"fmt"

	"golang.org/x/debug/internal/core"
)

type ReachableInfo struct {
	ReachableFrom map[core.Address]int64
	Roots         map[core.Address]bool
}

func computeRoots(c *Process) *ReachableInfo {
	dominatorSizes := make(map[core.Address]int64)
	root := make(map[core.Address]bool)

	done := make(map[core.Address]bool)

	c.ForEachObject(func(src Object) bool {
		c.ForEachPtr(src, func(srcOffset int64, dest Object, destOffset int64) bool {
			root[core.Address(dest)] = false
			return true
		})
		return true
	})
	c.ForEachObject(func(src Object) bool {
		if _, found := root[core.Address(src)]; !found {
			root[core.Address(src)] = true
		}
		return true
	})

	for {
		count := 0
		c.ForEachObject(func(src Object) bool {
			if done[core.Address(src)] {
				return true
			}
			isDone := true

			dominatorSize := int64(0)
			c.ForEachPtr(src, func(srcOffset int64, dest Object, destOffset int64) bool {
				if dest == 0 {
					return true
				}
				destAddr := core.Address(dest)
				if !done[destAddr] {
					isDone = false
					return false
				}
				dominatorSize += dominatorSizes[destAddr]
				return true
			})

			if isDone {
				size := c.Size(src)

				dominatorSizes[core.Address(src)] = size + dominatorSize
				done[core.Address(src)] = true
				count++
			}

			return true
		})
		fmt.Printf("marked %d\n", count)
		if count == 0 {
			break
		}
	}

	return &ReachableInfo{
		ReachableFrom: dominatorSizes,
		Roots:         root,
	}
}
