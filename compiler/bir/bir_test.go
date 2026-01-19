package bir

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

func main() {
	file, err := os.Open(filepath.Join("testdata", "sample.bir"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	bir := NewBir()

	err = bir.Read(kaitai.NewStream(file), nil, bir)
	if err != nil {
		panic(err)
	}

	for _, f := range bir.Module.Functions {
		if name, ok := bir.ConstantPool.ConstantPoolEntries[f.NameCpIndex].CpInfo.(*Bir_StringCpInfo); ok {
			fmt.Println(name.Value)
		}
	}
}
