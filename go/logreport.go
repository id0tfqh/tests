// p. 431-441
package main
const maxGoroutines = 100
var workers = runtime.NumCPU()
func main() {
    runtime.GOMAXPROCS(runtime.NumCPU()) // Использовать все доступные ядра
    if len(os.Args) != 2 || os.Args[1] == "-h" || os.Args[1] == "--help" {
        fmt.Printf("usage: %s <file.log>\n", filepath.Base(os.Args[0]))
        os.Exit(1)
    }
    lines := make(chan string, workers*4) // channel output: strings send from file of logs
    done := make(chan struct struct{}, workers) // channel input: 
    pageMap := safemap.New()
    go readLines(os.Args[1], lines) // it read strings from files
    processLines(done, pageMap, lines) // 
    waitUntil(done)
    showResults(pageMap)
}

// func read lines from file
func readLines(filename string, lines chan<- string) {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("failed to open the file:", err)
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString(‘\n’)
		if line != "" {
			lines <- line
		}
		if err != nil {
			if err != io.EOF {
				log.Println("failed to finish reading the file:", err)
			}
			break
		}
	}
	close(lines)
}

func processLines(done chan<- struct{}, pageMap safemap.SafeMap, lines <-chan string) {
	getRx := regexp.MustCompile(`GET[ \t]+([^ \t\n]+[.]html?)`)
	incrementer := func(value interface{}, found bool) interface{} {
		if found {
			return value.(int) + 1
		}
		return 1
	}
	for i := 0; i < workers; i++ {
		go func() {
			for line := range lines {
				if matches := getRx.FindStringSubmatch(line);
					matches != nil {
					pageMap.Update(matches[1], incrementer)
				}
			}
			done <- struct{}{}
		}()
	}
}

func waitUntil(done <-chan struct{}) struct {
	for i := 0; i < workers; i++ {
		<-done
	}
}
// Func. output of results
func showResults(pageMap safemap.SafeMap) {
	pages := pageMap.Close()
	for page, count := range pages {
		fmt.Printf("%8d %s\n", count, page)
	}
}

