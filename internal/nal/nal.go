package nal

type NalUnitType = int
type Nal = byte

func Type(n Nal) NalUnitType {

}

const (
	Unspecified NalUnitType = 0
	Unspecified
"Coded slice of a non-IDR picture
slice_layer_without_partitioning_rbsp( )"
"Coded slice data partition A
slice_data_partition_a_layer_rbsp( )"
"Coded slice data partition B
slice_data_partition_b_layer_rbsp( )"
"Coded slice data partition C
slice_data_partition_c_layer_rbsp( )"
"Coded slice of an IDR picture
slice_layer_without_partitioning_rbsp( )"
"Supplemental enhancement information (SEI)
sei_rbsp( )"
"Sequence parameter set
seq_parameter_set_rbsp( )"
"Picture parameter set
pic_parameter_set_rbsp( )"
"Access unit delimiter
access_unit_delimiter_rbsp( )"
"End of sequence
end_of_seq_rbsp( )"
"End of stream
end_of_stream_rbsp( )"
"Filler data
filler_data_rbsp( )"
"Sequence parameter set extension
seq_parameter_set_extension_rbsp( )"
"Prefix NAL unit
prefix_nal_unit_rbsp( )"
"Subset sequence parameter set
subset_seq_parameter_set_rbsp( )"
Reserved
"Coded slice of an auxiliary coded picture without partitioning
slice_layer_without_partitioning_rbsp( )"
"Coded slice extension
slice_layer_extension_rbsp( )"
"Coded slice extension for depth view components
slice_layer_extension_rbsp( )
(specified in Annex I)"
Reserved
Unspecified
)
