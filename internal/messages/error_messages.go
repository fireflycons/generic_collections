package messages

const (
	COLLECTION_MODIFIED      = "Collection has been modified"
	COLLECTION_EMPTY         = "Cannot perform operation on empty collection"
	NEGATIVE_CAPACITY        = "Cannot create collection with negative capacity"
	FOREIGN_NODE             = "Node does not belong to this list"
	NIL_NODE                 = "Cannot perform operation on nil node"
	SET_POINTER_MODIFICATION = "Cannot modify set elements through pointer"
	COMP_FN_NIL              = "Comparer function cannot be nil"
	HASH_BUCKET_SIZE_INVALID = "Hash bucket size cannot be less than 1"
	COMPARER_INVALID_INT_FMT = "Unsupported integer byte size %d"
	COMPARER_INVALID_KEY_FMT = "Unsupported key type %T of kind %v. Supply instance of CompararFunc[T]"
	ARG_NIL_FMT              = "Argument %s cannot be nil"
	ARG_OUT_OF_RANGE_FMT     = "Argument %s out of range"
	SLICE_TOO_SMALL          = "Slice is too small to receive all elements"
	AGG_SLICE_EMPTY          = "Cannot compute aggregate of empty slice"
)
