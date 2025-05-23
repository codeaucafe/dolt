// Copyright 2022-2023 Dolthub, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package serial

import (
	"strconv"
)

type Encoding byte

const (
	EncodingNull             Encoding = 0
	EncodingInt8             Encoding = 1
	EncodingUint8            Encoding = 2
	EncodingInt16            Encoding = 3
	EncodingUint16           Encoding = 4
	EncodingInt32            Encoding = 7
	EncodingUint32           Encoding = 8
	EncodingInt64            Encoding = 9
	EncodingUint64           Encoding = 10
	EncodingFloat32          Encoding = 11
	EncodingFloat64          Encoding = 12
	EncodingBit64            Encoding = 13
	EncodingHash128          Encoding = 14
	EncodingYear             Encoding = 15
	EncodingDate             Encoding = 16
	EncodingTime             Encoding = 17
	EncodingDatetime         Encoding = 18
	EncodingEnum             Encoding = 19
	EncodingSet              Encoding = 20
	EncodingBytesAddr        Encoding = 21
	EncodingCommitAddr       Encoding = 22
	EncodingStringAddr       Encoding = 23
	EncodingJSONAddr         Encoding = 24
	EncodingCell             Encoding = 25
	EncodingGeomAddr         Encoding = 26
	EncodingExtendedAddr     Encoding = 27
	EncodingString           Encoding = 128
	EncodingBytes            Encoding = 129
	EncodingDecimal          Encoding = 130
	EncodingJSON             Encoding = 131
	EncodingGeometry         Encoding = 133
	EncodingExtended         Encoding = 134
	EncodingStringAdaptive   Encoding = 135
	EncodingBytesAdaptive    Encoding = 136
	EncodingExtendedAdaptive Encoding = 137
)

var EnumNamesEncoding = map[Encoding]string{
	EncodingNull:             "Null",
	EncodingInt8:             "Int8",
	EncodingUint8:            "Uint8",
	EncodingInt16:            "Int16",
	EncodingUint16:           "Uint16",
	EncodingInt32:            "Int32",
	EncodingUint32:           "Uint32",
	EncodingInt64:            "Int64",
	EncodingUint64:           "Uint64",
	EncodingFloat32:          "Float32",
	EncodingFloat64:          "Float64",
	EncodingBit64:            "Bit64",
	EncodingHash128:          "Hash128",
	EncodingYear:             "Year",
	EncodingDate:             "Date",
	EncodingTime:             "Time",
	EncodingDatetime:         "Datetime",
	EncodingEnum:             "Enum",
	EncodingSet:              "Set",
	EncodingBytesAddr:        "BytesAddr",
	EncodingCommitAddr:       "CommitAddr",
	EncodingStringAddr:       "StringAddr",
	EncodingJSONAddr:         "JSONAddr",
	EncodingCell:             "Cell",
	EncodingGeomAddr:         "GeomAddr",
	EncodingExtendedAddr:     "ExtendedAddr",
	EncodingString:           "String",
	EncodingBytes:            "Bytes",
	EncodingDecimal:          "Decimal",
	EncodingJSON:             "JSON",
	EncodingGeometry:         "Geometry",
	EncodingExtended:         "Extended",
	EncodingStringAdaptive:   "StringAdaptive",
	EncodingBytesAdaptive:    "BytesAdaptive",
	EncodingExtendedAdaptive: "ExtendedAdaptive",
}

var EnumValuesEncoding = map[string]Encoding{
	"Null":             EncodingNull,
	"Int8":             EncodingInt8,
	"Uint8":            EncodingUint8,
	"Int16":            EncodingInt16,
	"Uint16":           EncodingUint16,
	"Int32":            EncodingInt32,
	"Uint32":           EncodingUint32,
	"Int64":            EncodingInt64,
	"Uint64":           EncodingUint64,
	"Float32":          EncodingFloat32,
	"Float64":          EncodingFloat64,
	"Bit64":            EncodingBit64,
	"Hash128":          EncodingHash128,
	"Year":             EncodingYear,
	"Date":             EncodingDate,
	"Time":             EncodingTime,
	"Datetime":         EncodingDatetime,
	"Enum":             EncodingEnum,
	"Set":              EncodingSet,
	"BytesAddr":        EncodingBytesAddr,
	"CommitAddr":       EncodingCommitAddr,
	"StringAddr":       EncodingStringAddr,
	"JSONAddr":         EncodingJSONAddr,
	"Cell":             EncodingCell,
	"GeomAddr":         EncodingGeomAddr,
	"ExtendedAddr":     EncodingExtendedAddr,
	"String":           EncodingString,
	"Bytes":            EncodingBytes,
	"Decimal":          EncodingDecimal,
	"JSON":             EncodingJSON,
	"Geometry":         EncodingGeometry,
	"Extended":         EncodingExtended,
	"StringAdaptive":   EncodingStringAdaptive,
	"BytesAdaptive":    EncodingBytesAdaptive,
	"ExtendedAdaptive": EncodingExtendedAdaptive,
}

func (v Encoding) String() string {
	if s, ok := EnumNamesEncoding[v]; ok {
		return s
	}
	return "Encoding(" + strconv.FormatInt(int64(v), 10) + ")"
}
