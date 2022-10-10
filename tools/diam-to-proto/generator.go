// go run generator.go -help
// Usage of generator:
//   -d value
//         Comma separated list of folders to load (default ./dict)
//   -intf string
//         Comma separated list (no spaces) of interfaces from [gx, gy, rx, sh, sy] (default "gx,gy")
//   -numberFormat string
//         Filed number format: seq or avpcode (default "seq")
// Example: go run generator.go -d ./dict -d ./custom -intf gx,gy,rx

package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"unicode"

	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
)

var fields []CompositeField
var parsedFields map[string]CompositeField = make(map[string]CompositeField)

var apps = map[string]uint32{
	"gy": 4,
	"sh": 16777217,
	"rx": 16777236,
	"gx": 16777238,
	"sy": 16777303,
}

type newDictionaryBuilder interface {
	load(paths *FlagSet) error
	build(name string, priority int, node *Node) CompositeField
	search(appId, vendorId uint32, code interface{}) (*dict.AVP, error)
}

type Dictionary struct {
	P *dict.Parser
	newDictionaryBuilder
}

type Node struct {
	appId    uint32
	rules    []*dict.Rule
	vendorId uint32
}

type FlagSet struct {
	elements map[string]bool
}

func (l *FlagSet) String() string {
	var names string
	for name := range l.elements {
		names += name
	}
	return names
}

func (l *FlagSet) Set(name string) error {
	l.elements[name] = true
	return nil
}

func main() {

	folders := &FlagSet{elements: map[string]bool{"./dict": true}}
	intf := flag.String("intf", "gx,gy", "Comma separated list (no spaces) of interfaces from [gx, gy, rx, sh, sy]")
	protoNumberFormat := flag.String("numberFormat", "seq", "Field number format: seq or avpcode")
	flag.Var(folders, "d", "Comma separated list of folders to load")
	flag.Parse()

	var enabledApps = make(map[uint32]bool)
	for _, id := range strings.Split(*intf, ",") {
		enabledApps[apps[id]] = true
	}

	dictionary := &Dictionary{}
	if err := dictionary.load(folders); err != nil {
		log.Fatalf("Failed to load dictionaries: %s", err)
	}

	var priority int = 0

	for _, app := range dictionary.P.Apps() {
		if enabledApps[app.ID] {
			for _, command := range app.Command {
				request := fmt.Sprintf("%s%s", app.Name, command.Name)
				replacer := strings.NewReplacer("TGPP", "", " ", "", "-", "")
				request = replacer.Replace(request)

				var vendorId = uint32(dict.UndefinedVendorID)
				if len(app.Vendor) > 0 {
					vendorId = app.Vendor[0].ID
				}
				reqField := dictionary.build(request+"RequestPB", priority,
					&Node{appId: app.ID, rules: command.Request.Rule, vendorId: vendorId},
				)
				fields = append(fields, reqField)
				priority++
				ansField := dictionary.build(request+"AnswerPB", priority,
					&Node{appId: app.ID, rules: command.Answer.Rule, vendorId: vendorId},
				)
				fields = append(fields, ansField)
				priority++
			}
		}
	}

	for _, parsedField := range parsedFields {
		fields = append(fields, parsedField)
	}

	sort.SliceStable(fields, func(i, j int) bool {
		diff := fields[i].priority - fields[j].priority
		if diff == 0 {
			diff = len(fields[j].fields) - len(fields[i].fields)
			if diff == 0 {
				return fields[i].name < fields[j].name
			}
		}
		return diff < 0
	})

	for _, v := range fields {
		// fmt.Printf("%s %s {\n", v.protoDataType, v.name)
		fmt.Printf("message %s {\n", v.name)
		if v.protoDataType == "enum" {
			fmt.Println("\tvalue Value = 1;")
			fmt.Println("\tenum value {")
		}

		// ascending sort fields based on avp codes
		if *protoNumberFormat == "avpcode" {
			sort.SliceStable(v.fields, func(i, j int) bool {
				return v.fields[i].GetCode() < v.fields[j].GetCode()
			})
		}

		for i, f := range v.fields {
			if *protoNumberFormat == "avpcode" {
				f.SetIndex(int(f.GetCode()))
			} else {
				f.SetIndex(i + 1)
			}
			fmt.Println(f)
		}
		if v.protoDataType == "enum" {
			fmt.Println("\t}")
		}
		fmt.Println("}")
		fmt.Println()
	}
}

func (d *Dictionary) load(paths *FlagSet) error {
	if d.P == nil {
		d.P, _ = dict.NewParser()
	}
	for path := range paths.elements {
		log.Printf("Loading dictionaries from %s", path)
		err := filepath.WalkDir(path, func(path string, info fs.DirEntry, err error) error {
			if err != nil {
				log.Println(err)
				return err
			}
			if info.IsDir() {
				return nil
			}
			dictErr := d.P.LoadFile(path)
			if dictErr != nil {
				log.Printf("Failed to load dictionary: %s: %s", path, dictErr)
				return dictErr
			}
			return nil
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Dictionary) build(name string, priority int, node *Node) CompositeField {
	composite := CompositeField{name: name, priority: priority, protoDataType: "message"}
	for _, r := range node.rules {
		avp, err := d.search(node.appId, node.vendorId, r.AVP)
		if err != nil {
			continue
		}
		typeName := kebabToCamelCase(avp.Name)
		a := []rune(typeName)
		a[0] = unicode.ToLower(a[0])
		varName := string(a)
		field := &GeneralField{
			varName:       varName,
			avpCode:       avp.Code,
			jsonFieldName: avp.Name,
			repeated:      r.Max != 1,
			required:      r.Required,
		}
		switch avp.Data.Type {

		case datatype.OctetStringType:
			// field.dataType = "bytes"
			field.dataType = "string"

		case datatype.UTF8StringType, datatype.DiameterIdentityType,
			datatype.AddressType, datatype.DiameterURIType, datatype.IPFilterRuleType:
			field.dataType = "string"

		case datatype.EnumeratedType:
			if len(avp.Data.Enum) == 0 {
				log.Panic("Enum with no values")
				continue
			}
			field.dataType = typeName + "Enum"
			enumField := processEnumField(field.dataType, avp.Data.Enum)
			checkConflictAndResolve(field.dataType, enumField)
		case datatype.GroupedType:
			field.dataType = typeName
			groupField := d.build(typeName, 50, &Node{appId: node.appId, rules: avp.Data.Rule, vendorId: node.vendorId})
			checkConflictAndResolve(field.dataType, groupField)
		case datatype.Unsigned32Type:
			field.dataType = "uint32"
			if !field.required {
				field.dataType = "google.protobuf.UInt32Value"
			}
		case datatype.Unsigned64Type:
			field.dataType = "uint64"
			if !field.required {
				field.dataType = "google.protobuf.UInt64Value"
			}
		case datatype.Integer32Type:
			field.dataType = "int32"
			if !field.required {
				field.dataType = "google.protobuf.Int32Value"
			}
		case datatype.Integer64Type:
			field.dataType = "int64"
			if !field.required {
				field.dataType = "google.protobuf.Int64Value"
			}
		case datatype.TimeType:
			field.dataType = "google.protobuf.Timestamp"
		default:
			field.dataType = avp.Data.TypeName
			log.Panicf("%s data type not supported yet", field.dataType)
		}
		composite.fields = append(composite.fields, field)
	}
	return composite
}

func (d *Dictionary) search(appId, vendorId uint32, code interface{}) (*dict.AVP, error) {
	var message = "+ Found AVP with VendorId [%s]"
	avp, err := d.P.FindAVPWithVendor(appId, code, vendorId)
	if err != nil {
		message = "- Failed to find AVP with VendorId [ %s ]"
		avp, err = d.P.FindAVP(appId, code)
		if err != nil {
			message = "-- Failed to find AVP without VendorId [ %s ]"
			avp, err = d.P.ScanAVP(code)
			if err != nil {
				message = "--- Failed to find AVP globally [ %s ]"
				// log.Fatalf("AVP [ %s ] not found !!!", code)
			}
		}
	}
	log.Printf(message, code)
	return avp, err
}

func processEnumField(name string, enums []*dict.Enum) CompositeField {
	composite := CompositeField{name: name, priority: 10, protoDataType: "enum"}
	if enums[0].Code != 0 {
		first := strings.Split(enums[0].Name, "_")
		second := strings.Split(enums[len(enums)-1].Name, "_")
		var field *EnumField
		if first[0] == second[0] {
			field = &EnumField{name: fmt.Sprintf("_%s_UNDEFINED", first[0]), code: 0}
		} else if first[len(first)-1] == second[len(second)-1] {
			field = &EnumField{name: fmt.Sprintf("_UNDEFINED_%s", first[len(first)-1]), code: 0}
		} else {
			field = &EnumField{name: fmt.Sprintf("_%s_UNDEFINED", name), code: 0}
		}
		composite.fields = append(composite.fields, field)
	}
	for _, enum := range enums {
		field := &EnumField{name: kebabToCamelCase(enum.Name), code: uint32(enum.Code)}
		composite.fields = append(composite.fields, field)
	}
	return composite
}

func checkConflictAndResolve(dataType string, compField CompositeField) {
	if parsedField, ok := parsedFields[dataType]; ok {
		log.Printf("Type %s already processed", dataType)
		if !reflect.DeepEqual(compField, parsedField) {
			log.Printf("*** Type %s has mismatching fields", dataType)
			if len(compField.fields) > len(parsedField.fields) {
				parsedFields[dataType] = compField
			} else if len(compField.fields) == len(parsedField.fields) {
				log.Fatalf("****** Type %s has deep mismatching fields. Needs manual intervention", dataType)
			}
		}
		return
	}
	parsedFields[dataType] = compField
}

func kebabToCamelCase(kebab string) (camelCase string) {
	isToUpper := false
	isFirstLetter := true
	for _, runeValue := range kebab {
		if isFirstLetter {
			isFirstLetter = false
			// special case where protoc fails to compile enums whose names starting with numeric char, mainly '3' as in 3GPP
			if string(runeValue) == "3" {
				camelCase += "T"
				continue
			}
		}
		if isToUpper {
			camelCase += strings.ToUpper(string(runeValue))
			isToUpper = false
		} else {
			if runeValue == '-' {
				isToUpper = true
			} else {
				camelCase += string(runeValue)
			}
		}
	}
	return
}

type CompositeField struct {
	priority      int
	name          string
	protoDataType string
	fields        []Field
}

type Field interface {
	GetCode() uint32
	SetIndex(int)
}

type EnumField struct {
	name string
	code uint32
}

func (f *EnumField) GetCode() uint32 {
	return f.code
}

func (f *EnumField) SetIndex(int) {
}

func (f *EnumField) String() string {
	// return fmt.Sprintf("\t%s = %d;", f.name, f.code)
	return fmt.Sprintf("\t\t%s = %d;", f.name, f.code)
}

type GeneralField struct {
	index         int
	dataType      string
	varName       string
	avpCode       uint32
	jsonFieldName string
	comment       string
	isAlternative bool
	repeated      bool
	required      bool
	nonnull       bool
}

func (f *GeneralField) GetCode() uint32 {
	return f.avpCode
}

func (f *GeneralField) SetIndex(i int) {
	f.index = i
}

func (f *GeneralField) String() string {
	s := "\t"
	nullExtension := ""
	if f.isAlternative {
		s += "// "
	}
	if f.repeated {
		s += "repeated "
	}
	if f.nonnull {
		nullExtension = ", (gogoproto.nullable) = false"
	}
	s += fmt.Sprintf("%s %s = %d [json_name = \"%s\"%s];", f.dataType, f.varName, f.index, f.jsonFieldName, nullExtension)
	if f.comment != "" {
		s += " // " + f.comment
	}
	return s
}
