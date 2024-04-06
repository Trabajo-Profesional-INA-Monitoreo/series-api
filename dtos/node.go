package dtos

type Node struct {
	Name              string
	Id                uint64
	ConfiguredStreams *[]ConfiguredStream
}
