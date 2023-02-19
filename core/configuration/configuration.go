package configuration

type LonginusSignatureConfiguration struct {
	Name              string `yaml:"name"`
	Signature         string `yaml:"signature"`
	InstructionOffset int    `yaml:"instruction_offset"`
	IsRelative        bool   `yaml:"is_relative"`
}

type LonginusExecutableConfiguration struct {
	Name       string                           `yaml:"name"`
	Signatures []LonginusSignatureConfiguration `yaml:"signatures"`
}

type LonginusConfiguration struct {
	Executables []LonginusExecutableConfiguration `yaml:"executables"`
}
