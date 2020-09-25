/*Copyright [2019] housepower

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package column

import (
	"github.com/housepower/clickhouse_sinker/column/impls"
)

var (
	columns = map[string]IColumn{}
)

type creator func() IColumn

func regist(name string, creator creator) {
	columns[name] = creator()
}

// GetColumnByName get the IColumn by the name of type
func GetColumnByName(name string) IColumn {
	return columns[name]
}

// init register column types for different data types
func init() {
	regist("UInt8", func() IColumn {
		return impls.NewIntColumn(8, false)
	})
	regist("UInt16", func() IColumn {
		return impls.NewIntColumn(16, false)
	})
	regist("UInt32", func() IColumn {
		return impls.NewIntColumn(32, false)
	})
	regist("UInt64", func() IColumn {
		return impls.NewIntColumn(64, false)
	})

	regist("Int8", func() IColumn {
		return impls.NewIntColumn(8, false)
	})
	regist("Int16", func() IColumn {
		return impls.NewIntColumn(16, false)
	})
	regist("Int32", func() IColumn {
		return impls.NewIntColumn(32, false)
	})
	regist("Int64", func() IColumn {
		return impls.NewIntColumn(64, false)
	})

	regist("Date", func() IColumn {
		return impls.NewIntColumn(16, true)
	})

	regist("DateTime", func() IColumn {
		return impls.NewIntColumn(32, true)
	})

	regist("DateTime64", func() IColumn {
		return impls.NewIntColumn(64, true)
	})

	regist("Float32", func() IColumn {
		return impls.NewFloatColumn(32)
	})
	regist("Float64", func() IColumn {
		return impls.NewFloatColumn(64)
	})

	regist("String", func() IColumn {
		return impls.NewStringColumn()
	})

	regist("FixedString", func() IColumn {
		return impls.NewStringColumn()
	})
}
