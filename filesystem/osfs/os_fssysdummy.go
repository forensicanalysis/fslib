// +build !go1.12

package osfs

// Sys returns a map of item attributes.
func (i *Info) Sys() interface{} {
	return nil
}
