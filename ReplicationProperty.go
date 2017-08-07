package main
import "github.com/google/gopacket"
import "github.com/davecgh/go-spew/spew"
import "errors"

type ReplicationProperty struct {
	Schema *PropertySchemaItem
	Value PropertyValue
	IsDefault bool
}


const (
	ROUND_JOINDATA	= iota
	ROUND_STRINGS	= iota
	ROUND_OTHER		= iota
	ROUND_UPDATE	= iota
)

func (schema *PropertySchemaItem) Decode(round int, thisBitstream *ExtendedReader, context *CommunicationContext, packet gopacket.Packet) (*ReplicationProperty, error) {
	var err error
	if !schema.Replicates {
		return nil, nil
	}
	isJoinData := round == 0

	isStringObject := true
	if schema.Type != "Object" &&
	schema.Type != "string" &&
	schema.Type != "ProtectedString" &&
	schema.Type != "BinaryString" &&
	schema.Type != "SystemAddress" &&
	schema.Type != "Content" {
		isStringObject = false
	}

	if round == ROUND_STRINGS && !isStringObject {
		return nil, nil
	}
	if round == ROUND_OTHER && isStringObject {
		return nil, nil
	}

	Property := &ReplicationProperty{Schema: schema}
	if schema.Type == "bool" {
		if err != nil {
			return Property, err
		}
		Property.Value, err = thisBitstream.ReadPBool()
		println(DebugInfo2(context, packet, isJoinData), "Read bool", schema.Name, Property.Value)
		if err != nil {
			return Property, err
		}
	} else {
		if round != ROUND_UPDATE {
			Property.IsDefault, err = thisBitstream.ReadBool()
			if err != nil {
				return Property, err
			}
			if Property.IsDefault {
				println(DebugInfo2(context, packet, isJoinData), "Read", schema.Name, "1 bit: default")
				return Property, nil
			}
		}
		switch schema.Type {
		case "string":
			Property.Value, err = thisBitstream.ReadPString(isJoinData, context)
			break
		case "ProtectedString":
			Property.Value, err = thisBitstream.ReadProtectedString(isJoinData, context)
			break
		case "BinaryString":
			Property.Value, err = thisBitstream.ReadBinaryString()
			break
		case "int":
			Property.Value, err = thisBitstream.ReadPInt()
			break
		case "float":
			Property.Value, err = thisBitstream.ReadPFloat()
			break
		case "double":
			Property.Value, err = thisBitstream.ReadPDouble()
			break
		case "Axes":
			Property.Value, err = thisBitstream.ReadAxes()
			break
		case "Faces":
			Property.Value, err = thisBitstream.ReadFaces()
			break
		case "BrickColor":
			Property.Value, err = thisBitstream.ReadBrickColor()
			break
		case "Object":
			Property.Value, err = thisBitstream.ReadObject(isJoinData, context)
			break
		case "UDim":
			Property.Value, err = thisBitstream.ReadUDim()
			break
		case "UDim2":
			Property.Value, err = thisBitstream.ReadUDim2()
			break
		case "Vector2":
			Property.Value, err = thisBitstream.ReadVector2()
			break
		case "Vector3":
			Property.Value, err = thisBitstream.ReadVector3()
			break
		case "Vector2uint16":
			Property.Value, err = thisBitstream.ReadVector2uint16()
			break
		case "Vector3uint16":
			Property.Value, err = thisBitstream.ReadVector3uint16()
			break
		case "Ray":
			Property.Value, err = thisBitstream.ReadRay()
			break
		case "Color3":
			Property.Value, err = thisBitstream.ReadColor3()
			break
		case "Color3uint8":
			Property.Value, err = thisBitstream.ReadColor3uint8()
			break
		case "CoordinateFrame":
			Property.Value, err = thisBitstream.ReadCFrame()
			break
		case "Content":
			Property.Value, err = thisBitstream.ReadContent()
			break
		default:
			if schema.IsEnum {
				Property.Value, err = thisBitstream.ReadEnumValue(schema.BitSize)
			} else {
				return Property, errors.New("property parser encountered unknown type")
			}
		}
		println(DebugInfo2(context, packet, isJoinData), "Read", schema.Name, spew.Sdump(Property.Value))
	}
	return Property, nil
}
