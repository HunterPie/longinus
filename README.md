## longinus

Longinus is a byte signature tree generator and scanner for finding byte patterns in binary files. In Brazil, Saint Longinus is known for its power of finding missing objects.

### Usage

The point of this application is to be agnostic to whatever executable you are scanning, so all you need to do is define a configuration for Longinus and run the command:

```shell
./longinus -executable <PATH> -config <PATH>
```

#### Configuration

The configuration is a `yaml` file containing the executable name and a list of signatures with its properties, the available properties are:

| Name               | Description                                                                                                                               |
|--------------------|-------------------------------------------------------------------------------------------------------------------------------------------|
| name               | Name of the signature                                                                                                                     |
| signature          | Byte array of the signature, it supports wildcards as `??`                                                                                |
| instruction_offset | The offset that will be added to the address where the signature was found                                                                |
|  is_relative       | If true, the address where the signature was found will be used. If the value is false, the value will be `address + *(address + offset)` |

```yaml
executables:
  - name: executable_name.exe
    signatures:
      - name: "PATTERN_NAME"                 
        signature: "48 8B ?? ?? ?? ?? 00 ??" 
        instruction_offset: 3                 
        is_relative: true                     
```

You can also find an example under the `./configuration/default.yaml` folder.

### Details

The signatures will be converted to a linked list and merged into a tree, which means the following signatures:

- `48 8B 05 ?? 02 00`
- `48 8B 15 ?? ??`
- `40 53 48 83`

Will be merged into this:

```
         ┌──────┐
         │ root │
         └──┬───┘
      ┌─────┴─────┐
    ┌─┴──┐      ┌─┴──┐
    │ 48 │      │ 40 │
    └─┬──┘      └─┬──┘
      │           │
    ┌─┴──┐      ┌─┴──┐
    │ 8B │      │ 53 │
    └─┬──┘      └─┬──┘
  ┌───┴───┐       │
┌─┴──┐  ┌─┴──┐  ┌─┴──┐
│ 05 │  │ 15 │  │ 48 │
└─┬──┘  └─┬──┘  └────┘
  │       │
┌─┴──┐  ┌─┴──┐
│ ?? │  │ ?? │
└─┬──┘  └────┘
  │
┌─┴──┐
│ 02 │
└────┘
```

