// Copyright (c) 2026, WSO2 LLC. (http://www.wso2.com).
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

package main

import (
	"fmt"
	"os"
	"path/filepath"

	"ballerina-lang-go/compiler/bir"
	birmodel "ballerina-lang-go/compiler/bir/model"

	"github.com/kaitai-io/kaitai_struct_go_runtime/kaitai"
)

func main() {
	file, err := os.Open(filepath.Join("..", "testdata", "bal_workspace.bir"))
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Method 1: Direct Kaitai read (original approach)
	fmt.Println("=== Method 1: Direct Kaitai Read ===")
	file1, _ := os.Open(filepath.Join("..", "testdata", "bal_workspace.bir"))
	birModel := bir.NewBir()
	err = birModel.Read(kaitai.NewStream(file1), nil, birModel)
	if err != nil {
		panic(err)
	}
	file1.Close()

	fmt.Println("ImportCount", birModel.Module.ImportCount)
	for _, imp := range birModel.Module.Imports {
		if name, ok := birModel.ConstantPool.ConstantPoolEntries[imp.PackageNameIndex].CpInfo.(*bir.Bir_StringCpInfo); ok {
			fmt.Printf("Import: %s\n", name.Value)
		}
	}
	for _, f := range birModel.Module.Functions {
		if name, ok := birModel.ConstantPool.ConstantPoolEntries[f.NameCpIndex].CpInfo.(*bir.Bir_StringCpInfo); ok {
			fmt.Printf("Function: %s\n", name.Value)
		}
	}

	fmt.Println("\n=== Method 2: Using Loader (BIRNode Model) ===")
	// Method 2: Using our loader to populate BIRNode model
	file2, _ := os.Open(filepath.Join("..", "testdata", "bal_workspace.bir"))
	pkg, err := bir.LoadBIRPackageFromReader(file2)
	if err != nil {
		panic(err)
	}
	file2.Close()

	fmt.Printf("Package: %s/%s:%s\n",
		pkg.GetPackageID().GetOrgName().GetValue(),
		pkg.GetPackageID().GetName().GetValue(),
		pkg.GetPackageID().GetPackageVersion().GetValue())

	fmt.Printf("Imports: %d\n", len(*pkg.GetImportModules()))
	fmt.Printf("Functions: %d\n", len(*pkg.GetFunctions()))

	for _, imp := range *pkg.GetImportModules() {
		fmt.Printf("Import: %s/%s:%s\n",
			imp.GetPackageID().GetOrgName().GetValue(),
			imp.GetPackageID().GetName().GetValue(),
			imp.GetPackageID().GetPackageVersion().GetValue())
	}

	for _, fn := range *pkg.GetFunctions() {
		fmt.Printf("Function: %s (original: %s, worker: %s, flags: %d)\n",
			fn.GetName().GetValue(),
			fn.GetOriginalName().GetValue(),
			fn.GetWorkerName().GetValue(),
			fn.GetFlags())

		params := fn.GetRequiredParams()
		if params != nil && len(*params) > 0 {
			fmt.Printf("  Required params:\n")
			for _, p := range *params {
				fmt.Printf("    - %s (flags: %d)\n", p.GetName().GetValue(), p.GetFlags())
			}
		}

		// Print basic blocks and instructions
		if fn.GetBasicBlocks() != nil {
			fmt.Printf("  Basic Blocks: %d\n", len(*fn.GetBasicBlocks()))
			for _, bb := range *fn.GetBasicBlocks() {
				fmt.Printf("    Basic Block: %s\n", bb.GetId().GetValue())

				// Print non-terminator instructions with details
				if bb.GetInstructions() != nil {
					fmt.Printf("      Non-Terminator Instructions: %d\n", len(*bb.GetInstructions()))
					for i, ins := range *bb.GetInstructions() {
						// Get kind from abstract instruction
						if absIns, ok := ins.(birmodel.BIRAbstractInstruction); ok {
							kindStr := getInstructionKindString(absIns.GetKind())
							fmt.Printf("        [%d] %s", i+1, kindStr)
						} else {
							fmt.Printf("        [%d] <unknown>", i+1)
						}

						// Print instruction-specific details
						printNonTerminatorDetails(ins)
						fmt.Println()
					}
				}

				// Print terminator with details
				if bb.GetTerminator() != nil {
					if absIns, ok := bb.GetTerminator().(birmodel.BIRAbstractInstruction); ok {
						kindStr := getInstructionKindString(absIns.GetKind())
						fmt.Printf("      Terminator: %s", kindStr)
						printTerminatorDetails(bb.GetTerminator())
						fmt.Println()
					} else {
						fmt.Printf("      Terminator: <unknown>\n")
					}
				}
			}
		}
	}
}

// printNonTerminatorDetails prints details specific to each instruction type
func printNonTerminatorDetails(ins birmodel.BIRNonTerminator) {
	switch instruction := ins.(type) {
	case *birmodel.BIRNonTerminatorMove:
		if instruction.GetLhsOperand() != nil {
			fmt.Printf(" -> LHS: %s", getOperandString(instruction.GetLhsOperand()))
		}
		if instruction.RhsOp != nil {
			fmt.Printf(", RHS: %s", getOperandString(instruction.RhsOp))
		}

	case *birmodel.BIRNonTerminatorBinaryOp:
		if instruction.GetLhsOperand() != nil {
			fmt.Printf(" -> LHS: %s", getOperandString(instruction.GetLhsOperand()))
		}
		if instruction.RhsOp1 != nil {
			fmt.Printf(", RHS1: %s", getOperandString(instruction.RhsOp1))
		}
		if instruction.RhsOp2 != nil {
			fmt.Printf(", RHS2: %s", getOperandString(instruction.RhsOp2))
		}

	case *birmodel.BIRNonTerminatorUnaryOP:
		if instruction.GetLhsOperand() != nil {
			fmt.Printf(" -> LHS: %s", getOperandString(instruction.GetLhsOperand()))
		}
		if instruction.RhsOp != nil {
			fmt.Printf(", RHS: %s", getOperandString(instruction.RhsOp))
		}

	case *birmodel.BIRNonTerminatorConstantLoad:
		if instruction.GetLhsOperand() != nil {
			fmt.Printf(" -> LHS: %s", getOperandString(instruction.GetLhsOperand()))
		}
		if instruction.Value != nil {
			fmt.Printf(", Value: %v", instruction.Value)
		}
		if instruction.Type != nil {
			fmt.Printf(", Type: %s", instruction.Type.String())
		}

	case *birmodel.BIRNonTerminatorNewStructure:
		if instruction.GetLhsOp() != nil {
			fmt.Printf(" -> LHS: %s", getOperandString(instruction.GetLhsOp()))
		}
		if instruction.RhsOp != nil {
			fmt.Printf(", RHS: %s", getOperandString(instruction.RhsOp))
		}
		if len(instruction.InitialValues) > 0 {
			fmt.Printf(", InitialValues: %d entries", len(instruction.InitialValues))
		}

	case *birmodel.BIRNonTerminatorNewArray:
		if instruction.GetLhsOp() != nil {
			fmt.Printf(" -> LHS: %s", getOperandString(instruction.GetLhsOp()))
		}
		if instruction.SizeOp != nil {
			fmt.Printf(", Size: %s", getOperandString(instruction.SizeOp))
		}
		if len(instruction.Values) > 0 {
			fmt.Printf(", Values: %d entries", len(instruction.Values))
		}
		if instruction.Type != nil {
			fmt.Printf(", Type: %s", instruction.Type.String())
		}

	case *birmodel.BIRNonTerminatorFieldAccess:
		if instruction.GetLhsOp() != nil {
			fmt.Printf(" -> LHS: %s", getOperandString(instruction.GetLhsOp()))
		}
		if instruction.KeyOp != nil {
			fmt.Printf(", Key: %s", getOperandString(instruction.KeyOp))
		}
		if instruction.RhsOp != nil {
			fmt.Printf(", RHS: %s", getOperandString(instruction.RhsOp))
		}
		if instruction.OptionalFieldAccess {
			fmt.Printf(" [optional]")
		}
		if instruction.FillingRead {
			fmt.Printf(" [filling]")
		}

	case *birmodel.BIRNonTerminatorNewError:
		if instruction.GetLhsOp() != nil {
			fmt.Printf(" -> LHS: %s", getOperandString(instruction.GetLhsOp()))
		}
		if instruction.MessageOp != nil {
			fmt.Printf(", Message: %s", getOperandString(instruction.MessageOp))
		}
		if instruction.CauseOp != nil {
			fmt.Printf(", Cause: %s", getOperandString(instruction.CauseOp))
		}
		if instruction.DetailOp != nil {
			fmt.Printf(", Detail: %s", getOperandString(instruction.DetailOp))
		}

	case *birmodel.BIRNonTerminatorTypeCast:
		if instruction.GetLhsOp() != nil {
			fmt.Printf(" -> LHS: %s", getOperandString(instruction.GetLhsOp()))
		}
		if instruction.RhsOp != nil {
			fmt.Printf(", RHS: %s", getOperandString(instruction.RhsOp))
		}
		if instruction.Type != nil {
			fmt.Printf(", CastType: %s", instruction.Type.String())
		}
		if instruction.CheckTypes {
			fmt.Printf(" [checkTypes]")
		}

	case *birmodel.BIRNonTerminatorIsLike:
		if instruction.GetLhsOp() != nil {
			fmt.Printf(" -> LHS: %s", getOperandString(instruction.GetLhsOp()))
		}
		if instruction.RhsOp != nil {
			fmt.Printf(", RHS: %s", getOperandString(instruction.RhsOp))
		}
		if instruction.Type != nil {
			fmt.Printf(", Type: %s", instruction.Type.String())
		}

	case *birmodel.BIRNonTerminatorTypeTest:
		if instruction.GetLhsOp() != nil {
			fmt.Printf(" -> LHS: %s", getOperandString(instruction.GetLhsOp()))
		}
		if instruction.RhsOp != nil {
			fmt.Printf(", RHS: %s", getOperandString(instruction.RhsOp))
		}
		if instruction.Type != nil {
			fmt.Printf(", Type: %s", instruction.Type.String())
		}

	default:
		// For other instruction types, just show if it has LHS operand
		if assignIns, ok := ins.(birmodel.BIRAssignInstruction); ok {
			if assignIns.GetLhsOperand() != nil {
				fmt.Printf(" -> LHS: %s", getOperandString(assignIns.GetLhsOperand()))
			}
		}
	}
}

// getInstructionKindString returns a readable string for instruction kind
func getInstructionKindString(kind birmodel.InstructionKind) string {
	kindMap := map[birmodel.InstructionKind]string{
		birmodel.INSTRUCTION_KIND_MOVE:                "MOVE",
		birmodel.INSTRUCTION_KIND_CONST_LOAD:          "CONST_LOAD",
		birmodel.INSTRUCTION_KIND_NEW_STRUCTURE:       "NEW_STRUCTURE",
		birmodel.INSTRUCTION_KIND_NEW_ARRAY:           "NEW_ARRAY",
		birmodel.INSTRUCTION_KIND_MAP_LOAD:            "MAP_LOAD",
		birmodel.INSTRUCTION_KIND_ARRAY_LOAD:          "ARRAY_LOAD",
		birmodel.INSTRUCTION_KIND_MAP_STORE:           "MAP_STORE",
		birmodel.INSTRUCTION_KIND_ARRAY_STORE:         "ARRAY_STORE",
		birmodel.INSTRUCTION_KIND_NEW_ERROR:           "NEW_ERROR",
		birmodel.INSTRUCTION_KIND_TYPE_CAST:           "TYPE_CAST",
		birmodel.INSTRUCTION_KIND_IS_LIKE:             "IS_LIKE",
		birmodel.INSTRUCTION_KIND_TYPE_TEST:           "TYPE_TEST",
		birmodel.INSTRUCTION_KIND_ADD:                 "ADD",
		birmodel.INSTRUCTION_KIND_SUB:                 "SUB",
		birmodel.INSTRUCTION_KIND_MUL:                 "MUL",
		birmodel.INSTRUCTION_KIND_DIV:                 "DIV",
		birmodel.INSTRUCTION_KIND_MOD:                 "MOD",
		birmodel.INSTRUCTION_KIND_EQUAL:               "EQUAL",
		birmodel.INSTRUCTION_KIND_NOT_EQUAL:           "NOT_EQUAL",
		birmodel.INSTRUCTION_KIND_GREATER_THAN:        "GREATER_THAN",
		birmodel.INSTRUCTION_KIND_GREATER_EQUAL:       "GREATER_EQUAL",
		birmodel.INSTRUCTION_KIND_LESS_THAN:           "LESS_THAN",
		birmodel.INSTRUCTION_KIND_LESS_EQUAL:          "LESS_EQUAL",
		birmodel.INSTRUCTION_KIND_AND:                 "AND",
		birmodel.INSTRUCTION_KIND_OR:                  "OR",
		birmodel.INSTRUCTION_KIND_TYPEOF:              "TYPEOF",
		birmodel.INSTRUCTION_KIND_NOT:                 "NOT",
		birmodel.INSTRUCTION_KIND_NEGATE:              "NEGATE",
		birmodel.INSTRUCTION_KIND_GOTO:                "GOTO",
		birmodel.INSTRUCTION_KIND_RETURN:              "RETURN",
		birmodel.INSTRUCTION_KIND_BRANCH:              "BRANCH",
		birmodel.INSTRUCTION_KIND_CALL:                "CALL",
		birmodel.INSTRUCTION_KIND_ASYNC_CALL:          "ASYNC_CALL",
		birmodel.INSTRUCTION_KIND_FP_CALL:             "FP_CALL",
		birmodel.INSTRUCTION_KIND_LOCK:                "LOCK",
		birmodel.INSTRUCTION_KIND_FIELD_LOCK:          "FIELD_LOCK",
		birmodel.INSTRUCTION_KIND_UNLOCK:              "UNLOCK",
		birmodel.INSTRUCTION_KIND_PANIC:               "PANIC",
		birmodel.INSTRUCTION_KIND_WAIT:                "WAIT",
		birmodel.INSTRUCTION_KIND_FLUSH:               "FLUSH",
		birmodel.INSTRUCTION_KIND_WK_RECEIVE:          "WK_RECEIVE",
		birmodel.INSTRUCTION_KIND_WK_SEND:             "WK_SEND",
		birmodel.INSTRUCTION_KIND_WK_ALT_RECEIVE:      "WK_ALT_RECEIVE",
		birmodel.INSTRUCTION_KIND_WK_MULTIPLE_RECEIVE: "WK_MULTIPLE_RECEIVE",
		birmodel.INSTRUCTION_KIND_WAIT_ALL:            "WAIT_ALL",
	}
	if str, ok := kindMap[kind]; ok {
		return str
	}
	return fmt.Sprintf("INSTRUCTION_KIND_%d", kind)
}

// printTerminatorDetails prints details specific to each terminator type
func printTerminatorDetails(term birmodel.BIRTerminator) {
	switch t := term.(type) {
	case *birmodel.BIRTerminatorGOTO:
		if t.TargetBB != nil {
			fmt.Printf(" -> Target: %s", t.TargetBB.GetId().GetValue())
		}

	case *birmodel.BIRTerminatorCall:
		if t.Name.GetValue() != "" {
			fmt.Printf(" -> Call: %s", t.Name.GetValue())
		}
		if len(t.Args) > 0 {
			fmt.Printf(", Args: %d", len(t.Args))
		}
		if t.GetLhsOperand() != nil {
			fmt.Printf(", LHS: %s", getOperandString(t.GetLhsOperand()))
		}
		if t.ThenBB != nil {
			fmt.Printf(", ThenBB: %s", t.ThenBB.GetId().GetValue())
		}

	case *birmodel.BIRTerminatorAsyncCall:
		if t.Name.GetValue() != "" {
			fmt.Printf(" -> AsyncCall: %s", t.Name.GetValue())
		}
		if len(t.Args) > 0 {
			fmt.Printf(", Args: %d", len(t.Args))
		}
		if len(t.AnnotAttachments) > 0 {
			fmt.Printf(", Annots: %d", len(t.AnnotAttachments))
		}

	case *birmodel.BIRTerminatorFPCall:
		if t.Fp != nil {
			fmt.Printf(" -> FP: %s", getOperandString(t.Fp))
		}
		if len(t.Args) > 0 {
			fmt.Printf(", Args: %d", len(t.Args))
		}
		if t.IsAsync {
			fmt.Printf(" [async]")
		}

	case *birmodel.BIRTerminatorReturn:
		// No additional details for return

	case *birmodel.BIRTerminatorBranch:
		if t.Op != nil {
			fmt.Printf(" -> Op: %s", getOperandString(t.Op))
		}
		if t.TrueBB != nil {
			fmt.Printf(", TrueBB: %s", t.TrueBB.GetId().GetValue())
		}
		if t.FalseBB != nil {
			fmt.Printf(", FalseBB: %s", t.FalseBB.GetId().GetValue())
		}

	case birmodel.BIRTerminatorLock:
		if lockedBB := t.GetThenBB(); lockedBB != nil {
			fmt.Printf(" -> LockedBB: %s", lockedBB.GetId().GetValue())
		}

	case *birmodel.BIRTerminatorFieldLock:
		if t.Field != "" {
			fmt.Printf(" -> Field: %s", t.Field)
		}
		if t.LockedBB != nil {
			fmt.Printf(", LockedBB: %s", t.LockedBB.GetId().GetValue())
		}

	case *birmodel.BIRTerminatorUnlock:
		if t.UnlockBB != nil {
			fmt.Printf(" -> UnlockBB: %s", t.UnlockBB.GetId().GetValue())
		}

	case *birmodel.BIRTerminatorPanic:
		if t.ErrorOp != nil {
			fmt.Printf(" -> Error: %s", getOperandString(t.ErrorOp))
		}

	case *birmodel.BIRTerminatorWait:
		if len(t.ExprList) > 0 {
			fmt.Printf(" -> Exprs: %d", len(t.ExprList))
		}
		if t.GetLhsOp() != nil {
			fmt.Printf(", LHS: %s", getOperandString(t.GetLhsOp()))
		}

	case *birmodel.BIRTerminatorFlush:
		if len(t.Channels) > 0 {
			fmt.Printf(" -> Channels: %d", len(t.Channels))
		}
		if t.GetLhsOp() != nil {
			fmt.Printf(", LHS: %s", getOperandString(t.GetLhsOp()))
		}

	case *birmodel.BIRTerminatorWorkerReceive:
		if t.WorkerName.GetValue() != "" {
			fmt.Printf(" -> Worker: %s", t.WorkerName.GetValue())
		}
		if t.IsSameStrand {
			fmt.Printf(" [same strand]")
		}

	case *birmodel.BIRTerminatorWorkerSend:
		if t.Channel.GetValue() != "" {
			fmt.Printf(" -> Channel: %s", t.Channel.GetValue())
		}
		if t.Data != nil {
			fmt.Printf(", Data: %s", getOperandString(t.Data))
		}
		if t.IsSync {
			fmt.Printf(" [sync]")
		}

	case *birmodel.BIRTerminatorWorkerAlternateReceive:
		if len(t.Channels) > 0 {
			fmt.Printf(" -> Channels: %d", len(t.Channels))
		}

	case *birmodel.BIRTerminatorWorkerMultipleReceive:
		if len(t.ReceiveFields) > 0 {
			fmt.Printf(" -> Fields: %d", len(t.ReceiveFields))
		}

	case *birmodel.BIRTerminatorWaitAll:
		if len(t.Keys) > 0 {
			fmt.Printf(" -> Keys: %d", len(t.Keys))
		}
		if len(t.ValueExprs) > 0 {
			fmt.Printf(", ValueExprs: %d", len(t.ValueExprs))
		}
	}
}

// getOperandString returns a string representation of an operand
func getOperandString(op birmodel.BIROperand) string {
	if op == nil {
		return "<nil>"
	}
	if varDcl := op.GetVariableDcl(); varDcl != nil {
		name := varDcl.GetName()
		if name.GetValue() != "" {
			return name.GetValue()
		}
		return "<unnamed>"
	}
	return "<no variable>"
}
