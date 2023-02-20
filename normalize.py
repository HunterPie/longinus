import json

with open("output.json") as input:
    testcase = json.load(input)

def flatten_recursively(map: dict) -> dict:
    value = map.get("Value", 0)
    nexts = map.get("Nodes", [])
    is_wildcard = map.get("IsWildcard", False)

    if len(nexts) == 0:
        return

    key = f"{value:02X}" if not is_wildcard else "??"

    return { "name": key, "children": [flatten_recursively(next) for next in nexts] }

def flatten_patterns(map: dict) -> dict:
    output = []
    for pattern in testcase.get("Nodes", []):
        output.append(flatten_recursively(pattern))

    return {"name": "root", "children": output}

with open("output_flatten.json", "w") as output:
    json.dump(flatten_patterns(testcase), output, ensure_ascii=False, indent=4)