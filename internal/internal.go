package internal

// LoadMapping loads function name mapping from a YAML file.
func LoadMapping(path string) (Mapping[string, string], error) {
	return NewYAMLMappingLoader[string, string]().Load(path)
}

// ProcessDir processes all .gno files in the given directory with the provided mapping.
func ProcessDir(dir string, mapping Mapping[string, string]) error {
	return process(dir, mapping)
}

func process(dir string, mapping Mapping[string, string]) error {
	transformer := newStdFunctionTransformer(mapping)
	processor := newGnoFileProcessor(mapping, transformer)
	dirProcessor := newDirectoryProcessor(processor)
	return dirProcessor.ProcessDir(dir)
}
