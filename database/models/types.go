package models

import "github.com/lib/pq"

func ConvertUint32ToInt32Array(a []uint32) pq.Int32Array {
	r := make([]int32, len(a))
	for i, u := range a {
		r[i] = int32(u)
	}
	return pq.Int32Array(r)
}
