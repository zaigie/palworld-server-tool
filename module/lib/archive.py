import io
import os
import struct
import uuid
from typing import Any, Callable, Optional, Union


class UUID:
    """Wrapper around uuid.UUID to delay evaluation of UUIDs until necessary"""

    __slots__ = ("raw_bytes", "parsed_uuid")
    raw_bytes: bytes
    parsed_uuid: uuid.UUID

    def __init__(self, raw_bytes: bytes) -> None:
        self.raw_bytes = raw_bytes
        self.parsed_uuid = None

    def from_str(s: str) -> "UUID":
        b = uuid.UUID(s).bytes
        return UUID(
            bytes(
                [
                    b[0x3],
                    b[0x2],
                    b[0x1],
                    b[0x0],
                    b[0x7],
                    b[0x6],
                    b[0x5],
                    b[0x4],
                    b[0xB],
                    b[0xA],
                    b[0x9],
                    b[0x8],
                    b[0xF],
                    b[0xE],
                    b[0xD],
                    b[0xC],
                ]
            )
        )

    def __str__(self) -> str:
        if not self.parsed_uuid:
            b = self.raw_bytes
            uuid_int = (
                b[0xC]
                + (b[0xD] << 8)
                + (b[0xE] << 16)
                + (b[0xF] << 24)
                + (b[0x8] << 32)
                + (b[0x9] << 40)
                + (b[0xA] << 48)
                + (b[0xB] << 56)
                + (b[0x4] << 64)
                + (b[0x5] << 72)
                + (b[0x6] << 80)
                + (b[0x7] << 88)
                + (b[0x0] << 96)
                + (b[0x1] << 104)
                + (b[0x2] << 112)
                + (b[0x3] << 120)
            )
            self.parsed_uuid = uuid.UUID(int=uuid_int)
        return str(self.parsed_uuid)

    def __eq__(self, __value: object) -> bool:
        return str(self) == str(__value)


def instance_id_reader(reader: "FArchiveReader") -> dict[str, UUID]:
    return {
        "guid": reader.guid(),
        "instance_id": reader.guid(),
    }


def uuid_reader(reader: "FArchiveReader") -> UUID:
    b = reader.read(16)
    if len(b) != 16:
        raise Exception("could not read 16 bytes for uuid")
    return UUID(b)


class FArchiveReader:
    data: io.BytesIO
    size: int
    type_hints: dict[str, str]
    custom_properties: dict[str, tuple[Callable, Callable]]
    debug: bool

    def __init__(
        self,
        data,
        type_hints: dict[str, str] = {},
        custom_properties: dict[str, tuple[Callable, Callable]] = {},
        debug: bool = os.environ.get("DEBUG", "0") == "1",
    ):
        self.data = io.BytesIO(data)
        self.size = len(data)
        self.type_hints = type_hints
        self.custom_properties = custom_properties
        self.debug = debug

    def __enter__(self):
        self.data.seek(0)
        return self

    def __exit__(self, type, value, traceback):
        self.data.close()

    def get_type_or(self, path: str, default: str):
        if path in self.type_hints:
            return self.type_hints[path]
        else:
            print(f"Struct type for {path} not found, assuming {default}")
            return default

    def eof(self) -> bool:
        return self.data.tell() >= self.size

    def read(self, size: int) -> bytes:
        return self.data.read(size)

    def read_to_end(self) -> bytes:
        return self.data.read(self.size - self.data.tell())

    def bool(self) -> bool:
        return self.byte() > 0

    def fstring(self) -> str:
        # in the hot loop, avoid function calls
        reader = self.data
        (size,) = self.unpack_i32(reader.read(4))

        if size == 0:
            return ""

        data: bytes
        encoding: str
        if size < 0:
            size = -size
            data = reader.read(size * 2)[:-2]
            encoding = "utf-16-le"
        else:
            data = reader.read(size)[:-1]
            encoding = "ascii"

        try:
            return data.decode(encoding)
        except Exception as e:
            try:
                escaped = data.decode(encoding, errors="surrogatepass")
                print(
                    f"Error decoding {encoding} string of length {size}, data loss may occur! {bytes(data)}"
                )
                return escaped
            except Exception as e:
                raise Exception(
                    f"Error decoding {encoding} string of length {size}: {bytes(data)}"
                ) from e

    unpack_i16 = struct.Struct("h").unpack

    def i16(self) -> int:
        return self.unpack_i16(self.data.read(2))[0]

    unpack_u16 = struct.Struct("H").unpack

    def u16(self) -> int:
        return self.unpack_u16(self.data.read(2))[0]

    unpack_i32 = struct.Struct("i").unpack

    def i32(self) -> int:
        return self.unpack_i32(self.data.read(4))[0]

    unpack_u32 = struct.Struct("I").unpack

    def u32(self) -> int:
        return self.unpack_u32(self.data.read(4))[0]

    unpack_i64 = struct.Struct("q").unpack

    def i64(self) -> int:
        return self.unpack_i64(self.data.read(8))[0]

    unpack_u64 = struct.Struct("Q").unpack

    def u64(self) -> int:
        return self.unpack_u64(self.data.read(8))[0]

    unpack_float = struct.Struct("f").unpack

    def float(self) -> float:
        return self.unpack_float(self.data.read(4))[0]

    unpack_double = struct.Struct("d").unpack

    def double(self) -> float:
        return self.unpack_double(self.data.read(8))[0]

    unpack_byte = struct.Struct("B").unpack

    def byte(self) -> int:
        return self.unpack_byte(self.data.read(1))[0]

    def byte_list(self, size: int) -> list[int]:
        return struct.unpack(str(size) + "B", self.data.read(size))

    def skip(self, size: int) -> None:
        self.data.read(size)

    def guid(self) -> UUID:
        # in the hot loop, avoid function calls
        return UUID(self.data.read(16))

    def optional_guid(self) -> Optional[UUID]:
        # in the hot loop, avoid function calls
        if self.data.read(1)[0]:
            return UUID(self.data.read(16))
        return None

    def tarray(
        self, type_reader: Callable[["FArchiveReader"], dict[str, Any]]
    ) -> list[dict[str, Any]]:
        count = self.u32()
        array = []
        for _ in range(count):
            array.append(type_reader(self))
        return array

    def properties_until_end(self, path: str = "") -> dict[str, Any]:
        properties = {}
        while True:
            name = self.fstring()
            if name == "None":
                break
            type_name = self.fstring()
            size = self.u64()
            properties[name] = self.property(type_name, size, f"{path}.{name}")
        return properties

    def property(
        self, type_name: str, size: int, path: str, nested_caller_path: str = ""
    ) -> dict[str, Any]:
        value = {}
        if path in self.custom_properties and (
            path is not nested_caller_path or nested_caller_path == ""
        ):
            value = self.custom_properties[path][0](self, type_name, size, path)
            value["custom_type"] = path
        elif type_name == "StructProperty":
            value = self.struct(path)
        elif type_name == "IntProperty":
            value = {
                "id": self.optional_guid(),
                "value": self.i32(),
            }
        elif type_name == "Int64Property":
            value = {
                "id": self.optional_guid(),
                "value": self.i64(),
            }
        elif type_name == "FixedPoint64Property":
            value = {
                "id": self.optional_guid(),
                "value": self.i32(),
            }
        elif type_name == "FloatProperty":
            value = {
                "id": self.optional_guid(),
                "value": self.float(),
            }
        elif type_name == "StrProperty":
            value = {
                "id": self.optional_guid(),
                "value": self.fstring(),
            }
        elif type_name == "NameProperty":
            value = {
                "id": self.optional_guid(),
                "value": self.fstring(),
            }
        elif type_name == "EnumProperty":
            enum_type = self.fstring()
            _id = self.optional_guid()
            enum_value = self.fstring()
            value = {
                "id": _id,
                "value": {
                    "type": enum_type,
                    "value": enum_value,
                },
            }
        elif type_name == "BoolProperty":
            value = {
                "value": self.bool(),
                "id": self.optional_guid(),
            }
        elif type_name == "ArrayProperty":
            array_type = self.fstring()
            value = {
                "array_type": array_type,
                "id": self.optional_guid(),
                "value": self.array_property(array_type, size - 4, path),
            }
        elif type_name == "MapProperty":
            key_type = self.fstring()
            value_type = self.fstring()
            _id = self.optional_guid()
            self.u32()
            count = self.u32()
            values = {}
            key_path = path + ".Key"
            if key_type == "StructProperty":
                key_struct_type = self.get_type_or(key_path, "Guid")
            else:
                key_struct_type = None
            value_path = path + ".Value"
            if value_type == "StructProperty":
                value_struct_type = self.get_type_or(value_path, "StructProperty")
            else:
                value_struct_type = None
            values = []
            for _ in range(count):
                key = self.prop_value(key_type, key_struct_type, key_path)
                value = self.prop_value(value_type, value_struct_type, value_path)
                values.append(
                    {
                        "key": key,
                        "value": value,
                    }
                )
            value = {
                "key_type": key_type,
                "value_type": value_type,
                "key_struct_type": key_struct_type,
                "value_struct_type": value_struct_type,
                "id": _id,
                "value": values,
            }
        else:
            raise Exception(f"Unknown type: {type_name} ({path})")
        value["type"] = type_name
        return value

    def prop_value(self, type_name: str, struct_type_name: str, path: str):
        if type_name == "StructProperty":
            return self.struct_value(struct_type_name, path)
        elif type_name == "EnumProperty":
            return self.fstring()
        elif type_name == "NameProperty":
            return self.fstring()
        elif type_name == "IntProperty":
            return self.i32()
        elif type_name == "BoolProperty":
            return self.bool()
        else:
            raise Exception(f"Unknown property value type: {type_name} ({path})")

    def struct(self, path: str) -> dict[str, Any]:
        struct_type = self.fstring()
        struct_id = self.guid()
        _id = self.optional_guid()
        value = self.struct_value(struct_type, path)
        return {
            "struct_type": struct_type,
            "struct_id": struct_id,
            "id": _id,
            "value": value,
        }

    def struct_value(self, struct_type: str, path: str = ""):
        if struct_type == "Vector":
            return {
                "x": self.double(),
                "y": self.double(),
                "z": self.double(),
            }
        elif struct_type == "DateTime":
            return self.u64()
        elif struct_type == "Guid":
            return self.guid()
        elif struct_type == "Quat":
            return {
                "x": self.double(),
                "y": self.double(),
                "z": self.double(),
                "w": self.double(),
            }
        elif struct_type == "LinearColor":
            return {
                "r": self.float(),
                "g": self.float(),
                "b": self.float(),
                "a": self.float(),
            }
        else:
            if self.debug:
                print(f"Assuming struct type: {struct_type} ({path})")
            return self.properties_until_end(path)

    def array_property(self, array_type: str, size: int, path: str):
        count = self.u32()
        value = {}
        if array_type == "StructProperty":
            prop_name = self.fstring()
            prop_type = self.fstring()
            self.u64()
            type_name = self.fstring()
            _id = self.guid()
            self.skip(1)
            prop_values = []
            for _ in range(count):
                prop_values.append(self.struct_value(type_name, f"{path}.{prop_name}"))
            value = {
                "prop_name": prop_name,
                "prop_type": prop_type,
                "values": prop_values,
                "type_name": type_name,
                "id": _id,
            }
        else:
            value = {
                "values": self.array_value(array_type, count, size, path),
            }
        return value

    def array_value(self, array_type: str, count: int, size: int, path: str):
        values = []
        decode_func: Callable
        if array_type == "EnumProperty":
            decode_func = self.fstring
        elif array_type == "NameProperty":
            decode_func = self.fstring
        elif array_type == "Guid":
            decode_func = self.guid
        elif array_type == "ByteProperty":
            if size == count:
                # Special case this and read faster in one go
                return self.byte_list(count)
            else:
                raise Exception("Labelled ByteProperty not implemented")
        else:
            raise Exception(f"Unknown array type: {array_type} ({path})")
        for _ in range(count):
            values.append(decode_func())

        return values

    def compressed_short_rotator(self) -> tuple[float, float, float]:
        short_pitch = self.u16() if self.bool() else 0
        short_yaw = self.u16() if self.bool() else 0
        short_roll = self.u16() if self.bool() else 0
        pitch = short_pitch * (360.0 / 65536.0)
        yaw = short_yaw * (360.0 / 65536.0)
        roll = short_roll * (360.0 / 65536.0)
        return [pitch, yaw, roll]

    def serializeint(self, component_bit_count: int) -> int:
        b = bytearray(self.read((component_bit_count + 7) // 8))
        if (component_bit_count % 8) != 0:
            b[-1] &= (1 << (component_bit_count % 8)) - 1
        value = int.from_bytes(b, "little")
        return value

    def packed_vector(self, scale_factor: int) -> tuple[float, float, float]:
        component_bit_count_and_extra_info = self.u32()
        component_bit_count = component_bit_count_and_extra_info & 63
        extra_info = component_bit_count_and_extra_info >> 6
        if component_bit_count > 0:
            x = self.serializeint(component_bit_count)
            y = self.serializeint(component_bit_count)
            z = self.serializeint(component_bit_count)
            sign_bit = 1 << (component_bit_count - 1)
            x = (x & (sign_bit - 1)) - (x & sign_bit)
            y = (y & (sign_bit - 1)) - (y & sign_bit)
            z = (z & (sign_bit - 1)) - (z & sign_bit)

            if extra_info:
                x /= scale_factor
                y /= scale_factor
                z /= scale_factor
            return (x, y, z)
        else:
            received_scaler_type_size = 8 if extra_info else 4
            if received_scaler_type_size == 8:
                x = self.double()
                y = self.double()
                z = self.double()
                return (x, y, z)
            else:
                x = self.float()
                y = self.float()
                z = self.float()
                return (x, y, z)

    def ftransform(self) -> dict[str, dict[str, float]]:
        return {
            "rotation": {
                "x": self.double(),
                "y": self.double(),
                "z": self.double(),
                "w": self.double(),
            },
            "translation": {
                "x": self.double(),
                "y": self.double(),
                "z": self.double(),
            },
            "scale3d": {
                "x": self.double(),
                "y": self.double(),
                "z": self.double(),
            },
        }


def uuid_writer(writer, s: Union[str, uuid.UUID, UUID]):
    if isinstance(s, str):
        s = uuid.UUID(s)
    if isinstance(s, uuid.UUID):
        b = s.bytes
        ub = bytes(
            [
                b[0x3],
                b[0x2],
                b[0x1],
                b[0x0],
                b[0x7],
                b[0x6],
                b[0x5],
                b[0x4],
                b[0xB],
                b[0xA],
                b[0x9],
                b[0x8],
                b[0xF],
                b[0xE],
                b[0xD],
                b[0xC],
            ]
        )
    elif isinstance(s, UUID):
        ub = s.raw_bytes
    writer.write(ub)


def instance_id_writer(writer, d):
    uuid_writer(writer, d["guid"])
    uuid_writer(writer, d["instance_id"])


class FArchiveWriter:
    data: io.BytesIO
    size: int
    custom_properties: dict[str, tuple[Callable, Callable]]
    debug: bool

    def __init__(
        self,
        custom_properties: dict[str, tuple[Callable, Callable]] = {},
        debug: bool = os.environ.get("DEBUG", "0") == "1",
    ):
        self.data = io.BytesIO()
        self.custom_properties = custom_properties
        self.debug = debug

    def __enter__(self):
        self.data.seek(0)
        return self

    def __exit__(self, type, value, traceback):
        self.data.close()

    def copy(self) -> "FArchiveWriter":
        return FArchiveWriter(self.custom_properties)

    def bytes(self) -> bytes:
        pos = self.data.tell()
        self.data.seek(0)
        b = self.data.read()
        self.data.seek(pos)
        return b

    def write(self, data: bytes):
        self.data.write(data)

    def bool(self, bool: bool):
        self.data.write(struct.pack("?", bool))

    def fstring(self, string: str) -> int:
        start = self.data.tell()
        if string == "":
            self.i32(0)
        elif string.isascii():
            str_bytes = string.encode("ascii")
            self.i32(len(str_bytes) + 1)
            self.data.write(str_bytes)
            self.data.write(b"\x00")
        else:
            str_bytes = string.encode("utf-16-le", errors="surrogatepass")
            assert len(str_bytes) % 2 == 0
            self.i32(-((len(str_bytes) // 2) + 1))
            self.data.write(str_bytes)
            self.data.write(b"\x00\x00")
        return self.data.tell() - start

    def i16(self, i: int):
        self.data.write(struct.pack("h", i))

    def u16(self, i: int):
        self.data.write(struct.pack("H", i))

    def i32(self, i: int):
        self.data.write(struct.pack("i", i))

    def u32(self, i: int):
        self.data.write(struct.pack("I", i))

    def i64(self, i: int):
        self.data.write(struct.pack("q", i))

    def u64(self, i: int):
        self.data.write(struct.pack("Q", i))

    def float(self, i: float):
        self.data.write(struct.pack("f", i))

    def double(self, i: float):
        self.data.write(struct.pack("d", i))

    def byte(self, b: int):
        self.data.write(bytes([b]))

    def u(self, b: int):
        self.data.write(struct.pack("B", b))

    def guid(self, u: Union[str, uuid.UUID, UUID]):
        uuid_writer(self, u)

    def optional_guid(self, u: Optional[Union[str, uuid.UUID, UUID]]):
        if u is None:
            self.bool(False)
        else:
            self.bool(True)
            uuid_writer(self, u)

    def tarray(
        self, type_writer: Callable[["FArchiveWriter", dict[str, Any]], None], array
    ):
        self.u32(len(array))
        for i in range(len(array)):
            type_writer(self, array[i])

    def properties(self, properties: dict[str, Any]):
        for key in properties:
            self.fstring(key)
            self.property(properties[key])
        self.fstring("None")

    def property(self, property: dict[str, Any]):
        # write type_name
        self.fstring(property["type"])
        nested_writer = self.copy()
        size: int
        property_type = property["type"]
        size = nested_writer.property_inner(property_type, property)
        buf = nested_writer.bytes()
        # write size
        self.u64(size)
        self.write(buf)

    def property_inner(self, property_type: str, property: dict[str, Any]) -> int:
        if "custom_type" in property:
            if property["custom_type"] in self.custom_properties:
                size = self.custom_properties[property["custom_type"]][1](
                    self, property_type, property
                )
            else:
                raise Exception(
                    f"Unknown custom property type: {property['custom_type']}"
                )
        elif property_type == "StructProperty":
            size = self.struct(property)
        elif property_type == "IntProperty":
            self.optional_guid(property.get("id", None))
            self.i32(property["value"])
            size = 4
        elif property_type == "Int64Property":
            self.optional_guid(property.get("id", None))
            self.i64(property["value"])
            size = 8
        elif property_type == "FixedPoint64Property":
            self.optional_guid(property.get("id", None))
            self.i32(property["value"])
            size = 4
        elif property_type == "FloatProperty":
            self.optional_guid(property.get("id", None))
            self.float(property["value"])
            size = 4
        elif property_type == "StrProperty":
            self.optional_guid(property.get("id", None))
            size = self.fstring(property["value"])
        elif property_type == "NameProperty":
            self.optional_guid(property.get("id", None))
            size = self.fstring(property["value"])
        elif property_type == "EnumProperty":
            self.fstring(property["value"]["type"])
            self.optional_guid(property.get("id", None))
            size = self.fstring(property["value"]["value"])
        elif property_type == "BoolProperty":
            self.bool(property["value"])
            self.optional_guid(property.get("id", None))
            size = 0
        elif property_type == "ArrayProperty":
            self.fstring(property["array_type"])
            self.optional_guid(property.get("id", None))
            array_writer = self.copy()
            array_writer.array_property(property["array_type"], property["value"])
            array_buf = array_writer.bytes()
            size = len(array_buf)
            self.write(array_buf)
        elif property_type == "MapProperty":
            self.fstring(property["key_type"])
            self.fstring(property["value_type"])
            self.optional_guid(property.get("id", None))
            map_writer = self.copy()
            map_writer.u32(0)
            map_writer.u32(len(property["value"]))
            for entry in property["value"]:
                map_writer.prop_value(
                    property["key_type"], property["key_struct_type"], entry["key"]
                )
                map_writer.prop_value(
                    property["value_type"],
                    property["value_struct_type"],
                    entry["value"],
                )
            map_buf = map_writer.bytes()
            size = len(map_buf)
            self.write(map_buf)
        else:
            raise Exception(f"Unknown property type: {property_type}")
        return size

    def struct(self, property: dict[str, Any]) -> int:
        self.fstring(property["struct_type"])
        self.guid(property["struct_id"])
        self.optional_guid(property.get("id", None))
        start = self.data.tell()
        self.struct_value(property["struct_type"], property["value"])
        return self.data.tell() - start

    def struct_value(self, struct_type: str, value):
        if struct_type == "Vector":
            self.double(value["x"])
            self.double(value["y"])
            self.double(value["z"])
        elif struct_type == "DateTime":
            self.u64(value)
        elif struct_type == "Guid":
            self.guid(value)
        elif struct_type == "Quat":
            self.double(value["x"])
            self.double(value["y"])
            self.double(value["z"])
            self.double(value["w"])
        elif struct_type == "LinearColor":
            self.float(value["r"])
            self.float(value["g"])
            self.float(value["b"])
            self.float(value["a"])
        else:
            if self.debug:
                print(f"Assuming struct type: {struct_type}")
            return self.properties(value)

    def prop_value(self, type_name: str, struct_type_name: str, value):
        if type_name == "StructProperty":
            self.struct_value(struct_type_name, value)
        elif type_name == "EnumProperty":
            self.fstring(value)
        elif type_name == "NameProperty":
            self.fstring(value)
        elif type_name == "IntProperty":
            self.i32(value)
        elif type_name == "BoolProperty":
            self.bool(value)
        else:
            raise Exception(f"Unknown property value type: {type_name}")

    def array_property(self, array_type: str, value: dict[str, Any]):
        count = len(value["values"])
        self.u32(count)
        if array_type == "StructProperty":
            self.fstring(value["prop_name"])
            self.fstring(value["prop_type"])
            nested_writer = self.copy()
            for i in range(count):
                nested_writer.struct_value(value["type_name"], value["values"][i])
            data_buf = nested_writer.bytes()
            self.u64(len(data_buf))
            self.fstring(value["type_name"])
            self.guid(value["id"])
            self.u(0)
            self.write(data_buf)
        else:
            self.array_value(array_type, count, value["values"])

    def array_value(self, array_type: str, count: int, values: list[Any]):
        for i in range(count):
            if array_type == "IntProperty":
                self.i32(values[i])
            elif array_type == "Int64Property":
                self.i64(values[i])
            elif array_type == "FloatProperty":
                self.float(values[i])
            elif array_type == "StrProperty":
                self.fstring(values[i])
            elif array_type == "NameProperty":
                self.fstring(values[i])
            elif array_type == "EnumProperty":
                self.fstring(values[i])
            elif array_type == "BoolProperty":
                self.bool(values[i])
            elif array_type == "ByteProperty":
                self.byte(values[i])
            else:
                raise Exception(f"Unknown array type: {array_type}")

    def compressed_short_rotator(self, pitch: float, yaw: float, roll: float):
        short_pitch = round(pitch * (65536.0 / 360.0)) & 0xFFFF
        short_yaw = round(yaw * (65536.0 / 360.0)) & 0xFFFF
        short_roll = round(roll * (65536.0 / 360.0)) & 0xFFFF
        if short_pitch != 0:
            self.bool(True)
            self.u16(short_pitch)
        else:
            self.bool(False)
        if short_yaw != 0:
            self.bool(True)
            self.u16(short_yaw)
        else:
            self.bool(False)
        if short_roll != 0:
            self.bool(True)
            self.u16(short_roll)
        else:
            self.bool(False)

    @staticmethod
    def unreal_round_float_to_int(value: float) -> int:
        return int(value)

    @staticmethod
    def unreal_get_bits_needed(value: int) -> int:
        massaged_value = value ^ (value >> 63)
        return 65 - FArchiveWriter.count_leading_zeroes(massaged_value)

    @staticmethod
    def count_leading_zeroes(value: int) -> int:
        return 67 - len(bin(-value)) & ~value >> 64

    def serializeint(self, component_bit_count: int, value: int):
        self.write(
            int.to_bytes(value, (component_bit_count + 7) // 8, "little", signed=True)
        )

    def packed_vector(self, scale_factor: int, x: float, y: float, z: float):
        max_exponent_for_scaling = 52
        max_value_to_scale = 1 << max_exponent_for_scaling
        max_exponent_after_scaling = 62
        max_scaled_value = 1 << max_exponent_after_scaling
        scaled_x = x * scale_factor
        scaled_y = y * scale_factor
        scaled_z = z * scale_factor
        if max(abs(scaled_x), abs(scaled_y), abs(scaled_z)) < max_scaled_value:
            use_scaled_value = min(abs(x), abs(y), abs(z)) < max_value_to_scale
            if use_scaled_value:
                x = self.unreal_round_float_to_int(scaled_x)
                y = self.unreal_round_float_to_int(scaled_y)
                z = self.unreal_round_float_to_int(scaled_z)
            else:
                x = self.unreal_round_float_to_int(x)
                y = self.unreal_round_float_to_int(y)
                z = self.unreal_round_float_to_int(z)

            component_bit_count = max(
                self.unreal_get_bits_needed(x),
                self.unreal_get_bits_needed(y),
                self.unreal_get_bits_needed(z),
            )
            component_bit_count_and_scale_info = (
                1 << 6 if use_scaled_value else 0
            ) | component_bit_count
            self.u32(component_bit_count_and_scale_info)
            self.serializeint(component_bit_count, x)
            self.serializeint(component_bit_count, y)
            self.serializeint(component_bit_count, z)
        else:
            component_bit_count = 0
            component_bit_count_and_scale_info = (1 << 6) | component_bit_count
            self.u32(component_bit_count_and_scale_info)
            self.double(x)
            self.double(y)
            self.double(z)

    def ftransform(self, value: dict[str, dict[str, float]]):
        self.double(value["rotation"]["x"])
        self.double(value["rotation"]["y"])
        self.double(value["rotation"]["z"])
        self.double(value["rotation"]["w"])
        self.double(value["translation"]["x"])
        self.double(value["translation"]["y"])
        self.double(value["translation"]["z"])
        self.double(value["scale3d"]["x"])
        self.double(value["scale3d"]["y"])
        self.double(value["scale3d"]["z"])
