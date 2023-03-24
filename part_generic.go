package mc

import smc "github.com/surfinggo/mc"

func SliceContains[V comparable](slice []V, v V) bool {
	return smc.SliceContains(slice, v)
}

func PointerTo[V any](v V) *V {
	return smc.PointerTo(v)
}

func PtrTo[V any](v V) *V {
	return smc.PointerTo(v)
}

func VarOr[V comparable](v V, defaultValue V) V {
	return smc.VarOr(v, defaultValue)
}
