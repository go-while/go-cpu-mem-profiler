# go-cpu-mem-profiler
 a hopefully simple to use cpu and mem profiler for GOlang


```go
import (
	"github.com/go-while/go-cpu-mem-profiler"
)

var (
	Prof *prof.Profiler
)

func main() {
	Prof = prof.NewProf()

	// start a webserver
	go Prof.PprofWeb(":1234")

	// start a cpu profiler
	// _ is 'CPUfile' = the open file handle
	_, err := Prof.StartCPUProfile()
	if err != nil {
		os.Exit(1)
	}

	// stops the running cpu profiler
	Prof.StopCPUProfile()

	// starts a memory profiler for runtime
	// use waittime to delay the start
	Prof.MemProfile(runtime, waittime, true)
}
```


## Contributing

Contributions to this code are welcome.

If you have suggestions for improvements or find issues, please feel free to open an issue or submit a pull request.

## License

This code is provided under the MIT License. See the [LICENSE](LICENSE) file for details.
