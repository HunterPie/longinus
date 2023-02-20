## Longinus

Longinus is a byte signature tree generator and scanner for scanning byte patterns in statically compiled executables. In Brazil, Saint Longinus is known for its power of finding missing objects.

> **Attention** This application is still under development and cannot be used in its current state 

### Usage

The point of this application is to be agnostic to whatever executable you are scanning, so all you need to do is define a configuration for Longinus and run the command:

```shell
./Longinus -executable <PATH> -config <PATH>
```

#### Configuration

The configuration is a `yaml` file containing the executable name and a list of signatures. Here's the schema:

```yaml
executables:
  - name: EXECUTABLE_NAME.exe
    signatures:
      - name: "PATTERN_NAME"
        signature: "48 8B ?? ?? ?? ?? 00 ??"
        instruction_offset: 3
        is_relative: true # This field is optional, true by default
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

