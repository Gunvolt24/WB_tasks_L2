package helpers

import "os"

// WithRedirForBuiltin - временно подменяет os.Stdin/os.Stdout для builtin и восстанавливает их.
func WithRedirForBuiltin(inFile, outFile string, appendMode bool, fn func()) error {
	origIn, origOut := os.Stdin, os.Stdout
	var in, out *os.File
	var err error

	// stdin: <
	if inFile != "" {
		in, err = os.Open(inFile)
		if err != nil {
			return err
		}
		os.Stdin = in
	}

	// stdout: >, >>
	if outFile != "" {
		flags := os.O_CREATE | os.O_WRONLY
		if appendMode {
			flags |= os.O_APPEND
		} else {
			flags |= os.O_TRUNC
		}
		out, err = os.OpenFile(outFile, flags, 0o644)
		if err != nil {
			if in != nil {
				_ = in.Close()
				os.Stdin = origIn
			}
			return err
		}
		os.Stdout = out
	}

	// вызываем builtin
	fn()

	// восстанавливаем дескрипторы
	if out != nil {
		_ = out.Close()
		os.Stdout = origOut
	}
	if in != nil {
		_ = in.Close()
		os.Stdin = origIn
	}
	return nil
}
