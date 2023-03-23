package mc

import "github.com/surfinggo/mc"

func SliceContains[V comparable](slice []V, v V) bool {
	return mc.SliceContains(slice, v)
}

func PointerTo[V any](v V) *V {
	return mc.PointerTo(v)
}

func PtrTo[V any](v V) *V {
	return mc.PointerTo(v)
}

func VarOr[V comparable](v V, defaultValue V) V {
	return mc.VarOr(v, defaultValue)
}
