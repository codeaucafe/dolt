// Copyright 2020 Dolthub, Inc.
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

package typeinfo

import (
	"context"
	"encoding/gob"
	"fmt"
	"strings"

	"github.com/dolthub/go-mysql-server/sql"
	gmstypes "github.com/dolthub/go-mysql-server/sql/types"

	"github.com/dolthub/dolt/go/store/types"
)

const (
	enumTypeParam_Collation = "collate"
	enumTypeParam_Values    = "vals"
)

// This is a dolt implementation of the MySQL type Enum, thus most of the functionality
// within is directly reliant on the go-mysql-server implementation.
type enumType struct {
	sqlEnumType sql.EnumType
}

var _ TypeInfo = (*enumType)(nil)

func CreateEnumTypeFromParams(params map[string]string) (TypeInfo, error) {
	var collation sql.CollationID
	var err error
	if collationStr, ok := params[enumTypeParam_Collation]; ok {
		collation, err = sql.ParseCollation("", collationStr, false)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf(`create enum type info is missing param "%v"`, enumTypeParam_Collation)
	}
	var values []string
	if valuesStr, ok := params[enumTypeParam_Values]; ok {
		dec := gob.NewDecoder(strings.NewReader(valuesStr))
		err = dec.Decode(&values)
		if err != nil {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf(`create enum type info is missing param "%v"`, enumTypeParam_Values)
	}
	sqlEnumType, err := gmstypes.CreateEnumType(values, collation)
	if err != nil {
		return nil, err
	}
	return CreateEnumTypeFromSqlEnumType(sqlEnumType), nil
}

func CreateEnumTypeFromSqlEnumType(sqlEnumType sql.EnumType) TypeInfo {
	return &enumType{sqlEnumType}
}

// ConvertNomsValueToValue implements TypeInfo interface.
func (ti *enumType) ConvertNomsValueToValue(v types.Value) (interface{}, error) {
	if val, ok := v.(types.Uint); ok {
		return uint16(val), nil
	}
	if _, ok := v.(types.Null); ok || v == nil {
		return nil, nil
	}
	return nil, fmt.Errorf(`"%v" cannot convert NomsKind "%v" to a value`, ti.String(), v.Kind())
}

// ReadFrom reads a go value from a noms types.CodecReader directly
func (ti *enumType) ReadFrom(_ *types.NomsBinFormat, reader types.CodecReader) (interface{}, error) {
	k := reader.ReadKind()
	switch k {
	case types.UintKind:
		return uint16(reader.ReadUint()), nil
	case types.NullKind:
		return nil, nil
	}

	return nil, fmt.Errorf(`"%v" cannot convert NomsKind "%v" to a value`, ti.String(), k)
}

// ConvertValueToNomsValue implements TypeInfo interface.
func (ti *enumType) ConvertValueToNomsValue(ctx context.Context, vrw types.ValueReadWriter, v interface{}) (types.Value, error) {
	if v == nil {
		return types.NullValue, nil
	}
	val, _, err := ti.sqlEnumType.Convert(ctx, v)
	if err != nil {
		return nil, err
	}
	return types.Uint(val.(uint16)), nil
}

// Equals implements TypeInfo interface.
func (ti *enumType) Equals(other TypeInfo) bool {
	if other == nil {
		return false
	}
	if ti2, ok := other.(*enumType); ok && ti.sqlEnumType.NumberOfElements() == ti2.sqlEnumType.NumberOfElements() {
		tiVals := ti.sqlEnumType.Values()
		ti2Vals := ti2.sqlEnumType.Values()
		for i := range tiVals {
			if tiVals[i] != ti2Vals[i] {
				return false
			}
		}
		return true
	}
	return false
}

// FormatValue implements TypeInfo interface.
func (ti *enumType) FormatValue(v types.Value) (*string, error) {
	if _, ok := v.(types.Null); ok || v == nil {
		return nil, nil
	}
	convVal, err := ti.ConvertNomsValueToValue(v)
	if err != nil {
		return nil, err
	}
	val, ok := ti.sqlEnumType.At(int(convVal.(uint16)))
	if !ok {
		return nil, fmt.Errorf(`"%v" has unexpectedly encountered a value of type "%T" from embedded type`, ti.String(), v)
	}
	return &val, nil
}

// GetTypeIdentifier implements TypeInfo interface.
func (ti *enumType) GetTypeIdentifier() Identifier {
	return EnumTypeIdentifier
}

// GetTypeParams implements TypeInfo interface.
func (ti *enumType) GetTypeParams() map[string]string {
	var sb strings.Builder
	enc := gob.NewEncoder(&sb)
	err := enc.Encode(ti.sqlEnumType.Values())
	// this should never error, encoding an array of strings should always succeed
	if err != nil {
		panic(err)
	}
	return map[string]string{
		enumTypeParam_Collation: ti.sqlEnumType.Collation().String(),
		enumTypeParam_Values:    sb.String(),
	}
}

// IsValid implements TypeInfo interface.
func (ti *enumType) IsValid(v types.Value) bool {
	if val, ok := v.(types.Uint); ok {
		_, ok := ti.sqlEnumType.At(int(val))
		if !ok {
			return false
		}
		return true
	}
	if _, ok := v.(types.Null); ok || v == nil {
		return true
	}
	return false
}

// NomsKind implements TypeInfo interface.
func (ti *enumType) NomsKind() types.NomsKind {
	return types.UintKind
}

// Promote implements TypeInfo interface.
func (ti *enumType) Promote() TypeInfo {
	return &enumType{ti.sqlEnumType.Promote().(sql.EnumType)}
}

// String implements TypeInfo interface.
func (ti *enumType) String() string {
	return fmt.Sprintf(`Enum(Collation: %v, Values: %v)`, ti.sqlEnumType.Collation().String(), strings.Join(ti.sqlEnumType.Values(), ", "))
}

// ToSqlType implements TypeInfo interface.
func (ti *enumType) ToSqlType() sql.Type {
	return ti.sqlEnumType
}

// enumTypeConverter is an internal function for GetTypeConverter that handles the specific type as the source TypeInfo.
func enumTypeConverter(ctx context.Context, src *enumType, destTi TypeInfo) (tc TypeConverter, needsConversion bool, err error) {
	switch dest := destTi.(type) {
	case *bitType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *blobStringType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *boolType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *datetimeType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *decimalType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *enumType:
		return func(ctx context.Context, vrw types.ValueReadWriter, v types.Value) (types.Value, error) {
			if v == nil || v == types.NullValue {
				return types.NullValue, nil
			}
			val, ok := v.(types.Uint)
			if !ok {
				return nil, fmt.Errorf("unexpected type converting enum to other enum: %T", v)
			}
			valStr, ok := src.sqlEnumType.At(int(val))
			if !ok {
				return nil, fmt.Errorf("%s does not contain an equivalent value of %d", src.sqlEnumType.String(), val)
			}
			newVal, _, err := dest.sqlEnumType.Convert(ctx, valStr)
			if err != nil {
				return nil, err
			}
			return types.Uint(newVal.(uint16)), nil
		}, true, nil
	case *floatType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *geomcollType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *geometryType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *inlineBlobType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *intType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *jsonType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *linestringType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *multilinestringType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *multipointType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *multipolygonType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *pointType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *polygonType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *setType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *timeType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *uintType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *uuidType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *varBinaryType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *varStringType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	case *yearType:
		return wrapConvertValueToNomsValue(dest.ConvertValueToNomsValue)
	default:
		return nil, false, UnhandledTypeConversion.New(src.String(), destTi.String())
	}
}
