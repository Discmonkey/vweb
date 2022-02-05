package nal

type UnitType = int
type Unit = byte

func Type(n Unit) UnitType {
	// the lowest 5 bits specify the nal type
	return UnitType(n & 63)
}

const (
	Unspecified                             UnitType = 0
	NonIdrSliceLayerWithoutPartitioningRbsp          = 1
	SliceDataPartitionALayerRbsp                     = 2
	SliceDataPartitionBLayerRbsp                     = 3
	SliceDataPartitionCLayerRbsp                     = 4
	IdrSliceLayerWithoutPartitioningRbsp             = 5
	SeiRbsp                                          = 6
	SeqParameterSetRbsp                              = 7
	PicParameterSetRbsp                              = 8
	AccessUnitDelimiterRbsp                          = 9
	EndOfSeqRbsp                                     = 10
	EndOfStreamRbsp                                  = 11
	FillerDataRbsp                                   = 12
	SeqParameterSetExtensionRbsp                     = 13
	PrefixNalUnitRbsp                                = 14
	SubsetSeqParameterSetRbsp                        = 15
	AuxSliceLayerWithoutPartitioningRbsp             = 19
	SliceLayerExtensionRbsp                          = 20
	DepthSliceLayerExtensionRbsp                     = 21
	// Reserved actually takes multiple values
	Reserved = 16
)

func toString(i UnitType) string {
	value := "Unspecified"
	switch i {
	case NonIdrSliceLayerWithoutPartitioningRbsp:
		value = "Coded slice of a non-IDR picture"
	case SliceDataPartitionALayerRbsp:
		value = "Coded slice data partition A"
	case SliceDataPartitionBLayerRbsp:
		value = "Coded slice data partition B"
	case SliceDataPartitionCLayerRbsp:
		value = "Coded slice data partition C"
	case IdrSliceLayerWithoutPartitioningRbsp:
		value = "Coded slice of an IDR picture"
	case SeiRbsp:
		value = "Supplemental enhancement information (SEI)"
	case SeqParameterSetRbsp:
		value = "Sequence parameter set"
	case PicParameterSetRbsp:
		value = "Picture parameter set"
	case AccessUnitDelimiterRbsp:
		value = "Access unit delimiter"
	case EndOfSeqRbsp:
		value = "End of sequence"
	case EndOfStreamRbsp:
		value = "End of stream"
	case FillerDataRbsp:
		value = "Filler data"
	case SeqParameterSetExtensionRbsp:
		value = "Sequence parameter set extension"
	case PrefixNalUnitRbsp:
		value = "Prefix NAL unit"
	case SubsetSeqParameterSetRbsp:
		value = "Subset sequence parameter set"
	case AuxSliceLayerWithoutPartitioningRbsp:
		value = "Coded slice of an auxiliary coded picture without partitioning"
	case SliceLayerExtensionRbsp:
		value = "Coded slice extension"
	case DepthSliceLayerExtensionRbsp:
		value = "Coded slice extension for depth view components"
	case 16, 17, 18, 22, 23:
		value = "Reserved"
	}

	return value
}
