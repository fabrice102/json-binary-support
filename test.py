import json

orig = {"Data": b"\x00\x01\x80\x81\xff"}
origJSON = r"""{
  "Data": "\u0000\u0001\u0080\u0081\u00ff" 
}"""

print("Original object:")
print(orig)
print()

print("Data = {}".format(orig["Data"]))
print()

try:
    print("Serialize:")
    print(json.dumps(orig))
except Exception as e:
    print(f"  Exception: {e}")
print()

print("Original JSON:")
print(origJSON)
print()

print("Deserialize:")
print(json.loads(origJSON))
print()

print("----------")

print("Same as above after encoding the string as raw unicode escape")
origRawUnicodeEscape = {"Data": orig["Data"].decode("raw-unicode-escape")}

print("Original object:")
print(origRawUnicodeEscape)
print()

print("Data = {}".format(origRawUnicodeEscape["Data"]))
print()

print("Serialize:")
origEnc = json.dumps(origRawUnicodeEscape)
print(origEnc)
print()

print("Deserialize:")
origEncDec = json.loads(origEnc)
print(origEncDec)
print()

print("Converting back to bytes")
origEncDec["Data"] = origEncDec["Data"].encode("raw-unicode-escape")
print(origEncDec)
print()

print("Checking that raw-unicode-escape work for all characters 0x00 to 0xff")
b = bytes([i for i in range(256)])
b2 = json.loads(json.dumps(b.decode("raw-unicode-escape"))).encode("raw-unicode-escape")
print(b == b2)

